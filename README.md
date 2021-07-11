# Istio JWT Custom Provider

## Initialing the environment

Setting up a cluster with Calico CNI and MetalLB

```
# Using https://github.com/thekubeworld/k8s-local-dev
$ ./k8s-local-dev calico
$ kind get clusters
$ ~/go/src/github.com/k8sbykeshed/k8s-service-lb-validator/hack/install_metallb.sh`

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

### Install the authentication service

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

### Install the end-user accessible service

Applying the north-south service without any authentication filter,
accessed via the Ingress Gateway, the LB external IP is on `172.18.0.40`:  

```
$ kubectl apply -f service/httpbin.yaml
$ istioctl -nistio-system proxy-config routes deploy/istio-ingressgateway                                                                                                                                                           [13:08:10]
NAME        DOMAINS     MATCH                  VIRTUAL SERVICE
http.80     *           /*                     httpbin-virtualservice.default

$ curl -X GET "http://172.18.0.40/headers" -H "accept: application/json"                                                                                                                                                            [13:11:31]
{
  "headers": {
    "Accept": "application/json", 
    "Host": "172.18.0.40", 
    "User-Agent": "curl/7.68.0", 
    "X-B3-Sampled": "1", 
    "X-B3-Spanid": "165c8291946d9062", 
    "X-B3-Traceid": "60fa94b4ebc5385d165c8291946d9062", 
    "X-Envoy-Attempt-Count": "1", 
    "X-Envoy-Decorator-Operation": "httpbin.default.svc.cluster.local:80/*", 
    "X-Envoy-Internal": "true", 
    "X-Envoy-Peer-Metadata": "ChQKDkFQUF9DT05UQUlORVJTEgIaAAoaCgpDTFVTVEVSX0lEEgwaCkt1YmVybmV0ZXMKGQoNSVNUSU9fVkVSU0lPThIIGgYxLjEwLjIKvwMKBkxBQkVMUxK0AyqxAwodCgNhcHASFhoUaXN0aW8taW5ncmVzc2dhdGV3YXkKEwoFY2hhcnQSChoIZ2F0ZXdheXMKFAoIaGVyaXRhZ2USCBoGVGlsbGVyCjYKKWluc3RhbGwub3BlcmF0b3IuaXN0aW8uaW8vb3duaW5nLXJlc291cmNlEgkaB3Vua25vd24KGQoFaXN0aW8SEBoOaW5ncmVzc2dhdGV3YXkKGQoMaXN0aW8uaW8vcmV2EgkaB2RlZmF1bHQKMAobb3BlcmF0b3IuaXN0aW8uaW8vY29tcG9uZW50EhEaD0luZ3Jlc3NHYXRld2F5cwohChFwb2QtdGVtcGxhdGUtaGFzaBIMGgo1ZDU3OTU1NDU0ChIKB3JlbGVhc2USBxoFaXN0aW8KOQofc2VydmljZS5pc3Rpby5pby9jYW5vbmljYWwtbmFtZRIWGhRpc3Rpby1pbmdyZXNzZ2F0ZXdheQovCiNzZXJ2aWNlLmlzdGlvLmlvL2Nhbm9uaWNhbC1yZXZpc2lvbhIIGgZsYXRlc3QKIgoXc2lkZWNhci5pc3Rpby5pby9pbmplY3QSBxoFZmFsc2UKGgoHTUVTSF9JRBIPGg1jbHVzdGVyLmxvY2FsCi8KBE5BTUUSJxolaXN0aW8taW5ncmVzc2dhdGV3YXktNWQ1Nzk1NTQ1NC1qYm10cAobCglOQU1FU1BBQ0USDhoMaXN0aW8tc3lzdGVtCl0KBU9XTkVSElQaUmt1YmVybmV0ZXM6Ly9hcGlzL2FwcHMvdjEvbmFtZXNwYWNlcy9pc3Rpby1zeXN0ZW0vZGVwbG95bWVudHMvaXN0aW8taW5ncmVzc2dhdGV3YXkKFwoRUExBVEZPUk1fTUVUQURBVEESAioACicKDVdPUktMT0FEX05BTUUSFhoUaXN0aW8taW5ncmVzc2dhdGV3YXk=", 
    "X-Envoy-Peer-Metadata-Id": "router~192.168.142.70~istio-ingressgateway-5d57955454-jbmtp.istio-system~istio-system.svc.cluster.local"
  }
}
```

## 001 - EnvoyFilter CRD

Using EnvoyFilter CRD

## 010 - AuthorizationPolicy CUSTOM filter 

Using AuthorizationPolicy with `action: CUSTOM`

## Cleaning up

```
$ ./hack/cleanup.sh
$ kind delete cluster --name calico-2021-07-11-sg33lz
```

## References

* MetalLB - https://github.com/K8sbykeshed/k8s-service-lb-validator
* Kind bootstrap - https://github.com/thekubeworld/k8s-local-dev
* Istio In Action - https://github.com/istioinaction/book-source-code/
* EnvoyFilter example - https://github.com/thiagocaiubi/playground/tree/main/istio-okta
