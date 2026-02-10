package output

import (
	"fmt"
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

type TableFormatter struct{}

func (f *TableFormatter) Format(data interface{}, columns []Column) error {
	val := reflect.ValueOf(data)

	// Handle pointer
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// If it's not a slice, wrap it
	if val.Kind() != reflect.Slice {
		slice := reflect.MakeSlice(reflect.SliceOf(val.Type()), 1, 1)
		slice.Index(0).Set(val)
		val = slice
	}

	if val.Len() == 0 {
		fmt.Println("No results found.")
		return nil
	}

	table := tablewriter.NewTable(os.Stdout)

	// Set headers
	headers := make([]string, len(columns))
	for i, col := range columns {
		headers[i] = col.Header
	}
	table.Header(headers)

	// Fill rows
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		row := make([]string, len(columns))
		for j, col := range columns {
			row[j] = getFieldValue(item, col.Field)
		}
		table.Append(row)
	}

	table.Render()
	return nil
}

func getFieldValue(v reflect.Value, field string) string {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", v.Interface())
	}

	f := v.FieldByName(field)
	if !f.IsValid() {
		return ""
	}

	// Handle pointer fields
	if f.Kind() == reflect.Ptr {
		if f.IsNil() {
			return ""
		}
		f = f.Elem()
	}

	// Handle slice fields (e.g. RoleNames, CountriesOrRegions)
	if f.Kind() == reflect.Slice {
		var parts []string
		for i := 0; i < f.Len(); i++ {
			parts = append(parts, fmt.Sprintf("%v", f.Index(i).Interface()))
		}
		return fmt.Sprintf("%v", parts)
	}

	// Handle Money type
	if f.Kind() == reflect.Struct {
		if amount := f.FieldByName("Amount"); amount.IsValid() {
			currency := f.FieldByName("Currency")
			if currency.IsValid() {
				return fmt.Sprintf("%s %s", amount.Interface(), currency.Interface())
			}
		}
		return fmt.Sprintf("%v", f.Interface())
	}

	return fmt.Sprintf("%v", f.Interface())
}
