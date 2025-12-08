# PID Controller Library

A comprehensive, high-performance PID (Proportional-Integral-Derivative) controller implementation in Go with advanced features for robotics and control systems.

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
- **Derivative low-pass filtering** to reduce measurement noise

### Performance & Reliability

- **~64ns per update** - Excellent performance for real-time applications
- **92.1% test coverage** with comprehensive test suite
- **Thread-safe design** for concurrent applications
- **Robust error handling** for edge cases and invalid inputs

## Installation

```bash
go get github.com/rbrabson/control
```

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

## Advanced Usage

### Using Options Pattern

```go
// Create PID with advanced features
controller := pid.New(1.0, 0.1, 0.05,
    pid.WithFeedForward(5.0),                    // Add feed-forward term
    pid.WithIntegralResetOnZeroCross(),          // Reset integral at zero crossing
    pid.WithStabilityThreshold(2.0),             // Disable integral when derivative > 2.0
    pid.WithIntegralSumMax(10.0),                // Cap integral sum at ±10.0
    pid.WithDerivativeFilter(0.1),               // Apply 10% low-pass filter
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

#### Derivative Filtering

```go
// Option function
pid.WithDerivativeFilter(alpha float64) // 0-1, where 0=no filter, 1=max filter

// Runtime methods
controller.SetDerivativeFilter(alpha float64)
controller.GetDerivativeFilter() float64
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

## Performance

- **Update Rate**: ~64ns per update (15+ million updates/second)
- **Memory Usage**: Minimal allocation after initialization
- **Precision**: Full float64 precision for all calculations
- **Timing**: Microsecond-accurate time-based integral/derivative calculation

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

- Basic PID control loop
- Motor speed control
- Temperature regulation
- Position servo control

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
