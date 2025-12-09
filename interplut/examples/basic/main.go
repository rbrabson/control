// Basic InterpLUT example demonstrating smooth interpolation
// between control points.
package main

import (
	"fmt"
	"log"

	"control/interplut"
)

func main() {
	// Create a new interpolated lookup table
	lut := interplut.New()

	// Add control points (distance -> velocity mapping for a shooter)
	fmt.Println("Adding control points for shooter velocity vs distance...")
	lut.Add(1.1, 0.2)  // At 1.1m distance, use 20% velocity
	lut.Add(2.7, 0.5)  // At 2.7m distance, use 50% velocity
	lut.Add(3.6, 0.75) // At 3.6m distance, use 75% velocity
	lut.Add(4.1, 0.9)  // At 4.1m distance, use 90% velocity
	lut.Add(5.0, 1.0)  // At 5.0m distance, use 100% velocity

	// Create the lookup table (computes spline coefficients)
	err := lut.CreateLUT()
	if err != nil {
		log.Fatal("Error creating LUT:", err)
	}

	fmt.Println("\nInterpolated values:")
	fmt.Println("Distance (m) | Velocity (%)")
	fmt.Println("-------------|-------------")

	// Test interpolation at various distances
	distances := []float64{1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5}

	for _, distance := range distances {
		velocity, err := lut.Get(distance)
		if err != nil {
			fmt.Printf("Error at %.1f: %v\n", distance, err)
			continue
		}
		fmt.Printf("    %.1f     |    %.3f\n", distance, velocity)
	}

	// Demonstrate exact matches at control points
	fmt.Println("\nExact matches at control points:")
	fmt.Println("Distance (m) | Velocity (%) | Expected")
	fmt.Println("-------------|--------------|----------")

	controlPoints := []struct{ x, y float64 }{
		{1.1, 0.2},
		{2.7, 0.5},
		{3.6, 0.75},
		{4.1, 0.9},
		{5.0, 1.0},
	}

	for _, cp := range controlPoints {
		velocity, err := lut.Get(cp.x)
		if err != nil {
			fmt.Printf("Error at %.1f: %v\n", cp.x, err)
			continue
		}
		fmt.Printf("    %.1f     |    %.3f    |   %.3f\n", cp.x, velocity, cp.y)
	}

	// Show the spline representation for debugging
	fmt.Println("\nSpline representation:")
	fmt.Println(lut.String())
}
