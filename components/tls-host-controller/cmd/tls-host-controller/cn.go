package main

import (
	"errors"
	"strings"

	"github.com/slok/kubewebhook/pkg/log"
)

func parseCN(logger log.Logger, cn string) ([]string, error) {
	if len(cn) == 0 {
		return nil, nil
	}
	cns := strings.Split(cn, ",")
	cns2 := make([]string, 0, len(cn))
	for _, c := range cns {
		if len(c) > cnLimit {
			logger.Warningf("%s length > %d - cannot use as CN", c, cnLimit)
		} else {
			cns2 = append(cns2, c)
		}
	}
	return cns2, nil
}

func makeCN(logger log.Logger, hosts, cns []string) (string, error) {
	for _, cn := range cns {
		for _, host := range hosts {
			if strings.HasSuffix(host, "."+cn) {
				return cn, nil
			}
		}
	}
	// compat
	if len(hosts) > 0 {
		host := hosts[0]
		if len(host) <= cnLimit {
			return host, nil
		}
		offset := len(host) - cnLimit
		i := strings.Index(host[offset:], ".")
		if i >= 0 {
			return host[offset+i+1:], nil
		}
	}
	return "", errors.New("unable to match default CN to any of host rules")
}
