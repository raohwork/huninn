// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package huninn

import (
	"context"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/cup"
	"github.com/raohwork/huninn/pearl"
	"github.com/raohwork/huninn/tapioca"
	"github.com/raohwork/task"
)

// NSLIComponent returns bordered flavor of huninn UI.
//
// The NSLI (No Sugar, Light Ice) UI consists of three parts with borders:
//   - A status bar at the bottom, showing a single line of text.
//   - A task list at the top, showing multiple tasks with progress in percentage.
//   - A log panel in the middle, showing log messages.
//
// It balances screen space and visual appeal.
//
// The tlSize parameter specifies the height of the task list in rows. The
// logBufferSize parameter specifies the maximum number of log entries to keep.
// If tlSize is less than 3, the task list and log panel will have equal height.
// The logBufferSize is automatically adjusted to be at least 10.
//
// The component also handles two shortcuts to terminate the program:
//   - Ctrl+C
//   - q
func NSLIComponent(tlSize int, logBufferSize int) (
	tea.Model,
	func(send func(tea.Msg)) (
		setStatus func(string),
		taskManager pearl.TaskManager,
		w io.Writer,
		logScroller tapioca.ScrollController,
	),
) {
	status := pearl.NewSpan()
	tasks := pearl.NewTaskList()
	lp := pearl.NewLogPanel(logBufferSize)

	lpBox := cup.NewBorderedBoxWithCaption(lp, "Logs")
	lpBox.Left = false
	lpBox.Right = false

	mainBox := cup.FixedTopLayout(tlSize, tasks, lpBox)
	allBox := cup.FixedBottomLayout(1, mainBox, status)
	root := cup.NewBorderedBoxWithCaption(allBox, "Tasks")

	return noSuguarComponent{root}, func(send func(tea.Msg)) (func(string), pearl.TaskManager, io.Writer, tapioca.ScrollController) {
		setStatus := status.Setter(send)
		tm := tasks.CreateManager(send)
		w := lp.CreateWriter(send, nil)
		return setStatus, tm, w, lp.ScrollController()
	}
}

func nsli(tlSize, logBufferSize int, opts ...tea.ProgramOption) (
	prog *tea.Program,
	setStatus func(string),
	tm pearl.TaskManager,
	w io.Writer,
	s tapioca.ScrollController,
) {
	m, f := NSLIComponent(tlSize, logBufferSize)
	prog = tea.NewProgram(m, opts...)

	setStatus, tm, w, s = f(prog.Send)
	return
}

// NSLI (No Sugar, Light Ice) wraps NSLIComponent to provide a ready-to-use
// program.
//
// Cancelling the context will terminate the program, making it suitable for
// cooperating with signal.NotifyContext.
//
// The tlSize and LogBufferSize parameters are passed to [NSLIComponent].
//
// The opts parameters are passed to [tea.NewProgram].
func NSLI(tlSize, logBufferSize int, opts ...tea.ProgramOption) (
	prog func(context.Context) error,
	setStatus func(string),
	tm pearl.TaskManager,
	w io.Writer,
	s tapioca.ScrollController,
) {
	app, setStatus, tm, w, s := nsli(tlSize, logBufferSize, opts...)
	prog = progAsTask(app)
	return
}

// LSLI (Light Sugar, Light Ice) adds some features to NSLI.
//
// It accepts a JobFactory to create a job function, which will be run
// concurrently with the UI program. The job function is the main program
// of your application, and it should utilize the provided setStatus,
// task manager, and writer to interact with the UI.
//
// The tlSize and logBufferSize parameters are passed to [NSLIComponent].
//
// The opts parameters are passed to [tea.NewProgram].
//
// If you set wait to true, the UI will remain active after the job
// completes successfully, allowing the user to review the final status
// and logs. UI always remain active if the job ends with an error.
func LSLI(tlSize, logBufferSize int, factory JobFactory, wait bool, opts ...tea.ProgramOption) func(context.Context) error {
	app, setStatus, tm, w, s := nsli(tlSize, logBufferSize, opts...)
	job := factory(setStatus, tm, w, s, app.Quit)

	return func(ctx context.Context) error {
		jobStopped := make(chan struct{})
		jobEnd := task.Task(job).
			Defer(func() { close(jobStopped) }).
			Go(ctx)
		appEnd := progAsTask(app).Go(ctx)

		for {
			select {
			case err := <-jobEnd:
				if err == nil {
					w.Write([]byte("Job completed successfully.\n"))
					if !wait {
						app.Quit()
					}
					continue
				}

				w.Write([]byte("Job ended with error: " + err.Error() + "\n"))
				setStatus("Press q or Ctrl+C to exit.")
			case err := <-appEnd:
				<-jobStopped
				return err
			}
		}
	}
}
