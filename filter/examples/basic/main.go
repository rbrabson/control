// Package main demonstrates basic Kalman filter usage for signal estimation.
//
// This example shows how to use the Kalman filter to estimate a signal
// in the presence of measurement noise and system dynamics.
package main

import (
	"fmt"
	"math"
	"math/rand"

	"control/filter"
)

// simulateNoisyMeasurement adds Gaussian noise to a true signal value
func simulateNoisyMeasurement(trueValue, noiseStdDev float64) float64 {
	noise := rand.NormFloat64() * noiseStdDev
	return trueValue + noise
}

// trueSignal generates a slowly varying signal with trend
func trueSignal(t float64) float64 {
	// A signal that increases linearly with some sinusoidal variation
	return 10.0 + 0.5*t + 2.0*math.Sin(0.1*t)
}

func main() {
	fmt.Println("Kalman Filter Signal Estimation Example")
	fmt.Println("======================================")
	rand.Seed(42)

	// Create Kalman filter
	// Q (process noise): 0.01 - low process noise (signal changes slowly)
	// R (measurement noise): 0.5 - moderate measurement noise
	// N (history size): 5 - use 5 recent estimates for linear regression
	kf, err := filter.NewKalmanFilter(0.01, 0.5, 5)
	if err != nil {
		fmt.Printf("Error creating Kalman filter: %v\n", err)
		return
	}

	fmt.Printf("Filter Parameters: Q=%.2f, R=%.2f, N=%d\n", 0.01, 0.5, 5)
	fmt.Printf("Kalman Gain (K): %.4f\n", kf.GetK())
	fmt.Printf("Error Covariance (P): %.4f\n\n", kf.GetP())

	// Simulation parameters
	numSteps := 50
	timeStep := 0.2
	measurementNoise := 0.7 // Standard deviation of measurement noise

	fmt.Printf("%-8s %-8s %-8s %-8s %-8s\n", "Time", "True", "Noisy", "Kalman", "Error")
	fmt.Printf("%-8s %-8s %-8s %-8s %-8s\n", "----", "----", "-----", "------", "-----")

	var totalError, totalKalmanError float64
	warmupSteps := 5

	for i := 0; i < numSteps; i++ {
		t := float64(i) * timeStep

		// Generate true signal value
		trueValue := trueSignal(t)

		// Create noisy measurement
		noisyMeasurement := simulateNoisyMeasurement(trueValue, measurementNoise)

		// Apply Kalman filter
		kalmanEstimate := kf.Estimate(noisyMeasurement)

		// Calculate errors
		measurementError := math.Abs(noisyMeasurement - trueValue)
		kalmanError := math.Abs(kalmanEstimate - trueValue)

		if i >= warmupSteps {
			totalError += measurementError
			totalKalmanError += kalmanError
		}

		// Print results every 5 steps
		if i%5 == 0 || i == numSteps-1 {
			fmt.Printf("%-8.1f %-8.2f %-8.2f %-8.2f %-8.2f\n",
				t, trueValue, noisyMeasurement, kalmanEstimate, kalmanError)
		}
	}

	fmt.Println()
	fmt.Println("Performance Summary:")
	fmt.Println("===================")

	effectiveSteps := float64(numSteps - warmupSteps)
	avgMeasurementError := totalError / effectiveSteps
	avgKalmanError := totalKalmanError / effectiveSteps
	errorChangePercent := ((avgKalmanError - avgMeasurementError) / avgMeasurementError) * 100

	fmt.Printf("Metrics exclude first %d warmup samples\n", warmupSteps)
	fmt.Printf("Average Measurement Error: %.3f\n", avgMeasurementError)
	fmt.Printf("Average Kalman Error:      %.3f\n", avgKalmanError)
	if errorChangePercent <= 0 {
		fmt.Printf("Error Reduction:           %.1f%%\n", -errorChangePercent)
	} else {
		fmt.Printf("Error Increase:            %.1f%%\n", errorChangePercent)
	}

	fmt.Println()
	fmt.Println("Filter Details:")
	fmt.Println("==============")
	fmt.Printf("Final State Estimate: %.3f\n", kf.GetX())
	fmt.Printf("Kalman Gain:          %.4f\n", kf.GetK())
	fmt.Printf("Error Covariance:     %.4f\n", kf.GetP())

	fmt.Println()
	fmt.Println("How it works:")
	fmt.Println("============")
	fmt.Println("1. Process Model: Uses linear regression on recent estimates to predict next value")
	fmt.Println("2. Measurement Update: Combines prediction with noisy measurement using Kalman gain")
	fmt.Println("3. DARE Solution: Automatically calculates optimal Kalman gain for given noise levels")
	fmt.Println("4. Adaptive Filtering: Balances between trusting the model vs. measurements")

	fmt.Println()
	fmt.Println("Key Parameters:")
	fmt.Println("==============")
	fmt.Println("• Q (Process Noise): Lower values trust the model more")
	fmt.Println("• R (Measurement Noise): Higher values trust measurements less")
	fmt.Println("• N (History Size): Number of recent estimates used for trend prediction")
	fmt.Println("• Kalman Gain K: Automatically computed, determines measurement vs. model weighting")
}
