package executor

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type Distribution struct {
	Formula string  `yaml:"formula,omitempty"`
	Min     float64 `yaml:"min"`
	Max     float64 `yaml:"max"`
}

type ExecutorStep struct {
	Distribution string  `yaml:"distribution,omitempty"`
	Duration     float64 `yaml:"duration,omitempty"` // in sec
}

// ExecutionsPattern
const (
	Mixed  = "mixed"
	Linear = "linear"
)

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
			return fmt.Errorf("validation of distribution %s failed, %v", alias, dist)
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
		return fmt.Errorf("no formula provide")
	}
	if distribution.Min < 1 {
		return fmt.Errorf("invalid min value")
	}
	if distribution.Max <= distribution.Min {
		return fmt.Errorf("max must be >= than Min")
	}

	// Maybe do a formula check, but I might do that when I parse

	return nil
}
