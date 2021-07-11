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

### Compiling image

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

An example of usage, the IAT is set to expire in 1 hour.

```
 http localhost:8080/proxy Cookie:anyvaluehere aud=myaudience sub=system
HTTP/1.1 200 OK
Authorization: JWT eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJteWF1ZGllbmNlIiwiaWF0IjoxNjI2MDE4MDA5LCJzdWIiOiJzeXN0ZW0ifQ.6CQf3_YDZVSyvaqfdU4Wq7sOSBff0U6TL8SAyCa-p6s
Content-Length: 0
Date: Sun, 11 Jul 2021 14:40:09 GMT
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