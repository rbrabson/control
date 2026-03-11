// Package main demonstrates PID control for servo position tracking.
//
// This example shows position control with a PID controller,
// demonstrating precise positioning behavior.
package main

import (
	"fmt"

	"control/pid"
)

func main() {
	fmt.Println("Position Servo Control Example")
	fmt.Println("==============================")
	fmt.Println()

	fmt.Println("PID Position Tracking")
	fmt.Println("--------------------")
	fmt.Println("Demonstrates servo position control with PID.")
	fmt.Println()

	// Create PID controller for position
	// Output is servo command (-1.0 to 1.0)
	controller := pid.New(0.15, 0.05, 0.08,
		pid.WithOutputLimits(-1.0, 1.0))

	fmt.Println("Controller Configuration:")
	fmt.Println("  Kp = 0.15 (proportional gain)")
	fmt.Println("  Ki = 0.05 (integral gain)")
	fmt.Println("  Kd = 0.08 (derivative gain)")
	fmt.Println("  Output limits: -1.0 to 1.0 (servo command)")
	fmt.Println()

	// Test: Move to target position
	fmt.Println("Test: Move to Position 90°")
	fmt.Printf("\n%-6s %-11s %-12s %-10s %-8s\n", "Step", "Target (°)", "Current (°)", "Error (°)", "Command")
	fmt.Printf("%-6s %-11s %-12s %-10s %-8s\n", "----", "----------", "-----------", "---------", "-------")

	targetPosition := 90.0
	dt := 0.1

	// Simulated servo positions (approaching target)
	currentPositions := []float64{
		0.0, 20.0, 40.0, 55.0, 68.0, 78.0, 84.0, 87.5, 89.0, 89.8, 90.0,
	}
	controller.Calculate(targetPosition, currentPositions[0])

	for i, pos := range currentPositions {
		command := controller.CalculateWithDt(targetPosition, pos, dt)
		error := targetPosition - pos

		fmt.Printf("%-6d %-11.0f %-12.1f %-10.1f %-8.3f\n",
			i+1, targetPosition, pos, error, command)
	}

	// Test: Reverse direction
	controller.Reset()

	fmt.Println("\nTest: Return to Position 30°")
	fmt.Printf("\n%-6s %-11s %-12s %-10s %-8s\n", "Step", "Target (°)", "Current (°)", "Error (°)", "Command")
	fmt.Printf("%-6s %-11s %-12s %-10s %-8s\n", "----", "----------", "-----------", "---------", "-------")

	targetPosition = 30.0

	// Servo moving back from 90 to 30
	positionsBack := []float64{
		90.0, 75.0, 60.0, 48.0, 38.0, 32.5, 30.5, 30.0,
	}
	controller.Calculate(targetPosition, positionsBack[0])

	for i, pos := range positionsBack {
		command := controller.CalculateWithDt(targetPosition, pos, dt)
		error := targetPosition - pos

		fmt.Printf("%-6d %-11.0f %-12.1f %-10.1f %-8.3f\n",
			i+1, targetPosition, pos, error, command)
	}

	// Test: Small adjustment
	controller.Reset()

	fmt.Println("\nTest: Small Adjustment (30° → 35°)")
	fmt.Printf("\n%-6s %-11s %-12s %-10s %-8s\n", "Step", "Target (°)", "Current (°)", "Error (°)", "Command")
	fmt.Printf("%-6s %-11s %-12s %-10s %-8s\n", "----", "----------", "-----------", "---------", "-------")

	targetPosition = 35.0

	// Fine positioning
	finePositions := []float64{
		30.0, 31.5, 33.0, 34.0, 34.7, 35.0,
	}
	controller.Calculate(targetPosition, finePositions[0])

	for i, pos := range finePositions {
		command := controller.CalculateWithDt(targetPosition, pos, dt)
		error := targetPosition - pos

		fmt.Printf("%-6d %-11.0f %-12.1f %-10.1f %-8.3f\n",
			i+1, targetPosition, pos, error, command)
	}

	fmt.Println("\nAnalysis:")
	fmt.Println("---------")
	fmt.Println()
	fmt.Println("Large movement (0° → 90°):")
	fmt.Println("  • Strong initial command (high error)")
	fmt.Println("  • Command decreases as position approaches target")
	fmt.Println("  • D term helps slow down near target")
	fmt.Println()
	fmt.Println("Reverse movement (90° → 30°):")
	fmt.Println("  • Negative command for opposite direction")
	fmt.Println("  • Same control behavior, different direction")
	fmt.Println()
	fmt.Println("Fine adjustment (30° → 35°):")
	fmt.Println("  • Smaller commands for small errors")
	fmt.Println("  • Precise positioning capability")
	fmt.Println("  • I term ensures reaching exact target")
	fmt.Println()
	fmt.Println("Key Points:")
	fmt.Println("  • P term provides position correction force")
	fmt.Println("  • I term eliminates positioning error")
	fmt.Println("  • D term prevents overshoot")
	fmt.Println("  • Output limits protect servo mechanics")
}
