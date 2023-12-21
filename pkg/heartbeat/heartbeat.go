package heartbeat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type HeartbeatCore struct {
	Description   string `json:"description,omitempty"`
	Team          string `json:"ownerTeam,omitempty"`
	AlertMessage  string `json:"alertMessage,omitempty"`
	AlertTags     string `json:"alertTags,omitempty"`
	Interval      int    `json:"interval,omitempty"`
	IntervalUnit  string `json:"intervalUnit,omitempty""`
	AlertPriority string `json:"alertPriority,omitempty"`
	Enabled       bool   `json:"enabled,omitempty"`
}

type Heartbeat struct {
	HeartbeatCore
	Name string `json:"name,omitempty"`
}

type HeartbeatConfig struct {
	Heartbeat
	ApiKey  string
	BaseUrl string
}

var client = &http.Client{}

func New() HeartbeatConfig {
	heartbeat := HeartbeatConfig{}
	/*	if err := defaults.Set(heartbeat); err != nil {
		panic(err)
	}*/
	return heartbeat
}

func (h *HeartbeatConfig) CreateOrUpdate() {

	var request *http.Request
	exist := h.IsExist()
	createUpdate := "create"
	if !exist {
		log.WithField("HeartbeatName", h.Name).Info("Creating heartbeat!")
		b, err := json.Marshal(h.Heartbeat)
		if err != nil {
			log.WithField("Heartbeat", h.Heartbeat).Panicf("Could not marshall heartbeat!")
		}
		requestBody := bytes.NewBuffer(b)
		url := fmt.Sprintf("%s/v2/heartbeats", h.BaseUrl)
		request = mustCreateRequest("POST", url, requestBody)
	} else {
		log.WithField("HeartbeatName", h.Name).Info("Updating heartbeat!")
		b, err := json.Marshal(h.Heartbeat.HeartbeatCore)
		if err != nil {
			log.WithField("Heartbeat", h.Heartbeat).Panicf("Could not marshall heartbeat!")
		}
		requestBody := bytes.NewBuffer(b)
		url := fmt.Sprintf("%s/v2/heartbeats/%s", h.BaseUrl, h.Name)
		request = mustCreateRequest("PATCH", url, requestBody)
		createUpdate = "update"
	}

	request.Header.Add("Content-Type", "application/json")
	resp, err := doRequest(h.ApiKey, request)

	if err != nil {
		log.WithField("Heartbeat", h).WithError(err).Panicf("Could not %s heartbeat!", createUpdate)
	}
	defer resp.Body.Close()

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		log.WithField("StatusCode", resp.StatusCode).WithField("Status", resp.Status).WithField("Heartbeat", h).Panicf("Could not %s heartbeat!", createUpdate)
	}
}

func (h *HeartbeatConfig) IsExist() bool {
	url := fmt.Sprintf("%s/v2/heartbeats/%s", h.BaseUrl, h.Heartbeat.Name)
	request := mustCreateRequest("GET", url, nil)
	resp, err := doRequest(h.ApiKey, request)

	if err != nil {
		log.WithError(err).WithField("URL", request.URL.Path).Panic("Could not send request!")
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
		log.WithField("StatusCode", resp.StatusCode).WithField("Status", resp.Status).Error("Could not get heartbeat!")
		return false
	}

	return true
}

func (h *HeartbeatConfig) Ping() {
	url := fmt.Sprintf("%s/v2/heartbeats/%s/ping", h.BaseUrl, h.Name)

	request := mustCreateRequest("GET", url, nil)
	log.WithField("Heartbeat", h.Name).Info("Sending ping!")

	resp, err := doRequest(h.ApiKey, request)

	if err != nil {
		log.WithField("Heartbeat", h).WithError(err).Errorf("Could not send heartbeat!")
	}
	defer resp.Body.Close()

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
		log.WithField("StatusCode", resp.StatusCode).WithField("Status", resp.Status).WithField("Heartbeat", h).Errorf("Could not ping!")
	}
}

func (h *HeartbeatConfig) PingPeriodically() {

	var duration time.Duration
	switch h.IntervalUnit {
	case "minutes":
		duration = time.Minute
		break
	case "hours":
		duration = time.Hour
		break
	case "days":
		duration = time.Hour * 24
	default:
		log.WithField("Duration", h.IntervalUnit).Panic("Wrong IntervalUnit description!")
	}
	totalDuration := time.Duration(h.Interval) * duration

	// send ping before interval
	totalDuration -= time.Second * 30

	log.WithField(fmt.Sprintf("Interval(%s)", h.IntervalUnit), h.Interval).Warn("Setting up periodic heartbeat")

	ticker := time.NewTicker(totalDuration)

	for ; ; <-ticker.C {
		h.Ping()
	}
}

func mustCreateRequest(method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.WithError(err).Panic("Could not create request!")
	}
	return req
}
func doRequest(apiKey string, req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "GenieKey "+apiKey)
	resp, err := client.Do(req)
	return resp, err
}
