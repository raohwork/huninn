// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pearl

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/raohwork/huninn/tapioca"
)

// Block is a component that displays a block of text with scrollable capabilities.
//
// Shrinking will lost some content.
type Block struct {
	id   uuid.UUID
	impl *tapioca.Component
}

// NewBlock creates a new Block component with horizontal and vertical scrolling enabled.
func NewBlock() *Block {
	return &Block{
		id: uuid.New(),
		impl: tapioca.NewComponent(
			1,
			tapioca.HorizontalScrollable(),
			tapioca.VerticalScrollable(),
		),
	}
}

// BlockSetContentMsg is a message type used to set the content of a Block component.
type BlockSetContentMsg struct {
	id   uuid.UUID
	data []string
}

// GetSetter returns a function that sends a BlockSetContentMsg to update the Block's content.
func (b *Block) Setter(send func(tea.Msg)) func(...string) {
	return func(s ...string) {
		send(BlockSetContentMsg{
			id:   b.id,
			data: s,
		})
	}
}

func (b *Block) Init() tea.Cmd { return nil }

func (b *Block) View() string { return b.impl.View() }

func (b *Block) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b.UpdateInto(msg)
}
func (b *Block) UpdateInto(msg tea.Msg) (*Block, tea.Cmd) {
	switch m := msg.(type) {
	case BlockSetContentMsg:
		if m.id == b.id {
			b.impl.Clear()
			for _, line := range m.data {
				b.impl.Append(line)
			}
		}
		return b, nil
	default:
		newImpl, cmd := b.impl.UpdateInto(msg)
		b.impl = newImpl
		return b, cmd
	}
}
