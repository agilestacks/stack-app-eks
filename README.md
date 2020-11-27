# Kubernetes Application Development Stack on EKS

App stack deploys External DNS, Traefik v2 ingress controller, TLS via Cert-manager with Let's Encrypt, Dex for SSO, Kubernetes Dashboard v2, Prometheus Operator, Harbor, PostgreSQL, and optionally Minio for application development.


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

Please make sure `certManager` and `externalDNS` addon policies are deployed.

Also, confirm `AWS_PROFILE` setting used to deploy the cluster (and generate Kubeconfig) match profile for the next step. Using different profiles - even if they're pointing to the same AWS account - may not work.

### 2. Install Hub CLI

Install [Hub CLI](https://docs.agilestacks.com/article/zrban5vpb5-install-toolbox#hub_cli) and then install extensions:

```
$ hub extensions install
```

Hub CLI Extensions require [AWS CLI], [kubectl], [eksctl], [jq], [yq v3]. Optionally install [Node.js] and NPM for `hub pull` extension.

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

### 5. Access Kubernetes Dashboard

You can find the URL of deployed Kubernetes Dashboard using `hub show` command:

```
$ hub show
```

The command will return the values for all stack parameters. Kubernetes Dashboard URL is shown as value for parameter `kube-dashboard:component.kubernetes-dashboard.url` with login username `component.dex.passwordDb.email` and password `component.dex.passwordDb.password`.


## What's next?

### Kubernetes Application

Deploy a Python application [with Skaffold and Hub CLI Tutorial](https://docs.agilestacks.com/article/3pbulps5n7-simplifying-kubernetes-for-developers-with-hub-cli-and-skaffold). Then you can enable the application to [access PostgreSQL](https://docs.agilestacks.com/article/j4cysq9ka5-201-python-efficient-development-for-kubernetes-enable-database).

### Redeploy one or more stack components

```
$ hub stack undeploy -c harbor,prometheus
$ hub stack deploy -c harbor,prometheus
```

### Deploy components starting from a particular component

```
$ hub stack deploy -o prometheus
```

## Teardown

```
$ hub stack undeploy
$ eksctl delete cluster -f etc/eks-cluster.yaml
```


[AWS CLI]: https://aws.amazon.com/cli/
[kubectl]: https://kubernetes.io/docs/reference/kubectl/overview/
[eksctl]: https://eksctl.io
[jq]: https://stedolan.github.io/jq/
[yq v3]: https://github.com/mikefarah/yq
[Node.js]: https://nodejs.org
