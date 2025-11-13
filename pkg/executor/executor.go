package executor

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
