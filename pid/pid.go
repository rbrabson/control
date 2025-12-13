package pid

import (
	"control/filter"
	"log/slog"
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

	// Options
	feedForward              float64       // Feed-forward value added to PID output
	integralResetOnZeroCross bool          // Reset integral when error crosses zero
	stabilityThreshold       float64       // Derivative threshold to disable integral calculation
	integralSumMax           float64       // Maximum absolute value of integral sum
	outputMin                float64       // Minimum output value
	outputMax                float64       // Maximum output value
	filter                   filter.Filter // Filter for derivative term

	// Internal state
	integral      float64   // Accumulated integral term
	lastReference float64   // Previous reference for derivative calculation
	lastError     float64   // Previous error for derivative calculation
	prevTime      time.Time // Previous update time
	initialized   bool      // Flag to track first update
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
		filter:                   nil,        // No derivative filter by default
		stabilityThreshold:       math.NaN(), // No stability threshold by default
		integralSumMax:           math.NaN(), // No integral sum cap by default
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

// WithFilter sets a filter for the derivative term. Examples are a low pass filter or a kalman filter.
func WithFilter(f filter.Filter) Option {
	return func(p *PID) {
		p.filter = f
	}
}

// WithOutputLimits sets the minimum and maximum output limits
func WithOutputLimits(min, max float64) Option {
	return func(p *PID) {
		if min <= max {
			p.outputMin = min
			p.outputMax = max
		}
	}
}

// WithDampening configures the derivative gain (kd) based on desired dampening characteristics, with an
// optional percent overshoot (po). If po is 0, critical dampening is used.
func WithDampening(ka, kv, po float64) Option {
	return func(p *PID) {
		// If this inequality is true, kd will be knegative and there will be a scary non-minimum phase system.
		// There isn't a good way to return the error here, so we log the error and return.
		if p.kp < kv*kv/4*ka {
			slog.Error("invalid kp, kv, and ka values for PID.WithDampening",
				"kp", p.kp,
				"kv", kv,
				"ka", ka,
				"kv^2 / 4*ka", kv*kv/4*ka,
			)
			return
		}

		if po == 0 {
			// Critical dampening
			p.kd = 2*math.Sqrt(ka*kv) - ka
		} else {
			// Calculate damping ratio from percent overshoot
			po = math.Max(po/100, 0.01) // Prevent log(0)
			poLog := math.Log(po)
			zeta := -poLog / math.Sqrt(math.Pi*math.Pi+poLog*poLog)
			p.kd = 2*zeta*math.Sqrt(ka*kv) - kv
		}
	}
}

// Calculate computes the PID output for the given reference (setpoint) and current state (measurement)
func (p *PID) Calculate(reference, state float64) float64 {
	now := time.Now()
	error := reference - state

	// Initialize on first call
	if !p.initialized {
		p.integral = 0
		p.lastReference = reference
		p.lastError = error
		p.prevTime = now
		p.initialized = true
	}

	// Reset integral on setpoint change to prevent windup
	if reference != p.lastReference {
		p.integral = 0
		p.lastReference = reference
	}

	// Calculate time delta
	dt := now.Sub(p.prevTime).Seconds()

	// Calculate PID terms
	proportional := p.calculateProportional(error)
	derivative := p.calcualteDerrivative(error, dt)
	integral := p.calculateIntegral(error, dt)

	// Calculate output
	output := proportional + integral + derivative + p.feedForward

	// Clamp output and handle integral windup
	clampedOutput := p.clamp(output)

	// Anti-windup: adjust integral if output is clamped
	if output != clampedOutput && p.ki != 0 {
		p.integral = (clampedOutput - proportional - derivative - p.feedForward) / p.ki
	}

	// Store values for next iteration
	p.lastError = error
	p.prevTime = now

	return clampedOutput
}

// calculateProportional computes the proportional term for a given error
func (p *PID) calculateProportional(error float64) float64 {
	proportional := p.kp * error
	return proportional
}

// calculateIntegral computes the integral term for a given error and time delta
func (p *PID) calculateIntegral(error, dt float64) float64 {
	// Check for zero crossover and reset integral if enabled
	if p.integralResetOnZeroCross && ((p.lastError > 0 && error < 0) || (p.lastError < 0 && error > 0)) {
		p.integral = 0
	}

	// Integral term with stability threshold check
	rawDerivative := p.calculateRawDerivative(error, dt)
	if math.IsNaN(p.stabilityThreshold) || math.Abs(rawDerivative) <= p.stabilityThreshold {
		p.integral += error * dt

		// Cap integral sum if enabled
		if !math.IsNaN(p.integralSumMax) {
			if p.integral > p.integralSumMax {
				p.integral = p.integralSumMax
			} else if p.integral < -p.integralSumMax {
				p.integral = -p.integralSumMax
			}
		}
	}

	integral := p.ki * p.integral
	return integral
}

// calcualteDerrivative computes the derivative term for a given error and time delta
func (p *PID) calcualteDerrivative(error, dt float64) float64 {
	rawDerivative := p.calculateRawDerivative(error, dt)
	derivative := p.kd * rawDerivative

	return derivative
}

// calculateRawDerivative computes the raw derivative term for a given error and time delta
func (p *PID) calculateRawDerivative(error, dt float64) float64 {
	if dt <= 0 {
		return 0
	}

	errorChange := error - p.lastError
	var currentEstimate float64
	if p.filter != nil {
		// Apply derivative filter if enabled
		currentEstimate = p.filter.Estimate(errorChange)
	} else {
		// No derivative filter
		currentEstimate = errorChange
	}
	rawDerivative := currentEstimate / dt
	return rawDerivative
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

// Reset the initialized state of the PID controller. When the PID output is calculated
// the next time, the internal state will be reset as well.
func (p *PID) Reset() *PID {
	p.integral = 0
	p.initialized = false
	p.lastError = 0
	if p.filter != nil {
		p.filter.Reset()
	}
	return p
}

// SetGains updates the PID gains
func (p *PID) SetGains(kp, ki, kd float64) *PID {
	p.kp = kp
	p.ki = ki
	p.kd = kd

	return p
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
func (p *PID) SetFeedForward(feedForward float64) *PID {
	p.feedForward = feedForward
	return p
}

// GetFeedForward returns the current feed-forward value
func (p *PID) GetFeedForward() float64 {
	return p.feedForward
}

// SetIntegralResetOnZeroCross enables or disables integral reset on zero crossover
func (p *PID) SetIntegralResetOnZeroCross(enabled bool) *PID {
	p.integralResetOnZeroCross = enabled
	return p
}

// GetIntegralResetOnZeroCross returns whether integral reset on zero crossover is enabled
func (p *PID) GetIntegralResetOnZeroCross() bool {
	return p.integralResetOnZeroCross
}

// SetStabilityThreshold sets the derivative threshold for disabling integral calculation
func (p *PID) SetStabilityThreshold(threshold float64) *PID {
	p.stabilityThreshold = math.Abs(threshold)
	return p
}

// GetStabilityThreshold returns the current stability threshold
func (p *PID) GetStabilityThreshold() float64 {
	return p.stabilityThreshold
}

// SetIntegralSumMax sets the maximum absolute value of the integral sum
func (p *PID) SetIntegralSumMax(maxSum float64) *PID {
	p.integralSumMax = math.Abs(maxSum)
	return p
}

// GetIntegralSumMax returns the current integral sum maximum
func (p *PID) GetIntegralSumMax() float64 {
	return p.integralSumMax
}

// SetFilter sets a filter for the derivative term. Examples are a low pass filter or a kalman filter.
func (p *PID) SetFilter(f filter.Filter) *PID {
	p.filter = f
	return p
}

// GetFilter returns the current filter used for the derivative term, or nil if no filter is set.
func (p *PID) GetFilter() filter.Filter {
	return p.filter
}

// SetOutputLimits sets the minimum and maximum output values
func (p *PID) SetOutputLimits(min, max float64) *PID {
	if min > max {
		return p
	}
	p.outputMin = min
	p.outputMax = max

	// Clamp integral to prevent windup
	p.integral = p.clamp(p.integral)
	return p
}

// GetOutputLimits returns the current output limits
func (p *PID) GetOutputLimits() (min, max float64) {
	return p.outputMin, p.outputMax
}
