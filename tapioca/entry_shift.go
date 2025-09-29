// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"sort"
	"strings"
	"unicode/utf8"
)

func (e *Entry) Shift(startCol, width int) string {
	runeEndOffsets := e.runeEndOffsets()
	if len(runeEndOffsets) == 0 {
		return ""
	}
	if width < 2 {
		width = 2
	}

	// Check if begin is beyond all characters
	cl := len(runeEndOffsets)
	if startCol >= runeEndOffsets[cl-1] {
		return ""
	}

	// find left boundary
	leftRuneIdx := 0
	var hasPrefix bool

	if startCol > 0 {
		leftRuneIdx = sort.Search(cl, func(i int) bool { return runeEndOffsets[i] >= startCol })

		if leftRuneIdx >= cl {
			if hasPrefix {
				return " "
			}
			return ""
		}

		if runeEndOffsets[leftRuneIdx] != startCol {
			hasPrefix = true
		}
		leftRuneIdx++
	}

	// find right boundry
	rightRuneIdx := sort.Search(cl, func(i int) bool { return runeEndOffsets[i] >= startCol+width })
	if rightRuneIdx >= cl {
		rightRuneIdx = cl - 1
	} else if runeEndOffsets[rightRuneIdx] != startCol+width {
		rightRuneIdx--
	}

	size := 0
	for i := leftRuneIdx; i <= rightRuneIdx; i++ {
		size += utf8.RuneLen(e.styledData[i].Rune)
	}

	buf := &strings.Builder{}
	buf.Grow(size + 1)

	if hasPrefix {
		buf.WriteByte(' ')
	}
	for i := leftRuneIdx; i <= rightRuneIdx; i++ {
		buf.WriteRune(e.styledData[i].Rune)
	}

	return buf.String()
}

// StyledShift 支援帶樣式的水平捲動功能
func (e *Entry) StyledShift(startCol, width int) string {
	runeEndOffsets := e.runeEndOffsets()
	if len(runeEndOffsets) == 0 {
		return ""
	}
	if width < 2 {
		width = 2
	}

	// Check if begin is beyond all characters
	cl := len(runeEndOffsets)
	if startCol >= runeEndOffsets[cl-1] {
		return ""
	}

	// find left boundary
	leftRuneIdx := 0
	var hasPrefix bool

	if startCol > 0 {
		leftRuneIdx = sort.Search(cl, func(i int) bool { return runeEndOffsets[i] >= startCol })

		if leftRuneIdx >= cl {
			if hasPrefix {
				return " "
			}
			return ""
		}

		if runeEndOffsets[leftRuneIdx] != startCol {
			hasPrefix = true
		}
		leftRuneIdx++
	}

	// find right boundry
	rightRuneIdx := sort.Search(cl, func(i int) bool { return runeEndOffsets[i] >= startCol+width })
	if rightRuneIdx >= cl {
		rightRuneIdx = cl - 1
	} else if runeEndOffsets[rightRuneIdx] != startCol+width {
		rightRuneIdx--
	}

	// Use styledSubstring logic for styled output
	start := leftRuneIdx
	end := rightRuneIdx + 1

	buf := &strings.Builder{}
	lastStyle := (*style)(nil)

	// Handle prefix space with initial style
	if hasPrefix {
		if start < len(e.styledData) {
			initialStyle := e.styledData[start].Style
			buf.WriteString(initialStyle.String())
			lastStyle = initialStyle
		}
		buf.WriteByte(' ')
	}

	// Apply styles for the content
	for i := start; i < end && i < len(e.styledData); i++ {
		sr := e.styledData[i]
		if sr.Style != lastStyle {
			buf.WriteString(sr.Style.String())
			lastStyle = sr.Style
		}
		buf.WriteRune(sr.Rune)
	}

	// Append a reset code at the end
	buf.WriteString("\x1b[m")
	return buf.String()
}
