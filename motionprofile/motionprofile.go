package motionprofile

import (
	"math"
)

// State represents a motion profile state at a specific time
type State struct {
	Position     float64 // Position in user units
	Velocity     float64 // Velocity in user units per second
	Acceleration float64 // Acceleration in user units per second^2
	Time         float64 // Time since profile start in seconds
}

// Constraints defines the motion constraints for profile generation
type Constraints struct {
	MaxVelocity     float64 // Maximum velocity in user units per second
	MaxAcceleration float64 // Maximum acceleration in user units per second^2
}

// MotionProfile generates and tracks a trapezoidal motion profile
type MotionProfile struct {
	constraints Constraints
	initial     State
	goal        State

	// Profile timing parameters
	accelerationTime float64 // Time to reach max velocity
	cruiseTime       float64 // Time spent at max velocity
	decelerationTime float64 // Time to decelerate from max velocity
	totalTime        float64 // Total profile time

	// Profile characteristics
	cruiseVelocity       float64 // Actual cruise velocity (may be less than max)
	accelerationDistance float64 // Distance covered during acceleration
	cruiseDistance       float64 // Distance covered during cruise
	decelerationDistance float64 // Distance covered during deceleration
}

// New creates a new trapezoidal motion profile
func New(constraints Constraints, initial, goal State) *MotionProfile {
	profile := &MotionProfile{
		constraints: constraints,
		initial:     initial,
		goal:        goal,
	}

	profile.calculateProfile()
	return profile
}

// calculateProfile computes the trapezoidal profile parameters
func (mp *MotionProfile) calculateProfile() {
	// Calculate the total displacement
	displacement := mp.goal.Position - mp.initial.Position

	// Handle zero displacement case
	if math.Abs(displacement) < 1e-10 {
		mp.totalTime = math.Abs(mp.goal.Velocity-mp.initial.Velocity) / mp.constraints.MaxAcceleration
		mp.accelerationTime = mp.totalTime
		mp.cruiseTime = 0
		mp.decelerationTime = 0
		mp.cruiseVelocity = mp.goal.Velocity
		return
	}

	// Determine direction and work in that direction
	direction := 1.0
	if displacement < 0 {
		direction = -1.0
	}

	// Target velocity in the direction of motion
	maxVel := mp.constraints.MaxVelocity * direction

	// Velocity constraints at start and end
	vStart := mp.initial.Velocity
	vEnd := mp.goal.Velocity

	// Calculate if we can reach max velocity with a trapezoidal profile
	// Distance needed to accelerate from vStart to maxVel
	accelDist := (maxVel*maxVel - vStart*vStart) / (2 * mp.constraints.MaxAcceleration * direction)
	// Distance needed to decelerate from maxVel to vEnd
	decelDist := (vEnd*vEnd - maxVel*maxVel) / (-2 * mp.constraints.MaxAcceleration * direction)

	if direction*(accelDist+decelDist) <= direction*displacement {
		// Trapezoidal profile - we can reach max velocity
		mp.accelerationTime = (maxVel - vStart) / (mp.constraints.MaxAcceleration * direction)
		mp.decelerationTime = (vEnd - maxVel) / (-mp.constraints.MaxAcceleration * direction)
		mp.cruiseVelocity = maxVel
		mp.accelerationDistance = accelDist
		mp.decelerationDistance = decelDist
		mp.cruiseDistance = displacement - accelDist - decelDist
		mp.cruiseTime = mp.cruiseDistance / maxVel
	} else {
		// Triangle profile - cannot reach max velocity
		// Find peak velocity using quadratic formula
		// s = (v^2 - v0^2)/(2*a) + (vf^2 - v^2)/(2*(-a))
		// Solving for v: v = sqrt((v0^2 + vf^2 + 2*a*s) / 2)
		discriminant := vStart*vStart + vEnd*vEnd + 2*mp.constraints.MaxAcceleration*displacement*direction
		if discriminant < 0 {
			// This shouldn't happen with valid constraints, but handle it gracefully
			mp.cruiseVelocity = vStart
			mp.totalTime = 0
			return
		}

		peakVel := direction * math.Sqrt(discriminant/2)
		mp.cruiseVelocity = peakVel
		mp.accelerationTime = (peakVel - vStart) / (mp.constraints.MaxAcceleration * direction)
		mp.decelerationTime = (vEnd - peakVel) / (-mp.constraints.MaxAcceleration * direction)
		mp.cruiseTime = 0
		mp.accelerationDistance = (peakVel*peakVel - vStart*vStart) / (2 * mp.constraints.MaxAcceleration * direction)
		mp.decelerationDistance = (vEnd*vEnd - peakVel*peakVel) / (-2 * mp.constraints.MaxAcceleration * direction)
		mp.cruiseDistance = 0
	}

	mp.totalTime = mp.accelerationTime + mp.cruiseTime + mp.decelerationTime
}

// Calculate returns the motion profile state at the given time
func (mp *MotionProfile) Calculate(t float64) State {
	// Clamp time to profile bounds
	if t <= 0 {
		return mp.initial
	}
	if t >= mp.totalTime {
		return mp.goal
	}

	direction := 1.0
	if mp.goal.Position < mp.initial.Position {
		direction = -1.0
	}

	var state State
	state.Time = t

	switch {
	case t <= mp.accelerationTime:
		// Acceleration phase
		state.Acceleration = mp.constraints.MaxAcceleration * direction
		state.Velocity = mp.initial.Velocity + state.Acceleration*t
		state.Position = mp.initial.Position + mp.initial.Velocity*t + 0.5*state.Acceleration*t*t
	case t <= mp.accelerationTime+mp.cruiseTime:
		// Cruise phase
		state.Acceleration = 0
		state.Velocity = mp.cruiseVelocity
		cruiseT := t - mp.accelerationTime
		state.Position = mp.initial.Position + mp.accelerationDistance + mp.cruiseVelocity*cruiseT
	default:
		// Deceleration phase
		decelT := t - mp.accelerationTime - mp.cruiseTime
		state.Acceleration = -mp.constraints.MaxAcceleration * direction
		state.Velocity = mp.cruiseVelocity + state.Acceleration*decelT
		state.Position = mp.initial.Position + mp.accelerationDistance + mp.cruiseDistance +
			mp.cruiseVelocity*decelT + 0.5*state.Acceleration*decelT*decelT
	}

	return state
}

// IsFinished returns true if the profile has completed at the given time
func (mp *MotionProfile) IsFinished(t float64) bool {
	return t >= mp.totalTime
}

// TotalTime returns the total time for the motion profile
func (mp *MotionProfile) TotalTime() float64 {
	return mp.totalTime
}

// TimeLeftUntil returns the time remaining until the given position is reached
func (mp *MotionProfile) TimeLeftUntil(targetPosition float64) float64 {
	direction := 1.0
	if mp.goal.Position < mp.initial.Position {
		direction = -1.0
	}

	targetDistance := (targetPosition - mp.initial.Position)
	if direction*targetDistance <= 0 {
		return 0
	}

	totalDistance := mp.goal.Position - mp.initial.Position
	if direction*targetDistance >= direction*totalDistance {
		return mp.totalTime
	}

	// Check which phase the target is in
	switch {
	case direction*targetDistance <= direction*mp.accelerationDistance:
		// Target is in acceleration phase
		// s = v0*t + 0.5*a*t^2, solve for t
		a := 0.5 * mp.constraints.MaxAcceleration * direction
		b := mp.initial.Velocity
		c := -targetDistance
		discriminant := b*b - 4*a*c
		if discriminant < 0 {
			return 0
		}
		return (-b + math.Sqrt(discriminant)) / (2 * a)
	case direction*targetDistance <= direction*(mp.accelerationDistance+mp.cruiseDistance):
		// Target is in cruise phase
		cruiseDistance := targetDistance - mp.accelerationDistance
		return mp.accelerationTime + cruiseDistance/mp.cruiseVelocity
	default:
		// Target is in deceleration phase
		remainingDistance := totalDistance - targetDistance
		// Working backward from end using deceleration
		a := 0.5 * mp.constraints.MaxAcceleration * direction
		b := mp.goal.Velocity
		c := remainingDistance
		discriminant := b*b - 4*a*c
		if discriminant < 0 {
			return mp.totalTime
		}
		timeFromEnd := (-b + math.Sqrt(discriminant)) / (2 * a)
		return mp.totalTime - timeFromEnd
	}
}
