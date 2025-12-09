package filter

import (
	"errors"
)

// LowPassFilter implements a simple first-order low-pass filter.
//
// The filter uses the formula: estimate = gain * previousEstimate + (1-gain) * measurement
// where gain is between 0 and 1.
//
// High values of gain (closer to 1) are smoother but have more phase lag.
// Low values of gain (closer to 0) allow more noise but respond faster to changes.
type LowPassFilter struct {
	gain             float64 // Filter gain (0 < gain < 1)
	previousEstimate float64 // Previous filtered value
	initialized      bool    // Whether the filter has been initialized
}

// NewLowPassFilter creates a new low-pass filter with the specified gain.
//
// Parameters:
//   - gain: Filter gain (0 < gain < 1). Higher values = smoother but more lag.
//
// Returns an error if gain is not in the valid range (0, 1).
func NewLowPassFilter(gain float64) (*LowPassFilter, error) {
	if gain <= 0 || gain >= 1 {
		return nil, errors.New("gain must be between 0 and 1 (exclusive)")
	}

	return &LowPassFilter{
		gain:             gain,
		previousEstimate: 0.0,
		initialized:      false,
	}, nil
}

// Estimate processes a measurement through the low-pass filter.
// This implements the Filter interface.
//
// On the first call, the filter is initialized with the measurement value.
// Subsequent calls apply the low-pass filtering formula.
func (lpf *LowPassFilter) Estimate(measurement float64) float64 {
	if !lpf.initialized {
		// Initialize with first measurement
		lpf.previousEstimate = measurement
		lpf.initialized = true
		return measurement
	}

	// Apply low-pass filter: estimate = gain * previous + (1-gain) * measurement
	estimate := lpf.gain*lpf.previousEstimate + (1-lpf.gain)*measurement
	lpf.previousEstimate = estimate
	return estimate
}

// GetGain returns the current filter gain.
func (lpf *LowPassFilter) GetGain() float64 {
	return lpf.gain
}

// SetGain updates the filter gain.
// Returns an error if the new gain is not in the valid range (0, 1).
func (lpf *LowPassFilter) SetGain(gain float64) error {
	if gain <= 0 || gain >= 1 {
		return errors.New("gain must be between 0 and 1 (exclusive)")
	}
	lpf.gain = gain
	return nil
}

// Reset resets the filter to its uninitialized state.
// The next call to Estimate will initialize the filter with the provided measurement.
func (lpf *LowPassFilter) Reset() {
	lpf.previousEstimate = 0.0
	lpf.initialized = false
}

// GetLastEstimate returns the last filtered value.
// Returns 0.0 if the filter hasn't been initialized yet.
func (lpf *LowPassFilter) GetLastEstimate() float64 {
	if !lpf.initialized {
		return 0.0
	}
	return lpf.previousEstimate
}

// IsInitialized returns whether the filter has processed at least one measurement.
func (lpf *LowPassFilter) IsInitialized() bool {
	return lpf.initialized
}
