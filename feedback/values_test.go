package feedback

import (
	"math"
	"testing"
)

func TestValues(t *testing.T) {
	tests := []struct {
		name   string
		values Values
		length int
	}{
		{
			name:   "Empty values",
			values: Values{},
			length: 0,
		},
		{
			name:   "Single value",
			values: Values{1.5},
			length: 1,
		},
		{
			name:   "Multiple values",
			values: Values{1.0, 2.5, -3.7, 0.0},
			length: 4,
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

func TestValuesOperations(t *testing.T) {
	// Test various operations on Values
	v := Values{1.0, 2.0, 3.0}

	// Test indexing
	if v[0] != 1.0 || v[1] != 2.0 || v[2] != 3.0 {
		t.Error("Values indexing failed")
	}

	// Test modification
	v[1] = 5.0
	if v[1] != 5.0 {
		t.Error("Values modification failed")
	}

	// Test append
	v = append(v, 4.0)
	if len(v) != 4 || v[3] != 4.0 {
		t.Error("Values append failed")
	}
}

func TestNoFeedback(t *testing.T) {
	nf := &NoFeedback{}

	// Test that NoFeedback implements Feedback interface
	var _ Feedback = nf

	tests := []struct {
		name        string
		setpoint    float64
		measurement float64
		expected    float64
	}{
		{
			name:        "Positive values",
			setpoint:    5.0,
			measurement: 3.0,
			expected:    0.0,
		},
		{
			name:        "Negative values",
			setpoint:    -2.0,
			measurement: -5.0,
			expected:    0.0,
		},
		{
			name:        "Zero values",
			setpoint:    0.0,
			measurement: 0.0,
			expected:    0.0,
		},
		{
			name:        "Large values",
			setpoint:    1000.0,
			measurement: 999.9,
			expected:    0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := nf.Calculate(tt.setpoint, tt.measurement)

			if output != tt.expected {
				t.Errorf("NoFeedback.Calculate() = %f, expected %f", output, tt.expected)
			}
		})
	}
}

func TestNoFeedbackThroughInterface(t *testing.T) {
	// Test that NoFeedback implements the Feedback interface
	nf := &NoFeedback{}
	var _ Feedback = nf

	// Test basic functionality through interface
	var feedback Feedback = nf
	output := feedback.Calculate(5.0, 3.0)

	if output != 0.0 {
		t.Errorf("NoFeedback through interface = %f, expected 0.0", output)
	}
}

func BenchmarkNoFeedback(b *testing.B) {
	nf := &NoFeedback{}
	setpoint := 10.0
	measurement := 8.5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nf.Calculate(setpoint, measurement)
	}
}

// Test for NaN values
func TestNoFeedbackWithSpecialValues(t *testing.T) {
	nf := &NoFeedback{}

	tests := []struct {
		name        string
		setpoint    float64
		measurement float64
		expected    float64
	}{
		{
			name:        "Infinity values",
			setpoint:    math.Inf(1),
			measurement: math.Inf(-1),
			expected:    0.0,
		},
		{
			name:        "NaN values",
			setpoint:    math.NaN(),
			measurement: 5.0,
			expected:    0.0,
		},
		{
			name:        "Very large values",
			setpoint:    1e100,
			measurement: -1e100,
			expected:    0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := nf.Calculate(tt.setpoint, tt.measurement)

			if output != tt.expected {
				t.Errorf("NoFeedback.Calculate() = %f, expected %f", output, tt.expected)
			}
		})
	}
}
