package main

import (
	"fmt"
	"math"
	"time"

	"control/feedforward"
)

// BasicFeedforwardExample demonstrates a simple feedforward controller
// for a motor system with velocity and acceleration feedforward terms.
func main() {
	fmt.Println("=== Basic Feedforward Control Example ===")
	fmt.Println("Simulating motor control with velocity and acceleration feedforward")
	fmt.Println()

	// Create a basic feedforward controller
	// kS = 0.1 (static gain), kV = 0.8 (velocity gain), kA = 0.05 (acceleration gain)
	ff := feedforward.New(0.1, 0.8, 0.05)

	// Simulate a motion profile with changing velocity and acceleration
	timeStep := 0.1 // seconds
	duration := 3.0 // seconds

	fmt.Printf("%-8s %-12s %-12s %-12s %-12s\n", "Time", "Position", "Velocity", "Accel", "FF Output")
	fmt.Println("--------------------------------------------------------------")

	for t := 0.0; t <= duration; t += timeStep {
		// Generate a smooth motion profile (sinusoidal)
		position := 2.0 * (1 - math.Cos(math.Pi*t/duration))
		velocity := (2.0 * math.Pi / duration) * math.Sin(math.Pi*t/duration)
		acceleration := (2.0 * math.Pi * math.Pi / (duration * duration)) * math.Cos(math.Pi*t/duration)

		// Calculate feedforward output
		ffOutput := ff.Calculate(position, velocity, acceleration)

		fmt.Printf("%-8.1f %-12.3f %-12.3f %-12.3f %-12.3f\n",
			t, position, velocity, acceleration, ffOutput)

		time.Sleep(50 * time.Millisecond) // Simulate real-time
	}

	fmt.Println("\n=== Analysis ===")
	fmt.Println("The feedforward controller provides:")
	fmt.Println("- Velocity compensation: Reduces steady-state error during motion")
	fmt.Println("- Acceleration compensation: Improves transient response")
	fmt.Println("- No position feedback: Pure feedforward control")
}
