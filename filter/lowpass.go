package filter

import (
	"errors"
)

// LowPassFilter implements a simple first-order low-pass filter.
//
// The filter uses the formula: estimate = alpha * previousEstimate + (1-alpha) * measurement
// where alpha is between 0 and 1.
//
// High values of alpha (closer to 1) are smoother but have more phase lag.
// Low values of alpha (closer to 0) allow more noise but respond faster to changes.
type LowPassFilter struct {
	alpha            float64 // Filter alpha (0 < alpha < 1)
	previousEstimate float64 // Previous filtered value
	initialized      bool    // Whether the filter has been initialized
}

// NewLowPassFilter creates a new low-pass filter with the specified alpha.
//
// Parameters:
//   - alpha: Filter alpha (0 < alpha < 1). Higher values = smoother but more lag.
//
// Returns an error if alpha is not in the valid range (0, 1).
func NewLowPassFilter(alpha float64) (*LowPassFilter, error) {
	if alpha <= 0 || alpha >= 1 {
		return nil, errors.New("alpha must be between 0 and 1 (exclusive)")
	}

	return &LowPassFilter{
		alpha:            alpha,
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

	// Apply low-pass filter: estimate = alpha * previous + (1-alpha) * measurement
	estimate := lpf.alpha*lpf.previousEstimate + (1-lpf.alpha)*measurement
	lpf.previousEstimate = estimate
	return estimate
}

// GetAlpha returns the current filter alpha.
func (lpf *LowPassFilter) GetAlpha() float64 {
	return lpf.alpha
}

// GetGain returns the current filter alpha (alias for GetAlpha).
// This method exists to satisfy the Filter interface.
func (lpf *LowPassFilter) GetGain() float64 {
	return lpf.alpha
}

// SetAlpha updates the filter alpha.
// Returns an error if the new alpha is not in the valid range (0, 1).
func (lpf *LowPassFilter) SetAlpha(alpha float64) error {
	if alpha <= 0 || alpha >= 1 {
		return errors.New("alpha must be between 0 and 1 (exclusive)")
	}
	lpf.alpha = alpha
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
