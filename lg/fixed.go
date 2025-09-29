// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package lg

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
)

type fixedInfo struct {
	size int
	tea.Model
}

// FixedLayout is a layout that reserves a fixed amount of space for one component
// and gives the rest of the space to the other component.
type FixedLayout struct {
	reserve    int
	components [2]fixedInfo

	horizontal bool
	end        bool
}

// FixedLeftLayout reserves space on the left side of the layout.
func FixedLeftLayout(reserve int, left, right tea.Model) *FixedLayout {
	return newFixed(true, false, reserve, left, right)
}

// FixedRightLayout reserves space on the right side of the layout.
func FixedRightLayout(reserve int, left, right tea.Model) *FixedLayout {
	return newFixed(true, true, reserve, left, right)
}

// FixedTopLayout reserves space on the top side of the layout.
func FixedTopLayout(reserve int, top, bottom tea.Model) *FixedLayout {
	return newFixed(false, false, reserve, top, bottom)
}

// FixedBottomLayout reserves space on the bottom side of the layout.
func FixedBottomLayout(reserve int, top, bottom tea.Model) *FixedLayout {
	return newFixed(false, true, reserve, top, bottom)
}

type emptyLayout struct {
	w, h int
}

func (e *emptyLayout) Init() tea.Cmd { return nil }
func (e *emptyLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tapioca.ResizeMsg); ok {
		e.w = msg.Width
		e.h = msg.Height
	}
	return e, nil
}
func (e *emptyLayout) View() string {
	b := strings.Builder{}
	b.WriteString(strings.Repeat(" ", e.w))
	for i := 1; i < e.h; i++ {
		b.Write([]byte{'\n'})
		b.WriteString(strings.Repeat(" ", e.w))
	}
	return b.String()
}

// PadTop is a convenience function to add padding to the top of a component.
func PadTop(reserve int, component tea.Model) *FixedLayout {
	return FixedTopLayout(reserve, &emptyLayout{}, component)
}

// PadBottom is a convenience function to add padding to the bottom of a component.
func PadBottom(reserve int, component tea.Model) *FixedLayout {
	return FixedBottomLayout(reserve, component, &emptyLayout{})
}

// PadLeft is a convenience function to add padding to the left of a component.
func PadLeft(reserve int, component tea.Model) *FixedLayout {
	return FixedLeftLayout(reserve, &emptyLayout{}, component)
}

// PadRight is a convenience function to add padding to the right of a component.
func PadRight(reserve int, component tea.Model) *FixedLayout {
	return FixedRightLayout(reserve, component, &emptyLayout{})
}

func newFixed(horizontal, end bool, reserve int, components ...tea.Model) *FixedLayout {
	if len(components) != 2 {
		panic("fixed can only contain up to 2 components")
	}
	return &FixedLayout{
		reserve: reserve,
		components: [2]fixedInfo{
			{size: 0, Model: components[0]},
			{size: 0, Model: components[1]},
		},
		horizontal: horizontal,
		end:        end,
	}
}

func (f *FixedLayout) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, 2)
	for _, c := range f.components {
		cmd := c.Init()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	if len(cmds) == 0 {
		return nil
	}
	return tea.Batch(cmds...)
}

func (f *FixedLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmds = append(cmds, f.handleResize(msg.Width, msg.Height)...)
	case tapioca.ResizeMsg:
		cmds = append(cmds, f.handleResize(msg.Width, msg.Height)...)
	default:
		for i := range f.components {
			m, cmd := f.components[i].Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			f.components[i].Model = m
		}
	}
	if len(cmds) == 0 {
		return f, nil
	}
	return f, tea.Batch(cmds...)
}

func (f *FixedLayout) handleResize(w, h int) (ret []tea.Cmd) {
	newSize := h
	if f.horizontal {
		newSize = w
	}

	newMsg := func(s int) tapioca.ResizeMsg {
		if f.horizontal {
			return tapioca.ResizeMsg{Width: s, Height: h}
		}
		return tapioca.ResizeMsg{Width: w, Height: s}
	}

	reserveAt := 0
	if f.end {
		reserveAt = 1
	}
	if f.components[reserveAt].size != f.reserve {
		f.components[reserveAt].size = f.reserve
		m, cmd := f.components[reserveAt].Update(newMsg(f.reserve))
		if cmd != nil {
			ret = append(ret, cmd)
		}
		f.components[reserveAt].Model = m
	}

	rest := max(newSize-f.reserve, 0)

	f.components[1-reserveAt].size = rest
	if rest > 0 {
		m, cmd := f.components[1-reserveAt].Update(newMsg(rest))
		if cmd != nil {
			ret = append(ret, cmd)
		}
		f.components[1-reserveAt].Model = m
	}

	return ret
}

func (f *FixedLayout) View() string {
	if f.components[0].size == 0 || f.components[1].size == 0 {
		return "Terminal too small"
	}

	if f.horizontal {
		return renderHorizontal(f.components[0].Model.View(), f.components[1].Model.View())
	}

	return strings.TrimRight(f.components[0].Model.View(), "\n") + "\n" + f.components[1].Model.View()
}

func renderHorizontal(left, right string) string {
	left = strings.TrimRight(left, "\n")
	right = strings.TrimRight(right, "\n")
	bytes := len(left) + len(right)
	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")

	var b strings.Builder
	b.Grow(bytes)

	for i, l, r := 0, len(leftLines), len(rightLines); i < max(l, r); i++ {
		if i < l {
			b.WriteString(leftLines[i])
		} else {
			b.WriteString(strings.Repeat(" ", len(leftLines[0])))
		}
		if i < r {
			b.WriteString(rightLines[i])
		} else {
			b.WriteString(strings.Repeat(" ", len(rightLines[0])))
		}
	}
	return b.String()
}
