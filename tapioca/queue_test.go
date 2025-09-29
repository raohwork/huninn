// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
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
		cb.Append(1)
		cb.Append(2)
		assert.Equal(t, 1, cb.Size(), "buffer with size 0 should default to capacity 1")
	})

	t.Run("create buffer with negative size defaults to size 1", func(t *testing.T) {
		cb := NewCircularBuffer[int](-5)
		require.NotNil(t, cb)
		cb.Append(1)
		cb.Append(2)
		assert.Equal(t, 1, cb.Size(), "buffer with negative size should default to capacity 1")
	})
}

func TestCircularBuffer_AppendAndSize(t *testing.T) {
	t.Run("add single item to empty buffer", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		cb.Append(42)
		assert.Equal(t, 1, cb.Size(), "buffer should contain one item")
	})

	t.Run("add multiple items without exceeding capacity", func(t *testing.T) {
		cb := NewCircularBuffer[string](3)
		cb.Append("first")
		cb.Append("second")
		assert.Equal(t, 2, cb.Size(), "buffer should contain two items")

		cb.Append("third")
		assert.Equal(t, 3, cb.Size(), "buffer should contain three items")
	})

	t.Run("add items exceeding capacity should overwrite oldest", func(t *testing.T) {
		cb := NewCircularBuffer[int](3)
		cb.Append(1)
		cb.Append(2)
		cb.Append(3)
		assert.Equal(t, 3, cb.Size(), "buffer should be full")

		cb.Append(4)
		assert.Equal(t, 3, cb.Size(), "buffer size should remain at capacity")

		items := cb.GetAll()
		assert.Equal(t, []int{2, 3, 4}, items, "oldest item should be overwritten")
	})

	t.Run("add to buffer of size 1", func(t *testing.T) {
		cb := NewCircularBuffer[int](1)
		cb.Append(1)
		assert.Equal(t, 1, cb.Size())

		cb.Append(2)
		assert.Equal(t, 1, cb.Size())

		items := cb.GetAll()
		assert.Equal(t, []int{2}, items, "should contain only the latest item")
	})

	t.Run("add many items beyond capacity", func(t *testing.T) {
		cb := NewCircularBuffer[int](2)
		for i := 1; i <= 10; i++ {
			cb.Append(i)
		}
		assert.Equal(t, 2, cb.Size(), "buffer size should remain at capacity")

		items := cb.GetAll()
		assert.Equal(t, []int{9, 10}, items, "should contain only the last two items")
	})
}

func TestCircularBuffer_AppendAndGet(t *testing.T) {
	t.Run("empty buffer returns nil for both methods", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		assert.Nil(t, cb.GetAll(), "GetAll should return nil for empty buffer")
	})

	t.Run("single item buffer", func(t *testing.T) {
		cb := NewCircularBuffer[string](3)
		cb.Append("only")

		all := cb.GetAll()

		assert.Equal(t, []string{"only"}, all, "GetAll should return single item")
	})

	t.Run("partially filled buffer maintains order", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		cb.Append(1)
		cb.Append(2)
		cb.Append(3)

		all := cb.GetAll()

		assert.Equal(t, []int{1, 2, 3}, all, "GetAll should return items in insertion order")
	})

	t.Run("full buffer without wraparound", func(t *testing.T) {
		cb := NewCircularBuffer[string](3)
		cb.Append("a")
		cb.Append("b")
		cb.Append("c")

		all := cb.GetAll()

		assert.Equal(t, []string{"a", "b", "c"}, all)
	})

	t.Run("buffer after wraparound", func(t *testing.T) {
		cb := NewCircularBuffer[int](3)
		cb.Append(1)
		cb.Append(2)
		cb.Append(3)
		cb.Append(4) // overwrites 1
		cb.Append(5) // overwrites 2

		all := cb.GetAll()

		assert.Equal(t, []int{3, 4, 5}, all, "GetAll should return remaining items in order")
	})

	t.Run("buffer with complex wraparound pattern", func(t *testing.T) {
		cb := NewCircularBuffer[int](4)
		// Fill buffer completely
		for i := 1; i <= 4; i++ {
			cb.Append(i)
		}
		// Add more items to cause multiple wraparounds
		for i := 5; i <= 10; i++ {
			cb.Append(i)
		}

		all := cb.GetAll()

		assert.Equal(t, 4, len(all), "should contain exactly 4 items")
		assert.Equal(t, []int{7, 8, 9, 10}, all, "should contain last 4 items in order")
	})
}

func TestCircularBuffer_PrependAndGet(t *testing.T) {
	t.Run("empty buffer returns nil for both methods", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		assert.Nil(t, cb.GetAll(), "GetAll should return nil for empty buffer")
	})

	t.Run("single item buffer", func(t *testing.T) {
		cb := NewCircularBuffer[string](3)
		cb.Prepend("only")

		all := cb.GetAll()

		assert.Equal(t, []string{"only"}, all, "GetAll should return single item")
	})

	t.Run("partially filled buffer maintains order", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		cb.Prepend(1)
		cb.Prepend(2)
		cb.Prepend(3)

		all := cb.GetAll()

		assert.Equal(t, []int{3, 2, 1}, all, "GetAll should return items in insertion order")
	})

	t.Run("full buffer without wraparound", func(t *testing.T) {
		cb := NewCircularBuffer[string](3)
		cb.Prepend("a")
		cb.Prepend("b")
		cb.Prepend("c")

		all := cb.GetAll()

		assert.Equal(t, []string{"c", "b", "a"}, all)
	})

	t.Run("buffer after wraparound", func(t *testing.T) {
		cb := NewCircularBuffer[int](3)
		cb.Prepend(1)
		cb.Prepend(2)
		cb.Prepend(3)
		cb.Prepend(4) // overwrites 1
		cb.Prepend(5) // overwrites 2

		all := cb.GetAll()

		assert.Equal(t, []int{5, 4, 3}, all, "GetAll should return remaining items in order")
	})

	t.Run("buffer with complex wraparound pattern", func(t *testing.T) {
		cb := NewCircularBuffer[int](4)
		// Fill buffer completely
		for i := 1; i <= 4; i++ {
			cb.Prepend(i)
		}
		// Add more items to cause multiple wraparounds
		for i := 5; i <= 10; i++ {
			cb.Prepend(i)
		}

		all := cb.GetAll()

		assert.Equal(t, 4, len(all), "should contain exactly 4 items")
		assert.Equal(t, []int{10, 9, 8, 7}, all, "should contain last 4 items in order")
	})
}

func TestCircularBuffer_IntegrationScenarios(t *testing.T) {
	t.Run("continuous add and get operations", func(t *testing.T) {
		cb := NewCircularBuffer[int](3)

		// First phase: fill buffer
		cb.Append(1)
		cb.Append(2)
		assert.Equal(t, []int{1, 2}, cb.GetAll())
		assert.Equal(t, 2, cb.Size())

		// Second phase: fill completely
		cb.Append(3)
		assert.Equal(t, []int{1, 2, 3}, cb.GetAll())
		assert.Equal(t, 3, cb.Size())

		// Third phase: cause wraparound
		cb.Append(4)
		assert.Equal(t, []int{2, 3, 4}, cb.GetAll())
		assert.Equal(t, 3, cb.Size())
	})

	t.Run("edge case with different data types", func(t *testing.T) {
		type testStruct struct {
			ID   int
			Name string
		}

		cb := NewCircularBuffer[testStruct](2)
		cb.Append(testStruct{1, "first"})
		cb.Append(testStruct{2, "second"})
		cb.Append(testStruct{3, "third"}) // should overwrite first

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
		cb.Append(1)
		assert.Equal(t, 1, cb.Size())
		assert.Equal(t, []int{1}, cb.GetAll())

		// Fill buffer
		cb.Append(2)
		assert.Equal(t, 2, cb.Size())
		assert.Equal(t, []int{1, 2}, cb.GetAll())

		// Cause wraparound
		cb.Append(3)
		assert.Equal(t, 2, cb.Size(), "size should remain constant after wraparound")
		assert.Equal(t, []int{2, 3}, cb.GetAll())
	})
}

func TestCircularBuffer_Resize(t *testing.T) {
	t.Run("shrink buffer smaller than current size", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		cb.Append(1)
		cb.Append(2)
		cb.Append(3)
		cb.Append(4)

		cb.Resize(2)
		assert.Equal(t, 2, cb.Capacity(), "capacity should be updated")
		assert.Equal(t, 2, cb.Size(), "size should be truncated")
		assert.Equal(t, []int{1, 2}, cb.GetAll(), "should keep the first elements")
	})

	t.Run("shrink buffer larger than current size", func(t *testing.T) {
		cb := NewCircularBuffer[int](10)
		cb.Append(1)
		cb.Append(2)
		cb.Append(3)

		cb.Resize(5)
		assert.Equal(t, 5, cb.Capacity(), "capacity should be updated")
		assert.Equal(t, 3, cb.Size(), "size should be preserved")
		assert.Equal(t, []int{1, 2, 3}, cb.GetAll(), "elements should be preserved")
	})

	t.Run("expand buffer", func(t *testing.T) {
		cb := NewCircularBuffer[int](3)
		cb.Append(1)
		cb.Append(2)

		cb.Resize(5)
		assert.Equal(t, 5, cb.Capacity(), "capacity should be expanded")
		assert.Equal(t, 2, cb.Size(), "size should be preserved")
		assert.Equal(t, []int{1, 2}, cb.GetAll(), "elements should be preserved")

		cb.Append(3)
		cb.Append(4)
		cb.Append(5)
		assert.Equal(t, 5, cb.Size(), "should be able to fill new capacity")
		assert.Equal(t, []int{1, 2, 3, 4, 5}, cb.GetAll())
	})

	t.Run("resize to same size", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		cb.Append(1)
		cb.Append(2)
		cb.Append(3)

		cb.Resize(5)
		assert.Equal(t, 5, cb.Capacity(), "capacity should not change")
		assert.Equal(t, 3, cb.Size(), "size should not change")
		assert.Equal(t, []int{1, 2, 3}, cb.GetAll(), "elements should not change")
	})

	t.Run("resize to zero or less defaults to 1", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		cb.Append(1)
		cb.Append(2)

		cb.Resize(0)
		assert.Equal(t, 1, cb.Capacity(), "capacity should be 1")
		assert.Equal(t, 1, cb.Size(), "size should be 1")
		assert.Equal(t, []int{1}, cb.GetAll(), "should keep first element")
	})

	t.Run("resize empty buffer", func(t *testing.T) {
		cb := NewCircularBuffer[int](5)
		cb.Resize(10)
		assert.Equal(t, 10, cb.Capacity(), "capacity should be updated")
		assert.Equal(t, 0, cb.Size(), "buffer should remain empty")
		assert.Nil(t, cb.GetAll(), "buffer should remain empty")
	})

	t.Run("resize full buffer with wraparound", func(t *testing.T) {
		cb := NewCircularBuffer[int](3)
		cb.Append(1)
		cb.Append(2)
		cb.Append(3)
		cb.Append(4) // Wraps around, state is [2, 3, 4]

		cb.Resize(2)
		assert.Equal(t, 2, cb.Capacity(), "capacity should be shrunk")
		assert.Equal(t, 2, cb.Size(), "size should be shrunk")
		assert.Equal(t, []int{2, 3}, cb.GetAll(), "should keep correct elements after wraparound")
	})

	t.Run("expand buffer with wraparound", func(t *testing.T) {
		cb := NewCircularBuffer[int](3)
		cb.Append(1)
		cb.Append(2)
		cb.Append(3)
		cb.Append(4) // Wraps around, state is [2, 3, 4]

		cb.Resize(5)
		assert.Equal(t, 5, cb.Capacity(), "capacity should be expanded")
		assert.Equal(t, 3, cb.Size(), "size should be preserved")
		assert.Equal(t, []int{2, 3, 4}, cb.GetAll(), "elements should be preserved in order")

		cb.Append(5)
		cb.Append(6)
		assert.Equal(t, 5, cb.Size(), "should be able to fill new capacity")
		assert.Equal(t, []int{2, 3, 4, 5, 6}, cb.GetAll())
	})
}
