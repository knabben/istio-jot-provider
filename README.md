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

### Compiling image and installation

Compile the proxy image (authentication service) and install the deployment, the service
will be installed under `istio-system`.

```
$ make push  # knabben/proxy:latest
$ make install
```

Is expected to receive a `Cookie: anyvalue` header and the reponse will return an
`Authorization: JWT ...` header, the payload can control the following claims: 

```
Audience  string `json:"aud,omitempty"`
ExpiresAt int64  `json:"exp,omitempty"`
Id        string `json:"jti,omitempty"`
Issuer    string `json:"iss,omitempty"`
NotBefore int64  `json:"nbf,omitempty"`
Subject   string `json:"sub,omitempty"`
```

Install the client helper:

```
kubectl apply -f https://raw.githubusercontent.com/istioinaction/book-source-code/8dad67e5f1939b0f9dc22b7d6204709cf96d1cc1/ch9/sleep.yaml
```

Validating the authorizer service usage, the IAT is set to expire in 1 hour.

```
❯ kubectl -n default exec -it deploy/sleep -- curl auth-proxy.istio-system:8080/proxy \
    -H "Cookie: value" -d'{"sub": "system"}' -v

*   Trying 10.96.173.134:8080...
* Connected to auth-proxy.istio-system (10.96.173.134) port 8080 (#0)
> POST /proxy HTTP/1.1
> Host: auth-proxy.istio-system:8080
> User-Agent: curl/7.77.0
> Accept: */*
> Cookie: value
> Content-Length: 17
> Content-Type: application/x-www-form-urlencoded
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Authorization: JWT eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2MjYwMTkzMjYsInN1YiI6InN5c3RlbSJ9.WY2a4EbIoqR6kgwLOaAH9nslc2Yu9GzDT7DOyiguRRE
< Date: Sun, 11 Jul 2021 15:02:06 GMT
< Content-Length: 0
< 
```

## 001 - EnvoyFilter CRD

Using EnvoyFilter CRD

## 010 - AuthorizationFilter CUSTOM

Using AuthorizationPolicy with `action: CUSTOM`

# Cleaning up

```
$ ./hack/cleanup.sh
$ kind delete cluster --name calico-2021-07-11-sg33lz
```