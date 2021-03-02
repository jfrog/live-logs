package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMillisToDuration(t *testing.T) {
	tests := []struct {
		name           string
		timeInMillis   int64
		wantedDuration time.Duration
	}{
		{
			name:           "100 millis",
			timeInMillis:   int64(100),
			wantedDuration: 100 * time.Millisecond,
		},
		{
			name:           "0 millis",
			timeInMillis:   int64(0),
			wantedDuration: 0 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDuration := MillisToDuration(tt.timeInMillis)
			assert.Equal(t, tt.wantedDuration, gotDuration)
		})
	}
}

func TestSliceToCsv(t *testing.T) {
	tests := []struct {
		name      string
		values    []string
		wantedCsv string
	}{
		{
			name:      "single value",
			values:    []string{"a"},
			wantedCsv: "a",
		},
		{
			name:      "multiple value",
			values:    []string{"a", "b", "c"},
			wantedCsv: "a,b,c",
		},
		{
			name:      "no value",
			values:    []string{},
			wantedCsv: "",
		},
		{
			name:      "nil value",
			values:    nil,
			wantedCsv: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCsv := SliceToCsv(tt.values)
			assert.Equal(t, tt.wantedCsv, gotCsv)
		})
	}
}

func TestInSlice(t *testing.T) {
	tests := []struct {
		name        string
		slice       []string
		val         string
		wantedFound bool
	}{
		{
			name:        "single value",
			slice:       []string{"a"},
			val:         "a",
			wantedFound: true,
		},
		{
			name:        "multiple values",
			slice:       []string{"a", "b"},
			val:         "a",
			wantedFound: true,
		},
		{
			name:        "multiple values not found",
			slice:       []string{"a", "b"},
			val:         "c",
			wantedFound: false,
		},
		{
			name:        "nil slice",
			slice:       nil,
			val:         "c",
			wantedFound: false,
		},
		{
			name:        "empty slice",
			slice:       []string{},
			val:         "c",
			wantedFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := InSlice(tt.slice, tt.val)
			assert.Equal(t, tt.wantedFound, found)
		})
	}
}
