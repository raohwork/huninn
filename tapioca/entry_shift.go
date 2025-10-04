// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"sort"
	"strings"
)

// StyledMove returns a styled substring of the entry, as if viewport
// is moving over the entry.
//
// startCol and width are in terms of display width, not string length. For
// example:
//
//	(0, 4): start from 0, and 4 columns wide
//	for "01三五七89" it returns "01三" (三 is 2 columns wide)
//	for "0123456789" it returns "0123"
//
//	(-1, 5): start from -1, and 5 columns wide
//	  for "01三五七89" it returns "  01三" (pad one space at beginning)
//	  for "0123456789" it returns " 01234" (pad one space at beginning)
//
//	(8, 4): start from 8, and 4 columns wide
//	  for "01三五七89" it returns "89  " (pad two space at end)
//	  for "0123456789" it returns "89  " (pad two space at end)
func (e *Entry) StyledMove(startCol, width int) string {
	// algo:
	//   1. compute if we have to pad spaces at beginning (startCol < 0)
	//   2. compute if we have to pad spaces at end (startCol+width > entry width)
	//   3. compute the substring to extract from the entry
	//   4. combine them together
	buf := &strings.Builder{}
	totalWidth := e.Width()
	prefix := 0
	suffix := 0

	if startCol < 0 {
		prefix = -startCol
		startCol = 0
		width -= prefix
	}
	if startCol+width > totalWidth {
		suffix = startCol + width - totalWidth
		width -= suffix
	}

	str := e.StyledShift(startCol, width)
	buf.Grow(len(str) + prefix + suffix)
	for i := 0; i < prefix; i++ {
		buf.WriteRune(' ')
	}
	buf.WriteString(str)
	for i := 0; i < suffix; i++ {
		buf.WriteRune(' ')
	}
	return buf.String()
}

// StyledShift returns a styled substring of the entry, starting at startCol
//
// startCol and width are in terms of display width, not string length. For
// example:
//
//	(0, 4): start from 0, and 4 columns wide
//	  for "01三五七89" it returns "01三" (三 is 2 columns wide)
//	  for "0123456789" it returns "0123"
//
//	(1, 4): start from 1, and 4 columns wide
//	  for "01三五七89" it returns "1三 " (五 is cutted)
//	  for "0123456789" it returns "1234"
//
// It will not shift across entry boundary, so if the requested position is
// larger than the entry width, it will return the whole entry.
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

func computeRuneSizeFromOffset(offsets []int, idx int) int {
	l := len(offsets)
	if idx < 0 || idx >= l {
		return 1
	}
	if idx == 0 {
		return offsets[0]
	}
	return offsets[idx] - offsets[idx-1]
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
	curSize := computeRuneSizeFromOffset(offsets, startIdx)
	// special case for width 1 with wide char at beginning
	if width == 1 && curSize > 1 {
		return startIdx, startIdx + 1, false, false
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
	curSize = computeRuneSizeFromOffset(offsets, endIdx)
	// wide        && cut
	if curSize > 1 && offsets[endIdx]-1 > endCol {
		endIdx--
		hasSuffix = true
	}
	endIdx++

	return
}
