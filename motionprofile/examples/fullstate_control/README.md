# Motion Profile with Full-State Feedback Control

This example demonstrates how to combine a trapezoidal motion profile generator with a full-state feedback controller to achieve precise trajectory tracking.

## Overview

The example shows:

- **Motion Profile Generation**: Creates smooth trapezoidal velocity profiles
- **Full-State Feedback Control**: Simultaneous position and velocity control
- **System Simulation**: Physics-based simulation with mass and damping
- **Performance Analysis**: Tracking error metrics and system behavior

## Key Concepts

### Full-State Feedback

Unlike PID control which typically uses only position error, full-state feedback uses both position and velocity errors simultaneously:

``` go
control_output = Kp * position_error + Kd * velocity_error
```

This approach provides:

- **Predictive Control**: Responds to reference velocity changes before position errors develop
- **Better Tracking**: Optimal following of smooth trajectories
- **Natural Damping**: Velocity feedback provides inherent system damping
- **Single Controller**: One controller handles the entire trajectory

### Motion Profiles

Trapezoidal motion profiles provide:

- **Smooth Motion**: Continuous velocity with bounded acceleration
- **Predictable Timing**: Known completion time and intermediate states
- **Optimal Efficiency**: Maximum use of available velocity and acceleration

## Running the Example

```bash
cd motionprofile/examples/fullstate_control
go run main.go
```

## Example Output

``` code
Motion Profile with Full-State Feedback Control Example
======================================================

Motion Profile Generated:
- Distance: 5.0 m
- Max Velocity: 2.0 m/s
- Max Acceleration: 1.0 m/s²
- Total Time: 4.50 s

Performance Analysis:
====================
Maximum Tracking Error: 0.082 m (1.6%)
RMS Tracking Error: 0.051 m (1.0%)
Final Settling Error: 0.000 m (0.00%)
✅ Motion completed successfully!
```

## System Configuration

### Motion Profile

- **Distance**: 5.0 meters
- **Max Velocity**: 2.0 m/s
- **Max Acceleration**: 1.0 m/s²

### Controller Gains

- **Position Gain (Kp)**: 50 N/m - Controls position tracking stiffness
- **Velocity Gain (Kd)**: 20 N·s/m - Provides damping and velocity tracking

### Simulated System

- **Mass**: 1.0 kg
- **Damping**: 2.0 N·s/m
- **Integration**: Euler method with 10ms timestep

## Tuning Guidelines

### Position Gain (Kp)

- **Higher Values**: Faster position response, may cause overshoot
- **Lower Values**: Slower response, better stability
- **Typical Range**: 10-100 N/m depending on system stiffness requirements

### Velocity Gain (Kd)

- **Higher Values**: More damping, smoother response, may slow response
- **Lower Values**: Less damping, faster response, may oscillate
- **Typical Range**: 5-50 N·s/m depending on system damping needs

### Balancing Guidelines

1. **Start Conservative**: Begin with lower gains and increase gradually
2. **Test Tracking**: Verify good reference following during motion
3. **Check Stability**: Ensure no oscillations or instability
4. **Optimize Performance**: Adjust for your specific accuracy requirements

## Applications

This control approach is ideal for:

- **CNC Machines**: Precise tool path following
- **Robotic Arms**: Smooth joint trajectory execution
- **Conveyor Systems**: Coordinated material handling
- **Positioning Systems**: High-accuracy motion control
- **3D Printers**: Precise extruder head positioning

## Comparison with PID

| Aspect | Full-State Feedback | PID Control |
|--------|-------------------|-------------|
| **Inputs** | Position + Velocity Error | Position Error + Derivatives |
| **Tracking** | Excellent for smooth trajectories | Good for step responses |
| **Tuning** | 2 gains (Kp, Kd) | 3 gains (Kp, Ki, Kd) |
| **Complexity** | Moderate | Higher (integral windup, derivative noise) |
| **Performance** | Optimal for known references | General purpose |

Full-state feedback excels when you have smooth reference trajectories and want optimal tracking performance with simpler tuning.
