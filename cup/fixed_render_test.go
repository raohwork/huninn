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
)

func splitFixedRows(output string, width, height int, horizontal bool) []string {
	if horizontal {
		rows := make([]string, 0, height)
		for start := 0; start < len(output); start += width {
			end := min(start+width, len(output))
			rows = append(rows, output[start:end])
		}
		return rows
	}

	rows := strings.Split(output, "\n")
	if len(rows) > height && rows[len(rows)-1] == "" {
		rows = rows[:len(rows)-1]
	}
	return rows
}

func TestFixedLayout_View_NormalScenarios(t *testing.T) {
	const (
		termWidth  = 8
		termHeight = 8
		reserve    = 3
	)

	tests := []struct {
		name         string
		makeLayout   func(reserveComponent, restComponent tea.Model) *FixedLayout
		horizontal   bool
		expectedRows []string
		reserveChar  string
		restChar     string
	}{
		{
			name: "left",
			makeLayout: func(reserveComponent, restComponent tea.Model) *FixedLayout {
				return FixedLeftLayout(reserve, reserveComponent, restComponent)
			},
			horizontal: true,
			expectedRows: []string{
				"LLLRRRRR",
				"LLLRRRRR",
				"LLLRRRRR",
				"LLLRRRRR",
				"LLLRRRRR",
				"LLLRRRRR",
				"LLLRRRRR",
				"LLLRRRRR",
			},
			reserveChar: "L",
			restChar:    "R",
		},
		{
			name: "right",
			makeLayout: func(reserveComponent, restComponent tea.Model) *FixedLayout {
				return FixedRightLayout(reserve, restComponent, reserveComponent)
			},
			horizontal: true,
			expectedRows: []string{
				"LLLLLRRR",
				"LLLLLRRR",
				"LLLLLRRR",
				"LLLLLRRR",
				"LLLLLRRR",
				"LLLLLRRR",
				"LLLLLRRR",
				"LLLLLRRR",
			},
			reserveChar: "R",
			restChar:    "L",
		},
		{
			name: "top",
			makeLayout: func(reserveComponent, restComponent tea.Model) *FixedLayout {
				return FixedTopLayout(reserve, reserveComponent, restComponent)
			},
			horizontal: false,
			expectedRows: []string{
				"TTTTTTTT",
				"TTTTTTTT",
				"TTTTTTTT",
				"BBBBBBBB",
				"BBBBBBBB",
				"BBBBBBBB",
				"BBBBBBBB",
				"BBBBBBBB",
			},
			reserveChar: "T",
			restChar:    "B",
		},
		{
			name: "bottom",
			makeLayout: func(reserveComponent, restComponent tea.Model) *FixedLayout {
				return FixedBottomLayout(reserve, restComponent, reserveComponent)
			},
			horizontal: false,
			expectedRows: []string{
				"TTTTTTTT",
				"TTTTTTTT",
				"TTTTTTTT",
				"TTTTTTTT",
				"TTTTTTTT",
				"BBBBBBBB",
				"BBBBBBBB",
				"BBBBBBBB",
			},
			reserveChar: "B",
			restChar:    "T",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reserveWidth := reserve
			restWidth := termWidth - reserve
			reserveHeight := termHeight
			restHeight := termHeight
			if !tt.horizontal {
				reserveWidth = termWidth
				restWidth = termWidth
				reserveHeight = reserve
				restHeight = termHeight - reserve
			}

			reserveMock := &MockRenderComponent{}
			reserveMock.On("Update", tapioca.ResizeMsg{Width: reserveWidth, Height: reserveHeight}).Return(reserveMock, nil)
			reserveMock.On("View").Return(createExactSizeString(strings.Repeat(tt.reserveChar, reserveWidth*reserveHeight), reserveWidth, reserveHeight))

			restMock := &MockRenderComponent{}
			restMock.On("Update", tapioca.ResizeMsg{Width: restWidth, Height: restHeight}).Return(restMock, nil)
			restMock.On("View").Return(createExactSizeString(strings.Repeat(tt.restChar, restWidth*restHeight), restWidth, restHeight))

			layout := tt.makeLayout(reserveMock, restMock)
			layout.handleResize(termWidth, termHeight)

			rows := splitFixedRows(layout.View(), termWidth, termHeight, tt.horizontal)
			assert.Equal(t, termHeight, len(rows), "should produce %d rows", termHeight)
			assert.Equal(t, tt.expectedRows, rows)
		})
	}
}

func TestFixedLayout_View_TerminalTooSmall(t *testing.T) {
	const (
		termWidth  = 5
		termHeight = 5
		reserve    = 5
	)

	tests := []struct {
		name       string
		makeLayout func(reserveComponent, restComponent tea.Model) *FixedLayout
		horizontal bool
	}{
		{
			name: "left",
			makeLayout: func(reserveComponent, restComponent tea.Model) *FixedLayout {
				return FixedLeftLayout(reserve, reserveComponent, restComponent)
			},
			horizontal: true,
		},
		{
			name: "right",
			makeLayout: func(reserveComponent, restComponent tea.Model) *FixedLayout {
				return FixedRightLayout(reserve, restComponent, reserveComponent)
			},
			horizontal: true,
		},
		{
			name: "top",
			makeLayout: func(reserveComponent, restComponent tea.Model) *FixedLayout {
				return FixedTopLayout(reserve, reserveComponent, restComponent)
			},
			horizontal: false,
		},
		{
			name: "bottom",
			makeLayout: func(reserveComponent, restComponent tea.Model) *FixedLayout {
				return FixedBottomLayout(reserve, restComponent, reserveComponent)
			},
			horizontal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reserveWidth := reserve
			reserveHeight := termHeight
			if !tt.horizontal {
				reserveWidth = termWidth
				reserveHeight = reserve
			}

			reserveMock := &MockRenderComponent{}
			reserveMock.On("Update", tapioca.ResizeMsg{Width: reserveWidth, Height: reserveHeight}).Return(reserveMock, nil)

			restMock := &MockRenderComponent{}

			layout := tt.makeLayout(reserveMock, restMock)
			layout.handleResize(termWidth, termHeight)

			assert.Equal(t, "Terminal too small", layout.View())
		})
	}
}
