package filter

import "math"

// LinearRegression provides least squares linear regression functionality.
type LinearRegression struct {
	data      []float64
	slope     float64
	intercept float64
	hasRun    bool
}

// NewLinearRegression creates a new linear regression with the given data.
func NewLinearRegression(data []float64) *LinearRegression {
	return &LinearRegression{
		data: make([]float64, len(data)),
	}
}

// RunLeastSquares performs least squares regression on the data.
// Uses the formula: y = mx + b where x represents time indices.
func (lr *LinearRegression) RunLeastSquares() error {
	n := len(lr.data)
	if n < 2 {
		lr.slope = 0
		lr.intercept = 0
		lr.hasRun = true
		return nil
	}

	// Calculate sums for least squares formula
	var sumX, sumY, sumXY, sumXX float64

	for i, y := range lr.data {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}

	nf := float64(n)

	// Calculate slope and intercept
	denominator := nf*sumXX - sumX*sumX
	if math.Abs(denominator) < 1e-10 {
		lr.slope = 0
		if n > 0 {
			lr.intercept = sumY / nf
		} else {
			lr.intercept = 0
		}
	} else {
		lr.slope = (nf*sumXY - sumX*sumY) / denominator
		lr.intercept = (sumY - lr.slope*sumX) / nf
	}

	lr.hasRun = true
	return nil
}

// PredictNextValue predicts the next value in the sequence based on the linear regression.
func (lr *LinearRegression) PredictNextValue() float64 {
	if !lr.hasRun {
		lr.RunLeastSquares()
	}

	// For single data point, return that value
	if len(lr.data) == 1 {
		return lr.data[0]
	}

	nextX := float64(len(lr.data))
	return lr.slope*nextX + lr.intercept
}

// UpdateData updates the regression data and marks it as needing to be recalculated.
func (lr *LinearRegression) UpdateData(data []float64) {
	lr.data = make([]float64, len(data))
	copy(lr.data, data)
	lr.hasRun = false
}
