package feedback

type NoFeedback struct{}

// Calculate always returns zero, indicating no feedback correction
func (nf *NoFeedback) Calculate(setpoint, measurement float64) float64 {
	return 0.0
}
