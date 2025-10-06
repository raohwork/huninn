// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pearl

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
)

// Init implements the tea.Model interface.
func (c *BufferedBlock) Init() tea.Cmd { return nil }

// Update implements the tea.Model interface. It handles ResizeMsg and scroll-related
// messages. When embedding this Component in your own model, you should forward
// ResizeMsg and scroll messages to this method while handling your own messages
// separately.
func (c *BufferedBlock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c.UpdateInto(msg)
}

// UpdateInto is identical to Update, but returns a *Component instead of a tea.Model.
//
// This is useful when embedding this Component in your own model, as it avoids
// the need for type assertions.
func (c *BufferedBlock) UpdateInto(msg tea.Msg) (*BufferedBlock, tea.Cmd) {
	var cmd []tea.Cmd
	switch msg := msg.(type) {
	case tapioca.ResizeMsg:
		c.HandleEvent(msg)
		c.recomputeCachedInfo()
	case tapioca.ScrollUpMsg,
		tapioca.ScrollDownMsg,
		tapioca.ScrollLeftMsg,
		tapioca.ScrollRightMsg,
		tapioca.ScrollTopMsg,
		tapioca.ScrollBottomMsg,
		tapioca.ScrollBeginMsg,
		tapioca.ScrollEndMsg,
		tapioca.ScrollToMsg:
		c.HandleEvent(msg)
	}

	// Ensure constraints are maintained
	if !c.vScroll {
		c.ScrollToTop()
	}
	if !c.hScroll {
		c.ScrollToBegin()
	}

	if len(cmd) == 0 {
		return c, nil
	}
	return c, tea.Batch(cmd...)
}

// View implements the tea.Model interface. It returns a string representation
// of the current viewport, applying the current scroll position and wrapping
// behavior based on the component's configuration.
func (c *BufferedBlock) View() string {
	if c.Width() <= 0 || c.Height() <= 0 {
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

func (c *BufferedBlock) blankScreen() string {
	line := strings.Repeat(" ", c.Width())
	lines := make([]string, c.Height())
	for i := range lines {
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func (c *BufferedBlock) viewWrap(entries []*tapioca.Entry) string {
	lines := make([]string, 0, c.Y()+c.Height())
	totalEntries := len(entries)

	curLine, curIdx := 0, 0
	// fill lines
	for curLine < c.Y()+c.Height() && curIdx < totalEntries {
		entry := entries[curIdx]
		l := entry.StyledBlock(c.Width())
		h := len(l)

		want := min(h, c.Y()+c.Height()-curLine)
		lines = append(lines, l[:want]...)
		curIdx++
		curLine += want
	}
	if curLine < c.Y()+c.Height() {
		padLine := strings.Repeat(" ", c.Width())
		for curLine < c.Y()+c.Height() {
			lines = append(lines, padLine)
			curLine++
		}
	}

	return strings.Join(lines[c.Y():], "\n")
}

func (c *BufferedBlock) viewNoWrap(entries []*tapioca.Entry) string {
	if c.X()+c.Width() > c.maxLineWidth {
		c.ScrollToBegin()
		c.ScrollRight(c.maxLineWidth - c.Width())
	}

	wantedEntries := entries[c.Y():min(c.Y()+c.Height(), len(entries))]
	lines := make([]string, c.Height())
	for i := range wantedEntries {
		lines[i] = wantedEntries[i].StyledMove(c.X(), c.Width())
	}
	for i := len(wantedEntries); i < c.Height(); i++ {
		lines[i] = strings.Repeat(" ", c.Width())
	}

	return strings.Join(lines, "\n")
}

func (c *BufferedBlock) recomputeCachedInfo() {
	entries := c.entries.GetAll()
	c.recomputeLines(entries)
	c.recomputeMaxLineWidth(entries)
}

func (c *BufferedBlock) recomputeLines(entries []*tapioca.Entry) {
	if c.hScroll {
		// When horizontal scrolling is enabled, no wrapping occurs,
		// so virtual screen line count equals number of entries
		c.lines = c.entries.Size()
	} else {
		// When horizontal scrolling is disabled, entries wrap,
		// so virtual screen line count is total lines after wrapping
		c.lines = 0
		for _, e := range entries {
			c.lines += e.Lines(c.Width())
		}
	}

	// Virtual screen line count should be at least the physical screen height
	c.lines = max(c.lines, c.Height())
}

func (c *BufferedBlock) recomputeMaxLineWidth(entries []*tapioca.Entry) {
	c.maxLineWidth = c.Width()
	if !c.hScroll {
		return
	}

	for _, e := range entries {
		if l := e.Width(); c.maxLineWidth < l {
			c.maxLineWidth = l
		}
	}
}
