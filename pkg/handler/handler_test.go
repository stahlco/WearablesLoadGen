package handler

import (
	"WearablesLoadGen/pkg/generator"
	"os"
	"path"
	"testing"
)

func TestGenerateHandlerFromYAML(t *testing.T) {
	p := path.Join("test", "01_handler.yml")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("unexpected error reading from handler file")
	}

	// Measurements should work
	measurementPath := path.Join("..", "generator", "test", "01_measurements.yml")
	measurementData, err := os.ReadFile(measurementPath)
	if err != nil {
		t.Fatalf("unexpected error reading from measurement file")
	}
	c, err := generator.ParseYAML(measurementData)
	if err != nil || c == nil {
		t.Fatalf("unexprected error parsing measurement types")
	}
	handler, err := GenerateHandlerFromYAML(data, c.GetAllMeasurementBlueprints())
	if err != nil {
		t.Fatalf("unexpected error reading from Handler: %v", err)
	}

	if handler.Topic != "wearables-raw" {
		t.Fatalf("wrong topic")
	}

	if handler.BrokerURL != "tcp://localhost:1883" {
		t.Fatalf("wrong broker url")
	}

	if handler.DeviceCount != 10000 {
		t.Fatalf("wrong deviceCount")
	}

	if handler.CsvFilePath != "./test/plot/plot.csv" {
		t.Fatalf("FilePath should be: \"./plot/plot.csv\", but is: \"%s\"", handler.CsvFilePath)
	}
}
