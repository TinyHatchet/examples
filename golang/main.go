package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	tinyhatchet "github.com/tinyhatchet/go-tinyhatchet"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags"`
}

func main() {
	// logs can be sent by hand
	sendManual()

	// or using the go-tinyhatchet library
	sendLibrary()
}

func doAThing() error {
	log.Println("did a thing")
	return errors.New("kaboom")
}

func sendManual() {
	err := doAThing()

	entry := LogEntry{
		Timestamp: time.Now(),
		Text:      err.Error(),
		Tags:      []string{"tinyhatchet", "example"},
	}
	b, _ := json.Marshal(entry)
	buffer := bytes.NewBuffer(b)
	req, _ := http.NewRequest(http.MethodPost, "http://tinyhatchet.com/injest.json", buffer)
	req.SetBasicAuth("apiTokenID", "apiTokenSecret")
	http.DefaultClient.Do(req)
}

func sendLibrary() {
	err := doAThing()
	log := tinyhatchet.Logger{
		APIToken:    "ApiTokenID",
		APISecret:   "ApiTokenSecret",
		DefaultTags: []string{"tinyhatchet", "example"},
		AutoTagger: func(tags []string, arg interface{}) []string {
			switch arg := arg.(type) {
			case error:
				// This will add to the default tags if the argument to the log commands is sql.ErrNoRows
				if errors.Is(arg, sql.ErrNoRows) {
					tags = append(tags, "database", "not found")
				}
			}
			return tags
		},
	}
	log.Println(err)
}
