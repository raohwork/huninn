// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCircularBuffer(t *testing.T) {
	t.Run("create buffer with positive size", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		require.NotNil(t, cb)
		assert.Equal(t, 0, cb.Size(), "new buffer should be empty")
		assert.Nil(t, cb.GetAll(), "empty buffer should return nil")
	})

	t.Run("create buffer with size 1", func(t *testing.T) {
		cb := NewCircularBuffer[string](1)
		require.NotNil(t, cb)
		assert.Equal(t, 0, cb.Size(), "new buffer should be empty")
	})

	t.Run("create buffer with zero size defaults to size 1", func(t *testing.T) {
		cb := NewCircularBuffer[int](0)
		require.NotNil(t, cb)
		cb.Add(1)
		cb.Add(2)
		assert.Equal(t, 1, cb.Size(), "buffer with size 0 should default to capacity 1")
	})

	t.Run("create buffer with negative size defaults to size 1", func(t *testing.T) {
		cb := NewCircularBuffer[int](-5)
		require.NotNil(t, cb)
		cb.Add(1)
		cb.Add(2)
		assert.Equal(t, 1, cb.Size(), "buffer with negative size should default to capacity 1")
	})
}

func TestCircularBuffer_AddAndSize(t *testing.T) {
	t.Run("add single item to empty buffer", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		cb.Add(42)
		assert.Equal(t, 1, cb.Size(), "buffer should contain one item")
	})

	t.Run("add multiple items without exceeding capacity", func(t *testing.T) {
		cb := NewCircularBuffer[string](3)
		cb.Add("first")
		cb.Add("second")
		assert.Equal(t, 2, cb.Size(), "buffer should contain two items")

		cb.Add("third")
		assert.Equal(t, 3, cb.Size(), "buffer should contain three items")
	})

	t.Run("add items exceeding capacity should overwrite oldest", func(t *testing.T) {
		cb := NewCircularBuffer[int](3)
		cb.Add(1)
		cb.Add(2)
		cb.Add(3)
		assert.Equal(t, 3, cb.Size(), "buffer should be full")

		cb.Add(4)
		assert.Equal(t, 3, cb.Size(), "buffer size should remain at capacity")

		items := cb.GetAll()
		assert.Equal(t, []int{2, 3, 4}, items, "oldest item should be overwritten")
	})

	t.Run("add to buffer of size 1", func(t *testing.T) {
		cb := NewCircularBuffer[int](1)
		cb.Add(1)
		assert.Equal(t, 1, cb.Size())

		cb.Add(2)
		assert.Equal(t, 1, cb.Size())

		items := cb.GetAll()
		assert.Equal(t, []int{2}, items, "should contain only the latest item")
	})

	t.Run("add many items beyond capacity", func(t *testing.T) {
		cb := NewCircularBuffer[int](2)
		for i := 1; i <= 10; i++ {
			cb.Add(i)
		}
		assert.Equal(t, 2, cb.Size(), "buffer size should remain at capacity")

		items := cb.GetAll()
		assert.Equal(t, []int{9, 10}, items, "should contain only the last two items")
	})
}

func TestCircularBuffer_GetAllAndGetAllReverse(t *testing.T) {
	t.Run("empty buffer returns nil for both methods", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		assert.Nil(t, cb.GetAll(), "GetAll should return nil for empty buffer")
		assert.Nil(t, cb.GetAllReverse(), "GetAllReverse should return nil for empty buffer")
	})

	t.Run("single item buffer", func(t *testing.T) {
		cb := NewCircularBuffer[string](3)
		cb.Add("only")

		all := cb.GetAll()
		allReverse := cb.GetAllReverse()

		assert.Equal(t, []string{"only"}, all, "GetAll should return single item")
		assert.Equal(t, all, allReverse, "both methods should return same result for single item")
	})

	t.Run("partially filled buffer maintains order", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		cb.Add(1)
		cb.Add(2)
		cb.Add(3)

		all := cb.GetAll()
		allReverse := cb.GetAllReverse()

		assert.Equal(t, []int{1, 2, 3}, all, "GetAll should return items in insertion order")
		reversed := slices.Clone(all)
		slices.Reverse(reversed)
		assert.Equal(t, reversed, allReverse, "GetAllReverse should be exact reverse of GetAll")
	})

	t.Run("full buffer without wraparound", func(t *testing.T) {
		cb := NewCircularBuffer[string](3)
		cb.Add("a")
		cb.Add("b")
		cb.Add("c")

		all := cb.GetAll()
		allReverse := cb.GetAllReverse()

		assert.Equal(t, []string{"a", "b", "c"}, all)
		reversed := slices.Clone(all)
		slices.Reverse(reversed)
		assert.Equal(t, reversed, allReverse)
	})

	t.Run("buffer after wraparound", func(t *testing.T) {
		cb := NewCircularBuffer[int](3)
		cb.Add(1)
		cb.Add(2)
		cb.Add(3)
		cb.Add(4) // overwrites 1
		cb.Add(5) // overwrites 2

		all := cb.GetAll()
		allReverse := cb.GetAllReverse()

		assert.Equal(t, []int{3, 4, 5}, all, "GetAll should return remaining items in order")
		reversed := slices.Clone(all)
		slices.Reverse(reversed)
		assert.Equal(t, reversed, allReverse)
	})

	t.Run("buffer with complex wraparound pattern", func(t *testing.T) {
		cb := NewCircularBuffer[int](4)
		// Fill buffer completely
		for i := 1; i <= 4; i++ {
			cb.Add(i)
		}
		// Add more items to cause multiple wraparounds
		for i := 5; i <= 10; i++ {
			cb.Add(i)
		}

		all := cb.GetAll()
		allReverse := cb.GetAllReverse()

		assert.Equal(t, 4, len(all), "should contain exactly 4 items")
		assert.Equal(t, []int{7, 8, 9, 10}, all, "should contain last 4 items in order")
		reversed := slices.Clone(all)
		slices.Reverse(reversed)
		assert.Equal(t, reversed, allReverse)
	})
}

func TestCircularBuffer_IntegrationScenarios(t *testing.T) {
	t.Run("continuous add and get operations", func(t *testing.T) {
		cb := NewCircularBuffer[int](3)

		// First phase: fill buffer
		cb.Add(1)
		cb.Add(2)
		assert.Equal(t, []int{1, 2}, cb.GetAll())
		assert.Equal(t, 2, cb.Size())

		// Second phase: fill completely
		cb.Add(3)
		assert.Equal(t, []int{1, 2, 3}, cb.GetAll())
		assert.Equal(t, 3, cb.Size())

		// Third phase: cause wraparound
		cb.Add(4)
		assert.Equal(t, []int{2, 3, 4}, cb.GetAll())
		assert.Equal(t, 3, cb.Size())
	})

	t.Run("edge case with different data types", func(t *testing.T) {
		type testStruct struct {
			ID   int
			Name string
		}

		cb := NewCircularBuffer[testStruct](2)
		cb.Add(testStruct{1, "first"})
		cb.Add(testStruct{2, "second"})
		cb.Add(testStruct{3, "third"}) // should overwrite first

		all := cb.GetAll()
		expected := []testStruct{{2, "second"}, {3, "third"}}
		assert.Equal(t, expected, all)
		assert.Equal(t, 2, cb.Size())
	})

	t.Run("size consistency across operations", func(t *testing.T) {
		cb := NewCircularBuffer[int](2)

		// Empty state
		assert.Equal(t, 0, cb.Size())
		assert.Nil(t, cb.GetAll())

		// Add one item
		cb.Add(1)
		assert.Equal(t, 1, cb.Size())
		assert.Equal(t, []int{1}, cb.GetAll())

		// Fill buffer
		cb.Add(2)
		assert.Equal(t, 2, cb.Size())
		assert.Equal(t, []int{1, 2}, cb.GetAll())

		// Cause wraparound
		cb.Add(3)
		assert.Equal(t, 2, cb.Size(), "size should remain constant after wraparound")
		assert.Equal(t, []int{2, 3}, cb.GetAll())
	})
}
