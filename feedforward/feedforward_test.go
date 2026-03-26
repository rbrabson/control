package feedforward

import (
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		kS       float64
		kV       float64
		kA       float64
		opts     []Option
		expected *FeedForward
	}{
		{
			name: "basic feedforward without options",
			kS:   0.1,
			kV:   0.5,
			kA:   0.02,
			opts: nil,
			expected: &FeedForward{
				kS:   0.1,
				kV:   0.5,
				kA:   0.02,
				kCos: 0.0,
			},
		},
		// Removed gravity gain test case
		{
			name: "feedforward with cosine gain",
			kS:   0.15,
			kV:   0.6,
			kA:   0.03,
			opts: []Option{WithCosineGain(2.5)},
			expected: &FeedForward{
				kS:   0.15,
				kV:   0.6,
				kA:   0.03,
				kCos: 2.5,
			},
		},
		// Removed gravity+cosine gain test case
		{
			name: "feedforward with zero gains",
			kS:   0.0,
			kV:   0.0,
			kA:   0.0,
			opts: []Option{WithCosineGain(0.0)},
			expected: &FeedForward{
				kS:   0.0,
				kV:   0.0,
				kA:   0.0,
				kCos: 0.0,
			},
		},
		{
			name: "feedforward with negative gains",
			kS:   -9.93,
			kV:   -0.5,
			kA:   -0.02,
			opts: []Option{WithCosineGain(-2.5)},
			expected: &FeedForward{
				kS:   -9.93,
				kV:   -0.5,
				kA:   -0.02,
				kCos: -2.5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := New(tt.kS, tt.kV, tt.kA, tt.opts...)

			if ff.kS != tt.expected.kS {
				t.Errorf("kS = %f, expected %f", ff.kS, tt.expected.kS)
			}
			if ff.kV != tt.expected.kV {
				t.Errorf("kV = %f, expected %f", ff.kV, tt.expected.kV)
			}
			if ff.kA != tt.expected.kA {
				t.Errorf("kA = %f, expected %f", ff.kA, tt.expected.kA)
			}
			if ff.kCos != tt.expected.kCos {
				t.Errorf("kCos = %f, expected %f", ff.kCos, tt.expected.kCos)
			}
		})
	}
}

// TestWithGravityGain removed: WithGravityGain is no longer supported

func TestWithCosineGain(t *testing.T) {
	tests := []struct {
		name     string
		kCos     float64
		expected float64
	}{
		{
			name:     "positive cosine gain",
			kCos:     2.5,
			expected: 2.5,
		},
		{
			name:     "negative cosine gain",
			kCos:     -2.5,
			expected: -2.5,
		},
		{
			name:     "zero cosine gain",
			kCos:     0.0,
			expected: 0.0,
		},
		{
			name:     "small cosine gain",
			kCos:     0.01,
			expected: 0.01,
		},
		{
			name:     "large cosine gain",
			kCos:     50.0,
			expected: 50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := New(0.1, 0.5, 0.02, WithCosineGain(tt.kCos))
			if ff.kCos != tt.expected {
				t.Errorf("kCos = %f, expected %f", ff.kCos, tt.expected)
			}
		})
	}
}

func TestCalculate(t *testing.T) {
	const tolerance = 1e-9

	tests := []struct {
		name         string
		ff           *FeedForward
		position     float64
		velocity     float64
		acceleration float64
		expected     float64
	}{
		{
			name: "basic calculation without options",
			ff: &FeedForward{
				kS:   0.1,
				kV:   0.5,
				kA:   0.02,
				kCos: 0.0,
			},
			position:     1.0,
			velocity:     2.0,
			acceleration: 3.0,
			expected:     0.1 + 0.5*2.0 + 0.02*3.0, // kS + kV*velocity + kA*acceleration
		},
		// Removed calculation with gravity gain
		{
			name: "calculation with cosine gain",
			ff: &FeedForward{
				kS:   0.1,
				kV:   0.5,
				kA:   0.02,
				kCos: 2.0,
			},
			position:     0.0, // cos(0) = 1
			velocity:     1.0,
			acceleration: 0.5,
			expected:     0.1 + 0.5*1.0 + 0.02*0.5 + 2.0*math.Cos(0.0), // kS + kV*velocity + kA*acceleration + kCos*cos(position)
		},
		{
			name: "calculation with cosine gain at π/2",
			ff: &FeedForward{
				kS:   0.1,
				kV:   0.5,
				kA:   0.02,
				kCos: 2.0,
			},
			position:     math.Pi / 2, // cos(π/2) ≈ 0
			velocity:     1.0,
			acceleration: 0.5,
			expected:     0.1 + 0.5*1.0 + 0.02*0.5 + 2.0*math.Cos(math.Pi/2), // kS + kV*velocity + kA*acceleration + kCos*cos(position)
		},
		// Removed calculation with both gravity and cosine gains
		// Removed calculation with zero velocity and acceleration (gravity gain)
		// Removed calculation with negative values (gravity gain)
		{
			name: "calculation with zero cosine gain (should not add cosine term)",
			ff: &FeedForward{
				kS:   0.1,
				kV:   0.5,
				kA:   0.02,
				kCos: 0.0,
			},
			position:     math.Pi / 2,
			velocity:     2.0,
			acceleration: 1.5,
			expected:     0.1 + 0.5*2.0 + 0.02*1.5, // kS + kV*velocity + kA*acceleration
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ff.Calculate(tt.position, tt.velocity, tt.acceleration)
			if math.Abs(result-tt.expected) > tolerance {
				t.Errorf("Calculate() = %f, expected %f (difference: %e)", result, tt.expected, math.Abs(result-tt.expected))
			}
		})
	}
}

func TestCalculateWithConstructor(t *testing.T) {
	const tolerance = 1e-9

	tests := []struct {
		name         string
		kS           float64
		kV           float64
		kA           float64
		opts         []Option
		position     float64
		velocity     float64
		acceleration float64
		expected     float64
	}{
		{
			name:         "basic feedforward",
			kS:           0.1,
			kV:           0.5,
			kA:           0.02,
			opts:         nil,
			position:     1.0,
			velocity:     2.0,
			acceleration: 1.0,
			expected:     0.1 + 0.5*2.0 + 0.02*1.0,
		},
		// Removed with gravity compensation
		{
			name:         "with cosine compensation",
			kS:           0.1,
			kV:           0.5,
			kA:           0.02,
			opts:         []Option{WithCosineGain(2.0)},
			position:     math.Pi,
			velocity:     0.0,
			acceleration: 1.0,
			expected:     0.1 + 0.5*0.0 + 0.02*1.0 + 2.0*math.Cos(math.Pi),
		},
		// Removed with both compensations (gravity+cosine)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := New(tt.kS, tt.kV, tt.kA, tt.opts...)
			result := ff.Calculate(tt.position, tt.velocity, tt.acceleration)
			if math.Abs(result-tt.expected) > tolerance {
				t.Errorf("Calculate() = %f, expected %f (difference: %e)", result, tt.expected, math.Abs(result-tt.expected))
			}
		})
	}
}

// Benchmark tests
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(0.1, 0.5, 0.02)
	}
}

func BenchmarkNewWithOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(9.86, 0.5, 0.02, WithCosineGain(2.5))
	}
}

func BenchmarkCalculateBasic(b *testing.B) {
	ff := New(0.1, 0.5, 0.02)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ff.Calculate(1.0, 2.0, 1.5)
	}
}

func BenchmarkCalculateWithGravity(b *testing.B) {
	ff := New(9.86, 0.5, 0.02)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ff.Calculate(1.0, 2.0, 1.5)
	}
}

func BenchmarkCalculateWithCosine(b *testing.B) {
	ff := New(0.1, 0.5, 0.02, WithCosineGain(2.5))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ff.Calculate(math.Pi/4, 2.0, 1.5)
	}
}

func BenchmarkCalculateWithBoth(b *testing.B) {
	ff := New(9.86, 0.5, 0.02, WithCosineGain(2.5))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ff.Calculate(math.Pi/4, 2.0, 1.5)
	}
}
