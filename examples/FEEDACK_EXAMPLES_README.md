# Feedback Control Examples

This directory contains examples demonstrating how to use the `control/feedback` package.

## Examples Overview

### 1. Basic Feedback Control (`feedback_control/`)

Demonstrates the fundamental usage of feedback controllers:

- **NoFeedback Controller**: Always returns zero output, useful for open-loop control
- **Single State Feedback**: Simple proportional control for position
- **Multi-State Feedback**: Position + velocity control for better stability  
- **Three State Feedback**: Position, velocity, and acceleration control
- **Controller Comparison**: Side-by-side comparison of different controllers
- **Error Handling**: Proper error handling for mismatched vector dimensions
- **Motion Profile Integration**: Trajectory following simulation

**Run the example:**

```bash
cd examples/feedback_control
go run main.go
```

### 2. Combined PID + Feedback Control (`combined_control/`)

Shows integration of PID controllers with state feedback:

- **Robot Arm Simulation**: PID for position control + state feedback for stabilization
- **Adaptive Control**: Switching between different feedback strategies based on conditions
- **Performance Comparison**: Shows how combining controllers can improve performance

**Run the example:**

```bash
cd examples/combined_control
go run main.go
```

## Key Concepts Demonstrated

### NoFeedback Controller

```go
noFeedback := &feedback.NoFeedback{}
output := noFeedback.Calculate(setpoint, measurement) // Always returns 0.0
```

### FullStateFeedback Controller

```go
// Single dimension (position only)
gains := feedback.Values{2.5}
controller := feedback.NewFullStateFeedback(gains)

target := feedback.Values{10.0}
current := feedback.Values{8.2}
output, err := controller.Calculate(target, current)
```

### Multi-Dimensional Control

```go
// Position and velocity feedback
gains := feedback.Values{1.8, 0.4} // [position_gain, velocity_gain]
controller := feedback.NewFullStateFeedback(gains)

target := feedback.Values{5.0, 0.0}    // [position, velocity]
current := feedback.Values{3.5, 2.1}   // [current_pos, current_vel]
output, err := controller.Calculate(target, current)
```

### Error Handling

```go
// Always check for errors with FullStateFeedback
output, err := controller.Calculate(target, current)
if err != nil {
    // Handle error (usually mismatched vector lengths)
    fmt.Printf("Error: %v\n", err)
    return
}
```

## Common Use Cases

1. **Position Control**: Single-dimension feedback for simple positioning
2. **Motion Control**: Multi-dimension feedback for position + velocity
3. **Stabilization**: Adding feedback to existing PID controllers
4. **Trajectory Following**: Using state feedback for path tracking
5. **Open-Loop Control**: Using NoFeedback as a placeholder or for pure feedforward control

## Tips for Usage

- **Vector Dimensions**: Ensure target and current vectors have the same length as the gain vector
- **Gain Tuning**: Start with smaller gains and increase gradually
- **Error Checking**: Always check for errors when using FullStateFeedback
- **Interface Compatibility**: Note that NoFeedback and FullStateFeedback have different interfaces
- **Performance**: NoFeedback is extremely fast (~0.3 ns/op), FullStateFeedback scales with dimension

## Related Packages

- `control/pid`: PID controllers that can be combined with feedback controllers
- Standard library packages used in examples: `fmt`, `math`, `time`
