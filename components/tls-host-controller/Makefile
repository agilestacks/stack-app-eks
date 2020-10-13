.DEFAULT_GOAL := deploy

export COMPONENT_NAME ?= tls-host-controller
export NAMESPACE      ?= kube-system
REGISTRY              ?= agilestacks
IMAGE                 ?= $(REGISTRY)/$(COMPONENT_NAME)
IMAGE_VERSION         ?= $(shell git rev-parse HEAD | colrm 7)
IMAGE_TAG             ?= latest
REGISTRY_PASS         ?= ~/.docker/agilestacks.txt

docker  := docker
kubectl := kubectl --context="$(DOMAIN_NAME)" --namespace="$(NAMESPACE)"

ifneq (,$(filter tls-ingress,$(HUB_PROVIDES)))
deploy: purge create install
else
deploy:
	@echo tls-ingress must be provided to install tls-host-controller
endif

undeploy: purge

purge: export DOMAIN_NAME ?= $(error DOMAIN_NAME must be set)

create:
	deploy/create_cm_issuer_and_cert.sh

install:
	$(kubectl) apply -f deploy/manifests.yaml

purge:
	-$(kubectl) delete certificate tls-host-controller
	-$(kubectl) delete secret tls-host-controller-certs
	-$(kubectl) delete secret cm-util-ca
	-$(kubectl) delete mutatingwebhookconfiguration tls-host-controller
	-$(kubectl) delete issuer util-ca
	-$(kubectl) delete -f deploy/manifests.yaml

build:
	$(docker) build -f Dockerfile -t $(IMAGE):$(IMAGE_VERSION) -t $(IMAGE):$(IMAGE_TAG) .
.PHONY: build

push: login push-version push-tag
.PHONY: push

push-version:
	$(docker) push $(IMAGE):$(IMAGE_VERSION)
.PHONY: push-version

push-tag:
	$(docker) tag $(IMAGE):$(IMAGE_VERSION) $(IMAGE):$(IMAGE_TAG)
	$(docker) push $(IMAGE):$(IMAGE_TAG)
.PHONY: push-tag

pull-latest:
	docker pull $(IMAGE):latest
.PHONY: pull-latest

push-stable: pull-latest
	$(MAKE) push-tag IMAGE_VERSION=latest IMAGE_TAG=stable
.PHONY: push-stable

push-stage: pull-latest
	$(MAKE) push-tag IMAGE_VERSION=latest IMAGE_TAG=stage
.PHONY: push-stage

push-preview: pull-latest
	$(MAKE) push-tag IMAGE_VERSION=latest IMAGE_TAG=preview
.PHONY: push-preview

login:
	@ touch $(REGISTRY_PASS)
	@ echo "Please put Docker Hub password into $(REGISTRY_PASS)"
	cat $(REGISTRY_PASS) | docker login --username agilestacks --password-stdin
.PHONY: login
