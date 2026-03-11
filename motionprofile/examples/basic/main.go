// Package main demonstrates motion profile goal reaching behavior.
//
// This example shows the basic operation of trapezoidal motion profiles,
// matching the behavior tested in MotionProfileTest.java.
package main

import (
	"fmt"
	"math"

	"control/motionprofile"
)

func main() {
	fmt.Println("Basic Motion Profile Example")
	fmt.Println("============================")
	fmt.Println()

	// Test: reachesGoalAtTotalTime
	// Similar to Java test: MotionProfileTest.reachesGoalAtTotalTime()
	fmt.Println("Test: Reaches Goal at Total Time")
	fmt.Println("--------------------------------")

	// Create motion constraints
	constraints := motionprofile.Constraints{
		MaxVelocity:     2.0, // units per second
		MaxAcceleration: 1.0, // units per second²
	}

	// Define initial and goal states (matching Java test)
	initial := motionprofile.State{
		Position: 0.0,
		Velocity: 0.0,
	}
	goal := motionprofile.State{
		Position: 5.0,
		Velocity: 0.0,
	}

	fmt.Println("Configuration:")
	fmt.Printf("  Max Velocity: %.1f units/s\n", constraints.MaxVelocity)
	fmt.Printf("  Max Acceleration: %.1f units/s²\n", constraints.MaxAcceleration)
	fmt.Printf("  Initial: position=%.1f, velocity=%.1f\n",
		initial.Position, initial.Velocity)
	fmt.Printf("  Goal: position=%.1f, velocity=%.1f\n",
		goal.Position, goal.Velocity)

	// Create the motion profile
	profile := motionprofile.New(constraints, initial, goal)

	totalTime := profile.TotalTime()
	fmt.Printf("\n  Total Time: %.3f seconds\n", totalTime)

	// Calculate state at total time
	endState := profile.Calculate(totalTime)

	fmt.Println("\nState at Total Time:")
	fmt.Printf("  Time: %.3f s\n", endState.Time)
	fmt.Printf("  Position: %.6f\n", endState.Position)
	fmt.Printf("  Velocity: %.6f\n", endState.Velocity)
	fmt.Printf("  Acceleration: %.6f\n", endState.Acceleration)

	// Verify goal reached
	positionError := math.Abs(endState.Position - 5.0)
	fmt.Printf("\n Verify Goal Reached:")
	fmt.Printf("  Expected position: %.1f\n", goal.Position)
	fmt.Printf("  Actual position: %.9f\n", endState.Position)
	fmt.Printf("  Error: %.9f\n", positionError)

	if positionError < 1e-9 {
		fmt.Println("  ✓ Position matches goal (error < 1e-9)")
	} else {
		fmt.Printf("  ✗ Position does not match goal!\n")
	}

	// Verify profile finished
	isFinished := profile.IsFinished(totalTime)
	fmt.Printf("\n  IsFinished at total time: %v\n", isFinished)
	if isFinished {
		fmt.Println("  ✓ Profile is finished at total time")
	} else {
		fmt.Println("  ✗ Profile is not finished!")
	}

	fmt.Println()

	// Additional demonstration: Profile phases
	fmt.Println("Profile Phases Demonstration")
	fmt.Println("---------------------------")
	fmt.Println("Trapezoidal profile has 3 phases:")
	fmt.Println("  1. Acceleration: velocity increases from 0 to max")
	fmt.Println("  2. Cruise: constant max velocity")
	fmt.Println("  3. Deceleration: velocity decreases from max to 0")

	fmt.Println("\nSample Points:")
	fmt.Printf("%-8s %-10s %-10s %-10s %-14s\n", "Time", "Position", "Velocity", "Accel", "Phase")
	fmt.Printf("%-8s %-10s %-10s %-10s %-14s\n", "----", "--------", "--------", "-----", "-----")

	// Sample at key points
	sampleTimes := []float64{0.0, 0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, totalTime}

	for _, t := range sampleTimes {
		if t > totalTime {
			continue
		}
		state := profile.Calculate(t)

		// Determine phase
		phase := ""
		switch {
		case math.Abs(state.Acceleration) > 0.01 && state.Acceleration > 0:
			phase = "Accelerating"
		case math.Abs(state.Acceleration) > 0.01:
			phase = "Decelerating"
		case math.Abs(state.Velocity) > 0.01:
			phase = "Cruising"
		default:
			phase = "At rest"
		}

		fmt.Printf("%-8.1f %-10.3f %-10.3f %-10.3f %-14s\n",
			t, state.Position, state.Velocity, state.Acceleration, phase)
	}

	fmt.Println()

	// Additional example: Different goal
	fmt.Println("Example 2: Short Move (no cruise phase)")
	fmt.Println("---------------------------------------")

	initial2 := motionprofile.State{Position: 0.0, Velocity: 0.0}
	goal2 := motionprofile.State{Position: 1.5, Velocity: 0.0}

	profile2 := motionprofile.New(constraints, initial2, goal2)
	totalTime2 := profile2.TotalTime()
	endState2 := profile2.Calculate(totalTime2)

	fmt.Printf("  Goal: %.1f units\n", goal2.Position)
	fmt.Printf("  Total time: %.3f seconds\n", totalTime2)
	fmt.Printf("  Final position: %.6f\n", endState2.Position)
	fmt.Printf("  Position error: %.9f\n", math.Abs(endState2.Position-goal2.Position))

	if math.Abs(endState2.Position-goal2.Position) < 1e-9 {
		fmt.Println("  ✓ Short move reaches goal accurately")
	}

	fmt.Println("\nKey Points:")
	fmt.Println("• Motion profiles ensure smooth acceleration/deceleration")
	fmt.Println("• Profile always reaches exact goal position at total time")
	fmt.Println("• IsFinished() returns true when profile is complete")
	fmt.Println("• Respects velocity and acceleration constraints")
	fmt.Println("• Suitable for trajectory planning and motion control")
}
