// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cup

import (
	"strings"
	"testing"

	"github.com/raohwork/huninn/pearl"
	"github.com/raohwork/huninn/tapioca"
	"github.com/stretchr/testify/assert"
)

func TestBorderedBox_View_ErrorScenarios(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "too small width",
			width:  3,
			height: 3,
		},
		{
			name:   "too small height",
			width:  8,
			height: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			innerMock := &MockRenderComponent{}
			innerMock.On("Init").Return(nil)
			innerMock.On("Update", tapioca.ResizeMsg{Width: tt.width - 2, Height: tt.height - 2}).Return(innerMock, nil)

			box := NewBorderedBox(innerMock)
			box.Init()
			box.Update(tapioca.ResizeMsg{Width: tt.width, Height: tt.height})

			result := box.View()
			assert.Equal(t, "too small", result)
		})
	}
}

func TestBorderedBox_View_NormalScenario(t *testing.T) {
	const (
		termWidth   = 8
		termHeight  = 3
		innerWidth  = 6 // 8 - 2 (left and right borders)
		innerHeight = 1 // 3 - 2 (top and bottom borders)
	)

	innerMock := &MockRenderComponent{}
	innerMock.On("Init").Return(nil)
	innerMock.On("Update", tapioca.ResizeMsg{Width: innerWidth, Height: innerHeight}).Return(innerMock, nil)
	innerMock.On("View").Return(createExactSizeString(strings.Repeat("X", innerWidth*innerHeight), innerWidth, innerHeight))

	box := NewBorderedBox(innerMock)
	box.Init()
	box.Update(tapioca.ResizeMsg{Width: termWidth, Height: termHeight})

	result := box.View()

	expectedRows := []string{
		"┌──────┐",
		"│XXXXXX│",
		"└──────┘",
	}

	actualRows := strings.Split(strings.TrimRight(result, "\n"), "\n")
	assert.Equal(t, len(expectedRows), len(actualRows), "should produce %d rows", len(expectedRows))
	assert.Equal(t, expectedRows, actualRows)
}

func TestBorderedBox_View_MissingOneBorder(t *testing.T) {
	const (
		termWidth   = 8
		termHeight  = 3
		innerWidth  = 6
		innerHeight = 1
	)

	tests := []struct {
		name     string
		config   BorderConfig
		expected []string
	}{
		{
			name: "missing left border",
			config: BorderConfig{
				Left:              false,
				Top:               true,
				Right:             true,
				Bottom:            true,
				VerticalLine:      '│',
				HorizontalLine:    '─',
				TopLeftCorner:     '┌',
				TopRightCorner:    '┐',
				BottomLeftCorner:  '└',
				BottomRightCorner: '┘',
			},
			expected: []string{
				"───────┐",
				"XXXXXXX│",
				"───────┘",
			},
		},
		{
			name: "missing top border",
			config: BorderConfig{
				Left:              true,
				Top:               false,
				Right:             true,
				Bottom:            true,
				VerticalLine:      '│',
				HorizontalLine:    '─',
				TopLeftCorner:     '┌',
				TopRightCorner:    '┐',
				BottomLeftCorner:  '└',
				BottomRightCorner: '┘',
			},
			expected: []string{
				"│XXXXXX│",
				"│XXXXXX│",
				"└──────┘",
			},
		},
		{
			name: "missing right border",
			config: BorderConfig{
				Left:              true,
				Top:               true,
				Right:             false,
				Bottom:            true,
				VerticalLine:      '│',
				HorizontalLine:    '─',
				TopLeftCorner:     '┌',
				TopRightCorner:    '┐',
				BottomLeftCorner:  '└',
				BottomRightCorner: '┘',
			},
			expected: []string{
				"┌───────",
				"│XXXXXXX",
				"└───────",
			},
		},
		{
			name: "missing bottom border",
			config: BorderConfig{
				Left:              true,
				Top:               true,
				Right:             true,
				Bottom:            false,
				VerticalLine:      '│',
				HorizontalLine:    '─',
				TopLeftCorner:     '┌',
				TopRightCorner:    '┐',
				BottomLeftCorner:  '└',
				BottomRightCorner: '┘',
			},
			expected: []string{
				"┌──────┐",
				"│XXXXXX│",
				"│XXXXXX│",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			innerMock := &MockRenderComponent{}
			innerMock.On("Init").Return(nil)

			// Calculate inner dimensions based on which borders are present
			innerW := termWidth
			innerH := termHeight
			if tt.config.Left {
				innerW--
			}
			if tt.config.Right {
				innerW--
			}
			if tt.config.Top {
				innerH--
			}
			if tt.config.Bottom {
				innerH--
			}

			innerMock.On("Update", tapioca.ResizeMsg{Width: innerW, Height: innerH}).Return(innerMock, nil)
			innerMock.On("View").Return(createExactSizeString(strings.Repeat("X", innerW*innerH), innerW, innerH))

			box := NewBorderedBox(innerMock)
			box.BorderConfig = tt.config
			box.Init()
			box.Update(tapioca.ResizeMsg{Width: termWidth, Height: termHeight})

			result := box.View()
			actualRows := strings.Split(strings.TrimRight(result, "\n"), "\n")
			assert.Equal(t, len(tt.expected), len(actualRows), "should produce %d rows", len(tt.expected))
			assert.Equal(t, tt.expected, actualRows)
		})
	}
}

func TestBorderedBox_View_OnlyOneBorder(t *testing.T) {
	const (
		termWidth   = 8
		termHeight  = 3
		innerWidth  = 7
		innerHeight = 2
	)

	tests := []struct {
		name     string
		config   BorderConfig
		expected []string
	}{
		{
			name: "only left border",
			config: BorderConfig{
				Left:              true,
				Top:               false,
				Right:             false,
				Bottom:            false,
				VerticalLine:      '│',
				HorizontalLine:    '─',
				TopLeftCorner:     '┌',
				TopRightCorner:    '┐',
				BottomLeftCorner:  '└',
				BottomRightCorner: '┘',
			},
			expected: []string{
				"│XXXXXXX",
				"│XXXXXXX",
				"│XXXXXXX",
			},
		},
		{
			name: "only top border",
			config: BorderConfig{
				Left:              false,
				Top:               true,
				Right:             false,
				Bottom:            false,
				VerticalLine:      '│',
				HorizontalLine:    '─',
				TopLeftCorner:     '┌',
				TopRightCorner:    '┐',
				BottomLeftCorner:  '└',
				BottomRightCorner: '┘',
			},
			expected: []string{
				"────────",
				"XXXXXXXX",
				"XXXXXXXX",
			},
		},
		{
			name: "only right border",
			config: BorderConfig{
				Left:              false,
				Top:               false,
				Right:             true,
				Bottom:            false,
				VerticalLine:      '│',
				HorizontalLine:    '─',
				TopLeftCorner:     '┌',
				TopRightCorner:    '┐',
				BottomLeftCorner:  '└',
				BottomRightCorner: '┘',
			},
			expected: []string{
				"XXXXXXX│",
				"XXXXXXX│",
				"XXXXXXX│",
			},
		},
		{
			name: "only bottom border",
			config: BorderConfig{
				Left:              false,
				Top:               false,
				Right:             false,
				Bottom:            true,
				VerticalLine:      '│',
				HorizontalLine:    '─',
				TopLeftCorner:     '┌',
				TopRightCorner:    '┐',
				BottomLeftCorner:  '└',
				BottomRightCorner: '┘',
			},
			expected: []string{
				"XXXXXXXX",
				"XXXXXXXX",
				"────────",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			innerMock := &MockRenderComponent{}
			innerMock.On("Init").Return(nil)

			// Calculate inner dimensions based on which borders are present
			innerW := termWidth
			innerH := termHeight
			if tt.config.Left {
				innerW--
			}
			if tt.config.Right {
				innerW--
			}
			if tt.config.Top {
				innerH--
			}
			if tt.config.Bottom {
				innerH--
			}

			innerMock.On("Update", tapioca.ResizeMsg{Width: innerW, Height: innerH}).Return(innerMock, nil)
			innerMock.On("View").Return(createExactSizeString(strings.Repeat("X", innerW*innerH), innerW, innerH))

			box := NewBorderedBox(innerMock)
			box.BorderConfig = tt.config
			box.Init()
			box.Update(tapioca.ResizeMsg{Width: termWidth, Height: termHeight})

			result := box.View()
			actualRows := strings.Split(strings.TrimRight(result, "\n"), "\n")
			assert.Equal(t, len(tt.expected), len(actualRows), "should produce %d rows", len(tt.expected))
			assert.Equal(t, tt.expected, actualRows)
		})
	}
}

func TestBorderedBox_RenderCaption(t *testing.T) {
	cases := []struct {
		name              string
		caption           string
		expectedString    string
		expectedRestWidth int
	}{
		{
			name:              "no caption",
			caption:           "",
			expectedString:    "",
			expectedRestWidth: 8,
		},
		{
			name:              "caption fits within width",
			caption:           "abc",
			expectedString:    "─ abc ─",
			expectedRestWidth: 1,
		},
		{
			name:              "caption exactly fits width",
			caption:           "abcd",
			expectedString:    "─ abcd ─",
			expectedRestWidth: 0,
		},
		{
			name:              "caption exceeds width",
			caption:           "abcdefgh",
			expectedString:    "─ ab… ─",
			expectedRestWidth: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inner := pearl.NewBlock()
			box := NewBorderedBox(inner)
			box.BorderConfig = BorderConfig{
				Top:            true,
				Left:           true,
				Right:          true,
				VerticalLine:   '│',
				HorizontalLine: '─',
				TopLeftCorner:  '┌',
				TopRightCorner: '┐',
			}
			box.Init()
			box.Update(tapioca.ResizeMsg{Width: 10, Height: 3})
			box.SetCaption(tc.caption)
			buf := &strings.Builder{}
			restWidth := box.renderCaption(buf)
			result := buf.String()
			assert.Equal(t, tc.expectedString, result)
			assert.Equal(t, tc.expectedRestWidth, restWidth)
		})
	}
}

func TestBorderedBox_CaptionCases(t *testing.T) {
	f := func(expected []string) string {
		return strings.Join(expected, "\n")
	}

	t.Run("no top border", func(t *testing.T) {
		inner := pearl.NewBlock()
		box := NewBorderedBox(inner)
		box.BorderConfig.Top = false
		box.Init()
		box.Update(tapioca.ResizeMsg{Width: 8, Height: 3})
		box.SetCaption("a")
		expected := f([]string{
			"│      │",
			"│      │",
			"└──────┘",
		})

		result := box.View()
		assert.Equal(t, expected, result)
	})

	t.Run("with top border", func(t *testing.T) {
		inner := pearl.NewBlock()
		box := NewBorderedBox(inner)
		box.Init()
		box.Update(tapioca.ResizeMsg{Width: 8, Height: 3})
		box.SetCaption("a")
		expected := f([]string{
			"┌─ a ──┐",
			"│      │",
			"└──────┘",
		})

		result := box.View()
		assert.Equal(t, expected, result)
	})
}
