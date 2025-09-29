// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

import (
	"regexp"
	"strings"
)

type style struct {
	fg, bg    string
	bold      bool
	faint     bool
	italic    bool
	underline bool
	strike    bool
	blink     bool
	reverse   bool
	hidden    bool
}

func (s style) Clone() *style {
	return &style{
		fg:        s.fg,
		bg:        s.bg,
		bold:      s.bold,
		faint:     s.faint,
		italic:    s.italic,
		underline: s.underline,
		strike:    s.strike,
		blink:     s.blink,
		reverse:   s.reverse,
		hidden:    s.hidden,
	}
}

// needsReset checks if transitioning from prevStyle to this style requires a reset
func (s *style) needsReset(prevStyle *style) bool {
	if prevStyle == nil {
		return false
	}
	// Check if any boolean attribute goes from true to false
	return (prevStyle.bold && !s.bold) ||
		(prevStyle.faint && !s.faint) ||
		(prevStyle.italic && !s.italic) ||
		(prevStyle.underline && !s.underline) ||
		(prevStyle.blink && !s.blink) ||
		(prevStyle.reverse && !s.reverse) ||
		(prevStyle.hidden && !s.hidden) ||
		(prevStyle.strike && !s.strike)
}

func (s *style) String() string {
	if s == nil {
		return "\x1b[m"
	}

	b := &strings.Builder{}
	w := func(param string) {
		b.WriteString("\x1b[")
		b.WriteString(param)
		b.WriteString("m")
	}
	w("0")

	str := func(param string) {
		if len(param) > 0 {
			w(param)
		}
	}
	bool := func(enabled bool, param string) {
		if enabled {
			w(param)
		}
	}

	str(s.fg)
	str(s.bg)
	bool(s.bold, "1")
	bool(s.faint, "2")
	bool(s.italic, "3")
	bool(s.underline, "4")
	bool(s.blink, "5")
	bool(s.reverse, "7")
	bool(s.hidden, "8")
	bool(s.strike, "9")

	return b.String()
}

var (
	ansiStyleRegex = regexp.MustCompile(`\x1b\[([0-9]{1,3}(;[0-9]{1,3})*)?m`)
	ansiOtherRegex = regexp.MustCompile(`\x1b\[([0-9]+(;[0-9]+)*)?[ABCDEFGHJKSTfsuhl]`)
)

// parseAnsiCode takes a full escape sequence (e.g., "\x1b[32m") and the previous style.
// It returns a new style object.
func parseAnsiCode(code string, previousStyle *style) *style {
	// Crucial: Clone the previous style to avoid mutation.
	newStyle := previousStyle.Clone()

	// Strip off "\x1b[" and "m"
	paramsStr := strings.TrimSuffix(strings.TrimPrefix(code, "\x1b["), "m")
	if paramsStr == "" {
		// e.g., \x1b[m is a reset
		return &style{}
	}

	parts := strings.Split(paramsStr, ";")
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		switch part {
		case "0": // Reset
			return &style{}
		case "1":
			newStyle.bold = true
		case "2":
			newStyle.faint = true
		case "3":
			newStyle.italic = true
		case "4":
			newStyle.underline = true
		case "5":
			newStyle.blink = true
		case "7":
			newStyle.reverse = true
		case "8":
			newStyle.hidden = true
		case "9":
			newStyle.strike = true
		case "22": // Normal intensity (turn off bold and faint)
			newStyle.bold = false
			newStyle.faint = false
		case "23": // No italic
			newStyle.italic = false
		case "24": // No underline
			newStyle.underline = false
		case "25": // No blink
			newStyle.blink = false
		case "27": // No reverse
			newStyle.reverse = false
		case "28": // No hidden
			newStyle.hidden = false
		case "29": // No strikethrough
			newStyle.strike = false
		case "39": // Default foreground
			newStyle.fg = ""
		case "49": // Default background
			newStyle.bg = ""
		case "30", "31", "32", "33", "34", "35", "36", "37": // Standard foreground colors
			newStyle.fg = part
		case "40", "41", "42", "43", "44", "45", "46", "47": // Standard background colors
			newStyle.bg = part
		case "90", "91", "92", "93", "94", "95", "96", "97": // Bright foreground colors
			newStyle.fg = part
		case "100", "101", "102", "103", "104", "105", "106", "107": // Bright background colors
			newStyle.bg = part
		case "38": // 256-color or RGB foreground
			if i+2 < len(parts) && parts[i+1] == "5" {
				// 256-color: 38;5;n
				newStyle.fg = "38;5;" + parts[i+2]
				i += 2 // Skip the next two parts
			} else if i+4 < len(parts) && parts[i+1] == "2" {
				// RGB: 38;2;r;g;b
				newStyle.fg = "38;2;" + parts[i+2] + ";" + parts[i+3] + ";" + parts[i+4]
				i += 4 // Skip the next four parts
			}
		case "48": // 256-color or RGB background
			if i+2 < len(parts) && parts[i+1] == "5" {
				// 256-color: 48;5;n
				newStyle.bg = "48;5;" + parts[i+2]
				i += 2 // Skip the next two parts
			} else if i+4 < len(parts) && parts[i+1] == "2" {
				// RGB: 48;2;r;g;b
				newStyle.bg = "48;2;" + parts[i+2] + ";" + parts[i+3] + ";" + parts[i+4]
				i += 4 // Skip the next four parts
			}
			// Ignore unknown codes (error handling strategy: continue with next codes)
		}
	}
	return newStyle
}
