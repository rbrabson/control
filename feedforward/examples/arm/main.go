// Package main demonstrates robotic arm feed-forward with cosine compensation.
//
// This example shows how cosine compensation handles varying gravitational
// torque as the arm rotates through different angles.
package main

import (
	"fmt"
	"math"

	"control/feedforward"
)

func main() {
	fmt.Println("Robotic Arm Feed-Forward Example")
	fmt.Println("================================")
	fmt.Println()

	fmt.Println("Cosine Compensation for Rotating Arms")
	fmt.Println("-------------------------------------")
	fmt.Println("Cosine term compensates for varying gravitational")
	fmt.Println("torque as arm angle changes.")
	fmt.Println()

	// Create controller with cosine compensation
	ff := feedforward.New(
		0.02,                            // kS: static friction
		0.5,                             // kV: velocity gain
		0.03,                            // kA: acceleration gain
		feedforward.WithCosineGain(2.5), // kCos: gravity compensation
	)

	fmt.Println("Controller Configuration:")
	fmt.Printf("  kS = %.2f (static friction)\n", 0.02)
	fmt.Printf("  kV = %.1f (velocity gain)\n", 0.5)
	fmt.Printf("  kA = %.2f (acceleration gain)\n", 0.03)
	fmt.Printf("  kCos = %.1f (cosine compensation)\n\n", 2.5)

	// Test at key arm positions with constant motion
	fmt.Println("Test: Arm Positions with Constant Motion")
	fmt.Println("(velocity = 1.0, acceleration = 0.5)")
	fmt.Printf("\n%-18s %-8s %-10s %-10s\n", "Position", "Angle", "Cos(θ)", "Output")
	fmt.Printf("%-18s %-8s %-10s %-10s\n", "--------", "-----", "------", "------")

	positions := []struct {
		name  string
		angle float64
	}{
		{"Horizontal", 0.0},
		{"45° Up", math.Pi / 4},
		{"Vertical", math.Pi / 2},
		{"135°", 3 * math.Pi / 4},
		{"Horizontal Left", math.Pi},
	}

	velocity := 1.0
	accel := 0.5

	for _, pos := range positions {
		output := ff.Calculate(pos.angle, velocity, accel)
		cosValue := math.Cos(pos.angle)
		angleDeg := pos.angle * 180.0 / math.Pi

		fmt.Printf("%-18s %-8.0f %-10.3f %-10.3f\n",
			pos.name, angleDeg, cosValue, output)
	}

	// Detailed breakdown
	fmt.Println("\nDetailed Breakdown:")
	fmt.Println("------------------")

	fmt.Println("\nAt Horizontal (0°):")
	output0 := ff.Calculate(0.0, velocity, accel)
	fmt.Printf("  Output = kV*v + kA*a + kCos*cos(0°)\n")
	fmt.Printf("  Output = %.2f*%.1f + %.2f*%.1f + %.1f*%.1f\n",
		0.5, velocity, 0.03, accel, 2.5, 1.0)
	fmt.Printf("  Output = %.2f + %.2f + %.1f = %.3f\n",
		0.5*velocity, 0.03*accel, 2.5*1.0, output0)

	fmt.Println("\nAt Vertical (90°):")
	outputPi2 := ff.Calculate(math.Pi/2, velocity, accel)
	fmt.Printf("  Output = kV*v + kA*a + kCos*cos(90°)\n")
	fmt.Printf("  Output = %.2f*%.1f + %.2f*%.1f + %.1f*%.4f\n",
		0.5, velocity, 0.03, accel, 2.5, math.Cos(math.Pi/2))
	fmt.Printf("  Output = %.2f + %.2f + %.4f = %.3f\n",
		0.5*velocity, 0.03*accel, 2.5*math.Cos(math.Pi/2), outputPi2)

	fmt.Println("\nKey Points:")
	fmt.Println("• Cosine compensation varies with arm angle")
	fmt.Println("• Maximum torque at horizontal (cos(0°) = 1.0)")
	fmt.Println("• Zero torque at vertical (cos(90°) = 0.0)")
	fmt.Println("• Negative torque opposite side (cos(180°) = -1.0)")
	fmt.Println("• Compensates for changing gravitational load")
}
