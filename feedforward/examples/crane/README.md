# Crane Control Example

This example demonstrates feedforward control applied to a crane or gantry system, showing how to manage payload swing and positioning.

## What This Example Shows

- Crane position and height control
- Payload swing dynamics and damping
- Cable tension management
- Smooth motion to prevent swinging
- Load stability during movement
- Realistic crane physics simulation

## Running the Example

```bash
cd feedforward/examples/crane
go run main.go
```

## Key Learning Points

### Crane Control Challenges

The example demonstrates:

- **Swing Dynamics**: Payload swings like a pendulum
- **Coupling**: Horizontal and vertical motion affect swing
- **Oscillations**: Swing period depends on cable length
- **Damping**: Slow motion reduces swing amplitude
- **Landing Stability**: Controlled descent for safe placement

### Payload Physics

Key concepts:

- **Pendulum Motion**: Payload swings at natural frequency
- **Energy Transfer**: Jerky motions excite oscillations
- **Resonance**: Moving at certain speeds amplifies swing
- **Damping**: Friction and air resistance slowly reduce swing
- **Control Input**: Crane motion directly affects payload state

## Output Interpretation

The example displays:

- **Crane Position**: Horizontal location (X coordinate)
- **Crane Height**: Vertical position (Z coordinate)
- **Cable Length**: Distance from crane to payload
- **Payload Position**: Actual position of hanging load
- **Payload Swing**: Angular displacement from vertical
- **Cable Tension**: Force in suspension cables

## System Parameters

The example typically uses:

- **Crane Position**: Horizontal range of motion
- **Crane Height**: Vertical range (typically large)
- **Cable Length**: Distance to payload
- **Payload Mass**: Weight of suspended load
- **Gravity**: Constant downward acceleration
- **Air Resistance**: Damping of swing motion

## Further Exploration

Try modifying:

- `payloadMass` - Test with different loads
- `cableLength` - Shorter cables swing less
- Motion speed - Slower motion reduces swing
- Target positions - Different approach patterns
- Damping coefficient - Air resistance effect
- Gravity - Test on different planets

## Real-World Applications

Crane control is critical for:

- **Port Cranes**: Loading/unloading cargo ships
- **Construction Cranes**: Moving materials on building sites
- **Industrial Cranes**: Heavy manufacturing and assembly
- **Automated Warehouses**: Goods movement and storage
- **Shipyard Operations**: Moving large ship components
- **Utility Pole Work**: Lifting poles and equipment
- **Amusement Rides**: Swinging-based attractions
- **Rescue Operations**: Controlled lowering of personnel

## Related Examples

- `../basic/` - Simpler feedforward without swing
- `../elevator/` - Vertical motion without swing
- `../arm/` - Multi-axis control
- `../compare/` - Comparison with other control types

## Swing Control Techniques

### Anti-Sway Control

The primary goal is to minimize payload swing:

1. **Trajectory Planning**: Design crane paths to minimize swing excitation
2. **Velocity Limiting**: Slower motion reduces swing amplitude
3. **Smooth Acceleration**: Ramped acceleration prevents jerky motion
4. **Input Shaping**: Pre-computed command sequences that inherently cancel swing
5. **Active Damping**: Using movement to damp oscillations

### Optimal Motion Profiles

Research shows that specific acceleration profiles minimize swing:

- **S-curve profiles**: Smooth acceleration, cruise, deceleration
- **Fuel-optimal**: Minimum-time motion that still controls swing
- **Input-shaped**: Pre-computed sequence that cancels resonant frequency

## Implementation Challenges

**Nonlinear Dynamics**: System behavior depends on swing amplitude
**Long Cable Effects**: Swing period depends on cable length
**Load Uncertainty**: Unknown payload mass affects swing frequency
**Structural Compliance**: Crane structure deflection complicates motion
**Wind Disturbances**: External forces excite swing motion

## Control Performance Metrics

- **Positioning Accuracy**: How close to target position
- **Swing Damping**: Time to settle oscillations
- **Energy Efficiency**: Smooth paths use less power
- **Safety**: Swing limits to prevent collision
- **Throughput**: Time to complete movements

## Advanced Concepts

- **Optimal Control**: Minimize time and swing simultaneously
- **State Feedback**: Using swing angle to control motion
- **Adaptive Control**: Learning optimal parameters for different loads
- **Model Predictive Control**: Predicting swing to prevent excitation
- **Reinforcement Learning**: Learning optimal trajectories

## Physics of Swing Motion

A hanging payload acts like an inverted pendulum:

- **Natural Frequency**: Depends on cable length (shorter = faster)
- **Period**: T = 2π√(L/g), where L is cable length
- **Resonance**: Moving at frequency matching period amplifies swing
- **Damping**: Energy dissipation slows oscillations

## Safety Considerations

- **Swing Limits**: Prevent excessive payload tilt
- **Anti-Collision**: Monitor space around crane and payload
- **Emergency Stop**: Instant stop capability
- **Overload Protection**: Prevent exceeding cable strength
- **Visibility**: Spotter systems for safe operation
- **Weather Limits**: Wind restrictions for safety

## Performance Expectations

Well-designed crane control achieves:

- Positioning accuracy of 10-30 cm
- Swing amplitude < 5° at destination
- Settling time < 30 seconds for typical moves
- Smooth, predictable motion
- Safe operation under all conditions
