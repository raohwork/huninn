// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pearl

import "github.com/raohwork/huninn/tapioca"

// BufferedBlock is the default implementation of a huninn component that provides
// scrollable text display functionality.
//
// If you need proper control over text display, you might want to take a look
// at [tapioca.Entry].
type BufferedBlock struct {
	// enable horizontal scroll, will disable line wrap
	hScroll bool
	// enable vertical scroll
	vScroll bool

	width  int
	height int
	x, y   int

	entries *tapioca.CircularBuffer[*tapioca.Entry]

	// cached info

	// total number of lines, used and updated only if hScroll is
	// false (line wrap is enabled)
	// this is used to calculate vertical scroll range
	lines int
	// max line width, used and updated only if hScroll is true
	// this is used to calculate horizontal scroll range
	maxLineWidth int
}

// Append adds a new entry to the end of the virtual screen. In a log panel context,
// this would add a new log message at the bottom (older messages are shown by default
// since the viewport starts at position 0,0).
func (c *BufferedBlock) Append(str string) {
	c.entries.Append(tapioca.NewEntry(str))
	c.recomputeCachedInfo()
}

// Prepend adds a new entry to the beginning of the virtual screen. In a log panel
// context, this would add a new log message at the top (newer messages are shown
// by default since the viewport starts at position 0,0).
func (c *BufferedBlock) Prepend(str string) {
	c.entries.Prepend(tapioca.NewEntry(str))
	c.recomputeCachedInfo()
}

// Clear removes all entries from the component.
func (c *BufferedBlock) Clear() {
	c.entries.Reset()
	c.recomputeCachedInfo()
}

// Entries returns all entries currently stored in the component.
func (c *BufferedBlock) Entries() []*tapioca.Entry {
	return c.entries.GetAll()
}

// Capacity returns the maximum number of entries the component can hold.
func (c *BufferedBlock) Capacity() int {
	return c.entries.Capacity()
}

// ResizeBuffer changes the capacity of the entries buffer to newSize.
func (c *BufferedBlock) ResizeBuffer(newSize int) {
	c.entries.Resize(newSize)
}

// Width returns the current width of the component.
func (c *BufferedBlock) Width() int {
	return c.width
}

// Height returns the current height of the component.
func (c *BufferedBlock) Height() int {
	return c.height
}

// NewBufferedBlock creates a new component with the specified entry capacity.
// The size parameter determines how many entries the circular buffer can hold.
// When the buffer is full, adding new entries will overwrite the oldest ones.
func NewBufferedBlock(size int, hScroll, vScroll bool) *BufferedBlock {
	ret := &BufferedBlock{
		entries: tapioca.NewCircularBuffer[*tapioca.Entry](size),
		hScroll: hScroll,
		vScroll: vScroll,
	}
	return ret
}
