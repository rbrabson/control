package main

import (
	"control/feedback"
	"control/motionprofile"
	"fmt"
	"math"
)

// SimulatedSystem represents a simple mass system for demonstration
type SimulatedSystem struct {
	mass     float64 // System mass (kg)
	position float64 // Current position (m)
	velocity float64 // Current velocity (m/s)
	damping  float64 // Damping coefficient
}

// NewSimulatedSystem creates a new simulated system
func NewSimulatedSystem(mass, damping float64) *SimulatedSystem {
	return &SimulatedSystem{
		mass:     mass,
		position: 0.0,
		velocity: 0.0,
		damping:  damping,
	}
}

// Update simulates the system dynamics given a control force
func (sys *SimulatedSystem) Update(force, dt float64) {
	// Simple physics: F = ma, with damping
	acceleration := (force - sys.damping*sys.velocity) / sys.mass

	// Euler integration
	sys.velocity += acceleration * dt
	sys.position += sys.velocity * dt
}

// GetState returns the current system state
func (sys *SimulatedSystem) GetState() feedback.Values {
	return feedback.Values{sys.position, sys.velocity}
}

// SetState sets the system state (for initialization)
func (sys *SimulatedSystem) SetState(position, velocity float64) {
	sys.position = position
	sys.velocity = velocity
}

func main() {
	fmt.Println("Motion Profile with Full-State Feedback Control Example")
	fmt.Println("======================================================")
	fmt.Println()

	// Create motion profile constraints
	constraints := motionprofile.Constraints{
		MaxVelocity:     2.0, // 2 m/s
		MaxAcceleration: 1.0, // 1 m/s²
	}

	// Define initial and goal states
	initial := motionprofile.State{
		Position: 0.0,
		Velocity: 0.0,
	}

	goal := motionprofile.State{
		Position: 5.0, // Move 5 meters
		Velocity: 0.0, // Come to rest
	}

	// Create motion profile
	profile := motionprofile.New(constraints, initial, goal)
	totalTime := profile.TotalTime()

	fmt.Printf("Motion Profile Generated:\n")
	fmt.Printf("- Distance: %.1f m\n", goal.Position-initial.Position)
	fmt.Printf("- Max Velocity: %.1f m/s\n", constraints.MaxVelocity)
	fmt.Printf("- Max Acceleration: %.1f m/s²\n", constraints.MaxAcceleration)
	fmt.Printf("- Total Time: %.2f s\n", totalTime)
	fmt.Println()

	// Create full-state feedback controller
	// Gains for [position_error, velocity_error] -> force
	// Kp = 50 N/m (position gain)
	// Kd = 20 N·s/m (velocity/damping gain)
	gains := feedback.Values{50.0, 20.0}
	controller := feedback.New(gains)

	// Create simulated system (1 kg mass, 2.0 N·s/m damping)
	system := NewSimulatedSystem(1.0, 2.0)

	// Simulation parameters
	dt := 0.01                        // 10ms timestep
	simulationTime := totalTime + 1.0 // Run a bit longer to see settling

	fmt.Printf("Full-State Feedback Control Simulation:\n")
	fmt.Printf("Controller Gains: Kp=%.0f N/m, Kd=%.0f N·s/m\n", gains[0], gains[1])
	fmt.Printf("System: Mass=%.1f kg, Damping=%.1f N·s/m\n", system.mass, system.damping)
	fmt.Println()

	// Tracking arrays for analysis
	var times []float64
	var positions []float64
	var velocities []float64
	var references []float64
	var errors []float64

	fmt.Printf("Time(s)  Ref Pos  Act Pos  Ref Vel  Act Vel  Error(m)  Force(N)\n")
	fmt.Printf("------  -------  -------  -------  -------  --------  --------\n")

	// Simulation loop
	for t := 0.0; t <= simulationTime; t += dt {
		// Get reference state from motion profile
		refState := profile.Calculate(t)
		reference := feedback.Values{refState.Position, refState.Velocity}

		// Get current system state
		measurement := system.GetState()

		// Calculate control output using full-state feedback
		force, err := controller.Calculate(reference, measurement)
		if err != nil {
			fmt.Printf("Controller error: %v\n", err)
			break
		}

		// Update system dynamics
		system.Update(force, dt)

		// Calculate tracking error
		positionError := refState.Position - measurement[0]

		// Store data for analysis
		times = append(times, t)
		positions = append(positions, measurement[0])
		velocities = append(velocities, measurement[1])
		references = append(references, refState.Position)
		errors = append(errors, math.Abs(positionError))

		// Print progress every 0.1 seconds
		if math.Mod(t, 0.1) < dt {
			fmt.Printf("%6.2f  %7.2f  %7.2f  %7.2f  %7.2f  %8.3f  %8.1f\n",
				t, refState.Position, measurement[0],
				refState.Velocity, measurement[1],
				positionError, force)
		}
	}

	fmt.Println()

	// Performance analysis
	fmt.Println("Performance Analysis:")
	fmt.Println("====================")

	// Find maximum tracking error
	maxError := 0.0
	for _, err := range errors {
		if err > maxError {
			maxError = err
		}
	}

	// Calculate final settling error
	finalError := math.Abs(goal.Position - system.position)

	// Calculate RMS error over motion profile time
	rmsError := 0.0
	profileSamples := int(totalTime / dt)
	for i := 0; i < profileSamples && i < len(errors); i++ {
		rmsError += errors[i] * errors[i]
	}
	rmsError = math.Sqrt(rmsError / float64(profileSamples))

	fmt.Printf("Maximum Tracking Error: %.3f m (%.1f%%)\n", maxError, 100*maxError/goal.Position)
	fmt.Printf("RMS Tracking Error: %.3f m (%.1f%%)\n", rmsError, 100*rmsError/goal.Position)
	fmt.Printf("Final Settling Error: %.3f m (%.2f%%)\n", finalError, 100*finalError/goal.Position)

	// Check if motion completed successfully
	if finalError < 0.01 { // Within 1cm
		fmt.Printf("✅ Motion completed successfully!\n")
	} else {
		fmt.Printf("⚠️  Large settling error detected\n")
	}

	fmt.Println()

	// Profile characteristics
	fmt.Println("Motion Profile Characteristics:")
	fmt.Println("==============================")
	fmt.Printf("Total Profile Time: %.2f s\n", totalTime)
	fmt.Println("Note: Detailed phase timing not available in current API")
	fmt.Println("Profile follows trapezoidal velocity shape with acceleration, cruise, and deceleration phases")

	// Control system insights
	fmt.Println()
	fmt.Println("Full-State Feedback Benefits:")
	fmt.Println("============================")
	fmt.Println("• Simultaneous position and velocity control")
	fmt.Println("• Predictive response to reference velocity")
	fmt.Println("• Optimal tracking of smooth motion profiles")
	fmt.Println("• Natural damping through velocity feedback")
	fmt.Println("• Single controller for complete trajectory following")

	fmt.Println()
	fmt.Println("Tuning Guidelines:")
	fmt.Println("=================")
	fmt.Printf("• Position Gain (Kp): %.0f N/m - Controls position tracking stiffness\n", gains[0])
	fmt.Printf("• Velocity Gain (Kd): %.0f N·s/m - Provides damping and velocity tracking\n", gains[1])
	fmt.Println("• Higher Kp: Faster position response, may cause overshoot")
	fmt.Println("• Higher Kd: More damping, smoother response, may slow response")
	fmt.Println("• Balance both gains for optimal tracking vs. stability")

}
