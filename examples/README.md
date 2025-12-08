# Control Systems Examples

This directory contains examples demonstrating the control systems library. The examples are organized across different control packages, each showcasing specific control techniques and applications.

## Overview of Available Examples

The control systems library includes examples in several locations:

### PID Control Examples (`../pid/examples/`)

Advanced PID controller implementations with various features:

- **`basic_control_loop/`** - Fundamental PID usage with a simple first-order system
- **`motor_speed/`** - Motor speed control with realistic dynamics and disturbances  
- **`position_servo/`** - Position control system with error tracking
- **`temperature_control/`** - Temperature regulation with thermal dynamics

### Feedback Control Examples (`../feedback/examples/`)

State-space and vector-based feedback control systems:

- **`feedback_control/`** - Full-state feedback controller implementation

### Feedforward Control Examples (`../feedforward/examples/`)

Predictive control for improved performance:

- **`basic/`** - Simple motor control with velocity/acceleration compensation
- **`elevator/`** - Elevator system with gravity compensation
- **`arm/`** - Robotic arm with angle-dependent cosine compensation
- **`crane/`** - Complex crane control with combined compensations
- **`compare/`** - Side-by-side comparison of all feedforward controller types

### Combined Control Examples (This Directory)

Advanced examples combining multiple control techniques:

- **`combined_control/`** - Integration of PID and feedback control methods

## Running the Examples

Each example can be run independently from its directory:

```bash
# Navigate to any example directory and run
cd basic_control_loop  # or any other example directory
go run main.go         # or the appropriate main file
```

### Quick Navigation

```bash
# PID Examples
cd ../pid/examples/basic_control_loop && go run main.go

# Feedback Examples  
cd ../feedback/examples/feedback_control && go run main.go

# Feedforward Examples
cd ../feedforward/examples/basic && go run .

# Combined Examples
cd combined_control && go run main.go
```

## Example Categories

### **Learning Path: Beginner**

Start with these examples to understand basic concepts:

1. `../pid/examples/basic_control_loop/` - Learn PID fundamentals
2. `../feedforward/examples/basic/` - Understand feedforward control
3. `../feedback/examples/feedback_control/` - Explore state feedback

### **Learning Path: Intermediate**

Build on the basics with realistic applications:

1. `../pid/examples/motor_speed/` - Practical PID tuning
2. `../feedforward/examples/elevator/` - Gravity compensation
3. `combined_control/` - Multiple control techniques

### **Learning Path: Advanced**

Complex systems and optimization:

1. `../feedforward/examples/crane/` - Multi-compensation systems
2. `../feedforward/examples/compare/` - Controller selection
3. `../pid/examples/temperature_control/` - Non-linear dynamics

## Key Control Concepts Demonstrated

### **PID Control**

- Proportional, Integral, Derivative action
- Output limiting and anti-windup
- Parameter tuning techniques
- Real-world disturbance handling

### **Feedforward Control**

- Velocity and acceleration compensation
- Gravity compensation for vertical systems
- Cosine compensation for rotating systems
- Combined compensation strategies

### **Feedback Control**

- Full-state feedback design
- Vector-based control systems
- Multi-input, multi-output systems

### **Combined Control**

- Feedforward + Feedback integration
- PID + State feedback combinations
- Performance comparison studies

## Performance Characteristics

The examples demonstrate varying performance characteristics:

- **PID Controllers**: ~50-100 ns/op (depends on features enabled)
- **Feedforward Controllers**: ~2-5 ns/op (ultra-fast calculations)
- **Feedback Controllers**: ~0.3-10 ns/op (scales with system dimension)

## Getting Started

1. **Install Go** (version 1.19 or later recommended)
2. **Clone the repository** and navigate to any example directory
3. **Run the example**: `go run main.go` (or `go run .`)
4. **Examine the output** to understand controller behavior
5. **Modify parameters** to see how they affect performance

## Tips for Exploration

- **Start Simple**: Begin with basic examples before moving to complex ones
- **Compare Outputs**: Run similar examples with different parameters
- **Read the Code**: Each example is well-commented and educational
- **Experiment**: Modify gain values and observe the effects
- **Combine Techniques**: Use insights from multiple examples in your own projects

## Integration Examples

The examples show how to integrate different control techniques:

```go
// Typical combined control structure
totalOutput := pidOutput + feedforwardOutput + feedbackOutput
```

This approach leverages the strengths of each control method:

- **PID**: Handles steady-state errors and disturbances
- **Feedforward**: Provides predictive compensation
- **Feedback**: Manages multi-variable system states

## Further Learning

Each example directory contains additional documentation and implementation details. The examples progress from basic concepts to advanced real-world applications, making them suitable for both learning and reference purposes.
