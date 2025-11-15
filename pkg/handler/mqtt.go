package handler

import (
	"WearablesLoadGen/pkg/generator"
	"WearablesLoadGen/pkg/plotter"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

func init() {
	RegisterHandler("mqtt", NewMQTTHandlerFromYAML)
}

type MQTTConfig struct {
	Topic       string         `yaml:"topic,omitempty"`
	Broker      string         `yaml:"broker,omitempty"`
	DeviceCount int            `yaml:"device-count,omitempty"`
	PlotterType plotter.Config `yaml:"plotter,omitempty"`
}

type MQTTHandler struct {
	Base
	mqttClient       paho.Client
	conf             MQTTConfig
	measurementTypes []*generator.MeasurementBlueprint
}

// NewMQTTHandlerFromYAML will be used by the Factory to generate a MQTTHandler
func NewMQTTHandlerFromYAML(data []byte, measurementTypes []*generator.MeasurementBlueprint) (Handler, error) {
	// We need to create BaseHandler
	var config MQTTConfig

	if err := yaml.UnmarshalWithOptions(data, &config, yaml.Strict()); err != nil {
		return nil, err
	}

	var handler MQTTHandler

	handler.conf = config
	handler.measurementTypes = measurementTypes

	return &handler, nil
}

func (h *MQTTHandler) Init() error {
	// validate the Config
	if err := validateMQTTConfig(&h.conf); err != nil {
		log.Printf("validating mqtt config failed with error: %v", err)
		return err
	}

	// init mqtt client
	u := uuid.New().String()

	opts := paho.NewClientOptions()
	opts.AddBroker(h.conf.Broker)
	opts.SetClientID(u)
	opts.SetCleanSession(true)

	h.mqttClient = paho.NewClient(opts)

	// init base -> and plotter
	p, err := plotter.NewPlotterFromConfig(&h.conf.PlotterType)
	if err != nil {
		return err
	}
	h.Plotter = p

	return nil
}

// GenerateLoad gets only invoked once a second
func (h *MQTTHandler) GenerateLoad(requests int) error {
	wg := sync.WaitGroup{}
	wg.Add(requests)
	errChan := make(chan error, requests)

	// Idk if correct, but I do it like that
	err := h.Plotter.PlotLoadOverSeconds(time.Now(), requests)
	if err != nil {
		return err
	}

	// Plotting the Outbound Throughput
	sizeChan := make(chan int, requests) // track bytes sent

	go func() {
		for bytes := range sizeChan {
			err := h.Plotter.PlotOutboundThroughput(time.Now(), bytes)
			if err != nil {
				return
			}
		}
	}()

	for i := 0; i < requests; i++ {
		go func() {
			defer wg.Done()

			// Generating a proper payload
			measurementType := h.measurementTypes[rand.Intn(len(h.measurementTypes))]
			deviceId := fmt.Sprintf("test-device-%d", rand.Intn(h.conf.DeviceCount))
			payload, err := generator.GenerateMockPayload(measurementType, deviceId)
			if err != nil {
				log.Printf("not able to generate mock payload, aborting: %v", err)
				errChan <- err
				return
			}

			// convert Payload to JSON
			data, err := json.Marshal(payload)
			if err != nil {
				log.Printf("failed to marshal payload: %v", err)
				errChan <- err
				return
			}

			sizeChan <- len(data)

			token := h.mqttClient.Publish(h.conf.Topic, 0, false, data)
			if token.Error() != nil {
				log.Printf("error publishing, skipping: %v", token.Error())
			}

			errChan <- nil
		}()

	}

	wg.Wait()
	close(errChan)
	close(sizeChan)

	var allErrors []error
	for err := range errChan {
		if err != nil {
			allErrors = append(allErrors, err)
		}
	}
	if len(allErrors) > 0 {
		return fmt.Errorf("encountered %d errors, first: %v", len(allErrors), allErrors[0])
	}

	return nil
}

func validateMQTTConfig(cfg *MQTTConfig) error {

	if cfg.Topic == "" {
		return fmt.Errorf("must provide a topic for the handler")
	}

	if cfg.Broker == "" {
		return fmt.Errorf("must provide a broker URL")
	}

	pattern := `^(tcp|ssl|ws|wss)://[a-zA-Z0-9.-]+:1883$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(cfg.Broker) {
		return fmt.Errorf("broker URL must match: %s", pattern)
	}

	if cfg.DeviceCount <= 0 {
		return fmt.Errorf("device count must be > 0, is: %d", cfg.DeviceCount)
	}

	// validate plotter config
	if err := plotter.ValidatePlotter(&cfg.PlotterType); err != nil {
		return err
	}

	return nil
}
