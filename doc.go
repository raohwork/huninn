// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package huninn provides tools to build fency CLI applications. It
// uses the Bubble Tea framework under the hood.
//
// # WIP
//
// Huninn is still a work in progress. The API may change without
// notice.
//
// # Overview
//
// It focuses on information displaying, not user interaction. A
// Huninn component is a tea.Model with extra restrictions:
//
//   - It accepts tapioca.ResizeMsg.
//   - View() method output a string exactly fitting its size.
//
// For example, a Huninn component with size 3x2 (specified by
// tapioca.ResizeMsg) must output exactly 3 columns and 2 rows of
// characters, like "abc\nd  ". Tailing newlines are optional.
//
// There are some predefined components in package pearl, and some
// tools to help you build your own components in package tapioca.
//
// In this package, we provide some preset for application creators.
// These presets are separated by "Ice" and "Suguar". The more ice, the
// "cooler" (fencier) the look is. The more sugar, the "sweeter" the
// application can use (feature-rich).
package huninn
