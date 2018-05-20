package editor

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLaunchTempFile(t *testing.T) {
	editor := NewEditor()

	expected := "something to be edited\n"

	// Simulate user making changes
	editor.LaunchFn = func(command, file string) error {
		return ioutil.WriteFile(file, []byte(expected), 0777)
	}

	contents, file, err := editor.LaunchTempFile("prefix", bytes.NewBufferString(expected))
	if err != nil {
		t.Fatalf("error launching temp file: %v", err)
	}

	// check if temp file still exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Fatalf("temp file doesn't exist: %s", file)
	}
	defer os.Remove(file)

	// check if filename is as expected
	if !strings.Contains(file, "prefix") {
		t.Errorf("filename doesn't contain prefix: %s", file)
	}

	// check if returned contents are as expected
	if string(contents) != expected {
		t.Errorf("returned contents don't match: %s", string(contents))
	}

	// check if temp file contents are as expected
	actual, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("unable to read temp file: %s", file)
	}
	if string(actual) != expected {
		t.Errorf("temp file contents don't match: %s", string(actual))
	}
}
