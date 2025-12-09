package main

import (
	"fmt"
	"math"
	"time"

	"control/feedforward"
)

// ControllerComparisonExample demonstrates the differences between
// all feedforward controller types available in the package.
func main() {
	fmt.Println("=== Feedforward Controller Comparison Example ===")
	fmt.Println("Comparing all available feedforward controller types")
	fmt.Println()

	// Create different controller types
	ffBasic := feedforward.New(0.1, 0.8, 0.05)
	ffGravity := feedforward.New(0.1, 0.8, 0.05, feedforward.WithGravityGain(5.0))
	ffCosine := feedforward.New(0.1, 0.8, 0.05, feedforward.WithCosineGain(3.0))
	ffCombined := feedforward.New(0.1, 0.8, 0.05,
		feedforward.WithGravityGain(5.0), feedforward.WithCosineGain(3.0))

	// Test motion profile
	timeStep := 0.4
	duration := 4.0

	fmt.Printf("%-6s %-8s %-8s %-8s %-8s %-8s\n",
		"Time", "Angle°", "Basic", "Gravity", "Cosine", "Combined")
	fmt.Println("--------------------------------------------------")

	for t := 0.0; t <= duration; t += timeStep {
		// Create a motion that showcases the differences
		// Position varies from 0 to π (0° to 180°)
		position := (math.Pi / 2) * (1 + math.Sin(2*math.Pi*t/duration))
		velocity := (math.Pi * math.Pi / duration) * math.Cos(2*math.Pi*t/duration)
		acceleration := (-2 * math.Pi * math.Pi * math.Pi / (duration * duration)) *
			math.Sin(2*math.Pi*t/duration)

		// Calculate outputs from all controllers
		outBasic := ffBasic.Calculate(position, velocity, acceleration)
		outGravity := ffGravity.Calculate(position, velocity, acceleration)
		outCosine := ffCosine.Calculate(position, velocity, acceleration)
		outCombined := ffCombined.Calculate(position, velocity, acceleration)

		// Convert angle to degrees for display
		angleDeg := position * 180.0 / math.Pi

		fmt.Printf("%-6.1f %-8.0f %-8.3f %-8.3f %-8.3f %-8.3f\n",
			t, angleDeg, outBasic, outGravity, outCosine, outCombined)

		time.Sleep(200 * time.Millisecond)
	}

	// Analyze specific positions
	fmt.Println("\n=== Detailed Analysis at Key Positions ===")
	positions := []struct {
		name  string
		angle float64
	}{
		{"Horizontal (0°)", 0.0},
		{"45° Up", math.Pi / 4},
		{"Vertical (90°)", math.Pi / 2},
		{"135°", 3 * math.Pi / 4},
		{"Horizontal (180°)", math.Pi},
	}

	// Test with constant motion parameters
	testVel := 1.0
	testAccel := 0.5

	fmt.Printf("\n%-15s %-8s %-8s %-8s %-8s\n",
		"Position", "Basic", "Gravity", "Cosine", "Combined")
	fmt.Println("-------------------------------------------")

	for _, pos := range positions {
		outBasic := ffBasic.Calculate(pos.angle, testVel, testAccel)
		outGravity := ffGravity.Calculate(pos.angle, testVel, testAccel)
		outCosine := ffCosine.Calculate(pos.angle, testVel, testAccel)
		outCombined := ffCombined.Calculate(pos.angle, testVel, testAccel)

		fmt.Printf("%-15s %-8.3f %-8.3f %-8.3f %-8.3f\n",
			pos.name, outBasic, outGravity, outCosine, outCombined)
	}

	fmt.Println("\n=== Controller Characteristics ===")
	fmt.Println("Basic FF:     kV*v + kA*a")
	fmt.Println("Gravity FF:   kV*v + kA*a + kG")
	fmt.Println("Cosine FF:    kV*v + kA*a + kCos*cos(θ)")
	fmt.Println("Combined FF:  kV*v + kA*a + kG + kCos*cos(θ)")

	fmt.Println("\n=== Use Cases ===")
	fmt.Println("Basic FF:     Simple systems without gravitational effects")
	fmt.Println("Gravity FF:   Elevators, vertical lifts, load handling")
	fmt.Println("Cosine FF:    Robotic arms, rotating machinery, pendulums")
	fmt.Println("Combined FF:  Cranes, construction equipment, complex robotics")
}
