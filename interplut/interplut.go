package interplut

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// InterpLUT performs spline interpolation given a set of control points.
// It creates a monotone cubic spline that is guaranteed to pass through
// each control point exactly. The spline maintains monotonicity if the
// control points are monotonic.
type InterpLUT struct {
	x []float64 // X coordinates (inputs)
	y []float64 // Y coordinates (outputs)
	m []float64 // Tangent slopes at each point
}

// New creates a new empty InterpLUT.
func New() *InterpLUT {
	return &InterpLUT{
		x: make([]float64, 0),
		y: make([]float64, 0),
		m: make([]float64, 0),
	}
}

// Add adds a control point to the lookup table.
// input: The X coordinate of the control point
// output: The Y coordinate of the control point
func (lut *InterpLUT) Add(input, output float64) {
	lut.x = append(lut.x, input)
	lut.y = append(lut.y, output)
}

// CreateLUT creates the monotone cubic spline from the added control points.
// This must be called after adding all control points and before using Get().
//
// Returns an error if there are fewer than 2 control points or if the
// X values are not strictly increasing after sorting.
func (lut *InterpLUT) CreateLUT() error {
	if len(lut.x) != len(lut.y) || len(lut.x) < 2 {
		return fmt.Errorf("there must be at least two control points and the arrays must be of equal length")
	}

	// Create pairs and sort by X values
	type point struct {
		x, y float64
	}
	points := make([]point, len(lut.x))
	for i := range lut.x {
		points[i] = point{lut.x[i], lut.y[i]}
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].x < points[j].x
	})

	// Extract sorted x and y values
	n := len(points)
	x := make([]float64, n)
	y := make([]float64, n)
	for i, p := range points {
		x[i] = p.x
		y[i] = p.y
	}

	// Ensure strictly increasing X values
	sort.Slice(x, func(i, j int) bool {
		// Check for sorting
		cond := x[i] < x[j]
		if cond {
			y[i], y[j] = y[j], y[i]
		}
		return cond
	})

	// Check for duplicate X values
	for i := 0; i < n-1; i++ {
		if x[i] == x[i+1] {
			return fmt.Errorf("the control points have duplicate X values")
		}
	}

	// Compute slopes of secant lines between successive points
	d := make([]float64, n-1)
	for i := 0; i < n-1; i++ {
		h := x[i+1] - x[i]
		d[i] = (y[i+1] - y[i]) / h
	}

	// Initialize the tangents as the average of the secants
	m := make([]float64, n)
	m[0] = d[0]
	for i := 1; i < n-1; i++ {
		m[i] = (d[i-1] + d[i]) * 0.5
	}
	m[n-1] = d[n-2]

	// Update the tangents to preserve monotonicity
	for i := 0; i < n-1; i++ {
		if d[i] == 0.0 { // successive Y values are equal
			m[i] = 0.0
			m[i+1] = 0.0
		} else {
			a := m[i] / d[i]
			b := m[i+1] / d[i]
			h := math.Hypot(a, b)
			if h > 9.0 {
				t := 3.0 / h
				m[i] = t * a * d[i]
				m[i+1] = t * b * d[i]
			}
		}
	}

	lut.x = x
	lut.y = y
	lut.m = m

	return nil
}

// Get interpolates the value of Y = f(X) for given X.
// The input is clamped to the domain of the spline.
//
// Returns an error if the input is outside the bounds of the lookup table
// or if CreateLUT() has not been called.
func (lut *InterpLUT) Get(input float64) (float64, error) {
	n := len(lut.x)
	if n == 0 {
		return 0, fmt.Errorf("CreateLUT() must be called before get()")
	}

	// Handle NaN input
	if math.IsNaN(input) {
		return input, nil
	}

	// Handle boundary cases
	if input < lut.x[0] || input > lut.x[n-1] {
		return 0, fmt.Errorf("user requested value outside of bounds of LUT. Bounds are: %f to %f. Value provided was: %f",
			lut.x[0], lut.x[n-1], input)
	}

	// Find the index 'i' of the last point with smaller X
	i := 0
	for input >= lut.x[i+1] && i < n-2 {
		i++
	}

	// Check for exact match
	if input == lut.x[i] {
		return lut.y[i], nil
	}

	// Perform cubic Hermite spline interpolation
	h := lut.x[i+1] - lut.x[i]
	t := (input - lut.x[i]) / h

	result := (lut.y[i]*(1+2*t)+h*lut.m[i]*t)*(1-t)*(1-t) +
		(lut.y[i+1]*(3-2*t)+h*lut.m[i+1]*(t-1))*t*t

	return result, nil
}

// String returns a string representation of the InterpLUT for debugging.
func (lut *InterpLUT) String() string {
	var str strings.Builder
	str.WriteString("[")
	n := len(lut.x)
	for i := 0; i < n; i++ {
		if i != 0 {
			str.WriteString(", ")
		}
		str.WriteString(fmt.Sprintf("(%g, %g: %g)", lut.x[i], lut.y[i], lut.m[i]))
	}
	str.WriteString("]")
	return str.String()
}
