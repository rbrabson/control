// Package main demonstrates full-state feedback control dot product calculation.
//
// This example shows the basic operation of full-state feedback,
// matching the behavior tested in FullStateFeedbackTest.java.
package main

import (
	"fmt"
	"math"

	"control/feedback"
)

func main() {
	fmt.Println("Full-State Feedback Control Example")
	fmt.Println("===================================\n")

	// Test 1: calculatesDotProductOfErrorAndGain
	// Similar to Java test: FullStateFeedbackTest.calculatesDotProductOfErrorAndGain()
	fmt.Println("Test 1: Dot Product Calculation")
	fmt.Println("-------------------------------")
	fmt.Println("Create feedback controller with gains [1.5, 0.3]")

	// Create controller with 2-state gains
	fsf := feedback.New(feedback.Values{1.5, 0.3})

	// Calculate control output
	reference := feedback.Values{10.0, 0.0}
	state := feedback.Values{8.0, 1.0}

	output, err := fsf.Calculate(reference, state)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("\nConfiguration:")
	fmt.Printf("  Gains: [%.1f, %.1f]\n", 1.5, 0.3)
	fmt.Printf("  Reference: [%.1f, %.1f]\n", reference[0], reference[1])
	fmt.Printf("  State: [%.1f, %.1f]\n", state[0], state[1])

	// Calculate error
	error0 := reference[0] - state[0]
	error1 := reference[1] - state[1]

	fmt.Println("\nCalculation:")
	fmt.Printf("  Error = Reference - State\n")
	fmt.Printf("  Error = [%.1f - %.1f, %.1f - %.1f]\n",
		reference[0], state[0], reference[1], state[1])
	fmt.Printf("  Error = [%.1f, %.1f]\n", error0, error1)

	fmt.Println("\nDot Product:")
	fmt.Printf("  Output = Gain · Error\n")
	fmt.Printf("  Output = [%.1f, %.1f] · [%.1f, %.1f]\n", 1.5, 0.3, error0, error1)
	fmt.Printf("  Output = %.1f × %.1f + %.1f × %.1f\n", 1.5, error0, 0.3, error1)
	fmt.Printf("  Output = %.1f + %.1f\n", 1.5*error0, 0.3*error1)
	fmt.Printf("  Output = %.1f\n", output)

	// Verify expected value
	expected := 2.7
	if math.Abs(output-expected) < 1e-9 {
		fmt.Printf("\n  ✓ Output matches expected value: %.1f\n", expected)
	} else {
		fmt.Printf("\n  ✗ Output mismatch! Expected %.1f, got %.1f\n", expected, output)
	}

	fmt.Println()

	// Test 2: throwsOnMismatchedVectorLengths
	// Similar to Java test: FullStateFeedbackTest.throwsOnMismatchedVectorLengths()
	fmt.Println("Test 2: Vector Length Validation")
	fmt.Println("--------------------------------")
	fmt.Println("Verify that mismatched vector lengths cause an error")

	fsf2 := feedback.New(feedback.Values{1.0, 1.0})

	// Try with mismatched lengths
	invalidRef := feedback.Values{1.0}      // Length 1
	validState := feedback.Values{1.0, 2.0} // Length 2

	_, err2 := fsf2.Calculate(invalidRef, validState)

	fmt.Printf("\n  Controller expects: 2 states\n")
	fmt.Printf("  Reference provided: %d state(s)\n", len(invalidRef))
	fmt.Printf("  State provided: %d state(s)\n", len(validState))

	if err2 != nil {
		fmt.Printf("\n  ✓ Error caught: %v\n", err2)
	} else {
		fmt.Println("\n  ✗ No error raised for mismatched lengths!")
	}

	fmt.Println()

	// Additional demonstration
	fmt.Println("Test 3: Multi-State Example")
	fmt.Println("---------------------------")
	fmt.Println("3-state system: [position, velocity, acceleration]")

	// Create 3-state controller
	fsf3 := feedback.New(feedback.Values{2.0, 1.5, 0.5})

	ref3 := feedback.Values{10.0, 0.0, 0.0}   // Target: at position 10, stopped
	state3 := feedback.Values{7.0, 2.0, -0.5} // Current: moving toward target

	output3, err3 := fsf3.Calculate(ref3, state3)
	if err3 != nil {
		fmt.Printf("Error: %v\n", err3)
		return
	}

	fmt.Printf("\n  Gains: [%.1f, %.1f, %.1f]\n", 2.0, 1.5, 0.5)
	fmt.Printf("  Reference: [%.1f, %.1f, %.1f]\n", ref3[0], ref3[1], ref3[2])
	fmt.Printf("  State: [%.1f, %.1f, %.1f]\n", state3[0], state3[1], state3[2])

	err0 := ref3[0] - state3[0]
	err1 := ref3[1] - state3[1]
	errorAccel := ref3[2] - state3[2]

	fmt.Printf("  Error: [%.1f, %.1f, %.1f]\n", err0, err1, errorAccel)
	fmt.Printf("  Output: %.1f × %.1f + %.1f × %.1f + %.1f × %.1f = %.1f\n",
		2.0, err0, 1.5, err1, 0.5, errorAccel, output3)

	fmt.Println("\nKey Points:")
	fmt.Println("• Full-state feedback computes dot product of gain and error vectors")
	fmt.Println("• All vectors (gains, reference, state) must have same dimensionality")
	fmt.Println("• Suitable for multi-variable control systems")
	fmt.Println("• Error = Reference - State")
	fmt.Println("• Output = K · (Reference - State)")
}
