package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/codyaray/go-editor"
)

func main() {
	edit := editor.NewEditor()
	contents, path, err := edit.LaunchTempFile("example", bytes.NewBufferString("something to be edited\n"))
	defer os.Remove(path)
	if err != nil {
		fmt.Println("error: " + err.Error())
		os.Exit(1)
	}

	fmt.Println(string(contents))
}
