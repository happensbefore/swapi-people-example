package services

import (
	"context"
	"fmt"
)

type Printer interface {
	Print(data []byte) error
}

type Serializer interface {
	Serialize(data string) []byte
}

const printBufCapacity = 16

type PrintService struct {
	receiver   chan string
	printer    Printer
	serializer Serializer
}

func NewPrintService(serializer Serializer, printer Printer) *PrintService {
	return &PrintService{
		receiver:   make(chan string, printBufCapacity),
		printer:    printer,
		serializer: serializer,
	}
}

func (s *PrintService) Print(data string) {
	s.receiver <- data
}

func (s *PrintService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-s.receiver:
			err := s.printer.Print(s.serializer.Serialize(data))
			if err != nil {
				return fmt.Errorf("can't print: %w", err)
			}
		}
	}
}
