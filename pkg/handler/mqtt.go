package handler

import (
	"WearablesLoadGen/pkg/generator"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHandler struct {
	mqttClient paho.Client
	// Those are used to create the messages
	topic            string
	measurementTypes []*generator.MeasurementBlueprint
	deviceCount      int
}

func NewMQTTHandler(
	topic string,
	measurementTypes []*generator.MeasurementBlueprint,
	deviceCount int,
	brokerURL string,
	clientID string,
) (*MQTTHandler, error) {

	opts := paho.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)

	client := paho.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect MQTT client: %w", token.Error())
	}

	return &MQTTHandler{
		mqttClient:       client,
		topic:            topic,
		measurementTypes: measurementTypes,
		deviceCount:      deviceCount,
	}, nil
}

func (h *MQTTHandler) GenerateLoad(requets int) error {
	wg := sync.WaitGroup{}
	wg.Add(requets)
	errChan := make(chan error, requets)

	for i := 0; i < requets; i++ {
		go func() {
			defer wg.Done()

			measurementType := h.measurementTypes[rand.Intn(len(h.measurementTypes))] // one of the existing
			device_ID := fmt.Sprintf("test-device-%d", rand.Intn(h.deviceCount))
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

			token := h.mqttClient.Publish(h.topic, 0, false, data)
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
