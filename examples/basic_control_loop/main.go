// Basic PID Control Loop Example
//
// This example demonstrates the fundamental usage of the PID controller
// with a simple simulated system. It shows how to:
// - Create a basic PID controller
// - Set output limits
// - Run a control loop with timing
// - Monitor system response

package main

import (
	"fmt"
	"math"
	"time"

	"control/pid"
)

// SimulatedSystem represents a simple first-order system for demonstration
type SimulatedSystem struct {
	position     float64 // Current position
	velocity     float64 // Current velocity
	timeConstant float64 // System time constant
	lastUpdate   time.Time
}

// NewSimulatedSystem creates a new simulated system
func NewSimulatedSystem(timeConstant float64) *SimulatedSystem {
	return &SimulatedSystem{
		position:     0.0,
		velocity:     0.0,
		timeConstant: timeConstant,
		lastUpdate:   time.Now(),
	}
}

// ApplyControl applies a control input and updates the system state
func (s *SimulatedSystem) ApplyControl(input float64) float64 {
	now := time.Now()
	dt := now.Sub(s.lastUpdate).Seconds()
	s.lastUpdate = now

	if dt > 0 {
		// First-order system dynamics: τ * dv/dt + v = input
		s.velocity += (input - s.velocity) * dt / s.timeConstant
		s.position += s.velocity * dt
	}

	return s.position
}

// GetPosition returns the current position
func (s *SimulatedSystem) GetPosition() float64 {
	return s.position
}

func main() {
	fmt.Println("Basic PID Control Loop Example")
	fmt.Println("==============================")
	fmt.Println()

	// Create PID controller with moderate gains and output limits
	controller := pid.New(2.0, 0.5, 0.1,
		pid.WithOutputLimits(-10.0, 10.0), // Set reasonable output limits during initialization
	)

	// Create a simulated system (time constant = 0.5 seconds)
	system := NewSimulatedSystem(0.5)

	// Control parameters
	setpoint := 5.0  // Desired position
	duration := 10.0 // Run for 10 seconds
	updateRate := 50 // 50Hz update rate
	interval := time.Duration(1000/updateRate) * time.Millisecond

	fmt.Printf("Setpoint: %.1f\n", setpoint)
	fmt.Printf("Update Rate: %d Hz\n", updateRate)
	fmt.Printf("Duration: %.1f seconds\n", duration)
	fmt.Println()
	fmt.Println("Time\tSetpoint\tPosition\tError\tOutput")
	fmt.Println("----\t--------\t--------\t-----\t------")

	startTime := time.Now()

	for time.Since(startTime).Seconds() < duration {
		// Get current position from system
		position := system.GetPosition()

		// Calculate error
		error := setpoint - position

		// Update PID controller
		output := controller.Update(error)

		// Apply control to system
		system.ApplyControl(output)

		// Print status every 0.2 seconds for readability
		elapsed := time.Since(startTime).Seconds()
		if math.Mod(elapsed, 0.2) < float64(interval.Seconds()) {
			fmt.Printf("%.2f\t%.2f\t\t%.3f\t\t%.3f\t%.3f\n",
				elapsed, setpoint, position, error, output)
		}

		time.Sleep(interval)
	}

	fmt.Println()
	fmt.Printf("Final position: %.3f (error: %.3f)\n",
		system.GetPosition(), setpoint-system.GetPosition())

	// Display controller gains
	kp, ki, kd := controller.GetGains()
	fmt.Printf("Controller gains - Kp: %.1f, Ki: %.1f, Kd: %.1f\n", kp, ki, kd)

	fmt.Printf("Final integral value: %.3f\n", controller.GetIntegral())
}
