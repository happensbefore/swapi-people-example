package services

import (
	"fmt"
)

type Printer interface {
	Print(data []byte) error
}

type Serializer interface {
	Serialize(data string) []byte
}

type PrintService struct {
	printer    Printer
	serializer Serializer
}

func NewPrintService(serializer Serializer, printer Printer) *PrintService {
	return &PrintService{
		printer:    printer,
		serializer: serializer,
	}
}

func (s *PrintService) Print(data string) error {
	err := s.printer.Print(s.serializer.Serialize(data))
	if err != nil {
		return fmt.Errorf("can't print: %w", err)
	}

	return nil
}
