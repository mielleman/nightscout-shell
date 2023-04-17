package nightscout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Status struct {
	Status   string   `json:"status"`
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Settings Settings `json:"settings"`
}

type Settings struct {
	Units      string     `json:"units"`
	TimeFormat int        `json:"timeFormat"`
	Thresholds Thresholds `json:"thresholds"`
}

type Thresholds struct {
	High         int `json:"bgHigh"`
	TargetTop    int `json:"bgTargetTop"`
	TargetBottom int `json:"bgTargetBottom"`
	Low          int `json:"bgLow"`
}

type Entry struct {
	Value     int    `json:"sgv"`
	Date      int    `json:"date"`
	Direction string `json:"direction"`
}

type Nightscout struct {
	url    string
	token  string
	client *http.Client
	Status *Status
}

func New(url string, token string) *Nightscout {
	return &Nightscout{
		url:    url,
		token:  token,
		client: &http.Client{},
	}
}

func (n *Nightscout) GetStatus() error {
	url := fmt.Sprintf("%s/api/v1/status?token=%s", n.url, n.token)
	log.Info(url)

	// Make the request
	data, err := n.get(url)
	if err != nil {
		return err
	}

	// Map the json (values, as an array) to the struct
	var status Status
	json.Unmarshal(data, &status)

	// Set the latest status
	n.Status = &status

	return nil
}

func (n *Nightscout) GetLastEntry() (*Entry, error) {
	url := fmt.Sprintf("%s/api/v1/entries?count=1&token=%s", n.url, n.token)
	log.Info(url)

	// Make the request
	data, err := n.get(url)
	if err != nil {
		return &Entry{}, err
	}

	// Map the json (values, as an array) to the struct
	var entries []Entry
	json.Unmarshal(data, &entries)

	if len(entries) != 1 {
		log.Error("Unknown response of entries received when retrieving latest entry from Nightscout")
		return &Entry{}, fmt.Errorf("received %d entries, expected 1", len(entries))
	}

	// return the first result (should be only one)
	return &entries[0], nil
}

func (n *Nightscout) get(url string) ([]byte, error) {
	// Make the request
	var data []byte
	request, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		log.Error("Failed to create a new HTTP request")
		return []byte{}, err
	}

	// Set the default headers
	request.Header.Set("accept", "application/json")

	// Do the request
	response, err := n.client.Do(request)
	if err != nil {
		log.Error("Failed to retrieve data from Nightscout")
		return []byte{}, err
	}

	// Read the body
	body, bodyErr := io.ReadAll(response.Body)

	// Check that we get a code 200
	if response.StatusCode != 200 {
		log.WithField("status", response.StatusCode).Error("Unknown status code received when retrieving data from Nightscout")
		return []byte{}, fmt.Errorf("unknown status code: %d", response.StatusCode)
	}

	// get the body and unpack it
	if bodyErr != nil {
		log.WithField("body", body).Error("Unknown response body received when retrieving data from Nightscout")
		return []byte{}, bodyErr
	}

	return body, nil
}
