package flusher

import (
	"context"
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
			f.flush()
			return nil
		default:
		}

		f.flush()
	}
}

func (f *CountFlusher) flush() {
	for _, data := range f.dataGetter.Get(batchSize) {
		f.dataPrinter.Print(data)
	}
}
