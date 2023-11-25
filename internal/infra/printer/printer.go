package printer

import (
	"fmt"
	"io"
)

type Printer struct {
	writer io.Writer
}

func New(writer io.Writer) *Printer {
	return &Printer{
		writer: writer,
	}
}

func (p *Printer) Print(data []byte) error {
	_, err := p.writer.Write(data)
	if err != nil {
		return fmt.Errorf("can't write: %w", err)
	}

	return nil
}
