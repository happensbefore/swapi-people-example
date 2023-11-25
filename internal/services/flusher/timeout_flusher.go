package flusher

import (
	"context"
	"time"
)

const timeout = 250 * time.Millisecond

type AllDataGetter interface {
	GetAll() []string
}

type TimeoutFlusher struct {
	dataGetter  AllDataGetter
	dataPrinter Printer
}

func NewTimeoutFlusher(dataGetter AllDataGetter, dataPrinter Printer) *TimeoutFlusher {
	return &TimeoutFlusher{
		dataGetter:  dataGetter,
		dataPrinter: dataPrinter,
	}
}

func (f *TimeoutFlusher) Run(ctx context.Context) error {
	for {
		time.Sleep(timeout)

		select {
		case <-ctx.Done():
			f.flush()
			return nil
		default:
		}

		f.flush()
	}
}

func (f *TimeoutFlusher) flush() {
	for _, data := range f.dataGetter.GetAll() {
		f.dataPrinter.Print(data)
	}
}
