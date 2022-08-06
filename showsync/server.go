package showsync

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetCompletedTorrents(u *url.URL) ([]string, error) {
	client := http.Client{Timeout: time.Second * 10}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
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

	entries, err := parseServerOutput(body)
	if err != nil {
		return nil, err
	}

	completedTorrentsMap := make(map[string]struct{})
	for _, entry := range entries {
		if strings.ToLower(entry.Status) == "completed" && entry.Protocol == "torrent" {
			completedTorrentsMap[entry.Title] = struct{}{}
		}
	}

	completedTorrents := make([]string, 0)
	for k, _ := range completedTorrentsMap {
		completedTorrents = append(completedTorrents, k)
	}

	return completedTorrents, nil
}

func parseServerOutput(data []byte) ([]Entry, error) {
	entries := make([]Entry, 0)
	if data[0] == '[' {
		if err := json.Unmarshal(data, &entries); err != nil {
			return nil, err
		}
	} else {
		radarr := Radarr{}
		if err := json.Unmarshal(data, &radarr); err != nil {
			return nil, err
		}
		entries = radarr.Records
	}
	return entries, nil
}

type Entry struct {
	Title    string `json:"title"`
	Status   string `json:"status"`
	Protocol string `json:"protocol"`
}

type Radarr struct {
	Records []Entry `json:"records"`
}
