package utils_test

import (
	"testing"

	"github.com/HarikrishnanBalagopal/asciidonut/utils"
)

func TestStep(t *testing.T) {
	t.Run("try multiple steps", func(t *testing.T) {
		time := 0.1
		for time < 10.0 {
			utils.Step(time)
			time += .1
		}
	})
}

func BenchmarkStep(t *testing.B) {
	t.Run("run multiple steps", func(t *testing.B) {
		time := 0.1
		for time < 10.0 {
			utils.Step(time)
			time += .1
		}
	})
}
