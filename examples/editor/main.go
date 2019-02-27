package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/confluentinc/go-editor"
)

func main() {
	edit := editor.NewEditor()
	contents, file, err := edit.LaunchTempFile("example", bytes.NewBufferString("something to be edited\n"))
	defer os.Remove(file)
	if err != nil {
		fmt.Println("error: " + err.Error())
		os.Exit(1)
	}

	fmt.Println(string(contents))
}
