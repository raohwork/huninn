// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pearl

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
)

// Block is a component that displays lines of text with scrollable capabilities.
type Block struct {
	id         int64
	entries    []*tapioca.Entry
	x, y, w, h int
	maxWidth   int
}

// NewBlock creates a new Block component with horizontal and vertical scrolling enabled.
func NewBlock() *Block {
	return &Block{
		id: tapioca.NewID(),
	}
}

func (b *Block) Width() int  { return b.w }
func (b *Block) Height() int { return b.h }
func (b *Block) X() int      { return b.x }
func (b *Block) Y() int      { return b.y }

// BlockSetContentMsg is a message type used to set the content of a Block component.
type BlockSetContentMsg struct {
	id   int64
	data []string
}

// SetContent sets the content of the Block to the provided lines of text.
//
// You should use it only when you are handling an event message.
func (b *Block) SetContent(data ...string) {
	b.entries = make([]*tapioca.Entry, len(data))
	b.maxWidth = 0
	for i, line := range data {
		b.entries[i] = tapioca.NewEntry(line)
		b.maxWidth = max(b.maxWidth, b.entries[i].Width())
	}
}

// Setter returns a function that sends a BlockSetContentMsg to update the Block's content.
func (b *Block) Setter(send func(tea.Msg)) func(...string) {
	return func(s ...string) {
		send(BlockSetContentMsg{
			id:   b.id,
			data: s,
		})
	}
}

func (b *Block) Init() tea.Cmd { return nil }

func (b *Block) View() string {
	lines := make([]string, 0, b.h)
	for i := 0; i < b.h && i < len(b.entries); i++ {
		lines = append(lines, b.entries[i].StyledMove(b.x, b.w))
	}
	for i := len(lines); i < b.h; i++ {
		lines = append(lines, strings.Repeat(" ", b.w))
	}
	return strings.Join(lines, "\n")
}

func (b *Block) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b.UpdateInto(msg)
}
func (b *Block) UpdateInto(msg tea.Msg) (*Block, tea.Cmd) {
	switch m := msg.(type) {
	case BlockSetContentMsg:
		if m.id != b.id {
			return b, nil
		}
		b.SetContent(m.data...)
	case tapioca.ResizeMsg:
		b.w, b.h = m.Width, m.Height
		// recompute x, y to be in bounds
		b.x = min(max(0, b.maxWidth-b.w), b.x)
		b.y = min(max(0, len(b.entries)-b.h), b.y)
	case tapioca.ScrollBeginMsg:
		b.x = 0
	case tapioca.ScrollEndMsg:
		b.x = max(0, b.maxWidth-b.w)
	case tapioca.ScrollLeftMsg:
		b.x = max(0, b.x-1)
	case tapioca.ScrollRightMsg:
		b.x = min(max(0, b.maxWidth-b.w), b.x+1)
	case tapioca.ScrollTopMsg:
		b.y = 0
	case tapioca.ScrollBottomMsg:
		b.y = max(0, len(b.entries)-b.h)
	case tapioca.ScrollUpMsg:
		b.y = max(0, b.y-1)
	case tapioca.ScrollDownMsg:
		b.y = min(max(0, len(b.entries)-b.h), b.y+1)
	}
	return b, nil
}
