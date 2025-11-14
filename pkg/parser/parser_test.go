package parser

import (
	"os"
	"path"
	"testing"
)

func TestSplitYAML(t *testing.T) {
	p := path.Join("test", "01_test.yml")

	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("unexpected error reading the file: %s with err: %v", p, err)
	}

	ex, gen, hand, err := SplitYAML(data)
	if err != nil {
		t.Fatalf("unexpected error parsing and splitting the yaml: %v", err)
	}

	// Executor name
	if ex.Executor.Name != "daily_health_simulation" {
		t.Fatalf("ex.Executor: wrong name")
	}

	if ex.Executor.ExecutionPattern != "linear" {
		t.Fatalf("ex.Executor: wrong pattern")
	}

	if _, ok := ex.Distributions["sinusoidal"]; !ok {
		t.Fatalf("sinusoidal not in distributions")
	}
	if ex.Distributions["sinusoidal"].Formula != "./scripts/sinusoidal.sh" {
		t.Fatalf("ex.Distributions: wrong formula")
	}

	if len(gen.MeasurementTypes) != 1 {
		t.Fatalf("unexpected length of measurement-types")
	}

	if gen.MeasurementTypes["heart-rate"].Type != "HKQuantityTypeIdentifierHeartRate" {
		t.Fatalf("gen.Types.HR: wrong type")
	}

	if hand.Topic != "wearables-raw" {
		t.Fatalf("hand: wrong topic")
	}

}
