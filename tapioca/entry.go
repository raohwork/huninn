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

// NewEntry creates a new Entry from the given string, parsing ANSI styles
// and handling East Asian wide characters.
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

// Entry represents a text entry with styled runes and provides methods to
// compute its display width, handle line wrapping, and retrieve styled substrings.
type Entry struct {
	styledData []StyledRune
	f          func() []int
}

// runeEndOffsets returns the cumulative display widths at each rune position
// Example: "ab你好cd" -> []int{1,2,4,6,7,8} (cumulative widths)
func (e *Entry) runeEndOffsets() []int {
	return e.f()
}

// Width returns the display width of the entry (considering East Asian wide characters)
func (e *Entry) Width() int {
	offsets := e.runeEndOffsets()
	if len(offsets) == 0 {
		return 0
	}
	return offsets[len(offsets)-1]
}

// String returns the plain text representation of the entry (without styles)
func (e *Entry) String() string {
	b := &strings.Builder{}
	b.Grow(len(e.styledData))
	for _, sr := range e.styledData {
		b.WriteRune(sr.Rune)
	}
	return b.String()
}

// StyledString returns the styled text representation of the entry
func (e *Entry) StyledString() string {
	return e.styledSubstring(0, len(e.styledData))
}

// Lines returns the number of lines the entry would occupy when wrapped at the given width
func (e *Entry) Lines(width int) int {
	return len(e.computeWarpPoints(width))
}

// RuneWidth returns the display width of a rune, considering East Asian wide characters.
func RuneWidth(r rune) int {
	props := width.LookupRune(r)
	if props.Kind() == width.EastAsianWide || props.Kind() == width.EastAsianFullwidth {
		return 2
	}
	return 1
}

func computeRuneEndOffsets(styledRunes []StyledRune) []int {
	// Pre-calculate cumulative widths for all runes
	// Example: "ab你好cd" -> []int{1,2,4,6,7,8} (cumulative widths)
	cumulativeWidths := make([]int, len(styledRunes))
	currentWidth := 0

	for i, sr := range styledRunes {
		runeWidth := RuneWidth(sr.Rune)

		currentWidth += runeWidth
		cumulativeWidths[i] = currentWidth
	}

	return cumulativeWidths
}

// styledSubstring returns the styled substring from start to end (exclusive)
//
// start and end are rune indices in e.styledData
func (e *Entry) styledSubstring(start, end int) string {
	if start >= end {
		return ""
	}
	if len(e.styledData) == 0 {
		return ""
	}

	b := &strings.Builder{}
	var lastStyle *style

	for i := start; i < end; i++ {
		sr := e.styledData[i]
		if sr.Style != lastStyle {
			// Use the Render method to properly transition between styles
			if lastStyle == nil {
				// First style in the substring
				b.WriteString(sr.Style.String())
			} else {
				b.WriteString(sr.Style.Render(lastStyle))
			}
			lastStyle = sr.Style
		}
		b.WriteRune(sr.Rune)
	}

	// Only append reset if we have any styling
	if lastStyle != nil && !lastStyle.isEmpty() {
		b.WriteString("\x1b[0m")
	}

	return b.String()
}

func (e *Entry) StyledLines(width int) []string {
	return e.StyledWarps(width)
}
