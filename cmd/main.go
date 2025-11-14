package main

import (
	"log"
	"os"

	paho "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	log.SetPrefix("loadgen: ")
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	log.Print("Hello World!")

	paho.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	paho.CRITICAL = log.New(os.Stdout, "[CRITICAL] ", 0)
	paho.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	paho.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	// GeneratorConfig

}
