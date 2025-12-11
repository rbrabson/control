# Controller Comparison Example

This example provides a comprehensive comparison of different control strategies (PID, Feedforward, and Feedback) applied to the same system, helping you understand the trade-offs and best use cases.

## What This Example Shows

- Side-by-side comparison of multiple control approaches
- Pure feedback (PID) control performance
- Pure feedforward control performance
- Combined feedback and feedforward control
- Sinusoidal trajectory tracking test
- Step response and disturbance rejection
- Performance metrics for each strategy

## Running the Example

```bash
cd feedforward/examples/compare
go run main.go
```

## Key Learning Points

### Control Strategy Comparison

The example demonstrates three approaches:

1. **Pure Feedback (PID)**
   - Reacts to errors
   - No model required
   - Handles disturbances and model errors
   - Slower to respond to setpoint changes
   - Most stable and robust

2. **Pure Feedforward**
   - Predicts required output
   - Requires accurate model
   - Excellent tracking of known trajectories
   - No disturbance rejection
   - Can be unstable if model is wrong

3. **Combined (Feedback + Feedforward)**
   - Best of both worlds
   - Uses feedforward for main action
   - Uses feedback to correct errors
   - Most common in industry
   - Requires both model and feedback

## Output Interpretation

The example displays:

- **Target Signal**: Reference trajectory
- **PID Output**: Response from feedback-only control
- **Feedforward Output**: Response from model-based control
- **Combined Output**: Response from hybrid approach
- **Tracking Error**: Error magnitude for each strategy
- **Performance Metrics**: Comparison statistics

## Test Scenarios

Typically includes:

1. **Sinusoidal Trajectory**: Smooth continuous motion test
2. **Step Response**: Sudden setpoint change test
3. **Disturbance Rejection**: How well each handles sudden disturbances
4. **Steady-State Tracking**: Long-term error in continuous operation

## Further Exploration

Try modifying:

- Trajectory shape - Different signal patterns
- Disturbance magnitude - Test robustness
- Model accuracy - Introduce errors in feedforward model
- PID gains - Tune for better performance
- System parameters - Different dynamics

## Real-World Decision Tree

**Use Pure Feedback (PID) when:**

- System dynamics are unknown or complex
- Disturbances are significant
- Robustness is more important than speed
- Implementation complexity should be minimized
- Tuning resources are limited

**Use Pure Feedforward when:**

- System model is well-known and accurate
- No significant unmeasured disturbances
- Reference signals are predictable
- Speed is critical
- System is inherently open-loop stable

**Use Combined Control when:**

- You have a good system model
- Disturbances are present but manageable
- Both speed and robustness are important
- You can afford the additional complexity
- Industrial robustness is required (most common case)

## Related Examples

- `../basic/` - Focused on feedforward only
- `../elevator/` - Gravity compensation without feedback
- `../arm/` - Multi-axis control combining approaches
- `../crane/` - Complex system with coupling
- `../../pid/examples/` - Pure feedback control studies
- `../../feedback/examples/` - Advanced feedback techniques

## Control Design Framework

### Step 1: Assess System Knowledge

- How well do you understand the system?
- Can you model it accurately?
- Are disturbances predictable?

### Step 2: Choose Strategy

- Good model + low disturbance → Feedforward
- Poor model + significant disturbance → Feedback
- Both → Combined (most cases)

### Step 3: Implement and Tune

- Feedback: Tune PID gains experimentally
- Feedforward: Validate model accuracy
- Combined: Balance feedforward compensation and feedback correction

### Step 4: Validate Performance

- Track target signals
- Handle disturbances
- Verify stability
- Optimize gains

## Performance Comparison Metrics

| Metric | PID | Feedforward | Combined |
|--------|-----|-------------|----------|
| Model Required | No | Yes | Yes |
| Disturbance Rejection | Good | Poor | Good |
| Tracking Accuracy | Moderate | Excellent | Excellent |
| Robustness | High | Low | High |
| Tuning Effort | High | Medium | Medium |
| Complexity | Low | Medium | Medium |
| Computation | Low | Low-Medium | Low-Medium |

## Implementation Considerations

### Model Accuracy

- Small errors in feedforward model are usually acceptable
- Feedback loop corrects for small model errors
- Large model errors make feedforward alone unreliable

### Disturbance Characterization

- Known disturbances: Include in feedforward
- Unmeasured disturbances: Rely on feedback
- Stochastic disturbances: Feedback essential

### Gain Tuning

- Feedback gains: Empirical tuning with step responses
- Feedforward gains: Based on system model
- Combined: Reduce feedback gains slightly when feedforward helps

## Common Pitfalls

**Trusting Feedforward Too Much**: Model errors accumulate over time
**PID with Unknown Disturbances**: Control effort increases without effect
**Poor Model in Feedforward**: Worse than feedback-only control
**Ignoring Feedback**: Can't handle real-world variations
**Over-Tuning Feedback**: Fighting against good feedforward signal

## Advanced Topics

- **Two-Degree-of-Freedom Control**: Separate setpoint and disturbance response
- **Adaptive Control**: Learning and adjusting based on performance
- **Robust Control**: Designing for worst-case model uncertainties
- **Optimal Control**: Minimizing energy or control effort
- **Learning Control**: Using neural networks for model and control

## Case Studies

### Elevator System

- Gravity provides primary feedforward
- Feedback handles load variation and friction
- Result: Smooth, energy-efficient operation

### Robotic Arm

- Inverse dynamics for main feedforward
- Feedback for precision and error correction
- Result: Fast, accurate positioning

### Aircraft Autopilot

- Commands mixed with feedforward from auto-throttle
- Feedback from sensors for stability
- Result: Smooth flight and fuel efficiency

## Practical Selection Guide

**Elevator**: 80% feedforward, 20% feedback
**Robotic Arm**: 70% feedforward, 30% feedback
**Temperature Control**: 20% feedforward, 80% feedback
**Aircraft**: 60% feedforward, 40% feedback
**Unknown System**: 0% feedforward, 100% feedback (start here)
