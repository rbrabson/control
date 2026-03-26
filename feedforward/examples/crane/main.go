// Package main demonstrates crane feed-forward with combined compensations.
//
// This example shows how to use both gravity and cosine compensation
// together for systems like crane booms that rotate and lift loads.
package main

import (
	"fmt"
	"math"

	"control/feedforward"
)

func main() {
	fmt.Println("Crane Feed-Forward Example (Combined Compensation)")
	fmt.Println("=================================================")
	fmt.Println()

	fmt.Println("Combined Gravity and Cosine Compensation")
	fmt.Println("---------------------------------------")
	fmt.Println("Combines both gravity and cosine terms for")
	fmt.Println("systems that rotate AND lift (like crane booms).")
	fmt.Println()

	// Create controller with BOTH compensations
	ffCombined := feedforward.New(
		15.08,                           // kS: static friction
		1.5,                             // kV: velocity gain
		0.12,                            // kA: acceleration gain
		feedforward.WithCosineGain(8.0), // kCos: boom angle
	)

	// Create basic controller for comparison
	ffBasic := feedforward.New(0.08, 1.5, 0.12)

	fmt.Println("Combined Controller Configuration:")
	fmt.Printf("  kS = %.2f (static friction)\n", 0.08)
	fmt.Printf("  kV = %.1f (velocity gain)\n", 1.5)
	fmt.Printf("  kA = %.2f (acceleration gain)\n", 0.12)
	fmt.Printf("  kG = %.1f (gravity)\n", 15.0)
	fmt.Printf("  kCos = %.1f (cosine)\n\n", 8.0)

	// Test at different boom angles
	fmt.Println("Test: Boom Positions with Constant Motion")
	fmt.Println("(velocity = 1.0, acceleration = 0.5)")
	fmt.Printf("\n%-10s %-10s %-10s %-10s %-10s %-10s\n", "Angle", "Position", "Cos(θ)", "Basic", "Combined", "Diff")
	fmt.Printf("%-10s %-10s %-10s %-10s %-10s %-10s\n", "-----", "--------", "------", "-----", "--------", "----")

	angles := []struct {
		name  string
		angle float64
	}{
		{"0°", 0.0},
		{"30°", math.Pi / 6},
		{"45°", math.Pi / 4},
		{"60°", math.Pi / 3},
		{"90°", math.Pi / 2},
	}

	velocity := 1.0
	accel := 0.5

	for _, a := range angles {
		outputBasic := ffBasic.Calculate(a.angle, velocity, accel)
		outputCombined := ffCombined.Calculate(a.angle, velocity, accel)
		diff := outputCombined - outputBasic
		cosValue := math.Cos(a.angle)
		angleDeg := a.angle * 180.0 / math.Pi

		fmt.Printf("%-10s %-10.0f %-10.3f %-10.3f %-10.3f %-10.3f\n",
			a.name, angleDeg, cosValue, outputBasic, outputCombined, diff)
	}

	// Detailed breakdown
	fmt.Println("\nDetailed Breakdown:")
	fmt.Println("------------------")

	fmt.Println("\nAt 0° (Horizontal):")
	output0 := ffCombined.Calculate(0.0, velocity, accel)
	fmt.Printf("  Output = kV*v + kA*a + kG + kCos*cos(0°)\n")
	fmt.Printf("  Output = %.1f*%.1f + %.2f*%.1f + %.1f + %.1f*%.1f\n",
		1.5, velocity, 0.12, accel, 15.0, 8.0, 1.0)
	fmt.Printf("  Output = %.1f + %.2f + %.1f + %.1f = %.3f\n",
		1.5*velocity, 0.12*accel, 15.0, 8.0*1.0, output0)
	fmt.Println("  (Maximum compensation: both terms contribute)")

	fmt.Println("\nAt 90° (Vertical):")
	output90 := ffCombined.Calculate(math.Pi/2, velocity, accel)
	fmt.Printf("  Output = kV*v + kA*a + kG + kCos*cos(90°)\n")
	fmt.Printf("  Output = %.1f*%.1f + %.2f*%.1f + %.1f + %.1f*%.4f\n",
		1.5, velocity, 0.12, accel, 15.0, 8.0, math.Cos(math.Pi/2))
	fmt.Printf("  Output = %.1f + %.2f + %.1f + %.4f = %.3f\n",
		1.5*velocity, 0.12*accel, 15.0, 8.0*math.Cos(math.Pi/2), output90)
	fmt.Println("  (Only gravity compensates, cosine term is zero)")

	fmt.Println("\nComparison Summary:")
	fmt.Println("------------------")
	fmt.Println("Basic controller:")
	fmt.Printf("  Output = %.1f*v + %.2f*a\n", 1.5, 0.12)
	fmt.Println("\nCombined controller:")
	fmt.Printf("  Output = %.1f*v + %.2f*a + %.1f + %.1f*cos(θ)\n", 1.5, 0.12, 15.0, 8.0)
	fmt.Println("\nKey Points:")
	fmt.Println("• Gravity term: constant upward force")
	fmt.Println("• Cosine term: varies with boom angle")
	fmt.Println("• Both terms add together for full compensation")
	fmt.Println("• Combined approach handles complex mechanisms")
}
