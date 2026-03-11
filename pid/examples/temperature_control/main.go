// Package main demonstrates PID control for temperature regulation.
//
// This example shows temperature control with PID, demonstrating
// heating control for thermal systems.
package main

import (
	"fmt"

	"control/pid"
)

func main() {
	fmt.Println("Temperature Control Example")
	fmt.Println("==========================")
	fmt.Println()

	fmt.Println("PID Temperature Regulation")
	fmt.Println("-------------------------")
	fmt.Println("Demonstrates heater control with PID.")
	fmt.Println()

	// Create PID controller for temperature
	// Output is heater power (0.0 to 1.0)
	controller := pid.New(0.08, 0.015, 0.25,
		pid.WithOutputLimits(0.0, 1.0))

	fmt.Println("Controller Configuration:")
	fmt.Println("  Kp = 0.08 (proportional gain)")
	fmt.Println("  Ki = 0.015 (integral gain)")
	fmt.Println("  Kd = 0.25 (derivative gain)")
	fmt.Println("  Output limits: 0.0 to 1.0 (heater power)")
	fmt.Println()

	// Test: Heat up to target temperature
	fmt.Println("Test: Heating to 75°C (from 20°C ambient)")
	fmt.Printf("\n%-6s %-12s %-13s %-11s %-10s\n", "Step", "Target (°C)", "Current (°C)", "Error (°C)", "Power (%)")
	fmt.Printf("%-6s %-12s %-13s %-11s %-10s\n", "----", "-----------", "------------", "----------", "---------")

	targetTemp := 75.0
	dt := 0.1

	// Simulated temperatures (gradually heating)
	currentTemps := []float64{
		20.0, 28.0, 38.0, 48.0, 56.0, 62.0, 67.0, 70.5, 72.5, 73.8, 74.5, 75.0,
	}
	controller.Calculate(targetTemp, currentTemps[0])

	for i, temp := range currentTemps {
		power := controller.CalculateWithDt(targetTemp, temp, dt)
		error := targetTemp - temp
		powerPct := power * 100.0

		fmt.Printf("%-6d %-12.0f %-13.1f %-11.1f %-10.1f\n",
			i+1, targetTemp, temp, error, powerPct)
	}

	// Test: Temperature drop recovery
	controller.Reset()

	fmt.Println("\nTest: Recovery from Temperature Drop")
	fmt.Println("(Door opened, temperature drops to 65°C)")
	fmt.Printf("\n%-6s %-12s %-13s %-11s %-10s\n", "Step", "Target (°C)", "Current (°C)", "Error (°C)", "Power (%)")
	fmt.Printf("%-6s %-12s %-13s %-11s %-10s\n", "----", "-----------", "------------", "----------", "---------")

	targetTemp = 75.0

	// Temperature recovering after disturbance
	tempsRecovery := []float64{
		75.0, // At setpoint
		65.0, // Sudden drop (disturbance)
		67.5, // Recovering
		70.0,
		72.0,
		73.5,
		74.5,
		75.0, // Back to setpoint
	}
	controller.Calculate(targetTemp, tempsRecovery[0])

	for i, temp := range tempsRecovery {
		power := controller.CalculateWithDt(targetTemp, temp, dt)
		error := targetTemp - temp
		powerPct := power * 100.0

		fmt.Printf("%-6d %-12.0f %-13.1f %-11.1f %-10.1f\n",
			i+1, targetTemp, temp, error, powerPct)
	}

	// Test: Setpoint change
	controller.Reset()

	fmt.Println("\nTest: Setpoint Change (75°C → 85°C)")
	fmt.Printf("\n%-6s %-12s %-13s %-11s %-10s\n", "Step", "Target (°C)", "Current (°C)", "Error (°C)", "Power (%)")
	fmt.Printf("%-6s %-12s %-13s %-11s %-10s\n", "----", "-----------", "------------", "----------", "---------")

	targetTemp = 85.0

	// Temperature increasing to new setpoint
	tempsIncrease := []float64{
		75.0, 77.5, 80.0, 82.0, 83.5, 84.5, 85.0,
	}
	controller.Calculate(targetTemp, tempsIncrease[0])

	for i, temp := range tempsIncrease {
		power := controller.CalculateWithDt(targetTemp, temp, dt)
		error := targetTemp - temp
		powerPct := power * 100.0

		fmt.Printf("%-6d %-12.0f %-13.1f %-11.1f %-10.1f\n",
			i+1, targetTemp, temp, error, powerPct)
	}

	fmt.Println("\nAnalysis:")
	fmt.Println("---------")
	fmt.Println()
	fmt.Println("Initial heating (20°C → 75°C):")
	fmt.Println("  • High power at start due to large error")
	fmt.Println("  • Power reduces as temperature approaches target")
	fmt.Println("  • Near setpoint: low power for maintenance")
	fmt.Println()
	fmt.Println("Disturbance recovery (drop to 65°C):")
	fmt.Println("  • Controller immediately increases power")
	fmt.Println("  • Quick response to temperature loss")
	fmt.Println("  • D term detects rapid change")
	fmt.Println()
	fmt.Println("Setpoint change (75°C → 85°C):")
	fmt.Println("  • Smooth transition to new temperature")
	fmt.Println("  • Controlled heating rate")
	fmt.Println()
	fmt.Println("Key Points:")
	fmt.Println("  • P term: Immediate response to temperature error")
	fmt.Println("  • I term: Ensures reaching exact temperature")
	fmt.Println("  • D term: Prevents overshoot and oscillation")
	fmt.Println("  • One-sided control: heater can't cool (min = 0%)")
}
