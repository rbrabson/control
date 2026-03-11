// Package main demonstrates basic PID controller usage.
//
// This example shows the core PID controller features,
// matching the behavior tested in PIDTest.java.
package main

import (
	"fmt"

	"control/filter"
	"control/pid"
)

func main() {
	fmt.Println("Basic PID Controller Example")
	fmt.Println("============================")
	fmt.Println()

	// Test 1: Output limits
	// Similar to Java test: respectsOutputLimits()
	fmt.Println("Test 1: Output Limits")
	fmt.Println("---------------------")
	fmt.Println("Create PID with Kp=10.0 and output limits [-5.0, 5.0]")

	pid1 := pid.New(10.0, 0.0, 0.0,
		pid.WithOutputLimits(-5.0, 5.0),
	)

	setpoint := 10.0
	state := 0.0
	pid1.Calculate(setpoint, state)
	output := pid1.CalculateWithDt(setpoint, state, 0.1)

	fmt.Printf("  Setpoint: %.1f\n", setpoint)
	fmt.Printf("  State: %.1f\n", state)
	fmt.Printf("  Error: %.1f\n", setpoint-state)
	fmt.Printf("  Raw P output (10.0 * 10.0): %.1f\n", 10.0*10.0)
	fmt.Printf("  Actual output (clamped): %.1f\n", output)

	if output <= 5.0 && output >= -5.0 {
		fmt.Println("  ✓ Output respects limits [-5.0, 5.0]")
	} else {
		fmt.Println("  ✗ Output exceeds limits!")
	}

	fmt.Println()

	// Test 2: Derivative filter
	// Similar to Java test: supportsDerivativeFilterOption()
	fmt.Println("Test 2: Derivative Filter Option")
	fmt.Println("--------------------------------")
	fmt.Println("Create PID with Kd=1.0 and low-pass filter (alpha=0.8)")

	lpFilter, err := filter.NewLowPassFilter(0.8)
	if err != nil {
		fmt.Printf("Error creating filter: %v\n", err)
		return
	}

	pid2 := pid.New(0.0, 0.0, 1.0,
		pid.WithFilter(lpFilter),
	)

	// Initialize PID
	pid2.Calculate(10.0, 0.0)

	// Apply step change
	output2 := pid2.CalculateWithDt(8.0, 0.0, 0.1)

	fmt.Printf("  First call: Calculate(10.0, 0.0) - initialize\n")
	fmt.Printf("  Second call: Calculate(8.0, 0.0) - error changed\n")
	fmt.Printf("  Output: %.6f\n", output2)

	// Check output is finite (not NaN or Inf)
	if !isFinite(output2) {
		fmt.Println("  ✗ Output is not finite!")
	} else {
		fmt.Println("  ✓ Derivative filter works, output is finite")
	}

	fmt.Println()

	// Additional demonstration: Combined features
	fmt.Println("Test 3: Combined Features")
	fmt.Println("-------------------------")
	fmt.Println("Create PID with all gains and both output limits and filter")

	filter3, _ := filter.NewLowPassFilter(0.5)
	pid3 := pid.New(2.0, 0.5, 0.1,
		pid.WithOutputLimits(-10.0, 10.0),
		pid.WithFilter(filter3),
	)

	kp, ki, kd := pid3.GetGains()
	outputMin, outputMax := pid3.GetOutputLimits()

	fmt.Printf("  Gains: Kp=%.1f, Ki=%.1f, Kd=%.1f\n", kp, ki, kd)
	fmt.Printf("  Output limits: [%.1f, %.1f]\n", outputMin, outputMax)
	fmt.Printf("  Filter alpha: %.1f\n", pid3.GetFilter().GetGain())

	// Run a few iterations
	fmt.Println("\n  Running control loop (setpoint=5.0):")
	fmt.Printf("  %-6s %-8s %-8s\n", "Step", "State", "Output")
	fmt.Printf("  %-6s %-8s %-8s\n", "----", "-----", "------")

	for i := 0; i < 5; i++ {
		state := float64(i) * 0.5
		output := pid3.CalculateWithDt(5.0, state, 0.1)
		fmt.Printf("  %-6d %-8.1f %-8.3f\n", i, state, output)
	}

	fmt.Println("\nKey Points:")
	fmt.Println("• Output limits clamp the controller output")
	fmt.Println("• Derivative filter smooths error rate changes")
	fmt.Println("• Multiple options can be combined")
	fmt.Println("• PID is suitable for real-time control loops")
}

// isFinite checks if a float64 is neither NaN nor Inf
func isFinite(x float64) bool {
	return !isNaN(x) && !isInf(x)
}

func isNaN(x float64) bool {
	return x != x
}

func isInf(x float64) bool {
	return x > 1e308 || x < -1e308
}
