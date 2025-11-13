package executer

type LoadPattern struct {
	Distributions map[string]Distribution `yaml:"distributions,omitempty"`
}

type Distribution struct {
	Formula string `yaml:"formula"`
}
