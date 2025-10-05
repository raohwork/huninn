// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pearl

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/raohwork/huninn/tapioca"
)

// TaskState represents the current state of a task.
type TaskState int

const (
	TaskPending TaskState = iota // task is pending
	TaskRunning                  // task is currently running
	TaskDone                     // task has completed successfully
	TaskFailed                   // task has failed
	invalidTaskState
)

func (s TaskState) IsValid() bool {
	return s >= TaskPending && s <= TaskFailed
}

type taskInfo struct {
	id       string
	desc     string
	state    TaskState
	progress float64
	spinner  spinner.Model
}

func newTaskInfo(id, desc string) *taskInfo {
	return &taskInfo{
		id:       id,
		desc:     desc,
		state:    TaskPending,
		progress: -1.0,
		spinner:  spinner.New(spinner.WithSpinner(spinner.Dot)),
	}
}

func (i *taskInfo) render() string {
	b := &strings.Builder{}
	b.Grow(len(i.desc) + 20)
	// icon (emoji)
	switch i.state {
	case TaskPending:
		b.WriteString(lipgloss.NewStyle().
			SetString(`🕓 `).
			Foreground(lipgloss.Color("15")).
			String())
	case TaskDone:
		b.WriteString(lipgloss.NewStyle().
			SetString(`✅ `).
			Foreground(lipgloss.Color("10")).
			String())
	case TaskFailed:
		b.WriteString(lipgloss.NewStyle().
			SetString(`❌ `).
			Foreground(lipgloss.Color("9")).
			String())
	case TaskRunning:
		b.WriteString(lipgloss.NewStyle().
			SetString(i.spinner.View()).
			Foreground(lipgloss.Color("14")).
			String())
	}

	// pad space as separator
	if l := lipgloss.Width(b.String()); l < 3 {
		b.Write([]byte(strings.Repeat(" ", 3-l)))
	}

	// progress counter
	if i.state == TaskRunning && i.progress >= 0.0 {
		if i.progress > 1.0 {
			i.progress = 1.0
		}
		fmt.Fprintf(b, "[%6.2f%%] ", i.progress*100.0)
	}

	// description
	b.WriteString(i.desc)

	return b.String()
}

// TaskList is a component that manages and displays a list of tasks with their states.
//
// You might send task message by your own, or use [TaskManager].
type TaskList struct {
	tasks map[string]*taskInfo
	impl  *Block
	id    int64

	// cached info
	pendingTasks []string // indexes of pending tasks
	runningTasks []string // indexes of running tasks
	completed    []string // indexes of completed tasks (done or failed)
}

func (t *TaskList) ID() int64 { return t.id }

// NewTaskList creates a new TaskList component.
func NewTaskList() *TaskList {
	return &TaskList{
		id:    tapioca.NewID(),
		tasks: make(map[string]*taskInfo),
		impl:  NewBlock(),
	}
}

// AddTaskMsg is a message to add a new task.
//
// If the id already exists, the task will not be added.
type AddTaskMsg struct {
	TaskListID int64
	ID         string
	Desc       string
}

// UpdateTaskStateMsg is a message to update the state of a task.
//
// If the id does not exist, the message will be ignored.
type UpdateTaskStateMsg struct {
	TaskListID int64
	ID         string
	State      TaskState

	// ignored if State is not TaskRunning
	// value between 0.0 and 1.0
	// if > 1.0, it will be set to 1.0
	// if < 0.0, hide progress counter
	Progress float64
}

// UpdateTaskDescMsg is a message to update the description of a task.
//
// If the id does not exist, the message will be ignored.
type UpdateTaskDescMsg struct {
	TaskListID int64
	ID         string
	Desc       string
}

// RemoveTaskMsg is a message to remove a task.
//
// If the id does not exist, the message will be ignored.
type RemoveTaskMsg struct {
	TaskListID int64
	ID         string
}
