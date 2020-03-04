package raspberrypi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/livecodecreator/docker-heroku/common"
)

type postStatusRequest struct {
	CPU      string `json:"cpu"`
	Disk     string `json:"disk"`
	Memory   string `json:"memory"`
	BootTime string `json:"bootTime"`
}

// PostStatusHandler is
func PostStatusHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%v\n", err)
		return
	}

	var req postStatusRequest
	err = json.Unmarshal(b, &req)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	common.LastRaspStatus.CPU = req.CPU
	common.LastRaspStatus.Disk = req.Disk
	common.LastRaspStatus.Memory = req.Memory
	common.LastRaspStatus.BootTime = req.BootTime
	common.LastRaspStatus.Timestamp = time.Now()
}
