package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/codyaray/go-editor"
)

// Example schema that expects an exact string match
type exactSchema struct {
	m string
}

func (s *exactSchema) ValidateBytes(data []byte) error {
	if string(data) != s.m {
		return fmt.Errorf("obj doesn't match")
	}
	return nil
}

func main() {
	schema := &exactSchema{m: "something else \n"}
	edit := editor.NewValidatingEditor(schema)

	obj := bytes.NewBufferString("something to be edited\n")

	contents, file, err := edit.LaunchTempFile("example", obj)
	defer os.Remove(file)
	if err != nil {
		fmt.Println("error: " + err.Error())
		os.Exit(1)
	}

	fmt.Println(string(contents))
}
