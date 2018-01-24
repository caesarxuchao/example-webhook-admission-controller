# Deprecated

This example works for Kubernetes v1.8. To use webhooks in 1.9, please refer to
the Kubernetes e2e test for the webhook
[setup](https://github.com/kubernetes/kubernetes/blob/release-1.9/test/e2e/apimachinery/webhook.go)
and
[implementation](https://github.com/kubernetes/kubernetes/tree/release-1.9/test/images/webhook).
Note that the authentication model has changed in 1.9, so the two-way tls shown
in this repository does not work in v1.9. The webhook in the e2e test uses
[one-way tls](https://github.com/kubernetes/kubernetes/blob/release-1.9/test/images/webhook/config.go#L48-L49).

# Kubernetes External Admission Webhook Example

The example shows how to build and deploy an external webhook that only admits
pods creation and update if the container images have the "grc.io" prefix.

## Prerequisites
Please use a Kubernetes release at least as new as v1.8.0 or v1.9.0-alpha.1,
because the generated server cert/key only works with Kubernetes release that
contains this [change](https://github.com/kubernetes/kubernetes/pull/50476).
Please checkout the `pre-v1.8` tag for an example that works with older
clusters.

Please enable the admission webhook feature
([doc](https://kubernetes.io/docs/admin/extensible-admission-controllers/#enable-external-admission-webhooks)).

## Build the code

```bash
make build
```

## Deploy the code

```bash
make deploy-only 
```

The Makefile assumes your cluster is created by the
[hack/local-up-cluster.sh](https://github.com/kubernetes/kubernetes/blob/master/hack/local-up-cluster.sh).
Please modify the Makefile accordingly if your cluster is created differently.

## Explanation on the CAs/Certs/Keys

The apiserver initiates a tls connection with the webhook, so the apiserver is
the tls client, and the webhook is the tls server.

The webhook proves its identity by the `serverCert` in the certs.go. The server
cert is signed by the CA in certs.go. To let the apiserver trust the `caCert`,
the webhook registers itself with the apiserver via the
`admissionregistration/v1alpha1/externalAdmissionHook` API, with
`clientConfig.caBundle=caCert`.

For maximum protection, this example webhook requires and verifies the client
(i.e., the apiserver in this case) cert. The cert presented by the apiserver is
signed by a client CA, whose cert is stored in the configmap
`extension-apiserver-authentication` in the `kube-system` namespace. See the
`getAPIServerCert` function for more information. Usually you don't need to
worry about setting up this CA cert. It's taken care of when the cluster is
created. You can disable the client cert verification by setting the
`tls.Config.ClientAuth` to `tls.NoClientCert` in `config.go`.
