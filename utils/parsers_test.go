package utils_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/zsamora/utils"
)

func TestFormatDate(t *testing.T) {
	input := time.Date(2024, time.March, 5, 0, 0, 0, 0, time.UTC)
	want := "05/03/2024"
	if got := utils.FormatDate(input); got != want {
		t.Errorf("FormatDate(%v) = %q, want %q", input, got, want)
	}
}

func TestFormatDatetime(t *testing.T) {
	input := time.Date(2024, time.March, 5, 14, 30, 45, 0, time.UTC)
	want := "05/03/2024 14:30:45"
	if got := utils.FormatDatetime(input); got != want {
		t.Errorf("FormatDatetime(%v) = %q, want %q", input, got, want)
	}
}

func TestFormatYear(t *testing.T) {
	input := time.Date(2024, time.March, 5, 0, 0, 0, 0, time.UTC)
	want := "2024"
	if got := utils.FormatYear(input); got != want {
		t.Errorf("FormatYear(%v) = %q, want %q", input, got, want)
	}
}

func TestFormatUint(t *testing.T) {
	tests := []struct {
		input uint64
		want  string
	}{
		{0, "0"},
		{255, "255"},
		{1234567890, "1234567890"},
	}
	for _, tt := range tests {
		if got := utils.FormatUint(tt.input); got != tt.want {
			t.Errorf("FormatUint(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"valid positive", "42", 42, false},
		{"valid negative", "-5", -5, false},
		{"empty string returns an error", "", 0, true},
		{"non-numeric returns an error", "abc", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ParseInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseInt(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseUUID(t *testing.T) {
	valid := uuid.New()
	tests := []struct {
		name  string
		input string
		want  uuid.UUID
	}{
		{"valid UUID round-trips", valid.String(), valid},
		{"invalid string silently becomes uuid.Nil", "not-a-uuid", uuid.Nil},
		{"empty string silently becomes uuid.Nil", "", uuid.Nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ParseUUID(tt.input); got != tt.want {
				t.Errorf("ParseUUID(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseDateForm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Time
	}{
		{"valid date", "2024-03-05", time.Date(2024, time.March, 5, 0, 0, 0, 0, time.UTC)},
		{"invalid date silently becomes zero time", "not-a-date", time.Time{}},
		{"empty string silently becomes zero time", "", time.Time{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ParseDateForm(tt.input); !got.Equal(tt.want) {
				t.Errorf("ParseDateForm(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseYearForm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Time
	}{
		{"valid year", "2024", time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"invalid year: fallback to \"0\" also fails silently", "abc", time.Time{}},
		{"empty string: fallback to \"0\" also fails silently", "", time.Time{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ParseYearForm(tt.input); !got.Equal(tt.want) {
				t.Errorf("ParseYearForm(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
