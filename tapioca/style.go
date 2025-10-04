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

func (s *style) Clone() *style {
	if s.isEmpty() {
		return nil
	}

	ret := *s
	return &ret
}

func (s *style) isEmpty() bool {
	if s == nil {
		return true
	}
	return s.fg == "" && s.bg == "" && !s.bold && !s.faint && !s.italic && !s.underline && !s.strike && !s.blink && !s.reverse && !s.hidden
}

// Render returns the ANSI escape sequence to transition from prevStyle to s.
//
// It uses a really simple strategy: if prevStyle is not empty, it resets
// all styles first, then applies s. This is not the most efficient way,
// but it's simple and works well enough in practice.
func (s *style) Render(prevStyle *style) string {
	if !prevStyle.isEmpty() {
		return "\x1b[0m" + s.String()
	}

	return s.String()
}

// String renders the style as an ANSI escape sequence.
func (s *style) String() string {
	if s.isEmpty() {
		return ""
	}

	b := &strings.Builder{}
	w := func(param string) {
		b.WriteString("\x1b[")
		b.WriteString(param)
		b.WriteString("m")
	}

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
	if newStyle == nil {
		newStyle = &style{}
	}

	// Strip off "\x1b[" and "m"
	paramsStr := strings.TrimSuffix(strings.TrimPrefix(code, "\x1b["), "m")
	if paramsStr == "" || paramsStr == "0" {
		// e.g., \x1b[m is a reset
		return nil
	}

	parts := strings.Split(paramsStr, ";")
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		switch part {
		case "0": // Reset
			if !newStyle.isEmpty() {
				newStyle = &style{}
			}
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
	if newStyle.isEmpty() {
		return nil
	}
	return newStyle
}
