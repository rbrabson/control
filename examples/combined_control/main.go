package main

import (
	"fmt"

	"control/feedback"
	"control/pid"
)

func main() {
	fmt.Println("=== Combined PID + Feedback Control Example ===")
	fmt.Println()

	// Example: Robot Arm Control with PID and State Feedback
	fmt.Println("Robot Arm Joint Control System")
	fmt.Println("- PID controller for primary position control")
	fmt.Println("- State feedback for stabilization (position + velocity)")
	fmt.Println()

	// Create PID controller for position control
	pidController := pid.New(2.0, 0.1, 0.05,
		pid.WithOutputLimits(-10.0, 10.0))

	// Create state feedback controller for stabilization
	// Gains for [position_error, velocity]
	stabilizationGains := feedback.Values{0.8, 0.3}
	stateController := feedback.New(stabilizationGains)

	// Simulation parameters
	targetPosition := 45.0 // Target angle in degrees
	currentPosition := 0.0 // Starting position
	currentVelocity := 0.0 // Starting velocity
	dt := 0.1              // 100ms time step

	fmt.Println("Time(s) | Target | Position | Velocity | PID Output | State FB | Total Output")
	fmt.Println("--------|--------|----------|----------|------------|----------|-------------")

	// Simulate control system for 3 seconds
	for step := 0; step <= 30; step++ {
		t := float64(step) * dt

		// PID controller calculates primary control signal
		pidOutput := pidController.Calculate(targetPosition, currentPosition)

		// State feedback provides stabilization
		// Use position error and current velocity for state feedback
		positionError := targetPosition - currentPosition
		stateVector := feedback.Values{positionError, currentVelocity}
		referenceVector := feedback.Values{0.0, 0.0} // Want zero error and zero velocity at target

		stateFeedback, err := stateController.Calculate(referenceVector, stateVector)
		if err != nil {
			fmt.Printf("State feedback error: %v\n", err)
			continue
		}

		// Combine PID and state feedback outputs
		totalOutput := pidOutput + stateFeedback

		// Simple plant simulation (second-order system)
		// Acceleration proportional to control input
		acceleration := totalOutput * 0.5
		currentVelocity += acceleration * dt
		currentPosition += currentVelocity * dt

		// Add some damping to make it realistic
		currentVelocity *= 0.98

		// Print results every few steps for readability
		if step%3 == 0 {
			fmt.Printf("%6.1f  | %6.1f | %8.2f | %8.2f | %10.2f | %8.2f | %11.2f\n",
				t, targetPosition, currentPosition, currentVelocity, pidOutput, stateFeedback, totalOutput)
		}

		// Stop if we're close enough to target with low velocity
		if abs(targetPosition-currentPosition) < 0.5 && abs(currentVelocity) < 0.1 {
			fmt.Printf("\nTarget reached at t=%.1fs!\n", t)
			break
		}
	}

	fmt.Println()

	// Show controller configurations
	fmt.Println("Controller Configurations:")
	fmt.Printf("PID Gains: P=%.1f, I=%.1f, D=%.2f\n", 2.0, 0.1, 0.05)

	minLimit, maxLimit := pidController.GetOutputLimits()
	fmt.Printf("PID Output Limits: [%.1f, %.1f]\n", minLimit, maxLimit)

	fmt.Printf("State Feedback Gains: Position Error=%.1f, Velocity=%.1f\n",
		stabilizationGains[0], stabilizationGains[1])

	fmt.Println()

	// Example 2: Switching between different feedback strategies
	fmt.Println("=== Adaptive Feedback Example ===")
	fmt.Println("Switching between NoFeedback and FullState based on conditions")
	fmt.Println()

	// Create different feedback controllers
	noFeedback := &feedback.NoFeedback{}
	aggressiveFeedback := feedback.New(feedback.Values{2.0, 0.8})
	gentleFeedback := feedback.New(feedback.Values{1.0, 0.2})

	// Simulate different operating conditions
	scenarios := []struct {
		name       string
		posError   float64
		velocity   float64
		controller string
	}{
		{"Startup (large error)", 10.0, 0.1, "Aggressive"},
		{"Approaching target", 2.0, 1.5, "Gentle"},
		{"At target", 0.1, 0.05, "None"},
		{"Overshoot", -1.5, 2.0, "Aggressive"},
		{"Settling", 0.3, 0.2, "Gentle"},
	}

	fmt.Println("Scenario                | Pos Error | Velocity | Controller | Output")
	fmt.Println("------------------------|-----------|----------|------------|--------")

	for _, scenario := range scenarios {
		var output float64
		var err error

		state := feedback.Values{scenario.posError, scenario.velocity}
		target := feedback.Values{0.0, 0.0} // Always want zero error and velocity

		switch scenario.controller {
		case "None":
			// Use NoFeedback for very small errors (open loop)
			output = noFeedback.Calculate(0.0, scenario.posError) // Will always be 0
		case "Gentle":
			output, err = gentleFeedback.Calculate(target, state)
		case "Aggressive":
			output, err = aggressiveFeedback.Calculate(target, state)
		}

		if err != nil {
			fmt.Printf("%-23s | %9.1f | %8.2f | %-10s | Error\n",
				scenario.name, scenario.posError, scenario.velocity, scenario.controller)
		} else {
			fmt.Printf("%-23s | %9.1f | %8.2f | %-10s | %6.2f\n",
				scenario.name, scenario.posError, scenario.velocity, scenario.controller, output)
		}
	}

	fmt.Println("\n=== End of Combined Examples ===")
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
