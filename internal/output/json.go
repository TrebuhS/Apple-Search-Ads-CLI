package output

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(data interface{}, columns []Column) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	return nil
}
