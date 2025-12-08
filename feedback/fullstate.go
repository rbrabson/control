package feedback

type Values []float64

// Full State feedback is an approach where we perform simultaneous feedback on each state
// (position, velocity, etc) of our system in parallel. This type of controller works especially
// well with motion profiles.
type FullStateFeedback struct {
	gain Values
}

// New creates a new FullStateFeedback controller with the specified gain values.
func New(gain Values) *FullStateFeedback {
	return &FullStateFeedback{
		gain: gain,
	}
}

// Calculate computes the control output based on the full state feedback
func (fsf *FullStateFeedback) Calculate(setpoint, measurement Values) (float64, error) {
	errorVec, err := minus(setpoint, measurement)
	if err != nil {
		return 0, err
	}
	return product(errorVec, fsf.gain)
}

// minus is a helper function that subtracts two vectors element-wise
func minus(a, b Values) (Values, error) {
	if len(a) != len(b) {
		return nil, ErrSlicessMustBeSameLength
	}
	result := make(Values, len(a))
	for i := range len(a) {
		result[i] = a[i] - b[i]
	}
	return result, nil
}

// product is a helper function that computes the dot product of two vectors
func product(a, b Values) (float64, error) {
	if len(a) != len(b) {
		return 0, ErrSlicessMustBeSameLength
	}
	sum := 0.0
	for i := range len(a) {
		sum += a[i] * b[i]
	}
	return sum, nil
}
