package ui

import (
	"testing"

	"github.com/becheran/roumon/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestFormatStatusSummary(t *testing.T) {
	routines := []model.Goroutine{
		{Status: "running"},
		{Status: "IO wait"},
		{Status: "running"},
		{Status: "chan receive"},
	}

	assert.Equal(t, "IO wait: 1\nchan receive: 1\nrunning: 2", formatStatusSummary(routines))
}

func TestFormatStatusSummaryEmpty(t *testing.T) {
	assert.Equal(t, "", formatStatusSummary(nil))
}
