package parser

import (
	"encoding/json"
	"log"
	"strings"
)

type Message struct {
	UserID  string
	Topic   string
	Payload *Payload
}

type Payload struct {
	Type          string `json:"type"`
	SourceName    string `json:"sourceName"`
	SourceVersion string `json:"sourceVersion"`
	Unit          string `json:"unit"`
	Timestamp     string `json:"endDate"`
	Value         string `json:"value"`
	DeviceID      string `json:"device"`
}

type PayloadGroup struct {
	Count    int        `json:"count"`
	Payloads []*Payload `json:"records"`
}

//Maybe I need my own implementation of the Unmarshal function

func ParseSingleJSONMeasurement(data []byte) (*Payload, error) {

	payload := Payload{}

	err := json.Unmarshal(data, &payload)
	if err != nil {
		log.Printf("failed to parse single object into payload: %v", err)
		return nil, err
	}

	return &payload, nil
}

func ParseArrayJSONMeasurements(data []byte) ([]*Payload, error) {
	cleaned := strings.ReplaceAll(string(data), "NaN", "null")

	var groups map[string]PayloadGroup
	if err := json.Unmarshal([]byte(cleaned), &groups); err != nil {
		log.Printf("failed to parse nested structure: %v", err)
		return nil, err
	}

	var all []*Payload
	for _, group := range groups {
		for _, payload := range group.Payloads {
			// Handle missing or null device gracefully
			all = append(all, payload)
		}
	}
	return all, nil
}
