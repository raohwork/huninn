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
	id      int64
	entries []*tapioca.Entry
	tapioca.Scrollable
	maxWidth int
}

// NewBlock creates a new Block component with horizontal and vertical scrolling enabled.
func NewBlock() *Block {
	ret := &Block{
		id: tapioca.NewID(),
	}
	ret.Scrollable = tapioca.NewScrollable(
		func() int { return ret.maxWidth },
		func() int { return len(ret.entries) },
	)
	return ret
}

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
	lines := make([]string, 0, b.Height())
	for i := 0; i < b.Height() && i < len(b.entries); i++ {
		lines = append(lines, b.entries[i].StyledMove(b.X(), b.Width()))
	}
	for i := len(lines); i < b.Height(); i++ {
		lines = append(lines, strings.Repeat(" ", b.Width()))
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
	default:
		b.HandleEvent(msg)
	}
	return b, nil
}
