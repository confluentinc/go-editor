package editor

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

var (
	editor = "vim"
)

func init() {
	if runtime.GOOS == "windows" {
		editor = "notepad"
	}
	if e := os.Getenv("VISUAL"); e != "" {
		editor = e
	} else if e := os.Getenv("EDITOR"); e != "" {
		editor = e
	}
}

type Editor struct {
	command string
}

// NewEditor creates an Editor with the users preferred text editor. The editor
// to use is determined by reading the $VISUAL and $EDITOR environment variables.
// If neither of these are present, vim or notepad (on Windows) is used.
func NewEditor() *Editor {
	return &Editor{
		command: editor,
	}
}

// Launch opens the given file path in the external editor or returns an error.
func (e *Editor) Launch(path string) error {
	cmd := exec.Command(e.command, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// LaunchTempFile launches the users preferred editor on a temporary file
// initialized with contents from the provided stream and named with the given
// prefix. Returns the modified data and the path to the temporary file so the
// caller can clean it up, or an error.
func (e *Editor) LaunchTempFile(prefix string, r io.Reader) ([]byte, string, error) {
	f, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	// seed the editor with the initial temp file contents
	if _, err := io.Copy(f, r); err != nil {
		os.Remove(f.Name())
		return nil, "", err
	}

	// close the fd to prevent the editor being unable to save file
	if err := f.Close(); err != nil {
		return nil, "", err
	}

	// launch the external editor on the temp file
	if err := e.Launch(f.Name()); err != nil {
		return nil, f.Name(), err
	}

	bytes, err := ioutil.ReadFile(f.Name())
	return bytes, f.Name(), err
}
