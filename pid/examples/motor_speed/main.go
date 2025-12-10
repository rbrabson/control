// Motor Speed Control Example
//
// This example demonstrates using a PID controller for motor speed control,
// which is common in robotics applications like robots. It shows:
// - Advanced PID features for motor control
// - Integral sum limiting to prevent motor saturation
// - Stability threshold to reduce overshoot
// - Derivative filtering for encoder noise
// - Realistic motor dynamics simulation

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"control/pid"
)

// Motor represents a simulated DC motor with realistic dynamics
type Motor struct {
	speed        float64 // Current speed (RPM)
	acceleration float64 // Current acceleration
	maxSpeed     float64 // Maximum motor speed
	timeConstant float64 // Motor time constant
	lastUpdate   time.Time
	noise        float64 // Encoder noise level
}

// NewMotor creates a new simulated motor
func NewMotor(maxSpeed, timeConstant, noiseLevel float64) *Motor {
	return &Motor{
		speed:        0.0,
		acceleration: 0.0,
		maxSpeed:     maxSpeed,
		timeConstant: timeConstant,
		lastUpdate:   time.Now(),
		noise:        noiseLevel,
	}
}

// ApplyPower applies power to the motor (-1.0 to 1.0)
func (m *Motor) ApplyPower(power float64) {
	now := time.Now()
	dt := now.Sub(m.lastUpdate).Seconds()
	m.lastUpdate = now

	// Clamp power to motor limits
	if power > 1.0 {
		power = 1.0
	} else if power < -1.0 {
		power = -1.0
	}

	if dt > 0 {
		// Target speed based on power
		targetSpeed := power * m.maxSpeed

		// First-order motor dynamics
		m.speed += (targetSpeed - m.speed) * dt / m.timeConstant

		// Add some acceleration dynamics for realism
		m.acceleration = (targetSpeed - m.speed) / m.timeConstant
	}
}

// GetSpeed returns current motor speed with encoder noise
func (m *Motor) GetSpeed() float64 {
	// Add realistic encoder noise
	noise := (rand.Float64() - 0.5) * m.noise * 2.0
	return m.speed + noise
}

// GetActualSpeed returns actual speed without noise (for comparison)
func (m *Motor) GetActualSpeed() float64 {
	return m.speed
}

func main() {
	fmt.Println("Motor Speed Control Example")
	fmt.Println("===========================")
	fmt.Println()

	// Create motor controller with advanced features for motor control
	controller := pid.New(0.8, 0.1, 0.02,
		pid.WithIntegralSumMax(1.0/0.1), // Ensure Ki * integralMax ≤ 1.0 for motor limits
		pid.WithStabilityThreshold(50),  // Disable integral during rapid speed changes
		pid.WithLowPassFilter(0.1),      // Filter encoder noise (10% filter)
	)

	// Set motor power limits (-1.0 to 1.0)
	controller.SetOutputLimits(-1.0, 1.0)

	// Create simulated motor (3000 RPM max, 0.2s time constant, 5 RPM noise)
	motor := NewMotor(3000.0, 0.2, 5.0)

	// Test multiple setpoints to demonstrate performance
	setpoints := []float64{1000, 2000, 1500, 500, -1000, 0}
	duration := 3.0   // 3 seconds per setpoint
	updateRate := 100 // 100Hz update rate (typical for motor control)
	interval := time.Duration(1000/updateRate) * time.Millisecond

	fmt.Printf("Update Rate: %d Hz\n", updateRate)
	fmt.Printf("Duration per setpoint: %.1f seconds\n", duration)
	fmt.Println()
	fmt.Println("Time\tSetpoint\tMeasured\tActual\t\tError\tPower")
	fmt.Println("----\t--------\t--------\t------\t\t-----\t-----")

	for _, setpoint := range setpoints {
		fmt.Printf("\nNew setpoint: %.0f RPM\n", setpoint)

		startTime := time.Now()

		for time.Since(startTime).Seconds() < duration {
			// Get current speed with noise (realistic encoder reading)
			measuredSpeed := motor.GetSpeed()
			actualSpeed := motor.GetActualSpeed()

			// Calculate PID output
			power := controller.Calculate(setpoint, measuredSpeed)

			// Apply power to motor
			motor.ApplyPower(power)

			// Print status every 0.1 seconds
			elapsed := time.Since(startTime).Seconds()
			if math.Mod(elapsed, 0.1) < float64(interval.Seconds()) {
				error := setpoint - measuredSpeed
				fmt.Printf("%.2f\t%.0f\t\t%.1f\t\t%.1f\t\t%.1f\t%.3f\n",
					elapsed, setpoint, measuredSpeed, actualSpeed, error, power)
			}

			time.Sleep(interval)
		}
	}

	fmt.Println()
	fmt.Println("Motor Control Completed")

	// Display final controller state
	kp, ki, kd := controller.GetGains()
	fmt.Printf("Final controller gains - Kp: %.1f, Ki: %.1f, Kd: %.3f\n", kp, ki, kd)
	fmt.Printf("Final integral value: %.3f\n", controller.GetIntegral())
	fmt.Printf("Integral sum max: %.1f\n", controller.GetIntegralSumMax())
	fmt.Printf("Stability threshold: %.0f\n", controller.GetStabilityThreshold())

	// Demonstrate runtime tuning
	fmt.Println("\nDemonstrating runtime gain tuning...")
	controller.SetGains(1.2, 0.15, 0.03)
	fmt.Println("Updated gains for more aggressive response")

	kp, ki, kd = controller.GetGains()
	fmt.Printf("New gains - Kp: %.1f, Ki: %.2f, Kd: %.3f\n", kp, ki, kd)
}
