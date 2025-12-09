// Temperature control example using InterpLUT for
// non-linear temperature response curves.
package main

import (
	"fmt"
	"log"
	"math"

	"control/interplut"
)

func main() {
	// Create lookup table for heater power vs desired temperature
	fmt.Println("Creating temperature control lookup table...")

	lut := interplut.New()

	// Add control points: temperature (°C) -> power (%)
	lut.Add(20.0, 0.0)   // Room temperature, no power needed
	lut.Add(50.0, 0.15)  // Low heat applications
	lut.Add(100.0, 0.45) // Boiling point
	lut.Add(150.0, 0.70) // Medium industrial heat
	lut.Add(200.0, 0.85) // High temperature applications
	lut.Add(250.0, 1.0)  // Maximum heater capacity

	err := lut.CreateLUT()
	if err != nil {
		log.Fatal("Error creating temperature LUT:", err)
	}

	fmt.Println("\nTemperature to Power Mapping:")
	fmt.Println("Temperature (°C) | Power (%) | Notes")
	fmt.Println("-----------------|-----------|--------")

	// Test at various temperatures
	temperatures := []struct {
		temp  float64
		notes string
	}{
		{25.0, "Room temp +5°C"},
		{60.0, "Hot water"},
		{85.0, "Coffee brewing"},
		{120.0, "Steam generation"},
		{175.0, "Industrial process"},
		{225.0, "High temp process"},
	}

	for _, t := range temperatures {
		power, err := lut.Get(t.temp)
		if err != nil {
			fmt.Printf("Error at %.1f°C: %v\n", t.temp, err)
			continue
		}
		fmt.Printf("      %.1f       |   %.3f   | %s\n", t.temp, power, t.notes)
	}

	// Show smooth curve by sampling many points
	fmt.Println("\nTemperature response curve (every 20°C):")
	for temp := 20.0; temp <= 250.0; temp += 20.0 {
		power, _ := lut.Get(temp)
		// Create simple ASCII bar chart
		barLength := int(power * 40) // Scale to 40 characters max
		bar := ""
		for i := 0; i < barLength; i++ {
			bar += "█"
		}
		fmt.Printf("%3.0f°C |%-40s| %.3f\n", temp, bar, power)
	}

	// Compare with simple linear approach
	fmt.Println("\nComparison at intermediate points:")
	testTemp := 125.0
	power, _ := lut.Get(testTemp)
	fmt.Printf("At %.1f°C: InterpLUT gives %.3f power\n", testTemp, power)
	fmt.Printf("(Smooth cubic spline interpolation between control points)")

	// Show monotonicity preservation
	fmt.Println("\nMonotonicity check (power should always increase):")
	prev := 0.0
	monotonic := true
	for temp := 20.0; temp <= 250.0; temp += 5.0 {
		power, _ := lut.Get(temp)
		if power < prev {
			monotonic = false
			break
		}
		prev = power
	}
	if monotonic {
		fmt.Println("✓ Monotonicity preserved - power increases smoothly with temperature")
	} else {
		fmt.Println("✗ Monotonicity violated - unexpected power decrease detected")
	}
}

// Helper function for demonstration
func _() {
	// Placeholder to avoid unused function error
	_ = math.Abs
}
