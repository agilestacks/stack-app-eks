# Kubernetes Application Development Stack on EKS

App stack deploys External DNS, Traefik v2 ingress controller, TLS via Cert-manager with Let's Encrypt, Dex for SSO, Kubernetes Dashboard v2, Prometheus Operator, Harbor, PostgreSQL, and optionally Minio for application development.


## Installation

Create EKS cluster in AWS Console or with [eksctl](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html) --config-file [etc/eks-cluster.yaml](etc/eks-cluster.yaml).

### 1. Create EKS cluster

Edit `etc/eks-cluster.yaml` to set cluster name and AWS region.

    $ eksctl create cluster -f etc/eks-cluster.yaml

In case you have an existing EKS cluster, then use the following example command to configure EKS cluster context:

    $ aws eks --region us-east-2 update-kubeconfig --name cluster-name

Please make sure `certManager` and `externalDNS` addon policies are deployed.

Also, confirm `AWS_PROFILE` setting used to deploy the cluster (and generate Kubeconfig) match profile for the next step. Using different profiles - even if they're pointing to the same AWS account - may not work.

### 2. Install Hub CLI

Install [Hub CLI](https://docs.agilestacks.com/article/zrban5vpb5-install-toolbox#hub_cli):

    curl -O https://controlplane.agilestacks.io/dist/hub-cli/hub.linux_amd64
    mv hub.linux_amd64 hub
    chmod +x hub
    sudo mv hub /usr/local/bin

There are [Linux amd64](https://controlplane.agilestacks.io/dist/hub-cli/hub.linux_amd64), [Linux arm64](https://controlplane.agilestacks.io/dist/hub-cli/hub.linux_arm64), and [macOS amd64](https://controlplane.agilestacks.io/dist/hub-cli/hub.darwin_amd64) binaries.

Install extensions:

    $ hub extensions install

Hub CLI Extensions require [AWS CLI], [kubectl], [eksctl], [jq], [yq v4]. Optionally install [Node.js] and NPM for `hub pull` extension.

Windows users please [read on](https://docs.agilestacks.com/article/u6a9cq5yya-hub-cli-on-windows).

#### macOS users

Depending on your machine Security & Privacy settings and macOS version (10.15+), you may get an error _cannot be opened because the developer cannot be verified_. Please [read on](https://github.com/hashicorp/terraform/issues/23033#issuecomment-542302933) for a simple workaround.

Alternativelly, to set global preference to _Allow apps downloaded from: Anywhere_, execute:

    $ sudo spctl --master-disable

For `tls-host-controller` to deploy we require `v3_ca` extension to be configured. Either your `openssl` binary is from OpenSSL project or you must add a few lines to `/etc/ssl/openssl.cnf` LibreSSL config:

    [ v3_ca ]
    basicConstraints = critical,CA:TRUE
    subjectKeyIdentifier = hash
    authorityKeyIdentifier = keyid:always,issuer:always

### 3. Configure stack

    $ hub configure -f hub.yaml

The command will use current Kubeconfig context (setup by `eksctl`), generates stack configuration, then stores it in environment file `.env` (a symlink to current active configuration):

* Unique sub-domain name under `bubble.superhub.io` domain
* Domain refresh key - a token authorizing domain name refresh
* AWS configuration - S3 bucket to store deployment state
* Stack parameters such as usernames and passwords

The obtained DNS subdomain (of `bubble.superhub.io`) is valid for 72h. To renew the DNS lease please run `hub configure -r aws --dns-update` every other day.  For a permanent DNS domain please contact support@agilestacks.com.

### 4. Deploy

    $ hub stack deploy

### 5. Access Kubernetes Dashboard

You can find the URL of deployed Kubernetes Dashboard using `hub show` command:

    $ hub show

The command will return the values for all stack parameters. Kubernetes Dashboard URL is shown as value for parameter `kube-dashboard:component.kubernetes-dashboard.url` with login username `component.dex.passwordDb.email` and password `component.dex.passwordDb.password`.


## What's next?

### Kubernetes Application

Deploy a Python application [with Skaffold and Hub CLI Tutorial](https://docs.agilestacks.com/article/3pbulps5n7-simplifying-kubernetes-for-developers-with-hub-cli-and-skaffold). Then you can enable the application to [access PostgreSQL](https://docs.agilestacks.com/article/j4cysq9ka5-201-python-efficient-development-for-kubernetes-enable-database).

### Redeploy one or more stack components

    $ hub stack undeploy -c harbor,prometheus
    $ hub stack deploy -c harbor,prometheus

### Deploy components starting from a particular component

    $ hub stack deploy -o prometheus

### Undeploy components up to a particular component

    $ hub stack undeploy -l prometheus

### Update elaborate file

If you change `hub.yaml`, or `params.yaml`, or any of `hub-component.yaml`-s, you may want to forcibly re-elaborate:

    $ hub stack elaborate


## Teardown

    $ hub stack undeploy
    $ eksctl delete cluster -f etc/eks-cluster.yaml


[AWS CLI]: https://aws.amazon.com/cli/
[kubectl]: https://kubernetes.io/docs/reference/kubectl/overview/
[eksctl]: https://eksctl.io
[jq]: https://stedolan.github.io/jq/
[yq v4]: https://github.com/mikefarah/yq
[Node.js]: https://nodejs.org
