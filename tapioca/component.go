// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

// ComponentOption is a function type used to configure Component behavior.
// Options are applied during component creation via NewComponent.
type ComponentOption func(*Component)

// HorizontalScrollable enables horizontal scrolling for the component.
//
// When enabled, text wrapping is disabled and entries are displayed as single
// lines that can be scrolled horizontally. This is useful for viewing long
// lines of text such as log entries or code.
//
// You have to send ScrollXXXMsg to [Component.Update] to actually scroll
// the content.
func HorizontalScrollable() ComponentOption {
	return func(c *Component) {
		c.hScroll = true
	}
}

// VerticalScrollable enables vertical scrolling for the component.
//
// When enabled, the component can scroll through entries vertically.
// This is useful when you have more content than can fit in the available
// height and want to allow users to navigate through the entries.
//
// You have to send ScrollXXXMsg to [Component.Update] to actually scroll
// the content.
func VerticalScrollable() ComponentOption {
	return func(c *Component) {
		c.vScroll = true
	}
}

// Component is the default implementation of a huninn component that provides
// scrollable text display functionality.
//
// If you need proper control over text display, you might want to take a look
// at [Entry].
//
// A huninn component MUST fulfill the following requirements:
//   - Implement tea.Model interface
//   - Handle ResizeMsg correctly
//   - View() must return a string that exactly fits in the given width
//     and height, like "a  \nb  \n   " for a 3x3 component
//
// This Component struct provides some extra functionalities including:
//   - Entry storage using a circular buffer as scroll buffer
//   - Horizontal and vertical scrolling
//   - Text wrapping
//
// # Embedding Component
//
// To implement your own component, embed this Component struct and override
// the Update() method. You should forward ResizeMsg and scroll messages to
// the embedded Component's Update() method while handling your own messages
// separately:
//
//	type MyComponent struct {
//	    Component
//	    // your fields
//	}
//
//	func (m MyComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	    switch msg := msg.(type) {
//	    case MyCustomMsg:
//	        // handle your messages
//	    default:
//	        // forward to embedded Component
//	        var cmd tea.Cmd
//	        m.Component, cmd = m.Component.Update(msg)
//	        return m, cmd
//	    }
//	}
//
// # Scrolling Behavior
//
// Scrolling behavior is controlled by ComponentOption functions:
//   - HorizontalScrollable(): Enables horizontal scrolling, disables text wrapping
//   - VerticalScrollable(): Enables vertical scrolling through entries
//
// When neither scrolling option is enabled, the component displays entries
// with text wrapping but without scrolling capability.
//
// # Triggering scrolling
//
// You have to send ScrollXXXMsg to [Component.Update] to actually scroll
// the content. The component does not scroll automatically when new entries
// are added.
//
// # Resize Behavior
//
// When a ResizeMsg is received, the component resets the scroll position to (0,0)
// to avoid complex position preservation logic that may not provide good UX.
type Component struct {
	// enable horizontal scroll, will disable line wrap
	hScroll bool
	// enable vertical scroll
	vScroll bool

	width  int
	height int
	x, y   int

	entries *CircularBuffer[*Entry]

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
func (c *Component) Append(str string) {
	c.entries.Append(NewEntry(str))
	c.recomputeCachedInfo()
}

// Prepend adds a new entry to the beginning of the virtual screen. In a log panel
// context, this would add a new log message at the top (newer messages are shown
// by default since the viewport starts at position 0,0).
func (c *Component) Prepend(str string) {
	c.entries.Prepend(NewEntry(str))
	c.recomputeCachedInfo()
}

// Clear removes all entries from the component.
func (c *Component) Clear() {
	c.entries.Reset()
	c.recomputeCachedInfo()
}

// Entries returns all entries currently stored in the component.
func (c *Component) Entries() []*Entry {
	return c.entries.GetAll()
}

// Capacity returns the maximum number of entries the component can hold.
func (c *Component) Capacity() int {
	return c.entries.Capacity()
}

// ResizeBuffer changes the capacity of the entries buffer to newSize.
func (c *Component) ResizeBuffer(newSize int) {
	c.entries.Resize(newSize)
}

// Width returns the current width of the component.
func (c *Component) Width() int {
	return c.width
}

// Height returns the current height of the component.
func (c *Component) Height() int {
	return c.height
}

// NewComponent creates a new component with the specified entry capacity.
// The size parameter determines how many entries the circular buffer can hold.
// When the buffer is full, adding new entries will overwrite the oldest ones.
func NewComponent(size int, options ...ComponentOption) *Component {
	ret := &Component{
		entries: NewCircularBuffer[*Entry](size),
	}
	for _, o := range options {
		o(ret)
	}
	return ret
}
