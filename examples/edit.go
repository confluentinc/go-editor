package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/codyaray/go-editor"
)

func main() {
	editor := editor.NewEditor()
	contents, path, err := editor.LaunchTempFile("example", bytes.NewBufferString("something to be edited\n"))
	if err != nil {
		panic(err)
	}
	defer os.Remove(path)

	fmt.Println(string(contents))
}
