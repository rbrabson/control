// Package main demonstrates PID derivative filtering to reduce derivative kick.
//
// This example shows how derivative filtering smooths the derivative term
// to prevent sudden spikes when the setpoint changes.
package main

import (
	"fmt"

	"control/filter"
	"control/pid"
)

func main() {
	fmt.Println("PID Derivative Filtering (Dampening) Example")
	fmt.Println("===========================================")

	fmt.Println("Derivative Kick Reduction")
	fmt.Println("-------------------------")
	fmt.Println("Demonstrates how derivative filtering prevents sudden")
	fmt.Println("spikes in control output during setpoint changes.")
	fmt.Println()

	// Create two controllers: one without filter, one with filter
	pidNoFilter := pid.New(1.0, 0.1, 0.5)

	lowpassFilter, _ := filter.NewLowPassFilter(0.3)
	pidWithFilter := pid.New(1.0, 0.1, 0.5,
		pid.WithFilter(lowpassFilter))

	fmt.Println("Controller Configuration:")
	fmt.Println("  Both: Kp=1.0, Ki=0.1, Kd=0.5")
	fmt.Println("  No Filter: Standard PID")
	fmt.Println("  With Filter: Low-pass filter (alpha=0.3)")
	fmt.Println()

	// Test: Setpoint change from 0 to 100
	fmt.Println("Test: Setpoint Change (0 → 100)")
	fmt.Println("State starts at 0, setpoint suddenly changes to 100")
	fmt.Printf("\n%-6s %-10s %-8s %-11s %-11s %-11s\n", "Step", "Setpoint", "State", "No Filter", "With Filter", "Diff")
	fmt.Printf("%-6s %-10s %-8s %-11s %-11s %-11s\n", "----", "--------", "-----", "---------", "-----------", "----")

	setpoint := 100.0
	state := 0.0
	dt := 0.1 // 100ms time step

	// Initialize controllers with first measurement
	pidNoFilter.Calculate(setpoint, state)
	pidWithFilter.Calculate(setpoint, state)

	// First step: Large derivative kick
	outNoFilter := pidNoFilter.CalculateWithDt(setpoint, state, dt)
	outWithFilter := pidWithFilter.CalculateWithDt(setpoint, state, dt)
	fmt.Printf("%-6d %-10.0f %-8.0f %-11.3f %-11.3f %-11.3f\n",
		1, setpoint, state, outNoFilter, outWithFilter, outNoFilter-outWithFilter)

	// Simulate state approaching setpoint
	steps := []float64{20.0, 40.0, 60.0, 80.0, 90.0, 95.0, 98.0, 99.0, 100.0}
	for i, s := range steps {
		state = s
		outNoFilter = pidNoFilter.CalculateWithDt(setpoint, state, dt)
		outWithFilter = pidWithFilter.CalculateWithDt(setpoint, state, dt)
		fmt.Printf("%-6d %-10.0f %-8.0f %-11.3f %-11.3f %-11.3f\n",
			i+2, setpoint, state, outNoFilter, outWithFilter, outNoFilter-outWithFilter)
	}

	// Reset controllers and test setpoint decrease
	pidNoFilter.Reset()
	pidWithFilter.Reset()

	fmt.Println("\nTest: Setpoint Decrease (100 → 50)")
	fmt.Printf("\n%-6s %-10s %-8s %-11s %-11s %-11s\n", "Step", "Setpoint", "State", "No Filter", "With Filter", "Diff")
	fmt.Printf("%-6s %-10s %-8s %-11s %-11s %-11s\n", "----", "--------", "-----", "---------", "-----------", "----")

	setpoint = 50.0
	state = 100.0

	// Initialize
	pidNoFilter.Calculate(setpoint, state)
	pidWithFilter.Calculate(setpoint, state)

	outNoFilter = pidNoFilter.CalculateWithDt(setpoint, state, dt)
	outWithFilter = pidWithFilter.CalculateWithDt(setpoint, state, dt)
	fmt.Printf("%-6d %-10.0f %-8.0f %-11.3f %-11.3f %-11.3f\n",
		1, setpoint, state, outNoFilter, outWithFilter, outNoFilter-outWithFilter)

	stepsDown := []float64{90.0, 80.0, 70.0, 60.0, 55.0, 52.0, 50.5, 50.0}
	for i, s := range stepsDown {
		state = s
		outNoFilter = pidNoFilter.CalculateWithDt(setpoint, state, dt)
		outWithFilter = pidWithFilter.CalculateWithDt(setpoint, state, dt)
		fmt.Printf("%-6d %-10.0f %-8.0f %-11.3f %-11.3f %-11.3f\n",
			i+2, setpoint, state, outNoFilter, outWithFilter, outNoFilter-outWithFilter)
	}

	fmt.Println("\nAnalysis:")
	fmt.Println("---------")
	fmt.Println("Without filtering:")
	fmt.Println("  • Large derivative kick on step 1 (setpoint change)")
	fmt.Println("  • Derivative term reacts instantly to error changes")
	fmt.Println("  • Can cause actuator saturation or system instability")
	fmt.Println()
	fmt.Println("With filtering (alpha=0.3):")
	fmt.Println("  • First sample can match raw derivative (filter initialization)")
	fmt.Println("  • Derivative builds up gradually over multiple steps")
	fmt.Println("  • Smoother control output, better for real systems")
	fmt.Println("  • Trade-off: Slightly slower response")
}
