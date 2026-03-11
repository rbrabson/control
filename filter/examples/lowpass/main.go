// Package main demonstrates first-order low-pass filter behavior.
//
// This example shows the basic operation of a low-pass filter,
// matching the behavior tested in LowPassFilterTest.java.
package main

import (
	"fmt"
	"math"

	"control/filter"
)

func main() {
	fmt.Println("Low-Pass Filter: First-Order Filtering")
	fmt.Println("======================================")
	fmt.Println()

	// Create a low-pass filter with alpha = 0.5
	// Formula: output = alpha * previous + (1-alpha) * measurement
	lpf, err := filter.NewLowPassFilter(0.5)
	if err != nil {
		fmt.Printf("Error creating filter: %v\n", err)
		return
	}

	fmt.Println("Filter Configuration:")
	fmt.Printf("  Alpha: %.1f\n", lpf.GetAlpha())
	fmt.Println("  Formula: output = 0.5 * previous + 0.5 * measurement")
	fmt.Println()

	// Test the filter behavior matching Java test
	fmt.Println("Testing filter behavior:")
	fmt.Printf("%-6s %-8s %-12s %-24s\n", "Step", "Input", "Output", "Expected")
	fmt.Printf("%-6s %-8s %-12s %-24s\n", "----", "-----", "------", "--------")

	// First estimate: initializes with input value
	input1 := 10.0
	output1 := lpf.Estimate(input1)
	fmt.Printf("%-6d %-8.1f %-12.6f %-24s\n", 1, input1, output1, fmt.Sprintf("%.1f (initialization)", input1))

	// Verify first output matches
	if math.Abs(output1-10.0) < 1e-9 {
		fmt.Printf("%-6s %-8s %-12s %-24s\n", "", "", "", "✓ Match")
	}

	// Second estimate: applies filtering
	// Expected: 0.5 * 10.0 + 0.5 * 20.0 = 15.0
	input2 := 20.0
	output2 := lpf.Estimate(input2)
	expected2 := 15.0
	fmt.Printf("%-6d %-8.1f %-12.6f %-24.1f\n", 2, input2, output2, expected2)

	// Verify second output matches
	if math.Abs(output2-15.0) < 1e-9 {
		fmt.Printf("%-6s %-8s %-12s %-24s\n", "", "", "", "✓ Match")
	}

	// Continue to show convergence
	fmt.Println("\nContinued filtering (converging to input):")
	fmt.Printf("%-6s %-8s %-12s\n", "Step", "Input", "Output")
	fmt.Printf("%-6s %-8s %-12s\n", "----", "-----", "------")

	for i := 3; i <= 10; i++ {
		output := lpf.Estimate(20.0)
		fmt.Printf("%-6d %-8.1f %-12.6f\n", i, 20.0, output)
	}

	// Demonstrate step response
	fmt.Println("\nStep Response (input changes from 20 to 40):")
	fmt.Printf("%-6s %-8s %-12s\n", "Step", "Input", "Output")
	fmt.Printf("%-6s %-8s %-12s\n", "----", "-----", "------")

	for i := 1; i <= 8; i++ {
		output := lpf.Estimate(40.0)
		fmt.Printf("%-6d %-8.1f %-12.6f\n", i, 40.0, output)
	}

	fmt.Println("\nKey Points:")
	fmt.Println("• First call initializes filter with input value")
	fmt.Println("• Subsequent calls apply exponential smoothing")
	fmt.Println("• Alpha controls smoothing (0.5 = balanced)")
	fmt.Println("• Output gradually converges to input value")
	fmt.Println("• Higher alpha = more smoothing, slower response")
	fmt.Println("• Lower alpha = less smoothing, faster response")
}
