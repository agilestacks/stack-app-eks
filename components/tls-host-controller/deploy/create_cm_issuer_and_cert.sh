#!/bin/bash -xe

function error_exit {
	echo "$1" >&2
	exit -1
}

KUBECTL="kubectl --context=${DOMAIN_NAME} --namespace=${NAMESPACE}"

[ -z "$NAMESPACE" ] && error_exit "NAMESPACE env var must be set"

# create the PK for our new CA
openssl ecparam -out ca.private.key -name prime256v1 -genkey -noout

#use the key to generate a ca cert to be used by cert-manager
openssl req -x509 -new -nodes -key ca.private.key -subj "/CN=cluster.local" -days 3650 -reqexts v3_req -extensions v3_ca -out ca.crt

# create the ca cert tls secret to be used by cert-manager
$KUBECTL create secret tls cm-util-ca \
	--cert=ca.crt \
	--key=ca.private.key \
	--namespace=${NAMESPACE}

cat <<EOF | $KUBECTL apply -f -
apiVersion: cert-manager.io/v1alpha2
kind: Issuer
metadata:
  name: util-ca
  namespace: ${NAMESPACE}
spec:
  ca:
    secretName: cm-util-ca
EOF

cat <<EOF | $KUBECTL apply -f -
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: tls-host-controller
  namespace: ${NAMESPACE}
spec:
  secretName: tls-host-controller-certs
  issuerRef:
    name: util-ca
    kind: Issuer
  commonName: tls-host-controller.kube-system.svc.cluster.local
  organization:
  - AgileStacks
  dnsNames:
  - tls-host-controller.kube-system.svc.cluster.local
  - tls-host-controller.kube-system.svc.cluster
  - tls-host-controller.kube-system.svc
  - tls-host-controller.kube-system
  - tls-host-controller
EOF

cat <<EOF | $KUBECTL apply -f -
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  name: tls-host-controller
  labels:
    app: tls-host-controller
webhooks:
- name: tls-host-controller.${NAMESPACE}.svc.cluster.local
  clientConfig:
    service:
      name: tls-host-controller
      namespace: ${NAMESPACE}
      path: "/mutate"
    caBundle: $(cat ca.crt | base64 | tr -d '\n')
  rules:
  - operations: [ "CREATE", "UPDATE" ]
    apiGroups: ["networking.k8s.io", "extensions"]
    apiVersions: ["v1beta1"]
    resources: ["ingresses"]
EOF
