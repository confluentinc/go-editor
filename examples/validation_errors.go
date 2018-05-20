package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/codyaray/go-editor"
	"encoding/json"
)

type validationError struct {
	errors []error
}

func (e *validationError) Add(err error) error {
	e.errors = append(e.errors, err)
	return e
}

func (e *validationError) Error() string {
	out := &bytes.Buffer{}
	for _, err := range e.errors {
		out.WriteString(fmt.Sprintf("\n\t* %s", err.Error()))
	}
	return out.String()
}

func (e *validationError) IsEmpty() bool {
	return len(e.errors) == 0
}

type Format string
const FormatJSON = "json"

// Example schema that expects a known format and specific set of field keys
type keyedSchema struct {
	format Format
	fields []string
}

func (s *keyedSchema) ValidateBytes(data []byte) error {
	var obj map[string]interface{}
	ve := &validationError{}

	switch s.format {
	case FormatJSON:
		err := json.Unmarshal(data, &obj)
		if err != nil {
			return ve.Add(fmt.Errorf("data cannot be unmarshalled as JSON: %v", err))
		}
	default:
		return ve.Add(fmt.Errorf("unknown data format"))
	}

	for _, k := range s.fields {
		if _, ok := obj[k]; !ok {
			ve.Add(fmt.Errorf("missing field %s", k))
		}
	}

	if !ve.IsEmpty() {
		return ve
	}
	return nil
}

func main() {
	schema := &keyedSchema{format: "json", fields: []string{"key1", "key2", "key3"}}
	edit := editor.NewValidatingEditor(schema)

	obj := bytes.NewBufferString(`{"key1":1, "key1":2, "key1":3}`+"\n")

	contents, file, err := edit.LaunchTempFile("example", obj)
	defer os.Remove(file)
	if err != nil {
		fmt.Println("error: " + err.Error())
		os.Exit(1)
	}

	fmt.Println(string(contents))
}
