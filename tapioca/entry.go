// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"strings"
	"unicode/utf8"

	"github.com/raohwork/task/noctx/ncaction"
	"golang.org/x/text/width"
)

type StyledRune struct {
	Rune  rune
	Style *style
}

func NewEntry(data string) *Entry {
	// First, clean unsupported ANSI CSI sequences
	data = ansiOtherRegex.ReplaceAllString(data, "")

	styledData := make([]StyledRune, 0, len(data))
	currentStyle := &style{} // Start with a default/reset style

	i := 0
	for i < len(data) {
		// Check for ANSI escape code
		if data[i] == '\x1b' && i+1 < len(data) && data[i+1] == '[' {
			m := ansiStyleRegex.FindStringIndex(data[i:])
			if m != nil {
				code := data[i : i+m[1]]
				// parseAnsiCode should be implemented in entry_ansi_style.go
				// It must return a new style object, not modify the old one.
				currentStyle = parseAnsiCode(code, currentStyle)
				i += m[1]
				continue
			}
		}

		// Handle a regular rune
		r, size := utf8.DecodeRuneInString(data[i:])
		styledData = append(styledData, StyledRune{Rune: r, Style: currentStyle})
		i += size
	}

	return &Entry{
		styledData: styledData,
		f:          ncaction.NoErrGet(computeRuneEndOffsets).By(styledData).Cached().NoErr(),
	}
}

type Entry struct {
	styledData []StyledRune
	f          func() []int
}

func (e *Entry) runeEndOffsets() []int {
	return e.f()
}

func (e *Entry) Len() int {
	return len(e.styledData)
}

func (e *Entry) String() string {
	b := &strings.Builder{}
	b.Grow(len(e.styledData))
	for _, sr := range e.styledData {
		b.WriteRune(sr.Rune)
	}
	return b.String()
}

func (e *Entry) Warps(width int) []string {
	return e.warpLineTextWidth(width)
}

func (e *Entry) warpPositions(width int) []int {
	if len(e.styledData) == 0 {
		return nil
	}
	if width <= 2 {
		width = 2

	}

	// 1,2,3,4,5,6,7,9,11
	charPos := e.runeEndOffsets()
	ret := make([]int, 1, len(charPos)/width+2)
	ret[0] = 0
	headPos := 0
	for idx := 1; idx < len(charPos); idx++ {
		pos := charPos[idx]
		if pos-headPos > width {
			headPos = charPos[idx-1]
			ret = append(ret, idx)
		}
	}
	return ret
}

func (e *Entry) Lines(width int) int {
	return max(1, len(e.warpPositions(width)))
}

func computeRuneEndOffsets(styledRunes []StyledRune) []int {
	// Pre-calculate cumulative widths for all runes
	// Example: "ab你好cd" -> []int{1,2,4,6,7,8} (cumulative widths)
	cumulativeWidths := make([]int, len(styledRunes))
	currentWidth := 0

	for i, sr := range styledRunes {
		props := width.LookupRune(sr.Rune)
		runeWidth := 1 // default width for most characters
		if props.Kind() == width.EastAsianWide || props.Kind() == width.EastAsianFullwidth {
			runeWidth = 2
		}

		currentWidth += runeWidth
		cumulativeWidths[i] = currentWidth
	}

	return cumulativeWidths
}

// warpLineTextWidth implements text wrapping using golang.org/x/text/width
// for efficient width calculations. It pre-calculates cumulative widths to
// reduce memory allocations and improve performance.
func (e *Entry) warpLineTextWidth(lineWidth int) []string {
	indexes := e.warpPositions(lineWidth)
	if len(indexes) <= 1 {
		return []string{e.String()}
	}

	var result []string
	for l, idx := len(indexes), 0; idx < l-1; idx++ {
		// Extract runes from styledData and convert to string
		runes := make([]rune, 0, indexes[idx+1]-indexes[idx])
		for i := indexes[idx]; i < indexes[idx+1]; i++ {
			runes = append(runes, e.styledData[i].Rune)
		}
		result = append(result, string(runes))
	}
	// Handle the last segment
	runes := make([]rune, 0, len(e.styledData)-indexes[len(indexes)-1])
	for i := indexes[len(indexes)-1]; i < len(e.styledData); i++ {
		runes = append(runes, e.styledData[i].Rune)
	}
	result = append(result, string(runes))

	return result
}

// styledSubstring 產生帶樣式的字串片段
func (e *Entry) styledSubstring(start, end int) string {
	if start >= end {
		return ""
	}

	b := &strings.Builder{}
	lastStyle := (*style)(nil) // Use nil to force initial style print

	// Ensure the line starts with the correct style
	initialStyle := e.styledData[start].Style
	b.WriteString(initialStyle.String())
	lastStyle = initialStyle

	for i := start; i < end; i++ {
		sr := e.styledData[i]
		if sr.Style != lastStyle {
			b.WriteString(sr.Style.String())
			lastStyle = sr.Style
		}
		b.WriteRune(sr.Rune)
	}

	// Append a reset code at the end of the line
	b.WriteString("\x1b[m")
	return b.String()
}

// StyledLines 產生帶樣式的字串陣列
func (e *Entry) StyledLines(width int) []string {
	indexes := e.warpPositions(width)
	if len(indexes) <= 1 {
		// No wrapping needed, but still need to apply style
		return []string{e.styledSubstring(0, len(e.styledData))}
	}

	var result []string
	for i := 0; i < len(indexes)-1; i++ {
		result = append(result, e.styledSubstring(indexes[i], indexes[i+1]))
	}
	// Append the last line
	result = append(result, e.styledSubstring(indexes[len(indexes)-1], len(e.styledData)))

	return result
}
