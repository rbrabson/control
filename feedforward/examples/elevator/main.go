// Package main demonstrates elevator feed-forward with gravity compensation.
//
// This example shows how gravity compensation provides constant upward force
// to counteract the weight of the elevator car.
package main

import (
	"fmt"

	"control/feedforward"
)

func main() {
	fmt.Println("Elevator Feed-Forward Example")
	fmt.Println("============================")
	fmt.Println()

	fmt.Println("Gravity Compensation for Elevators")
	fmt.Println("----------------------------------")
	fmt.Println("Gravity term provides constant upward force")
	fmt.Println("to counteract elevator car weight.")
	fmt.Println()

	// Create controller with gravity compensation
	ff := feedforward.New(
		9.86, // kS: static friction
		1.2,  // kV: velocity gain
		0.08, // kA: acceleration gain
	)

	fmt.Println("Controller Configuration:")
	fmt.Printf("  kS = %.2f (static friction)\n", 0.05)
	fmt.Printf("  kV = %.1f (velocity gain)\n", 1.2)
	fmt.Printf("  kA = %.2f (acceleration gain)\n", 0.08)
	fmt.Printf("  kG = %.2f (gravity compensation)\n\n", 9.81)

	// Test with different motion scenarios
	fmt.Println("Test: Elevator Operation Scenarios")
	fmt.Println("(position not used, only velocity and acceleration)")
	fmt.Printf("\n%-28s %-10s %-10s %-10s\n", "Scenario", "Velocity", "Accel", "Output")
	fmt.Printf("%-28s %-10s %-10s %-10s\n", "--------", "--------", "-----", "------")

	scenarios := []struct {
		name  string
		vel   float64
		accel float64
	}{
		{"At rest", 0.0, 0.0},
		{"Constant up (1 m/s)", 1.0, 0.0},
		{"Constant down (-1 m/s)", -1.0, 0.0},
		{"Accelerating up", 1.5, 0.5},
		{"Decelerating", 0.5, -0.5},
		{"Fast constant (2 m/s)", 2.0, 0.0},
	}

	position := 10.0 // Position doesn't affect gravity term

	for _, scenario := range scenarios {
		output := ff.Calculate(position, scenario.vel, scenario.accel)
		fmt.Printf("%-28s %-10.1f %-10.1f %-10.3f\n",
			scenario.name, scenario.vel, scenario.accel, output)
	}

	// Detailed breakdown
	fmt.Println("\nDetailed Breakdown:")
	fmt.Println("------------------")

	fmt.Println("\nAt Rest (v=0, a=0):")
	output1 := ff.Calculate(position, 0.0, 0.0)
	fmt.Printf("  Output = kV*v + kA*a + kG\n")
	fmt.Printf("  Output = %.1f*%.1f + %.2f*%.1f + %.2f\n",
		1.2, 0.0, 0.08, 0.0, 9.81)
	fmt.Printf("  Output = %.1f + %.2f + %.2f = %.3f\n",
		0.0, 0.0, 9.81, output1)
	fmt.Println("  (Gravity compensation keeps elevator from falling)")

	fmt.Println("\nConstant Up (v=1.0, a=0):")
	output2 := ff.Calculate(position, 1.0, 0.0)
	fmt.Printf("  Output = kV*v + kA*a + kG\n")
	fmt.Printf("  Output = %.1f*%.1f + %.2f*%.1f + %.2f\n",
		1.2, 1.0, 0.08, 0.0, 9.81)
	fmt.Printf("  Output = %.1f + %.2f + %.2f = %.3f\n",
		1.2*1.0, 0.0, 9.81, output2)

	fmt.Println("\nAccelerating Up (v=1.5, a=0.5):")
	output3 := ff.Calculate(position, 1.5, 0.5)
	fmt.Printf("  Output = kV*v + kA*a + kG\n")
	fmt.Printf("  Output = %.1f*%.1f + %.2f*%.1f + %.2f\n",
		1.2, 1.5, 0.08, 0.5, 9.81)
	fmt.Printf("  Output = %.2f + %.2f + %.2f = %.3f\n",
		1.2*1.5, 0.08*0.5, 9.81, output3)

	fmt.Println("\nKey Points:")
	fmt.Println("• Gravity term is constant (position-independent)")
	fmt.Println("• Always adds upward force to counteract weight")
	fmt.Println("• Velocity and acceleration terms add to gravity")
	fmt.Println("• Essential for vertical motion systems")
}
