package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/livecodecreator/docker-heroku/raspberrypi"
	"github.com/livecodecreator/docker-heroku/slack"
)

func responseLoggerMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request method: %v\n", r.Method)
		log.Printf("request requestURI: %v\n", r.RequestURI)
		log.Printf("request url scheme: %v\n", r.URL.Scheme)
		log.Printf("request url host: %v\n", r.URL.Host)
		log.Printf("request url path: %v\n", r.URL.Path)
		log.Printf("request url raw query: %v\n", r.URL.RawQuery)

		for k, v := range r.Header {
			log.Printf("request header: %v: %v\n", k, v)
		}

		next.ServeHTTP(w, r)
	})
}

// DefaultHandler is
func defaultHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNotFound)
}

func main() {
	log.SetFlags(log.Lshortfile)

	port := os.Getenv("PORT")
	r := mux.NewRouter()
	r.Use(responseLoggerMiddleware)
	r.HandleFunc("/", defaultHandler).Methods("GET")
	r.HandleFunc("/slack", slack.SlackHandler)
	r.HandleFunc("/status", defaultHandler).Methods("GET")
	r.HandleFunc("/rasp/status", raspberrypi.PostStatusHandler).Methods("POST")
	r.HandleFunc("/{wildcard}", defaultHandler)
	http.Handle("/", r)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Printf("%v\n", err)
	}
}
