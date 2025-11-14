package handler

import (
	"WearablesLoadGen/pkg/generator"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strconv"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

type Handler interface {
	GenerateLoad(requests int) error
	logToCSV(requests int) error
}

type Config struct {
	Handlers map[string]*MQTTHandler `yaml:"handlers,omitempty"`
}

type MQTTHandler struct {
	mqttClient paho.Client
	// Those are used to create the messages
	Topic            string `yaml:"topic,omitempty"`
	BrokerURL        string `yaml:"broker,omitempty"`
	DeviceCount      int    `yaml:"device-count,omitempty"`
	measurementTypes []*generator.MeasurementBlueprint

	// For Plotting
	CsvFilePath string `yaml:"plot-path,omitempty"`
	csvFile     *os.File
	csvWriter   *csv.Writer
	callCount   int
	mu          sync.Mutex
}

func GenerateHandlerFromYAML(data []byte, measurementTypes []*generator.MeasurementBlueprint) (*MQTTHandler, error) {
	var handler *MQTTHandler
	var config Config

	if err := yaml.UnmarshalWithOptions(data, &config, yaml.Strict()); err != nil {
		return nil, err
	}

	for alias, h := range config.Handlers {
		if alias == "mqtt" {
			handler = h
		}
	}

	if handler == nil {
		return nil, fmt.Errorf("not able to find a supported type in the config: Support Types: [mqtt]")
	}

	if err := validateMQTTHandler(handler); err != nil {
		return nil, err
	}

	u := uuid.New().String() //clientID

	handler.measurementTypes = measurementTypes

	opts := paho.NewClientOptions()
	opts.AddBroker(handler.BrokerURL)
	opts.SetClientID(u)
	opts.SetCleanSession(true)

	client := paho.NewClient(opts)

	handler.mqttClient = client

	file, err := os.OpenFile(handler.CsvFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	handler.csvFile = file
	handler.csvWriter = csv.NewWriter(file)
	handler.callCount = 0

	return handler, nil
}

func validateMQTTHandler(handler *MQTTHandler) error {
	// we then already have: topic, broker, device-count, plot-path
	if handler.Topic == "" {
		return fmt.Errorf("must provide a topic for the handler")
	}

	if handler.BrokerURL == "" {
		return fmt.Errorf("must provide a broker URL")
	}

	pattern := `^(tcp|ssl|ws|wss)://[a-zA-Z0-9.-]+:1883$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(handler.BrokerURL) {
		return fmt.Errorf("broker URL must match: %s", pattern)
	}

	if handler.CsvFilePath == "" {
		return fmt.Errorf("plot file path must at least provide a directory")
	}

	dir, file := path.Split(handler.CsvFilePath)
	if file == "" {
		handler.CsvFilePath = dir + "plot.csv"
	}

	return nil
}

func (h *MQTTHandler) GenerateLoad(requests int) error {
	wg := sync.WaitGroup{}
	wg.Add(requests)
	errChan := make(chan error, requests)

	for i := 0; i < requests; i++ {
		go func() {
			defer wg.Done()

			measurementType := h.measurementTypes[rand.Intn(len(h.measurementTypes))] // one of the existing
			device_ID := fmt.Sprintf("test-device-%d", rand.Intn(h.DeviceCount))
			payload, err := generator.GenerateMockPayload(measurementType, device_ID)
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

			token := h.mqttClient.Publish(h.Topic, 0, false, data)
			token.Wait()
			if token.Error() != nil {
				log.Printf("error publishing, skipping: %v", token.Error())
			}

			errChan <- nil
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *MQTTHandler) logToCSV(requests int) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.callCount++
	row := []string{
		strconv.Itoa(h.callCount),
		strconv.Itoa(requests),
		time.Now().Format(time.RFC3339),
	}

	if err := h.csvWriter.Write(row); err != nil {
		return err
	}

	return nil
}
