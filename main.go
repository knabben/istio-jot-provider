package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.
var hmacSecret = []byte("secret")

type Request struct {
	jwt.StandardClaims
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
		body  []byte
		token string
	)
	req := &Request{}

	if !checkCookieHeader(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.IssuedAt = time.Now().Add(1 * time.Hour).Unix()
	if err = json.Unmarshal(body, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Sign and get the complete encoded token as a string using the secret
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, req)
	if token, err = t.SignedString(hmacSecret); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("JWT %s", token))
	fmt.Fprintf(w, "")
}

func main() {
	http.HandleFunc("/proxy", jwtProxyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
