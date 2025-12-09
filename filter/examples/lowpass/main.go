// Package main demonstrates LowPass filter usage for signal smoothing.
//
// This example shows how to use the LowPass filter to smooth noisy signals
// and demonstrates the effect of different gain values on filter behavior.
package main

import (
	"fmt"
	"math"
	"math/rand"

	"control/filter"
)

// generateNoisySignal creates a signal with added random noise
func generateNoisySignal(t, amplitude, frequency, noiseLevel float64) float64 {
	signal := amplitude * math.Sin(2*math.Pi*frequency*t)
	noise := (rand.Float64() - 0.5) * noiseLevel
	return signal + noise
}

func main() {
	fmt.Println("Low-Pass Filter Signal Smoothing Example")
	fmt.Println("=======================================")

	// Create filters with different gains
	lowGainFilter, _ := filter.NewLowPassFilter(0.1)   // Fast response, less smoothing
	medGainFilter, _ := filter.NewLowPassFilter(0.5)   // Balanced
	highGainFilter, _ := filter.NewLowPassFilter(0.9)  // Slow response, more smoothing
	
	fmt.Println("Filter Configuration:")
	fmt.Printf("  Low Gain (0.1):  Fast response, minimal smoothing\n")
	fmt.Printf("  Med Gain (0.5):  Balanced response and smoothing\n")
	fmt.Printf("  High Gain (0.9): Slow response, maximum smoothing\n\n")

	// Signal parameters
	amplitude := 10.0
	frequency := 0.1  // Hz
	noiseLevel := 3.0
	timeStep := 0.1   // seconds
	numSamples := 50

	fmt.Println("Time\tNoisy\tLow(0.1)\tMed(0.5)\tHigh(0.9)")
	fmt.Println("----\t-----\t--------\t--------\t---------")

	var totalNoise, totalLowError, totalMedError, totalHighError float64

	for i := 0; i < numSamples; i++ {
		t := float64(i) * timeStep
		
		// Generate clean signal (for error calculation)
		cleanSignal := amplitude * math.Sin(2*math.Pi*frequency*t)
		
		// Generate noisy measurement
		noisySignal := generateNoisySignal(t, amplitude, frequency, noiseLevel)
		
		// Apply filters
		lowFiltered := lowGainFilter.Estimate(noisySignal)
		medFiltered := medGainFilter.Estimate(noisySignal)
		highFiltered := highGainFilter.Estimate(noisySignal)
		
		// Calculate errors compared to clean signal
		noiseError := math.Abs(noisySignal - cleanSignal)
		lowError := math.Abs(lowFiltered - cleanSignal)
		medError := math.Abs(medFiltered - cleanSignal)
		highError := math.Abs(highFiltered - cleanSignal)
		
		totalNoise += noiseError
		totalLowError += lowError
		totalMedError += medError
		totalHighError += highError
		
		// Print every 5 samples
		if i%5 == 0 || i == numSamples-1 {
			fmt.Printf("%.1f\t%.2f\t%.2f\t\t%.2f\t\t%.2f\n",
				t, noisySignal, lowFiltered, medFiltered, highFiltered)
		}
	}

	fmt.Println()
	fmt.Println("Performance Analysis:")
	fmt.Println("====================")

	avgNoise := totalNoise / float64(numSamples)
	avgLowError := totalLowError / float64(numSamples)
	avgMedError := totalMedError / float64(numSamples)
	avgHighError := totalHighError / float64(numSamples)

	fmt.Printf("Average Errors (vs clean signal):\n")
	fmt.Printf("  Raw Noisy:     %.3f\n", avgNoise)
	fmt.Printf("  Low Gain:      %.3f (%.1f%% improvement)\n", 
		avgLowError, (avgNoise-avgLowError)/avgNoise*100)
	fmt.Printf("  Medium Gain:   %.3f (%.1f%% improvement)\n", 
		avgMedError, (avgNoise-avgMedError)/avgNoise*100)
	fmt.Printf("  High Gain:     %.3f (%.1f%% improvement)\n", 
		avgHighError, (avgNoise-avgHighError)/avgNoise*100)

	fmt.Println()
	fmt.Println("Step Response Comparison:")
	fmt.Println("========================")

	// Reset filters and test step response
	lowGainFilter.Reset()
	medGainFilter.Reset()
	highGainFilter.Reset()

	// Initialize with 0
	lowGainFilter.Estimate(0.0)
	medGainFilter.Estimate(0.0)
	highGainFilter.Estimate(0.0)

	fmt.Println("Step\tLow(0.1)\tMed(0.5)\tHigh(0.9)")
	fmt.Println("----\t--------\t--------\t---------")

	// Apply step input of 10.0
	for step := 1; step <= 10; step++ {
		lowResp := lowGainFilter.Estimate(10.0)
		medResp := medGainFilter.Estimate(10.0)
		highResp := highGainFilter.Estimate(10.0)
		
		fmt.Printf("%d\t%.3f\t\t%.3f\t\t%.3f\n", step, lowResp, medResp, highResp)
		
		// Stop when all are close to final value
		if step >= 5 && lowResp > 9.9 && medResp > 9.9 && highResp > 9.5 {
			break
		}
	}

	fmt.Println()
	fmt.Println("Filter Characteristics:")
	fmt.Println("======================")
	fmt.Println("• Low Gain (0.1):  Fast response, reaches ~99% in 2-3 steps, minimal noise rejection")
	fmt.Println("• Medium Gain (0.5): Moderate response, reaches ~99% in 4-5 steps, good balance")
	fmt.Println("• High Gain (0.9):  Slow response, reaches ~99% in 8-10 steps, excellent noise rejection")

	fmt.Println()
	fmt.Println("Usage Guidelines:")
	fmt.Println("================")
	fmt.Println("• Use LOW gain (0.1-0.3) for:")
	fmt.Println("  - Fast-changing signals")
	fmt.Println("  - Real-time control systems requiring quick response")
	fmt.Println("  - When signal already has low noise")
	
	fmt.Println()
	fmt.Println("• Use MEDIUM gain (0.4-0.6) for:")
	fmt.Println("  - General-purpose filtering")
	fmt.Println("  - Balanced noise rejection and response speed")
	fmt.Println("  - Most sensor filtering applications")
	
	fmt.Println()
	fmt.Println("• Use HIGH gain (0.7-0.9) for:")
	fmt.Println("  - Very noisy signals")
	fmt.Println("  - Slowly changing processes")
	fmt.Println("  - When smooth output is more important than fast response")
}