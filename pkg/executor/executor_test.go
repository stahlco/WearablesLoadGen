package executor

import (
	"os"
	"path"
	"testing"
)

func TestParseExecutionConfigYAML(t *testing.T) {
	p := path.Join("config", "01_executor.yml")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("unexpected error while reading file %s: %v", p, err)
	}

	execCfg, err := ParseExecutionConfigYAML(data)
	if err != nil {
		t.Fatalf("unexpected error while parsing yaml: %v", err)
	}

	if execCfg == nil {
		t.Fatalf("expected non-nil ExecutionConfig, got nil")
	}

	if len(execCfg.Distributions) != 2 {
		t.Fatalf("expected 2 distributions, got %d", len(execCfg.Distributions))
	}

	equalDist, ok := execCfg.Distributions["equal_distribution"]
	if !ok {
		t.Fatalf("missing 'equal_distribution' in distributions")
	}
	if equalDist.Formula != "min + rand*(max-min)" {
		t.Fatalf("unexpected formula for equal_distribution: got %q", equalDist.Formula)
	}
	if equalDist.Min != 800 {
		t.Fatalf("unexpected min for equal_distribution: got %f", equalDist.Min)
	}
	if equalDist.Max != 1000 {
		t.Fatalf("unexpected max for equal_distribution: got %f", equalDist.Max)
	}

	highPeek, ok := execCfg.Distributions["high_peek"]
	if !ok {
		t.Fatalf("missing 'high_peek' in distributions")
	}
	if highPeek.Formula != "min + (rand^2)*(max-min)" {
		t.Fatalf("unexpected formula for high_peek: got %q", highPeek.Formula)
	}
	if highPeek.Min != 1500 {
		t.Fatalf("unexpected min for high_peek: got %f", highPeek.Min)
	}
	if highPeek.Max != 2000 {
		t.Fatalf("unexpected max for high_peek: got %f", highPeek.Max)
	}

	executor := execCfg.Executor
	if executor.Name != "daily_health_simulation" {
		t.Fatalf("unexpected executor name: got %q", executor.Name)
	}
	if executor.ExecutionPattern != "mixed" {
		t.Fatalf("unexpected execution pattern: got %q", executor.ExecutionPattern)
	}
	if executor.Duration != 250 {
		t.Fatalf("unexpected executor duration: got %f", executor.Duration)
	}

	if len(executor.ExecutionSteps) != 2 {
		t.Fatalf("expected 2 execution steps, got %d", len(executor.ExecutionSteps))
	}

	step1 := executor.ExecutionSteps[0]
	if step1.Distribution != "equal_distribution" {
		t.Fatalf("unexpected first step distribution: got %q", step1.Distribution)
	}
	if step1.Duration != 120 {
		t.Fatalf("unexpected first step duration: got %f", step1.Duration)
	}

	step2 := executor.ExecutionSteps[1]
	if step2.Distribution != "high_peek" {
		t.Fatalf("unexpected second step distribution: got %q", step2.Distribution)
	}
	if step2.Duration != 5 {
		t.Fatalf("unexpected second step duration: got %f", step2.Duration)
	}
}
