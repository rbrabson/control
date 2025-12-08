package main

import (
	"fmt"
	"math"
	"time"

	"control/feedforward"
)

// RoboticArmExample demonstrates feedforward control for a rotating robotic arm
// with cosine compensation to handle gravitational effects at different angles.
func main() {
	fmt.Println("=== Robotic Arm Feedforward Control Example ===")
	fmt.Println("Simulating robotic arm control with cosine compensation")
	fmt.Println()

	// Create feedforward controller with cosine compensation
	// The cosine gain compensates for varying gravitational torque
	ff := feedforward.New(
		0.02,                            // kS: static gain (friction)
		0.5,                             // kV: velocity gain
		0.03,                            // kA: acceleration gain
		feedforward.WithCosineGain(2.5), // Cosine compensation for gravity
	)

	// Simulate arm movement through different positions
	targetPositions := []struct {
		name  string
		angle float64 // radians
	}{
		{"Horizontal Right", 0.0},
		{"45° Up", math.Pi / 4},
		{"Vertical Up", math.Pi / 2},
		{"135° Up", 3 * math.Pi / 4},
		{"Horizontal Left", math.Pi},
	}

	fmt.Printf("%-8s %-15s %-12s %-12s %-12s %-12s %-12s\n",
		"Time", "Position", "Angle", "Velocity", "Accel", "Cos(θ)", "FF Output")
	fmt.Println("--------------------------------------------------------------------------------")

	timeStep := 0.15
	for i := 0; i < len(targetPositions)-1; i++ {
		startAngle := targetPositions[i].angle
		endAngle := targetPositions[i+1].angle

		// Handle angle wrapping
		angleDiff := endAngle - startAngle
		if angleDiff > math.Pi {
			angleDiff -= 2 * math.Pi
		} else if angleDiff < -math.Pi {
			angleDiff += 2 * math.Pi
		}

		moveTime := 2.0 // seconds
		fmt.Printf("Moving from %s to %s...\n",
			targetPositions[i].name, targetPositions[i+1].name)

		for t := 0.0; t <= moveTime; t += timeStep {
			// Smooth trapezoidal motion profile
			normalizedTime := t / moveTime
			var s, v, a float64

			accelTime := 0.3 // 30% of move time for acceleration
			decelTime := 0.7 // start deceleration at 70% of move time

			if normalizedTime <= accelTime {
				// Acceleration phase
				s = 0.5 * (normalizedTime / accelTime) * (normalizedTime / accelTime)
				v = normalizedTime / (accelTime * moveTime)
				a = 1.0 / (accelTime * moveTime * moveTime)
			} else if normalizedTime <= decelTime {
				// Constant velocity phase
				s = 0.5*accelTime + (normalizedTime-accelTime)/accelTime*0.5
				v = 1.0 / (accelTime * moveTime)
				a = 0.0
			} else {
				// Deceleration phase
				remaining := (1.0 - normalizedTime) / (1.0 - decelTime)
				s = 1.0 - 0.5*remaining*remaining*(1.0-decelTime)
				v = remaining / ((1.0 - decelTime) * moveTime)
				a = -1.0 / ((1.0 - decelTime) * moveTime * moveTime)
			}

			angle := startAngle + s*angleDiff
			angularVel := v * angleDiff
			angularAccel := a * angleDiff

			// Calculate feedforward output (includes cosine compensation)
			ffOutput := ff.Calculate(angle, angularVel, angularAccel)
			cosValue := math.Cos(angle)

			// Convert angle to degrees for display
			angleDeg := angle * 180.0 / math.Pi

			fmt.Printf("%-8.1f %-15s %-12.1f %-12.3f %-12.3f %-12.3f %-12.3f\n",
				t, fmt.Sprintf("%.1f°", angleDeg), angle, angularVel, angularAccel, cosValue, ffOutput)

			time.Sleep(80 * time.Millisecond)
		}
		fmt.Println()
	}

	fmt.Println("=== Analysis ===")
	fmt.Println("Cosine compensation provides:")
	fmt.Println("- Maximum torque at horizontal position (cos(0°) = 1.0)")
	fmt.Println("- Zero torque at vertical positions (cos(90°) = 0.0)")
	fmt.Println("- Smooth torque variation following cos(θ) profile")
	fmt.Println("- Compensates for arm weight at different angles")
}
