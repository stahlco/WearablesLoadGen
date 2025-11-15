package handler

import (
	"WearablesLoadGen/pkg/generator"
	"fmt"
	"log"

	yaml "gopkg.in/yaml.v3"
)

type Handler interface {
	Init() error
	GenerateLoad(requests int) error
	logToCSV(requests int) error
	Close() error
}

type GenericHandlerYAML struct {
	Handler struct {
		Type   string    `yaml:"type"`
		Config yaml.Node `yaml:"config"`
	} `yaml:"handler"`
}

type Factory func(data []byte, measurementTypes []*generator.MeasurementBlueprint) (Handler, error)

var handlerRegistry = map[string]Factory{}

func RegisterHandler(name string, factory Factory) {
	handlerRegistry[name] = factory
}

func CreateHandler(name string, data []byte, measurementTypes []*generator.MeasurementBlueprint) (Handler, error) {
	factory, ok := handlerRegistry[name]
	if !ok {
		return nil, fmt.Errorf("unknown handler type: %s, supported types: %v", name, handlerRegistry)
	}

	return factory(data, measurementTypes)
}

// GenerateHandlerFromYAML creates a generic handler
func GenerateHandlerFromYAML(data []byte, measurementTypes []*generator.MeasurementBlueprint) (Handler, error) {
	var parsed GenericHandlerYAML

	if err := yaml.Unmarshal(data, &parsed); err != nil {
		log.Printf("parsing the yaml failed with err: %v", err)
		return nil, err
	}

	cfgBytes, err := yaml.Marshal(parsed.Handler.Config)
	if err != nil {
		log.Printf("marshalling the config failed with error nil")
		return nil, err
	}

	handler, err := CreateHandler(parsed.Handler.Type, cfgBytes, measurementTypes)
	if err != nil {
		return nil, err
	}

	return handler, handler.Init()
}
