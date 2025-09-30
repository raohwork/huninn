// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package lg

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockComponent for testing resize messages
type MockComponent struct {
	mock.Mock
}

func (m *MockComponent) Init() tea.Cmd {
	args := m.Called()
	return args.Get(0).(tea.Cmd)
}

func (m *MockComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	args := m.Called(msg)
	cmd := args.Get(1)
	if cmd == nil {
		return args.Get(0).(tea.Model), nil
	}
	return args.Get(0).(tea.Model), cmd.(tea.Cmd)
}

func (m *MockComponent) View() string {
	args := m.Called()
	return args.String(0)
}

func TestGridLayout_HandleResize(t *testing.T) {
	tests := []struct {
		name                 string
		gridW, gridH         int
		terminalW, terminalH int
		components           []struct{ x, y, w, h int }
		expectedError        bool
		expectedCellWidths   []int
		expectedCellHeights  []int
		expectedResizeMsg    []tapioca.ResizeMsg
	}{
		{
			name:                "normal case - even division",
			gridW:               2,
			gridH:               2,
			terminalW:           20,
			terminalH:           10,
			components:          []struct{ x, y, w, h int }{{0, 0, 1, 1}, {1, 1, 1, 1}},
			expectedError:       false,
			expectedCellWidths:  []int{10, 10},
			expectedCellHeights: []int{5, 5},
			expectedResizeMsg:   []tapioca.ResizeMsg{{Width: 10, Height: 5}, {Width: 10, Height: 5}},
		},
		{
			name:                "remainder distribution - width example (82÷5)",
			gridW:               5,
			gridH:               1,
			terminalW:           82,
			terminalH:           10,
			components:          []struct{ x, y, w, h int }{{0, 0, 5, 1}},
			expectedError:       false,
			expectedCellWidths:  []int{17, 17, 16, 16, 16}, // 82÷5=16 remainder 2, distributed to first 2 cells
			expectedCellHeights: []int{10},
			expectedResizeMsg:   []tapioca.ResizeMsg{{Width: 82, Height: 10}},
		},
		{
			name:                "remainder distribution - height example",
			gridW:               1,
			gridH:               3,
			terminalW:           10,
			terminalH:           8,
			components:          []struct{ x, y, w, h int }{{0, 0, 1, 3}},
			expectedError:       false,
			expectedCellWidths:  []int{10},
			expectedCellHeights: []int{3, 3, 2}, // 8÷3=2 remainder 2, distributed to first 2 rows
			expectedResizeMsg:   []tapioca.ResizeMsg{{Width: 10, Height: 8}},
		},
		{
			name:                "spanning components",
			gridW:               3,
			gridH:               2,
			terminalW:           15,
			terminalH:           6,
			components:          []struct{ x, y, w, h int }{{0, 0, 2, 1}, {2, 0, 1, 2}},
			expectedError:       false,
			expectedCellWidths:  []int{5, 5, 5},
			expectedCellHeights: []int{3, 3},
			expectedResizeMsg:   []tapioca.ResizeMsg{{Width: 10, Height: 3}, {Width: 5, Height: 6}},
		},
		{
			name:          "terminal width too small",
			gridW:         2,
			gridH:         2,
			terminalW:     2, // <= 2
			terminalH:     10,
			components:    []struct{ x, y, w, h int }{{0, 0, 1, 1}},
			expectedError: true,
		},
		{
			name:          "terminal height too small",
			gridW:         2,
			gridH:         2,
			terminalW:     10,
			terminalH:     0, // < 1
			components:    []struct{ x, y, w, h int }{{0, 0, 1, 1}},
			expectedError: true,
		},
		{
			name:          "component size too small - width",
			gridW:         10,
			gridH:         1,
			terminalW:     15, // 15÷10=1 remainder 5, first 5 cells get width 2, last 5 get width 1
			terminalH:     5,
			components:    []struct{ x, y, w, h int }{{6, 0, 1, 1}}, // 6th cell has width 1 < 2
			expectedError: true,
		},
		{
			name:          "component size too small - height",
			gridW:         1,
			gridH:         10,
			terminalW:     10,
			terminalH:     5,                                        // 5÷10=0 remainder 5, first 5 rows get height 1, last 5 get height 0
			components:    []struct{ x, y, w, h int }{{0, 6, 1, 1}}, // 6th row has height 0 < 1
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create grid layout
			grid := NewGridLayout(tt.gridW, tt.gridH)

			// Create mock components
			var mockComps []*MockComponent
			for i, comp := range tt.components {
				mockComp := &MockComponent{}
				mockComps = append(mockComps, mockComp)

				if !tt.expectedError && i < len(tt.expectedResizeMsg) {
					// Expect the resize message
					mockComp.On("Update", tt.expectedResizeMsg[i]).Return(mockComp, nil)
				}

				// Add component to grid
				grid.Add(mockComp, comp.x, comp.y, comp.w, comp.h)
			}

			// Call handleResize
			cmds := grid.handleResize(tt.terminalW, tt.terminalH)

			// Check error state
			assert.Equal(t, tt.expectedError, grid.hasError, "hasError should match expected")

			if tt.expectedError {
				assert.Nil(t, cmds, "commands should be nil when error occurs")
				// For error cases, we don't need to check other expectations
				return
			}

			// Check cell dimensions cache
			assert.Equal(t, tt.expectedCellWidths, grid.gridMap.cellWidths, "cellWidths should match expected")
			assert.Equal(t, tt.expectedCellHeights, grid.gridMap.cellHeights, "cellHeights should match expected")

			// Verify all mock expectations were met
			for _, mockComp := range mockComps {
				mockComp.AssertExpectations(t)
			}
		})
	}
}

func TestGridLayout_HandleResize_CellSizeCalculation(t *testing.T) {
	// Test specific edge cases for cell size calculation
	tests := []struct {
		name                 string
		terminalW, terminalH int
		gridW, gridH         int
		expectedCellWidths   []int
		expectedCellHeights  []int
	}{
		{
			name:                "width remainder distribution",
			terminalW:           82,
			terminalH:           10,
			gridW:               5,
			gridH:               1,
			expectedCellWidths:  []int{17, 17, 16, 16, 16},
			expectedCellHeights: []int{10},
		},
		{
			name:                "height remainder distribution",
			terminalW:           10,
			terminalH:           8,
			gridW:               1,
			gridH:               3,
			expectedCellWidths:  []int{10},
			expectedCellHeights: []int{3, 3, 2},
		},
		{
			name:                "both width and height remainder",
			terminalW:           13,
			terminalH:           7,
			gridW:               3,
			gridH:               2,
			expectedCellWidths:  []int{5, 4, 4}, // 13÷3=4 remainder 1
			expectedCellHeights: []int{4, 3},    // 7÷2=3 remainder 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGridLayout(tt.gridW, tt.gridH)

			// Add a dummy component
			mockComp := &MockComponent{}
			mockComp.On("Update", mock.AnythingOfType("tapioca.ResizeMsg")).Return(mockComp, nil)
			grid.Add(mockComp, 0, 0, 1, 1)

			// Call handleResize
			grid.handleResize(tt.terminalW, tt.terminalH)

			// Check cell dimensions
			assert.Equal(t, tt.expectedCellWidths, grid.gridMap.cellWidths, "cellWidths should match expected")
			assert.Equal(t, tt.expectedCellHeights, grid.gridMap.cellHeights, "cellHeights should match expected")

			mockComp.AssertExpectations(t)
		})
	}
}
