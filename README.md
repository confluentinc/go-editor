# go-editor

Allow your CLI users to edit arbitrary data in their preferred editor.

Just like editing messages in `git commit` or resources with `kubectl edit`.

## Install

    go get github.com/codyaray/go-editor

## Usage

Provide any `io.Reader` with the initial contents:

	editor := editor.NewEditor()
	original := bytes.NewBufferString("something to be edited\n")
	edited, path, err := editor.LaunchTempFile("example", original)

The library leaves it up to you to cleanup the temp file. For example, this
allows your CLI to validate the edited data and prompt the user to continue
editing where they left off, rather than starting their changes over.

When you're done, be sure to clean up after yourself:

	defer os.Remove(path)

Happy editing!

## Acknowledgements

Thanks to these other projects and groups for pointing the way.

* [kubernetes/kubernetes](https://github.com/kubernetes/kubernetes)
* [AlecAivazis/survey](https://github.com/AlecAivazis/survey)
* [golang/nuts](https://groups.google.com/forum/#!topic/golang-nuts/cuAEvgqqYFU)
