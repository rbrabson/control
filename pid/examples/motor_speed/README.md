# Motor Speed Control Example

This example demonstrates PID control applied to a motor speed regulation system, a classic real-world control problem.

## What This Example Shows

- PID control of a motor's rotational speed
- Handling of load disturbances
- Steady-state error elimination with integral term
- Real-world constraints (maximum torque/power)
- Speed profile with transient and steady-state behavior

## Running the Example

```bash
cd pid/examples/motor_speed
go run main.go
```

## Key Learning Points

### Motor Control Challenges

The example demonstrates:

- **Speed Regulation**: Maintaining desired RPM despite load changes
- **Load Disturbance**: How the controller responds to sudden load increases
- **Steady-State Accuracy**: Using the integral term to eliminate error
- **Torque Limits**: Realistic physical constraints on motor output

### Real-World Considerations

- Motor inertia (resistance to acceleration)
- Load torque (friction and mechanical resistance)
- Speed feedback from encoder or tachometer
- Limited torque/current available

## Output Interpretation

The example displays:

- **Target Speed**: Desired motor speed in RPM or normalized units
- **Time**: Elapsed simulation time
- **Actual Speed**: Current motor speed
- **Error**: Speed deviation from target
- **Torque Output**: Control signal to motor driver

## System Parameters

The example typically uses:

- **Target Speed**: 100 RPM or normalized units
- **Motor Inertia**: Affects acceleration capability
- **Load Torque**: Represents friction and mechanical load
- **PID Gains**: Tuned for smooth speed tracking
- **Control Update Rate**: 10ms or 100Hz

## Further Exploration

Try modifying:

- `targetSpeed` - Test different speed profiles
- `loadTorque` - Simulate sudden load changes
- PID gains - Observe effects on transient response
- `motorInertia` - Simulate lighter or heavier rotors

## Real-World Applications

This control technique is used in:

- Electric motors in robotics and automation
- Vehicle traction control systems
- Industrial motor speed regulators
- Fan and pump speed controllers
- Conveyor belt systems

## Related Examples

- `../basic_control_loop/` - Simpler example to start with
- `../temperature_control/` - Similar control problem in thermal domain
- `../dampening/` - Advanced techniques to improve response
- `../../feedforward/examples/` - Combining with feedforward control

## Common Tuning Tips

**Sluggish Response**: Increase Kp for faster speed adjustment
**Speed Oscillation**: Reduce Ki or increase Kd
**Overshoot**: Increase Kd to add damping
**Steady-State Error**: Increase Ki for better integral action
