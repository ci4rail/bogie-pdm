package signalprocessing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRMS(t *testing.T) {
	assert := assert.New(t)

	data := []float64{1, 2, 3, 4, 5}
	assert.InDelta(3.31662, RMS(data), 0.001)

	data = []float64{10, 4, 6, 8}
	assert.InDelta(7.34847, RMS(data), 0.001)
}
