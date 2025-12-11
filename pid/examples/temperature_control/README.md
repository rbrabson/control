# Temperature Control Example

This example demonstrates PID control applied to a temperature regulation system, commonly used in industrial and laboratory applications.

## What This Example Shows

- PID control of system temperature
- Heating and cooling dynamics
- Temperature setpoint tracking with overshoot control
- Handling of thermal system lag
- Reaching and maintaining temperature targets

## Running the Example

```bash
cd pid/examples/temperature_control
go run main.go
```

## Key Learning Points

### Temperature Control Characteristics

The example demonstrates:

- **Thermal Inertia**: Slow response to heating/cooling inputs (high time constant)
- **Asymmetry**: Heating and cooling may have different rates
- **Overshoot Prevention**: Using Kd term to avoid excessive overshoot
- **Steady-State Precision**: Integral term ensures target temperature is reached
- **Environmental Coupling**: Heat loss to environment affects system

### Thermal Dynamics

Thermal systems have:

- Large time constants (slower than mechanical systems)
- Significant overshoot tendency
- Relatively stable dynamics
- Long settling times

## Output Interpretation

The example displays:

- **Target Temperature**: Desired temperature setpoint
- **Current Temperature**: Actual system temperature
- **Time**: Elapsed simulation time
- **Error**: Temperature deviation from target
- **Heating Power**: Control signal (0-100% heater output)

## System Parameters

The example typically uses:

- **Target Temperature**: 100°C (or normalized to 100 units)
- **Time Constant**: 5-10 seconds (slow thermal response)
- **Initial Temperature**: Room temperature (20-25°C)
- **PID Gains**: Conservative to prevent overshoot
- **Control Update Rate**: 100ms or slower (thermal systems are slow)

## Further Exploration

Try modifying:

- `targetTemperature` - Test different temperature setpoints
- `timeConstant` - Simulate different thermal masses
- `ambientTemperature` - Include environmental heat loss
- PID gains - Balance response speed vs. overshoot
- `coolingRate` - Add active cooling effects

## Real-World Applications

This control technique is used in:

- Laboratory ovens and furnaces
- Incubators and environmental chambers
- Refrigeration and HVAC systems
- Industrial process heating
- 3D printer nozzle/bed temperature
- Water heater thermostats
- Reflow ovens for electronics manufacturing

## Related Examples

- `../motor_speed/` - Control in mechanical domain
- `../basic_control_loop/` - Start here for fundamentals
- `../position_servo/` - Position control (faster dynamics)
- `../dampening/` - Advanced response shaping
- `../../interplut/examples/temperature/` - Non-linear temperature control

## Temperature Control Challenges

**Slow Convergence**: Thermal systems are inherently slow; high Kp can help but risk overshoot
**Overshoot**: Temperature controllers must avoid overshooting (especially with heating only)
**Oscillation Around Setpoint**: Tune Ki carefully to avoid sustained oscillations
**Hysteresis**: Real thermostats may have deadbands; this example ignores that for simplicity

## Tuning Guidelines for Thermal Systems

- Start with low Kp (temperature response is slow)
- Use moderate Ki to eliminate steady-state error
- Use higher Kd to prevent overshoot (derivative dominates in thermal systems)
- Expect settling times of 30-60 seconds even for well-tuned systems
