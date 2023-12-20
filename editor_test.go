package editor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestBasicEditor_LaunchTempFile(t *testing.T) {
	type fields struct {
		Command  string
		LaunchFn func(command, file string) error
	}
	type args struct {
		prefix   string
		original io.Reader
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData []byte
		wantFile bool
		wantErr  bool
		wantDisk []byte
	}{
		{
			name: "successful launch",
			fields: fields{
				LaunchFn: func(command, file string) error {
					return nil
				},
			},
			args: args{
				prefix:   "prefix",
				original: bytes.NewBufferString("some random text"),
			},
			wantData: []byte("some random text"),
			wantFile: true,
			wantErr:  false,
			wantDisk: []byte("some random text"),
		},
		{
			name: "failed launch",
			fields: fields{
				LaunchFn: func(command, file string) error {
					return fmt.Errorf("failure to launch")
				},
			},
			args: args{
				prefix:   "prefix",
				original: bytes.NewBufferString("some random text"),
			},
			wantData: []byte{},
			wantFile: true,
			wantErr:  true,
			wantDisk: []byte("some random text"),
		},
		{
			name:   "execs command",
			fields: fields{Command: getCatCommand()},
			args: args{
				prefix:   "prefix",
				original: bytes.NewBufferString("some random text\n"),
			},
			wantData: []byte("some random text\n"),
			wantFile: true,
			wantErr:  false,
			wantDisk: []byte("some random text\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEditor()
			e.Command = tt.fields.Command
			if tt.fields.LaunchFn != nil {
				e.LaunchFn = tt.fields.LaunchFn
			}
			data, file, err := e.LaunchTempFile(tt.args.prefix, tt.args.original)
			defer os.Remove(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("BasicEditor.LaunchTempFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !bytes.Equal(data, tt.wantData) {
				t.Errorf("BasicEditor.LaunchTempFile() data = '%v', wantData '%v'", string(data), string(tt.wantData))
			}
			if (file != "") != tt.wantFile {
				t.Errorf("BasicEditor.LaunchTempFile() file = %v, wantFile %v", file, tt.wantFile)
				return
			}
			if file != "" {
				if _, err := os.Stat(file); os.IsNotExist(err) {
					t.Fatalf("BasicEditor.LaunchTempFile() temp file doesn't exist: %s", file)
				}
				if !strings.Contains(file, tt.args.prefix) {
					t.Errorf("BasicEditor.LaunchTempFile() file = %v, wantPrefix = %v", file, tt.args.prefix)
				}
				actual, err := os.ReadFile(file)
				if err != nil {
					t.Errorf("BasicEditor.LaunchTempFile() unable to read temp file: %s", file)
				}
				if !bytes.Equal(actual, tt.wantDisk) {
					t.Errorf("BasicEditor.LaunchTempFile() disk = '%v', wantData '%v'", string(actual), string(tt.wantDisk))
				}
			}
		})
	}
}

func getCatCommand() string {
	if runtime.GOOS == "windows" {
		return "Get-Content"
	} else {
		return "cat"
	}
}
