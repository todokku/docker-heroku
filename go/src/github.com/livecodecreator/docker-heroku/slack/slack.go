package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/livecodecreator/docker-heroku/common"
)

const (
	slackEventTypeURLVerification = "url_verification"
	slackEventTypeCallback        = "event_callback"
	slackEventCallbackTypeMessage = "message"
	slackChatPostMessageEndpoint  = "https://slack.com/api/chat.postMessage"
)

type slackChallengeRequest struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

type slackChallengeResponse struct {
	Challenge string `json:"challenge"`
}

type slackEventCallbackRequest struct {
	Type  string                  `json:"type"`
	Event slackEventCallbackEvent `json:"event"`
}

type slackEventCallbackEvent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type slackLastRaspStatusRequest struct {
	Token   string `json:"token"`
	Channel string `json:channel`
	Text    string `json:text`
}

// SlackHandler is
func SlackHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%v\n", err)
		return
	}

	log.Println("request body:")
	log.Printf("%v\n", string(b))

	if slackChallengeRequestIfNeeded(w, r, b) {
		return
	}

	if slackEventCallbackRequestIfNeeded(w, r, b) {
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	return
}

func slackChallengeRequestIfNeeded(w http.ResponseWriter, r *http.Request, b []byte) bool {

	var req slackChallengeRequest
	err := json.Unmarshal(b, &req)
	if err != nil {
		log.Printf("%v\n", err)
		return false
	}

	if req.Type != slackEventTypeURLVerification {
		log.Printf("type is not %v\n", slackEventTypeURLVerification)
		return false
	}

	res := slackChallengeResponse{Challenge: req.Challenge}
	d, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v\n", err)
		return true
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(d))
	return true
}

func slackEventCallbackRequestIfNeeded(w http.ResponseWriter, r *http.Request, b []byte) bool {

	var req slackEventCallbackRequest
	err := json.Unmarshal(b, &req)
	if err != nil {
		log.Printf("%v\n", err)
		return false
	}

	if req.Type != slackEventTypeCallback {
		log.Printf("type is not %v\n", slackEventTypeCallback)
		return false
	}

	if req.Event.Type != slackEventCallbackTypeMessage {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("event.type is not %v\n", slackEventCallbackTypeMessage)
		log.Printf("event.type: %v¥n", req.Event.Type)
		return true
	}

	log.Printf("event.text: %v\n", req.Event.Text)

	// if req.Event.Text == "ok" {
	// 	common.LastStatus = "OK"
	// }

	// if req.Event.Text == "ng" {
	// 	common.LastStatus = "NG"
	// }

	if strings.Contains(req.Event.Text, "hello rasp") {
		log.Println("hello rasp matched!")
		postSlackLastRaspStatus(w, r)
	} else {
		log.Println("hello rasp mis matched!")
	}

	w.WriteHeader(http.StatusOK)
	return true
}

func postSlackLastRaspStatus(w http.ResponseWriter, r *http.Request) {

	d, err := json.Marshal(common.LastRaspStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v\n", err)
		return
	}

	req := slackLastRaspStatusRequest{
		Token:   os.Getenv("SLACK_TOKEN"),
		Channel: os.Getenv("SLACK_CHANNEL"),
		Text:    string(d),
	}

	p, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v\n", err)
		return
	}

	res, err := http.Post(slackChatPostMessageEndpoint, "application/json", bytes.NewReader(p))
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	log.Printf("slack api response: %+v\n", res)

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	log.Printf("slack api response body: %v\n", string(b))
}
