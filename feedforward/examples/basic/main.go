// Package main demonstrates feed-forward control with gravity and cosine terms.
//
// This example shows the basic operation of a feed-forward controller,
// matching the behavior tested in FeedForwardTest.java.
package main

import (
	"fmt"
	"math"

	"control/feedforward"
)

func main() {
	fmt.Println("Basic Feed-Forward Controller Example")
	fmt.Println("=====================================\n")

	// Test: includesGravityAndCosineTerms
	// Similar to Java test: FeedForwardTest.includesGravityAndCosineTerms()
	fmt.Println("Test: Gravity and Cosine Terms")
	fmt.Println("------------------------------")
	fmt.Println("Create feed-forward controller with gravity and cosine compensation")

	// Create controller with kS=0.0, kV=2.0, kA=3.0
	// and add gravity gain=5.0, cosine gain=2.0
	ff := feedforward.New(0.0, 2.0, 3.0,
		feedforward.WithGravityGain(5.0),
		feedforward.WithCosineGain(2.0),
	)

	fmt.Println("\nController Configuration:")
	fmt.Printf("  kS (static): %.1f\n", 0.0)
	fmt.Printf("  kV (velocity): %.1f\n", 2.0)
	fmt.Printf("  kA (acceleration): %.1f\n", 3.0)
	fmt.Printf("  kG (gravity): %.1f\n", 5.0)
	fmt.Printf("  kCos (cosine): %.1f\n", 2.0)

	// Calculate output with specific inputs
	position := math.Pi
	velocity := 1.0
	acceleration := 2.0

	output := ff.Calculate(position, velocity, acceleration)

	fmt.Println("\nCalculation:")
	fmt.Printf("  Position: π (%.6f radians)\n", position)
	fmt.Printf("  Velocity: %.1f\n", velocity)
	fmt.Printf("  Acceleration: %.1f\n", acceleration)
	fmt.Println("\nFormula:")
	fmt.Println("  output = kV*velocity + kA*acceleration + kG*1.0 + kCos*cos(position)")
	fmt.Printf("  output = %.1f*%.1f + %.1f*%.1f + %.1f*1.0 + %.1f*cos(π)\n",
		2.0, velocity, 3.0, acceleration, 5.0, 2.0)
	fmt.Printf("  output = %.1f + %.1f + %.1f + %.1f*(%.1f)\n",
		2.0, 6.0, 5.0, 2.0, math.Cos(math.Pi))
	fmt.Printf("  output = %.1f + %.1f + %.1f - %.1f\n", 2.0, 6.0, 5.0, 2.0)
	fmt.Printf("  output = %.1f\n", output)

	// Verify expected value
	expected := 11.0
	if math.Abs(output-expected) < 1e-9 {
		fmt.Printf("\n  ✓ Output matches expected value: %.1f\n", expected)
	} else {
		fmt.Printf("\n  ✗ Output mismatch! Expected %.1f, got %.1f\n", expected, output)
	}

	fmt.Println("\nAdditional Examples:")
	fmt.Println("-------------------")

	// Example 2: At position = 0 (cos(0) = 1.0)
	fmt.Println("\nExample 2: Position = 0 (horizontal, cos(0) = 1.0)")
	pos2 := 0.0
	out2 := ff.Calculate(pos2, velocity, acceleration)
	fmt.Printf("  Position: %.1f, Velocity: %.1f, Accel: %.1f\n", pos2, velocity, acceleration)
	fmt.Printf("  Output: %.1f = %.1f + %.1f + %.1f + %.1f\n",
		out2, 2.0, 6.0, 5.0, 2.0*math.Cos(pos2))

	// Example 3: At position = π/2 (cos(π/2) = 0)
	fmt.Println("\nExample 3: Position = π/2 (vertical, cos(π/2) = 0)")
	pos3 := math.Pi / 2
	out3 := ff.Calculate(pos3, velocity, acceleration)
	fmt.Printf("  Position: π/2 (%.6f), Velocity: %.1f, Accel: %.1f\n", pos3, velocity, acceleration)
	fmt.Printf("  Output: %.1f = %.1f + %.1f + %.1f + %.1f\n",
		out3, 2.0, 6.0, 5.0, 2.0*math.Cos(pos3))

	fmt.Println("\nKey Points:")
	fmt.Println("• Gravity term compensates for constant force (e.g., robot arm weight)")
	fmt.Println("• Cosine term compensates for position-dependent forces")
	fmt.Println("• Velocity term provides damping compensation")
	fmt.Println("• Acceleration term provides inertial compensation")
	fmt.Println("• Feed-forward improves tracking without feedback")
}
