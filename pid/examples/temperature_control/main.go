// Package main demonstrates PID control for temperature regulation.
//
// This example shows temperature control with PID, demonstrating
// heating control for thermal systems.
package main

import (
	"fmt"

	"control/pid"
)

type thermalPlant struct {
	temperature float64
	ambient     float64
	heatGain    float64
	cooling     float64
}

func (p *thermalPlant) update(heaterPower, dt float64) {
	delta := (p.heatGain*heaterPower - p.cooling*(p.temperature-p.ambient)) * dt
	p.temperature += delta
}

func (p *thermalPlant) applyDisturbance(deltaC float64) {
	p.temperature += deltaC
}

func main() {
	fmt.Println("Temperature Control Example")
	fmt.Println("==========================")
	fmt.Println()

	fmt.Println("PID Temperature Regulation")
	fmt.Println("-------------------------")
	fmt.Println("Demonstrates closed-loop heater control with a simple thermal plant model.")
	fmt.Println()

	// Create PID controller for temperature
	// Output is heater power (0.0 to 1.0)
	controller := pid.New(0.2, 0.04, 0.0,
		pid.WithOutputLimits(0.0, 1.0))

	fmt.Println("Controller Configuration:")
	fmt.Println("  Kp = 0.2 (proportional gain)")
	fmt.Println("  Ki = 0.04 (integral gain)")
	fmt.Println("  Kd = 0.0 (derivative gain)")
	fmt.Println("  Output limits: 0.0 to 1.0 (heater power)")
	fmt.Println()
	fmt.Println("Plant Model: dT/dt = heatGain*power - cooling*(T-ambient)")
	fmt.Println("  ambient = 20°C, heatGain = 4.0, cooling = 0.03")
	fmt.Println()

	dt := 0.5

	runScenario := func(title string, initialTemp float64, steps int, targetAtStep func(step int) float64, disturbanceAtStep map[int]float64) {
		plant := thermalPlant{
			temperature: initialTemp,
			ambient:     20.0,
			heatGain:    4.0,
			cooling:     0.03,
		}

		controller.Reset()
		initialTarget := targetAtStep(0)
		controller.CalculateWithDt(initialTarget, plant.temperature, dt)

		fmt.Println(title)
		fmt.Printf("\n%-6s %-8s %-12s %-13s %-11s %-10s\n", "Step", "Time", "Target (°C)", "Current (°C)", "Error (°C)", "Power (%)")
		fmt.Printf("%-6s %-8s %-12s %-13s %-11s %-10s\n", "----", "----", "-----------", "------------", "----------", "---------")

		for i := 0; i < steps; i++ {
			if delta, ok := disturbanceAtStep[i]; ok {
				plant.applyDisturbance(delta)
			}

			target := targetAtStep(i)
			power := controller.CalculateWithDt(target, plant.temperature, dt)
			error := target - plant.temperature

			if i%2 == 0 || i == steps-1 {
				fmt.Printf("%-6d %-8.1f %-12.1f %-13.2f %-11.2f %-10.1f\n",
					i+1, float64(i+1)*dt, target, plant.temperature, error, power*100)
			}

			plant.update(power, dt)
		}

		fmt.Println()
	}

	runScenario(
		"Test 1: Heat to 75°C (from 20°C ambient)",
		20.0,
		40,
		func(step int) float64 { return 75.0 },
		nil,
	)

	runScenario(
		"Test 2: Disturbance Recovery at 75°C (door opens at step 8)",
		75.0,
		30,
		func(step int) float64 { return 75.0 },
		map[int]float64{8: -10.0},
	)

	runScenario(
		"Test 3: Setpoint Change (75°C -> 85°C at step 12)",
		75.0,
		36,
		func(step int) float64 {
			if step < 12 {
				return 75.0
			}
			return 85.0
		},
		nil,
	)

	fmt.Println("\nAnalysis:")
	fmt.Println("---------")
	fmt.Println()
	fmt.Println("Initial heating (20°C -> 75°C):")
	fmt.Println("  • Controller saturates power initially, then tapers as temperature rises")
	fmt.Println("  • As error shrinks, output settles to the steady-state power needed to hold temperature")
	fmt.Println()
	fmt.Println("Disturbance recovery (10°C drop):")
	fmt.Println("  • Power jumps after disturbance, then smoothly decays as the room reheats")
	fmt.Println("  • Closed-loop dynamics show realistic thermal lag rather than scripted measurements")
	fmt.Println()
	fmt.Println("Setpoint change (75°C -> 85°C):")
	fmt.Println("  • Controller responds with a transient power increase and then rebalances near new setpoint")
	fmt.Println("  • Integral action removes residual offset once proportional response slows")
	fmt.Println()
	fmt.Println("Key Points:")
	fmt.Println("  • P term: Immediate response to temperature error")
	fmt.Println("  • I term: Ensures reaching exact temperature")
	fmt.Println("  • D term: Optional for thermal systems; this tuning uses PI for smoother hold behavior")
	fmt.Println("  • One-sided control: heater can't cool (min = 0%)")
}
