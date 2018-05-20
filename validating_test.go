package editor

import (
	"bytes"
	//"io"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type alwaysValidSchema struct{}

func (s *alwaysValidSchema) ValidateBytes([]byte) error {
	return nil
}

type alwaysInvalidSchema struct{}

func (s *alwaysInvalidSchema) ValidateBytes([]byte) error {
	return fmt.Errorf("invalid")
}

type compoundSchema struct {
	schemas []Schema
	count   int
}

func (s *compoundSchema) ValidateBytes(data []byte) error {
	ret := s.schemas[s.count]
	s.count++
	return ret.ValidateBytes(data)
}

func TestValidatingEditor_LaunchTempFile(t *testing.T) {
	type fields struct {
		Schema Schema
	}
	type args struct {
		prefix   string
		original string
		edited   []string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData string
		wantFile bool
		wantErr  string
	}{
		{
			name: "cancel on original unchanged",
			fields: fields{
				Schema: &alwaysInvalidSchema{},
			},
			args: args{
				original: "original data",
				edited:   []string{"original data"},
			},
			wantData: "",
			wantFile: false,
			wantErr:  msgCancelledNoOrigChanges,
		},
		{
			name: "cancel on invalid edit and then unchanged edit",
			fields: fields{
				Schema: &alwaysInvalidSchema{},
			},
			args: args{
				original: "original data",
				edited:   []string{"invalid data", "invalid data"},
			},
			wantData: "invalid data",
			wantFile: true,
			wantErr:  "invalid " + msgCancelledNoValidChanges,
		},
		{
			name: "cancel on invalid edit, different invalid edit, and then unchanged edit",
			fields: fields{
				Schema: &alwaysInvalidSchema{},
			},
			args: args{
				original: "original data",
				edited:   []string{"invalid data", "more invalid data", "more invalid data"},
			},
			wantData: "more invalid data",
			wantFile: true,
			wantErr:  "invalid " + msgCancelledNoValidChanges,
		},
		{
			name: "cancel on empty file",
			fields: fields{
				Schema: &alwaysInvalidSchema{},
			},
			args: args{
				original: "original data",
				edited:   []string{"invalid data", ""},
			},
			wantData: "",
			wantFile: false,
			wantErr:  msgCancelledEmptyFile,
		},
		{
			name: "cancel on invalid edit, and then empty file",
			fields: fields{
				Schema: &alwaysInvalidSchema{},
			},
			args: args{
				original: "original data",
				edited:   []string{"invalid data", "more invalid data", ""},
			},
			wantData: "",
			wantFile: false,
			wantErr:  msgCancelledEmptyFile,
		},
		{
			name: "cancel on comment-only file",
			fields: fields{
				Schema: &alwaysInvalidSchema{},
			},
			args: args{
				original: "original data",
				edited:   []string{"invalid data", "# foo\n  # world\n \t\r\n"},
			},
			wantData: "",
			wantFile: false,
			wantErr:  msgCancelledEmptyFile,
		},
		{
			name: "successful edit on first try",
			fields: fields{
				Schema: &alwaysValidSchema{},
			},
			args: args{
				original: "original data",
				edited:   []string{"new data"},
			},
			wantData: "new data",
			wantFile: true,
			wantErr:  "",
		},
		{
			name: "successful edit on third try",
			fields: fields{
				Schema: &compoundSchema{
					schemas: []Schema{
						&alwaysInvalidSchema{},
						&alwaysInvalidSchema{},
						&alwaysValidSchema{},
					},
				},
			},
			args: args{
				original: "original data",
				edited:   []string{"invalid data", "more invalid data", "new data"},
			},
			wantData: "new data",
			wantFile: true,
			wantErr:  msgCancelledEmptyFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewValidatingEditor(tt.fields.Schema)
			e.InvalidFn = func(e error) error { return fmt.Errorf("%s %s", e.Error(), msgCancelledNoValidChanges) }
			e.PreserveFileFn = func(data []byte, file string, err error) ([]byte, string, error) { return data, file, err }
			editCount := 0
			e.LaunchFn = func(command, file string) error {
				if editCount >= len(tt.args.edited) {
					return fmt.Errorf("EDITOR_NEVER_EXITED")
				}
				err := ioutil.WriteFile(file, []byte(tt.args.edited[editCount]), 0777)
				editCount++
				return err
			}
			data, file, err := e.LaunchTempFile(tt.args.prefix, bytes.NewBufferString(tt.args.original))
			defer os.Remove(file)
			if (err != nil) != (tt.wantErr != "") && (err != nil && (err.Error() != tt.wantErr)) {
				t.Errorf("ValidatingEditor.LaunchTempFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(string(data), tt.wantData) {
				t.Errorf("ValidatingEditor.LaunchTempFile() data = '%v', wantData '%v'", string(data), tt.wantData)
			}
			if (file != "") != tt.wantFile {
				t.Errorf("ValidatingEditor.LaunchTempFile() file = %v, wantFile %v", file, tt.wantFile)
			}
			if editCount != len(tt.args.edited) {
				t.Errorf("ValidatingEditor.LaunchTempFile() editCount = %v, wantEditCount %v", editCount, len(tt.args.edited))
			}
		})
	}
}

func TestValidatingEditor_isEmpty(t *testing.T) {
	type args struct {
		comments []string
		data     []byte
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "empty file",
			args: args{
				data: []byte(""),
			},
			want: true,
		},
		{
			name: "just whitespace",
			args: args{
				data: []byte("  \t  \r\n"),
			},
			want: true,
		},
		{
			name: "one line, commented with #",
			args: args{
				data: []byte("# hello"),
			},
			want: true,
		},
		{
			name: "one line, commented with //",
			args: args{
				data: []byte("// hello"),
			},
			want: true,
		},
		{
			name: "two lines, one commented",
			args: args{
				data: []byte("# hello\nworld"),
			},
			want: false,
		},
		{
			name: "two lines, both commented",
			args: args{
				data: []byte("# hello\n//world"),
			},
			want: true,
		},
		{
			name: "multiple lines, whitespace and comments",
			args: args{
				data: []byte("# hello\n//world\n\n  \t  \r\n"),
			},
			want: true,
		},
		{
			name: "multiple lines, whitespace and comments and text",
			args: args{
				data: []byte("# hello\n//world\n\n  \t  \r\n   not empty  "),
			},
			want: false,
		},
		{
			name: "not really commented",
			args: args{
				data: []byte("hello # world"),
			},
			want: false,
		},
		{
			name: "custom comment char",
			args: args{
				comments: []string{";"},
				data:     []byte("; hello"),
			},
			want: true,
		},
		{
			name: "custom comment char without comments",
			args: args{
				comments: []string{";"},
				data:     []byte("hello"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewValidatingEditor(&alwaysValidSchema{})
			if tt.args.comments != nil {
				e.CommentChars = tt.args.comments
			}
			e.LaunchFn = func(command, file string) error {
				return ioutil.WriteFile(file, tt.args.data, 0777)
			}
			got, err := e.isEmpty(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatingEditor.isEmpty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidatingEditor.isEmpty() = %v, wantData %v", got, tt.want)
			}
		})
	}
}
