package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Server struct {
	api string
	url string
}

func (r Server) GetEntries() ([]string, error) {
	log.Printf("http.GetEntries")
	client := http.Client{Timeout: time.Second * 10}

	url := fmt.Sprintf("%s/api/queue?apikey=%s", r.url, r.api)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}

	entries := make([]Entry, 0)
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, err
	}

	completedTorrents := make([]string, 0)
	for _, entry := range entries {
		if entry.Status == "Completed" && entry.Protocol == "torrent" {
			completedTorrents = append(completedTorrents, entry.Title)
		}
	}

	return completedTorrents, nil
}

type Entry struct {
	Title    string `json:"title"`
	Status   string `json:"status"`
	Protocol string `json:"protocol"`
}
