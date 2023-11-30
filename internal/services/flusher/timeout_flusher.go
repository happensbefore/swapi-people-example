package flusher

import (
	"context"
	"fmt"
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
			return f.flush()
		default:
		}

		err := f.flush()
		if err != nil {
			return fmt.Errorf("can't flush: %w", err)
		}
	}
}

func (f *TimeoutFlusher) flush() error {
	for _, data := range f.dataGetter.GetAll() {
		err := f.dataPrinter.Print(data)
		if err != nil {
			return err
		}
	}

	return nil
}
