// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pearl

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
)

// TaskController is used to control what info is shown for a task.
type TaskController interface {
	SetDesc(desc string)
	SetState(state TaskState, progress float64)
	Done()
	Fail()
	Remove()
}

type taskController struct {
	tid  int64
	send func(tea.Msg)
	id   string
}

func (tc *taskController) SetDesc(desc string) {
	tc.send(UpdateTaskDescMsg{TaskListID: tc.tid, ID: tc.id, Desc: desc})
}

func (tc *taskController) SetState(state TaskState, progress float64) {
	if !state.IsValid() {
		return
	}
	tc.send(UpdateTaskStateMsg{
		TaskListID: tc.tid,
		ID:         tc.id,
		State:      state,
		Progress:   progress,
	})
}

func (tc *taskController) Done() {
	tc.send(UpdateTaskStateMsg{
		TaskListID: tc.tid,
		ID:         tc.id,
		State:      TaskDone,
		Progress:   1,
	})
}

func (tc *taskController) Fail() {
	tc.send(UpdateTaskStateMsg{
		TaskListID: tc.tid,
		ID:         tc.id,
		State:      TaskFailed,
		Progress:   -1,
	})
}

func (tc *taskController) Remove() {
	tc.send(RemoveTaskMsg{TaskListID: tc.tid, ID: tc.id})
}

// TaskManager is used to create and manage tasks.
type TaskManager interface {
	AddTask(desc, id string) TaskController
}

type taskManager struct {
	id   int64
	send func(tea.Msg)
}

func (tm taskManager) AddTask(desc, id string) TaskController {
	if id == "" {
		id = "task#" + strconv.FormatInt(tapioca.NewID(), 10)
	}
	tm.send(AddTaskMsg{TaskListID: tm.id, ID: id, Desc: desc})
	return &taskController{tid: tm.id, send: tm.send, id: id}
}

// CreateManager creates a new TaskManager.
func (tl *TaskList) CreateManager(send func(tea.Msg)) TaskManager {
	return taskManager{id: tl.id, send: send}
}
