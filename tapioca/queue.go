// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

// CircularBuffer is a generic circular buffer implementation.
type CircularBuffer[T any] struct {
	data       []T
	start, end int
}

// NewCircularBuffer creates a new CircularBuffer with the given size.
func NewCircularBuffer[T any](size int) *CircularBuffer[T] {
	if size < 1 {
		size = 1
	}

	return &CircularBuffer[T]{
		data:  make([]T, size+1),
		start: 0,
		end:   0,
	}
}

// Reset clears the buffer
func (cb *CircularBuffer[T]) Reset() {
	cb.start = 0
	cb.end = 0
}

// Append adds an item to the end of the buffer, removing the oldest item if full
func (cb *CircularBuffer[T]) Append(item T) {
	cb.data[cb.end] = item
	newEnd := (cb.end + 1) % len(cb.data)

	if newEnd == cb.start {
		cb.start = (cb.start + 1) % len(cb.data)
	}
	cb.end = newEnd
}

// Prepend adds an item to the start of the buffer, removing the oldest item if full
func (cb *CircularBuffer[T]) Prepend(item T) {
	cb.start = (cb.start - 1 + len(cb.data)) % len(cb.data)
	cb.data[cb.start] = item

	if cb.start == cb.end {
		cb.end = (cb.end - 1 + len(cb.data)) % len(cb.data)
	}
}

// GetAll returns all elements in the buffer in order from oldest to newest
func (cb *CircularBuffer[T]) GetAll() []T {
	if cb.start == cb.end {
		return nil
	}

	var result []T
	if cb.end > cb.start {
		result = cb.data[cb.start:cb.end]
	} else {
		result = append(cb.data[cb.start:], cb.data[:cb.end]...)
	}
	return result
}

// Size returns the number of elements currently in the buffer
func (cb *CircularBuffer[T]) Size() int {
	if cb.end >= cb.start {
		return cb.end - cb.start
	}
	return len(cb.data) - cb.start + cb.end
}

// Capacity returns the maximum number of elements the buffer can hold
func (cb *CircularBuffer[T]) Capacity() int {
	return len(cb.data) - 1
}

// Resize changes the capacity of the buffer, preserving existing elements
func (cb *CircularBuffer[T]) Resize(newSize int) {
	if newSize < 1 {
		newSize = 1
	}

	c := cb.Capacity()
	if newSize == c {
		return
	}

	arr := cb.GetAll()
	l := len(arr)
	newData := make([]T, newSize+1)
	end := min(l, newSize)
	copy(newData, arr[:end])
	cb.data = newData
	cb.start = 0
	cb.end = end
}
