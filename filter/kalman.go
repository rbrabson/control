// Package filter provides signal filtering implementations for control systems.
//
// This package includes advanced filtering algorithms such as Kalman filters
// that can be used to estimate system states and reduce measurement noise.
package filter

import (
	"errors"
)

// Float64Stack is a type alias for SizedStack with float64 elements for backward compatibility.
type Float64Stack = SizedStack[float64]

// NewFloat64Stack creates a new float64 sized stack with the given capacity.
// This function maintains backward compatibility with existing code.
func NewFloat64Stack(capacity int) *Float64Stack {
	return NewSizedStack[float64](capacity)
}

// KalmanFilter implements a Kalman filter that uses least squares regression as its model.
type KalmanFilter struct {
	q          float64           // Sensor covariance
	r          float64           // Model covariance
	n          int               // Number of elements in the stack
	p          float64           // Error covariance estimate
	k          float64           // Kalman gain
	x          float64           // State estimate
	estimates  *Float64Stack     // Stack of recent estimates
	regression *LinearRegression // Linear regression for prediction
}

// NewKalmanFilter creates a new Kalman filter.
//
// Parameters:
//   - q: Sensor covariance (measurement noise)
//   - r: Model covariance (process noise)
//   - n: Number of elements to maintain in the estimates stack
func NewKalmanFilter(q, r float64, n int) (*KalmanFilter, error) {
	if n <= 0 {
		return nil, errors.New("stack size must be positive")
	}
	if q < 0 || r < 0 {
		return nil, errors.New("covariance values must be non-negative")
	}

	kf := &KalmanFilter{
		q:         q,
		r:         r,
		n:         n,
		p:         1.0,
		k:         0.0,
		x:         0.0,
		estimates: NewFloat64Stack(n),
	}

	// Initialize stack with zeros
	for i := 0; i < n; i++ {
		kf.estimates.Push(0.0)
	}

	// Initialize regression with the zero-filled stack
	kf.regression = NewLinearRegression(kf.estimates.ToArray())

	// Calculate initial Kalman gain using DARE
	kf.findK()

	return kf, nil
}

// SetX sets the current state estimate.
func (kf *KalmanFilter) SetX(x float64) {
	kf.x = x
}

// GetX returns the current state estimate.
func (kf *KalmanFilter) GetX() float64 {
	return kf.x
}

// GetK returns the current Kalman gain.
func (kf *KalmanFilter) GetK() float64 {
	return kf.k
}

// GetP returns the current error covariance estimate.
func (kf *KalmanFilter) GetP() float64 {
	return kf.p
}

// Estimate processes a measurement and returns the optimal state estimate.
// This implements the Filter interface.
func (kf *KalmanFilter) Estimate(measurement float64) float64 {
	// Run regression to get prediction
	kf.regression.RunLeastSquares()

	// Update state estimate using regression prediction
	prediction := kf.regression.PredictNextValue()
	kf.x += prediction - kf.estimates.Peek()

	// Apply Kalman filter update
	kf.x += kf.k * (measurement - kf.x)

	// Store new estimate
	kf.estimates.Push(kf.x)

	// Update regression with new data
	kf.regression.UpdateData(kf.estimates.ToArray())

	return kf.x
}

// findK iteratively computes the Kalman gain using the DARE.
func (kf *KalmanFilter) findK() {
	// Run 2000 iterations to converge to steady-state solution
	for i := 0; i < 2000; i++ {
		kf.solveDARE()
	}
}

// solveDARE solves the Discrete Algebraic Riccati Equation (DARE).
// This implements one iteration of the DARE solution for steady-state Kalman gain.
func (kf *KalmanFilter) solveDARE() {
	// Prediction step: P = P + Q
	kf.p += kf.q

	// Calculate Kalman gain: K = P / (P + R)
	kf.k = kf.p / (kf.p + kf.r)

	// Update step: P = (1 - K) * P
	kf.p = (1 - kf.k) * kf.p
}

// Reset resets the filter to its initial state.
func (kf *KalmanFilter) Reset() {
	// Store the converged values
	convergedP := kf.p
	convergedK := kf.k

	// Reset state estimate
	kf.x = 0.0

	// Reinitialize stack with zeros
	kf.estimates = NewFloat64Stack(kf.n)
	for i := 0; i < kf.n; i++ {
		kf.estimates.Push(0.0)
	}

	// Reinitialize regression
	kf.regression = NewLinearRegression(kf.estimates.ToArray())

	// Restore converged values (they should be the same for the same Q,R parameters)
	kf.p = convergedP
	kf.k = convergedK
}
