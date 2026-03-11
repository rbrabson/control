# Basic PID Controller Example

This example demonstrates the core PID controller features, matching the behavior tested in PIDTest.java.

## What This Example Shows

- Output limits that clamp controller output
- Derivative filter option for smoother control
- Combining multiple PID options
- Basic PID calculation usage

## Running the Example

```bash
cd pid/examples/basic_control_loop
go run main.go
```

## Key Learning Points

### Test 1: Output Limits

Shows how to constrain PID output to safe operating ranges:

- Creates PID with Kp=10.0 and limits [-5.0, 5.0]
- Even with large error (10.0), output is clamped
- Critical for real systems with actuator limits

### Test 2: Derivative Filter

Demonstrates derivative filtering to reduce noise sensitivity:

- Creates PID with Kd=1.0 and low-pass filter (alpha=0.8)
- Filter smooths derivative term calculations
- Produces stable, finite output despite error changes

### Test 3: Combined Features

Shows combining multiple PID options:

- All three gains (P, I, D)
- Output limits and derivative filter together
- Simple control loop demonstration

## Application

Use this pattern when:

- Validating PID controller setup
- Testing output limiting behavior
- Verifying filter integration
- Learning PID basics

## See Also

- [Simulation Example](../simulation/) - Time-stepping PID control
- [Temperature Control](../temperature_control/) - Real-world application
- [Filter Comparison](../filter_comparison/) - Different filter strategies

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
