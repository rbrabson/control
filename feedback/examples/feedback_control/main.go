package main

import (
	"fmt"
	"math"

	"control/feedback"
)

func main() {
	fmt.Println("=== Feedback Control Examples ===")
	fmt.Println()

	// Example 1: NoFeedback Controller
	fmt.Println("1. NoFeedback Controller:")
	fmt.Println("   Always returns zero, useful for open-loop control or as a placeholder")

	noFeedback := &feedback.NoFeedback{}
	output1 := noFeedback.Calculate(10.0, 8.5)
	output2 := noFeedback.Calculate(-5.0, 3.2)

	fmt.Printf("   NoFeedback.Calculate(10.0, 8.5) = %.2f\n", output1)
	fmt.Printf("   NoFeedback.Calculate(-5.0, 3.2) = %.2f\n", output2)
	fmt.Println("   (Always returns 0.0 regardless of input)")
	fmt.Println()

	// Example 2: Single State Feedback Controller
	fmt.Println("2. Single State Feedback (Position Control):")
	fmt.Println("   Simple proportional feedback for position control")

	positionGain := feedback.Values{2.5} // Proportional gain for position
	positionController := feedback.New(positionGain)

	targetPosition := feedback.Values{10.0}
	currentPosition := feedback.Values{8.2}

	output, err := positionController.Calculate(targetPosition, currentPosition)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Target: %.1f, Current: %.1f\n", targetPosition[0], currentPosition[0])
		fmt.Printf("   Control Output: %.2f\n", output)
		fmt.Printf("   (Gain × Error = 2.5 × (10.0 - 8.2) = %.2f)\n", 2.5*(10.0-8.2))
	}
	fmt.Println()

	// Example 3: Two State Feedback Controller (Position + Velocity)
	fmt.Println("3. Two State Feedback (Position + Velocity Control):")
	fmt.Println("   Feedback on both position and velocity for better stability")

	// Gains for [position, velocity]
	pvGain := feedback.Values{1.8, 0.4}
	pvController := feedback.New(pvGain)

	// Target: position=5.0, velocity=0.0 (stopped at target)
	target := feedback.Values{5.0, 0.0}
	// Current: position=3.5, velocity=2.1 (moving toward target)
	current := feedback.Values{3.5, 2.1}

	output, err = pvController.Calculate(target, current)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Target: [pos=%.1f, vel=%.1f]\n", target[0], target[1])
		fmt.Printf("   Current: [pos=%.1f, vel=%.1f]\n", current[0], current[1])
		fmt.Printf("   Control Output: %.2f\n", output)
		posError := target[0] - current[0]
		velError := target[1] - current[1]
		fmt.Printf("   (1.8×%.1f + 0.4×%.1f = %.2f)\n", posError, velError, 1.8*posError+0.4*velError)
	}
	fmt.Println()

	// Example 4: Three State Feedback Controller (Position + Velocity + Acceleration)
	fmt.Println("4. Three State Feedback (PVA Control):")
	fmt.Println("   Full state feedback for position, velocity, and acceleration")

	// Gains for [position, velocity, acceleration]
	pvaGain := feedback.Values{1.2, 0.8, 0.1}
	pvaController := feedback.New(pvaGain)

	// Target state: [position=0.0, velocity=0.0, acceleration=0.0]
	targetState := feedback.Values{0.0, 0.0, 0.0}
	// Current state: [position=-2.5, velocity=1.5, acceleration=-0.3]
	currentState := feedback.Values{-2.5, 1.5, -0.3}

	output, err = pvaController.Calculate(targetState, currentState)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Target: [pos=%.1f, vel=%.1f, acc=%.1f]\n",
			targetState[0], targetState[1], targetState[2])
		fmt.Printf("   Current: [pos=%.1f, vel=%.1f, acc=%.1f]\n",
			currentState[0], currentState[1], currentState[2])
		fmt.Printf("   Control Output: %.2f\n", output)

		posErr := targetState[0] - currentState[0]
		velErr := targetState[1] - currentState[1]
		accErr := targetState[2] - currentState[2]
		fmt.Printf("   (1.2×%.1f + 0.8×%.1f + 0.1×%.1f = %.2f)\n",
			posErr, velErr, accErr, 1.2*posErr+0.8*velErr+0.1*accErr)
	}
	fmt.Println()

	// Example 5: NoFeedback vs FullStateFeedback Comparison
	fmt.Println("5. Controller Comparison:")
	fmt.Println("   Comparing NoFeedback and FullStateFeedback outputs")

	noFeedbackCtrl := &feedback.NoFeedback{}
	fullStateCtrl := feedback.New(feedback.Values{1.5})

	setpoint := 10.0
	measurement := 7.0

	noOutput := noFeedbackCtrl.Calculate(setpoint, measurement)
	fullOutput, _ := fullStateCtrl.Calculate(feedback.Values{setpoint}, feedback.Values{measurement})

	fmt.Printf("   Setpoint: %.1f, Measurement: %.1f\n", setpoint, measurement)
	fmt.Printf("   NoFeedback output: %.2f\n", noOutput)
	fmt.Printf("   FullStateFeedback output: %.2f\n", fullOutput)
	fmt.Println()

	// Example 6: Error Handling
	fmt.Println("6. Error Handling:")
	fmt.Println("   Demonstrating error conditions with mismatched vector lengths")

	controller := feedback.New(feedback.Values{1.0, 0.5})

	// Correct usage
	validSetpoint := feedback.Values{5.0, 2.0}
	validMeasurement := feedback.Values{4.0, 2.5}
	output, err = controller.Calculate(validSetpoint, validMeasurement)
	fmt.Printf("   Valid input: Output=%.2f, Error=%v\n", output, err)

	// Error case: mismatched lengths
	invalidSetpoint := feedback.Values{5.0} // Wrong length
	validMeasurement = feedback.Values{4.0, 2.5}
	output, err = controller.Calculate(invalidSetpoint, validMeasurement)
	fmt.Printf("   Invalid input: Output=%.2f, Error=%v\n", output, err)
	fmt.Println()

	// Example 7: Motion Profile Integration
	fmt.Println("7. Motion Profile Integration Example:")
	fmt.Println("   Simulating a simple trajectory following scenario")

	// Controller gains for [position, velocity]
	gains := feedback.Values{3.0, 0.6}
	motionController := feedback.New(gains)

	// Simulate a simple trajectory (moving from 0 to 10 over time)
	fmt.Println("   Time | Target Pos | Current Pos | Target Vel | Current Vel | Control Output")
	fmt.Println("   -----|------------|-------------|------------|-------------|---------------")

	for t := 0; t <= 5; t++ {
		// Simple linear trajectory: position increases linearly, constant velocity
		targetPos := float64(t) * 2.0 // Move 2 units per time step
		targetVel := 2.0              // Constant desired velocity
		if targetPos > 10.0 {
			targetPos = 10.0
			targetVel = 0.0 // Stop at target
		}

		// Simulate some tracking error
		currentPos := targetPos - 0.3*math.Sin(float64(t)) // Some oscillation
		currentVel := targetVel - 0.2*math.Cos(float64(t)) // Velocity variation

		target := feedback.Values{targetPos, targetVel}
		current := feedback.Values{currentPos, currentVel}

		output, err := motionController.Calculate(target, current)
		if err != nil {
			fmt.Printf("   %d    | Error: %v\n", t, err)
		} else {
			fmt.Printf("   %d    |   %6.2f   |    %6.2f    |   %6.2f   |    %6.2f    |     %6.2f\n",
				t, targetPos, currentPos, targetVel, currentVel, output)
		}
	}

	fmt.Println("\n=== End of Examples ===")
}
