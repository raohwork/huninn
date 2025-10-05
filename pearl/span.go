// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pearl

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
)

type SpanSetContentMsg struct {
	id   int64
	data string
}

// Span is a single line text box that can be used to display a single entry of text.
//
// It warps the content to fit the width of the box, and truncates the content to fit the
// height of the box.
type Span struct {
	id    int64
	entry *tapioca.Entry
	w, h  int
}

// NewSpan creates a new Span.
func NewSpan() *Span {
	return &Span{
		id:    tapioca.NewID(),
		entry: tapioca.NewEntry(""),
	}
}

// SetContent sets the content of the span.
//
// You should use it only when you are handling an event message.
func (s *Span) SetContent(data string) {
	arr := strings.Split(data, "\n")
	s.entry = tapioca.NewEntry(arr[0])
}

// Setter returns a function that can be used to set the content of the span by sending
// a SpanSetContentMsg to bubble tea program.
func (s *Span) Setter(f func(tea.Msg)) func(string) {
	return func(data string) {
		f(SpanSetContentMsg{
			id:   s.id,
			data: data,
		})
	}
}

// UpdateInto is identical to Update but returns a *Span instead of tea.Model to prevent
// type assertion.
func (s *Span) UpdateInto(msg tea.Msg) (*Span, tea.Cmd) {
	switch m := msg.(type) {
	case tapioca.ResizeMsg:
		s.w, s.h = m.Width, m.Height
	case SpanSetContentMsg:
		if m.id != s.id {
			return s, nil
		}
		s.SetContent(m.data)
	}
	return s, nil
}

func (s *Span) Init() tea.Cmd                           { return nil }
func (s *Span) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s.UpdateInto(msg) }

func (s *Span) View() string {
	lines := s.entry.StyledBlock(s.w)
	if len(lines) > s.h {
		return strings.Join(lines[:s.h], "\n")
	}

	for i := len(lines); i < s.h; i++ {
		lines = append(lines, strings.Repeat(" ", s.w))
	}
	return strings.Join(lines, "\n")
}
