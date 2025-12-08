package main

import (
	"fmt"
	"math"
	"time"

	"control/feedforward"
)

// CombinedControlExample demonstrates a complex system using both gravity
// and cosine compensation, such as a crane or construction equipment.
func main() {
	fmt.Println("=== Combined Feedforward Control Example ===")
	fmt.Println("Simulating crane boom control with both gravity and cosine compensation")
	fmt.Println()

	// Create feedforward controller with both compensations
	// This simulates a crane boom that rotates and lifts
	ff := feedforward.New(
		0.08,                              // kS: static gain (bearing friction)
		1.5,                               // kV: velocity gain (higher for heavy machinery)
		0.12,                              // kA: acceleration gain
		feedforward.WithGravityGain(15.7), // Heavy load gravity compensation
		feedforward.WithCosineGain(8.2),   // Boom angle compensation
	)

	// Also create a basic feedforward for comparison
	ffBasic := feedforward.New(0.08, 1.5, 0.12)

	// Simulate crane operations
	operations := []struct {
		name       string
		startAngle float64
		endAngle   float64
		duration   float64
	}{
		{"Raise boom from horizontal", 0.0, math.Pi / 3, 4.0},
		{"Lower to mid position", math.Pi / 3, math.Pi / 6, 2.5},
		{"Raise to near vertical", math.Pi / 6, math.Pi * 0.45, 3.5},
	}

	fmt.Printf("%-6s %-25s %-8s %-8s %-8s %-10s %-10s %-10s\n",
		"Time", "Operation", "Angle°", "Vel", "Accel", "FF Basic", "FF Comb.", "Diff")
	fmt.Println("------------------------------------------------------------------------------")

	timeStep := 0.25
	for _, op := range operations {
		fmt.Printf("\n%s...\n", op.name)

		for t := 0.0; t <= op.duration; t += timeStep {
			// Smooth S-curve motion profile
			normalizedTime := t / op.duration
			var s, v, a float64

			// S-curve with smooth acceleration
			if normalizedTime < 0.25 {
				// Smooth start
				eta := 4 * normalizedTime
				s = (eta * eta) / 8.0
				v = eta / (2.0 * op.duration)
				a = 2.0 / (op.duration * op.duration)
			} else if normalizedTime < 0.75 {
				// Constant velocity
				s = 0.125 + (normalizedTime-0.25)*0.5
				v = 2.0 / op.duration
				a = 0.0
			} else {
				// Smooth stop
				eta := 4 * (1 - normalizedTime)
				s = 1.0 - (eta*eta)/8.0
				v = eta / (2.0 * op.duration)
				a = -2.0 / (op.duration * op.duration)
			}

			angleDiff := op.endAngle - op.startAngle
			angle := op.startAngle + s*angleDiff
			angularVel := v * angleDiff
			angularAccel := a * angleDiff

			// Calculate both feedforward outputs
			ffBasicOutput := ffBasic.Calculate(angle, angularVel, angularAccel)
			ffCombinedOutput := ff.Calculate(angle, angularVel, angularAccel)
			difference := ffCombinedOutput - ffBasicOutput

			// Convert to degrees for display
			angleDeg := angle * 180.0 / math.Pi

			fmt.Printf("%-6.1f %-25s %-8.1f %-8.3f %-8.3f %-10.3f %-10.3f %-10.3f\n",
				t, "", angleDeg, angularVel, angularAccel,
				ffBasicOutput, ffCombinedOutput, difference)

			time.Sleep(120 * time.Millisecond)
		}
	}

	fmt.Println("\n=== Analysis ===")
	fmt.Println("Combined compensation provides:")
	fmt.Println("- Gravity compensation: Constant upward force for load support")
	fmt.Println("- Cosine compensation: Variable torque based on boom angle")
	fmt.Println("- At horizontal (0°): Maximum cosine effect + full gravity")
	fmt.Println("- At vertical (90°): No cosine effect + full gravity")
	fmt.Println("- Significantly higher control effort than basic feedforward")
	fmt.Printf("- Gravity component: %.1f units\n", 15.7)
	fmt.Printf("- Max cosine component: %.1f units (at 0°)\n", 8.2)
	fmt.Printf("- Min cosine component: %.1f units (at 90°)\n", 8.2*math.Cos(math.Pi/2))
}
