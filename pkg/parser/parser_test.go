package parser

import (
	"log"
	"os"
	"path"
	"testing"
)

func TestParseSingleJSONMeasurement_Success(t *testing.T) {
	p := path.Join("examples", "01_heart_rate.json")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("unexpected error reading from file: %v", err)
	}
	log.Print(string(data))

	// Parsing to JSON
	payload, err := ParseSingleJSONMeasurement(data)
	if err != nil {
		t.Fatalf("unexprected error reading from file: %v", err)
	}
	if payload.Type != "HKQuantityTypeIdentifierHeartRate" {
		t.Fatalf("type of payload does not match: HKQuantityTypeIdentifierHeartRate, is: [%s]", payload.Type)
	}
	if payload.SourceName != "ESP32-Wecker" {
		t.Fatalf("sourceName of payload does not match: ESP32-Wecker, is [%s]", payload.SourceName)
	}
	if payload.SourceVersion != "9.0" {
		t.Fatalf("sourceVersion of payload does not match: 9.0 , is [%s]", payload.SourceVersion)
	}
	if payload.Unit != "count/min" {
		t.Fatalf("wrong unit: %s", payload.Unit)
	}
	if payload.Timestamp != "2022-09-17 16:01:03" {
		t.Fatalf("wrong timestamp: %s", payload.Timestamp)
	}
	if payload.DeviceID != "<<HKDevice: 0x7b3ae6ac0>, name:Apple Watch, manufacturer:Apple Inc., model:Watch, hardware:Watch6,15, software:9.0, creation date:2022-08-27 14:08:39 +0000>" {
		t.Fatalf("wrong device, got: %s", payload.DeviceID)
	}
}

func TestParseArrayJSONMeasurements_Success(t *testing.T) {
	p := path.Join("examples", "02_dietary_water.json")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("unexpected error while reading from file %s: %v", p, err)
	}

	log.Print(string(data))
	payloads, err := ParseArrayJSONMeasurements(data)
	if err != nil {
		t.Fatalf("unexprected error parsing data: %v", err)
	}

	if len(payloads) != 2 {
		t.Fatalf("wrong len(payloads), 2 is: [%d]", len(payloads))
	}
	payload := payloads[0]

	if payload.Type != "HKQuantityTypeIdentifierDietaryWater" {
		t.Fatalf("type of payload does not match: HKQuantityTypeIdentifierDietaryWater, is: [%s]", payload.Type)
	}
	if payload.SourceName != "Lifesum" {
		t.Fatalf("sourceName of payload does not match: Lifesum, is [%s]", payload.SourceName)
	}
	if payload.SourceVersion != "15" {
		t.Fatalf("sourceVersion of payload does not match: 15 , is [%s]", payload.SourceVersion)
	}
	if payload.Unit != "mL" {
		t.Fatalf("wrong unit: %s", payload.Unit)
	}
	if payload.Timestamp != "2023-04-17 18:29:08" {
		t.Fatalf("wrong timestamp: %s", payload.Timestamp)
	}
	if payload.Value != "4000" {
		t.Fatalf("wrong value: %s", payload.Value)
	}
	if payload.DeviceID != "" {
		t.Fatalf("wrong device, got: %s", payload.DeviceID)
	}

	payload = payloads[1]

	if payload.Type != "HKQuantityTypeIdentifierDietaryWater" {
		t.Fatalf("type of payload does not match: HKQuantityTypeIdentifierDietaryWater, is: [%s]", payload.Type)
	}
	if payload.SourceName != "Lifesum" {
		t.Fatalf("sourceName of payload does not match: Lifesum, is [%s]", payload.SourceName)
	}
	if payload.SourceVersion != "15" {
		t.Fatalf("sourceVersion of payload does not match: 15 , is [%s]", payload.SourceVersion)
	}
	if payload.Unit != "mL" {
		t.Fatalf("wrong unit: %s", payload.Unit)
	}
	if payload.Timestamp != "2023-04-18 20:55:42" {
		t.Fatalf("wrong timestamp: %s", payload.Timestamp)
	}
	if payload.Value != "2500" {
		t.Fatalf("wrong value: %s", payload.Value)
	}
	if payload.DeviceID != "" {
		t.Fatalf("wrong device, got: %s", payload.DeviceID)
	}

}
