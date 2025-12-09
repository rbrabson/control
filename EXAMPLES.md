# Control Systems Library - Complete Examples Guide

This document provides a comprehensive overview of all examples included in the Control Systems Library. The examples are organized by control type and demonstrate practical applications ranging from basic concepts to advanced real-world scenarios.

## Overview

The library includes **14 complete examples** across five control packages:

- **4 PID Controller Examples** - Classic feedback control with advanced features
- **1 Feedback Controller Example** - Multi-dimensional state feedback control  
- **5 Feedforward Controller Examples** - Predictive control with various compensation strategies
- **2 Motion Profile Examples** - Smooth trapezoidal motion generation
- **2 InterpLUT Examples** - Interpolating lookup tables with cubic splines

## Quick Start Guide

### Running Individual Examples

Each example is self-contained and can be run independently:

```bash
# Navigate to any example directory
cd pid/examples/basic_control_loop
go run main.go

# Or use the current directory approach
cd feedforward/examples/basic
go run .
```

### Running All Examples by Type

```bash
# PID Examples
for dir in pid/examples/*/; do (cd "$dir" && echo "=== $(basename "$dir") ===" && go run main.go); done

# Feedforward Examples  
for dir in feedforward/examples/*/; do (cd "$dir" && echo "=== $(basename "$dir") ===" && go run .); done

# Feedback Examples
for dir in feedback/examples/*/; do (cd "$dir" && echo "=== $(basename "$dir") ===" && go run main.go); done

# InterpLUT Examples
for dir in interplut/examples/*/; do (cd "$dir" && echo "=== $(basename "$dir") ===" && go run main.go); done
```

---

## PID Controller Examples

### 1. Basic Control Loop (`pid/examples/basic_control_loop/`)

**Purpose**: Fundamental PID usage with a simple simulated system

**Key Features**:

- Basic PID controller creation and configuration
- Output limiting with configurable bounds
- Simple first-order system simulation
- Real-time control loop with proper timing
- Performance monitoring and error tracking

**What You'll Learn**:

- PID controller initialization
- Basic control loop structure
- Output limiting importance
- System response characteristics
- Error convergence patterns

**Sample Output**:

``` bash
Time: 0.100s, Setpoint: 10.00, Position: 0.32, Error: 9.68, Output: 9.68
Time: 0.200s, Setpoint: 10.00, Position: 1.22, Error: 8.78, Output: 8.78
...
Final Error: 0.02, Settling Time: ~2.5s
```

**Run Command**: `cd pid/examples/basic_control_loop && go run main.go`

---

### 2. Motor Speed Control (`pid/examples/motor_speed/`)

**Purpose**: Realistic motor speed control with disturbances and load variations

**Key Features**:

- Advanced PID configuration with integral limits
- Load disturbance simulation
- Motor inertia and friction modeling
- Anti-windup protection demonstration
- Performance metrics calculation

**What You'll Learn**:

- Real-world PID tuning challenges
- Disturbance rejection techniques
- Integral windup prevention
- Motor dynamics modeling
- Performance optimization

**Sample Output**:

``` bash
RPM Target: 1500, Current: 1498, Error: 2, Load: 75%, Output: 85%
Load Change Detected: 50% → 90%
Disturbance Rejection: Recovered in 0.8s
```

**Run Command**: `cd pid/examples/motor_speed && go run main.go`

---

### 3. Temperature Control (`pid/examples/temperature_control/`)

**Purpose**: Temperature regulation system with thermal dynamics and environmental effects

**Key Features**:

- Thermal system modeling with heat transfer
- Environmental temperature variations
- Sensor noise simulation and filtering
- Setpoint ramping for gradual changes
- Energy efficiency monitoring

**What You'll Learn**:

- Temperature control challenges
- Thermal lag compensation
- Noise filtering techniques
- Energy-efficient control strategies
- Setpoint management

**Sample Output**:

``` bash
Target: 75.0°C, Current: 74.8°C, Ambient: 22°C, Heater: 45%
Environmental Change: 22°C → 18°C
Compensation Applied: +12% heater output
```

**Run Command**: `cd pid/examples/temperature_control && go run main.go`

---

### 4. Position Servo Control (`pid/examples/position_servo/`)

**Purpose**: Precise position control system with encoder feedback and trajectory following

**Key Features**:

- High-resolution encoder simulation
- Multi-point trajectory following
- Velocity and acceleration limiting
- Position accuracy verification
- Overshoot analysis

**What You'll Learn**:

- Precision position control
- Trajectory planning integration
- Encoder-based feedback
- Motion profiling
- Accuracy vs. speed trade-offs

**Sample Output**:

``` bash
Position: 45.23°, Target: 45.00°, Error: -0.23°, Velocity: 12.5°/s
Trajectory Point Reached: Position 2 of 5
Max Overshoot: 1.2°, Settling Time: 0.45s
```

**Run Command**: `cd pid/examples/position_servo && go run main.go`

---

## Feedback Controller Examples

### 5. Full State Feedback Control (`feedback/examples/feedback_control/`)

**Purpose**: Multi-dimensional state feedback control for complex systems

**Key Features**:

- Single-state feedback for position control
- Multi-state feedback for position/velocity control
- 3D system control (position/velocity/acceleration)
- Error handling for dimension mismatches
- Performance comparison between control strategies

**What You'll Learn**:

- State-space control concepts
- Multi-variable system control
- Gain matrix design principles
- System stability considerations
- Vector-based control calculations

**Sample Output**:

``` bash
2D System Control:
Target: [10.0, 0.0], Current: [2.5, 1.2], Output: 9.14

3D System Control: 
Target: [5.0, 0.0, 0.0], Current: [1.0, 0.5, 0.1], Output: 8.45

Performance: 0.31 ns/op (ultra-fast calculations)
```

**Run Command**: `cd feedback/examples/feedback_control && go run main.go`

---

## Feedforward Controller Examples

### 6. Basic Feedforward Control (`feedforward/examples/basic/`)

**Purpose**: Simple motor control with velocity and acceleration compensation

**Key Features**:

- Basic feedforward controller without options
- Velocity and acceleration feedforward terms
- Smooth sinusoidal motion profile simulation
- Transient response improvement demonstration
- Performance analysis vs. pure feedback

**What You'll Learn**:

- Feedforward control fundamentals
- Velocity/acceleration compensation
- Predictive control benefits
- Motion profile generation
- Feedforward vs. feedback comparison

**Sample Output**:

``` bash
Time     Position     Velocity     Accel        FF Output   
--------------------------------------------------------------
0.0      0.000        0.000        2.193        0.110       
0.1      0.011        0.219        2.181        0.284       
...
Analysis: Reduces steady-state error during motion
```

**Run Command**: `cd feedforward/examples/basic && go run .`

---

### 7. Elevator Control (`feedforward/examples/elevator/`)

**Purpose**: Elevator system with gravity compensation for vertical movement

**Key Features**:

- Gravity compensation using `WithGravityGain()` option
- Multi-floor elevator movement simulation
- S-curve motion profiles for passenger comfort
- Constant upward force demonstration
- Direction-independent performance

**What You'll Learn**:

- Gravity compensation techniques
- Vertical motion control
- S-curve profile benefits
- Constant force compensation
- Elevator-specific control challenges

**Sample Output**:

``` bash
Moving from Ground to Floor 2...
Time     Floor        Position     Velocity     FF Output   
------------------------------------------------------------------------
0.0      0.0m         0.00         0.000        9.880       
1.0      0.4m         0.44         0.875        10.930      
...
Base gravity compensation: 9.81 N
```

**Run Command**: `cd feedforward/examples/elevator && go run .`

---

### 8. Robotic Arm Control (`feedforward/examples/arm/`)

**Purpose**: Rotating robotic arm with angle-dependent cosine compensation

**Key Features**:

- Cosine compensation using `WithCosineGain()` option
- Full rotation simulation (0° to 180°)
- Gravitational torque variation demonstration
- Position-dependent force calculation
- Joint angle movement profiles

**What You'll Learn**:

- Angular system compensation
- Gravitational torque effects
- Cosine function applications
- Robotic arm dynamics
- Position-dependent control

**Sample Output**:

``` bash
Time     Position        Angle        Cos(θ)       FF Output   
--------------------------------------------------------------------------------
0.0      0.0°            0.0          1.000        2.520       
0.8      22.5°           0.4          0.924        2.526       
1.9      90.0°           1.6          0.001        0.037       
...
Maximum torque at horizontal (cos(0°) = 1.0)
Zero torque at vertical (cos(90°) = 0.0)
```

**Run Command**: `cd feedforward/examples/arm && go run .`

---

### 9. Crane Control (`feedforward/examples/crane/`)

**Purpose**: Complex crane boom control with combined gravity and cosine compensation

**Key Features**:

- Combined `WithGravityGain()` and `WithCosineGain()` options
- Heavy machinery simulation with realistic loads
- Multi-phase operation (raise, lower, position)
- Performance comparison: basic vs. combined feedforward
- Maximum control complexity demonstration

**What You'll Learn**:

- Multiple compensation strategies
- Complex machinery control
- Combined feedforward benefits
- Heavy load handling
- Control effort optimization

**Sample Output**:

``` bash
Time   Operation                 Angle°   FF Basic   FF Comb.   Diff      
------------------------------------------------------------------------------
0.0    Raise boom from horizontal 0.0     0.016      23.916     23.900    
2.0    At horizontal position     15.0    0.785      24.406     23.621    
4.0    Vertical position         60.0     -0.016     19.784     19.800    
...
Gravity component: 15.7 units, Max cosine component: 8.2 units
```

**Run Command**: `cd feedforward/examples/crane && go run .`

---

### 10. Controller Comparison (`feedforward/examples/compare/`)

**Purpose**: Side-by-side comparison of all feedforward controller types

**Key Features**:

- All controller types in one demonstration
- Performance analysis across different positions
- Educational comparison of each approach
- Output differences visualization
- Use case recommendations

**What You'll Learn**:

- Controller selection criteria
- Performance characteristics
- Application-specific benefits
- Comparative analysis techniques
- Decision-making guidelines

**Sample Output**:

``` bash
Time   Angle°   Basic    Gravity  Cosine   Combined None    
----------------------------------------------------------
0.0    90       1.974    6.974    1.974    6.974    0.000   
2.0    90       -1.974   3.026    -1.974   3.026    0.000   
4.0    90       1.974    6.974    1.974    6.974    0.000   

Use Cases:
Basic FF:     Simple systems without gravitational effects
Gravity FF:   Elevators, vertical lifts, load handling
Combined FF:  Cranes, construction equipment, complex robotics
```

**Run Command**: `cd feedforward/examples/compare && go run .`

---

## Performance Benchmarks

### Execution Speed Comparison

| Controller Type | Speed (ns/op) | Allocations | Use Case |
|-----------------|---------------|-------------|----------|
| PID Controller | ~64 | 0 | Feedback control loops |
| Feedback Control | ~0.3-10 | 0 | Multi-variable systems |
| Basic Feedforward | ~2.1 | 0 | Simple predictive control |
| Gravity Feedforward | ~2.1 | 0 | Vertical motion systems |
| Cosine Feedforward | ~5.2 | 0 | Rotating machinery |
| Combined Feedforward | ~5.2 | 0 | Complex multi-axis systems |

### Memory Efficiency

All controllers are designed for zero allocation during calculations:

- ✅ No memory allocations in calculation loops
- ✅ Minimal initialization overhead
- ✅ Optimal for real-time applications
- ✅ Suitable for embedded systems

---

## Learning Paths

### Beginner Path (Start Here)

1. **Basic Control Loop** - Learn PID fundamentals
2. **Basic Feedforward** - Understand predictive control
3. **Feedback Control** - Explore state-space methods

### Intermediate Path

1. **Motor Speed Control** - Real-world PID applications
2. **Elevator Control** - Gravity compensation techniques
3. **Temperature Control** - Environmental disturbances

### Advanced Path

1. **Position Servo** - Precision control systems
2. **Crane Control** - Multi-compensation strategies
3. **Controller Comparison** - Optimization and selection

### Integration Path

1. Study individual controller types
2. Compare performance characteristics
3. Design combined control strategies
4. Implement application-specific solutions

---

## Integration Examples

### Combined Control Architecture

```go
// Typical control system combining all three approaches
func CombinedControl(setpoint, current, velocity, acceleration, position float64) float64 {
    // PID for feedback control
    pidOutput := pidController.Calculate(setpoint, current, time.Now())
    
    // Feedforward for predictive control
    ffOutput := feedforwardController.Calculate(velocity, acceleration, position)
    
    // State feedback for multi-variable control  
    target := feedback.Values{setpoint, 0.0}
    state := feedback.Values{current, velocity}
    fbOutput, _ := feedbackController.Calculate(target, state)
    
    // Combine all control efforts
    totalOutput := pidOutput + ffOutput + fbOutput
    
    return totalOutput
}
```

### Application-Specific Combinations

**Robotic Arm Control**:

- PID: Joint position feedback
- Feedforward: Gravity + cosine compensation  
- State Feedback: Multi-joint coordination

**Elevator System**:

- PID: Floor positioning accuracy
- Feedforward: Gravity compensation
- State Feedback: Speed and acceleration control

**CNC Machine**:

- PID: Axis position control
- Feedforward: Acceleration compensation
- State Feedback: Multi-axis coordination

---

## Tips for Exploration

### Getting Started

1. **Start Simple**: Begin with basic examples before complex ones
2. **Understand Output**: Study the printed results carefully
3. **Modify Parameters**: Change gains and observe effects
4. **Compare Approaches**: Run similar examples with different controllers

### Experimentation Ideas

- Modify gain values in PID examples to see stability effects
- Change compensation parameters in feedforward examples
- Combine multiple controller outputs for hybrid control
- Add your own disturbance simulations

### Real-World Application

- Use examples as templates for your own systems
- Adapt controller parameters for your specific requirements
- Integrate multiple control techniques for robust performance
- Reference performance benchmarks for system design

### Troubleshooting

- Check controller parameter ranges (avoid negative values where inappropriate)
- Ensure proper time step sizing for stability
- Verify input/output scaling for your application
- Use performance metrics to validate controller effectiveness

---

## InterpLUT Examples

### 1. Basic Shooter Velocity (`interplut/examples/basic/`)

**Purpose**: Demonstrate smooth interpolation for robotics shooter velocity mapping

**Key Features Demonstrated**:

- Control point addition and LUT creation
- Smooth interpolation between discrete velocity settings
- Exact matches at control points
- Error handling for out-of-bounds inputs

**Application**: Robot shooter systems where velocity must vary smoothly with target distance

**Sample Output**:

```
Distance (m) | Velocity (%)
-------------|-------------
    1.5     |    0.272
    2.0     |    0.359
    2.5     |    0.456
    3.0     |    0.577
```

**Key Learning Points**:

- How to set up control points for non-linear mappings
- Understanding cubic spline interpolation benefits
- Proper error handling in robotics applications

### 2. Temperature Control (`interplut/examples/temperature/`)

**Purpose**: Non-linear temperature control with thermodynamic modeling

**Key Features Demonstrated**:

- Non-linear heater power curves
- Monotonicity preservation for stable control
- ASCII visualization of response curves
- Comparison with linear interpolation methods

**Application**: Industrial heating systems, thermal management, process control

**Sample Output**:

```
Temperature (°C) | Power (%) | Notes
-----------------|-----------|--------
      60.0       |   0.208   | Hot water
      85.0       |   0.362   | Coffee brewing
      120.0       |   0.558   | Steam generation
```

**Key Learning Points**:

- Modeling real-world non-linear systems
- Benefits of smooth interpolation over linear methods
- Visual representation of control curves

---

## Next Steps

After exploring these examples:

1. **Choose the Right Controller**: Match controller type to your application needs
2. **Tune Parameters**: Use examples as starting points for your own systems
3. **Combine Techniques**: Integrate multiple control approaches for optimal performance
4. **Build Applications**: Use the library in your own robotics or automation projects
5. **Contribute**: Add your own examples or improvements to the library

Each example includes detailed comments and explanations to help you understand both the theory and practical implementation of control systems in Go.
