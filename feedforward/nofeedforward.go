package feedforward

// NoFeedForward is a feedforward controller that always returns zero output. It is used to clarify
// cases where no feedforward action is desired.
type NoFeedForward struct{}

// Calculate always returns 0.0 for no feed-forward action.
func (n *NoFeedForward) Calculate(setpoint float64, velocity float64, acceleration float64) float64 {
	return 0.0
}
