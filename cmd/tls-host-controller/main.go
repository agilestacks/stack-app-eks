package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/cloudflare/certinel"
	"github.com/cloudflare/certinel/fswatcher"
	whhttp "github.com/slok/kubewebhook/pkg/http"
	"github.com/slok/kubewebhook/pkg/log"
	"github.com/slok/kubewebhook/pkg/webhook/mutating"
	v1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// sorted string slice impl
type byLength []string

func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

type nothing struct{}

const cnLimit = 63

// https://cert-manager.io/docs/usage/ingress/
var certManagerAnnotations = []string{"kubernetes.io/tls-acme", "cert-manager.io/issuer", "cert-manager.io/cluster-issuer"}

func main() {
	logger := &log.Std{Debug: true}

	var defaultCN string
	flag.StringVar(&defaultCN, "default-cn", "",
		"A comma separated list of CNs that will be considered to create certificate CN shorter than 64 bytes")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr,
			`Usage:
	tls-host-controller [-default-cn app.cluster.account.superhub.io]

	If no -default-cn is set then the shortest host rule well be chomped in front to create CN < 64 bytes long

Flags:
`)
		flag.PrintDefaults()
	}
	flag.Parse()

	cns, err := parseCN(logger, defaultCN)
	if err != nil {
		panic(err)
	}

	mt := mutating.MutatorFunc(func(_ context.Context, obj metav1.Object) (bool, error) {
		ingress := obj.(*v1beta1.Ingress)
		var name string
		if len(ingress.ObjectMeta.Name) == 0 && len(ingress.ObjectMeta.GenerateName) > 0 {
			name = ingress.ObjectMeta.GenerateName
		} else {
			name = ingress.ObjectMeta.Name
		}

		// cert-manager installed ingress
		logger.Debugf("checking ingress %s", name)
		if strings.HasPrefix(name, "cm-acme-http-solver") {
			logger.Debugf("skipping cert-manager installed ingress %s", name)
			return false, nil
		}

		spec := &ingress.Spec

		// don't interfere with explicit TLS spec
		if spec.TLS != nil {
			logger.Debugf("skipping %s as it has TLS block configured", name)
			return false, nil
		}

		rulesHosts := make(map[string]nothing)
		for _, r := range spec.Rules {
			if len(r.Host) > 0 {
				rulesHosts[r.Host] = nothing{}
			}
		}

		if len(rulesHosts) > 0 {
			hosts := make([]string, 0, len(rulesHosts))
			for host := range rulesHosts {
				hosts = append(hosts, host)
			}

			// there is a 63 char limit in the CN of cert-manager/LE
			// so we sort the slice of domain names so the shortest is first
			// if it is over 63 characters, we'll need to synthesize a new one and make it first
			sort.Sort(byLength(hosts))
			if len(hosts[0]) > cnLimit {
				cn, err := makeCN(logger, hosts, cns)
				if err != nil {
					logger.Warningf("unable to append %s tls block: %v", name, err)
					return false, nil
				}
				hosts = append([]string{cn}, hosts...)
			}

			// create the IngressTLS Object with our extra hosts and a custom secret
			secret := fmt.Sprintf("auto-%s-tls", name)
			newtls := v1beta1.IngressTLS{
				Hosts:      hosts,
				SecretName: secret,
			}
			spec.TLS = append(spec.TLS, newtls)
			logger.Debugf("appending %s tls block: %+v", name, newtls)

			// append Cert-manager annotation if not present already
			annotations := ingress.ObjectMeta.Annotations
			addAnnotation := true
			if len(annotations) > 0 {
				for _, annotationKey := range certManagerAnnotations {
					if _, exist := annotations[annotationKey]; exist {
						addAnnotation = false
						break
					}
				}
			}
			if addAnnotation {
				if ingress.ObjectMeta.Annotations == nil {
					ingress.ObjectMeta.Annotations = make(map[string]string)
				}
				ingress.ObjectMeta.Annotations[certManagerAnnotations[0]] = "true"
				logger.Debugf("appending %s annotation: %s: true", name, certManagerAnnotations[0])
			}
		}

		return false, nil
	})

	cfg := mutating.WebhookConfig{
		Name: "tls-host-controller",
		Obj:  &v1beta1.Ingress{},
	}

	wh, err := mutating.NewWebhook(cfg, mt, nil, nil, logger)
	if err != nil {
		panic(err)
	}

	// Get the handler for our webhook.
	whHandler, err := whhttp.HandlerFor(wh)
	if err != nil {
		panic(err)
	}

	watcher, err := fswatcher.New("/data/tls.crt", "/data/tls.key")
	if err != nil {
		logger.Errorf("unable to read server certificate. err='%s'", err)
		os.Exit(1)
	}
	sentinel := certinel.New(watcher, func(err error) {
		logger.Warningf("certinel was unable to reload the certificate. err='%s'", err)
	})

	sentinel.Watch()

	server := http.Server{
		Addr:    ":4443",
		Handler: whHandler,
		TLSConfig: &tls.Config{
			GetCertificate: sentinel.GetCertificate,
		},
	}

	logger.Infof("Listening on :4443")
	err = server.ListenAndServeTLS("", "")

	if err != nil {
		panic(err)
	}
}
