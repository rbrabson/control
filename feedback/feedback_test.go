package feedback

import (
	"testing"
)

// TestFeedbackInterface ensures that all feedback types implement the interface
func TestFeedbackInterface(t *testing.T) {
	tests := []struct {
		name           string
		controller     Feedback
		expectedOutput float64
	}{
		{
			name:           "NoFeedback implements Feedback",
			controller:     &NoFeedback{},
			expectedOutput: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the Calculate method can be called
			output := tt.controller.Calculate(10.0, 5.0)

			// Check expected output
			if output != tt.expectedOutput {
				t.Errorf("%s.Calculate() = %f, expected %f", tt.name, output, tt.expectedOutput)
			}
		})
	}
}

func TestValuesType(t *testing.T) {
	tests := []struct {
		name   string
		values Values
		length int
	}{
		{
			name:   "Empty Values",
			values: Values{},
			length: 0,
		},
		{
			name:   "Single value",
			values: Values{1.0},
			length: 1,
		},
		{
			name:   "Multiple values",
			values: Values{1.0, 2.0, 3.0},
			length: 3,
		},
		{
			name:   "Negative values",
			values: Values{-1.0, -2.5, 3.7},
			length: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.values) != tt.length {
				t.Errorf("Values length = %d, expected %d", len(tt.values), tt.length)
			}
		})
	}
}

func BenchmarkNoFeedbackCalculate(b *testing.B) {
	nf := &NoFeedback{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nf.Calculate(float64(i), float64(i*2))
	}
}
