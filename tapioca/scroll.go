// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ScrollController defines the interface for controlling scrolling behavior
// in a UI component. IT DOES NOT TRIGGER A RENDER UPDATE.
type ScrollController interface {
	X() int
	Y() int
	Width() int
	Height() int
	ScrollUp(lines int)
	ScrollDown(lines int)
	ScrollToTop()
	ScrollToBottom()
	ScrollLeft(cols int)
	ScrollRight(cols int)
	ScrollToBegin()
	ScrollToEnd()
	ScrollTo(col, row int)
}

// Scrollable provides a basic implementation of the ScrollController interface.
//
// It handles ResizeMsg and Scroll*Msg messages to update its state.
type Scrollable struct {
	x, y, w, h int
	maxW, maxH func() int
}

func NewScrollable(maxW, maxH func() int) Scrollable {
	return Scrollable{
		maxW: maxW,
		maxH: maxH,
	}
}

func (s *Scrollable) X() int      { return s.x }
func (s *Scrollable) Y() int      { return s.y }
func (s *Scrollable) Width() int  { return s.w }
func (s *Scrollable) Height() int { return s.h }

func (s *Scrollable) ScrollUp(lines int) {
	s.y = max(0, s.y-max(lines, 0))
}
func (s *Scrollable) ScrollDown(lines int) {
	s.y = min(max(0, s.maxH()-s.h), s.y+max(lines, 0))
}
func (s *Scrollable) ScrollToTop() {
	s.y = 0
}
func (s *Scrollable) ScrollToBottom() {
	s.y = max(0, s.maxH()-s.h)
}
func (s *Scrollable) ScrollLeft(cols int) {
	s.x = max(0, s.x-max(cols, 0))
}
func (s *Scrollable) ScrollRight(cols int) {
	s.x = min(max(0, s.maxW()-s.w), s.x+max(cols, 0))
}
func (s *Scrollable) ScrollToBegin() {
	s.x = 0
}
func (s *Scrollable) ScrollToEnd() {
	s.x = max(0, s.maxW()-s.w)
}
func (s *Scrollable) ScrollTo(col, row int) {
	s.ScrollToBegin()
	s.ScrollToTop()
	s.ScrollRight(col)
	s.ScrollDown(row)
}

func (s *Scrollable) HandleEvent(msg tea.Msg) {
	switch m := msg.(type) {
	case ResizeMsg:
		s.w, s.h = m.Width, m.Height
		// recompute x, y to be in bounds
		s.x = min(max(0, s.maxW()-s.w), s.x)
		s.y = min(max(0, s.maxH()-s.h), s.y)
	case ScrollBeginMsg:
		s.ScrollToBegin()
	case ScrollEndMsg:
		s.ScrollToEnd()
	case ScrollLeftMsg:
		s.ScrollLeft(int(m))
	case ScrollRightMsg:
		s.ScrollRight(int(m))
	case ScrollTopMsg:
		s.ScrollToTop()
	case ScrollBottomMsg:
		s.ScrollToBottom()
	case ScrollUpMsg:
		s.ScrollUp(int(m))
	case ScrollDownMsg:
		s.ScrollDown(int(m))
	case ScrollToMsg:
		s.ScrollTo(m.X, m.Y)
	}
}
