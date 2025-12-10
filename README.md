# Control Systems Library

A comprehensive control systems library in Go featuring PID controllers, feedback control, and feedforward control implementations for robotics and industrial automation.

## Features

### Core PID Functionality

- **Proportional, Integral, and Derivative** control with configurable gains
- **Time-based calculations** for accurate integral and derivative computation
- **Output clamping** with configurable minimum and maximum limits
- **Anti-windup protection** to prevent integral saturation

### Advanced Control Features

- **Feed-forward control** for predictive system response
- **Integral reset on zero crossover** to prevent wrong-direction movement
- **Stability threshold** to disable integral calculation during high-speed changes
- **Integral sum capping** for additional windup protection
- **Filter interface support** with LowPassFilter and KalmanFilter implementations for noise reduction
- **Combined dampening features** for enhanced stability in noisy environments
- **Gravity compensation** for vertical motion systems
- **Cosine compensation** for angular/rotating systems
- **Combined compensation** strategies for complex machinery

### Performance & Reliability

- **PID Controllers**: ~64ns per update - Excellent performance for real-time applications
- **Feedforward Controllers**: ~2-5ns per calculation - Ultra-fast predictive control
- **Feedback Controllers**: ~3-220ns depending on dimension, optimized for multi-variable systems
- **Filter Controllers**: Kalman ~44ns, LowPass ~4.7ns - Advanced signal filtering
- **Comprehensive test coverage** with extensive test suites
- **Thread-safe design** for concurrent applications
- **Robust error handling** for edge cases and invalid inputs

## Installation

```bash
go get github.com/rbrabson/control
```

## Packages

This library provides six main packages:

### PID Package (`control/pid`)

High-performance PID controller implementation with advanced features:

- Proportional, Integral, and Derivative control
- Feed-forward control and anti-windup protection
- Configurable output limits and pluggable filter interface
- Runtime parameter adjustment

### Feedback Package (`control/feedback`)

Full-state feedback control implementation:

- Multi-dimensional state feedback control
- Vector-based control calculations
- Error handling for dimension mismatches
- High-performance implementation

### Feedforward Package (`control/feedforward`)

Predictive feedforward controllers with advanced compensation:

- Basic velocity and acceleration feedforward
- Gravity compensation for vertical systems
- Cosine compensation for rotating machinery
- Combined compensation strategies
- Options pattern for flexible configuration
- Ultra-fast calculations (2-5 ns/op)

### Motion Profile Package (`control/motionprofile`)

Trapezoidal motion profile generation for smooth trajectory planning:

- Trapezoidal and triangle motion profiles
- Bidirectional motion support (forward/backward)
- Real-time state calculation (~50ns per calculation)
- Constraint validation and error handling
- Complete trajectory information (position, velocity, acceleration)

### InterpLUT Package (`control/interplut`)

Interpolating lookup tables with cubic Hermite spline interpolation:

- Smooth interpolation between control points using cubic splines
- Monotonicity preservation for stable control behavior
- High-performance lookups (~36ns per interpolation)
- Automatic sorting and error handling
- Perfect for robotics sensor calibration and non-linear mappings
- Adaptive PID control with dynamic coefficient lookup

Examples include basic shooter velocity mapping, non-linear temperature control, and adaptive PID control with dynamic coefficient lookup.

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    "control/pid"
)

func main() {
    // Create a basic PID controller with output limits
    controller := pid.New(1.0, 0.1, 0.05,
        pid.WithOutputLimits(-100, 100), // Set output limits during initialization
    )
    
    // Control loop
    setpoint := 50.0
    measurement := 0.0
    
    for i := 0; i < 100; i++ {
        output := controller.Calculate(setpoint, measurement)
        
        // Apply output to your system
        // measurement = yourSystem.ApplyControl(output)
        
        error := setpoint - measurement
        fmt.Printf("Error: %.2f, Output: %.2f\n", error, output)
        time.Sleep(10 * time.Millisecond)
    }
}
```

### Using Feedback Controllers

```go
package main

import (
    "fmt"
    "control/feedback"
)

func main() {
    // Full state feedback for [position, velocity] control
    gains := feedback.Values{1.5, 0.3}
    controller := feedback.New(gains)
    
    // Control loop
    setpoint := feedback.Values{10.0, 0.0}  // Target position=10, velocity=0
    
    for i := 0; i < 10; i++ {
        // Current state [position, velocity]
        current := feedback.Values{float64(i), 0.5}
        
        output, err := controller.Calculate(setpoint, current)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }
        
        fmt.Printf("State: %v, Output: %.3f\n", current, output)
    }
}
```

## Advanced Usage

### Using Options Pattern

```go
// Create a low-pass filter
filter, _ := filter.NewLowPassFilter(0.1)

// Create PID with advanced features
controller := pid.New(1.0, 0.1, 0.05,
    pid.WithFeedForward(5.0),                    // Add feed-forward term
    pid.WithIntegralResetOnZeroCross(),          // Reset integral at zero crossing
    pid.WithStabilityThreshold(2.0),             // Disable integral when derivative > 2.0
    pid.WithIntegralSumMax(10.0),                // Cap integral sum at ±10.0
    pid.WithFilter(filter),                      // Apply filter interface for noise reduction
)
```

### Runtime Configuration

```go
controller := pid.New(1.0, 0.1, 0.05)

// Update gains during runtime
controller.SetGains(1.5, 0.2, 0.08)

// Configure output limits
controller.SetOutputLimits(-50, 50)

// Modify advanced features
controller.SetFeedForward(3.0)
```

### Feedforward Controller Options

```go
// Basic feedforward
basicFF := feedforward.New()

// Elevator with gravity compensation  
elevatorFF := feedforward.New(feedforward.WithGravityGain(9.81))

// Robotic arm with cosine compensation
armFF := feedforward.New(feedforward.WithCosineGain(2.5))

// Crane with combined compensation
craneFF := feedforward.New(
    feedforward.WithGravityGain(15.7),
    feedforward.WithCosineGain(8.2),
)
controller.SetStabilityThreshold(1.5)
controller.SetIntegralSumMax(8.0)
controller.SetDerivativeFilter(0.2)
```

## API Reference

### Constructor

#### `New(kp, ki, kd float64, opts ...Option) *PID`

Creates a new PID controller with the specified gains and optional configurations.

**Parameters:**

- `kp`: Proportional gain
- `ki`: Integral gain  
- `kd`: Derivative gain
- `opts`: Optional configuration functions

### Core Methods

#### `Calculate(reference, state float64) float64`

Computes the PID output for the given reference (setpoint) and current state (measurement). The error is calculated internally as reference - state.

#### `Reset()`

Clears all internal state (integral, previous error, initialization flag).

#### `SetGains(kp, ki, kd float64)`

Updates the PID gains during runtime.

#### `GetGains() (kp, ki, kd float64)`

Returns the current PID gains.

### Output Control

```go
// Option function
pid.WithOutputLimits(min, max float64) // Set output limits during initialization

// Runtime method
controller.SetOutputLimits(min, max float64) // Set output limits after creation
```

Sets the minimum and maximum output values with anti-windup protection.

### Advanced Methods

#### Feed-Forward Control

```go
// Option function
pid.WithFeedForward(value float64)

// Runtime methods
controller.SetFeedForward(value float64)
controller.GetFeedForward() float64
```

#### Integral Reset on Zero Crossover

```go
// Option function
pid.WithIntegralResetOnZeroCross()

// Runtime methods
controller.SetIntegralResetOnZeroCross(enabled bool)
controller.GetIntegralResetOnZeroCross() bool
```

#### Stability Threshold

```go
// Option function
pid.WithStabilityThreshold(threshold float64)

// Runtime methods
controller.SetStabilityThreshold(threshold float64)
controller.GetStabilityThreshold() float64
```

#### Integral Sum Limiting

```go
// Option function
pid.WithIntegralSumMax(maxSum float64)

// Runtime methods
controller.SetIntegralSumMax(maxSum float64)
controller.GetIntegralSumMax() float64
```

#### Filter Interface

```go
// Create filter
filter, _ := filter.NewLowPassFilter(0.1) // 0-1 gain parameter

// Option function
pid.WithFilter(filter filter.Filter)

// Runtime methods
controller.SetFilter(filter filter.Filter)
controller.GetFilter() filter.Filter
```

## Feedback Package API

### FullStateFeedback

Full state feedback controller for multi-dimensional control systems.

#### FullStateFeedback Constructor

```go
func New(gain Values) *FullStateFeedback
```

Creates a new full state feedback controller with specified gain vector.

**Parameters:**

- `gain`: Vector of gain values for each state variable

#### Methods

```go
func (fsf *FullStateFeedback) Calculate(setpoint, measurement Values) (float64, error)
```

Computes control output based on full state feedback: u = K(r - x)

**Parameters:**

- `setpoint`: Reference state vector
- `measurement`: Current state vector

**Returns:**

- Control output value
- Error if vectors have different lengths

#### Usage Example

```go
import "control/feedback"

// Create gains for [position, velocity] state feedback
gains := feedback.Values{2.0, 0.5}
controller := feedback.New(gains)

// Control calculation
setpoint := feedback.Values{10.0, 0.0}  // Target: position=10, velocity=0
current := feedback.Values{8.5, 1.2}    // Current: position=8.5, velocity=1.2
output, err := controller.Calculate(setpoint, current)
if err != nil {
    // Handle error (e.g., mismatched vector lengths)
}
```

### Types

```go
type Values []float64  // Vector type for multi-dimensional states
```

## Use Cases

### Robotics Applications

```go
// Motor speed control for robots
speedController := pid.New(0.8, 0.1, 0.02,
    pid.WithIntegralSumMax(1.0/0.1),  // Ensure Ki * integralMax ≤ 1.0 for motor limits
    pid.WithStabilityThreshold(50),    // Reduce overshoot during rapid changes
    pid.WithDerivativeFilter(0.1),     // Filter encoder noise
)
speedController.SetOutputLimits(-1.0, 1.0) // Motor power limits
```

### Temperature Control

```go
// Temperature controller with feed-forward for ambient compensation
tempController := pid.New(2.0, 0.5, 0.1,
    pid.WithFeedForward(ambientCompensation),
    pid.WithIntegralResetOnZeroCross(),  // Prevent overshoot when crossing target
    pid.WithDerivativeFilter(0.2),       // Filter temperature sensor noise
)
```

### Position Control

```go
// Servo position control with anti-windup
positionController := pid.New(1.5, 0.2, 0.05,
    pid.WithIntegralSumMax(10.0),        // Prevent integral windup
    pid.WithStabilityThreshold(5.0),     // Disable integral during rapid movement
)
```

### Motion Profile Control

```go
// Generate smooth trapezoidal motion profile
constraints := motionprofile.Constraints{
    MaxVelocity: 5.0,      // max velocity units/sec
    MaxAcceleration: 10.0, // max acceleration units/sec²
}
start := motionprofile.State{Position: 0.0, Velocity: 0.0}
goal := motionprofile.State{Position: 20.0, Velocity: 0.0}

profile := motionprofile.New(constraints, goal, start)

// Use in control loop
for t := 0.0; !profile.IsFinished(t); t += 0.02 {
    setpoint := profile.Calculate(t)
    // Feed setpoint to position controller
    output := positionController.Calculate(setpoint.Position, currentPosition)
}
```

## Feedforward Control

The feedforward package provides predictive control for improved system performance:

### Controller Types

#### Basic Feedforward

```go
ff := feedforward.New()
// Output: kV*velocity + kA*acceleration
```

#### Gravity Compensation

```go
ff := feedforward.New(feedforward.WithGravityGain(9.81))
// Output: kV*velocity + kA*acceleration + kG
```

#### Cosine Compensation

```go
ff := feedforward.New(feedforward.WithCosineGain(2.5))
// Output: kV*velocity + kA*acceleration + kCos*cos(position)
```

#### Combined Compensation

```go
ff := feedforward.New(
    feedforward.WithGravityGain(15.7),
    feedforward.WithCosineGain(8.2),
)
// Output: kV*velocity + kA*acceleration + kG + kCos*cos(position)
```

### Applications

- **Basic FF**: Simple systems without gravitational effects
- **Gravity FF**: Elevators, vertical lifts, load handling
- **Cosine FF**: Robotic arms, rotating machinery, pendulums  
- **Combined FF**: Cranes, construction equipment, complex robotics

### Integration with Feedback Control

```go
// Typical control system structure
totalOutput := pidOutput + feedforwardOutput + feedbackOutput
```

The feedforward provides the bulk of the control effort, while feedback handles disturbances and model uncertainties.

## Performance

- **PID Controllers**: ~64ns per update (15+ million updates/second)
- **Feedforward Controllers**: ~2-5ns per calculation (200+ million calculations/second)
- **Feedback Controllers**: Scales with system dimension
- **Memory Usage**: Minimal allocation after initialization, zero allocations during calculations
- **Precision**: Full float64 precision for all calculations
- **Timing**: Microsecond-accurate time-based integral/derivative calculation

### Feedforward Benchmarks

```text
BenchmarkCalculateBasic-10         577752901    2.120 ns/op    0 B/op    0 allocs/op
BenchmarkCalculateWithGravity-10   579784483    2.066 ns/op    0 B/op    0 allocs/op
BenchmarkCalculateWithCosine-10    230321818    5.203 ns/op    0 B/op    0 allocs/op
BenchmarkCalculateWithBoth-10      229596540    5.189 ns/op    0 B/op    0 allocs/op
```

## Best Practices

### Gain Tuning

1. **Start with P-only**: Set Ki=0, Kd=0, tune Kp for basic response
2. **Add Integral**: Gradually increase Ki to eliminate steady-state error
3. **Add Derivative**: Add small Kd to reduce overshoot and improve stability
4. **Fine-tune**: Adjust all gains for optimal performance

### Feature Guidelines

- **Feed-forward**: Use when you know the expected output for a given setpoint
- **Zero crossing reset**: Useful for systems that can overshoot in either direction
- **Stability threshold**: Helps reduce overshoot in high-speed applications
- **Integral capping**: For systems with strict output limits (e.g., motor controllers)
- **Derivative filtering**: Always recommended for noisy sensor inputs

### Robot Recommendations

```go
// Typical motor controller setup
motorPID := pid.New(kp, ki, kd,
    pid.WithIntegralSumMax(1.0/ki),           // Prevent motor saturation
    pid.WithStabilityThreshold(encoderCPR/4), // Disable integral during rapid movement
    pid.WithDerivativeFilter(0.1),            // Filter encoder noise
)
motorPID.SetOutputLimits(-1.0, 1.0)
```

## Testing

The library includes comprehensive tests with 92.1% coverage:

```bash
# Run tests
go test ./pid

# Run tests with coverage
go test ./pid -cover

# Run benchmarks
go test ./pid -bench=.
```

## Examples

See the `examples/` directory for complete working examples:

### PID Controller Examples

- Basic PID control loop
- Motor speed control  
- Temperature regulation
- Position servo control

### Feedback Controller Examples

- Full state feedback control for multi-dimensional systems
- Vector-based control calculations

### Feedforward Controller Examples

- **Basic Control** (`feedforward/examples/basic/`) - Simple motor control with velocity/acceleration compensation
- **Elevator Control** (`feedforward/examples/elevator/`) - Gravity compensation for vertical movement
- **Robotic Arm** (`feedforward/examples/arm/`) - Cosine compensation for rotating joints
- **Crane Control** (`feedforward/examples/crane/`) - Combined gravity and cosine compensation
- **Controller Comparison** (`feedforward/examples/compare/`) - Side-by-side performance analysis

#### Feedforward Performance

- **Ultra-fast calculations**: 2-5 nanoseconds per operation
- **Zero memory allocations** during calculations
- **100% test coverage** with comprehensive benchmarks
- **Real-world applications**: Elevators, robotic arms, cranes, industrial automation

### Filter Examples

- **Kalman Filter** (`filter/examples/basic/`) - Advanced signal estimation with DARE solver
- **Low-Pass Filter** (`filter/examples/lowpass/`) - Signal smoothing with configurable response

#### Filter Performance

- **Kalman Filter**: ~44 nanoseconds per estimation with optimal gain calculation
- **Low-Pass Filter**: ~4.7 nanoseconds per operation for real-time smoothing
- **92.5% test coverage** with comprehensive validation
- **Real-world applications**: Sensor fusion, noise reduction, signal conditioning

## License

MIT License - see LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## Support

For questions, issues, or contributions, please open an issue on GitHub.
