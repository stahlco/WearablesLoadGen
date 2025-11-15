package handler

import (
	"WearablesLoadGen/pkg/plotter"
	"sync"
)

type Base struct {
	Plotter   plotter.Plotter
	CallCount int
	mu        sync.Mutex
}

func (b *Base) logToCSV(requests int) error {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Close() error {
	//TODO implement me
	panic("implement me")
}
