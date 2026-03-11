// Package main demonstrates PID control for motor speed regulation.
//
// This example shows basic motor speed control with PID, demonstrating
// how the controller tracks a target speed.
package main

import (
	"fmt"
	"math"

	"control/pid"
)

type motorPlant struct {
	speed float64
}

const maxNoLoadRPM = 1200.0
const motorTimeConstant = 0.45

func (m *motorPlant) update(power, dt float64) {
	targetSpeed := power * maxNoLoadRPM
	m.speed += (targetSpeed - m.speed) * dt / motorTimeConstant

	if math.Abs(m.speed) < 1e-9 {
		m.speed = 0
	}
}

func runScenario(controller *pid.PID, plant *motorPlant, targetSpeed float64, steps int, dt float64) {
	controller.SetFeedForward(targetSpeed / maxNoLoadRPM)
	controller.Calculate(targetSpeed, plant.speed)

	for i := 0; i < steps; i++ {
		power := controller.CalculateWithDt(targetSpeed, plant.speed, dt)
		error := targetSpeed - plant.speed

		fmt.Printf("%-6d %-13.0f %-14.1f %-10.1f %-8.3f\n",
			i+1, targetSpeed, plant.speed, error, power)

		plant.update(power, dt)
	}
}

func main() {
	fmt.Println("Motor Speed Control Example")
	fmt.Println("==========================")
	fmt.Println()

	fmt.Println("PID Speed Tracking")
	fmt.Println("------------------")
	fmt.Println("Demonstrates motor speed control with PID.")
	fmt.Println()
	fmt.Println("Plant Model: first-order motor response with inertia and drag")
	fmt.Println()

	// Create PID controller for motor speed
	// Output is motor power (-1.0 to 1.0)
	controller := pid.New(0.0014, 0.0025, 0.00006,
		pid.WithOutputLimits(-1.0, 1.0))

	fmt.Println("Controller Configuration:")
	fmt.Println("  Kp = 0.0014 (proportional gain)")
	fmt.Println("  Ki = 0.0025 (integral gain)")
	fmt.Println("  Kd = 0.00006 (derivative gain)")
	fmt.Println("  Feed-forward = target speed / max speed")
	fmt.Println("  Output limits: -1.0 to 1.0 (motor power)")
	fmt.Println("  Max motor speed: 1200 RPM at full power")
	fmt.Printf("  Motor time constant: %.2f s\n", motorTimeConstant)
	fmt.Println()

	// Test: Ramp up to target speed
	fmt.Println("Test: Speed Ramp to 1000 RPM")
	fmt.Printf("\n%-6s %-13s %-14s %-10s %-8s\n", "Step", "Target (RPM)", "Current (RPM)", "Error", "Power")
	fmt.Printf("%-6s %-13s %-14s %-10s %-8s\n", "----", "------------", "-------------", "-----", "-----")

	targetSpeed := 1000.0
	dt := 0.1
	plant := &motorPlant{}

	runScenario(controller, plant, targetSpeed, 32, dt)

	// Test: Speed change
	controller.Reset()

	fmt.Println("\nTest: Speed Change (1000 → 500 RPM)")
	fmt.Printf("\n%-6s %-13s %-14s %-10s %-8s\n", "Step", "Target (RPM)", "Current (RPM)", "Error", "Power")
	fmt.Printf("%-6s %-13s %-14s %-10s %-8s\n", "----", "------------", "-------------", "-----", "-----")

	targetSpeed = 500.0
	runScenario(controller, plant, targetSpeed, 32, dt)

	fmt.Println("\nAnalysis:")
	fmt.Println("---------")
	fmt.Println()
	fmt.Println("Acceleration phase (0 → 1000 RPM):")
	fmt.Println("  • High initial power due to large error")
	fmt.Println("  • Feed-forward supplies most of the nominal drive needed for the target speed")
	fmt.Println("  • The PID terms trim the remaining error, so the trace converges faster")
	fmt.Println("  • The trace shows a stable rise toward the target without the sign-flipping behavior from the old scripted example")
	fmt.Println()
	fmt.Println("Deceleration phase (1000 → 500 RPM):")
	fmt.Println("  • The controller initially commands braking to shed speed quickly")
	fmt.Println("  • The lower feed-forward bias helps the controller settle near the new operating point with smaller steady-state error")
	fmt.Println("  • This is now a true closed-loop example: controller output drives the plant state")
	fmt.Println()
	fmt.Println("Key Points:")
	fmt.Println("  • P term: Immediate response to error")
	fmt.Println("  • I term: Eliminates steady-state error")
	fmt.Println("  • D term: Reduces overshoot during changes")
	fmt.Println("  • Output limits prevent motor saturation")
	fmt.Println("  • The plant model is intentionally simple, but it produces realistic closed-loop behavior")
}
