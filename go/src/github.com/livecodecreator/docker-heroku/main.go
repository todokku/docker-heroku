package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	r := mux.NewRouter()
	r.HandleFunc("/", DefaultHandler)
	r.HandleFunc("/{wildcard}", DefaultHandler)
	http.Handle("/", r)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Printf("err: %v\n", err)
	}
}

// DefaultHandler is
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request method: %v\n", r.Method)
	log.Printf("request requestURI: %v\n", r.RequestURI)
	log.Printf("request url scheme: %v\n", r.URL.Scheme)
	log.Printf("request url host: %v\n", r.URL.Host)
	log.Printf("request url path: %v\n", r.URL.Path)
	log.Printf("request url raw query: %v\n", r.URL.RawQuery)

	for k, v := range r.Header {
		log.Printf("request header: %v: %v\n", k, v)
	}

	log.Printf("request body:\n")
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)
	log.Println(string(b))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
