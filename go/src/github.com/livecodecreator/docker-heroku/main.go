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

var logger *log.Logger

var lastStatus string

const (
	slackEventTypeURLVerification = "url_verification"
	slackEventTypeCallback        = "event_callback"
	slackEventCallbackTypeMessage = "message"
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

// StatusResponse is
type StatusResponse struct {
	LastStatus string `json:"lastStatus"`
}

// SlackHandler is
func SlackHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Printf("%v\n", err)
		return
	}

	logger.Println("request body:")
	logger.Println(string(b))

	if slackChallengeRequest(w, r, b) {
		return
	}

	if slackEventCallbackRequest(w, r, b) {
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	return
}

func slackChallengeRequest(w http.ResponseWriter, r *http.Request, b []byte) bool {

	var req SlackChallengeRequest
	err := json.Unmarshal(b, &req)
	if err != nil {
		logger.Printf("%v\n", err)
		return false
	}

	if req.Type != slackEventTypeURLVerification {
		logger.Printf("type is not %v\n", slackEventTypeURLVerification)
		return false
	}

	res := SlackChallengeResponse{Challenge: req.Challenge}
	d, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Printf("%v\n", err)
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
		if err, ok := err.(*json.SyntaxError); ok {
			logger.Printf("%v\n", string(b[err.Offset-15:err.Offset+15]))
		}
		logger.Printf("%v\n", err)
		return false
	}

	if req.Type != slackEventTypeCallback {
		logger.Printf("type is not %v\n", slackEventTypeCallback)
		return false
	}

	if req.Event.Type != slackEventCallbackTypeMessage {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Printf("event.type is not %v\n", slackEventCallbackTypeMessage)
		logger.Printf("event.type: %v¥n", req.Event.Type)
		return true
	}

	logger.Printf("event.message: %v\n", req.Event.Message)

	if req.Event.Message == "ok" {
		lastStatus = "OK"
	}

	if req.Event.Message == "ng" {
		lastStatus = "NG"
	}

	w.WriteHeader(http.StatusOK)
	return true
}

// StatusHandler is
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	res := StatusResponse{LastStatus: lastStatus}
	d, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Printf("%v\n", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(d))
}

// ResponseLoggerMiddleware is
func ResponseLoggerMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("request method: %v\n", r.Method)
		logger.Printf("request requestURI: %v\n", r.RequestURI)
		logger.Printf("request url scheme: %v\n", r.URL.Scheme)
		logger.Printf("request url host: %v\n", r.URL.Host)
		logger.Printf("request url path: %v\n", r.URL.Path)
		logger.Printf("request url raw query: %v\n", r.URL.RawQuery)

		for k, v := range r.Header {
			logger.Printf("request header: %v: %v\n", k, v)
		}

		next.ServeHTTP(w, r)
	})
}

// DefaultHandler is
func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func main() {

	logger = log.New(os.Stdout, "", log.Lshortfile)
	port := os.Getenv("PORT")
	r := mux.NewRouter()
	r.Use(ResponseLoggerMiddleware)
	r.HandleFunc("/", DefaultHandler)
	r.HandleFunc("/slack", SlackHandler)
	r.HandleFunc("/status", StatusHandler)
	r.HandleFunc("/{wildcard}", DefaultHandler)
	http.Handle("/", r)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		logger.Printf("%v\n", err)
	}
}
