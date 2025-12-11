# Basic Trapezoidal Motion Profile Example

This example demonstrates how to generate smooth, time-optimal motion profiles for controlled acceleration, constant velocity, and deceleration phases.

## What This Example Shows

- Trapezoidal velocity profile generation
- Smooth acceleration and deceleration phases
- Constant velocity cruise phase
- Position, velocity, and acceleration calculations at any time
- Time-optimal trajectory planning
- Profile constraints (max velocity, max acceleration)

## Running the Example

```bash
cd motionprofile/examples/basic
go run main.go
```

## Key Learning Points

### Motion Profile Characteristics

The example demonstrates:

- **Acceleration Phase**: Linear velocity increase from zero
- **Cruise Phase**: Constant velocity motion
- **Deceleration Phase**: Linear velocity decrease to zero
- **Time Optimality**: Minimum time while respecting constraints
- **Smoothness**: Continuous velocity eliminates jerk spikes

### Profile Constraints

The profile respects:

- **Maximum Velocity**: Speed limit for safety or practical reasons
- **Maximum Acceleration**: Limit based on motor/actuator capability
- **Initial and Final States**: Starting and ending at rest or specific velocities

## Output Interpretation

The example displays:

- **Time**: Current simulation or profile time
- **Position**: Distance traveled along trajectory
- **Velocity**: Current speed (positive = forward)
- **Acceleration**: Current rate of velocity change
- **Phase**: Which phase of the profile (acceleration/cruise/deceleration)

## System Parameters

The example typically uses:

- **Max Velocity**: 1.0 m/s (adjustable)
- **Max Acceleration**: 0.5 m/s² (adjustable)
- **Total Distance**: Initial position to final position
- **Sample Rate**: How often to query the profile (e.g., 10ms)

## Further Exploration

Try modifying:

- `maxVelocity` - Test different speed limits
- `maxAcceleration` - Simulate different actuators
- Distance or initial/final positions - Different trajectories
- Sample times - See how to query at arbitrary times
- Constraints - Test edge cases and special motions

## Real-World Applications

Motion profiles are essential in:

- **CNC Machining**: Precise tool path generation
- **Robotics**: Smooth joint motion and end-effector paths
- **Manufacturing**: Synchronized multi-axis motion
- **Printing**: Carriage and paper feed control
- **Elevators**: Comfortable acceleration/deceleration
- **Autonomous Vehicles**: Smooth acceleration profiles
- **Industrial Automation**: Coordinated equipment motion
- **Movie/Animation**: Smooth camera movements

## Related Examples

- `../triangle/` - Alternative (triangular) profile without cruise phase
- `../fullstate_control/` - Integrating profiles with feedback control
- `../../feedforward/examples/basic/` - Using profiles with feedforward control
- `../../pid/examples/` - Simple setpoint tracking without profiles

## Motion Profile Phases

### Phase 1: Acceleration

- Duration: $t_1 = v_{max} / a_{max}$
- Distance: $s_1 = \frac{v_{max}^2}{2 \cdot a_{max}}$
- Equation: $v(t) = a_{max} \cdot t$

### Phase 2: Cruise

- Velocity: Constant at $v_{max}$
- Distance: $s_2 = s_{total} - 2 \cdot s_1$
- Duration: $t_2 = s_2 / v_{max}$

### Phase 3: Deceleration

- Symmetric to acceleration phase
- Duration: $t_3 = t_1$
- Distance: $s_3 = s_1$

## Advantages of Trapezoidal Profiles

- **Time Optimal**: Minimum time while respecting constraints
- **Smooth**: No discontinuous acceleration (eliminates jerk spikes)
- **Simple**: Analytical equations, easy to compute
- **Predictable**: Exactly controlled acceleration
- **Universal**: Works for any 1D motion task

## Limitations

**No Jerk Limiting**: Sudden acceleration changes at phase transitions
**Discontinuity**: May not be optimal if intermediate constraints exist
**Overshoot Risk**: If feedback control overshoots profile values
**Symmetry Assumption**: Assumes forward and reverse acceleration are equal

## Implementation Details

### Calculating Profile State at Arbitrary Time

Given time $t$ and knowing phase durations:

1. Determine which phase: $t \in [0, t_1]$, $[t_1, t_1+t_2]$, or later
2. Apply appropriate equation for position and velocity
3. Acceleration is constant per phase

### Edge Cases

**Distance Too Small**: Can't reach max velocity (triangle profile)
**Zero Duration**: If distance is zero
**Already at Target**: No motion needed

## Advanced Concepts

- **S-Curve Profiles**: Adding jerk limits for smoother transitions
- **Multi-Axis Profiles**: Coordinating multiple moving axes
- **Path Planning**: Connecting multiple motion segments
- **Optimal Control**: Minimizing energy or time with constraints
- **Real-Time Generation**: Computing profiles on-demand

## Performance Expectations

A well-implemented motion profile:

- Reaches target position exactly
- Respects all velocity and acceleration constraints
- Provides smooth motion (no discontinuities)
- Can be queried at any time
- Computes quickly (< 1 microsecond)

## Integration with Control

Profiles are typically used with:

- **Feedforward Control**: Predict required command from profile
- **PID Control**: Track the profile velocity setpoint
- **State Feedback**: Use profile as reference for all states
- **Feedforward + Feedback**: Combine for best results
