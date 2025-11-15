package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	paho "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	log.SetPrefix("loadgen: ")
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	paho.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	paho.CRITICAL = log.New(os.Stdout, "[CRITICAL] ", 0)
	paho.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	paho.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	// Support both -c and --config
	configPathShort := flag.String("c", "", "Path to config file")
	configPathLong := flag.String("config", "", "Path to config file")

	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	configPath := *configPathLong
	if configPath == "" {
		configPath = *configPathShort
	}

	if configPath == "" {
		fmt.Println("Error: no config file provided.\nUsage:")
		fmt.Println("  -c <path> or --config <path>")
		os.Exit(1)
	}

	// GeneratorConfig

}
