// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import "strings"

type warpPoint struct {
	start     int
	end       int
	hasSuffix bool
}

func (e *Entry) computeWarpPoints(width int) []warpPoint {
	if width < 1 {
		return nil
	}
	w := e.Width()
	offsets := e.runeEndOffsets()
	n := len(offsets)
	if w <= width {
		return []warpPoint{{0, n, false}}
	}

	startPos := 0

	ret := make([]warpPoint, 0, w/width+1)

	for startPos < w {
		start, end, _, hasSuffix := computeStartAndEndForShift(
			offsets,
			startPos,
			width,
		)
		realWidth := offsets[end-1] - offsets[start] + computeRuneSizeFromOffset(offsets, start)
		if hasSuffix {
			realWidth++
		}
		ret = append(ret, warpPoint{start, end, hasSuffix})
		startPos += realWidth
		if hasSuffix {
			startPos-- // account for the space
		}
	}

	return ret
}

// StyledWarps returns the entry split into lines, each line wrapped at the given width
// without any padding.
//
// For width <= 0, it resets to 1.
//
// For width == 1 with wide characters, the lines contain wide character will have a
// width of 2.
func (e *Entry) StyledWarps(width int) []string {
	if len(e.styledData) < 1 {
		return []string{""}
	}

	width = max(1, width)
	points := e.computeWarpPoints(width)
	ret := make([]string, 0, len(points))
	for _, p := range points {
		if p.hasSuffix {
			ret = append(ret, e.styledSubstring(p.start, p.end)+" ")
		} else {
			ret = append(ret, e.styledSubstring(p.start, p.end))
		}
	}

	return ret
}

func (e *Entry) substringWidth(start, end int) int {
	offsets := e.runeEndOffsets()
	end = min(end-1, len(e.styledData))

	return offsets[end] - offsets[start] + computeRuneSizeFromOffset(offsets, start)
}

// StyledBlock returns the entry split into lines, each line wrapped at the given width
// and padded with spaces to ensure each line is exactly width characters wide
//
// For width <= 0, it resets to 1.
//
// For width == 1 with wide characters, the lines contain wide character will have a
// width of 2.
func (e *Entry) StyledBlock(width int) []string {
	if len(e.styledData) < 1 {
		return []string{""}
	}

	width = max(1, width)
	points := e.computeWarpPoints(width)
	ret := make([]string, 0, len(points))
	for _, p := range points[:len(points)-1] {
		if p.hasSuffix {
			ret = append(ret, e.styledSubstring(p.start, p.end)+" ")
		} else {
			ret = append(ret, e.styledSubstring(p.start, p.end))
		}
	}

	last := points[len(points)-1]
	lastLine := e.styledSubstring(last.start, last.end)
	buf := strings.Builder{}
	w := e.substringWidth(last.start, last.end)
	padding := max(width, w) - w
	buf.Grow(padding + len(lastLine))
	buf.WriteString(lastLine)
	for range padding {
		buf.WriteByte(' ')
	}
	ret = append(ret, buf.String())

	return ret
}
