package plotter

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Plotter interface {
	PlotLoadOverTime(t time.Time, load int) error
	PlotOutboundThroughput(t time.Time, bytes int) error
	PlotMessageSizePerMessage(message int, size int) error
	PlotMessageSizeDistribution(bucket string, count int) error
	writeRow(plotType string, row []string) error
}

type Config struct {
	Type string `yaml:"type,omitempty"`
	Path string `yaml:"path,omitempty"`
}

type CsvPlotter struct {
	ID         string
	CsvDir     string
	csvWriters map[string]*csv.Writer
}

func NewPlotterFromConfig(cfg *Config) (Plotter, error) {
	switch strings.ToLower(cfg.Type) {
	case "csv":
		plotter, err := newCsvPlotter(cfg)
		if err != nil {
			return nil, err
		}
		return plotter, nil
	default:
		return nil, fmt.Errorf("%s is not supported", cfg.Type)
	}
}

func (c CsvPlotter) PlotLoadOverTime(t time.Time, load int) error {
	return nil
}

func (c CsvPlotter) PlotOutboundThroughput(t time.Time, bytes int) error {
	//TODO implement me
	panic("implement me")
}

func (c CsvPlotter) PlotMessageSizePerMessage(message int, size int) error {
	//TODO implement me
	panic("implement me")
}

func (c CsvPlotter) PlotMessageSizeDistribution(bucket string, count int) error {
	//TODO implement me
	panic("implement me")
}

func (c CsvPlotter) writeRow(plotType string, row []string) error {
	return nil
}

func newCsvPlotter(cfg *Config) (Plotter, error) {
	var plotter CsvPlotter

	if _, err := os.ReadDir(cfg.Path); err != nil {
		err = os.MkdirAll(cfg.Path, 0777)
		if err != nil {
			return nil, err
		}
		plotter.CsvDir = cfg.Path
	}

	plotter.ID = uuid.New().String()

	// Create files for all different plots that will be created
	// <plotter-id>_<plot-type>.csv
	plotter.csvWriters = plotter.createPlotWriters()

	// Create csv-Writer

	return nil, nil
}

func (p CsvPlotter) createPlotWriters() map[string]*csv.Writer {
	files := map[string]string{
		"load_over_time":         fmt.Sprintf("%s-load_over_time.csv", p.ID),
		"outbound_throughput":    fmt.Sprintf("%s-outbound_throughput.csv", p.ID),
		"message_size_over_time": fmt.Sprintf("%s-message_size_over_time.csv", p.ID),
		"message_size_dist":      fmt.Sprintf("%s-message_size_distribution.csv", p.ID),
	}

	// File handler
	handles := make(map[string]*os.File, len(files))
	for key, fileName := range files {
		fullPath := filepath.Join(p.CsvDir, fileName)

		f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			// Close all already-opened files before panicking
			for _, h := range handles {
				h.Close()
			}
			panic(fmt.Errorf("failed to create plot file %s: %w", fullPath, err))
		}

		handles[key] = f
	}

	writers := make(map[string]*csv.Writer, len(handles))
	for key, file := range handles {
		writer := csv.NewWriter(file)
		writers[key] = writer
	}

	return writers
}

func ValidatePlotter(cfg *Config) error {
	switch strings.ToLower(cfg.Type) {
	case "csv":
		err := validateCsvPlotter(cfg)
		if err != nil {
			log.Printf("validating CsvPlotter failed with err: %v", err)
			return err
		}
		return nil
	default:
		return fmt.Errorf("%s is not supported", cfg.Type)
	}
}

func validateCsvPlotter(cfg *Config) error {
	return nil
}
