package flusher

import (
	"context"
	"fmt"
)

const batchSize = 10

type CountDataGetter interface {
	Get(n int) []string
}

type CountFlusher struct {
	dataGetter  CountDataGetter
	dataPrinter Printer
}

func NewCountFlusher(dataGetter CountDataGetter, dataPrinter Printer) *CountFlusher {
	return &CountFlusher{
		dataGetter:  dataGetter,
		dataPrinter: dataPrinter,
	}
}

func (f *CountFlusher) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return f.flush()
		default:
		}

		err := f.flush()
		if err != nil {
			return fmt.Errorf("can't flush: %w", err)
		}
	}
}

func (f *CountFlusher) flush() error {
	for _, data := range f.dataGetter.Get(batchSize) {
		err := f.dataPrinter.Print(data)
		if err != nil {
			return err
		}
	}

	return nil
}
