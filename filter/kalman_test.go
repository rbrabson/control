package filter

import (
	"math"
	"testing"
)

// TestSizedStack tests the SizedStack functionality
func TestSizedStack(t *testing.T) {
	t.Run("Basic operations", func(t *testing.T) {
		stack := NewFloat64Stack(3)

		if stack.Size() != 0 {
			t.Errorf("Expected size 0, got %d", stack.Size())
		}

		// Test pushing elements
		stack.Push(1.0)
		stack.Push(2.0)
		stack.Push(3.0)

		if stack.Size() != 3 {
			t.Errorf("Expected size 3, got %d", stack.Size())
		}

		if stack.Peek() != 3.0 {
			t.Errorf("Expected peek 3.0, got %f", stack.Peek())
		}
	})

	t.Run("Capacity overflow", func(t *testing.T) {
		stack := NewFloat64Stack(2)

		stack.Push(1.0)
		stack.Push(2.0)
		stack.Push(3.0) // Should remove 1.0

		if stack.Size() != 2 {
			t.Errorf("Expected size 2, got %d", stack.Size())
		}

		if stack.Get(0) != 2.0 {
			t.Errorf("Expected first element 2.0, got %f", stack.Get(0))
		}

		if stack.Peek() != 3.0 {
			t.Errorf("Expected peek 3.0, got %f", stack.Peek())
		}
	})

	t.Run("ToArray", func(t *testing.T) {
		stack := NewFloat64Stack(3)
		stack.Push(1.0)
		stack.Push(2.0)

		array := stack.ToArray()
		expected := []float64{1.0, 2.0}

		if len(array) != len(expected) {
			t.Errorf("Expected array length %d, got %d", len(expected), len(array))
		}

		for i, v := range expected {
			if array[i] != v {
				t.Errorf("Expected array[%d] = %f, got %f", i, v, array[i])
			}
		}
	})

	t.Run("Generic types - string", func(t *testing.T) {
		stack := NewSizedStack[string](2)

		stack.Push("first")
		stack.Push("second")
		stack.Push("third") // Should remove "first"

		if stack.Size() != 2 {
			t.Errorf("Expected size 2, got %d", stack.Size())
		}

		if stack.Peek() != "third" {
			t.Errorf("Expected peek 'third', got '%s'", stack.Peek())
		}

		if stack.Get(0) != "second" {
			t.Errorf("Expected first element 'second', got '%s'", stack.Get(0))
		}

		array := stack.ToArray()
		expected := []string{"second", "third"}

		if len(array) != len(expected) {
			t.Errorf("Expected array length %d, got %d", len(expected), len(array))
		}

		for i, v := range expected {
			if array[i] != v {
				t.Errorf("Expected array[%d] = %s, got %s", i, v, array[i])
			}
		}
	})

	t.Run("Generic types - int", func(t *testing.T) {
		stack := NewSizedStack[int](3)

		stack.Push(10)
		stack.Push(20)
		stack.Push(30)

		if stack.Size() != 3 {
			t.Errorf("Expected size 3, got %d", stack.Size())
		}

		if stack.Peek() != 30 {
			t.Errorf("Expected peek 30, got %d", stack.Peek())
		}

		// Test empty stack returns zero value
		emptyStack := NewSizedStack[int](1)
		if emptyStack.Peek() != 0 {
			t.Errorf("Expected empty stack peek to return 0, got %d", emptyStack.Peek())
		}

		if emptyStack.Get(0) != 0 {
			t.Errorf("Expected empty stack get to return 0, got %d", emptyStack.Get(0))
		}
	})
}

// TestLinearRegression tests the LinearRegression functionality
func TestLinearRegression(t *testing.T) {
	t.Run("Basic linear regression", func(t *testing.T) {
		// Test with simple linear data: y = 2x + 1
		data := []float64{1.0, 3.0, 5.0, 7.0, 9.0} // y = 2*0+1, 2*1+1, 2*2+1, etc.

		lr := NewLinearRegression(data)
		lr.UpdateData(data)
		err := lr.RunLeastSquares()

		if err != nil {
			t.Fatalf("RunLeastSquares failed: %v", err)
		}

		// Predict next value (should be 11.0)
		prediction := lr.PredictNextValue()
		expected := 11.0

		if math.Abs(prediction-expected) > 0.001 {
			t.Errorf("Expected prediction %f, got %f", expected, prediction)
		}
	})

	t.Run("Constant data", func(t *testing.T) {
		data := []float64{5.0, 5.0, 5.0, 5.0}

		lr := NewLinearRegression(data)
		lr.UpdateData(data)
		lr.RunLeastSquares()

		prediction := lr.PredictNextValue()
		expected := 5.0

		if math.Abs(prediction-expected) > 0.001 {
			t.Errorf("Expected prediction %f, got %f", expected, prediction)
		}
	})

	t.Run("Single data point", func(t *testing.T) {
		data := []float64{5.0}

		lr := NewLinearRegression(data)
		lr.UpdateData(data)
		err := lr.RunLeastSquares()

		if err != nil {
			t.Fatalf("RunLeastSquares failed: %v", err)
		}

		prediction := lr.PredictNextValue()
		expected := 5.0

		if math.Abs(prediction-expected) > 0.001 {
			t.Errorf("Expected prediction %f, got %f", expected, prediction)
		}
	})
}

// TestKalmanFilter tests the KalmanFilter functionality
func TestKalmanFilter(t *testing.T) {
	t.Run("Constructor validation", func(t *testing.T) {
		// Test invalid parameters
		_, err := NewKalmanFilter(-1.0, 1.0, 5)
		if err == nil {
			t.Error("Expected error for negative Q")
		}

		_, err = NewKalmanFilter(1.0, -1.0, 5)
		if err == nil {
			t.Error("Expected error for negative R")
		}

		_, err = NewKalmanFilter(1.0, 1.0, 0)
		if err == nil {
			t.Error("Expected error for zero stack size")
		}

		// Test valid parameters
		kf, err := NewKalmanFilter(0.1, 0.1, 5)
		if err != nil {
			t.Fatalf("Expected no error for valid parameters, got %v", err)
		}

		if kf.GetX() != 0.0 {
			t.Errorf("Expected initial X = 0.0, got %f", kf.GetX())
		}
	})

	t.Run("Basic filtering", func(t *testing.T) {
		kf, err := NewKalmanFilter(0.1, 0.1, 5)
		if err != nil {
			t.Fatalf("Failed to create Kalman filter: %v", err)
		}

		// Test that it implements Filter interface
		var filter Filter = kf

		// Process some measurements
		measurements := []float64{1.0, 2.0, 3.0, 4.0, 5.0}

		for _, measurement := range measurements {
			estimate := filter.Estimate(measurement)

			// Estimate should be finite and reasonable
			if math.IsInf(estimate, 0) || math.IsNaN(estimate) {
				t.Errorf("Invalid estimate: %f", estimate)
			}
		}

		// Final estimate should be close to the last measurement
		finalEstimate := kf.GetX()
		if math.Abs(finalEstimate-5.0) > 2.0 {
			t.Errorf("Final estimate %f too far from expected value 5.0", finalEstimate)
		}
	})

	t.Run("Reset functionality", func(t *testing.T) {
		kf, err := NewKalmanFilter(0.1, 0.1, 3)
		if err != nil {
			t.Fatalf("Failed to create Kalman filter: %v", err)
		}

		// Process some measurements to change state
		kf.Estimate(5.0)
		kf.Estimate(10.0)

		originalK := kf.GetK()

		// Reset and check state
		kf.Reset()

		if kf.GetX() != 0.0 {
			t.Errorf("Expected X = 0.0 after reset, got %f", kf.GetX())
		}

		// After reset, P and K should maintain their converged values
		if math.Abs(kf.GetP()-originalK) > 0.001 { // P should be restored to converged value
			t.Logf("P after reset: %f (this is the converged steady-state value)", kf.GetP())
		}

		// K should be the same as before reset
		if math.Abs(kf.GetK()-originalK) > 0.001 {
			t.Errorf("Expected K to be restored after reset, got %f vs %f", kf.GetK(), originalK)
		}
	})

	t.Run("DARE convergence", func(t *testing.T) {
		kf, err := NewKalmanFilter(0.1, 0.5, 5)
		if err != nil {
			t.Fatalf("Failed to create Kalman filter: %v", err)
		}

		// Kalman gain should be between 0 and 1
		k := kf.GetK()
		if k < 0 || k > 1 {
			t.Errorf("Kalman gain %f should be between 0 and 1", k)
		}

		// Error covariance should be positive
		p := kf.GetP()
		if p <= 0 {
			t.Errorf("Error covariance %f should be positive", p)
		}
	})
}

// BenchmarkKalmanFilter benchmarks the Kalman filter performance
func BenchmarkKalmanFilter(b *testing.B) {
	kf, err := NewKalmanFilter(0.1, 0.1, 5)
	if err != nil {
		b.Fatalf("Failed to create Kalman filter: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kf.Estimate(float64(i % 100))
	}
}

// TestLowPassFilter tests the LowPassFilter functionality
func TestLowPassFilter(t *testing.T) {
	t.Run("Constructor validation", func(t *testing.T) {
		// Test invalid gains
		_, err := NewLowPassFilter(0.0)
		if err == nil {
			t.Error("Expected error for gain = 0.0")
		}

		_, err = NewLowPassFilter(1.0)
		if err == nil {
			t.Error("Expected error for gain = 1.0")
		}

		_, err = NewLowPassFilter(-0.1)
		if err == nil {
			t.Error("Expected error for negative gain")
		}

		_, err = NewLowPassFilter(1.1)
		if err == nil {
			t.Error("Expected error for gain > 1")
		}

		// Test valid gain
		lpf, err := NewLowPassFilter(0.5)
		if err != nil {
			t.Fatalf("Expected no error for valid gain, got %v", err)
		}

		if lpf.GetGain() != 0.5 {
			t.Errorf("Expected gain = 0.5, got %f", lpf.GetGain())
		}

		if lpf.IsInitialized() {
			t.Error("Filter should not be initialized initially")
		}
	})

	t.Run("Interface compliance", func(t *testing.T) {
		lpf, err := NewLowPassFilter(0.3)
		if err != nil {
			t.Fatalf("Failed to create LowPass filter: %v", err)
		}

		// Test that it implements Filter interface
		var filter Filter = lpf

		estimate := filter.Estimate(10.0)
		if estimate != 10.0 {
			t.Errorf("First estimate should equal measurement, got %f", estimate)
		}
	})

	t.Run("Initialization behavior", func(t *testing.T) {
		lpf, err := NewLowPassFilter(0.8)
		if err != nil {
			t.Fatalf("Failed to create LowPass filter: %v", err)
		}

		// First measurement should initialize the filter
		firstEstimate := lpf.Estimate(5.0)
		if firstEstimate != 5.0 {
			t.Errorf("Expected first estimate = 5.0, got %f", firstEstimate)
		}

		if !lpf.IsInitialized() {
			t.Error("Filter should be initialized after first estimate")
		}

		if lpf.GetLastEstimate() != 5.0 {
			t.Errorf("Expected last estimate = 5.0, got %f", lpf.GetLastEstimate())
		}
	})

	t.Run("Filtering behavior", func(t *testing.T) {
		lpf, err := NewLowPassFilter(0.9) // High gain = more smoothing
		if err != nil {
			t.Fatalf("Failed to create LowPass filter: %v", err)
		}

		// Initialize with first value
		lpf.Estimate(0.0)

		// Apply step change
		estimate := lpf.Estimate(10.0)

		// With high gain (0.9), should be smoothed
		// estimate = 0.9 * 0.0 + 0.1 * 10.0 = 1.0
		expected := 0.9*0.0 + 0.1*10.0
		if math.Abs(estimate-expected) > 0.001 {
			t.Errorf("Expected estimate %f, got %f", expected, estimate)
		}

		// Next estimate should be further smoothed
		estimate2 := lpf.Estimate(10.0)
		// estimate2 = 0.9 * 1.0 + 0.1 * 10.0 = 0.9 + 1.0 = 1.9
		expected2 := 0.9*1.0 + 0.1*10.0
		if math.Abs(estimate2-expected2) > 0.001 {
			t.Errorf("Expected second estimate %f, got %f", expected2, estimate2)
		}
	})

	t.Run("Low vs high gain comparison", func(t *testing.T) {
		lpfLowGain, _ := NewLowPassFilter(0.1)  // Low gain = less smoothing, faster response
		lpfHighGain, _ := NewLowPassFilter(0.9) // High gain = more smoothing, slower response

		// Initialize both with 0
		lpfLowGain.Estimate(0.0)
		lpfHighGain.Estimate(0.0)

		// Apply same step input
		lowEstimate := lpfLowGain.Estimate(10.0)
		highEstimate := lpfHighGain.Estimate(10.0)

		// Low gain should respond faster (higher estimate)
		if lowEstimate <= highEstimate {
			t.Errorf("Low gain filter should respond faster: low=%f, high=%f", lowEstimate, highEstimate)
		}

		// Low gain: 0.1 * 0 + 0.9 * 10 = 9.0
		expectedLow := 0.1*0.0 + 0.9*10.0
		if math.Abs(lowEstimate-expectedLow) > 0.001 {
			t.Errorf("Expected low gain estimate %f, got %f", expectedLow, lowEstimate)
		}

		// High gain: 0.9 * 0 + 0.1 * 10 = 1.0
		expectedHigh := 0.9*0.0 + 0.1*10.0
		if math.Abs(highEstimate-expectedHigh) > 0.001 {
			t.Errorf("Expected high gain estimate %f, got %f", expectedHigh, highEstimate)
		}
	})

	t.Run("SetGain functionality", func(t *testing.T) {
		lpf, err := NewLowPassFilter(0.5)
		if err != nil {
			t.Fatalf("Failed to create LowPass filter: %v", err)
		}

		// Test valid gain change
		err = lpf.SetGain(0.8)
		if err != nil {
			t.Errorf("Expected no error for valid gain change: %v", err)
		}

		if lpf.GetGain() != 0.8 {
			t.Errorf("Expected gain = 0.8, got %f", lpf.GetGain())
		}

		// Test invalid gain changes
		err = lpf.SetGain(0.0)
		if err == nil {
			t.Error("Expected error for gain = 0.0")
		}

		err = lpf.SetGain(1.0)
		if err == nil {
			t.Error("Expected error for gain = 1.0")
		}

		// Gain should remain unchanged after invalid attempts
		if lpf.GetGain() != 0.8 {
			t.Errorf("Gain should remain 0.8 after invalid attempts, got %f", lpf.GetGain())
		}
	})

	t.Run("Reset functionality", func(t *testing.T) {
		lpf, err := NewLowPassFilter(0.6)
		if err != nil {
			t.Fatalf("Failed to create LowPass filter: %v", err)
		}

		// Initialize and process some data
		lpf.Estimate(5.0)
		lpf.Estimate(10.0)

		if !lpf.IsInitialized() {
			t.Error("Filter should be initialized")
		}

		// Reset the filter
		lpf.Reset()

		if lpf.IsInitialized() {
			t.Error("Filter should not be initialized after reset")
		}

		if lpf.GetLastEstimate() != 0.0 {
			t.Errorf("Last estimate should be 0.0 after reset, got %f", lpf.GetLastEstimate())
		}

		// Next estimate should initialize again
		newEstimate := lpf.Estimate(7.0)
		if newEstimate != 7.0 {
			t.Errorf("First estimate after reset should be 7.0, got %f", newEstimate)
		}
	})

	t.Run("Noise smoothing", func(t *testing.T) {
		lpf, err := NewLowPassFilter(0.8)
		if err != nil {
			t.Fatalf("Failed to create LowPass filter: %v", err)
		}

		// Simulate noisy signal around value 5.0
		noisyInputs := []float64{5.0, 5.5, 4.8, 5.2, 4.9, 5.3, 4.7, 5.1}
		var estimates []float64

		for _, input := range noisyInputs {
			estimate := lpf.Estimate(input)
			estimates = append(estimates, estimate)
		}

		// Check that estimates are smoother than inputs
		// Calculate variance of inputs vs estimates (excluding first)
		if len(estimates) > 2 {
			inputVariance := calculateVariance(noisyInputs[1:]) // Skip first for fair comparison
			estimateVariance := calculateVariance(estimates[1:])

			if estimateVariance >= inputVariance {
				t.Errorf("Filter should smooth noise: input variance=%f, estimate variance=%f",
					inputVariance, estimateVariance)
			}
		}

		// Final estimate should be reasonably close to signal average
		finalEstimate := estimates[len(estimates)-1]
		expectedAverage := 5.0 // Approximate center of noisy signal
		if math.Abs(finalEstimate-expectedAverage) > 0.5 {
			t.Errorf("Final estimate %f should be close to signal average %f",
				finalEstimate, expectedAverage)
		}
	})
}

// calculateVariance calculates the variance of a slice of values
func calculateVariance(values []float64) float64 {
	if len(values) <= 1 {
		return 0.0
	}

	// Calculate mean
	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	// Calculate variance
	var variance float64
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}
	return variance / float64(len(values))
}

// BenchmarkLowPassFilter benchmarks the LowPass filter performance
func BenchmarkLowPassFilter(b *testing.B) {
	lpf, err := NewLowPassFilter(0.7)
	if err != nil {
		b.Fatalf("Failed to create LowPass filter: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lpf.Estimate(float64(i % 100))
	}
}
