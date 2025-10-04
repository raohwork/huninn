// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"sort"
	"strings"
)

// StyledShift returns a styled substring of the entry, starting at startCol
//
// startCol and width are in terms of display width, not string length.
func (e *Entry) StyledShift(startCol, width int) string {
	ret, _, _ := e.styledShift(startCol, width)
	return ret
}
func (e *Entry) styledShift(startCol, width int) (string, bool, bool) {
	if len(e.styledData) == 0 {
		return "", false, false
	}
	offsets := e.runeEndOffsets()
	totalWidth := offsets[len(offsets)-1]

	// Normalize
	width = min(max(1, width), totalWidth)
	startCol = max(0, startCol)
	if startCol+width > totalWidth {
		startCol = totalWidth - width
	}

	// If width covers everything, return everything
	if startCol == 0 && width == totalWidth {
		return e.StyledString(), false, false
	}

	buf := &strings.Builder{}
	startIdx, endIdx, hasPrefix, hasSuffix := computeStartAndEndForShift(offsets, startCol, width)
	ret := e.styledSubstring(startIdx, endIdx)

	// output
	total := len(ret)
	if hasPrefix {
		total++
	}
	if hasSuffix {
		total++
	}
	buf.Grow(total)
	if hasPrefix {
		buf.WriteRune(' ')
	}
	buf.WriteString(ret)
	if hasSuffix {
		buf.WriteRune(' ')
	}

	return buf.String(), hasPrefix, hasSuffix
}

func computeStartAndEndForShift(offsets []int, startCol, width int) (startIdx, endIdx int, hasPrefix, hasSuffix bool) {
	n := len(offsets)

	// "01三五七89"
	//  123456789A  position
	//  0123456789  position-1
	//  01 2 3 456  rune index
	//
	// for wide chars
	//   if not cut (eg. startCol=2), offset[startIdx] should larger than startCol
	//   if cut (eg. startCol=3), offset[startIdx] should be equal to startCol
	// for narrow chars, never cut
	startIdx = sort.Search(n, func(i int) bool {
		return offsets[i]-1 >= startCol
	})
	curSize := offsets[startIdx]
	if startIdx > 0 {
		curSize -= offsets[startIdx-1]
	}
	// wide        && cut
	if curSize > 1 && offsets[startIdx]-1 == startCol {
		startIdx++
		hasPrefix = true
	}

	// "01三五七89"
	//  123456789A  position
	//  0123456789  position-1
	//  01 2 3 456  rune index
	//
	// endCol is inclusive here
	// for wide chars
	//   if not cut (eg. endCol=7), offset[endIdx] should equal endCol
	//   if cut (eg. endCol=6), offset[endIdx] should be larger than endCol
	// for narrow chars, never cut
	endCol := startCol + width - 1 // inclusive
	endIdx = sort.Search(n, func(i int) bool {
		return offsets[i]-1 >= endCol
	})
	endIdx = min(endIdx, n-1)
	curSize = offsets[endIdx]
	if endIdx > 0 {
		curSize -= offsets[endIdx-1]
	}
	// wide        && cut
	if curSize > 1 && offsets[endIdx]-1 > endCol {
		endIdx--
		hasSuffix = true
	}
	endIdx++

	return
}
