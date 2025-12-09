package motionprofile

import (
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		constraints Constraints
		initial     State
		goal        State
	}{
		{
			name: "basic forward motion",
			constraints: Constraints{
				MaxVelocity:     2.0,
				MaxAcceleration: 1.0,
			},
			initial: State{Position: 0, Velocity: 0},
			goal:    State{Position: 10, Velocity: 0},
		},
		{
			name: "backward motion",
			constraints: Constraints{
				MaxVelocity:     1.5,
				MaxAcceleration: 0.8,
			},
			initial: State{Position: 5, Velocity: 0},
			goal:    State{Position: -3, Velocity: 0},
		},
		{
			name: "with initial velocity",
			constraints: Constraints{
				MaxVelocity:     3.0,
				MaxAcceleration: 2.0,
			},
			initial: State{Position: 0, Velocity: 1.0},
			goal:    State{Position: 8, Velocity: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := New(tt.constraints, tt.initial, tt.goal)

			// Check that profile is created
			if profile == nil {
				t.Fatal("Profile should not be nil")
			}

			// Check that total time is positive
			if profile.TotalTime() <= 0 {
				t.Errorf("Total time should be positive, got %f", profile.TotalTime())
			}
		})
	}
}

func TestMotionProfileCalculate(t *testing.T) {
	constraints := Constraints{
		MaxVelocity:     2.0,
		MaxAcceleration: 1.0,
	}
	initial := State{Position: 0, Velocity: 0}
	goal := State{Position: 10, Velocity: 0}

	profile := New(constraints, initial, goal)

	// Test initial state
	initialState := profile.Calculate(0)
	if initialState.Position != initial.Position {
		t.Errorf("Initial position should be %f, got %f", initial.Position, initialState.Position)
	}
	if initialState.Velocity != initial.Velocity {
		t.Errorf("Initial velocity should be %f, got %f", initial.Velocity, initialState.Velocity)
	}

	// Test final state
	finalState := profile.Calculate(profile.TotalTime())
	tolerance := 1e-10
	if math.Abs(finalState.Position-goal.Position) > tolerance {
		t.Errorf("Final position should be %f, got %f", goal.Position, finalState.Position)
	}
	if math.Abs(finalState.Velocity-goal.Velocity) > tolerance {
		t.Errorf("Final velocity should be %f, got %f", goal.Velocity, finalState.Velocity)
	}
}

func TestMotionProfileTriangleProfile(t *testing.T) {
	// Create constraints that force a triangle profile (short distance)
	constraints := Constraints{
		MaxVelocity:     10.0, // High max velocity
		MaxAcceleration: 1.0,
	}
	initial := State{Position: 0, Velocity: 0}
	goal := State{Position: 2, Velocity: 0} // Short distance

	profile := New(constraints, initial, goal)

	// Verify final position is reached
	finalState := profile.Calculate(profile.TotalTime())
	tolerance := 1e-10
	if math.Abs(finalState.Position-goal.Position) > tolerance {
		t.Errorf("Final position should be %f, got %f", goal.Position, finalState.Position)
	}
}

func TestMotionProfileIsFinished(t *testing.T) {
	constraints := Constraints{
		MaxVelocity:     2.0,
		MaxAcceleration: 1.0,
	}
	initial := State{Position: 0, Velocity: 0}
	goal := State{Position: 5, Velocity: 0}

	profile := New(constraints, initial, goal)

	// Should not be finished at start
	if profile.IsFinished(0) {
		t.Error("Profile should not be finished at time 0")
	}

	// Should not be finished at middle
	if profile.IsFinished(profile.TotalTime() / 2) {
		t.Error("Profile should not be finished at middle time")
	}

	// Should be finished at end
	if !profile.IsFinished(profile.TotalTime()) {
		t.Error("Profile should be finished at total time")
	}

	// Should be finished after end
	if !profile.IsFinished(profile.TotalTime() + 1) {
		t.Error("Profile should be finished after total time")
	}
}

func TestMotionProfileTimeLeftUntil(t *testing.T) {
	constraints := Constraints{
		MaxVelocity:     2.0,
		MaxAcceleration: 1.0,
	}
	initial := State{Position: 0, Velocity: 0}
	goal := State{Position: 10, Velocity: 0}

	profile := New(constraints, initial, goal)

	// Time to reach start should be 0
	if profile.TimeLeftUntil(0) != 0 {
		t.Errorf("Time to reach start should be 0, got %f", profile.TimeLeftUntil(0))
	}

	// Time to reach goal should be total time
	tolerance := 1e-10
	if math.Abs(profile.TimeLeftUntil(10)-profile.TotalTime()) > tolerance {
		t.Errorf("Time to reach goal should be %f, got %f", profile.TotalTime(), profile.TimeLeftUntil(10))
	}

	// Time to reach middle should be less than total time
	timeToMiddle := profile.TimeLeftUntil(5)
	if timeToMiddle <= 0 || timeToMiddle >= profile.TotalTime() {
		t.Errorf("Time to middle should be between 0 and %f, got %f", profile.TotalTime(), timeToMiddle)
	}
}

func TestMotionProfileBackwardMotion(t *testing.T) {
	constraints := Constraints{
		MaxVelocity:     2.0,
		MaxAcceleration: 1.0,
	}
	initial := State{Position: 10, Velocity: 0}
	goal := State{Position: 0, Velocity: 0}

	profile := New(constraints, initial, goal)

	// Check that profile moves backward
	midTime := profile.TotalTime() / 2
	midState := profile.Calculate(midTime)

	if midState.Position >= initial.Position {
		t.Error("Profile should move backward")
	}

	if midState.Velocity >= 0 {
		t.Error("Velocity should be negative for backward motion")
	}

	// Verify final position
	finalState := profile.Calculate(profile.TotalTime())
	tolerance := 1e-10
	if math.Abs(finalState.Position-goal.Position) > tolerance {
		t.Errorf("Final position should be %f, got %f", goal.Position, finalState.Position)
	}
}

// Benchmark tests
func BenchmarkMotionProfileCalculate(b *testing.B) {
	constraints := Constraints{
		MaxVelocity:     2.0,
		MaxAcceleration: 1.0,
	}
	initial := State{Position: 0, Velocity: 0}
	goal := State{Position: 10, Velocity: 0}

	profile := New(constraints, initial, goal)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		profile.Calculate(float64(i%100) * 0.01) // Test various times
	}
}

func BenchmarkNew(b *testing.B) {
	constraints := Constraints{
		MaxVelocity:     2.0,
		MaxAcceleration: 1.0,
	}
	initial := State{Position: 0, Velocity: 0}
	goal := State{Position: 10, Velocity: 0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(constraints, initial, goal)
	}
}
