package editor_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/confluentinc/go-editor"
)

// Example schema that expects a string prefix match
type prefixSchema struct {
	prefix string
}

func (s *prefixSchema) ValidateBytes(data []byte) error {
	if !strings.HasPrefix(string(data), s.prefix) {
		return fmt.Errorf("data missing prefix")
	}
	return nil
}

func Example_validating() {
	schema := &prefixSchema{prefix: "something"}
	edit := editor.NewValidatingEditor(schema)

	// Simulate user making changes
	edit.LaunchFn = func(command, file string) error {
		return os.WriteFile(file, []byte("something else here"), 0777)
	}

	obj := bytes.NewBufferString("something else")

	contents, file, err := edit.LaunchTempFile("example", obj)
	defer os.Remove(file)
	if err != nil {
		fmt.Println("error: " + err.Error())
		os.Exit(1)
	}

	fmt.Println(string(contents))
	// Output:
	// something else here
}
