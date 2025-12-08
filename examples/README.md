# PID Controller Examples

This directory contains complete working examples demonstrating various uses of the PID controller library.

## Running the Examples

Each example is in its own subdirectory and can be run independently:

```bash
# Basic control loop example
cd basic_control_loop
go run basic_control_loop.go

# Motor speed control example  
cd motor_speed
go run main.go

# Temperature control example
cd temperature_control
go run main.go

# Position servo control example
cd position_servo  
go run main.go
```

## Example Descriptions

### Basic Control Loop (`basic_control_loop.go`)

Demonstrates fundamental PID usage with a simple simulated first-order system:

- Basic PID controller creation and configuration
- Output limiting
- Simple control loop with timing
- System response monitoring

### Motor Speed Control (`motor_speed/main.go`)

Shows advanced PID features for motor control applications:

- Integral sum limiting to prevent motor saturation
- Stability threshold to reduce overshoot during rapid changes
- Derivative filtering for encoder noise reduction
- Realistic motor dynamics simulation with noise
- Multiple setpoint changes demonstrating response characteristics

### Temperature Control (`temperature_control/main.go`)

Demonstrates thermal system control with environmental factors:

- Feed-forward control for ambient temperature compensation
- Integral reset on zero crossover to prevent overshoot
- Derivative filtering for temperature sensor noise
- Realistic thermal dynamics with heat capacity and loss
- Ambient temperature disturbances

### Position Servo Control (`position_servo/main.go`)

Illustrates precise positioning applications:

- Integral windup protection for large position errors
- Stability threshold during high-speed movement
- Encoder quantization effects
- Realistic servo dynamics with backlash and friction
- Multiple target positions with accuracy analysis

## Key Learning Points

Each example demonstrates different aspects of PID control:

1. **Basic Usage**: How to create, configure, and use a PID controller
2. **Advanced Features**: When and how to use integral reset, stability thresholds, and filtering
3. **Real-World Considerations**: Noise, saturation, mechanical effects, and disturbances
4. **Performance Tuning**: Different gain values for different system characteristics
5. **System Analysis**: Monitoring controller performance and system response

## Modifying the Examples

You can experiment with the examples by:

- Changing PID gains to see effects on response
- Modifying system parameters (time constants, noise levels, etc.)
- Adding or removing advanced PID features
- Implementing different setpoint profiles
- Adding disturbances to test robustness

## Example Output

Each example produces formatted output showing:

- Time progression
- Setpoint and measured values
- Error and control output
- System state information
- Final performance metrics

This helps visualize controller performance and understand system behavior under different conditions.
