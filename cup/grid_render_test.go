// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cup

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRenderComponent is a mock component that outputs exact size strings
type MockRenderComponent struct {
	mock.Mock
	width, height int
	content       string
}

func (m *MockRenderComponent) Init() tea.Cmd {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(tea.Cmd)
}

func (m *MockRenderComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	args := m.Called(msg)
	// Store resize dimensions for View() method
	if resizeMsg, ok := msg.(tapioca.ResizeMsg); ok {
		m.width = resizeMsg.Width
		m.height = resizeMsg.Height
	}
	cmd := args.Get(1)
	if cmd == nil {
		return args.Get(0).(tea.Model), nil
	}
	return args.Get(0).(tea.Model), cmd.(tea.Cmd)
}

func (m *MockRenderComponent) View() string {
	args := m.Called()
	return args.String(0)
}

// Helper function to create exact size string with specific content
func createExactSizeString(content string, width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	lines := make([]string, height)
	contentRunes := []rune(content)
	contentIdx := 0

	for i := range height {
		line := make([]rune, width)
		// Initialize with spaces
		for j := range width {
			line[j] = ' '
		}

		// Fill with content
		j := 0
		for j < width && contentIdx < len(contentRunes) {
			if contentRunes[contentIdx] == '\n' {
				contentIdx++
				break // Move to next line
			}
			line[j] = contentRunes[contentIdx]
			j++
			contentIdx++
		}
		lines[i] = string(line)
	}

	return strings.Join(lines, "\n")
}

func TestGridLayout_View_NormalCase(t *testing.T) {
	// Test case: 32x25 terminal with 3x3 grid
	// Components layout matches the comment:
	// A A B
	// C D B
	// C E E

	grid := NewGridLayout(3, 3)

	// Calculate expected cell sizes for 32x25 terminal
	// Width: 32÷3 = 10 remainder 2, so cells are [11, 11, 10]
	// Height: 25÷3 = 8 remainder 1, so cells are [9, 8, 8]

	// Component A: (0, 0, 2, 1) -> width: 11+11=22, height: 9
	mockA := &MockRenderComponent{}
	mockA.On("Update", tapioca.ResizeMsg{Width: 22, Height: 9}).Return(mockA, nil)
	mockA.On("View").Return(createExactSizeString(strings.Repeat("A", 22*9), 22, 9))

	// Component B: (2, 0, 1, 2) -> width: 10, height: 9+8=17
	mockB := &MockRenderComponent{}
	mockB.On("Update", tapioca.ResizeMsg{Width: 10, Height: 17}).Return(mockB, nil)
	mockB.On("View").Return(createExactSizeString(strings.Repeat("B", 10*17), 10, 17))

	// Component C: (0, 1, 1, 2) -> width: 11, height: 8+8=16
	mockC := &MockRenderComponent{}
	mockC.On("Update", tapioca.ResizeMsg{Width: 11, Height: 16}).Return(mockC, nil)
	mockC.On("View").Return(createExactSizeString(strings.Repeat("C", 11*16), 11, 16))

	// Component D: (1, 1, 1, 1) -> width: 11, height: 8
	mockD := &MockRenderComponent{}
	mockD.On("Update", tapioca.ResizeMsg{Width: 11, Height: 8}).Return(mockD, nil)
	mockD.On("View").Return(createExactSizeString(strings.Repeat("D", 11*8), 11, 8))

	// Component E: (1, 2, 2, 1) -> width: 11+10=21, height: 8
	mockE := &MockRenderComponent{}
	mockE.On("Update", tapioca.ResizeMsg{Width: 21, Height: 8}).Return(mockE, nil)
	mockE.On("View").Return(createExactSizeString(strings.Repeat("E", 21*8), 21, 8))

	// Add components to grid
	success := grid.Add(mockA, 0, 0, 2, 1)
	assert.True(t, success)
	success = grid.Add(mockB, 2, 0, 1, 2)
	assert.True(t, success)
	success = grid.Add(mockC, 0, 1, 1, 2)
	assert.True(t, success)
	success = grid.Add(mockD, 1, 1, 1, 1)
	assert.True(t, success)
	success = grid.Add(mockE, 1, 2, 2, 1)
	assert.True(t, success)

	// Trigger resize to setup component sizes
	grid.handleResize(32, 25)

	// Test that View() returns the correct rendered output
	// Note: Since current implementation returns "", this test will verify the expected behavior
	result := strings.TrimRight(grid.View(), "\n")

	expectedOutput := []string{
		//         L          C         R
		"AAAAAAAAAAAAAAAAAAAAAABBBBBBBBBB", // Row 0 (Top)
		"AAAAAAAAAAAAAAAAAAAAAABBBBBBBBBB", // Row 1
		"AAAAAAAAAAAAAAAAAAAAAABBBBBBBBBB", // Row 2
		"AAAAAAAAAAAAAAAAAAAAAABBBBBBBBBB", // Row 3
		"AAAAAAAAAAAAAAAAAAAAAABBBBBBBBBB", // Row 4
		"AAAAAAAAAAAAAAAAAAAAAABBBBBBBBBB", // Row 5
		"AAAAAAAAAAAAAAAAAAAAAABBBBBBBBBB", // Row 6
		"AAAAAAAAAAAAAAAAAAAAAABBBBBBBBBB", // Row 7
		"AAAAAAAAAAAAAAAAAAAAAABBBBBBBBBB", // Row 8
		"CCCCCCCCCCCDDDDDDDDDDDBBBBBBBBBB", // Row 0 (Center)
		"CCCCCCCCCCCDDDDDDDDDDDBBBBBBBBBB", // Row 1
		"CCCCCCCCCCCDDDDDDDDDDDBBBBBBBBBB", // Row 2
		"CCCCCCCCCCCDDDDDDDDDDDBBBBBBBBBB", // Row 3
		"CCCCCCCCCCCDDDDDDDDDDDBBBBBBBBBB", // Row 4
		"CCCCCCCCCCCDDDDDDDDDDDBBBBBBBBBB", // Row 5
		"CCCCCCCCCCCDDDDDDDDDDDBBBBBBBBBB", // Row 6
		"CCCCCCCCCCCDDDDDDDDDDDBBBBBBBBBB", // Row 7
		"CCCCCCCCCCCEEEEEEEEEEEEEEEEEEEEE", // Row 0 (Bottom)
		"CCCCCCCCCCCEEEEEEEEEEEEEEEEEEEEE", // Row 1
		"CCCCCCCCCCCEEEEEEEEEEEEEEEEEEEEE", // Row 2
		"CCCCCCCCCCCEEEEEEEEEEEEEEEEEEEEE", // Row 3
		"CCCCCCCCCCCEEEEEEEEEEEEEEEEEEEEE", // Row 4
		"CCCCCCCCCCCEEEEEEEEEEEEEEEEEEEEE", // Row 5
		"CCCCCCCCCCCEEEEEEEEEEEEEEEEEEEEE", // Row 6
		"CCCCCCCCCCCEEEEEEEEEEEEEEEEEEEEE", // Row 7
	}

	assert.Equal(t, strings.Join(expectedOutput, "\n"), result, "View() should return the expected grid layout")

	// Verify all mock expectations were met
	mockA.AssertExpectations(t)
	mockB.AssertExpectations(t)
	mockC.AssertExpectations(t)
	mockD.AssertExpectations(t)
	mockE.AssertExpectations(t)
}

func TestGridLayout_View_PartialCoverage(t *testing.T) {
	// Test case: 2x2 grid with only one component, rest should be filled with spaces
	// Terminal size 20x10, grid should have cells of 10x5 each
	// Component layout:
	// A .
	// . .
	// Where A is the component and . represents empty cells filled with spaces

	grid := NewGridLayout(2, 2)

	// Calculate expected cell sizes for 20x10 terminal
	// Width: 20÷2 = 10, Height: 10÷2 = 5
	// Each cell is 10x5

	// Component A: (0, 0, 1, 1) -> width: 10, height: 5
	mockA := &MockRenderComponent{}
	mockA.On("Update", tapioca.ResizeMsg{Width: 10, Height: 5}).Return(mockA, nil)
	mockA.On("View").Return(createExactSizeString(strings.Repeat("A", 10*5), 10, 5))

	// Add only one component to grid, leaving other cells empty
	success := grid.Add(mockA, 0, 0, 1, 1)
	assert.True(t, success)

	// Trigger resize to setup component sizes
	grid.handleResize(20, 10)

	// Test that View() returns the correct rendered output with spaces for empty cells
	result := strings.TrimRight(grid.View(), "\n")

	expectedOutput := []string{
		//        L         R
		"AAAAAAAAAA          ", // Row 0: A's content | empty cell
		"AAAAAAAAAA          ", // Row 1: A's content | empty cell
		"AAAAAAAAAA          ", // Row 2: A's content | empty cell
		"AAAAAAAAAA          ", // Row 3: A's content | empty cell
		"AAAAAAAAAA          ", // Row 4: A's content | empty cell
		"                    ", // Row 5: empty cell  | empty cell
		"                    ", // Row 6: empty cell  | empty cell
		"                    ", // Row 7: empty cell  | empty cell
		"                    ", // Row 8: empty cell  | empty cell
		"                    ", // Row 9: empty cell  | empty cell
	}

	// Verify that View() returns the correct grid layout with spaces for empty cells
	assert.Equal(t, strings.Join(expectedOutput, "\n"), result, "View() should return grid layout with spaces for empty cells")

	// Verify mock expectations were met
	mockA.AssertExpectations(t)
}

func TestGridLayout_View_ErrorCase(t *testing.T) {
	tests := []struct {
		name        string
		setupError  func(*GridLayout)
		expectError string
	}{
		{
			name: "terminal size too small - width",
			setupError: func(g *GridLayout) {
				// Simulate small terminal size that triggers hasError
				g.handleResize(2, 10) // width <= 2
			},
			expectError: "Terminal size too small",
		},
		{
			name: "terminal size too small - height",
			setupError: func(g *GridLayout) {
				// Simulate small terminal size that triggers hasError
				g.handleResize(10, 0) // height < 1
			},
			expectError: "Terminal size too small",
		},
		{
			name: "component size too small",
			setupError: func(g *GridLayout) {
				// Add a component and trigger resize that makes component too small
				mockComp := &MockRenderComponent{}
				g.Add(mockComp, 0, 0, 1, 1)
				// Create grid with many columns so each cell becomes too narrow
				g.w = 20

				g.h = 1
				g.gridMap = newGridMap(20, 1)
				g.handleResize(10, 5) // 10÷20 < 2, component width will be < 2
			},
			expectError: "Terminal size too small",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGridLayout(2, 2)

			// Setup error condition
			tt.setupError(grid)

			// Test that View() returns the error message
			result := grid.View()
			assert.Equal(t, tt.expectError, result)
			assert.True(t, grid.hasError, "hasError should be true")
		})
	}
}

func TestGridLayout_View_EmptyGrid(t *testing.T) {
	// Test grid with no components
	grid := NewGridLayout(3, 3)

	// Trigger resize with valid size
	grid.handleResize(30, 15)

	// Should not have error
	assert.False(t, grid.hasError)

	// View should return empty string (current implementation)
	result := grid.View()
	assert.Equal(t, "", result)
}

func TestGridLayout_View_ComponentExactSizing(t *testing.T) {
	// Test that components must return exact size strings
	// This test verifies the mock component helper function works correctly

	tests := []struct {
		name     string
		content  string
		width    int
		height   int
		expected string
	}{
		{
			name:     "simple case",
			content:  "hello",
			width:    10,
			height:   2,
			expected: "hello     \n          ",
		},
		{
			name:     "with newline",
			content:  "line1\nline2",
			width:    8,
			height:   3,
			expected: "line1   \nline2   \n        ",
		},
		{
			name:     "overflow content",
			content:  "very long content that exceeds width",
			width:    5,
			height:   2,
			expected: "very \nlong ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createExactSizeString(tt.content, tt.width, tt.height)
			assert.Equal(t, tt.expected, result)

			// Verify dimensions
			lines := strings.Split(result, "\n")
			assert.Len(t, lines, tt.height, "Should have correct number of lines")
			for i, line := range lines {
				assert.Len(t, line, tt.width, "Line %d should have correct width", i)
			}
		})
	}
}
