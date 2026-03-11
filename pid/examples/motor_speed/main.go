// Package main demonstrates PID control for motor speed regulation.
//
// This example shows basic motor speed control with PID, demonstrating
// how the controller tracks a target speed.
package main

import (
	"fmt"

	"control/pid"
)

func main() {
	fmt.Println("Motor Speed Control Example")
	fmt.Println("==========================")
	fmt.Println()

	fmt.Println("PID Speed Tracking")
	fmt.Println("------------------")
	fmt.Println("Demonstrates motor speed control with PID.")
	fmt.Println()

	// Create PID controller for motor speed
	// Output is motor power (-1.0 to 1.0)
	controller := pid.New(0.05, 0.02, 0.01,
		pid.WithOutputLimits(-1.0, 1.0))

	fmt.Println("Controller Configuration:")
	fmt.Println("  Kp = 0.05 (proportional gain)")
	fmt.Println("  Ki = 0.02 (integral gain)")
	fmt.Println("  Kd = 0.01 (derivative gain)")
	fmt.Println("  Output limits: -1.0 to 1.0 (motor power)")
	fmt.Println()

	// Test: Ramp up to target speed
	fmt.Println("Test: Speed Ramp to 1000 RPM")
	fmt.Printf("\n%-6s %-13s %-14s %-10s %-8s\n", "Step", "Target (RPM)", "Current (RPM)", "Error", "Power")
	fmt.Printf("%-6s %-13s %-14s %-10s %-8s\n", "----", "------------", "-------------", "-----", "-----")

	targetSpeed := 1000.0
	dt := 0.1

	// Simulated motor speeds (gradually approaching target)
	currentSpeeds := []float64{
		0, 150, 350, 550, 700, 820, 900, 950, 980, 995, 1000,
	}
	controller.Calculate(targetSpeed, currentSpeeds[0])

	for i, speed := range currentSpeeds {
		power := controller.CalculateWithDt(targetSpeed, speed, dt)
		error := targetSpeed - speed

		fmt.Printf("%-6d %-13.0f %-14.0f %-10.0f %-8.3f\n",
			i+1, targetSpeed, speed, error, power)
	}

	// Test: Speed change
	controller.Reset()

	fmt.Println("\nTest: Speed Change (1000 → 500 RPM)")
	fmt.Printf("\n%-6s %-13s %-14s %-10s %-8s\n", "Step", "Target (RPM)", "Current (RPM)", "Error", "Power")
	fmt.Printf("%-6s %-13s %-14s %-10s %-8s\n", "----", "------------", "-------------", "-----", "-----")

	targetSpeed = 500.0

	// Motor decelerating from 1000 to 500
	speedsDown := []float64{
		1000, 900, 800, 700, 600, 550, 520, 505, 500,
	}
	controller.Calculate(targetSpeed, speedsDown[0])

	for i, speed := range speedsDown {
		power := controller.CalculateWithDt(targetSpeed, speed, dt)
		error := targetSpeed - speed

		fmt.Printf("%-6d %-13.0f %-14.0f %-10.0f %-8.3f\n",
			i+1, targetSpeed, speed, error, power)
	}

	fmt.Println("\nAnalysis:")
	fmt.Println("---------")
	fmt.Println()
	fmt.Println("Acceleration phase (0 → 1000 RPM):")
	fmt.Println("  • High initial power due to large error")
	fmt.Println("  • Power decreases as error reduces")
	fmt.Println("  • At target: output near zero (maintenance power)")
	fmt.Println()
	fmt.Println("Deceleration phase (1000 → 500 RPM):")
	fmt.Println("  • Negative power to slow motor down")
	fmt.Println("  • Controller reverses direction")
	fmt.Println("  • Integral term prevents steady-state error")
	fmt.Println()
	fmt.Println("Key Points:")
	fmt.Println("  • P term: Immediate response to error")
	fmt.Println("  • I term: Eliminates steady-state error")
	fmt.Println("  • D term: Reduces overshoot during changes")
	fmt.Println("  • Output limits prevent motor saturation")
}
