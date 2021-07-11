package main

import (
	"io/ioutil"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var keyData []byte

type Request struct {
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}

func init() {
	var err error
	if keyData, err = ioutil.ReadFile("/private.key"); err != nil {
		log.Fatal(err)
	}
}

func (r *Request) Valid() error {
	return nil
}

// checkCookieHeader - it's possible to externally check a
// Django session_id like sessionid=ytlxxsjjs0jepar3qt9q1eyp12drg0mc from
// a cached session, in this example only checking if the cookie is set,
// otherwise we return a 401.
func checkCookieHeader(r *http.Request) bool {
	_, ok := r.Header["Cookie"]
	return ok
}

// jwtProxyHandler validates a Cookie and returns the JWT parsed value.
func jwtProxyHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		token string
	)

	if !checkCookieHeader(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req := &Request{
		Subject: r.Header.Get("sub"),
		Issuer: r.Header.Get("iss"),
		Audience: r.Header.Get("aud"),
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		IssuedAt: time.Now().Unix(),
	}

	log.Print(fmt.Sprintf("Request: %+v", req))

	key, _ := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	// Sign and get the complete encoded token as a string using the secret
	t := jwt.NewWithClaims(jwt.GetSigningMethod("RS512"), req)
	if token, err = t.SignedString(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	fmt.Fprintf(w, "")
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(jwtProxyHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
