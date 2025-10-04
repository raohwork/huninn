// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ToppingTestSpec is a specification for testing a huninn component.
type ToppingTestSpec struct {
	Width  int
	Height int
	tea.Model
	EventLoopMaxTimes int
}

// IsThisTopping is a test helper to check if a given huninn component fulfills
// the basic requirements.
//
// It returns an empty string if all is good, otherwise a descriptive error message.
//
// Check [Component] for the requirements.
func IsThisTopping(spec ToppingTestSpec) string {
	max := max(0, spec.EventLoopMaxTimes)
	m, cmd := spec.Model.Update(ResizeMsg{Width: spec.Width, Height: spec.Height})
	for cmd != nil && max > 0 {
		m, cmd = m.Update(cmd())
		max--
	}

	got := m.View()

	// check rows (height)
	lines := strings.Split(got, "\n")
	if len(lines) != spec.Height {
		return fmt.Sprintf("expected %d lines, got %d lines", spec.Height, len(lines))
	}

	// check columns (width)
	for i, line := range lines {
		e := NewEntry(line)
		if e.Width() != spec.Width {
			if spec.Width == 1 && e.Width() == 2 && e.runeEndOffsets()[0] == 2 {
				continue // special case: a single wide rune
			}
			return fmt.Sprintf("line %d: expected width %d, got width %d (content: %q)", i, spec.Width, e.Width(), line)
		}
	}

	return "" // all good
}
