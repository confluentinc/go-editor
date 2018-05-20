package editor

import "io"

// Editor launches an instance of the users preferred editor.
// The editor to use is determined by reading the $VISUAL or $EDITOR environment
// variables. If neither of these are present, vim or notepad (on Windows) is used.
type Editor interface {
	Launch(path string) error
	LaunchTempFile(prefix string, r io.Reader) ([]byte, string, error)
}

// Schema is an interface for validating data
type Schema interface {
	ValidateBytes(data []byte) error
}

// ValidationFailedFn is a function with which you can check, modify, or handle the validation error
type ValidationFailedFn func(error) error

// CancelEditingFn is a function with which you cancel editing and provide a suitable error message
type CancelEditingFn func() (bool, error)

// PreserveFileFn is a function with which you can inspect the preserved file, edited data, and resulting error
type PreserveFileFn func(data []byte, path string, err error) ([]byte, string, error)
