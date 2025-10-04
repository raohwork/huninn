// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

type warpPoint struct {
	start     int
	end       int
	hasSuffix bool
}

func (e *Entry) computeWarpPoints(width int) []warpPoint {
	w := e.Width()
	if w <= width {
		return []warpPoint{{0, w, false}}
	}

	offsets := e.runeEndOffsets()
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
