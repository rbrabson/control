// Package main demonstrates comparing different derivative filter options in PID.
//
// This example compares no filter, low-pass filter, and Kalman filter
// to show the effect on control output smoothness.
package main

import (
	"fmt"

	"control/filter"
	"control/pid"
)

func main() {
	fmt.Println("PID Filter Comparison Example")
	fmt.Println("============================")

	fmt.Println("Comparing Derivative Filter Types")
	fmt.Println("---------------------------------")
	fmt.Println("Demonstrates differences between no filter,")
	fmt.Println("low-pass filter, and Kalman filter.")
	fmt.Println()

	// Create controllers with different filters
	pidNoFilter := pid.New(1.0, 0.2, 0.8)

	lowpassFilter, _ := filter.NewLowPassFilter(0.5)
	pidLowpass := pid.New(1.0, 0.2, 0.8,
		pid.WithFilter(lowpassFilter))

	kalmanFilter, _ := filter.NewKalmanFilter(0.1, 1.0, 5)
	pidKalman := pid.New(1.0, 0.2, 0.8,
		pid.WithFilter(kalmanFilter))

	fmt.Println("Controller Configuration:")
	fmt.Println("  All: Kp=1.0, Ki=0.2, Kd=0.8")
	fmt.Println("  No Filter: Raw derivative")
	fmt.Println("  Low-Pass: alpha=0.5")
	fmt.Println("  Kalman: process_var=0.1, measurement_var=1.0")
	fmt.Println()

	// Test with noisy state measurements
	fmt.Println("Test: Approaching Setpoint with Measurement Noise")
	fmt.Printf("\n%-6s %-10s %-8s %-11s %-11s %-11s\n", "Step", "Setpoint", "State", "No Filter", "Low-Pass", "Kalman")
	fmt.Printf("%-6s %-10s %-8s %-11s %-11s %-11s\n", "----", "--------", "-----", "---------", "--------", "------")

	setpoint := 100.0
	dt := 0.1

	// Simulated noisy measurements (ideal state with noise)
	noisyStates := []float64{
		0.0,   // Start
		18.5,  // Should be 20, but has noise
		41.2,  // Should be 40
		58.8,  // Should be 60
		81.5,  // Should be 80
		89.2,  // Should be 90
		94.5,  // Should be 95
		97.8,  // Should be 98
		99.3,  // Should be 99
		100.2, // Should be 100
	}

	for i, state := range noisyStates {
		outNoFilter := pidNoFilter.CalculateWithDt(setpoint, state, dt)
		outLowpass := pidLowpass.CalculateWithDt(setpoint, state, dt)
		outKalman := pidKalman.CalculateWithDt(setpoint, state, dt)

		fmt.Printf("%-6d %-10.0f %-8.1f %-11.3f %-11.3f %-11.3f\n",
			i+1, setpoint, state, outNoFilter, outLowpass, outKalman)
	}

	fmt.Println("\nOutput Smoothness Comparison:")
	fmt.Println("-----------------------------")

	// Compare derivative response to sudden change
	pidNoFilter.Reset()
	pidLowpass.Reset()
	pidKalman.Reset()

	fmt.Println("\nSudden State Jump (0 → 50):")
	state1 := 0.0
	state2 := 50.0

	pidNoFilter.Calculate(setpoint, state1)
	pidLowpass.Calculate(setpoint, state1)
	pidKalman.Calculate(setpoint, state1)

	outNo := pidNoFilter.CalculateWithDt(setpoint, state2, dt)
	outLow := pidLowpass.CalculateWithDt(setpoint, state2, dt)
	outKal := pidKalman.CalculateWithDt(setpoint, state2, dt)

	fmt.Printf("No filter output: %.3f (raw derivative spike)\n", outNo)
	fmt.Printf("Low-pass output:  %.3f (first-step response)\n", outLow)
	fmt.Printf("Kalman output:    %.3f (predicted and filtered)\n", outKal)

	fmt.Println("\nAnalysis:")
	fmt.Println("---------")
	fmt.Println("No Filter:")
	fmt.Println("  • Most responsive to changes")
	fmt.Println("  • Amplifies measurement noise")
	fmt.Println("  • Large derivative kicks")
	fmt.Println()
	fmt.Println("Low-Pass Filter:")
	fmt.Println("  • Simple and computationally cheap")
	fmt.Println("  • Good noise reduction")
	fmt.Println("  • Introduces phase lag")
	fmt.Println()
	fmt.Println("Kalman Filter:")
	fmt.Println("  • Optimal for Gaussian noise")
	fmt.Println("  • Estimates true state and derivative")
	fmt.Println("  • More complex but best performance")
	fmt.Println("  • Requires tuning process/measurement variance")
}
