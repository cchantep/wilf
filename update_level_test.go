package main

import (
	"testing"
)

func TestParseUpdateLevel(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    UpdateLevel
		wantErr bool
	}{
		{
			name:    "valid patch",
			input:   "patch",
			want:    Patch,
			wantErr: false,
		},
		{
			name:    "valid minor",
			input:   "minor",
			want:    Minor,
			wantErr: false,
		},
		{
			name:    "valid major",
			input:   "major",
			want:    Major,
			wantErr: false,
		},
		{
			name:    "invalid update level",
			input:   "invalid",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUpdateLevel(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUpdateLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ParseUpdateLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
