package reader

import (
	"strings"
	"testing"
)

func TestReader_Next(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLines []string
		wantErr   bool
	}{
		{
			name:      "single line",
			input:     "line1",
			wantLines: []string{"line1"},
			wantErr:   false,
		},
		{
			name:      "multiple lines",
			input:     "line1\nline2\nline3",
			wantLines: []string{"line1", "line2", "line3"},
			wantErr:   false,
		},
		{
			name:      "empty input",
			input:     "",
			wantLines: []string{},
			wantErr:   false,
		},
		{
			name:      "line with tabs",
			input:     "field1\tfield2\tfield3",
			wantLines: []string{"field1\tfield2\tfield3"},
			wantErr:   false,
		},
		{
			name:      "line with custom delimiter",
			input:     "field1,field2,field3",
			wantLines: []string{"field1,field2,field3"},
			wantErr:   false,
		},
		{
			name:      "empty lines",
			input:     "\n\n\n",
			wantLines: []string{"", "", ""},
			wantErr:   false,
		},
		{
			name:      "trailing newline",
			input:     "line1\nline2\n",
			wantLines: []string{"line1", "line2"},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(strings.NewReader(tt.input))

			gotLines := []string{}
			for {
				line, ok, err := r.Next()
				if err != nil {
					if !tt.wantErr {
						t.Errorf("Next() error = %v, wantErr %v", err, tt.wantErr)
					}
					break
				}
				if !ok {
					break
				}
				gotLines = append(gotLines, line)
			}

			if len(gotLines) != len(tt.wantLines) {
				t.Errorf("Next() read %d lines, want %d", len(gotLines), len(tt.wantLines))
			}

			for i, got := range gotLines {
				if i < len(tt.wantLines) && got != tt.wantLines[i] {
					t.Errorf("Next() line %d = %q, want %q", i, got, tt.wantLines[i])
				}
			}
		})
	}
}

func TestReader_New(t *testing.T) {
	r := New(strings.NewReader("test"))
	if r == nil {
		t.Error("New() returned nil")
	}
	if r.scanner == nil {
		t.Error("New() scanner is nil")
	}
}

