# Basic Shooter Velocity Mapping Example

This example demonstrates how to use interpolating lookup tables (InterpLUT) for smooth velocity mapping in robotic shooters or similar systems.

## What This Example Shows

- Creating smooth interpolation tables with cubic splines
- Mapping input (distance) to output (velocity)
- Achieving exact values at control points
- Smooth interpolation between discrete settings
- Error handling for out-of-bounds inputs
- Practical robotics application

## Running the Example

```bash
cd interplut/examples/basic
go run main.go
```

## Key Learning Points

### Interpolation Fundamentals

The example demonstrates:

- **Control Points**: Discrete input-output pairs you define
- **Smooth Interpolation**: Filling in values between control points
- **Cubic Splines**: Smooth curve fitting without oscillation
- **Exact Matching**: Output exactly matches control points
- **Monotonicity**: Preserving increasing/decreasing behavior

### Shooter Velocity Mapping

In robotic shooters:

- Shooting distance varies (different field positions)
- Different distances require different velocities
- Manual tuning for each distance is tedious
- Interpolation allows smooth velocity variation
- Eliminates need for discrete velocity settings

## Output Interpretation

The example displays:

- **Distance (m)**: Target distance in meters
- **Velocity (%)**: Required shooter speed as percentage or units
- **Accuracy**: How close to expected values at control points

## System Parameters

The example typically uses:

- **Control Points**: Calibrated distance-to-velocity pairs
- **Distance Range**: Minimum to maximum distances
- **Velocity Range**: Minimum to maximum shooter speeds

## Further Exploration

Try modifying:

- Add more control points - Improve accuracy across range
- Change distance values - Test different field sizes
- Modify velocity values - Simulate different shooter configurations
- Query new distances - See interpolated values
- Test edge cases - Out-of-bounds behavior

## Real-World Applications

Interpolation is used in:

- **Robot Shooters**: Turrets for game pieces (FIRST Robotics)
- **Gun Fire Control**: Weapon elevation and distance compensation
- **Manufacturing**: Tool selection based on material
- **Motors**: Speed control based on load
- **Dosing Systems**: Chemical or fluid quantities based on parameters
- **Audio**: Volume/frequency adjustment curves
- **Gaming**: Difficulty scaling based on progression
- **Thermal**: Temperature control based on load

## Related Examples

- `../temperature/` - Non-linear temperature control
- `../adaptive_pid/` - Adaptive PID using interpolation
- `../../pid/examples/` - Pure feedback control
- `../../feedforward/examples/` - Model-based control

## Advantages of Interpolation

- **Flexibility**: Easy to adjust behavior by changing control points
- **Smoothness**: Smooth transitions without lookup artifacts
- **Accuracy**: Can match any non-linear function
- **Simplicity**: No complex model or tuning needed
- **Real-Time**: Fast computation for control loops

## Cubic Spline Properties

The cubic splines used in interpolation have:

- **Continuous Second Derivative**: Smooth curvature
- **Local Support**: Each spline segment depends only on nearby points
- **Uniqueness**: One solution for given control points
- **Natural Boundaries**: Optimized endpoints for open/closed curves

## Implementation Considerations

- **Control Point Selection**: Choose points across full operating range
- **Accuracy**: More points provide better accuracy but add computation
- **Validation**: Test interpolated values against actual system
- **Bounds Checking**: Handle distances outside calibrated range
- **Updates**: Recalibrate control points as system changes

## Common Pitfalls

**Too Few Points**: Inaccurate interpolation between distant points
**Noisy Data**: Control points from noisy measurements create wiggles
**Extrapolation**: Values beyond calibrated range may be unreliable
**Discontinuity**: Sudden shooter velocity changes are undesirable

## Performance Expectations

A well-designed interpolation system:

- Evaluates queries in < 1 microsecond
- Matches control point values exactly
- Produces smooth transitions
- Handles hundreds of control points efficiently

## Advanced Techniques

- **2D Interpolation**: Map multiple inputs (distance, angle) to output
- **Reverse Lookup**: Find input that produces desired output
- **Constrained Splines**: Enforce monotonicity or smoothness
- **Adaptive Points**: Add points in high-curvature regions
- **Error Estimation**: Quantify interpolation accuracy
