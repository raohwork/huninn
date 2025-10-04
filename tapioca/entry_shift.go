// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

// StyledShift 支援帶樣式的水平捲動功能
func (e *Entry) StyledShift(startCol, width int) string {
	width = max(1, width)
	startCol = max(0, startCol)
	w := e.Width()
	if startCol+width > w {
		startCol = max(0, w-width)
		width = w - startCol
	}

	return e.styledSubstring(startCol, startCol+width)
}
