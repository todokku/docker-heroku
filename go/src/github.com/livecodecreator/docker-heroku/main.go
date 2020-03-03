package main

import (
	"encoding/json"
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
	r.Use(ResponseLoggerMiddleware)
	r.HandleFunc("/", DefaultHandler)
	r.HandleFunc("/slack", SlackHandler)
	r.HandleFunc("/{wildcard}", DefaultHandler)
	http.Handle("/", r)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Printf("err: %v\n", err)
	}
}

const (
	slackEventTypeURLVerification = "url_verification"
	slackEventTypeCallback        = "event_callback"
	slackEventCallbackTypeText    = "text"
)

// SlackChallengeRequest is
type SlackChallengeRequest struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

// SlackChallengeResponse is
type SlackChallengeResponse struct {
	Challenge string `json:"challenge"`
}

// SlackEventCallbackRequest is
type SlackEventCallbackRequest struct {
	Type  string                  `json:"type"`
	Event SlackEventCallbackEvent `json:"event"`
}

// SlackEventCallbackEvent is
type SlackEventCallbackEvent struct {
	Type    string `json:"type"`
	Message string `json:"text"`
}

// SlackHandler is
func SlackHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "E0001")
		log.Println("E0001")
		return
	}

	if slackChallengeRequest(w, r, b) {
		return
	}

	if slackEventCallbackRequest(w, r, b) {
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "E0000")
	log.Println("E0000")
	return
}

func slackChallengeRequest(w http.ResponseWriter, r *http.Request, b []byte) bool {
	var req SlackChallengeRequest
	err := json.Unmarshal(b, &req)
	if err != nil {
		log.Println("E100A")
		return false
	}

	if req.Type != slackEventTypeURLVerification {
		log.Println("E100B")
		return false
	}
	log.Println("Execute ChallengeRequest")

	scres := SlackChallengeResponse{Challenge: req.Challenge}
	d, err := json.Marshal(scres)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "E1001")
		log.Println("E1001")
		return true
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(d))
	return true
}

func slackEventCallbackRequest(w http.ResponseWriter, r *http.Request, b []byte) bool {
	var req SlackEventCallbackRequest
	err := json.Unmarshal(b, &req)
	if err != nil {
		log.Println("E200A")
		return false
	}

	if req.Type != slackEventTypeCallback {
		log.Println("E100B")
		return false
	}
	log.Println("Execute EventCallbackRequest")

	if req.Event.Type != slackEventCallbackTypeText {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "E2001")
		log.Println("E2001")
		return true
	}

	log.Printf("SlackCallbackMessage: %v\n", req.Event.Message)

	w.WriteHeader(http.StatusOK)
	return true
}

// ResponseLoggerMiddleware is
func ResponseLoggerMiddleware(next http.Handler) http.Handler {
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

		log.Printf("request body:\n")
		defer r.Body.Close()
		b, _ := ioutil.ReadAll(r.Body)
		log.Println(string(b))

		next.ServeHTTP(w, r)
	})
}

// DefaultHandler is
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
