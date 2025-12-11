// Package main demonstrates a comparison between Kalman and LowPass filters in PID control.
//
// This example shows how different filter types affect PID controller performance
// in noisy environments, demonstrating trade-offs between filtering effectiveness
// and computational cost.
package main

import (
	"fmt"
	"math"
	"math/rand"

	"control/filter"
	"control/pid"
)

// SystemSimulator represents a first-order system with noise
type SystemSimulator struct {
	currentState   float64
	timeConstant   float64
	noiseMagnitude float64
}

// NewSystemSimulator creates a new system simulator
func NewSystemSimulator(initialState, timeConstant, noiseMagnitude float64) *SystemSimulator {
	return &SystemSimulator{
		currentState:   initialState,
		timeConstant:   timeConstant,
		noiseMagnitude: noiseMagnitude,
	}
}

// Update simulates one time step of the system
// Uses first-order dynamics: dx/dt = (u - x) / tau
func (s *SystemSimulator) Update(controlInput, dt float64) float64 {
	// System dynamics
	derivative := (controlInput - s.currentState) / s.timeConstant
	s.currentState += derivative * dt

	return s.currentState
}

// GetMeasurement returns the current state with added measurement noise
func (s *SystemSimulator) GetMeasurement() float64 {
	noise := (rand.Float64() - 0.5) * 2 * s.noiseMagnitude
	return s.currentState + noise
}

// ControllerMetrics tracks controller performance
type ControllerMetrics struct {
	name         string
	totalError   float64
	maxError     float64
	outputs      []float64
	measurements []float64
	filterType   string
}

// NewControllerMetrics creates a new metrics tracker
func NewControllerMetrics(name, filterType string) *ControllerMetrics {
	return &ControllerMetrics{
		name:         name,
		filterType:   filterType,
		outputs:      make([]float64, 0),
		measurements: make([]float64, 0),
	}
}

// UpdateMetrics records performance data
func (m *ControllerMetrics) UpdateMetrics(error, output, measurement float64) {
	m.totalError += math.Abs(error)
	if math.Abs(error) > m.maxError {
		m.maxError = math.Abs(error)
	}
	m.outputs = append(m.outputs, output)
	m.measurements = append(m.measurements, measurement)
}

// CalculateVariance computes output variance (smoothness metric)
func (m *ControllerMetrics) CalculateVariance() float64 {
	if len(m.outputs) == 0 {
		return 0
	}

	var sum float64
	for _, v := range m.outputs {
		sum += v
	}
	mean := sum / float64(len(m.outputs))

	var variance float64
	for _, v := range m.outputs {
		diff := v - mean
		variance += diff * diff
	}
	return variance / float64(len(m.outputs))
}

// PrintResults displays the metrics
func (m *ControllerMetrics) PrintResults() {
	fmt.Printf("\n%s (%s Filter)\n", m.name, m.filterType)
	fmt.Println("----------------------------------------")
	fmt.Printf("Total Absolute Error:    %.4f\n", m.totalError)
	fmt.Printf("Maximum Error:           %.4f\n", m.maxError)
	fmt.Printf("Output Variance:         %.4f (lower is smoother)\n", m.CalculateVariance())
	fmt.Printf("Average Output:          %.4f\n", calculateMean(m.outputs))
}

// calculateMean computes the mean of a slice
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// calculateStdDev computes the standard deviation
func calculateStdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	mean := calculateMean(values)
	var variance float64
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))
	return math.Sqrt(variance)
}

func main() {
	fmt.Println("Kalman Filter vs LowPass Filter - PID Control Comparison")
	fmt.Println("=======================================================")

	// Simulation parameters
	setpoint := 50.0
	duration := 10.0 // Simulation duration in seconds
	dt := 0.01       // Time step in seconds
	iterations := int(duration / dt)

	// System parameters
	initialState := 10.0
	timeConstant := 1.0
	noiseMagnitude := 2.0

	// Create controllers with different filters
	// LowPass filter controller
	lpFilter, _ := filter.NewLowPassFilter(0.3)
	lpPID := pid.New(0.5, 0.1, 0.05,
		pid.WithFilter(lpFilter),
		pid.WithOutputLimits(-100.0, 100.0),
	)
	lpMetrics := NewControllerMetrics("LowPass Filter Controller", "LowPass")

	// Kalman filter controller (sensor covariance 0.05, model covariance 0.1, 10 states)
	kfFilter, _ := filter.NewKalmanFilter(0.05, 0.1, 10)
	kfPID := pid.New(0.5, 0.1, 0.05,
		pid.WithFilter(kfFilter),
		pid.WithOutputLimits(-100.0, 100.0),
	)
	kfMetrics := NewControllerMetrics("Kalman Filter Controller", "Kalman")

	// No filter controller for comparison
	nfPID := pid.New(0.5, 0.1, 0.05,
		pid.WithOutputLimits(-100.0, 100.0),
	)
	nfMetrics := NewControllerMetrics("No Filter Controller", "None")

	// Create systems for each controller
	lpSystem := NewSystemSimulator(initialState, timeConstant, noiseMagnitude)
	kfSystem := NewSystemSimulator(initialState, timeConstant, noiseMagnitude)
	nfSystem := NewSystemSimulator(initialState, timeConstant, noiseMagnitude)

	fmt.Printf("Simulation Parameters:\n")
	fmt.Printf("  Setpoint:          %.1f\n", setpoint)
	fmt.Printf("  Duration:          %.1f seconds\n", duration)
	fmt.Printf("  Time Step:         %.3f seconds\n", dt)
	fmt.Printf("  System Time Const: %.1f seconds\n", timeConstant)
	fmt.Printf("  Noise Magnitude:   %.1f\n\n", noiseMagnitude)

	// Run simulation
	for i := 0; i < iterations; i++ {
		// Update LowPass filter system
		lpOutput := lpPID.Calculate(setpoint, lpSystem.GetMeasurement())
		lpSystem.Update(lpOutput, dt)
		lpError := setpoint - lpSystem.currentState
		lpMetrics.UpdateMetrics(lpError, lpOutput, lpSystem.GetMeasurement())

		// Update Kalman filter system
		kfOutput := kfPID.Calculate(setpoint, kfSystem.GetMeasurement())
		kfSystem.Update(kfOutput, dt)
		kfError := setpoint - kfSystem.currentState
		kfMetrics.UpdateMetrics(kfError, kfOutput, kfSystem.GetMeasurement())

		// Update No filter system
		nfOutput := nfPID.Calculate(setpoint, nfSystem.GetMeasurement())
		nfSystem.Update(nfOutput, dt)
		nfError := setpoint - nfSystem.currentState
		nfMetrics.UpdateMetrics(nfError, nfOutput, nfSystem.GetMeasurement())
	}

	// Print results
	lpMetrics.PrintResults()
	kfMetrics.PrintResults()
	nfMetrics.PrintResults()

	// Comparative analysis
	fmt.Println("\nComparative Analysis")
	fmt.Println("-------------------")

	lpVar := lpMetrics.CalculateVariance()
	kfVar := kfMetrics.CalculateVariance()
	nfVar := nfMetrics.CalculateVariance()

	fmt.Printf("\nOutput Smoothness (Variance):\n")
	fmt.Printf("  LowPass:  %.4f\n", lpVar)
	fmt.Printf("  Kalman:   %.4f\n", kfVar)
	fmt.Printf("  None:     %.4f\n", nfVar)

	if kfVar < lpVar {
		reduction := ((lpVar - kfVar) / lpVar) * 100
		fmt.Printf("\n✓ Kalman filter is %.1f%% smoother than LowPass\n", reduction)
	} else {
		reduction := ((kfVar - lpVar) / kfVar) * 100
		fmt.Printf("\n✓ LowPass filter is %.1f%% smoother than Kalman\n", reduction)
	}

	fmt.Printf("\nError Tracking (Lower is Better):\n")
	fmt.Printf("  LowPass:  %.4f\n", lpMetrics.totalError)
	fmt.Printf("  Kalman:   %.4f\n", kfMetrics.totalError)
	fmt.Printf("  None:     %.4f\n", nfMetrics.totalError)

	fmt.Printf("\nFinal State Values (Should approach %.1f):\n", setpoint)
	fmt.Printf("  LowPass:  %.4f\n", lpSystem.currentState)
	fmt.Printf("  Kalman:   %.4f\n", kfSystem.currentState)
	fmt.Printf("  None:     %.4f\n", nfSystem.currentState)

	// Performance summary
	fmt.Println("\nPerformance Summary")
	fmt.Println("-------------------")

	lpStdDev := calculateStdDev(lpMetrics.outputs)
	kfStdDev := calculateStdDev(kfMetrics.outputs)
	nfStdDev := calculateStdDev(nfMetrics.outputs)

	fmt.Printf("\nOutput Standard Deviation (Control Effort Consistency):\n")
	fmt.Printf("  LowPass:  %.6f\n", lpStdDev)
	fmt.Printf("  Kalman:   %.6f\n", kfStdDev)
	fmt.Printf("  None:     %.6f\n", nfStdDev)

	fmt.Println("\nKey Observations:")
	fmt.Println("-----------------")
	fmt.Println("• Kalman Filter: Advanced state estimation, adapts to system dynamics")
	fmt.Println("• LowPass Filter: Simple, computationally efficient, fixed smoothing")
	fmt.Println("• No Filter: Most responsive but noisiest output")
	fmt.Println("\nChoose based on your application:")
	fmt.Println("- Use Kalman when you have good knowledge of system dynamics")
	fmt.Println("- Use LowPass for simplicity and predictable behavior")
	fmt.Println("- Use No Filter when response speed is critical")
}
