# Triangle Motion Profile Example

This example demonstrates an alternative motion profile that uses only acceleration and deceleration phases, without a constant velocity cruise phase.

## What This Example Shows

- Triangle velocity profile generation
- Symmetric acceleration and deceleration
- No constant velocity phase
- Shorter total motion time than trapezoidal profiles
- Maximum velocity achieved at the midpoint
- Time-optimal for shorter distances

## Running the Example

```bash
cd motionprofile/examples/triangle
go run main.go
```

## Key Learning Points

### Triangle Profile Characteristics

The example demonstrates:

- **Acceleration Phase**: Linear velocity increase to peak
- **Deceleration Phase**: Linear velocity decrease to zero
- **Peak Velocity**: Achieved at midpoint of motion
- **Symmetry**: Acceleration and deceleration phases are identical
- **No Cruise**: Transitions directly from acceleration to deceleration

### When to Use Triangle Profiles

- Short distances where peak velocity isn't reached
- Systems with tight space constraints
- Avoiding unnecessary constant-velocity motion
- Simpler implementation (no three-phase logic)
- Faster point-to-point motion

## Output Interpretation

The example displays:

- **Time**: Current profile time
- **Position**: Distance traveled from start
- **Velocity**: Current speed (peaks at midpoint)
- **Acceleration**: Always maximum or negative (only two values)
- **Progress**: Percentage of total motion completed

## System Parameters

The example typically uses:

- **Max Acceleration**: Maximum acceleration rate
- **Total Distance**: Starting position to final position
- **Sample Rate**: Frequency of profile queries

## Further Exploration

Try modifying:

- `maxAcceleration` - Test different accelerator capabilities
- Distance values - Try various motion distances
- Sample times - Query the profile at different intervals
- Compare with trapezoidal - Notice timing differences

## Real-World Applications

Triangle profiles are used in:

- **Point-to-Point Motion**: Quick positioning without dwelling
- **Compact Systems**: Space-limited machinery
- **Quick Pick-and-Place**: Robotic assembly operations
- **Short-Distance Moves**: Within tight spaces
- **Energy-Efficient**: Minimizes time and energy
- **Simple Hardware**: Less sophisticated motion control

## Related Examples

- `../basic/` - Trapezoidal profile (adds cruise phase)
- `../fullstate_control/` - Integrating profiles with feedback control
- `../../feedforward/examples/basic/` - Using profiles with feedforward control

## Comparison: Triangle vs. Trapezoidal

| Aspect | Triangle | Trapezoidal |
|--------|----------|-------------|
| Phases | 2 | 3 |
| Cruise Phase | No | Yes |
| Max Velocity | At midpoint | Sustained |
| Peak Accel | Higher | Lower |
| Best For | Short distance | Long distance |
| Implementation | Simpler | More complex |
| Motion Time | Potentially faster | May be faster if cruise is long |

## Triangle Profile Equations

For motion from position $s_0$ to $s_f$ with maximum acceleration $a_{max}$:

- **Total Distance**: $\Delta s = s_f - s_0$
- **Acceleration Time**: $t_a = \sqrt{\Delta s / a_{max}}$
- **Peak Velocity**: $v_{max} = a_{max} \cdot t_a$
- **Total Time**: $t_{total} = 2 \cdot t_a$

At time $t$ during acceleration phase:

- $v(t) = a_{max} \cdot t$
- $s(t) = s_0 + \frac{1}{2} a_{max} \cdot t^2$

At time $t$ during deceleration phase (with $\tau = t - t_a$):

- $v(t) = v_{max} - a_{max} \cdot \tau$
- $s(t) = s_{midpoint} + v_{max} \cdot \tau - \frac{1}{2} a_{max} \cdot \tau^2$

## Advantages

- **Simplicity**: Only two phases makes logic straightforward
- **Fast**: No wasted time in cruise phase
- **Predictable**: Analytical equations are simple
- **Symmetric**: Forward and reverse accelerations are equal

## Limitations

- **High Peak Acceleration**: May require stronger actuators
- **Not Suitable for Long Distances**: Velocities can get very high
- **Physical Constraints**: May exceed velocity limits
- **Jerk Spikes**: Discontinuities at phase transitions (like trapezoidal)

## Implementation Considerations

### Handling Velocity Limits

If the calculated peak velocity exceeds the maximum allowed:

- Switch to trapezoidal profile
- Reduce acceleration
- Limit peak velocity and adjust distance

### Querying at Arbitrary Times

The profile can be evaluated at any time, automatically:

1. Determine which phase: $t < t_a$ (acceleration) or $t \geq t_a$ (deceleration)
2. Apply appropriate equations
3. Return position, velocity, acceleration

## Performance Expectations

Triangle profile motion:

- Reaches target position exactly
- Respects acceleration constraints
- Can be queried at any time
- Computes instantly
- Works well for short-distance moves

## Advanced Variations

- **Asymmetric Triangle**: Different acceleration and deceleration rates
- **Velocity-Limited**: Triangles that exceed max velocity cap
- **Multi-Axis**: Coordinating triangles on multiple axes
- **Path Planning**: Concatenating multiple triangle profiles

## When to Use

**Choose Triangle When:**

- Moving short distances
- Actuators can handle high acceleration
- Smooth constant-velocity cruise isn't needed
- Simplicity is important
- Point-to-point speed is critical

**Choose Trapezoidal When:**

- Moving longer distances
- Need sustained constant velocity
- Peak acceleration needs to be limited
- Comfort or smoothness is important
- Complex motion requires coordination
