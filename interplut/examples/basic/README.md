# InterpLUT Basic Example

This example demonstrates that InterpLUT returns exact values at control points, matching the behavior tested in the Java implementation.

## What This Example Shows

- Creating an interpolated lookup table with control points
- Exact value returns at control points
- Smooth interpolation between points
- Simple, focused demonstration of core functionality

## Running the Example

```bash
cd interplut/examples/basic
go run main.go
```

## Key Learning Points

### Core Functionality

The example demonstrates:

- **Control Points**: Define input-output pairs
- **Exact Matching**: Returns exact values at control points (within floating-point precision)
- **Interpolation**: Smooth values between control points
- **Error Handling**: Proper error checking

### Interpolation Fundamentals

InterpLUT uses cubic spline interpolation:

- Creates smooth curves through control points
- Guarantees passing through each control point exactly
- Maintains monotonicity if control points are monotonic
- No oscillation between points

## Output Interpretation

The example displays:

- **Exact Control Point Returns**: Verifies that querying at control points returns exact values
- **Interpolated Values**: Shows smooth interpolation between points
- **Match Indicators**: Visual confirmation of exact matches (✓/✗)

## Use Cases

Interpolated lookup tables are useful for:

- Feedforward control (e.g., voltage vs. velocity)
- Nonlinear mappings (e.g., sensor calibration)
- Shooter velocity tables (distance to power)
- Elevator position to voltage mappings
- Any scenario requiring smooth value lookups

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
