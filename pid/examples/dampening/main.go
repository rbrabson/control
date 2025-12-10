// Package main demonstrates PID dampening features for noisy control systems.
//
// This example shows how derivative filtering and stability threshold can improve
// control performance in real-world scenarios with measurement noise and rapid
// system changes.
package main

import (
	"fmt"
	"math"
	"math/rand"

	"control/filter"
	"control/pid"
)

// simulateNoisyMeasurement adds realistic sensor noise to a measurement
func simulateNoisyMeasurement(actual float64, noiseLevel float64) float64 {
	noise := (rand.Float64() - 0.5) * 2 * noiseLevel
	return actual + noise
}

// simulateSystemResponse models a simple first-order system response
func simulateSystemResponse(currentValue, controlOutput, timeStep float64) float64 {
	// Simple exponential approach: x(t+1) = x(t) + (output * timeStep)
	responseRate := 0.1 // How quickly system responds to control
	return currentValue + (controlOutput * responseRate * timeStep)
}

func main() {
	fmt.Println("PID Dampening Features Demonstration")
	fmt.Println("===================================")

	// Simulation parameters
	targetValue := 100.0
	initialValue := 20.0
	noiseLevel := 2.0 // Sensor noise amplitude
	timeStep := 0.1   // Time step in seconds
	duration := 20.0  // Simulation duration in seconds

	// Create three PID controllers for comparison
	basicPID := pid.New(2.0, 0.5, 0.8,
		pid.WithOutputLimits(-50.0, 50.0),
	)

	f, _ := filter.NewLowPassFilter(0.4)
	filteredPID := pid.New(2.0, 0.5, 0.8,
		pid.WithFilter(f), // Low-pass filter derivative
		pid.WithOutputLimits(-50.0, 50.0),
	)

	f2, _ := filter.NewLowPassFilter(0.4)
	dampedPID := pid.New(2.0, 0.5, 0.8,
		pid.WithFilter(f2),              // Filter derivative noise
		pid.WithStabilityThreshold(3.0), // Disable integral during instability
		pid.WithOutputLimits(-50.0, 50.0),
	)

	// Initialize system states
	basicValue := initialValue
	filteredValue := initialValue
	dampedValue := initialValue

	fmt.Printf("Target: %.1f, Initial: %.1f, Noise: ±%.1f\n", targetValue, initialValue, noiseLevel)
	fmt.Println("Time\tBasic\tFiltered\tDamped\tB_Out\tF_Out\tD_Out")
	fmt.Println("----\t-----\t--------\t------\t-----\t-----\t-----")

	steps := int(duration / timeStep)

	for i := 0; i <= steps; i++ {
		currentTime := float64(i) * timeStep

		// Add measurement noise to all controllers
		basicMeasurement := simulateNoisyMeasurement(basicValue, noiseLevel)
		filteredMeasurement := simulateNoisyMeasurement(filteredValue, noiseLevel)
		dampedMeasurement := simulateNoisyMeasurement(dampedValue, noiseLevel)

		// Calculate control outputs
		basicOutput := basicPID.Calculate(targetValue, basicMeasurement)
		filteredOutput := filteredPID.Calculate(targetValue, filteredMeasurement)
		dampedOutput := dampedPID.Calculate(targetValue, dampedMeasurement)

		// Simulate system responses
		basicValue = simulateSystemResponse(basicValue, basicOutput, timeStep)
		filteredValue = simulateSystemResponse(filteredValue, filteredOutput, timeStep)
		dampedValue = simulateSystemResponse(dampedValue, dampedOutput, timeStep)

		// Print results every second
		if i%int(1.0/timeStep) == 0 {
			fmt.Printf("%.1f\t%.1f\t%.1f\t\t%.1f\t%.1f\t%.1f\t%.1f\n",
				currentTime, basicValue, filteredValue, dampedValue,
				basicOutput, filteredOutput, dampedOutput)
		}
	}

	fmt.Println()
	fmt.Println("Performance Analysis:")
	fmt.Println("====================")

	// Calculate final errors
	basicError := math.Abs(basicValue - targetValue)
	filteredError := math.Abs(filteredValue - targetValue)
	dampedError := math.Abs(dampedValue - targetValue)

	fmt.Printf("Final Errors:\n")
	fmt.Printf("  Basic PID:      %.2f (%.1f%% of target)\n", basicError, basicError/targetValue*100)
	fmt.Printf("  Filtered PID:   %.2f (%.1f%% of target)\n", filteredError, filteredError/targetValue*100)
	fmt.Printf("  Damped PID:     %.2f (%.1f%% of target)\n", dampedError, dampedError/targetValue*100)

	fmt.Println()
	fmt.Println("Dampening Features Explained:")
	fmt.Println("============================")
	filterMsg := "Derivative Filter: Applies low-pass filtering to derivative term"
	if filteredPID.GetFilter() != nil {
		filterMsg = "Derivative Filter (enabled): Applies low-pass filtering to derivative term"
	}
	fmt.Println(filterMsg)
	fmt.Println("  - Reduces impact of high-frequency measurement noise")
	fmt.Println("  - Smooths control output for more stable system response")
	fmt.Println("  - Range: 0.0 (no filtering) to 1.0 (maximum filtering)")

	fmt.Printf("\nStability Threshold (%.1f): Disables integral when derivative exceeds threshold\n", dampedPID.GetStabilityThreshold())
	fmt.Println("  - Prevents integral windup during rapid system changes")
	fmt.Println("  - Automatically re-enables when system stabilizes")
	fmt.Println("  - Helps maintain control stability in dynamic conditions")

	fmt.Println()
	fmt.Println("When to Use Dampening:")
	fmt.Println("======================")
	fmt.Println("• Derivative Filter:")
	fmt.Println("  - Noisy measurement sensors")
	fmt.Println("  - High-frequency disturbances")
	fmt.Println("  - Jerky control outputs")
	fmt.Println()
	fmt.Println("• Stability Threshold:")
	fmt.Println("  - Systems with rapid setpoint changes")
	fmt.Println("  - Aggressive tuning prone to overshoot")
	fmt.Println("  - Variable system dynamics")
}
