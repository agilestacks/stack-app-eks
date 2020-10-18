# Kubernetes Application Development Stack on EKS

App stack deploys Traefik ingress controller, External DNS, TLS via Cert-manager, Dex for SSO, Prometheus, Kube Dashboard, Harbor, PostgreSQL, and optionally Minio for application development.


## Installation

Create EKS cluster in AWS Console or with [eksctl](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html) --config-file [etc/eks-cluster.yaml](etc/eks-cluster.yaml).

### 1. Create EKS cluster

Edit `etc/eks-cluster.yaml` to set cluster name and AWS region.

```
$ eksctl create cluster -f etc/eks-cluster.yaml
```

In case you have an existing EKS cluster, then use the following example command to configure EKS cluster context:

```
$ aws eks --region us-east-2 update-kubeconfig --name cluster-name
```

Also, please make sure `certManager` and `externalDNS` addon policies are deployed.

### 2. Install Hub CLI

Install [Hub CLI](https://docs.agilestacks.com/article/zrban5vpb5-install-toolbox#hub_cli) and then install extensions:

```
$ hub extensions install
```

Hub CLI Extensions require [AWS CLI](https://aws.amazon.com/cli/), [KUBECTL](https://kubernetes.io/docs/reference/kubectl/overview/), [EKSCTL](https://eksctl.io), [JQ](https://stedolan.github.io/jq/), [YQ v3](https://github.com/mikefarah/yq). Optionally install [Node.js](https://nodejs.org) and NPM for `hub pull` extension.

### 3. Configure stack

```
$ hub configure -f hub.yaml
```

The command will use current Kubeconfig context (setup by `eksctl`), generates stack configuration, then stores it in environment file `.env` (a symlink to current active configuration):

* Unique sub-domain name under `bubble.superhub.io` domain
* Domain refresh key - a token authorizing domain name refresh
* AWS configuration - S3 bucket to store deployment state
* Stack parameters such as usernames and passwords

The obtained DNS subdomain (of `bubble.superhub.io`) is valid for 72h. To renew the DNS lease please run `hub configure -r aws --dns-update` every other day.  For a permanent DNS domain please contact support@agilestacks.com.

### 4. Deploy

```
$ hub stack deploy
```

## Access the Kubernetes Dashboard

You can find the URL of deployed Kubeflow Dashboard using the `hub show` command:
```bash
$ hub show
```
Note: `hub show` command will the return values for all stack parameters. Kubernetes Dashboard URL is shown as value for parameter `kube-dashboard:component.kubernetes-dashboard.url` with login username `component.dex.passwordDb.email` and password `component.dex.passwordDb.password`

## What's next?

### Kubernetes Application

**Work in progress** We have a tutorial on a simple Kubernetes application [with Skaffold](https://docs.agilestacks.com/article/4b2q2dcof9-development-workflow-on-kubernetes-with-skaffold). Then a follow up tutorial to enable application to [access PostgreSQL](https://docs.agilestacks.com/article/j4cysq9ka5-201-python-efficient-development-for-kubernetes-enable-database).


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
