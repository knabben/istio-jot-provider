# Istio JWT Custom Provider

## Initialing the environment

Setting up a cluster with Calico CNI

```
# Using https://github.com/thekubeworld/k8s-local-dev
$ ./k8s-local-dev calico
$ kind get clusters
calico-2021-07-11-sg33lz
```

Install istiod and export the path

```
$ ./hack/setup-istio.sh
✔ No issues found when checking the cluster. Istio is safe to install or upgrade!
  To get started, check out https://istio.io/latest/docs/setup/getting-started/
✔ Istio core installed                                                                                                 
✔ Istiod installed                                                                                                     
✔ Egress gateways installed                                                                                            
✔ Ingress gateways installed                                                                                           
✔ Installation complete                                                                                                

Thank you for installing Istio 1.10.  Please take a few minutes to tell us about your install/upgrade experience!  https://forms.gle/KjkrDnMPByq7akrYA
```

# Cleaning up

```
$ ./hack/cleanup.sh
$ kind delete cluster --name calico-2021-07-11-sg33lz
```