package pid

import (
	"math"
	"time"
)

// Option is a function type for configuring PID controller options
type Option func(*PID)

// PID represents a PID controller with proportional, integral, and derivative gains
type PID struct {
	// PID gains
	kp float64 // Proportional gain
	ki float64 // Integral gain
	kd float64 // Derivative gain

	// Control limits
	outputMin float64 // Minimum output value
	outputMax float64 // Maximum output value

	// Options
	feedForward              float64 // Feed-forward value added to PID output
	integralResetOnZeroCross bool    // Reset integral when error crosses zero
	stabilityThreshold       float64 // Derivative threshold to disable integral calculation
	integralSumMax           float64 // Maximum absolute value of integral sum
	derivativeFilterAlpha    float64 // Low-pass filter alpha for derivative (0-1, 0=no filter, 1=max filter)

	// Internal state
	integral           float64   // Accumulated integral term
	prevError          float64   // Previous error for derivative calculation
	filteredDerivative float64   // Filtered derivative for low-pass filtering
	prevTime           time.Time // Previous update time
	initialized        bool      // Flag to track first update
}

// New creates a new PID controller with the specified gains and optional configurations
func New(kp, ki, kd float64, opts ...Option) *PID {
	pid := &PID{
		kp:          kp,
		ki:          ki,
		kd:          kd,
		outputMin:   -math.Inf(1),
		outputMax:   math.Inf(1),
		initialized: false,
		// Set default values for optional features (disabled by default)
		integralResetOnZeroCross: false,
		stabilityThreshold:       0, // 0 means disabled
		integralSumMax:           0, // 0 means no limit
		derivativeFilterAlpha:    0, // 0 means no filtering
	}

	// Apply options
	for _, opt := range opts {
		opt(pid)
	}

	return pid
}

// WithFeedForward returns an option function that adds a feed-forward value to the PID output
func WithFeedForward(feedForward float64) Option {
	return func(p *PID) {
		p.feedForward = feedForward
	}
}

// WithIntegralResetOnZeroCross enables resetting the integral when error crosses zero
func WithIntegralResetOnZeroCross() Option {
	return func(p *PID) {
		p.integralResetOnZeroCross = true
	}
}

// WithStabilityThreshold sets a derivative threshold above which integral calculation is disabled
func WithStabilityThreshold(threshold float64) Option {
	return func(p *PID) {
		p.stabilityThreshold = math.Abs(threshold)
	}
}

// WithIntegralSumMax sets the maximum absolute value of the integral sum to prevent windup
func WithIntegralSumMax(maxSum float64) Option {
	return func(p *PID) {
		p.integralSumMax = math.Abs(maxSum)
	}
}

// WithDerivativeFilter applies a low-pass filter to the derivative term
// alpha should be between 0 (no filtering) and 1 (maximum filtering)
func WithDerivativeFilter(alpha float64) Option {
	return func(p *PID) {
		if alpha < 0 {
			alpha = 0
		} else if alpha > 1 {
			alpha = 1
		}
		p.derivativeFilterAlpha = alpha
	}
}

// SetOutputLimits sets the minimum and maximum output values
func (p *PID) SetOutputLimits(min, max float64) {
	if min > max {
		return
	}
	p.outputMin = min
	p.outputMax = max

	// Clamp integral to prevent windup
	p.integral = p.clamp(p.integral)
}

// SetGains updates the PID gains
func (p *PID) SetGains(kp, ki, kd float64) {
	p.kp = kp
	p.ki = ki
	p.kd = kd
}

// Update computes the PID output given the current error (setpoint - measurement)
func (p *PID) Update(error float64) float64 {
	now := time.Now()

	// Initialize on first call
	if !p.initialized {
		p.prevError = error
		p.prevTime = now
		p.initialized = true
		return 0
	}

	// Calculate time delta
	dt := now.Sub(p.prevTime).Seconds()
	if dt <= 0 {
		proportional := p.calculateProportional(error)
		return p.clamp(proportional)
	}

	// Calculate PID terms
	proportional := p.calculateProportional(error)
	derivative := p.calcualteDerrivative(error, dt)
	integral := p.calcualteIntegral(error, dt)

	// Calculate output
	output := proportional + integral + derivative + p.feedForward

	// Clamp output and handle integral windup
	clampedOutput := p.clamp(output)

	// Anti-windup: adjust integral if output is clamped
	if output != clampedOutput && p.ki != 0 {
		p.integral = (clampedOutput - proportional - derivative - p.feedForward) / p.ki
	}

	// Store values for next iteration
	p.prevError = error
	p.prevTime = now

	return clampedOutput
}

// calculateProportional computes the proportional term for a given error
func (p *PID) calculateProportional(error float64) float64 {
	proportional := p.kp * error
	return proportional
}

// calcualteDerrivative computes the derivative term for a given error and time delta
func (p *PID) calcualteDerrivative(error, dt float64) float64 {
	rawDerivative := p.calculateRawDerivative(error, dt)

	// Apply low-pass filter to derivative if enabled
	var derivative float64
	if p.derivativeFilterAlpha > 0 {
		p.filteredDerivative = p.derivativeFilterAlpha*p.filteredDerivative + (1-p.derivativeFilterAlpha)*rawDerivative
		derivative = p.kd * p.filteredDerivative
	} else {
		derivative = p.kd * rawDerivative
	}
	return derivative
}

// calcualteIntegral computes the integral term for a given error and time delta
func (p *PID) calcualteIntegral(error, dt float64) float64 {
	// Check for zero crossover and reset integral if enabled
	if p.integralResetOnZeroCross && ((p.prevError > 0 && error < 0) || (p.prevError < 0 && error > 0)) {
		p.integral = 0
	}

	// Integral term with stability threshold check
	rawDerivative := p.calculateRawDerivative(error, dt)
	var integral float64
	if p.stabilityThreshold == 0 || math.Abs(rawDerivative) <= p.stabilityThreshold {
		p.integral += error * dt

		// Cap integral sum if enabled
		if p.integralSumMax > 0 {
			if p.integral > p.integralSumMax {
				p.integral = p.integralSumMax
			} else if p.integral < -p.integralSumMax {
				p.integral = -p.integralSumMax
			}
		}

		integral = p.ki * p.integral
	} else {
		integral = p.ki * p.integral // Don't accumulate integral when above stability threshold
	}
	return integral
}

// calculateRawDerivative computes the raw derivative term for a given error and time delta
func (p *PID) calculateRawDerivative(error, dt float64) float64 {
	rawDerivative := (error - p.prevError) / dt
	return rawDerivative
}

// Reset clears the internal state of the PID controller
func (p *PID) Reset() {
	p.integral = 0
	p.prevError = 0
	p.filteredDerivative = 0
	p.initialized = false
}

// GetGains returns the current PID gains
func (p *PID) GetGains() (kp, ki, kd float64) {
	return p.kp, p.ki, p.kd
}

// GetIntegral returns the current integral value
func (p *PID) GetIntegral() float64 {
	return p.integral
}

// SetFeedForward sets the feed-forward value
func (p *PID) SetFeedForward(feedForward float64) {
	p.feedForward = feedForward
}

// GetFeedForward returns the current feed-forward value
func (p *PID) GetFeedForward() float64 {
	return p.feedForward
}

// SetIntegralResetOnZeroCross enables or disables integral reset on zero crossover
func (p *PID) SetIntegralResetOnZeroCross(enabled bool) {
	p.integralResetOnZeroCross = enabled
}

// GetIntegralResetOnZeroCross returns whether integral reset on zero crossover is enabled
func (p *PID) GetIntegralResetOnZeroCross() bool {
	return p.integralResetOnZeroCross
}

// SetStabilityThreshold sets the derivative threshold for disabling integral calculation
func (p *PID) SetStabilityThreshold(threshold float64) {
	p.stabilityThreshold = math.Abs(threshold)
}

// GetStabilityThreshold returns the current stability threshold
func (p *PID) GetStabilityThreshold() float64 {
	return p.stabilityThreshold
}

// SetIntegralSumMax sets the maximum absolute value of the integral sum
func (p *PID) SetIntegralSumMax(maxSum float64) {
	p.integralSumMax = math.Abs(maxSum)
}

// GetIntegralSumMax returns the current integral sum maximum
func (p *PID) GetIntegralSumMax() float64 {
	return p.integralSumMax
}

// SetDerivativeFilter sets the low-pass filter alpha for the derivative term
func (p *PID) SetDerivativeFilter(alpha float64) {
	if alpha < 0 {
		alpha = 0
	} else if alpha > 1 {
		alpha = 1
	}
	p.derivativeFilterAlpha = alpha
}

// GetDerivativeFilter returns the current derivative filter alpha value
func (p *PID) GetDerivativeFilter() float64 {
	return p.derivativeFilterAlpha
}

// clamp restricts the value to the output limits
func (p *PID) clamp(value float64) float64 {
	if value > p.outputMax {
		return p.outputMax
	}
	if value < p.outputMin {
		return p.outputMin
	}
	return value
}
