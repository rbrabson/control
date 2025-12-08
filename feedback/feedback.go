package feedback

// Feedback is an interface for feedback controllers
type Feedback interface {
	Calculate(setpoint, measurement float64) float64
}
