package interplut

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	lut := New()
	if lut == nil {
		t.Fatal("New() returned nil")
	}
	if len(lut.x) != 0 || len(lut.y) != 0 || len(lut.m) != 0 {
		t.Error("New InterpLUT should have empty slices")
	}
}

func TestAdd(t *testing.T) {
	lut := New()

	lut.Add(1.0, 2.0)
	lut.Add(2.0, 4.0)
	lut.Add(3.0, 6.0)

	if len(lut.x) != 3 || len(lut.y) != 3 {
		t.Errorf("Expected 3 points, got x:%d, y:%d", len(lut.x), len(lut.y))
	}

	expectedX := []float64{1.0, 2.0, 3.0}
	expectedY := []float64{2.0, 4.0, 6.0}

	for i, x := range lut.x {
		if x != expectedX[i] {
			t.Errorf("X[%d] = %f, expected %f", i, x, expectedX[i])
		}
	}

	for i, y := range lut.y {
		if y != expectedY[i] {
			t.Errorf("Y[%d] = %f, expected %f", i, y, expectedY[i])
		}
	}
}

func TestCreateLUT_InsufficientData(t *testing.T) {
	tests := []struct {
		name   string
		points [][]float64
		want   string
	}{
		{
			name:   "empty",
			points: [][]float64{},
			want:   "there must be at least two control points",
		},
		{
			name:   "single point",
			points: [][]float64{{1.0, 1.0}},
			want:   "there must be at least two control points",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lut := New()
			for _, p := range tt.points {
				lut.Add(p[0], p[1])
			}

			err := lut.CreateLUT()
			if err == nil {
				t.Error("Expected error, got nil")
			}
			if !contains(err.Error(), tt.want) {
				t.Errorf("Error = %q, want to contain %q", err.Error(), tt.want)
			}
		})
	}
}

func TestCreateLUT_DuplicateX(t *testing.T) {
	lut := New()
	lut.Add(1.0, 1.0)
	lut.Add(1.0, 3.0) // Duplicate X value
	lut.Add(4.0, 1.0)

	err := lut.CreateLUT()
	if err == nil {
		t.Fatal("Expected error for duplicate X values")
	}

	expected := "strictly increasing X values"
	if !contains(err.Error(), expected) {
		t.Errorf("Error = %q, want to contain %q", err.Error(), expected)
	}
}

func TestCreateLUT_Success(t *testing.T) {
	lut := New()
	lut.Add(1.0, 1.0)
	lut.Add(2.0, 2.0)
	lut.Add(3.0, 3.0)
	lut.Add(4.0, 4.0)

	err := lut.CreateLUT()
	if err != nil {
		t.Fatalf("CreateLUT() error = %v", err)
	}

	// Check that data is sorted and tangents are calculated
	if len(lut.m) != 4 {
		t.Errorf("Expected 4 tangent values, got %d", len(lut.m))
	}
}

func TestGet_ExactMatch(t *testing.T) {
	lut := New()
	lut.Add(1.0, 1.0)
	lut.Add(2.0, 2.0)
	lut.Add(3.0, 3.0)
	lut.Add(4.0, 4.0)

	err := lut.CreateLUT()
	if err != nil {
		t.Fatalf("CreateLUT() error = %v", err)
	}

	// Test exact matches at control points
	for i, x := range lut.x {
		result, err := lut.Get(x)
		if err != nil {
			t.Errorf("Get(%f) error = %v", x, err)
		}

		if result != lut.y[i] {
			t.Errorf("Get(%f) = %f, expected %f", x, result, lut.y[i])
		}
	}
}

func TestGet_ComplexCase(t *testing.T) {
	lut := New()
	lut.Add(1.1, 0.2)
	lut.Add(2.7, 0.5)
	lut.Add(3.6, 0.75)
	lut.Add(4.1, 0.9)
	lut.Add(5.0, 1.0)

	err := lut.CreateLUT()
	if err != nil {
		t.Fatalf("CreateLUT() error = %v", err)
	}

	// Test interpolation at a point between control points
	result, err := lut.Get(3.0)
	if err != nil {
		t.Errorf("Get(3.0) error = %v", err)
	}

	// Result should be between the neighboring control points
	if result < 0.5 || result > 0.75 {
		t.Errorf("Get(3.0) = %f, expected between 0.5 and 0.75", result)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmarks
func BenchmarkCreateLUT(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				lut := New()
				for j := 0; j < size; j++ {
					lut.Add(float64(j), float64(j*j))
				}
				lut.CreateLUT()
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	lut := New()
	for i := 0; i < 100; i++ {
		lut.Add(float64(i), float64(i*i))
	}
	lut.CreateLUT()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lut.Get(50.5)
	}
}
