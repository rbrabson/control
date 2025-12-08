package main

import (
	"fmt"
	"time"

	"control/feedforward"
)

// ElevatorControlExample demonstrates feedforward control for an elevator system
// with gravity compensation using the WithGravityGain option.
func main() {
	fmt.Println("=== Elevator Feedforward Control Example ===")
	fmt.Println("Simulating elevator control with gravity compensation")
	fmt.Println()

	// Create feedforward controller with gravity compensation
	// The gravity gain counteracts the weight of the elevator car
	ff := feedforward.New(
		0.05,                              // kS: static gain
		1.2,                               // kV: velocity gain (higher for elevator dynamics)
		0.08,                              // kA: acceleration gain
		feedforward.WithGravityGain(9.81), // Gravity compensation
	)

	// Simulate elevator moving between floors
	floors := []struct {
		name   string
		height float64 // meters
	}{
		{"Ground", 0.0},
		{"Floor 2", 3.5},
		{"Floor 5", 14.0},
		{"Floor 8", 28.0},
	}

	fmt.Printf("%-8s %-12s %-12s %-12s %-12s %-12s\n",
		"Time", "Floor", "Position", "Velocity", "Accel", "FF Output")
	fmt.Println("------------------------------------------------------------------------")

	timeStep := 0.2
	for i := 0; i < len(floors)-1; i++ {
		startHeight := floors[i].height
		endHeight := floors[i+1].height
		distance := endHeight - startHeight
		moveTime := 4.0 // seconds to move between floors

		fmt.Printf("Moving from %s to %s...\n", floors[i].name, floors[i+1].name)

		for t := 0.0; t <= moveTime; t += timeStep {
			// S-curve motion profile for smooth elevator movement
			normalizedTime := t / moveTime
			var s, v, a float64

			if normalizedTime < 0.5 {
				// Acceleration phase
				s = 2 * normalizedTime * normalizedTime
				v = 4 * normalizedTime / moveTime
				a = 4 / (moveTime * moveTime)
			} else {
				// Deceleration phase
				s = 1 - 2*(1-normalizedTime)*(1-normalizedTime)
				v = 4 * (1 - normalizedTime) / moveTime
				a = -4 / (moveTime * moveTime)
			}

			position := startHeight + s*distance
			velocity := v * distance
			acceleration := a * distance

			// Calculate feedforward output (includes gravity compensation)
			ffOutput := ff.Calculate(position, velocity, acceleration)

			fmt.Printf("%-8.1f %-12s %-12.2f %-12.3f %-12.3f %-12.3f\n",
				t, fmt.Sprintf("%.1fm", position), position, velocity, acceleration, ffOutput)

			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println()
	}

	fmt.Println("=== Analysis ===")
	fmt.Println("Gravity compensation provides:")
	fmt.Println("- Constant upward force to counteract elevator weight")
	fmt.Println("- Reduced motor effort during upward movement")
	fmt.Println("- More consistent performance regardless of direction")
	fmt.Printf("- Base gravity compensation: %.2f N\n", 9.81)
}
