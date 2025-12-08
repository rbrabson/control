// Position Servo Control Example
//
// This example demonstrates using a PID controller for servo position control,
// common in robotics for precise positioning. It shows:
// - Position control with integral windup protection
// - Stability threshold for high-speed movement
// - Output limiting for servo constraints
// - Realistic servo dynamics with backlash and friction

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"control/pid"
)

// ServoSystem represents a simulated servo motor system
type ServoSystem struct {
	position     float64 // Current position (degrees)
	velocity     float64 // Current velocity (degrees/second)
	acceleration float64 // Current acceleration
	maxVelocity  float64 // Maximum servo velocity
	maxAccel     float64 // Maximum servo acceleration
	backlash     float64 // Mechanical backlash (degrees)
	friction     float64 // Friction coefficient
	lastUpdate   time.Time
	encoderRes   float64 // Encoder resolution (degrees per count)
}

// NewServoSystem creates a new simulated servo system
func NewServoSystem(maxVel, maxAccel, backlash, friction, encoderRes float64) *ServoSystem {
	return &ServoSystem{
		position:     0.0,
		velocity:     0.0,
		acceleration: 0.0,
		maxVelocity:  maxVel,
		maxAccel:     maxAccel,
		backlash:     backlash,
		friction:     friction,
		lastUpdate:   time.Now(),
		encoderRes:   encoderRes,
	}
}

// ApplyCommand applies a position command (-1.0 to 1.0 normalized)
func (s *ServoSystem) ApplyCommand(command float64) {
	now := time.Now()
	dt := now.Sub(s.lastUpdate).Seconds()
	s.lastUpdate = now

	// Clamp command to valid range
	if command > 1.0 {
		command = 1.0
	} else if command < -1.0 {
		command = -1.0
	}

	if dt > 0 {
		// Convert command to desired acceleration
		desiredAccel := command * s.maxAccel

		// Apply friction (opposing current motion)
		frictionForce := -math.Copysign(s.friction, s.velocity)

		// Net acceleration
		s.acceleration = desiredAccel + frictionForce

		// Update velocity with acceleration limits
		newVelocity := s.velocity + s.acceleration*dt
		if newVelocity > s.maxVelocity {
			newVelocity = s.maxVelocity
		} else if newVelocity < -s.maxVelocity {
			newVelocity = -s.maxVelocity
		}
		s.velocity = newVelocity

		// Update position
		s.position += s.velocity * dt

		// Simulate backlash (simple model)
		if math.Abs(s.velocity) < 0.1 { // Near zero velocity
			backlashError := (rand.Float64() - 0.5) * s.backlash
			s.position += backlashError * 0.1 // Small backlash effect
		}
	}
}

// GetPosition returns current position with encoder quantization
func (s *ServoSystem) GetPosition() float64 {
	// Quantize to encoder resolution
	counts := math.Round(s.position / s.encoderRes)
	return counts * s.encoderRes
}

// GetActualPosition returns actual position without quantization
func (s *ServoSystem) GetActualPosition() float64 {
	return s.position
}

// GetVelocity returns current velocity
func (s *ServoSystem) GetVelocity() float64 {
	return s.velocity
}

func main() {
	fmt.Println("Position Servo Control Example")
	fmt.Println("==============================")
	fmt.Println()

	// Create position controller with anti-windup and stability features
	controller := pid.New(1.5, 0.2, 0.05,
		pid.WithIntegralSumMax(10.0),     // Prevent integral windup
		pid.WithStabilityThreshold(50.0), // Disable integral during rapid movement (>50 deg/s)
	)

	// Set servo command limits (-1.0 to 1.0)
	controller.SetOutputLimits(-1.0, 1.0)

	// Create simulated servo system
	// (max vel: 180 deg/s, max accel: 360 deg/s², backlash: 0.1°, friction: 10, encoder: 0.1°)
	servo := NewServoSystem(180.0, 360.0, 0.1, 10.0, 0.1)

	// Test multiple position targets
	targets := []float64{45.0, 90.0, -30.0, 180.0, 0.0, -90.0}
	duration := 4.0   // 4 seconds per target
	updateRate := 100 // 100Hz update rate
	interval := time.Duration(1000/updateRate) * time.Millisecond

	fmt.Printf("Update Rate: %d Hz\n", updateRate)
	fmt.Printf("Duration per target: %.1f seconds\n", duration)
	fmt.Println()
	fmt.Println("Time\tTarget\tMeasured\tActual\t\tVelocity\tError\tCommand")
	fmt.Println("----\t------\t--------\t------\t\t--------\t-----\t-------")

	totalTime := 0.0

	for i, target := range targets {
		fmt.Printf("\nMoving to target %d: %.1f°\n", i+1, target)

		startTime := time.Now()

		for time.Since(startTime).Seconds() < duration {
			// Get current position (with encoder quantization)
			measuredPos := servo.GetPosition()
			actualPos := servo.GetActualPosition()
			velocity := servo.GetVelocity()

			// Calculate PID output
			command := controller.Calculate(target, measuredPos)

			// Apply command to servo
			servo.ApplyCommand(command)

			// Print status every 0.1 seconds
			elapsed := time.Since(startTime).Seconds()
			if math.Mod(elapsed, 0.1) < float64(interval.Seconds()) {
				error := target - measuredPos
				fmt.Printf("%.2f\t%.1f\t%.2f\t\t%.2f\t\t%.1f\t\t%.2f\t%.3f\n",
					totalTime+elapsed, target, measuredPos, actualPos, velocity, error, command)
			}

			time.Sleep(interval)
		}

		totalTime += duration
	}

	fmt.Println()
	fmt.Println("Position Control Completed")

	// Display final system state
	fmt.Printf("Final position: %.2f° (measured: %.2f°)\n",
		servo.GetActualPosition(), servo.GetPosition())
	fmt.Printf("Final velocity: %.1f°/s\n", servo.GetVelocity())

	// Display controller configuration
	kp, ki, kd := controller.GetGains()
	fmt.Printf("Controller gains - Kp: %.1f, Ki: %.1f, Kd: %.2f\n", kp, ki, kd)
	fmt.Printf("Integral sum max: %.1f\n", controller.GetIntegralSumMax())
	fmt.Printf("Stability threshold: %.1f°/s\n", controller.GetStabilityThreshold())
	fmt.Printf("Final integral: %.3f\n", controller.GetIntegral())

	// Demonstrate settling time analysis
	fmt.Println("\nPosition servo control demonstrates:")
	fmt.Println("- Precise positioning with encoder feedback")
	fmt.Println("- Integral windup protection for large errors")
	fmt.Println("- Stability threshold to prevent oscillation")
	fmt.Println("- Realistic servo dynamics with backlash and friction")

	// Show final accuracy
	finalError := math.Abs(targets[len(targets)-1] - servo.GetPosition())
	fmt.Printf("Final positioning accuracy: ±%.3f°\n", finalError)
}
