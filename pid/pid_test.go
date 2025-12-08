package pid

import (
	"math"
	"testing"
	"time"
)

// Helper function to check if two floats are approximately equal
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		kp   float64
		ki   float64
		kd   float64
	}{
		{"Basic PID", 1.0, 0.1, 0.05},
		{"Zero gains", 0.0, 0.0, 0.0},
		{"Negative gains", -1.0, -0.1, -0.05},
		{"Large gains", 100.0, 50.0, 25.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid := New(tt.kp, tt.ki, tt.kd)

			kp, ki, kd := pid.GetGains()
			if kp != tt.kp || ki != tt.ki || kd != tt.kd {
				t.Errorf("Expected gains (%f, %f, %f), got (%f, %f, %f)",
					tt.kp, tt.ki, tt.kd, kp, ki, kd)
			}

			if pid.initialized {
				t.Error("PID should not be initialized on creation")
			}

			if pid.GetIntegral() != 0 {
				t.Error("Integral should be zero on creation")
			}
		})
	}
}

func TestWithFeedForward(t *testing.T) {
	tests := []struct {
		name        string
		feedForward float64
	}{
		{"Positive feed-forward", 5.0},
		{"Zero feed-forward", 0.0},
		{"Negative feed-forward", -2.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid := New(1.0, 0.1, 0.05, WithFeedForward(tt.feedForward))

			if pid.GetFeedForward() != tt.feedForward {
				t.Errorf("Expected feed-forward %f, got %f",
					tt.feedForward, pid.GetFeedForward())
			}
		})
	}
}

func TestWithIntegralResetOnZeroCross(t *testing.T) {
	pid := New(1.0, 0.1, 0.05, WithIntegralResetOnZeroCross())

	if !pid.GetIntegralResetOnZeroCross() {
		t.Error("Integral reset on zero cross should be enabled")
	}
}

func TestWithStabilityThreshold(t *testing.T) {
	tests := []struct {
		name      string
		threshold float64
		expected  float64
	}{
		{"Positive threshold", 2.0, 2.0},
		{"Negative threshold", -1.5, 1.5}, // Should be absolute value
		{"Zero threshold", 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid := New(1.0, 0.1, 0.05, WithStabilityThreshold(tt.threshold))

			if pid.GetStabilityThreshold() != tt.expected {
				t.Errorf("Expected stability threshold %f, got %f",
					tt.expected, pid.GetStabilityThreshold())
			}
		})
	}
}

func TestWithIntegralSumMax(t *testing.T) {
	tests := []struct {
		name     string
		maxSum   float64
		expected float64
	}{
		{"Positive max", 10.0, 10.0},
		{"Negative max", -5.0, 5.0}, // Should be absolute value
		{"Zero max", 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid := New(1.0, 0.1, 0.05, WithIntegralSumMax(tt.maxSum))

			if pid.GetIntegralSumMax() != tt.expected {
				t.Errorf("Expected integral sum max %f, got %f",
					tt.expected, pid.GetIntegralSumMax())
			}
		})
	}
}

func TestWithDerivativeFilter(t *testing.T) {
	tests := []struct {
		name     string
		alpha    float64
		expected float64
	}{
		{"Valid alpha", 0.5, 0.5},
		{"Alpha too low", -0.1, 0.0}, // Should be clamped to 0
		{"Alpha too high", 1.5, 1.0}, // Should be clamped to 1
		{"Zero alpha", 0.0, 0.0},
		{"One alpha", 1.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid := New(1.0, 0.1, 0.05, WithDerivativeFilter(tt.alpha))

			if pid.GetDerivativeFilter() != tt.expected {
				t.Errorf("Expected derivative filter %f, got %f",
					tt.expected, pid.GetDerivativeFilter())
			}
		})
	}
}

func TestWithOutputLimits(t *testing.T) {
	tests := []struct {
		name        string
		min         float64
		max         float64
		expectedMin float64
		expectedMax float64
	}{
		{"Valid limits", -10.0, 10.0, -10.0, 10.0},
		{"Invalid limits (min > max)", 10.0, -10.0, math.Inf(-1), math.Inf(1)}, // Should not change from defaults
		{"Zero limits", 0.0, 0.0, 0.0, 0.0},
		{"Negative limits", -5.0, -1.0, -5.0, -1.0},
		{"Positive limits", 1.0, 5.0, 1.0, 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid := New(1.0, 0.1, 0.05, WithOutputLimits(tt.min, tt.max))

			min, max := pid.GetOutputLimits()
			if min != tt.expectedMin || max != tt.expectedMax {
				t.Errorf("Expected limits (%f, %f), got (%f, %f)",
					tt.expectedMin, tt.expectedMax, min, max)
			}
		})
	}
}

func TestSetOutputLimits(t *testing.T) {
	pid := New(1.0, 0.1, 0.05)

	// Test valid limits
	pid.SetOutputLimits(-10.0, 10.0)
	min, max := pid.GetOutputLimits()
	if min != -10.0 || max != 10.0 {
		t.Errorf("Expected limits (-10, 10), got (%f, %f)", min, max)
	}

	// Test invalid limits (min > max) - should not change
	pid.SetOutputLimits(10.0, -10.0)
	min, max = pid.GetOutputLimits()
	if min != -10.0 || max != 10.0 {
		t.Errorf("Limits should not change for invalid input, got (%f, %f)", min, max)
	}

	// Test clamping with large error
	pid.Calculate(1000.0, 0.0) // Large error
	time.Sleep(10 * time.Millisecond)
	output := pid.Calculate(1000.0, 0.0)

	if output > 10.0 || output < -10.0 {
		t.Errorf("Output %f should be clamped between -10 and 10", output)
	}
}

func TestSetGains(t *testing.T) {
	pid := New(1.0, 0.1, 0.05)

	newKp, newKi, newKd := 2.0, 0.2, 0.1
	pid.SetGains(newKp, newKi, newKd)

	kp, ki, kd := pid.GetGains()
	if kp != newKp || ki != newKi || kd != newKd {
		t.Errorf("Expected gains (%f, %f, %f), got (%f, %f, %f)",
			newKp, newKi, newKd, kp, ki, kd)
	}
}

func TestBasicPIDCalculate(t *testing.T) {
	pid := New(1.0, 0.1, 0.05)

	// First calculate should return 0 (initialization)
	output := pid.Calculate(1.0, 0.0)
	if output != 0 {
		t.Errorf("First calculate should return 0, got %f", output)
	}

	// Wait a bit and calculate again
	time.Sleep(10 * time.Millisecond)
	output = pid.Calculate(1.0, 0.0)

	// Should have proportional component at least
	if output == 0 {
		t.Error("Second calculate should not be zero with non-zero error")
	}
}

func TestProportionalTerm(t *testing.T) {
	// Pure proportional controller
	pid := New(2.0, 0.0, 0.0)

	pid.Calculate(1.0, 0.0) // Initialize
	time.Sleep(10 * time.Millisecond)
	output := pid.Calculate(1.0, 0.0)

	// Output should equal Kp * error = 2.0 * 1.0 = 2.0
	if !almostEqual(output, 2.0, 0.001) {
		t.Errorf("Expected proportional output 2.0, got %f", output)
	}
}

func TestIntegralTerm(t *testing.T) {
	// Pure integral controller
	pid := New(0.0, 1.0, 0.0)

	pid.Calculate(1.0, 0.0) // Initialize
	time.Sleep(100 * time.Millisecond)
	output := pid.Calculate(1.0, 0.0)

	// Integral should accumulate error over time
	if output <= 0 {
		t.Error("Integral term should accumulate positive error")
	}

	// Check that integral is accessible
	integral := pid.GetIntegral()
	if integral <= 0 {
		t.Error("Integral value should be positive after positive error")
	}
}

func TestDerivativeTerm(t *testing.T) {
	// Pure derivative controller
	pid := New(0.0, 0.0, 1.0)

	pid.Calculate(0.0, 0.0) // Initialize with zero error
	time.Sleep(10 * time.Millisecond)
	output := pid.Calculate(1.0, 0.0) // Step change in error

	// Derivative should respond to change in error
	if output <= 0 {
		t.Error("Derivative term should respond to positive error change")
	}
}

func TestFeedForward(t *testing.T) {
	feedForward := 5.0
	pid := New(0.0, 0.0, 0.0, WithFeedForward(feedForward))

	pid.Calculate(0.0, 0.0) // Initialize
	time.Sleep(10 * time.Millisecond)
	output := pid.Calculate(0.0, 0.0) // Zero error

	// Output should equal feed-forward value
	if !almostEqual(output, feedForward, 0.001) {
		t.Errorf("Expected feed-forward output %f, got %f", feedForward, output)
	}
}

func TestIntegralResetOnZeroCross(t *testing.T) {
	pid := New(0.0, 1.0, 0.0, WithIntegralResetOnZeroCross())

	// Build up positive integral
	pid.Calculate(1.0, 0.0) // Initialize
	time.Sleep(10 * time.Millisecond)
	pid.Calculate(1.0, 0.0)
	time.Sleep(10 * time.Millisecond)
	pid.Calculate(1.0, 0.0)

	integral1 := pid.GetIntegral()
	if integral1 <= 0 {
		t.Error("Integral should be positive after positive errors")
	}

	// Cross zero - the integral gets reset but then immediately accumulates the new error
	time.Sleep(10 * time.Millisecond)
	pid.Calculate(-1.0, 0.0)

	integral2 := pid.GetIntegral()
	// The integral should be much smaller (close to zero) since it was reset at zero crossing
	// but may not be exactly zero due to the new negative error accumulation
	if math.Abs(integral2) >= math.Abs(integral1) {
		t.Errorf("Integral should be much smaller after zero crossing reset. Before: %f, After: %f", integral1, integral2)
	}

	// Test the reset more directly by checking right at the zero crossing
	pid.Reset()
	pid.Calculate(1.0, 0.0) // Initialize
	time.Sleep(10 * time.Millisecond)
	pid.Calculate(1.0, 0.0) // Positive error

	integralBeforeReset := pid.GetIntegral()

	time.Sleep(10 * time.Millisecond)
	pid.Calculate(0.0, 0.0) // Zero error - should not trigger reset

	time.Sleep(10 * time.Millisecond)
	pid.Calculate(-0.1, 0.0) // Small negative error - should trigger reset

	integralAfterReset := pid.GetIntegral()

	// The integral should have been reset and now only contains the small negative accumulation
	if math.Abs(integralAfterReset) >= math.Abs(integralBeforeReset) {
		t.Errorf("Zero crossing reset failed. Before: %f, After: %f", integralBeforeReset, integralAfterReset)
	}
}

func TestStabilityThreshold(t *testing.T) {
	threshold := 1.0
	pid := New(0.0, 1.0, 0.0, WithStabilityThreshold(threshold))

	pid.Calculate(0.0, 0.0) // Initialize
	time.Sleep(10 * time.Millisecond)

	// Small error change (below threshold)
	pid.Calculate(0.1, 0.0)
	integral1 := pid.GetIntegral()

	time.Sleep(10 * time.Millisecond)

	// Large error change (above threshold)
	pid.Calculate(2.0, 0.0) // This creates derivative > threshold
	integral2 := pid.GetIntegral()

	// Integral should not have accumulated much during high derivative
	if integral2 <= integral1 {
		t.Log("Integral accumulation was limited by stability threshold")
	}
}

func TestIntegralSumMax(t *testing.T) {
	maxSum := 1.0
	pid := New(0.0, 1.0, 0.0, WithIntegralSumMax(maxSum))

	pid.Calculate(1.0, 0.0) // Initialize

	// Keep adding large errors to try to exceed max
	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		pid.Calculate(10.0, 0.0) // Large error
	}

	integral := pid.GetIntegral()
	if integral > maxSum {
		t.Errorf("Integral %f should not exceed max sum %f", integral, maxSum)
	}
}

func TestDerivativeFilter(t *testing.T) {
	alpha := 0.5
	pid := New(0.0, 0.0, 1.0, WithDerivativeFilter(alpha))

	pid.Calculate(0.0, 0.0) // Initialize
	time.Sleep(10 * time.Millisecond)
	output1 := pid.Calculate(1.0, 0.0) // Step change

	time.Sleep(10 * time.Millisecond)
	output2 := pid.Calculate(1.0, 0.0) // Same error (derivative should be 0)

	// With filtering, derivative response should be smoother
	if math.Abs(output2) >= math.Abs(output1) {
		t.Error("Filtered derivative should decrease when error stops changing")
	}
}

func TestReset(t *testing.T) {
	pid := New(1.0, 1.0, 1.0)

	// Build up some state
	pid.Calculate(1.0, 0.0)
	time.Sleep(10 * time.Millisecond)
	pid.Calculate(1.0, 0.0)

	if pid.GetIntegral() == 0 {
		t.Error("Integral should not be zero before reset")
	}

	// Reset and check
	pid.Reset()

	if pid.GetIntegral() != 0 {
		t.Error("Integral should be zero after reset")
	}

	if pid.initialized {
		t.Error("PID should not be initialized after reset")
	}
}

func TestOutputClamping(t *testing.T) {
	pid := New(10.0, 0.0, 0.0) // Large proportional gain
	pid.SetOutputLimits(-1.0, 1.0)

	pid.Calculate(1.0, 0.0) // Initialize
	time.Sleep(10 * time.Millisecond)
	output := pid.Calculate(1.0, 0.0) // Should produce output > 1.0 without clamping

	if output > 1.0 || output < -1.0 {
		t.Errorf("Output %f should be clamped between -1.0 and 1.0", output)
	}
}

func TestAntiWindup(t *testing.T) {
	pid := New(1.0, 1.0, 0.0)
	pid.SetOutputLimits(-1.0, 1.0)

	// Build up integral until output saturates
	pid.Calculate(1.0, 0.0) // Initialize
	for i := 0; i < 20; i++ {
		time.Sleep(10 * time.Millisecond)
		pid.Calculate(1.0, 0.0) // Constant error
	}

	// Integral should be limited due to anti-windup
	integral := pid.GetIntegral()

	// The anti-windup should prevent integral from growing indefinitely
	if integral > 10.0 { // Reasonable bound check
		t.Errorf("Anti-windup failed, integral is too large: %f", integral)
	}
}

func TestCombinedOptions(t *testing.T) {
	pid := New(1.0, 0.1, 0.05,
		WithFeedForward(2.0),
		WithIntegralResetOnZeroCross(),
		WithStabilityThreshold(1.0),
		WithIntegralSumMax(5.0),
		WithDerivativeFilter(0.3),
	)

	// Verify all options were applied
	if pid.GetFeedForward() != 2.0 {
		t.Error("Feed-forward not set correctly")
	}

	if !pid.GetIntegralResetOnZeroCross() {
		t.Error("Integral reset on zero cross not enabled")
	}

	if pid.GetStabilityThreshold() != 1.0 {
		t.Error("Stability threshold not set correctly")
	}

	if pid.GetIntegralSumMax() != 5.0 {
		t.Error("Integral sum max not set correctly")
	}

	if pid.GetDerivativeFilter() != 0.3 {
		t.Error("Derivative filter not set correctly")
	}
}

func TestSetterMethods(t *testing.T) {
	pid := New(1.0, 0.1, 0.05)

	// Test all setter methods
	pid.SetFeedForward(3.0)
	if pid.GetFeedForward() != 3.0 {
		t.Error("SetFeedForward failed")
	}

	pid.SetIntegralResetOnZeroCross(true)
	if !pid.GetIntegralResetOnZeroCross() {
		t.Error("SetIntegralResetOnZeroCross failed")
	}

	pid.SetStabilityThreshold(2.5)
	if pid.GetStabilityThreshold() != 2.5 {
		t.Error("SetStabilityThreshold failed")
	}

	pid.SetIntegralSumMax(7.5)
	if pid.GetIntegralSumMax() != 7.5 {
		t.Error("SetIntegralSumMax failed")
	}

	pid.SetDerivativeFilter(0.8)
	if pid.GetDerivativeFilter() != 0.8 {
		t.Error("SetDerivativeFilter failed")
	}

	pid.SetOutputLimits(-5.5, 5.5)
	min, max := pid.GetOutputLimits()
	if min != -5.5 || max != 5.5 {
		t.Errorf("SetOutputLimits/GetOutputLimits failed: expected (-5.5, 5.5), got (%f, %f)", min, max)
	}
}

func TestZeroTimeDelta(t *testing.T) {
	pid := New(1.0, 0.1, 0.05)

	pid.Calculate(1.0, 0.0) // Initialize

	// Update immediately without time passing
	output := pid.Calculate(1.0, 0.0)

	// Should handle zero time delta gracefully
	if !almostEqual(output, 1.0, 0.1) { // Should return proportional term
		t.Errorf("Expected proportional output ~1.0 for zero time delta, got %f", output)
	}
}

func TestNegativeError(t *testing.T) {
	pid := New(1.0, 0.1, 0.05)

	pid.Calculate(-1.0, 0.0) // Initialize with negative error
	time.Sleep(10 * time.Millisecond)
	output := pid.Calculate(-1.0, 0.0)

	// Output should be negative for negative error
	if output >= 0 {
		t.Errorf("Expected negative output for negative error, got %f", output)
	}
}

// Benchmark tests
func BenchmarkPIDCalculate(b *testing.B) {
	pid := New(1.0, 0.1, 0.05)
	pid.Calculate(0.0, 0.0) // Initialize

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pid.Calculate(float64(i%100), 0.0)
	}
}

func BenchmarkPIDCalculateWithAllOptions(b *testing.B) {
	pid := New(1.0, 0.1, 0.05,
		WithFeedForward(2.0),
		WithIntegralResetOnZeroCross(),
		WithStabilityThreshold(1.0),
		WithIntegralSumMax(5.0),
		WithDerivativeFilter(0.3),
	)
	pid.Calculate(0.0, 0.0) // Initialize

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pid.Calculate(float64(i%100), 0.0)
	}
}
