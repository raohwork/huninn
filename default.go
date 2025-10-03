// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package huninn

import (
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/cup"
	"github.com/raohwork/huninn/pearl"
)

type defaultModel struct {
	tea.Model
}

func (m defaultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		m.Model, cmd = m.Model.Update(msg)
		return m, cmd
	}

	if s := k.String(); s == "ctrl+c" || s == "q" {
		return m, tea.Quit
	}

	m.Model, cmd = m.Model.Update(msg)
	return m, cmd
}

// Default returns default flavor of huninn UI. This will be dropped in next minor version.
//
// The default UI consists of three parts (from top to bottom):
//   - Task list with given height (in number of lines). tlSize < 3 will split the screen into two equal parts.
//   - Log panel
//   - Status bar: single line showing status.
//
// Typical usage is like
//
//	m, f := huninn.Default(10) // task list with 10 lines height
//	app := tea.NewProgram(m, opts...)
//	setStatus, taskManager, writer := f(app.Send)
//	logger := log.New(writer, "", log.LstdFlags)
//	go func() {
//	    setStatus("Starting...")
//	    task := taskManager.AddTask("Example Task")
//	    for i := 0; i <= 100; i += 10 {
//	        logger.Println("Progress", i)
//	        task.SetState(pearl.TaskRunning, float64(i)/100)
//	        time.Sleep(time.Second)
//	    }
//	    task.Done()
//	    setStatus("All tasks done.")
//	}()
//	app.Run()
func Default(tlSize int) (m tea.Model, f func(func(tea.Msg)) (setStatus func(...string), taskManager pearl.TaskManager, w io.Writer)) {
	status := pearl.NewBlock()
	tl := pearl.NewTaskList()
	logger := pearl.NewLogPanel(1000)

	var main tea.Model
	if tlSize < 3 {
		x := cup.NewGridLayout(1, 2)
		x.Add(tl, 0, 0, 1, 1)
		x.Add(logger, 1, 0, 1, 1)
		main = x
	} else {
		main = cup.FixedTopLayout(tlSize, tl, logger)
	}
	root := cup.FixedBottomLayout(1, main, status)

	return defaultModel{root}, func(send func(tea.Msg)) (setStatus func(...string), taskManager pearl.TaskManager, w io.Writer) {
		return status.Setter(send),
			tl.CreateManager(send),
			logger.CreateWriter(send, nil)
	}
}
