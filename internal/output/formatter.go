package output

import (
	"fmt"
	"os"
)

type Format string

const (
	FormatJSON  Format = "json"
	FormatTable Format = "table"
)

type Formatter interface {
	Format(data interface{}, columns []Column) error
}

type Column struct {
	Header string
	Field  string
	Width  int
}

func NewFormatter(format Format) Formatter {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}
	case FormatTable:
		return &TableFormatter{}
	default:
		return &TableFormatter{}
	}
}

func Print(format Format, data interface{}, columns []Column) {
	f := NewFormatter(format)
	if err := f.Format(data, columns); err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		os.Exit(1)
	}
}
