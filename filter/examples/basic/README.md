# Kalman Filter Signal Estimation Example

This example demonstrates the Kalman filter, an optimal recursive algorithm for state estimation from noisy measurements.

## What This Example Shows

- Kalman filter fundamentals for signal smoothing
- State estimation from noisy measurements
- Optimal balance between measurement and model
- First-order system filtering application
- Variance and error metrics
- Comparison with other filtering approaches

## Running the Example

```bash
cd filter/examples/basic
go run main.go
```

## Key Learning Points

### Kalman Filter Principles

The example demonstrates:

- **Recursive Algorithm**: Processes one measurement at a time
- **Optimal**: Minimizes error variance (statistically optimal)
- **Adaptive**: Balances measurement noise vs. model uncertainty
- **State Estimation**: Estimates true state from noisy data
- **Linear Systems**: Classic formulation for linear dynamics

### Signal Estimation Problem

Given:

- Noisy measurements of a signal
- Model of how signal evolves
- Statistics of noise and model error

The Kalman filter computes:

- Best estimate of true signal
- Uncertainty in the estimate
- Prediction of future state

## Output Interpretation

The example displays:

- **Time**: Sample number or elapsed time
- **True Signal**: Actual value (if known)
- **Measurement**: Noisy observation
- **Kalman Estimate**: Filtered estimate
- **Error**: Difference from true value
- **Variance**: Estimate uncertainty

## System Parameters

The example typically uses:

- **Process Model**: System dynamics parameters (e.g., time constant)
- **Noise Covariance Q**: Uncertainty in process model
- **Measurement Noise R**: Variance of sensor noise
- **Initial Conditions**: Starting state and uncertainty

## Further Exploration

Try modifying:

- `Q` value - Increase to trust measurements more
- `R` value - Increase to trust model more
- Initial state and covariance
- System dynamics (time constant)
- Measurement noise characteristics
- Initial uncertainty estimate

## Real-World Applications

Kalman filters are used in:

- **GPS Navigation**: Combining noisy position measurements
- **Missile Guidance**: Tracking moving targets
- **Robot Localization**: Position and orientation estimation
- **Stock Price Prediction**: Financial time series
- **Spacecraft Attitude**: Orientation estimation
- **Aircraft Navigation**: Integrated navigation systems
- **Seismic Monitoring**: Earthquake detection and location
- **Medical Signals**: ECG and EEG filtering
- **Autonomous Vehicles**: Sensor fusion

## Kalman Filter Equations

### Prediction Step

```math
\hat{x}_{k|k-1} = A \hat{x}_{k-1|k-1} + B u_k
P_{k|k-1} = A P_{k-1|k-1} A^T + Q
```

### Update Step

```math
K_k = P_{k|k-1} H^T (H P_{k|k-1} H^T + R)^{-1}
\hat{x}_{k|k} = \hat{x}_{k|k-1} + K_k (z_k - H \hat{x}_{k|k-1})
P_{k|k} = (I - K_k H) P_{k|k-1}
```

Where:

- $\hat{x}$: State estimate
- $P$: Estimation covariance
- $K$: Kalman gain
- $z$: Measurement
- $A, B, H$: System matrices
- $Q, R$: Process and measurement noise covariances

## Advantages

- **Optimal**: Statistically optimal for linear systems
- **Adaptive**: Automatically balances sources of uncertainty
- **Recursive**: Memory-efficient, constant computation
- **Proven**: Decades of real-world validation
- **Flexible**: Handles multi-dimensional problems

## Limitations

- **Linear Systems Only**: Non-linear version (EKF, UKF) more complex
- **Model Required**: Must know system dynamics
- **Gaussian Noise**: Assumes noise is normally distributed
- **Noise Statistics**: Must estimate Q and R values
- **Tuning Required**: Gain values affect performance

## Tuning Q and R

### Q Matrix (Process Noise)

- Higher Q: Trust model less, weight measurements more
- Lower Q: Trust model more, smooth measurements heavily
- Diagonal elements: Uncertainty per state variable

### R Matrix (Measurement Noise)

- Higher R: Measurement sensor is noisy, discount measurements
- Lower R: Sensor is accurate, weight measurements heavily
- Scalar or diagonal for single-variable systems

### Tuning Procedure

1. **Start Conservative**: High R, low Q
2. **Observe Response**: Filter output should be smooth but responsive
3. **Adjust Q**: Increase if filter lags behind true signal
4. **Adjust R**: Decrease if filter follows noise too much
5. **Validate**: Test with known signals

## Extensions of Kalman Filtering

### Extended Kalman Filter (EKF)

- Handles non-linear systems
- Uses linearization around current estimate
- More complex but applicable to wider range of problems

### Unscented Kalman Filter (UKF)

- Better non-linear performance than EKF
- Uses carefully chosen sample points
- Avoids explicit Jacobian computation

### Ensemble Kalman Filter

- For high-dimensional systems
- Computational advantages for large problems
- Used in weather forecasting and geophysics

## Implementation Considerations

- **Numerical Stability**: Covariance matrix must remain positive definite
- **P Matrix Update**: Use Joseph form for guaranteed stability
- **Initialization**: Proper initial covariance is important
- **Square-Root Form**: Alternative formulation for better numerics
- **Adaptive Forms**: Automatically estimate noise statistics

## Performance Expectations

A well-tuned Kalman filter:

- Estimates true signal better than raw measurements
- Produces smooth output without lag
- Adapts to changing noise characteristics
- Maintains consistent performance over time
- Provides realistic uncertainty estimates

## Comparison with Other Filters

| Filter | Linear | Non-linear | Optimal | Complexity |
|--------|--------|-----------|---------|-----------|
| Kalman | Yes | No | Yes | Low |
| EKF | Yes | Yes | No | Medium |
| UKF | Yes | Yes | Better | Medium |
| Moving Avg | Yes | Yes | No | Very Low |
| Low-Pass | Yes | Yes | No | Very Low |

## Use Kalman Filter When

- You have a mathematical model of the system
- Measurements are noisy
- You need optimal estimation
- Real-time computation is needed
- Multi-dimensional state estimation required

## Use Simpler Filter When

- No system model available
- Simplicity is critical
- Computational resources limited
- Performance of simple filter is adequate
- Development time is limited
