# Kubernetes Application Development Stack on EKS

App stack deploys Traefik ingress controller, External DNS, TLS via Cert-manager, Dex for SSO, Prometheus, Kube Dashboard, Harbor, PostgreSQL, and optionally Minio for application development.

## Installation

Create EKS cluster in AWS Console or with [eksctl](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html) --config-file [etc/eks-cluster.yaml](etc/eks-cluster.yaml). In case you have an existing EKS cluster, please make sure `certManager` and `externalDNS` addon policies are deployed.

### 1. Create EKS cluster

Edit `etc/eks-cluster.yaml` to set cluster name and AWS region.

```
$ eksctl create cluster -f etc/eks-cluster.yaml
```

### 2. Install Hub CLI

Install [Hub CLI](https://docs.agilestacks.com/article/zrban5vpb5-install-toolbox#hub_cli) and then install extensions:

```
$ hub extensions install
```

### 3. Configure stack

```
$ hub configure --current-kubecontext --force
```

The command will use current Kubeconfig context (setup by `eksctl`), generates stack configuration, then stores it in environment file `.env` (a symlink to current active configuration):

* Obtain an unique sub-domain name under `devops.delivery` domain (valid for 72 hours, refreshed by `hub aws init` command)
* Domain refresh key - a token authorizing domain name refresh
* AWS configuration - S3 bucket to store deployment state

Source configuration into current shell:

```
source .env
```

### 4. Install cloud resources

This is one-time operation per cloud account to create S3 bucket and one-time operation per stack to grab an unique DNS domain.

```
$ hub aws init
```

The obtained DNS subdomain (of `devops.delivery`) is valid for 72h. To renew the lease please re-run `hub aws init` every other day.

### 5. Deploy current stack

```
$ hub stack deploy
```


## What's next?

### Kubernetes Application

We have a tutorial on a simple Kubernetes application [with Skaffold](https://docs.agilestacks.com/article/4b2q2dcof9-development-workflow-on-kubernetes-with-skaffold). Then a follow up tutorial to enable application to [access PostgreSQL](https://docs.agilestacks.com/article/j4cysq9ka5-201-python-efficient-development-for-kubernetes-enable-database).


### Redeploy one or more stack components

```
$ hub stack undeploy -c harbor,prometheus
$ hub stack deploy -c harbor,prometheus
```

### Deploy components starting from a particular component

```bash
$ hub stack deploy -o prometheus
```

## Teardown

```
$ hub stack undeploy
$ eksctl delete cluster -f etc/eks-cluster.yaml
```
