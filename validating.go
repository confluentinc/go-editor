package editor

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var (
	ErrValidationFailed        = fmt.Errorf("The edited file failed validation")
	ErrCancelledNoValidChanges = fmt.Errorf("Edit cancelled, no valid changes were saved.")
	ErrCancelledNoOrigChanges  = fmt.Errorf("Edit cancelled, no changes made.")
	ErrCancelledEmptyFile      = fmt.Errorf("Edit cancelled, saved file was empty.")

	defaultInvalidFn = func(e error) error {
		fmt.Printf("%s: %v\n", ErrValidationFailed, e)
		return fmt.Errorf("%s\n", ErrCancelledNoValidChanges)
	}
	defaultNoChangesFn    = func() (bool, error) { return true, ErrCancelledNoOrigChanges }
	defaultEmptyFileFn    = func() (bool, error) { return true, ErrCancelledEmptyFile }
	defaultPreserveFileFn = func(data []byte, path string, err error) ([]byte, string, error) {
		fmt.Printf("A copy of your changes has been stored to %s\n", path)
		return data, path, err
	}
	defaultCommentChars = []string{"#", "//"}
)

// ValidatingEditor is an Editor which validates edited data and retries invalid edits until successful or canceled.
type ValidatingEditor struct {
	BasicEditor

	// Schema is used to validate the edited data
	Schema Schema

	// InvalidFn is called when a Schema fails to validate data
	InvalidFn ValidationFailedFn
	// OriginalNoChangesFn is called when no changes were made from the original
	OriginalNoChangesFn CancelEditingFn
	// EmptyFileFn is called when the edited data is (effectively) empty; this means the file doesn't have any uncommented lines (ignoring whitespace)
	EmptyFileFn CancelEditingFn
	// PreserveFileFn is called when a non-recoverable error has occurred to give you the chance to inform the user that their edits have been preserved and where to find them
	PreserveFileFn PreserveFileFn

	// CommentChars is a list of comment string prefixes and defaults to "#" and "//"
	CommentChars []string
}

// NewValidatingEditor creates an ValidatingEditor with the users preferred text editor.
// The editor to use is determined by reading the $VISUAL and $EDITOR environment variables.
// If neither of these are present, vim or notepad (on Windows) is used.
func NewValidatingEditor(schema Schema) *ValidatingEditor {
	return &ValidatingEditor{
		BasicEditor:         BasicEditor{Command: editor},
		Schema:              schema,
		InvalidFn:           defaultInvalidFn,
		OriginalNoChangesFn: defaultNoChangesFn,
		EmptyFileFn:         defaultEmptyFileFn,
		PreserveFileFn:      defaultPreserveFileFn,
		CommentChars:        defaultCommentChars,
	}
}

// LaunchTempFile launches the users preferred editor on a temporary file
// initialized with contents from the provided stream and named with the given
// prefix. Returns the modified data and the path to the temporary file so the
// caller can clean it up, or an error.
//
// The last byte of `obj` must be a newline to cancel editing if no changes are made.
// (This is because many editors like vim automatically add a newline when saving.)
func (e *ValidatingEditor) LaunchTempFile(prefix string, obj io.Reader) ([]byte, string, error) {
	editor := NewEditor()

	var (
		prevErr  error
		original []byte
		edited   []byte
		file     string
		err      error
	)

	originalObj, err := ioutil.ReadAll(obj)
	if err != nil {
		return nil, "", err
	}

	// loop until we succeed or cancel editing
	for {
		// Create the file to edit
		buf := &bytes.Buffer{}
		if prevErr == nil {
			buf.Write(originalObj)
			original = buf.Bytes()
		} else {
			// Preserve the edited file
			buf.Write(edited)
		}

		// Launch the editor
		editedDiff := edited
		edited, file, err = editor.LaunchTempFile(prefix, buf)
		if err != nil {
			return e.PreserveFileFn(edited, file, err)
		}

		// If we're retrying the loop because of an error, and no change was made in the file, short-circuit
		if prevErr != nil && bytes.Equal(editedDiff, edited) {
			return e.PreserveFileFn(edited, file, e.InvalidFn(prevErr))
		}

		// Compare contents for changes
		if bytes.Equal(original, edited) {
			cancel, err := e.OriginalNoChangesFn()
			if cancel {
				os.Remove(file)
				return nil, "", err
			}
		}

		// Check for an (effectively) empty file
		empty, err := e.isEmpty(edited)
		if err != nil {
			return e.PreserveFileFn(edited, file, err)
		}
		if empty {
			cancel, err := e.EmptyFileFn()
			if cancel {
				os.Remove(file)
				return nil, "", err
			}
		}

		// Apply validation
		err = e.Schema.ValidateBytes(edited)
		if err != nil {
			prevErr = err
			os.Remove(file)
			continue
		}

		return edited, file, nil
	}
}

// isEmpty returns true if the file doesn't have any uncommented lines (ignoring whitespace)
func (e *ValidatingEditor) isEmpty(data []byte) (bool, error) {
	empty := true
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			commented := false
			for _, c := range e.CommentChars {
				if strings.HasPrefix(line, c) {
					commented = true
				}
			}
			if !commented {
				empty = false
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return empty, nil
}
