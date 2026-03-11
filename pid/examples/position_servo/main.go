// Package main demonstrates PID control for servo position tracking.
//
// This example shows position control with a PID controller,
// demonstrating precise positioning behavior.
package main

import (
	"fmt"

	"control/pid"
)

type servoPlant struct {
	position float64
	rate     float64
	maxRate  float64
	response float64
}

func (s *servoPlant) update(command, dt float64) {
	desiredRate := command * s.maxRate
	s.rate += (desiredRate - s.rate) * s.response * dt
	s.position += s.rate * dt
}

func main() {
	fmt.Println("Position Servo Control Example")
	fmt.Println("==============================")
	fmt.Println()

	fmt.Println("PID Position Tracking")
	fmt.Println("--------------------")
	fmt.Println("Demonstrates closed-loop servo position control with a simple plant model.")
	fmt.Println()

	// Create PID controller for position
	// Output is servo command (-1.0 to 1.0)
	controller := pid.New(0.07, 0.025, 0.02,
		pid.WithOutputLimits(-1.0, 1.0))

	fmt.Println("Controller Configuration:")
	fmt.Println("  Kp = 0.07 (proportional gain)")
	fmt.Println("  Ki = 0.025 (integral gain)")
	fmt.Println("  Kd = 0.02 (derivative gain)")
	fmt.Println("  Output limits: -1.0 to 1.0 (servo command)")
	fmt.Println()
	fmt.Println("Plant Model: velocity-limited servo with first-order rate response")
	fmt.Println("  maxRate = 120 deg/s, response = 8.0 1/s")
	fmt.Println()

	dt := 0.05

	runScenario := func(title string, initialPosition, initialRate float64, steps int, targetAtStep func(step int) float64) {
		plant := servoPlant{
			position: initialPosition,
			rate:     initialRate,
			maxRate:  120.0,
			response: 8.0,
		}

		controller.Reset()
		initialTarget := targetAtStep(0)
		controller.CalculateWithDt(initialTarget, plant.position, dt)

		fmt.Println(title)
		fmt.Printf("\n%-6s %-8s %-11s %-12s %-10s %-10s %-10s\n", "Step", "Time", "Target (°)", "Current (°)", "Error (°)", "Command", "Rate (°/s)")
		fmt.Printf("%-6s %-8s %-10s %-11s %-9s %-10s %-9s\n", "----", "----", "----------", "-----------", "---------", "-------", "---------")

		for i := 0; i < steps; i++ {
			target := targetAtStep(i)
			command := controller.CalculateWithDt(target, plant.position, dt)
			error := target - plant.position

			if i%4 == 0 || i == steps-1 {
				fmt.Printf("%-6d %-8.2f %-11.1f %-12.2f %-10.2f %-10.3f %-10.2f\n",
					i+1, float64(i+1)*dt, target, plant.position, error, command, plant.rate)
			}

			plant.update(command, dt)
		}

		fmt.Println()
	}

	runScenario(
		"Test 1: Move to Position 90°",
		0.0,
		0.0,
		320,
		func(step int) float64 { return 90.0 },
	)

	runScenario(
		"Test 2: Return to Position 30°",
		90.0,
		0.0,
		300,
		func(step int) float64 { return 30.0 },
	)

	runScenario(
		"Test 3: Small Adjustment (30° -> 35°)",
		30.0,
		0.0,
		50,
		func(step int) float64 { return 35.0 },
	)

	fmt.Println("\nAnalysis:")
	fmt.Println("---------")
	fmt.Println()
	fmt.Println("Large movement (0° -> 90°):")
	fmt.Println("  • Command saturates early to accelerate quickly")
	fmt.Println("  • As position approaches target, command tapers and the rate decreases")
	fmt.Println("  • Rate column shows actuator speed limiting and response lag")
	fmt.Println()
	fmt.Println("Reverse movement (90° -> 30°):")
	fmt.Println("  • Controller commands negative torque to decelerate and reverse direction")
	fmt.Println("  • The rate state transitions smoothly through zero as the servo changes direction")
	fmt.Println()
	fmt.Println("Fine adjustment (30° -> 35°):")
	fmt.Println("  • Smaller error yields much smaller command than large-setpoint moves")
	fmt.Println("  • PID settles with low residual rate and small steady-state error")
	fmt.Println()
	fmt.Println("Key Points:")
	fmt.Println("  • P term provides position correction force")
	fmt.Println("  • I term eliminates positioning error")
	fmt.Println("  • D term damps fast command changes near target")
	fmt.Println("  • Output limits protect servo mechanics")
}
