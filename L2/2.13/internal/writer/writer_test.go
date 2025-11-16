package writer

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriter_Write(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "single line",
			input:    "line1",
			expected: "line1\n",
			wantErr:  false,
		},
		{
			name:     "empty line",
			input:    "",
			expected: "\n",
			wantErr:  false,
		},
		{
			name:     "line with tabs",
			input:    "field1\tfield2",
			expected: "field1\tfield2\n",
			wantErr:  false,
		},
		{
			name:     "line with custom delimiter",
			input:    "field1,field2,field3",
			expected: "field1,field2,field3\n",
			wantErr:  false,
		},
		{
			name:     "long line",
			input:    strings.Repeat("a", 1000),
			expected: strings.Repeat("a", 1000) + "\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := New(buf)

			err := w.Write(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got := buf.String(); got != tt.expected {
				t.Errorf("Write() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestWriter_WriteMultiple(t *testing.T) {
	buf := &bytes.Buffer{}
	w := New(buf)

	lines := []string{"line1", "line2", "line3"}
	expected := strings.Join(lines, "\n") + "\n"

	for _, line := range lines {
		if err := w.Write(line); err != nil {
			t.Fatalf("Write() error = %v", err)
		}
	}

	if got := buf.String(); got != expected {
		t.Errorf("Write() multiple = %q, want %q", got, expected)
	}
}

func TestWriter_New(t *testing.T) {
	buf := &bytes.Buffer{}
	w := New(buf)
	if w == nil {
		t.Error("New() returned nil")
	}
	if w.w == nil {
		t.Error("New() writer is nil")
	}
}

