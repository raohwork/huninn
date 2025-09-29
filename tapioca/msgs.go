// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tapioca

// ResizeMsg denotes the layout component is asking target to resize
type ResizeMsg struct {
	Width  int
	Height int
}

// ScrollLeftMsg tells the component to scroll left (towards start of line)
type ScrollLeftMsg uint

// ScrollRightMsg tells the component to scroll right (towards end of line)
type ScrollRightMsg uint

// ScrollUpMsg tells the component to scroll up (toward top of component)
type ScrollUpMsg uint

// ScrollDownMsg tells the component to scroll down (towards bottom of component)
type ScrollDownMsg uint

// ScrollTopMsg tells the component to scroll to the top of the component
type ScrollTopMsg struct{}

// ScrollBottomMsg tells the component to scroll to the bottom of the component
type ScrollBottomMsg struct{}

// ScrollBeginMsg tells the component to scroll to the beginning of the line
type ScrollBeginMsg struct{}

// ScrollEndMsg tells the component to scroll to the end of the line
type ScrollEndMsg struct{}
