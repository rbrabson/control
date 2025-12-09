package main

import (
	"fmt"
	"time"

	"control/motionprofile"
)

func main() {
	fmt.Println("=== Basic Motion Profile Example ===")
	fmt.Println("Demonstrating trapezoidal motion profile generation")
	fmt.Println()

	// Define motion constraints
	constraints := motionprofile.Constraints{
		MaxVelocity:     2.0, // units per second
		MaxAcceleration: 1.0, // units per second²
	}

	// Define initial and goal states
	initial := motionprofile.State{
		Position: 0.0,
		Velocity: 0.0,
	}
	goal := motionprofile.State{
		Position: 10.0,
		Velocity: 0.0,
	}

	// Create the motion profile
	profile := motionprofile.New(constraints, initial, goal)

	fmt.Printf("Profile Parameters:\n")
	fmt.Printf("Total Time: %.3f seconds\n", profile.TotalTime())
	fmt.Printf("Initial Position: %.1f, Goal Position: %.1f\n", initial.Position, goal.Position)
	fmt.Printf("Max Velocity: %.1f, Max Acceleration: %.1f\n", constraints.MaxVelocity, constraints.MaxAcceleration)
	fmt.Println()

	// Sample the profile at regular intervals
	fmt.Printf("Time     Position     Velocity     Acceleration\n")
	fmt.Printf("------   ----------   ----------   ------------\n")

	sampleTime := 0.1
	for t := 0.0; t <= profile.TotalTime(); t += sampleTime {
		state := profile.Calculate(t)
		fmt.Printf("%-8.1f %-12.3f %-12.3f %-12.3f\n",
			state.Time, state.Position, state.Velocity, state.Acceleration)

		// Add a small delay to visualize the motion
		time.Sleep(100 * time.Millisecond)
	}

	// Test the final state
	finalState := profile.Calculate(profile.TotalTime())
	fmt.Println()
	fmt.Printf("Final State: Position=%.6f, Velocity=%.6f\n",
		finalState.Position, finalState.Velocity)
	fmt.Printf("Goal State:  Position=%.6f, Velocity=%.6f\n",
		goal.Position, goal.Velocity)

	// Test IsFinished method
	fmt.Println()
	fmt.Printf("Profile finished at t=%.1f? %v\n", profile.TotalTime()/2,
		profile.IsFinished(profile.TotalTime()/2))
	fmt.Printf("Profile finished at t=%.1f? %v\n", profile.TotalTime(),
		profile.IsFinished(profile.TotalTime()))
	fmt.Printf("Profile finished at t=%.1f? %v\n", profile.TotalTime()+1,
		profile.IsFinished(profile.TotalTime()+1))

	// Test TimeLeftUntil method
	fmt.Println()
	fmt.Printf("Time to reach position 5.0: %.3f seconds\n",
		profile.TimeLeftUntil(5.0))
	fmt.Printf("Time to reach position 7.5: %.3f seconds\n",
		profile.TimeLeftUntil(7.5))
	fmt.Printf("Time to reach goal position: %.3f seconds\n",
		profile.TimeLeftUntil(goal.Position))
}
