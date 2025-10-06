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

// noSuguarComponent provides second-to-none features.
type noSuguarComponent struct {
	tea.Model
}

func (m noSuguarComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Model, cmd = m.Model.Update(tapioca.ResizeMsg{
			Width:  msg.Width,
			Height: msg.Height,
		})
	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" {
			cmd = tea.Quit
		}
	default:
		m.Model, cmd = m.Model.Update(msg)
	}

	return m, cmd
}

// NSNIComponent returns simplified flavor of huninn UI.
//
// The NSNI (No Sugar, No Ice) UI consists of three parts without any decorations:
//   - A status bar at the bottom, showing a single line of text.
//   - A task list at the top, showing multiple tasks with progress in percentage.
//   - A log panel in the middle, showing log messages.
//
// It maximizes the use of screen space, but can be less visually appealing than
// the xxLI (Light Ice) component.
//
// The tlSize parameter specifies the height of the task list in rows. The
// logBufferSize parameter specifies the maximum number of log entries to keep.
// If tlSize is less than 3, the task list and log panel will have equal height.
// The logBufferSize is automatically adjusted to be at least 10.
//
// The component also handles two shortcuts to terminate the program:
//   - Ctrl+C
//   - q
func NSNIComponent(tlSize, logBufferSize int) (
	m tea.Model, f func(func(tea.Msg)) (
		setStatus func(string),
		taskManager pearl.TaskManager,
		w io.Writer,
		logScroller tapioca.ScrollController,
	),
) {
	status := pearl.NewSpan()
	tl := pearl.NewTaskList()
	logger := pearl.NewLogPanel(max(logBufferSize, 10))

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

	return noSuguarComponent{root}, func(send func(tea.Msg)) (setStatus func(string), taskManager pearl.TaskManager, w io.Writer, logScroller tapioca.ScrollController) {
		return status.Setter(send),
			tl.CreateManager(send),
			logger.CreateWriter(send, nil),
			logger.ScrollController()
	}
}

func nsni(tlSize, logBufferSize int, opts ...tea.ProgramOption) (
	prog *tea.Program,
	setStatus func(string),
	tm pearl.TaskManager,
	w io.Writer,
	s tapioca.ScrollController,
) {
	m, f := NSNIComponent(tlSize, logBufferSize)
	prog = tea.NewProgram(m, opts...)

	setStatus, tm, w, s = f(prog.Send)
	return
}

func progAsTask(app *tea.Program) task.Task {
	return task.FromServer(func() error {
		_, err := app.Run()
		return err
	}, app.Quit)
}

// NSNI (No Sugar, No Ice) wraps NSNIComponent to provide a ready-to-run program.
//
// Cancelling the context will stop the program, making it suitable for
// cooperating with signal.NotifyContext.
//
// The tlSize and logBufferSize parameters are passed to [NSNIComponent].
//
// The opts parameters are passed to [tea.NewProgram].
func NSNI(tlSize, logBufferSize int, opts ...tea.ProgramOption) (
	prog func(context.Context) error,
	setStatus func(string),
	tm pearl.TaskManager,
	w io.Writer,
	s tapioca.ScrollController,
) {
	app, setStatus, tm, w, s := nsni(tlSize, logBufferSize, opts...)
	prog = progAsTask(app)
	return
}

// JobFactory wraps you program logic to create a job function.
//
// It accepts four parameters:
//   - setStatus: a function to update the status bar.
//   - tskManager: a task manager to control task list component.
//   - w: an io.Writer to write log messages.
//   - logScroller: a scroll controller to control log panel component.
//   - quit: a function to terminate the UI program.
//
// You SHOULD NOT use quit() in most cases, returning from the job
// function is enough. But when you have to, you MUST remember returning
// from the job function after calling quit() ASAP, or the program
// may hang.
type JobFactory func(
	setStatus func(string),
	tskManager pearl.TaskManager,
	w io.Writer,
	logScroller tapioca.ScrollController,
	quit func(),
) func(context.Context) error

// LSNI (Light Sugar, No Ice) adds some features to NSNI.
//
// It accepts a JobFactory to create a job function, which will be run
// concurrently with the UI program. The job function is the main program
// of your application, and it should utilize the provided setStatus,
// task manager, and writer to interact with the UI.
//
// The tlSize and logBufferSize parameters are passed to [NSNIComponent].
//
// The opts parameters are passed to [tea.NewProgram].
//
// If you set wait to true, the UI will remain active after the job
// completes successfully, allowing the user to review the final status
// and logs. UI always remain active if the job ends with an error.
func LSNI(tlSize, logBufferSize int, factory JobFactory, wait bool, opts ...tea.ProgramOption) func(context.Context) error {
	app, setStatus, tm, w, s := nsni(tlSize, logBufferSize, opts...)
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
