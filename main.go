package main

import (
	"fmt"
	"log"
	"net/http"
)

type Request struct {

}

func jwtProxyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/proxy", jwtProxyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
