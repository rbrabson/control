# Elevator Control Example

This example demonstrates feedforward control applied to an elevator system, showing how to compensate for gravity in vertical motion control.

## What This Example Shows

- Gravity compensation feedforward for vertical systems
- Weight and load handling
- Smooth acceleration and deceleration
- Elevator position tracking with realistic physics
- Tension and force analysis
- Energy-efficient control strategies

## Running the Example

```bash
cd feedforward/examples/elevator
go run main.go
```

## Key Learning Points

### Gravity Compensation

The example demonstrates:

- **Constant Force**: Gravity provides constant downward force (mg)
- **Feedforward Compensation**: Pre-computing required elevator tension
- **Load Changes**: Handling passenger weight variations
- **Energy Efficiency**: Minimizing unnecessary forces
- **Smooth Motion**: Avoiding jerky acceleration/deceleration

### Elevator Physics

Key concepts:

- **Weight**: Total mass including cabin and passengers
- **Gravity**: Constant acceleration (9.8 m/s²)
- **Tension**: Force from cables/motors (must overcome weight)
- **Acceleration**: Change in velocity (smooth ramping preferred)

## Output Interpretation

The example displays:

- **Target Floor**: Desired destination
- **Current Position**: Elevator position in building
- **Target Height**: Height corresponding to floor
- **Velocity**: Elevator speed (positive = upward)
- **Required Tension**: Force needed from motor
- **Error**: Position tracking error

## System Parameters

The example typically uses:

- **Building Height**: Total height or number of floors
- **Floor Spacing**: Distance between floors
- **Cabin Mass**: Empty elevator weight
- **Max Passengers**: Weight capacity
- **Max Acceleration**: Comfort and safety limits
- **Max Velocity**: Safety speed limit

## Further Exploration

Try modifying:

- `passengerWeight` - Simulate loading/unloading
- `maxAcceleration` - Adjust for comfort vs. speed trade-off
- `maxVelocity` - Safety limits and efficiency
- Gravity value - Simulate different planetary conditions
- `floorHeight` - Different building configurations

## Real-World Applications

This control technique is used in:

- **Elevators**: Building vertical transportation
- **Cranes**: Lifting and lowering loads
- **Hoists**: Industrial material handling
- **Cable Cars**: Inclined plane transportation
- **Ski Lifts**: Passenger transport up slopes
- **Parking Lifts**: Automated parking systems
- **Stage Equipment**: Theater rigging systems
- **Offshore Equipment**: Winches and hoists

## Related Examples

- `../basic/` - Simpler feedforward without gravity
- `../crane/` - Similar but more complex with swing dynamics
- `../arm/` - Multi-axis robotic systems
- `../compare/` - Comparison with other control types
- `../../pid/examples/temperature_control/` - Different domain with similar feedforward

## Gravity Compensation Formula

For a vertical system:

- Constant Feedforward = Weight × g
- Dynamic Feedforward = Mass × required_acceleration
- Total = Weight × g + Mass × required_acceleration

This ensures the elevator can accelerate smoothly despite gravity.

## Safety Considerations

- **Emergency Stop**: Instant stopping and holding
- **Overspeed Limits**: Preventing dangerous speeds
- **Backlash**: Smooth tension reversal to avoid shock
- **Counterweight**: Often used to reduce motor load
- **Redundancy**: Multiple safety systems

## Comfort Requirements

Good elevator control provides:

- Smooth acceleration (avoid jerky starts)
- Constant velocity cruise phase
- Smooth deceleration (no sudden stops)
- No oscillation around floor positions
- Quiet operation

## Implementation Challenges

**Load Variation**: Passenger weight changes compensation requirement
**Friction**: Varies with load, affects compensation accuracy
**Cable Stretch**: Elastic effects in tall buildings
**Door Timing**: Coordination with floor doors
**Energy Recovery**: Regenerating power during descent

## Performance Expectations

Well-designed elevator control achieves:

- Smooth motion without noticeable jerks
- Accurate floor stopping (± few cm)
- Fast response to call buttons
- Comfortable acceleration limits (< 1 m/s²)
- Efficient energy use (especially with regeneration)
