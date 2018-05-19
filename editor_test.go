package editor

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLaunchTempFile(t *testing.T) {
	editor := Editor{command: "cat"}
	expected := "something to be edited\n"

	contents, path, err := editor.LaunchTempFile("prefix", bytes.NewBufferString(expected))
	if err != nil {
		t.Fatalf("error launching temp file: %v", err)
	}

	// check if temp file still exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("temp file doesn't exist: %s", path)
	}
	defer os.Remove(path)

	// check if filename is as expected
	if !strings.Contains(path, "prefix") {
		t.Errorf("filename doesn't contain prefix: %s", path)
	}

	// check if returned contents are as expected
	if string(contents) != expected {
		t.Errorf("returned contents don't match: %s", string(contents))
	}

	// check if temp file contents are as expected
	actual, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("unable to read temp file: %s", path)
	}
	if string(actual) != expected {
		t.Errorf("temp file contents don't match: %s", string(actual))
	}
}
