package main

import (
	"fmt"
	"time"

	"control/motionprofile"
)

func main() {
	fmt.Println("=== Triangle Motion Profile Example ===")
	fmt.Println("Demonstrating triangle profile (short distance, can't reach max velocity)")
	fmt.Println()

	// Define motion constraints with high max velocity
	constraints := motionprofile.Constraints{
		MaxVelocity:     10.0, // High max velocity
		MaxAcceleration: 2.0,  // units per second²
	}

	// Define initial and goal states with short distance
	initial := motionprofile.State{
		Position: 0.0,
		Velocity: 0.0,
	}
	goal := motionprofile.State{
		Position: 3.0, // Short distance - can't reach max velocity
		Velocity: 0.0,
	}

	// Create the motion profile
	profile := motionprofile.New(constraints, initial, goal)

	fmt.Printf("Profile Parameters:\n")
	fmt.Printf("Total Time: %.3f seconds\n", profile.TotalTime())
	fmt.Printf("Initial Position: %.1f, Goal Position: %.1f\n", initial.Position, goal.Position)
	fmt.Printf("Max Velocity: %.1f, Max Acceleration: %.1f\n", constraints.MaxVelocity, constraints.MaxAcceleration)
	fmt.Println("Note: This creates a triangle profile because the distance is too short to reach max velocity")
	fmt.Println()

	// Sample the profile at regular intervals
	fmt.Printf("Time     Position     Velocity     Acceleration   Phase\n")
	fmt.Printf("------   ----------   ----------   ------------   -----------\n")

	sampleTime := 0.1
	totalTime := profile.TotalTime()
	for t := 0.0; t <= totalTime; t += sampleTime {
		state := profile.Calculate(t)

		// Determine which phase we're in
		phase := ""
		if t < totalTime/2 {
			phase = "Accel"
		} else {
			phase = "Decel"
		}

		fmt.Printf("%-8.1f %-12.3f %-12.3f %-12.3f   %-11s\n",
			state.Time, state.Position, state.Velocity, state.Acceleration, phase)

		time.Sleep(150 * time.Millisecond)
	}

	// Test the final state
	finalState := profile.Calculate(profile.TotalTime())
	fmt.Println()
	fmt.Printf("Final State: Position=%.6f, Velocity=%.6f\n",
		finalState.Position, finalState.Velocity)
	fmt.Printf("Goal State:  Position=%.6f, Velocity=%.6f\n",
		goal.Position, goal.Velocity)

	// Compare with a trapezoidal profile
	fmt.Println("\n=== Comparison: Trapezoidal Profile ===")

	// Same constraints but longer distance
	goalFar := motionprofile.State{
		Position: 20.0, // Longer distance - can reach max velocity
		Velocity: 0.0,
	}

	profileTrap := motionprofile.New(constraints, initial, goalFar)

	fmt.Printf("Longer distance profile:\n")
	fmt.Printf("Distance: %.1f, Total Time: %.3f seconds\n", goalFar.Position, profileTrap.TotalTime())
	fmt.Println("This creates a trapezoidal profile with acceleration, cruise, and deceleration phases")

	// Show key time comparisons
	fmt.Println("\nKey time comparisons:")
	fmt.Printf("Triangle profile (distance=%.1f): %.3f seconds\n", goal.Position, profile.TotalTime())
	fmt.Printf("Trapezoidal profile (distance=%.1f): %.3f seconds\n", goalFar.Position, profileTrap.TotalTime())
}
