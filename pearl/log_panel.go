// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pearl

import (
	"bytes"
	"io"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
)

// LogPanel displays log messages in an area.
//
// You must create LogPanel with NewLogPanel().
//
// LogPanel does not support manually scrolling.
type LogPanel struct {
	// if true, new log messages are placed at the top
	// by default, new log messages are placed at the bottom
	Reverse bool

	impl *tapioca.Component
}

// LogMsg denotes a logger has written a log message to LogPanel.
type LogMsg []byte

// NewLogPanel creates a new LogPanel.
func NewLogPanel(size int) *LogPanel {
	if size < 10 {
		size = 10
	}
	lp := &LogPanel{
		impl: tapioca.NewComponent(size, tapioca.VerticalScrollable()),
	}
	return lp
}

func (lp *LogPanel) Init() tea.Cmd {
	return lp.impl.Init()
}

func (lp *LogPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return lp.UpdateInto(msg)
}

// UpdateInto is identical to Update, but returns *LogPanel instead of tea.Model.
func (lp *LogPanel) UpdateInto(msg tea.Msg) (*LogPanel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case LogMsg:
		lines := bytes.Split(msg, []byte{'\n'})
		add := lp.impl.Append
		if lp.Reverse {
			add = lp.impl.Prepend
		}

		for _, line := range lines {
			add(string(line))
		}

		if !lp.Reverse {
			var cmd tea.Cmd
			lp.impl, cmd = lp.impl.UpdateInto(tapioca.ScrollBottomMsg{})
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	case tapioca.ResizeMsg:
		var cmd tea.Cmd
		lp.impl, cmd = lp.impl.UpdateInto(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		if !lp.Reverse {
			lp.impl, cmd = lp.impl.UpdateInto(tapioca.ScrollBottomMsg{})
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	default:
		var cmd tea.Cmd
		lp.impl, cmd = lp.impl.UpdateInto(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) == 0 {
		return lp, nil
	}
	return lp, tea.Batch(cmds...)
}

func (lp *LogPanel) View() string {
	return lp.impl.View()
}

type logWriterImpl struct {
	also io.Writer
	lp   *LogPanel
	send func(tea.Msg)
	lock sync.Mutex
}

// CreateWriter returns an io.Writer that writes log messages to the given LogPanel.
//
// The send function is used to send LogMsg messages to the Bubble Tea program.
//
// If also is not nil, the returned io.Writer also writes to also.
//
// You can use LogPanelWriter like this:
//
//	var logPanel *pearl.LogPanel
//	var logFile *os.File
//	...
//	logger := log.New(pearl.LogPanelWriter(logPanel, send, logFile), "", log.LstdFlags)
//	logger.Println("This log message is written to both logPanel and logFile")
func (lp *LogPanel) CreateWriter(send func(tea.Msg), also io.Writer) io.Writer {
	if lp == nil {
		panic("lp is nil")
	}
	if send == nil {
		panic("send is nil")
	}
	return &logWriterImpl{
		lp:   lp,
		send: send,
		also: also,
	}
}

func (w *logWriterImpl) WriteString(p string) (n int, err error) {
	if w.also == nil {
		return w.Write([]byte(p))
	}

	ww, ok := w.also.(io.StringWriter)
	if !ok {
		return w.Write([]byte(p))
	}

	w.lock.Lock()
	n, err = ww.WriteString(p)
	w.lock.Unlock()

	if err == nil {
		w.send(LogMsg([]byte(strings.TrimRight(p, "\n"))))
	}
	return
}

func (w *logWriterImpl) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	n = len(p)
	if w.also != nil {
		n, err = w.also.Write(p)
	}
	w.lock.Unlock()

	if err == nil {
		w.send(LogMsg(bytes.TrimRight(p, "\n")))
	}

	return
}
