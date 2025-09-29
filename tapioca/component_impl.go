// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Init implements the tea.Model interface.
func (c *Component) Init() tea.Cmd { return nil }

// Update implements the tea.Model interface. It handles ResizeMsg and scroll-related
// messages. When embedding this Component in your own model, you should forward
// ResizeMsg and scroll messages to this method while handling your own messages
// separately.
func (c *Component) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c.UpdateInto(msg)
}

// UpdateInto is identical to Update, but returns a *Component instead of a tea.Model.
//
// This is useful when embedding this Component in your own model, as it avoids
// the need for type assertions.
func (c *Component) UpdateInto(msg tea.Msg) (*Component, tea.Cmd) {
	var cmd []tea.Cmd
	switch msg := msg.(type) {
	case ResizeMsg:
		c.width, c.height = msg.Width, msg.Height
		c.x, c.y = 0, 0 // Reset position after resize
		c.recomputeCachedInfo()
	case ScrollUpMsg:
		if c.vScroll {
			c.y -= int(msg)
			if c.y < 0 {
				c.y = 0
			}
		}
	case ScrollDownMsg:
		if c.vScroll {
			c.y += int(msg)
			if c.y > c.lines-c.height {
				c.y = c.lines - c.height
			}
		}
	case ScrollLeftMsg:
		if c.hScroll {
			c.x -= int(msg)
			if c.x < 0 {
				c.x = 0
			}
		}
	case ScrollRightMsg:
		if c.hScroll {
			c.x += int(msg)
			if c.x > c.maxLineWidth-c.width {
				c.x = c.maxLineWidth - c.width
			}
		}
	case ScrollTopMsg:
		if c.vScroll {
			c.y = 0
		}
	case ScrollBottomMsg:
		if c.vScroll {
			c.y = c.lines - c.height
		}
	case ScrollBeginMsg:
		if c.hScroll {
			c.x = 0
		}
	case ScrollEndMsg:
		if c.hScroll {
			c.x = c.maxLineWidth - c.width
		}
	}

	// Ensure constraints are maintained
	if !c.vScroll {
		c.y = 0
	}
	if !c.hScroll {
		c.x = 0
	}

	if len(cmd) == 0 {
		return c, nil
	}
	return c, tea.Batch(cmd...)
}

// View implements the tea.Model interface. It returns a string representation
// of the current viewport, applying the current scroll position and wrapping
// behavior based on the component's configuration.
func (c *Component) View() string {
	if c.width <= 0 || c.height <= 0 {
		return ""
	}

	entries := c.entries.GetAll()
	if len(entries) == 0 {
		// No entries, return blank screen
		return c.blankScreen()
	}

	// hScroll controls whether wrapping is enabled
	if c.hScroll {
		// No-wrap mode (may have horizontal scrolling)
		return c.viewNoWrap(entries)
	} else {
		// Wrap mode (may have vertical scrolling)
		return c.viewWrap(entries)
	}
}

func (c *Component) blankScreen() string {
	line := strings.Repeat(" ", c.width)
	lines := make([]string, c.height)
	for i := range lines {
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func (c *Component) viewWrap(entries []*Entry) string {
	virtualLines := c.buildVirtualScreen(entries, false)
	lines := make([]string, c.height)

	startY := c.y
	for i := 0; i < c.height; i++ {
		lineIndex := startY + i
		if lineIndex < len(virtualLines) {
			lines[i] = c.padLine(virtualLines[lineIndex])
		} else {
			lines[i] = strings.Repeat(" ", c.width)
		}
	}

	return strings.Join(lines, "\n")
}

func (c *Component) viewNoWrap(entries []*Entry) string {
	virtualLines := c.buildVirtualScreen(entries, true)
	lines := make([]string, c.height)

	startY := c.y
	for i := 0; i < c.height; i++ {
		lineIndex := startY + i
		if lineIndex < len(virtualLines) {
			line := virtualLines[lineIndex]
			// Apply horizontal scrolling
			if c.x < len(line) {
				end := min(c.x+c.width, len(line))
				line = line[c.x:end]
			} else {
				line = ""
			}
			lines[i] = c.padLine(line)
		} else {
			lines[i] = strings.Repeat(" ", c.width)
		}
	}

	return strings.Join(lines, "\n")
}

func (c *Component) buildVirtualScreen(entries []*Entry, noWrap bool) []string {
	var virtualLines []string

	for _, entry := range entries {
		if noWrap {
			// No-wrap mode: each entry becomes one line
			virtualLines = append(virtualLines, entry.String())
		} else {
			// Wrap mode: each entry may become multiple lines
			wrappedLines := entry.Warps(c.width)
			virtualLines = append(virtualLines, wrappedLines...)
		}
	}

	return virtualLines
}

func (c *Component) padLine(line string) string {
	if len(line) >= c.width {
		return line[:c.width]
	}
	return line + strings.Repeat(" ", c.width-len(line))
}

func (c *Component) recomputeCachedInfo() {
	entries := c.entries.GetAll()
	c.recomputeLines(entries)
	c.recomputeMaxLineWidth(entries)
}

func (c *Component) recomputeLines(entries []*Entry) {
	if c.hScroll {
		// When horizontal scrolling is enabled, no wrapping occurs,
		// so virtual screen line count equals number of entries
		c.lines = c.entries.Size()
	} else {
		// When horizontal scrolling is disabled, entries wrap,
		// so virtual screen line count is total lines after wrapping
		c.lines = 0
		for _, e := range entries {
			c.lines += e.Lines(c.width)
		}
	}

	// Virtual screen line count should be at least the physical screen height
	c.lines = max(c.lines, c.height)
}

func (c *Component) recomputeMaxLineWidth(entries []*Entry) {
	c.maxLineWidth = c.width
	if !c.hScroll {
		return
	}

	for _, e := range entries {
		if l := e.Len(); c.maxLineWidth < l {
			c.maxLineWidth = l
		}
	}
}
