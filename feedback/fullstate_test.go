package feedback

import (
	"testing"
)

func TestNewFullStateFeedback(t *testing.T) {
	tests := []struct {
		name string
		gain Values
	}{
		{
			name: "Single gain",
			gain: Values{1.0},
		},
		{
			name: "Multiple gains",
			gain: Values{1.5, 0.3, 2.1},
		},
		{
			name: "Zero gains",
			gain: Values{0.0, 0.0},
		},
		{
			name: "Negative gains",
			gain: Values{-1.0, -0.5},
		},
		{
			name: "Empty gains",
			gain: Values{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewFullStateFeedback(tt.gain)

			if controller == nil {
				t.Error("NewFullStateFeedback returned nil")
				return
			}

			if len(controller.gain) != len(tt.gain) {
				t.Errorf("Gain length = %d, expected %d", len(controller.gain), len(tt.gain))
			}

			for i, v := range tt.gain {
				if controller.gain[i] != v {
					t.Errorf("Gain[%d] = %f, expected %f", i, controller.gain[i], v)
				}
			}
		})
	}
}

func TestFullStateFeedbackCalculate(t *testing.T) {
	tests := []struct {
		name        string
		gain        Values
		setpoint    Values
		measurement Values
		expected    float64
		shouldError bool
	}{
		{
			name:        "Single dimension control",
			gain:        Values{2.0},
			setpoint:    Values{10.0},
			measurement: Values{8.0},
			expected:    4.0, // gain * (setpoint - measurement) = 2.0 * (10.0 - 8.0) = 4.0
			shouldError: false,
		},
		{
			name:        "Two dimension control",
			gain:        Values{1.5, 0.3},
			setpoint:    Values{10.0, 0.0},
			measurement: Values{8.5, 1.2},
			expected:    1.89, // 1.5*(10-8.5) + 0.3*(0-1.2) = 1.5*1.5 - 0.3*1.2 = 2.25 - 0.36 = 1.89
			shouldError: false,
		},
		{
			name:        "Three dimension control",
			gain:        Values{2.0, 0.5, 1.0},
			setpoint:    Values{5.0, 2.0, 1.0},
			measurement: Values{4.0, 3.0, 0.5},
			expected:    2.0, // 2.0*1.0 + 0.5*(-1.0) + 1.0*0.5 = 2.0 - 0.5 + 0.5 = 2.0
			shouldError: false,
		},
		{
			name:        "Zero error - no output",
			gain:        Values{1.0, 2.0},
			setpoint:    Values{5.0, 3.0},
			measurement: Values{5.0, 3.0},
			expected:    0.0,
			shouldError: false,
		},
		{
			name:        "Negative gains",
			gain:        Values{-1.0, -0.5},
			setpoint:    Values{2.0, 4.0},
			measurement: Values{1.0, 5.0},
			expected:    -0.5, // -1.0*1.0 + (-0.5)*(-1.0) = -1.0 + 0.5 = -0.5
			shouldError: false,
		},
		{
			name:        "Mismatched setpoint length",
			gain:        Values{1.0, 2.0},
			setpoint:    Values{5.0}, // Wrong length
			measurement: Values{3.0, 4.0},
			expected:    0.0,
			shouldError: true,
		},
		{
			name:        "Mismatched measurement length",
			gain:        Values{1.0, 2.0},
			setpoint:    Values{5.0, 6.0},
			measurement: Values{3.0}, // Wrong length
			expected:    0.0,
			shouldError: true,
		},
		{
			name:        "Empty vectors",
			gain:        Values{},
			setpoint:    Values{},
			measurement: Values{},
			expected:    0.0,
			shouldError: false,
		},
		{
			name:        "Gain length mismatch",
			gain:        Values{1.0}, // Wrong length
			setpoint:    Values{5.0, 6.0},
			measurement: Values{3.0, 4.0},
			expected:    0.0,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewFullStateFeedback(tt.gain)
			output, err := controller.Calculate(tt.setpoint, tt.measurement)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !almostEqual(output, tt.expected, 0.001) {
				t.Errorf("Calculate() = %f, expected %f", output, tt.expected)
			}
		})
	}
}

func TestFullStateFeedbackCalculateErrorConditions(t *testing.T) {
	controller := NewFullStateFeedback(Values{1.0, 0.5})

	// Test various error conditions
	testCases := []struct {
		name        string
		setpoint    Values
		measurement Values
		expectedErr string
	}{
		{
			name:        "Setpoint too short",
			setpoint:    Values{5.0},
			measurement: Values{3.0, 4.0},
			expectedErr: "vectors must be of same length",
		},
		{
			name:        "Measurement too short",
			setpoint:    Values{5.0, 6.0},
			measurement: Values{3.0},
			expectedErr: "vectors must be of same length",
		},
		{
			name:        "Both vectors wrong length",
			setpoint:    Values{5.0, 6.0, 7.0},
			measurement: Values{3.0, 4.0, 5.0},
			expectedErr: "vectors must be of same length",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := controller.Calculate(tc.setpoint, tc.measurement)
			if err == nil {
				t.Error("Expected error, got nil")
			} else if err.Error() != tc.expectedErr {
				t.Errorf("Error = %q, expected %q", err.Error(), tc.expectedErr)
			}
		})
	}
}

// Helper function for floating point comparison
func almostEqual(a, b, tolerance float64) bool {
	if a == b {
		return true
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < tolerance
}

func BenchmarkFullStateFeedbackCalculate(b *testing.B) {
	gain := Values{1.5, 0.3, 2.1, 0.8}
	controller := NewFullStateFeedback(gain)
	setpoint := Values{10.0, 5.0, 2.0, 8.0}
	measurement := Values{8.5, 6.0, 1.5, 7.2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		controller.Calculate(setpoint, measurement)
	}
}

func BenchmarkFullStateFeedbackCalculateLarge(b *testing.B) {
	// Test with larger vectors
	size := 100
	gain := make(Values, size)
	setpoint := make(Values, size)
	measurement := make(Values, size)

	for i := 0; i < size; i++ {
		gain[i] = float64(i%10) * 0.1
		setpoint[i] = float64(i)
		measurement[i] = float64(i) * 0.9
	}

	controller := NewFullStateFeedback(gain)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		controller.Calculate(setpoint, measurement)
	}
}
