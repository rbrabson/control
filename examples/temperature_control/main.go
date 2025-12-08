// Temperature Control Example
//
// This example demonstrates using a PID controller for temperature regulation,
// such as in an oven, heater, or climate control system. It shows:
// - Feed-forward control for ambient temperature compensation
// - Integral reset on zero crossover to prevent overshoot
// - Derivative filtering for temperature sensor noise
// - Realistic thermal system dynamics

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"control/pid"
)

// ThermalSystem represents a simulated heating system
type ThermalSystem struct {
	temperature    float64   // Current temperature (°C)
	ambientTemp    float64   // Ambient temperature (°C)
	heatCapacity   float64   // Thermal mass
	heatLoss       float64   // Heat loss coefficient
	heaterPower    float64   // Current heater power (0-1)
	maxHeaterPower float64   // Maximum heater power (watts)
	lastUpdate     time.Time
	sensorNoise    float64   // Temperature sensor noise (°C)
}

// NewThermalSystem creates a new simulated thermal system
func NewThermalSystem(ambientTemp, heatCapacity, heatLoss, maxPower, sensorNoise float64) *ThermalSystem {
	return &ThermalSystem{
		temperature:    ambientTemp, // Start at ambient temperature
		ambientTemp:    ambientTemp,
		heatCapacity:   heatCapacity,
		heatLoss:       heatLoss,
		heaterPower:    0.0,
		maxHeaterPower: maxPower,
		lastUpdate:     time.Now(),
		sensorNoise:    sensorNoise,
	}
}

// ApplyHeaterPower applies heater power (0.0 to 1.0)
func (ts *ThermalSystem) ApplyHeaterPower(power float64) {
	now := time.Now()
	dt := now.Sub(ts.lastUpdate).Seconds()
	ts.lastUpdate = now

	// Clamp power to valid range
	if power > 1.0 {
		power = 1.0
	} else if power < 0.0 {
		power = 0.0
	}

	ts.heaterPower = power

	if dt > 0 {
		// Heat input from heater
		heatInput := power * ts.maxHeaterPower

		// Heat loss to ambient (proportional to temperature difference)
		heatLoss := ts.heatLoss * (ts.temperature - ts.ambientTemp)

		// Net heat change
		netHeat := heatInput - heatLoss

		// Temperature change (dT = Q / (mass * specific_heat))
		tempChange := netHeat * dt / ts.heatCapacity

		ts.temperature += tempChange
	}
}

// GetTemperature returns current temperature with sensor noise
func (ts *ThermalSystem) GetTemperature() float64 {
	// Add realistic sensor noise
	noise := (rand.Float64() - 0.5) * ts.sensorNoise * 2.0
	return ts.temperature + noise
}

// GetActualTemperature returns actual temperature without noise
func (ts *ThermalSystem) GetActualTemperature() float64 {
	return ts.temperature
}

// SetAmbientTemperature changes the ambient temperature (disturbance)
func (ts *ThermalSystem) SetAmbientTemperature(temp float64) {
	ts.ambientTemp = temp
}

func main() {
	fmt.Println("Temperature Control Example")
	fmt.Println("===========================")
	fmt.Println()

	// Ambient temperature compensation (feed-forward)
	ambientTemp := 20.0 // °C
	feedForwardGain := 0.02 // Rough estimate of power needed per degree above ambient

	// Create temperature controller with advanced features
	controller := pid.New(0.5, 0.1, 0.02,
		pid.WithFeedForward(ambientTemp*feedForwardGain), // Ambient compensation
		pid.WithIntegralResetOnZeroCross(),               // Prevent overshoot when crossing target
		pid.WithDerivativeFilter(0.2),                    // Filter temperature sensor noise (20% filter)
	)

	// Set heater power limits (0% to 100%)
	controller.SetOutputLimits(0.0, 1.0)

	// Create simulated thermal system
	// (ambient: 20°C, heat capacity: 1000 J/°C, heat loss: 50 W/°C, max power: 2000W, noise: 0.5°C)
	thermalSystem := NewThermalSystem(20.0, 1000.0, 50.0, 2000.0, 0.5)

	// Test scenario: heat up, then change setpoint, then ambient disturbance
	scenarios := []struct {
		setpoint    float64
		duration    float64
		description string
		ambientTemp float64
	}{
		{60.0, 15.0, "Initial heating to 60°C", 20.0},
		{80.0, 10.0, "Increase setpoint to 80°C", 20.0},
		{80.0, 10.0, "Ambient temperature rises to 25°C", 25.0},
		{70.0, 8.0, "Reduce setpoint to 70°C", 25.0},
	}

	updateRate := 10 // 10Hz update rate (typical for thermal systems)
	interval := time.Duration(1000/updateRate) * time.Millisecond

	fmt.Printf("Update Rate: %d Hz\n", updateRate)
	fmt.Printf("Ambient Temperature: %.1f°C\n", ambientTemp)
	fmt.Println()
	fmt.Println("Time\tSetpoint\tMeasured\tActual\t\tError\tPower%")
	fmt.Println("----\t--------\t--------\t------\t\t-----\t------")

	totalTime := 0.0

	for _, scenario := range scenarios {
		fmt.Printf("\n%s\n", scenario.description)
		
		// Update ambient temperature if changed
		thermalSystem.SetAmbientTemperature(scenario.ambientTemp)
		
		// Update feed-forward for new ambient temperature
		controller.SetFeedForward(scenario.ambientTemp * feedForwardGain)
		
		startTime := time.Now()
		
		for time.Since(startTime).Seconds() < scenario.duration {
			// Get current temperature with noise
			measuredTemp := thermalSystem.GetTemperature()
			actualTemp := thermalSystem.GetActualTemperature()
			
			// Calculate error
			error := scenario.setpoint - measuredTemp
			
			// Update PID controller
			heaterPower := controller.Update(error)
			
			// Apply power to thermal system
			thermalSystem.ApplyHeaterPower(heaterPower)
			
			// Print status every 0.5 seconds
			elapsed := time.Since(startTime).Seconds()
			if math.Mod(elapsed, 0.5) < float64(interval.Seconds()) {
				fmt.Printf("%.1f\t%.1f\t\t%.2f\t\t%.2f\t\t%.2f\t%.1f%%\n", 
					totalTime+elapsed, scenario.setpoint, measuredTemp, actualTemp, error, heaterPower*100)
			}
			
			time.Sleep(interval)
		}
		
		totalTime += scenario.duration
	}

	fmt.Println()
	fmt.Println("Temperature Control Completed")
	
	// Display final system state
	fmt.Printf("Final temperature: %.2f°C\n", thermalSystem.GetActualTemperature())
	fmt.Printf("Final ambient: %.1f°C\n", thermalSystem.ambientTemp)
	
	// Display controller configuration
	kp, ki, kd := controller.GetGains()
	fmt.Printf("Controller gains - Kp: %.1f, Ki: %.1f, Kd: %.3f\n", kp, ki, kd)
	fmt.Printf("Feed-forward value: %.3f\n", controller.GetFeedForward())
	fmt.Printf("Integral reset enabled: %t\n", controller.GetIntegralResetOnZeroCross())
	fmt.Printf("Derivative filter: %.1f\n", controller.GetDerivativeFilter())

	fmt.Println("\nTemperature control demonstrates:")
	fmt.Println("- Feed-forward compensation for ambient temperature")
	fmt.Println("- Integral reset to prevent overshoot when crossing target")
	fmt.Println("- Derivative filtering to handle sensor noise")
	fmt.Println("- Response to setpoint changes and disturbances")
}