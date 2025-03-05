package InMemory

import (
	"strconv"
	"testing"
)

func TestMultipleTimes(t *testing.T) {
	const numberOfRuns = 100
	for i := 0; i < numberOfRuns; i++ {
		t.Run("Iteration:"+strconv.Itoa(i), TestWithSimulation)
	}
}
