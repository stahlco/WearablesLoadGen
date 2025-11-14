package executor

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
)

// ExecutionsPattern
const (
	Mixed  = "mixed"
	Linear = "linear"
)

type Distribution struct {
	Formula string  `yaml:"formula,omitempty"`
	Base    float64 `yaml:"base"`
	Amp     float64 `yaml:"amp"`
}

type ExecutorStep struct {
	Distribution string  `yaml:"distribution,omitempty"`
	Duration     float64 `yaml:"duration,omitempty"` // in sec
}

type ExecutorConfig struct {
	Name             string         `yaml:"name,omitempty"`
	ExecutionPattern string         `yaml:"execution-pattern,omitempty"`
	Duration         float64        `yaml:"duration,omitempty"` // in sec
	ExecutionSteps   []ExecutorStep `yaml:"steps,omitempty"`
}

type ExecutionConfig struct {
	Distributions map[string]Distribution `yaml:"distributions,omitempty"`
	Executor      ExecutorConfig          `yaml:"executor,omitempty"`
}

func ParseExecutionConfigYAML(data []byte) (*ExecutionConfig, error) {
	var config ExecutionConfig

	if err := yaml.UnmarshalWithOptions(data, &config, yaml.Strict()); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(cfg *ExecutionConfig) error {
	distributions := cfg.Distributions
	executor := cfg.Executor

	if len(distributions) == 0 {
		return fmt.Errorf("no distributions provided")
	}
	for alias, dist := range distributions {
		err := validateDistributions(dist)
		if err != nil {
			return fmt.Errorf("validation of distribution %s failed, %v: %v", alias, dist, err)
		}
	}

	if err := validateExecutor(executor, distributions); err != nil {
		return fmt.Errorf("validation executor failed, %v", err)
	}

	return nil
}

func validateExecutor(executor ExecutorConfig, distributions map[string]Distribution) error {
	if executor.Name == "" {
		return fmt.Errorf("no name for executor provided")
	}
	if executor.Duration < 1 {
		return fmt.Errorf("execution duration is eiter negative or too small")
	}
	if executor.ExecutionPattern != Mixed && executor.ExecutionPattern != Linear {
		return fmt.Errorf("execution pattern can either be \"mixed\" or \"linear\"")
	}
	if len(executor.ExecutionSteps) == 0 {
		return fmt.Errorf("must provide minimum a step")
	}

	for _, step := range executor.ExecutionSteps {
		if err := validateExecutionStep(step, distributions); err != nil {
			return err
		}
	}

	return nil
}

func validateExecutionStep(step ExecutorStep, distributions map[string]Distribution) error {
	if step.Distribution == "" {
		return fmt.Errorf("missing distribution in execution step")
	}

	if _, ok := distributions[step.Distribution]; !ok {
		return fmt.Errorf("unknown distribution '%s'", step.Distribution)
	}

	if step.Duration <= 0 {
		return fmt.Errorf("invalid duration %.2f — must be > 0", step.Duration)
	}

	return nil
}

func validateDistributions(distribution Distribution) error {
	if distribution.Formula == "" {
		return fmt.Errorf("no script provide")
	}
	if distribution.Base < 1 {
		return fmt.Errorf("invalid min value")
	}
	if distribution.Amp >= distribution.Base {
		return fmt.Errorf("amp must be <= than base, but is: Amp: %.2f, Base: %.2f", distribution.Amp, distribution.Base)
	}

	// Maybe do a formula check, but I might do that when I parse

	return nil
}

func EvaluateDistribution(d Distribution, t int) int {
	if d.Formula == "" {
		log.Printf("No formula/script defined, returning base value (fallback is equal-dist)")
		return int(d.Base)
	}

	cmd := exec.Command(
		d.Formula,
		fmt.Sprintf("--base=%d", int(d.Base)),
		fmt.Sprintf("--amp=%d", int(d.Amp)),
		fmt.Sprintf("--t=%d", t),
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running script '%s': %v", d.Formula, err)
		return int(d.Base)
	}

	outputStr := strings.TrimSpace(out.String())
	val, err := strconv.ParseFloat(outputStr, 64)
	if err != nil {
		log.Printf("Error parsing script output '%s': %v", outputStr, err)
		return int(d.Base)
	}

	return int(val)
}
