package feedforward

import "math"

// Option represents a configuration option for the feedforward controller.
type Option func(*FeedForward)

type FeedForward struct {
	kS   float64 // Static gain
	kV   float64 // Velocity gain
	kA   float64 // Acceleration gain
	kG   float64 // Gravity constant, for elevator / lift systems this is used to directly counter gravity.
	kCos float64 // For rotating arm systems, this is multiplied by the Cos of the target angle to counter nonlinear effects.
}

// New creates a new FeedForward controller with the given gains and optional parameters.
func New(kS, kV, kA float64, opts ...Option) *FeedForward {
	ff := &FeedForward{
		kS:   kS,
		kV:   kV,
		kA:   kA,
		kG:   0.0,
		kCos: 0.0,
	}

	// Apply options
	for _, opt := range opts {
		opt(ff)
	}

	return ff
}

// WithGravityGain sets the gravity constant gain for elevator/lift systems.
func WithGravityGain(kG float64) Option {
	return func(ff *FeedForward) {
		ff.kG = kG
	}
}

// WithCosineGain sets the cosine gain for rotating arm systems to counter nonlinear effects.
func WithCosineGain(kCos float64) Option {
	return func(ff *FeedForward) {
		ff.kCos = kCos
	}
}

// Calculate computes the feedforward output based on the position, velocity, and acceleration.
func (ff *FeedForward) Calculate(position float64, velocity float64, acceleration float64) float64 {
	output := ff.kV*velocity + ff.kA*acceleration + ff.kG
	if ff.kCos != 0.0 {
		output += ff.kCos * math.Cos(position)
	}

	return output
}
