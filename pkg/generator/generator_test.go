package generator

import (
	"os"
	"path"
	"testing"
)

func TestParseYAML(t *testing.T) {
	p := path.Join("config", "01_measurements.yml")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("unexpected failure reading file, %v", err)
	}

	cfg, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("unexpected failure parsing file: %v", err)
	}

	if len(cfg.MeasurementTypes) != 2 {
		t.Fatalf("config has not the correct size of test")
	}

	measurementType := cfg.MeasurementTypes["heart-rate"]
	if measurementType.Type != "HKQuantityTypeIdentifierHeartRate" {
		t.Fatalf("wrong type: %s", measurementType.Type)
	}
	if measurementType.SourceName != "ESP32-Wecker" {
		t.Fatalf("wrong sourceName: %s", measurementType.SourceName)
	}
	if measurementType.SourceVersion != "9.0" {
		t.Fatalf("wrong sourceVersion: %s", measurementType.SourceVersion)
	}
	if measurementType.Min != 100 {
		t.Fatalf("wrong min: %.2f", measurementType.Min)
	}

	if measurementType.Max != 200 {
		t.Fatalf("wrong max: %.2f", measurementType.Max)
	}

	if measurementType.Unit != "count/min" {
		t.Fatalf("wrong unit: %s", measurementType.Unit)
	}

	measurementType = cfg.MeasurementTypes["dietary-water"]
	if measurementType.Type != "HKQuantityTypeIdentifierDietaryWater" {
		t.Fatalf("wrong type: %s", measurementType.Type)
	}
	if measurementType.SourceName != "Lifesum" {
		t.Fatalf("wrong sourceName: %s", measurementType.SourceName)
	}
	if measurementType.SourceVersion != "15" {
		t.Fatalf("wrong sourceVersion: %s", measurementType.SourceVersion)
	}
	if measurementType.Min != 2000 {
		t.Fatalf("wrong min: %.2f", measurementType.Min)
	}

	if measurementType.Max != 5000 {
		t.Fatalf("wrong max: %.2f", measurementType.Max)
	}

	if measurementType.Unit != "mL" {
		t.Fatalf("wrong unit: %s", measurementType.Unit)
	}
}

func TestGenerateMockPayload(t *testing.T) {

}
