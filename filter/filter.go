package filter

// Filter defines the interface for signal filtering algorithms.
type Filter interface {
	// Estimate processes a measurement and returns the optimal state estimate.
	Estimate(measurement float64) float64
}
