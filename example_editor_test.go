package editor_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codyaray/go-editor"
)

func Example_basic() {
	edit := editor.NewEditor()

	// Simulate user making changes
	edit.LaunchFn = func(command, file string) error {
		return ioutil.WriteFile(file, []byte("something else here"), 0777)
	}

	contents, file, err := edit.LaunchTempFile("example", bytes.NewBufferString("something to be edited\n"))
	defer os.Remove(file)
	if err != nil {
		fmt.Println("error: " + err.Error())
		os.Exit(1)
	}

	fmt.Println(string(contents))
	// Output:
	// something else here
}
