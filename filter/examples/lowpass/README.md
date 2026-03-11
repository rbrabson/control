# Low-Pass Filter Basic Example

This example demonstrates first-order low-pass filtering behavior, matching the test cases in the Java implementation.

## What This Example Shows

- Basic first-order exponential smoothing
- Filter initialization with first measurement
- Gradual convergence to input value
- Step response behavior
- Simple, focused demonstration

## Running the Example

```bash
cd filter/examples/lowpass
go run main.go
```

## Key Learning Points

### First-Order Filter Equation

The recursive equation is:

```math
y_k = \alpha \cdot y_{k-1} + (1 - \alpha) \cdot z_k
```

Where:

- $y_k$: Filtered output at time k
- $z_k$: Measurement at time k  
- $\alpha$: Filter alpha (0 to 1)
- $y_{k-1}$: Previous filtered output

### Behavior

- **First call**: Initializes filter with input value (returns input)
- **Subsequent calls**: Applies exponential smoothing
- **Convergence**: Output gradually approaches input value
- **Alpha parameter**:
  - High α (e.g., 0.9): More smoothing, slower response
  - Low α (e.g., 0.1): Less smoothing, faster response
  - α = 0.5: Balanced (50% old, 50% new)

## Output Interpretation

The example shows:

1. **Initialization**: First estimate returns the input value exactly
2. **Filtering**: Second estimate shows the smoothing formula in action
3. **Convergence**: Multiple calls with constant input show gradual convergence
4. **Step Response**: Behavior when input suddenly changes

## Use Cases

Low-pass filters are used for:

- Sensor noise reduction
- Signal smoothing
- Derivative filtering (in PID control)
- Moving averages
- Data cleanup before processing

### Relationship to Cutoff Frequency

For a first-order filter:

```math
f_c = \frac{\alpha}{2\pi} f_s
```

Where $f_s$ is the sampling frequency.

### Time Constant

Effective time constant:

```math
\tau = \frac{1 - \alpha}{\alpha \cdot f_s}
```

Larger time constant = more smoothing.

## Advantages

- **Simplicity**: Single line of code
- **Speed**: Microsecond computation
- **Minimal Memory**: Requires only one previous value
- **Stable**: Never diverges
- **Intuitive**: Easy to understand and tune
- **Real-Time**: Works perfectly in embedded systems

## Limitations

- **Phase Lag**: Introduces delay proportional to smoothing
- **Not Optimal**: Kalman filter is better if model is available
- **Sharp Cutoff**: Not as steep as digital filters
- **Limited Rejection**: High attenuation only at very high frequencies
- **Fixed Gain**: Can't adapt to changing noise characteristics

## Design Procedure

### Step 1: Identify Noise Frequency

- Measure or estimate the dominant noise frequency
- Typically 2-10x faster than signal of interest

### Step 2: Choose Cutoff Frequency

- Set cutoff at least 2x signal frequency
- Set below noise frequency to attenuate noise
- Compromise between smoothing and lag

### Step 3: Calculate α

From cutoff frequency relationship:

```math
\alpha = \frac{2\pi f_c}{f_s + 2\pi f_c}
```

### Step 4: Validate

- Test with noisy measurement
- Verify smoothing is adequate
- Check that lag is acceptable
- Adjust α if needed

## Alternative Forms

### Discrete Time Constant Form

Specify desired time constant and sample rate:

```math
\alpha = \frac{f_s \cdot dt}{f_s \cdot dt + 1}
```

Where dt is desired time constant.

### Butterworth Form

For sharper frequency response:

```math
y_k = \frac{\alpha z_k + (1-\alpha)y_{k-1} + (1-\alpha)y_{k-2}}{2-\alpha}
```

More complex but better filtering.

## Cascaded Filters

Use multiple stages for sharper cutoff:

```text
Input → Filter 1 → Filter 2 → Filter 3 → Output
```

Effect:

- Each stage reduces noise more
- Lag increases proportionally
- Computation still very fast
- Good practical compromise

## Comparison with Other Approaches

| Method | Simplicity | Speed | Optimality | Lag |
|--------|-----------|-------|-----------|-----|
| No Filter | Very High | N/A | Poor | None |
| Moving Average | High | Fast | Fair | Moderate |
| Low-Pass (α=0.5) | Very High | Very Fast | Fair | Moderate |
| Low-Pass (α=0.1) | Very High | Very Fast | Fair | High |
| Kalman | Low | Fast | Optimal | Low |
| FIR Filter | Medium | Medium | Good | Moderate |
| IIR Filter | Medium | Medium | Good | Low |

## Real-World Tuning Examples

### High-Speed Motor (1000 RPM)

- Noise: 100 Hz vibration
- Signal bandwidth: 10 Hz
- Cutoff frequency: 20 Hz
- α ≈ 0.3

### Temperature Sensor (100 Hz sample)

- Noise: 50 Hz electrical
- Signal bandwidth: 0.1 Hz
- Cutoff frequency: 0.5 Hz
- α ≈ 0.001

### Position Feedback (50 Hz sample)

- Noise: Quantization at 5 Hz
- Signal bandwidth: 2 Hz
- Cutoff frequency: 5 Hz
- α ≈ 0.16

## Performance Expectations

A low-pass filter achieves:

- Smooth signal without chatter
- Reduced but not eliminated noise
- Minimal computational overhead
- Consistent behavior over time
- Easy integration into existing systems

## When to Use Low-Pass Filter

**Choose low-pass when:**

- Simplicity is more important than optimality
- Computational resources are limited
- No system model is available
- Real-time response with minimal lag needed
- Embedded system with limited memory

**Consider Kalman filter when:**

- System model is available
- Optimal performance is needed
- Faster convergence is important
- Complex sensor fusion required

## Advanced Topics

- **Adaptive Filtering**: Auto-adjust α based on signal characteristics
- **Notch Filters**: Target specific noise frequencies
- **Chebyshev/Butterworth**: Higher-order designs for sharper cutoff
- **Digital Signal Processing**: Advanced filter design theory
- **State-Space Form**: Equivalent higher-order formulations
