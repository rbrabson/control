# Filter Comparison Example

This example demonstrates the differences between **Kalman Filter** and **LowPass Filter** when used with PID controllers in noisy environments.

## What This Example Shows

- Side-by-side comparison of three controller configurations:
  1. **Kalman Filter**: Advanced state estimation with adaptive filtering
  2. **LowPass Filter**: Simple, computationally efficient filtering
  3. **No Filter**: Direct measurement for response speed reference

- Performance metrics including:
  - Total tracking error
  - Maximum error magnitude
  - Output variance (smoothness)
  - Output standard deviation
  - Final state convergence

## Running the Example

```bash
cd pid/examples/filter_comparison
go run main.go
```

## Key Learning Points

### When to Use Kalman Filter

- You have a good model of system dynamics
- You need adaptive filtering that learns system behavior
- Maximum accuracy is important
- You can afford slightly higher computational cost

### When to Use LowPass Filter

- Simple, predictable filtering behavior needed
- Computational efficiency is critical
- System dynamics are unknown or complex
- You want a fast, straightforward solution

### When to Use No Filter

- Measurement noise is minimal
- Response speed is critical (no filtering lag)
- You prefer handling noise at a higher level

## Output Interpretation

The example displays:

- **Total Absolute Error**: Sum of absolute errors over the simulation (lower is better)
- **Maximum Error**: Peak error during the simulation
- **Output Variance**: Measure of control output smoothness (lower is smoother)
- **Output Standard Deviation**: Consistency of control effort

## System Parameters

The example uses:

- Time constant: 1.0 second (first-order system)
- Noise magnitude: 2.0 (realistic sensor noise)
- Simulation duration: 10.0 seconds
- Sample time: 0.01 seconds

You can modify these parameters in the code to test different scenarios.

## Further Exploration

Try modifying:

- `noiseMagnitude` - Increase to see how filters handle more noise
- `timeConstant` - Change system dynamics
- Filter parameters (gain for LowPass, q/r for Kalman)
- PID gains to see impact on performance

## Related Examples

- `../dampening/` - Demonstrates PID dampening features
- `../motor_speed/` - Real-world motor control example
- `../temperature_control/` - Thermal system control with filtering
