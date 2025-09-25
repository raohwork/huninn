// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

type CircularBuffer[T any] struct {
	data       []T
	start, end int
}

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

func (cb *CircularBuffer[T]) Add(item T) {
	cb.data[cb.end] = item
	newEnd := (cb.end + 1) % len(cb.data)

	if newEnd == cb.start {
		cb.start = (cb.start + 1) % len(cb.data)
	}
	cb.end = newEnd
}

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

// GetAllReverse returns all elements in reverse order
func (cb *CircularBuffer[T]) GetAllReverse() []T {
	if cb.start == cb.end {
		return nil
	}

	var result []T
	if cb.end > cb.start {
		for i := cb.end - 1; i >= cb.start; i-- {
			result = append(result, cb.data[i])
		}
	} else {
		for i := cb.end - 1; i >= 0; i-- {
			result = append(result, cb.data[i])
		}
		for i := len(cb.data) - 1; i >= cb.start; i-- {
			result = append(result, cb.data[i])
		}
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
