# Full State Feedback Control Example

This example demonstrates state feedback control for multi-dimensional systems, a powerful technique for controlling complex systems with multiple state variables.

## What This Example Shows

- Single-state feedback control for simple systems
- Multi-state feedback (position and velocity) for improved control
- 3D system control (position, velocity, and acceleration)
- Error handling for dimension mismatches
- Performance comparison between control strategies
- Gain matrix design for optimal control

## Running the Example

```bash
cd feedback/examples/feedback_control
go run main.go
```

## Key Learning Points

### State Feedback Concepts

The example demonstrates:

- **State Vector**: Representing system with multiple variables
- **Gain Matrix**: Feedback gains for each state component
- **Full State Knowledge**: Assumes all states are measurable
- **Decoupling**: Controlling multiple outputs independently
- **Closed-Loop Poles**: Placing system dynamics through feedback

### Single vs. Multi-State Feedback

1. **1D Feedback**: Simple proportional control of one variable
2. **2D Feedback**: Position and velocity control (most common)
3. **3D Feedback**: Position, velocity, and acceleration (advanced)

## Output Interpretation

The example displays:

- **Target State**: Desired values for each state component
- **Current State**: Actual measured state
- **Control Output**: Calculated control signal
- **Dimension**: Number of state variables (1D, 2D, or 3D)
- **Performance**: Time to calculate (ns/op)

## System Parameters

The example typically uses:

- **State Variables**: 1 to 3 dimensions
- **Gain Vector**: Different gain for each state component
- **Sample States**: Typical system values at different operating points

## Further Exploration

Try modifying:

- Gain values - Observe effect on control signal
- Target states - Test different setpoints
- State dimensions - Add or remove state variables
- Initial state values - Test convergence from different starting points

## Real-World Applications

State feedback is used in:

- Aircraft autopilot and stability augmentation
- Spacecraft attitude control
- Rocket launch vehicle guidance
- Quadcopter/drone stabilization
- Robotic manipulator control
- Active suspension systems
- Power system frequency control
- Industrial process control with multiple variables

## Related Examples

- `../../pid/examples/` - Simpler single-variable control
- `../../motionprofile/examples/fullstate_control/` - State feedback with motion planning
- `../../feedforward/examples/` - Combining with feedforward

## Advanced Concepts

- **Linear Quadratic Regulator (LQR)**: Optimal gain computation
- **State Observer**: Estimating unmeasured states
- **Pole Placement**: Designing desired closed-loop dynamics
- **Decoupling**: Managing coupling between multiple outputs

## Comparison with PID Control

| Aspect | PID | State Feedback |
|--------|-----|-----------------|
| State Variables | 1 | Multiple |
| Design Method | Empirical Tuning | Pole Placement/LQR |
| Model Requirement | Minimal | Detailed Model |
| Implementation | Simpler | More Complex |
| Multi-Output | Difficult | Natural |
| Optimal Performance | Approximate | Provable |

## Implementation Considerations

- **State Measurement**: All states must be measurable or estimated
- **Sampling Rate**: Fast enough to capture system dynamics
- **Computation Load**: Minimal for 1-3 states, grows with dimensionality
- **Robustness**: Sensitive to model errors and unmodeled dynamics
- **Noise Sensitivity**: May need observer for noisy measurements
