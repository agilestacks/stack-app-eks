
# TLS Host Controller

This is an admission webhook mutator. This routine receives a manifest submission in front of the API. If it meets filter criteria, will mutate the manifest, and passes the mutated manifest along.

This specific webhook looks only for `Ingress` manifests.  Its logic is as follows:

```
if spec.TLS exists:
   exit

if spec.rules.hosts exists:
  collect every value in spec.rules[*].hosts into a list
  create a new `IngressTLS` struct
  add all of the hosts from the list to the TLS hosts
  if lenght of hosts[0] > 63 (ACME/LE CN limit):
     synthesize new CN and put it in front of the list so that it became the CN of the cert
  if no any TLS annotation exists:
    append `kubernetes.io/tls-acme` annotation
```

Note that `rules.host` will never be modified under any circumstance, even if there are rules in `tls.host` that are not present in `rules.host`.
Note that any existing `tls` blocks or `tls.host` entries will not be modified. The change happens only if `tls` block is absent.

# Deploying

This container is intended to be deployed into customer infrastructure (namely on-prem, but should work anywhere).
Therefore, this container is hosted as a public resource in hub.docker.com/agilestacks/tls-host-controller

# Building and Pushing

This is meant to be built and deployed in a container. To create a new docker image, run `make build`:

```
docker login #(into agilestacks account)
make build push
```
