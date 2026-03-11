// Package main demonstrates comparing different feed-forward configurations.
//
// This example shows the output differences between basic, gravity,
// cosine, and combined feed-forward controllers.
package main

import (
	"fmt"
	"math"

	"control/feedforward"
)

func main() {
	fmt.Println("Feed-Forward Controller Comparison")
	fmt.Println("=================================")
	fmt.Println()

	fmt.Println("Comparing All Controller Types")
	fmt.Println("------------------------------")
	fmt.Println("Shows how different compensation terms affect output.")
	fmt.Println()

	// Create all four controller types
	ffBasic := feedforward.New(0.1, 0.8, 0.05)

	ffGravity := feedforward.New(0.1, 0.8, 0.05,
		feedforward.WithGravityGain(5.0))

	ffCosine := feedforward.New(0.1, 0.8, 0.05,
		feedforward.WithCosineGain(3.0))

	ffCombined := feedforward.New(0.1, 0.8, 0.05,
		feedforward.WithGravityGain(5.0),
		feedforward.WithCosineGain(3.0))

	fmt.Println("Controller Configurations:")
	fmt.Println("-------------------------")
	fmt.Println("All controllers share: kS=0.1, kV=0.8, kA=0.05")
	fmt.Println("  Basic:    (no additional terms)")
	fmt.Println("  Gravity:  + kG=5.0")
	fmt.Println("  Cosine:   + kCos=3.0")
	fmt.Println("  Combined: + kG=5.0, kCos=3.0")
	fmt.Println()

	// Test at specific positions with constant motion
	fmt.Println("Test: Constant Motion at Different Positions")
	fmt.Println("(velocity = 1.0, acceleration = 0.5)")
	fmt.Printf("\n%-8s %-10s %-8s %-8s %-8s %-10s\n", "Angle", "Position", "Basic", "Gravity", "Cosine", "Combined")
	fmt.Printf("%-8s %-10s %-8s %-8s %-8s %-10s\n", "-----", "--------", "-----", "-------", "------", "--------")

	positions := []struct {
		name  string
		angle float64
	}{
		{"0°", 0.0},
		{"45°", math.Pi / 4},
		{"90°", math.Pi / 2},
		{"135°", 3 * math.Pi / 4},
		{"180°", math.Pi},
	}

	velocity := 1.0
	accel := 0.5

	for _, pos := range positions {
		outBasic := ffBasic.Calculate(pos.angle, velocity, accel)
		outGravity := ffGravity.Calculate(pos.angle, velocity, accel)
		outCosine := ffCosine.Calculate(pos.angle, velocity, accel)
		outCombined := ffCombined.Calculate(pos.angle, velocity, accel)
		angleDeg := pos.angle * 180.0 / math.Pi

		fmt.Printf("%-8s %-10.0f %-8.3f %-8.3f %-8.3f %-10.3f\n",
			pos.name, angleDeg, outBasic, outGravity, outCosine, outCombined)
	}

	// Analysis section
	fmt.Println("\nDetailed Analysis:")
	fmt.Println("-----------------")

	fmt.Println("\nAt 0° (Horizontal):")
	fmt.Println("  Basic:    0.1*0 + 0.8*1.0 + 0.05*0.5           = 0.825")
	fmt.Println("  Gravity:  0.825 + 5.0                         = 5.825")
	fmt.Println("  Cosine:   0.825 + 3.0*cos(0°) = 0.825 + 3.0   = 3.825")
	fmt.Println("  Combined: 0.825 + 5.0 + 3.0                   = 8.825")

	fmt.Println("\nAt 90° (Vertical):")
	fmt.Println("  Basic:    0.825")
	fmt.Println("  Gravity:  0.825 + 5.0                         = 5.825")
	fmt.Println("  Cosine:   0.825 + 3.0*cos(90°) = 0.825 + 0    = 0.825")
	fmt.Println("  Combined: 0.825 + 5.0 + 0                     = 5.825")

	fmt.Println("\nComparison Summary:")
	fmt.Println("------------------")
	fmt.Println("Basic controller:")
	fmt.Println("  • Simplest form: only velocity and acceleration")
	fmt.Println("  • Formula: kS*sign(v) + kV*v + kA*a")
	fmt.Println()
	fmt.Println("Gravity controller:")
	fmt.Println("  • Adds constant upward force")
	fmt.Println("  • Formula: Basic + kG")
	fmt.Println("  • Use for: elevators, vertical lifts")
	fmt.Println()
	fmt.Println("Cosine controller:")
	fmt.Println("  • Adds angle-dependent force")
	fmt.Println("  • Formula: Basic + kCos*cos(position)")
	fmt.Println("  • Use for: rotating arms, pendulums")
	fmt.Println()
	fmt.Println("Combined controller:")
	fmt.Println("  • Uses both compensation terms")
	fmt.Println("  • Formula: Basic + kG + kCos*cos(position)")
	fmt.Println("  • Use for: crane booms, complex mechanisms")
}
