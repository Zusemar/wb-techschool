package core

import (
	"bytes"
	"cut/internal/reader"
	"cut/internal/writer"
	"strings"
	"testing"
)

// TestIntegration_FlagCombinations tests various flag combinations as required by the spec
func TestIntegration_FlagCombinations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "basic -f with single field",
			args:     []string{"-f", "1"},
			input:    "field1\tfield2\tfield3\n",
			expected: "field1\n",
			wantErr:  false,
		},
		{
			name:     "-f with range",
			args:     []string{"-f", "2-4"},
			input:    "field1\tfield2\tfield3\tfield4\tfield5\n",
			expected: "field2\tfield3\tfield4\n",
			wantErr:  false,
		},
		{
			name:     "-f with multiple fields and ranges",
			args:     []string{"-f", "1,3-5,7"},
			input:    "field1\tfield2\tfield3\tfield4\tfield5\tfield6\tfield7\n",
			expected: "field1\tfield3\tfield4\tfield5\tfield7\n",
			wantErr:  false,
		},
		{
			name:     "-f and -d (comma delimiter)",
			args:     []string{"-f", "1,3", "-d", ","},
			input:    "field1,field2,field3\n",
			expected: "field1,field3\n",
			wantErr:  false,
		},
		{
			name:     "-f, -d, and -s (skip lines without delimiter)",
			args:     []string{"-f", "1", "-d", ",", "-s"},
			input:    "field1,field2\nfield3 field4\nfield5,field6\n",
			expected: "field1\nfield5\n",
			wantErr:  false,
		},
		{
			name:     "-f and -s (tab delimiter, skip lines without delimiter)",
			args:     []string{"-f", "2", "-s"},
			input:    "field1\tfield2\tfield3\nfield4 field5\nfield6\tfield7\n",
			expected: "field2\nfield7\n",
			wantErr:  false,
		},
		{
			name:     "all flags together with colon delimiter",
			args:     []string{"-f", "1,3-4", "-d", ":", "-s"},
			input:    "field1:field2:field3:field4\nfield5 field6\nfield7:field8:field9\n",
			expected: "field1:field3:field4\nfield7:field9\n", // field 4 doesn't exist in second line, so only fields 1 and 3
			wantErr:  false,
		},
		{
			name:     "out of bounds field (should be ignored)",
			args:     []string{"-f", "1,10"},
			input:    "field1\tfield2\tfield3\n",
			expected: "field1\n",
			wantErr:  false,
		},
		{
			name:     "multiple lines with range",
			args:     []string{"-f", "2-3"},
			input:    "a\tb\tc\td\nx\ty\tz\tw\n",
			expected: "b\tc\ny\tz\n",
			wantErr:  false,
		},
		{
			name:     "custom delimiter (semicolon)",
			args:     []string{"-f", "1,3", "-d", ";"},
			input:    "field1;field2;field3\n",
			expected: "field1;field3\n",
			wantErr:  false,
		},
		{
			name:     "fields with spaces in list",
			args:     []string{"-f", "1, 3 , 5"},
			input:    "field1\tfield2\tfield3\tfield4\tfield5\n",
			expected: "field1\tfield3\tfield5\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := New(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			// Replace stdin/stdout with test buffers
			input := strings.NewReader(tt.input)
			output := &bytes.Buffer{}

			core.r = reader.New(input)
			core.w = writer.New(output)

			err = core.Run()
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			if output.String() != tt.expected {
				t.Errorf("Run() output = %q, want %q", output.String(), tt.expected)
			}
		})
	}
}

