package parser

import (
	"WearablesLoadGen/pkg/executor"
	"WearablesLoadGen/pkg/generator"
	"WearablesLoadGen/pkg/handler"
	"bytes"
	"fmt"
	"log"
	"strings"
)

func SplitYAML(data []byte) (*executor.ExecutionConfig, *generator.GeneratorConfig, *handler.MQTTHandler2, error) {

	docs := bytes.Split(data, []byte("\n---"))

	var execDoc, genDoc, handlerDoc []byte

	for _, doc := range docs {
		lines := bytes.Split(doc, []byte("\n"))
		if len(lines) == 0 {
			continue
		}

		var kind string
		var startIndex int

		for i, line := range lines {
			lineStr := strings.TrimSpace(string(line))
			if strings.HasPrefix(lineStr, "kind:") {
				parts := strings.SplitN(lineStr, ":", 2)
				if len(parts) == 2 {
					kind = strings.TrimSpace(parts[1])
				}
				startIndex = i + 1
				break
			}
		}

		if kind == "" {
			return nil, nil, nil, fmt.Errorf("document missing 'kind' field")
		}

		body := bytes.Join(lines[startIndex:], []byte("\n"))

		switch strings.ToLower(kind) {
		case "executor":
			execDoc = body
		case "generator":
			genDoc = body
		case "handler":
			handlerDoc = body
		default:
			return nil, nil, nil, fmt.Errorf("unknown kind: %s", kind)
		}
	}

	// Forward the docs
	execConfig, err := executor.ParseExecutionConfigYAML(execDoc)
	if err != nil {
		log.Printf("failed to parse execution test")
		return nil, nil, nil, err
	}

	genConfig, err := generator.ParseYAML(genDoc)
	if err != nil {
		log.Printf("failed to parse generator test")
		return nil, nil, nil, err
	}

	handlerConfig, err := handler.GenerateMQTTHandlerFromYAML(handlerDoc, genConfig.GetAllMeasurementBlueprints())
	if err != nil {
		log.Printf("failed to parse execution test")
		return nil, nil, nil, err
	}

	return execConfig, genConfig, handlerConfig, nil
}
