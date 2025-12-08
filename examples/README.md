# Control Systems Examples

This directory contains complete working examples demonstrating various uses of both the PID controller and feedback control libraries.

## Running the Examples

Each example is in its own subdirectory and can be run independently:

### Running PID Control Examples

```bash
# Basic control loop example
cd basic_control_loop
go run basic_control_loop.go

# Motor speed control example  
cd motor_speed
go run main.go

# Temperature control example
cd temperature_control
go run main.go

# Position servo control example
cd position_servo  
go run main.go
```

### Running Feedback Control Examples

```bash
# Basic feedback control example
cd feedback_control
go run main.go

# Combined PID + feedback control example
cd combined_control
go run main.go
```

## Example Descriptions

### PID Control Examples

#### Basic Control Loop (`basic_control_loop/`)

Demonstrates fundamental PID usage with a simple simulated first-order system:

- Basic PID controller creation and configuration
- Output limiting with options pattern
- Simple control loop with timing
- System response monitoring

#### Motor Speed Control (`motor_speed/`)

Shows advanced PID features for motor control applications:

- Integral sum limiting to prevent motor saturation
- Stability threshold to reduce overshoot during rapid changes
- Derivative filtering for encoder noise reduction
- Realistic motor dynamics simulation with noise
- Multiple setpoint changes demonstrating response characteristics

#### Temperature Control (`temperature_control/`)

Demonstrates thermal system control with environmental factors:

- Feed-forward control for ambient temperature compensation
- Integral reset on zero crossover to prevent overshoot
- Derivative filtering for temperature sensor noise
- Realistic thermal dynamics with heat capacity and loss
- Ambient temperature disturbances

#### Position Servo Control (`position_servo/`)

Illustrates precise positioning applications:

- Integral windup protection for large position errors
- Stability threshold during high-speed movement
- Encoder quantization effects
- Realistic servo dynamics with backlash and friction
- Multiple target positions with accuracy analysis

### Feedback Control Examples

#### Basic Feedback Control (`feedback_control/`)

Demonstrates the fundamental usage of feedback controllers:

- **NoFeedback Controller**: Always returns zero output, useful for open-loop control
- **Single State Feedback**: Simple proportional control for position
- **Multi-State Feedback**: Position + velocity control for better stability  
- **Three State Feedback**: Position, velocity, and acceleration control
- **Controller Comparison**: Side-by-side comparison of different controllers
- **Error Handling**: Proper error handling for mismatched vector dimensions
- **Motion Profile Integration**: Trajectory following simulation

#### Combined PID + Feedback Control (`combined_control/`)

Shows integration of PID controllers with state feedback:

- **Robot Arm Simulation**: PID for position control + state feedback for stabilization
- **Adaptive Control**: Switching between different feedback strategies based on conditions
- **Performance Analysis**: Shows how combining controllers can improve performance

## Key Usage Patterns

### PID Controller Usage

```go
// Basic PID controller with options
controller := pid.New(2.0, 0.5, 0.1, 
    pid.WithOutputLimits(-10.0, 10.0))

// Control loop using Calculate method
for {
    output := controller.Calculate(setpoint, measurement)
    // Apply output to system
}
```

### Feedback Controller Usage

```go
// NoFeedback controller (always returns 0)
noFeedback := &feedback.NoFeedback{}
output := noFeedback.Calculate(setpoint, measurement)

// FullStateFeedback controller
gains := feedback.Values{2.5, 0.4}  // position, velocity gains
controller := feedback.NewFullStateFeedback(gains)

target := feedback.Values{10.0, 0.0}    // target position, velocity
current := feedback.Values{8.2, 1.5}    // current position, velocity
output, err := controller.Calculate(target, current)
```

### Combined Control Usage

```go
// Combine PID and state feedback
pidController := pid.New(2.0, 0.1, 0.05, pid.WithOutputLimits(-10.0, 10.0))
stateFeedback := feedback.NewFullStateFeedback(feedback.Values{0.8, 0.3})

// In control loop
pidOutput := pidController.Calculate(setpoint, measurement)
stateOutput, _ := stateFeedback.Calculate(
    feedback.Values{0.0, 0.0},           // target error and velocity  
    feedback.Values{error, velocity})     // current error and velocity
totalOutput := pidOutput + stateOutput
```

## Key Learning Points

Each example demonstrates different aspects of control systems:

### PID Control Concepts

1. **Basic Usage**: How to create, configure, and use a PID controller
2. **Advanced Features**: When and how to use integral reset, stability thresholds, and filtering
3. **Real-World Considerations**: Noise, saturation, mechanical effects, and disturbances
4. **Performance Tuning**: Different gain values for different system characteristics
5. **System Analysis**: Monitoring controller performance and system response

### Feedback Control Concepts

1. **Interface Design**: Clean abstractions for different control strategies
2. **Multi-Dimensional Control**: Handling position, velocity, and acceleration feedback
3. **Error Handling**: Proper validation for vector-based controllers
4. **Performance Scaling**: Understanding computational costs for different approaches
5. **Integration Patterns**: Combining multiple control strategies effectively

## Modifying the Examples

You can experiment with the examples by:

### PID Examples

- Changing PID gains to see effects on response
- Modifying system parameters (time constants, noise levels, etc.)
- Adding or removing advanced PID features
- Implementing different setpoint profiles
- Adding disturbances to test robustness

### Feedback Examples

- Adjusting feedback gains for different response characteristics
- Changing vector dimensions for different control scenarios
- Experimenting with different controller combinations
- Implementing custom feedback strategies
- Testing error conditions and edge cases

## Common Use Cases

1. **Position Control**: Single-dimension feedback for simple positioning
2. **Motion Control**: Multi-dimension feedback for position + velocity
3. **Stabilization**: Adding feedback to existing PID controllers
4. **Trajectory Following**: Using state feedback for path tracking
5. **Open-Loop Control**: Using NoFeedback as a placeholder or for pure feedforward control

## Example Output

Each example produces formatted output showing:

- Time progression
- Setpoint and measured values
- Error and control output
- System state information
- Final performance metrics

This helps visualize controller performance and understand system behavior under different conditions.

## Tips for Usage

### PID Controllers

- Start with proportional gain only, then add integral and derivative
- Use output limits to prevent actuator saturation
- Enable integral reset for systems with frequent setpoint changes
- Apply derivative filtering in noisy environments

### Feedback Controllers

- **Vector Dimensions**: Ensure target and current vectors have the same length as the gain vector
- **Gain Tuning**: Start with smaller gains and increase gradually
- **Error Checking**: Always check for errors when using FullStateFeedback
- **Interface Compatibility**: Note that NoFeedback and FullStateFeedback have different interfaces
- **Performance**: NoFeedback is extremely fast (~0.3 ns/op), FullStateFeedback scales with dimension
