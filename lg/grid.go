// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package lg

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/raohwork/huninn/tapioca"
)

type gridSpec struct {
	x, y, w, h int
	comp       tea.Model
}

type gridMap struct {
	grid        [][]int
	cellWidths  []int
	cellHeights []int
}

func (m *gridMap) add(x, y, w, h, idx int) bool {
	for i := y; i < y+h && i < len(m.grid); i++ {
		for j := x; j < x+w && j < len(m.grid[i]); j++ {
			if m.grid[i][j] != -1 {
				return false
			}
		}
	}

	for i := y; i < y+h && i < len(m.grid); i++ {
		for j := x; j < x+w && j < len(m.grid[i]); j++ {
			m.grid[i][j] = idx
		}
	}
	return true
}

// GridLayout is a layout manager that arranges components in a grid.
//
// It separates the available space into a grid of cells and places components
// into these cells evenly based on their specified positions and spans.
//
// For example, a terminal screen of width 80 and height 24, with a grid layout
// of 4 columns and 3 rows, would divide the screen into cells of width 20 and
// height 8. A component added at position (0, 0) with a span of (2, 1) would
// occupy the area from (0, 0) to (39, 7) inclusive.
type GridLayout struct {
	components []gridSpec
	w, h       int
	hasError   bool
	*gridMap
}

// NewGridLayout creates a new GridLayout with the specified number of columns (w)
// and rows (h). Both w and h must be greater than zero.
func NewGridLayout(w, h int) *GridLayout {
	if w*h == 0 {
		panic("GridLayout requires both width and height to be greater than zero")
	}
	ret := &GridLayout{
		w: w,
		h: h,
	}
	ret.gridMap = newGridMap(w, h)
	return ret
}

func newGridMap(w, h int) *gridMap {
	grid := make([][]int, h)
	for i := range grid {
		grid[i] = make([]int, w)
		for idx := range grid[i] {
			grid[i][idx] = -1
		}
	}
	return &gridMap{grid: grid}
}

func (g *GridLayout) Add(comp tea.Model, x, y, w, h int) bool {
	if !g.gridMap.add(x, y, w, h, len(g.components)) {
		return false
	}

	g.components = append(g.components, gridSpec{x: x, y: y, w: w, h: h, comp: comp})
	return true
}

func (g *GridLayout) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range g.components {
		if cmd := c.comp.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	if len(cmds) > 0 {
		return tea.Batch(cmds...)
	}
	return nil
}

func (g *GridLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if c := g.handleResize(msg.Width, msg.Height); len(c) > 0 {
			cmds = append(cmds, c...)
		}
	case tapioca.ResizeMsg:
		if c := g.handleResize(msg.Width, msg.Height); len(c) > 0 {
			cmds = append(cmds, c...)
		}
	default:
		for i, c := range g.components {
			newComp, cmd := c.comp.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			g.components[i].comp = newComp
		}
	}
	if len(cmds) > 0 {
		return g, tea.Batch(cmds...)
	}
	return g, nil
}

// break down the new width and height into grid cells and send
// resize messages to each component
func (g *GridLayout) handleResize(w, h int) []tea.Cmd {
	// edge cases:
	//  - g.w or g.h is zero: it will be checked at NewGridLayout, ignore here
	//  - w or h is zero: create a cache field in GridLayout to indicate
	//    that something is wrong, cannot render normally.
	//  - remainder is distributed to columns/rows one by one from left to right,
	//    top to bottom.

	// handle zero terminal size
	if w <= 2 || h < 1 {
		g.hasError = true
		return nil
	}

	// compute cell widths and heights
	g.gridMap.cellWidths = make([]int, g.w)
	g.gridMap.cellHeights = make([]int, g.h)

	// distribute width with remainder going to leftmost cells
	baseWidth := w / g.w
	widthRemainder := w % g.w
	for i := 0; i < g.w; i++ {
		g.gridMap.cellWidths[i] = baseWidth
		if i < widthRemainder {
			g.gridMap.cellWidths[i]++
		}
	}

	// distribute height with remainder going to topmost cells
	baseHeight := h / g.h
	heightRemainder := h % g.h
	for i := 0; i < g.h; i++ {
		g.gridMap.cellHeights[i] = baseHeight
		if i < heightRemainder {
			g.gridMap.cellHeights[i]++
		}
	}

	// check minimum size requirement and collect resize commands
	var cmds []tea.Cmd
	for i, spec := range g.components {
		// calculate component's actual pixel size
		compWidth := 0
		for col := spec.x; col < spec.x+spec.w && col < len(g.gridMap.cellWidths); col++ {
			compWidth += g.gridMap.cellWidths[col]
		}
		compHeight := 0
		for row := spec.y; row < spec.y+spec.h && row < len(g.gridMap.cellHeights); row++ {
			compHeight += g.gridMap.cellHeights[row]
		}

		// check minimum size (2, 1)
		if compWidth < 2 || compHeight < 1 {
			g.hasError = true
			return nil
		}

		// send resize message to component
		newComp, cmd := spec.comp.Update(tapioca.ResizeMsg{Width: compWidth, Height: compHeight})
		g.components[i].comp = newComp
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// reset error flag if we got this far
	g.hasError = false
	return cmds
}

func (g *GridLayout) View() string {
	if g.hasError {
		return "Terminal size too small"
	}

	// If no components or cell dimensions not set, return blank
	if len(g.components) == 0 || len(g.gridMap.cellWidths) == 0 || len(g.gridMap.cellHeights) == 0 {
		return ""
	}

	// 1. Prerender all components
	componentLines := g.prerenderAllComponents()

	// 2. Organize component lines by grid rows
	gridRows := g.organizeComponentLinesByGridRows(componentLines)

	// 3. Calculate total capacity and preallocate
	totalCapacity := g.calculateTotalCapacity()
	result := strings.Builder{}
	result.Grow(totalCapacity)

	// 4. Write sequentially by grid rows
	for gridRow := 0; gridRow < g.h; gridRow++ {
		cellHeight := g.gridMap.cellHeights[gridRow]
		for physicalRow := 0; physicalRow < cellHeight; physicalRow++ {
			// Write all cells in this row
			for gridCol := 0; gridCol < g.w; gridCol++ {
				g.writeCell(&result, gridRows[gridRow][gridCol], physicalRow, g.gridMap.cellWidths[gridCol])
			}
			// Add newline (except for the last line)
			if gridRow < g.h-1 || physicalRow < cellHeight-1 {
				result.WriteByte('\n')
			}
		}
	}

	return result.String()
}

// prerenderAllComponents prerenders all components and organizes lines by component index
func (g *GridLayout) prerenderAllComponents() [][]string {
	componentLines := make([][]string, len(g.components))

	for i, spec := range g.components {
		content := spec.comp.View()
		if content == "" {
			// Empty content, create appropriate sized empty lines
			compHeight := g.calculateComponentHeight(spec)
			componentLines[i] = make([]string, compHeight)
			for j := range componentLines[i] {
				componentLines[i][j] = ""
			}
		} else {
			componentLines[i] = strings.Split(content, "\n")
		}
	}

	return componentLines
}

// organizeComponentLinesByGridRows reorganizes component lines by grid rows
func (g *GridLayout) organizeComponentLinesByGridRows(componentLines [][]string) [][][]string {
	// gridRows[gridRow][gridCol] = all lines of component in that cell
	gridRows := make([][][]string, g.h)
	for gridRow := range gridRows {
		gridRows[gridRow] = make([][]string, g.w)
	}

	// Traverse grid and assign corresponding component lines to each cell
	for gridRow := 0; gridRow < g.h; gridRow++ {
		for gridCol := 0; gridCol < g.w; gridCol++ {
			compIdx := g.gridMap.grid[gridRow][gridCol]
			if compIdx == -1 {
				// Empty cell, create empty lines
				cellHeight := g.gridMap.cellHeights[gridRow]
				gridRows[gridRow][gridCol] = make([]string, cellHeight)
				for i := range gridRows[gridRow][gridCol] {
					gridRows[gridRow][gridCol][i] = ""
				}
			} else {
				// Cell with component, calculate relative position within component
				spec := g.components[compIdx]

				// Calculate relative row position of this cell within component
				relativeRow := gridRow - spec.y
				startLineInComponent := 0
				for r := 0; r < relativeRow; r++ {
					if spec.y+r < len(g.gridMap.cellHeights) {
						startLineInComponent += g.gridMap.cellHeights[spec.y+r]
					}
				}

				// Extract lines corresponding to this cell
				cellHeight := g.gridMap.cellHeights[gridRow]
				cellLines := make([]string, cellHeight)

				for i := 0; i < cellHeight; i++ {
					lineIdx := startLineInComponent + i
					if lineIdx < len(componentLines[compIdx]) {
						cellLines[i] = componentLines[compIdx][lineIdx]
					} else {
						cellLines[i] = ""
					}
				}

				gridRows[gridRow][gridCol] = cellLines
			}
		}
	}

	return gridRows
}

// calculateTotalCapacity calculates the total capacity of the final string
func (g *GridLayout) calculateTotalCapacity() int {
	totalWidth := 0
	for _, w := range g.gridMap.cellWidths {
		totalWidth += w
	}

	totalHeight := 0
	for _, h := range g.gridMap.cellHeights {
		totalHeight += h
	}

	// width * height + number of newlines
	return totalWidth*totalHeight + totalHeight - 1
}

// writeCell writes a single line of a cell
func (g *GridLayout) writeCell(result *strings.Builder, cellLines []string, physicalRow, cellWidth int) {
	var line string
	if physicalRow < len(cellLines) {
		line = cellLines[physicalRow]
	} else {
		line = ""
	}

	// Truncate or pad to correct width
	if len(line) >= cellWidth {
		result.WriteString(line[:cellWidth])
	} else {
		result.WriteString(line)
		// Fill remaining space
		if remaining := cellWidth - len(line); remaining > 0 {
			result.WriteString(strings.Repeat(" ", remaining))
		}
	}
}

// calculateComponentHeight calculates the actual height of a component
func (g *GridLayout) calculateComponentHeight(spec gridSpec) int {
	height := 0
	for row := spec.y; row < spec.y+spec.h && row < len(g.gridMap.cellHeights); row++ {
		height += g.gridMap.cellHeights[row]
	}
	return height
}
