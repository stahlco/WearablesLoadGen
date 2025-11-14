package executor

import (
	"os"
	"path"
	"testing"
)

func TestParseExecutionConfigYAML(t *testing.T) {
	p := path.Join("test", "01_executor.yml")
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
	if equalDist.Formula != "./scripts/equal_distribution.sh" {
		t.Fatalf("unexpected formula for equal_distribution: got %q", equalDist.Formula)
	}
	if equalDist.Base != 800 {
		t.Fatalf("unexpected min for equal_distribution: got %f", equalDist.Base)
	}
	if equalDist.Amp != 200 {
		t.Fatalf("unexpected amp for equal_distribution: got %f", equalDist.Amp)
	}

	highPeek, ok := execCfg.Distributions["sinusoidal"]
	if !ok {
		t.Fatalf("missing 'sinusoidal' in distributions")
	}
	if highPeek.Formula != "./scripts/sinusoidal.sh" {
		t.Fatalf("unexpected formula for high_peek: got %q", highPeek.Formula)
	}
	if highPeek.Base != 2000 {
		t.Fatalf("unexpected min for high_peek: got %f", highPeek.Base)
	}
	if highPeek.Amp != 500 {
		t.Fatalf("unexpected max for high_peek: got %f", highPeek.Amp)
	}

	executor := execCfg.Executor
	if executor.Name != "daily_health_simulation" {
		t.Fatalf("unexpected executor name: got %q", executor.Name)
	}
	if executor.ExecutionPattern != "linear" {
		t.Fatalf("unexpected execution pattern: got %q", executor.ExecutionPattern)
	}
	if executor.Duration != 500 {
		t.Fatalf("unexpected executor duration: got %f", executor.Duration)
	}

	if len(executor.ExecutionSteps) != 2 {
		t.Fatalf("expected 2 execution steps, got %d", len(executor.ExecutionSteps))
	}

	step1 := executor.ExecutionSteps[0]
	if step1.Distribution != "equal_distribution" {
		t.Fatalf("unexpected first step distribution: got %q", step1.Distribution)
	}
	if step1.Duration != 100 {
		t.Fatalf("unexpected first step duration: got %f", step1.Duration)
	}

	step2 := executor.ExecutionSteps[1]
	if step2.Distribution != "sinusoidal" {
		t.Fatalf("unexpected second step distribution: got %q", step2.Distribution)
	}
	if step2.Duration != 400 {
		t.Fatalf("unexpected second step duration: got %f", step2.Duration)
	}
}

func TestEvaluateDistribution(t *testing.T) {
	dist := Distribution{
		Formula: "./scripts/sinusoidal.sh",
		Base:    2000,
		Amp:     500,
	}

	res := EvaluateDistribution(dist, 100)
	if res != 2321 {
		t.Fatalf("wrong result should be 2349, but is: %d", res)
	}

}
