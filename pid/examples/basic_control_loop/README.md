# Basic Control Loop Example

This example demonstrates the fundamental usage of the PID controller with a simple simulated system.

## What This Example Shows

- Creating a basic PID controller
- Setting output limits for safety and realism
- Running a control loop with proper timing
- Monitoring system response and convergence
- Simple first-order system dynamics simulation

## Running the Example

```bash
cd pid/examples/basic_control_loop
go run main.go
```

## Key Learning Points

### PID Controller Basics

The example demonstrates:

- **PID Gains**: How proportional, integral, and derivative terms affect control
- **Output Limits**: Constraining controller output (max ±50 in this example)
- **Feedback Loop**: Real-time adjustment based on error
- **Convergence**: How the system reaches the setpoint

### Control Loop Timing

- Uses Go's `time` package for accurate loop timing
- Demonstrates real-time control concepts
- Shows sample-time impact on control performance

## Output Interpretation

The example displays:

- **Target**: The desired setpoint
- **Time**: Elapsed simulation time in seconds
- **Position**: Current system state
- **Error**: Difference between target and actual position
- **Output**: Control signal from PID controller

## System Parameters

The example uses:

- **Setpoint**: 100.0 units
- **Time Constant**: 2.0 seconds (system response speed)
- **PID Gains**: Kp=1.0, Ki=0.1, Kd=0.5 (tuned for stable response)
- **Output Limits**: ±50 units (representing physical constraints)
- **Sample Time**: 10ms (0.01 seconds)

## Further Exploration

Try modifying:

- `setpoint` - Change the target value
- `Kp, Ki, Kd` - Tune the controller for different response characteristics
- `maxOutput` - Adjust saturation limits
- `timeConstant` - Simulate different system dynamics

## Comparison with Other Examples

- `../motor_speed/` - Real-world application with motor control
- `../temperature_control/` - Another real-world example with thermal dynamics
- `../position_servo/` - Multi-turn servo control
- `../dampening/` - Advanced dampening techniques
- `../filter_comparison/` - Adding filtering for noisy measurements

## Common Challenges

**Overshoot**: If the controller overshoots, increase Kd or decrease Kp
**Slow Response**: If convergence is too slow, increase Kp
**Oscillation**: If the output oscillates, reduce Ki or increase Kd
