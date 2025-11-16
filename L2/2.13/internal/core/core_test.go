package core

import (
	"bytes"
	"cut/internal/parser"
	"cut/internal/reader"
	"cut/internal/writer"
	"strings"
	"testing"
)

func TestCore_Process(t *testing.T) {
	tests := []struct {
		name     string
		opts     parser.Options
		input    string
		expected string
		wantOut  bool
	}{
		{
			name: "single field with tab delimiter",
			opts: parser.Options{
				Delimiter: '\t',
				Mask:      map[int]bool{0: true},
			},
			input:    "field1\tfield2\tfield3",
			expected: "field1",
			wantOut:  true,
		},
		{
			name: "multiple fields with tab delimiter",
			opts: parser.Options{
				Delimiter: '\t',
				Mask:      map[int]bool{0: true, 2: true},
			},
			input:    "field1\tfield2\tfield3",
			expected: "field1\tfield3",
			wantOut:  true,
		},
		{
			name: "field range",
			opts: parser.Options{
				Delimiter: '\t',
				Mask:      map[int]bool{1: true, 2: true, 3: true},
			},
			input:    "field1\tfield2\tfield3\tfield4\tfield5",
			expected: "field2\tfield3\tfield4",
			wantOut:  true,
		},
		{
			name: "custom delimiter (comma)",
			opts: parser.Options{
				Delimiter: ',',
				Mask:      map[int]bool{0: true, 2: true},
			},
			input:    "field1,field2,field3",
			expected: "field1,field3",
			wantOut:  true,
		},
		{
			name: "separated flag - line with delimiter",
			opts: parser.Options{
				Delimiter: '\t',
				Separated: true,
				Mask:      map[int]bool{0: true},
			},
			input:    "field1\tfield2",
			expected: "field1",
			wantOut:  true,
		},
		{
			name: "separated flag - line without delimiter",
			opts: parser.Options{
				Delimiter: '\t',
				Separated: true,
				Mask:      map[int]bool{0: true},
			},
			input:    "field1 field2",
			expected: "",
			wantOut:  false,
		},
		{
			name: "out of bounds field - ignored",
			opts: parser.Options{
				Delimiter: '\t',
				Mask:      map[int]bool{0: true, 10: true},
			},
			input:    "field1\tfield2\tfield3",
			expected: "field1",
			wantOut:  true,
		},
		{
			name: "no matching fields",
			opts: parser.Options{
				Delimiter: '\t',
				Mask:      map[int]bool{10: true},
			},
			input:    "field1\tfield2\tfield3",
			expected: "",
			wantOut:  false,
		},
		{
			name: "all fields",
			opts: parser.Options{
				Delimiter: '\t',
				Mask:      map[int]bool{0: true, 1: true, 2: true},
			},
			input:    "field1\tfield2\tfield3",
			expected: "field1\tfield2\tfield3",
			wantOut:  true,
		},
		{
			name: "empty line with delimiter",
			opts: parser.Options{
				Delimiter: '\t',
				Mask:      map[int]bool{0: true},
			},
			input:    "\t",
			expected: "",
			wantOut:  false, // empty field produces no output
		},
		{
			name: "custom delimiter colon",
			opts: parser.Options{
				Delimiter: ':',
				Mask:      map[int]bool{1: true},
			},
			input:    "field1:field2:field3",
			expected: "field2",
			wantOut:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &Core{
				opts: tt.opts,
			}

			result, gotOut := core.Process(tt.input)
			if gotOut != tt.wantOut {
				t.Errorf("Process() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
			if gotOut && result != tt.expected {
				t.Errorf("Process() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCore_Run(t *testing.T) {
	tests := []struct {
		name     string
		opts     parser.Options
		input    string
		expected string
	}{
		{
			name: "single line",
			opts: parser.Options{
				Delimiter: '\t',
				Mask:      map[int]bool{0: true},
			},
			input:    "field1\tfield2\tfield3\n",
			expected: "field1\n",
		},
		{
			name: "multiple lines",
			opts: parser.Options{
				Delimiter: '\t',
				Mask:      map[int]bool{1: true},
			},
			input:    "field1\tfield2\tfield3\nfield4\tfield5\tfield6\n",
			expected: "field2\nfield5\n",
		},
		{
			name: "with separated flag",
			opts: parser.Options{
				Delimiter: '\t',
				Separated: true,
				Mask:      map[int]bool{0: true},
			},
			input:    "field1\tfield2\nfield3 field4\nfield5\tfield6\n",
			expected: "field1\nfield5\n",
		},
		{
			name: "custom delimiter",
			opts: parser.Options{
				Delimiter: ',',
				Mask:      map[int]bool{0: true, 2: true},
			},
			input:    "field1,field2,field3\nfield4,field5,field6\n",
			expected: "field1,field3\nfield4,field6\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			output := &bytes.Buffer{}

			core := &Core{
				r:    reader.New(input),
				w:    writer.New(output),
				opts: tt.opts,
			}

			err := core.Run()
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			if output.String() != tt.expected {
				t.Errorf("Run() output = %q, want %q", output.String(), tt.expected)
			}
		})
	}
}

func TestCore_New(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid args",
			args:    []string{"-f", "1"},
			wantErr: false,
		},
		{
			name:    "invalid args",
			args:    []string{"-f"},
			wantErr: true,
		},
		{
			name:    "empty args",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := New(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && core == nil {
				t.Error("New() returned nil core when no error expected")
			}
		})
	}
}

