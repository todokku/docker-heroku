package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	fmt.Println("OK")
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Printf("err: %v\n", err)
	}
}

// HomeHandler is
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request url: %s\n", r.URL)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World")
	fmt.Printf("%v\n", r.Body)
}
