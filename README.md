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
Issuer    string `json:"iss,omitempty"`
Subject   string `json:"sub,omitempty"`
```

Install the client helper:

```
kubectl apply -f https://raw.githubusercontent.com/istioinaction/book-source-code/8dad67e5f1939b0f9dc22b7d6204709cf96d1cc1/ch9/sleep.yaml
```

Validating the authorizer service usage, the IAT is set to expire in 1 hour.

```
❯ kubectl -n default exec -it deploy/sleep -- curl auth-proxy.istio-system:8080/proxy \
    -H "Cookie: value" -H "Sub: subject" -H "Aud: audience" -v

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
< Authorization: JWT eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhdWRpZW5jZSIsImV4cCI6MTYyNjAyNjg2NywiaWF0IjoxNjI2MDIzMjY3LCJzdWIiOiJzdWJqZWN0In0.gliLIgNPNop5arA_CONY0ewgGnEDF1jnRyd0RQ4xB_A
< Date: Sun, 11 Jul 2021 15:02:06 GMT
< Content-Length: 0
< 
```

### Install the end-user accessible service

Applying the north-south service without any authentication filter,
accessed via the Ingress Gateway, the LB external IP is on `172.18.0.40`:  

```
$ kubectl apply -f service/httpbin.yaml
$ istioctl -nistio-system proxy-config routes deploy/istio-ingressgateway
NAME        DOMAINS     MATCH                  VIRTUAL SERVICE
http.80     *           /*                     httpbin-virtualservice.default

$ curl -X GET "http://172.18.0.40/headers" -H "accept: application/json"
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

Using EnvoyFilter CRD in the wildcards request path of the authentication service:

```
$ kubectl apply -f 01-envoy-filter/envoy-filter.yaml
$ kubectl -n istio-system get envoyfilter ext-authz-filter -o yam

...
typed_config:
  '@type': type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthz
  http_service:
    authorization_request:
      allowed_headers:
        patterns:
        - exact: cookie
        - exact: sub
        - exact: aud
        - exact: iss
    authorization_response:
      allowed_upstream_headers:
        patterns:
        - exact: authorization
    server_uri:
      cluster: outbound|8080||auth-proxy.istio-system.svc.cluster.local
      timeout: 1s
      uri: http://auth-proxy.istio-system:8080
...
```

Testing the request in the httpbin:

```
# -- Without authentication

$ http 172.18.0.40/anything
HTTP/1.1 401 Unauthorized

# -- With authentication

$ http 172.18.0.40/anything cookie:anyvalue iss:my-issuer sub:my-subject
HTTP/1.1 200 OK

{
    "headers": {
        "Authorization": "Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjYwNDI4NjksImlhdCI6MTYyNjAzOTI2OSwiaXNzIjoibXktaXNzdWVyIiwic3ViIjoibXktc3ViamVjdCJ9.Y5Caeltn-1ZbDfEt0byq9_MW-N0_axPI4Rq6UBc0z1bY_BXGeiY_K0H3XLqtnvjw6hjwhz-yD9utPSUoBRpjpI9s9rUsMRO1zvxWz81DgiJUOcjftoVIshLdl1rHZYC1PqP6cJ2oadpHuiGpp7afIrrYBHULMyX27qYYuruwgz05552a5SdHtezTyFEKKvBVrFS_p7upkRsCfIKK3MHBcmeRZBoUXUtZ19fbW_vaGdxo-VU4TqJhWoYbwi93F0bdlkIIf_Eec-1Qww9hB2IjVw-ezruN_oCfQh1uw4Cp1mVEx-ElxoUJb-sajOTaaeClNBPKlh8jGS9HE-boV1i9DcPm-OI9A63KJuZqq2Zco3HoyoUiWODRLHoNGrp8blLBfpvAisexiQtKAiLXPcfIP1LBsWIkHqg1NthIiGpC5DKzN9oNy_dFx17JR1MymoXSTT9KwpTdD7CVWWnuoKKysOTDcZnhF2FFfMo0ubEYhkB89gbSaUUubTbrxbZJlBhX-IpHaK-cbzRWeDDLIufyojc2dlQayA-9ku_uswO6S-JfIBaMDB4cmxBAutxvWRyGdOzx9yIzxpgMEsU9VwcHU1QQh_xRiPl6FmONLj1mK_u0e-XFpJbnc3k5D5FnQYzVJHovn1-z2RbzJjM41fvgJQbNahP-WSCVK0HzH7ITu9A"
    }
}
```

Analyzing the JWT data generated:

```
# Header

{
  "alg": "RS512",
  "typ": "JWT"
}

# Payload

{
  "exp": 1626042869,
  "iat": 1626039269,
  "iss": "my-issuer",
  "sub": "my-subject"
}

```

Forcing the policy in the service:

```
$ kubectl apply -f 01-envoy-filter/req-auth.yaml
```

This will apply some objects in the httpbin service:

```
# RequestAuthentication with a JWT set, and an issuer defined.
selector:
    matchLabels:
      app: httpbin
  jwtRules:
    - issuer: "my-issuer"
      jwks: |
        {"keys": [{"kty": ...
        
# AuthorizationPolicy - Denying all requests
# AuthorizationPolicy - Allowing a requestPrincipals as issuer and subject: [ "my-issuer/subject" ]
```

* NOTE: Generate the JWK, with the `keys/generate_jwks.py`
* NOTE2: DONT USE THESE KEYS, generate using openssl.

To test the setup is working try a few a combination for the claims values:

```
$ http 172.18.0.40/anything cookie:mytest iss:my-issuer sub:my-subject
200

# Wrong subject
$ http 172.18.0.40/anything cookie:mytest iss:my-issuer sub:wrong-subject
403 Forbidden - RBAC: access denied

# Wrong issuer
$ http 172.18.0.40/anything cookie:mytest iss:wrong-issuer sub:wrong-subject
401 Unauthorized - Jwt is not configured 
```

Cleaning up the added final filter.

```
$ kubectl -n istio-system delete envoyfilter ext-authz-filter
$ kubectl delete authorizationpolicy,requestauthentication --all
```

## 010 - CUSTOM External authorization policy

Apply the custom policy:

```
$ kubectl apply -f 10-custom/auth-policy.yaml
$ http http://172.18.0.40/headers
403
```

Looking the istiod logs, it needs the the service yet. It is introduced
by the global mesh options, set here:

https://istio.io/latest/docs/tasks/security/authorization/authz-custom/#define-the-external-authorizer

```
2021-07-11T22:22:54.199583Z     error   authorization   Processed authorization policy for httpbin-74fb669cc6-lcg4s.default, 2 errors occurred:

    * failed to process CUSTOM action: available providers are [] but found "ext-authz-service"
    * failed to parse CUSTOM action, will generate a deny all config: available providers are [] but found "ext-authz-service"
```

The configuration reference:

https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/

Editing the configmap:
```
$ kubectl edit configmap istio -n istio-system
data:                                                                                                                                                                                                                                          
  mesh: "extensionProviders:
    - name: "ext-authz-service"
      envoyExtAuthzHttp:                                                                                                                                                            
        service: "auth-proxy.istio-system.svc.cluster.local"
        port: "8080"                                                                                                                                                           
      includeHeadersInCheck: ["cookie", "sub", "aud", "iss"] 
      headersToUpstreamOnAllow: ["authorization"]
 
# Restart istiod    
$ kubectl rollout restart deployment/istiod -n istio-system   
```

At this point we should be able to send the Bearer token, using
the authentication service under `auth-proxy.istio-system`, apply the
same authorization policies and request authentication again: 

```
# Before applying the AuthorizationPolicy and RequestAuthentication
$ http 172.18.0.40/anything cookie:mytest
200

$ kubectl apply -f 01-envoy-filter/req-auth.yaml

$ http 172.18.0.40/anything cookie:mytest
403 Forbidden

$ http 172.18.0.40/anything cookie:mytest iss:my-issuer sub:subject1
403 Forbidden

$ http 172.18.0.40/anything cookie:mytest iss:my-issuer sub:subject
200
```

## Cleaning up istio-system

```
$ ./hack/cleanup.sh
$ kind delete cluster --name calico-2021-07-11-sg33lz
```

## References

* MetalLB - https://github.com/K8sbykeshed/k8s-service-lb-validator
* Kind bootstrap - https://github.com/thekubeworld/k8s-local-dev
* Istio In Action - https://github.com/istioinaction/book-source-code/
* AuthZ JWT https://istio.io/latest/docs/tasks/security/authorization/authz-jwt/
* EnvoyFilter example - https://github.com/thiagocaiubi/playground/tree/main/istio-okta
