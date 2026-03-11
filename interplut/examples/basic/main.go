// Package main demonstrates InterpLUT returning exact values at control points.
//
// This example shows that the interpolated lookup table returns exact values
// at the control points, matching the behavior tested in InterpLUTTest.java.
package main

import (
	"fmt"
	"math"

	"control/interplut"
)

func main() {
	fmt.Println("InterpLUT: Exact Control Point Returns")
	fmt.Println("======================================")
	fmt.Println()

	// Create a lookup table with two simple control points
	lut := interplut.New()
	lut.Add(0.0, 0.0)
	lut.Add(1.0, 1.0)

	err := lut.CreateLUT()
	if err != nil {
		fmt.Printf("Error creating LUT: %v\n", err)
		return
	}

	// Test that control points return exact values
	fmt.Println("Testing exact control point returns:")
	fmt.Printf("%-8s %-10s %-10s %-5s\n", "Input", "Expected", "Actual", "Match")
	fmt.Printf("%-8s %-10s %-10s %-5s\n", "-----", "--------", "------", "-----")

	tests := []struct {
		input    float64
		expected float64
	}{
		{0.0, 0.0},
		{1.0, 1.0},
	}

	for _, test := range tests {
		actual, err := lut.Get(test.input)
		if err != nil {
			fmt.Printf("%-8.1f Error: %v\n", test.input, err)
			continue
		}

		match := math.Abs(actual-test.expected) < 1e-9
		matchStr := "✓"
		if !match {
			matchStr = "✗"
		}

		fmt.Printf("%-8.1f %-10.6f %-10.6f %-5s\n", test.input, test.expected, actual, matchStr)
	}

	// Also test interpolated values
	fmt.Println("\nInterpolated values between control points:")
	fmt.Printf("%-8s %-10s\n", "Input", "Output")
	fmt.Printf("%-8s %-10s\n", "-----", "------")

	for x := 0.0; x <= 1.0; x += 0.25 {
		y, err := lut.Get(x)
		if err != nil {
			fmt.Printf("%-8.2f Error: %v\n", x, err)
			continue
		}
		fmt.Printf("%-8.2f %-10.4f\n", x, y)
	}

	fmt.Println("\nKey Points:")
	fmt.Println("• InterpLUT returns exact values at control points")
	fmt.Println("• Interpolation is smooth between points")
	fmt.Println("• Suitable for lookup tables in control systems")
}
