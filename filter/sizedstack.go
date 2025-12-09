package filter

import "slices"

// SizedStack represents a fixed-size stack that maintains the most recent N elements.
type SizedStack[T any] []T

// NewSizedStack creates a new generic sized stack with the given capacity.
func NewSizedStack[T any](capacity int) *SizedStack[T] {
	s := make(SizedStack[T], 0, capacity)
	return &s
}

// Push adds an element to the stack, removing the oldest if at capacity.
func (s *SizedStack[T]) Push(value T) {
	if len(*s) < cap(*s) {
		// Still have capacity, append the element
		*s = append(*s, value)
	} else {
		// At capacity, shift elements left and add new element at the end
		copy(*s, (*s)[1:])
		(*s)[len(*s)-1] = value
	}
}

// Peek returns the most recently added element without removing it.
// Returns zero value of T if stack is empty.
func (s *SizedStack[T]) Peek() T {
	var zero T
	if len(*s) == 0 {
		return zero
	}
	return (*s)[len(*s)-1]
}

// Get returns the element at the given index (0 = oldest, size-1 = newest).
// Returns zero value of T if index is out of bounds.
func (s *SizedStack[T]) Get(index int) T {
	var zero T
	if index < 0 || index >= len(*s) {
		return zero
	}
	return (*s)[index]
}

// Size returns the current number of elements in the stack.
func (s *SizedStack[T]) Size() int {
	return len(*s)
}

// ToArray returns a copy of the stack data as a slice.
func (s *SizedStack[T]) ToArray() []T {
	return slices.Clone(*s)
}
