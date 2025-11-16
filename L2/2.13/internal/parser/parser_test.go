package parser

import (
	"testing"
)

func TestParse_BasicFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(*testing.T, Options)
	}{
		{
			name:    "empty args",
			args:    []string{},
			wantErr: true,
		},
		{
			name: "single field",
			args: []string{"-f", "1"},
			check: func(t *testing.T, opts Options) {
				if opts.Delimiter != '\t' {
					t.Errorf("expected delimiter '\t', got %c", opts.Delimiter)
				}
				if !opts.Mask[0] {
					t.Error("field 1 should be in mask")
				}
				if len(opts.Mask) != 1 {
					t.Errorf("expected 1 field in mask, got %d", len(opts.Mask))
				}
				if opts.Separated {
					t.Error("separated should be false")
				}
			},
		},
		{
			name: "multiple fields",
			args: []string{"-f", "1,3,5"},
			check: func(t *testing.T, opts Options) {
				if !opts.Mask[0] || !opts.Mask[2] || !opts.Mask[4] {
					t.Error("fields 1, 3, 5 should be in mask")
				}
				if opts.Mask[1] || opts.Mask[3] {
					t.Error("fields 2, 4 should not be in mask")
				}
				if len(opts.Mask) != 3 {
					t.Errorf("expected 3 fields in mask, got %d", len(opts.Mask))
				}
			},
		},
		{
			name: "field range",
			args: []string{"-f", "3-5"},
			check: func(t *testing.T, opts Options) {
				if !opts.Mask[2] || !opts.Mask[3] || !opts.Mask[4] {
					t.Error("fields 3-5 should be in mask")
				}
				if opts.Mask[1] || opts.Mask[5] {
					t.Error("fields 2, 6 should not be in mask")
				}
				if len(opts.Mask) != 3 {
					t.Errorf("expected 3 fields in mask, got %d", len(opts.Mask))
				}
			},
		},
		{
			name: "mixed fields and ranges",
			args: []string{"-f", "1,3-5,7"},
			check: func(t *testing.T, opts Options) {
				if !opts.Mask[0] || !opts.Mask[2] || !opts.Mask[3] || !opts.Mask[4] || !opts.Mask[6] {
					t.Error("fields 1, 3-5, 7 should be in mask")
				}
				if len(opts.Mask) != 5 {
					t.Errorf("expected 5 fields in mask, got %d", len(opts.Mask))
				}
			},
		},
		{
			name: "custom delimiter",
			args: []string{"-f", "1", "-d", ","},
			check: func(t *testing.T, opts Options) {
				if opts.Delimiter != ',' {
					t.Errorf("expected delimiter ',', got %c", opts.Delimiter)
				}
			},
		},
		{
			name: "separated flag",
			args: []string{"-f", "1", "-s"},
			check: func(t *testing.T, opts Options) {
				if !opts.Separated {
					t.Error("separated should be true")
				}
			},
		},
		{
			name: "all flags together",
			args: []string{"-f", "1,3-5", "-d", ":", "-s"},
			check: func(t *testing.T, opts Options) {
				if opts.Delimiter != ':' {
					t.Errorf("expected delimiter ':', got %c", opts.Delimiter)
				}
				if !opts.Separated {
					t.Error("separated should be true")
				}
				if !opts.Mask[0] || !opts.Mask[2] || !opts.Mask[3] || !opts.Mask[4] {
					t.Error("fields 1, 3-5 should be in mask")
				}
			},
		},
		{
			name:    "missing -f value",
			args:    []string{"-f"},
			wantErr: true,
		},
		{
			name:    "invalid field number (negative)",
			args:    []string{"-f", "-1"},
			wantErr: true,
		},
		{
			name:    "invalid field number (zero)",
			args:    []string{"-f", "0"},
			wantErr: true,
		},
		{
			name:    "invalid range (from > to)",
			args:    []string{"-f", "5-3"},
			wantErr: true,
		},
		{
			name:    "invalid range format (empty from)",
			args:    []string{"-f", "-5"},
			wantErr: true,
		},
		{
			name:    "invalid range format (empty to)",
			args:    []string{"-f", "5-"},
			wantErr: true,
		},
		{
			name:    "invalid delimiter (empty)",
			args:    []string{"-f", "1", "-d", ""},
			wantErr: true,
		},
		{
			name:    "invalid delimiter (multiple chars)",
			args:    []string{"-f", "1", "-d", "ab"},
			wantErr: true,
		},
		{
			name:    "unknown flag",
			args:    []string{"-f", "1", "-x"},
			wantErr: true,
		},
		{
			name:    "missing -f flag",
			args:    []string{"-d", ","},
			wantErr: true,
		},
		{
			name: "delimiter without value",
			args: []string{"-f", "1", "-d"},
			check: func(t *testing.T, opts Options) {
				// Should default to tab
				if opts.Delimiter != '\t' {
					t.Errorf("expected delimiter '\t', got %c", opts.Delimiter)
				}
			},
		},
		{
			name: "fields with spaces",
			args: []string{"-f", "1, 3 , 5"},
			check: func(t *testing.T, opts Options) {
				if !opts.Mask[0] || !opts.Mask[2] || !opts.Mask[4] {
					t.Error("fields 1, 3, 5 should be in mask")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, opts)
			}
		})
	}
}

func TestParse_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(*testing.T, Options)
	}{
		{
			name: "single character range",
			args: []string{"-f", "1-1"},
			check: func(t *testing.T, opts Options) {
				if !opts.Mask[0] {
					t.Error("field 1 should be in mask")
				}
				if len(opts.Mask) != 1 {
					t.Errorf("expected 1 field in mask, got %d", len(opts.Mask))
				}
			},
		},
		{
			name: "large range",
			args: []string{"-f", "1-10"},
			check: func(t *testing.T, opts Options) {
				if len(opts.Mask) != 10 {
					t.Errorf("expected 10 fields in mask, got %d", len(opts.Mask))
				}
				for i := 0; i < 10; i++ {
					if !opts.Mask[i] {
						t.Errorf("field %d should be in mask", i+1)
					}
				}
			},
		},
		{
			name: "out of bounds field (ignored)",
			args: []string{"-f", "1,100"},
			check: func(t *testing.T, opts Options) {
				if !opts.Mask[0] {
					t.Error("field 1 should be in mask")
				}
				if !opts.Mask[99] {
					t.Error("field 100 should be in mask")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, opts)
			}
		})
	}
}

