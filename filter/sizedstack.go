package filter

// SizedStack represents a fixed-size stack that maintains the most recent N elements.
type SizedStack struct {
	data     []float64
	size     int
	capacity int
}

// NewSizedStack creates a new sized stack with the given capacity.
func NewSizedStack(capacity int) *SizedStack {
	return &SizedStack{
		data:     make([]float64, 0, capacity),
		capacity: capacity,
	}
}

// Push adds an element to the stack, removing the oldest if at capacity.
func (s *SizedStack) Push(value float64) {
	if s.size < s.capacity {
		s.data = append(s.data, value)
		s.size++
	} else {
		// Shift all elements left and add new element at the end
		copy(s.data, s.data[1:])
		s.data[s.size-1] = value
	}
}

// Peek returns the most recently added element without removing it.
func (s *SizedStack) Peek() float64 {
	if s.size == 0 {
		return 0.0
	}
	return s.data[s.size-1]
}

// Get returns the element at the given index (0 = oldest, size-1 = newest).
func (s *SizedStack) Get(index int) float64 {
	if index < 0 || index >= s.size {
		return 0.0
	}
	return s.data[index]
}

// Size returns the current number of elements in the stack.
func (s *SizedStack) Size() int {
	return s.size
}

// ToArray returns a copy of the stack data as a slice.
func (s *SizedStack) ToArray() []float64 {
	result := make([]float64, s.size)
	copy(result, s.data[:s.size])
	return result
}
