// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package pearl

import (
	"slices"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/raohwork/huninn/tapioca"
)

func (l *TaskList) addTask(id, desc string) {
	if _, ok := l.tasks[id]; ok {
		return
	}
	l.tasks[id] = newTaskInfo(id, desc)
	l.pendingTasks = append(l.pendingTasks, id)
}

func (l *TaskList) updateTaskState(id string, state TaskState, progress float64) (ret tea.Cmd) {
	task, ok := l.tasks[id]
	if !ok {
		return nil
	}

	if task.state != state {
		l.recomputeCache(task, state)
		if state == TaskRunning {
			ret = task.spinner.Tick
		}
		task.state = state
	}
	if state == TaskRunning {
		task.progress = min(progress, 1.0)
	}
	return
}

func (l *TaskList) recomputeCache(task *taskInfo, newState TaskState) {
	l.removeFromCache(task.id, task.state)
	l.addToCache(task.id, newState)
}

func (l *TaskList) removeFromCache(id string, state TaskState) {
	// remove from old state
	switch state {
	case TaskPending:
		idx := slices.Index(l.pendingTasks, id)
		if idx >= 0 {
			l.pendingTasks = append(l.pendingTasks[:idx], l.pendingTasks[idx+1:]...)
		}
	case TaskRunning:
		idx := slices.Index(l.runningTasks, id)
		if idx >= 0 {
			l.runningTasks = append(l.runningTasks[:idx], l.runningTasks[idx+1:]...)
		}
	case TaskDone, TaskFailed:
		idx := slices.Index(l.completed, id)
		if idx >= 0 {
			l.completed = append(l.completed[:idx], l.completed[idx+1:]...)
		}
	}
}

func (l *TaskList) addToCache(id string, newState TaskState) {
	// add to new state
	switch newState {
	case TaskPending:
		l.pendingTasks = append(l.pendingTasks, id)
	case TaskRunning:
		l.runningTasks = append(l.runningTasks, id)
	case TaskDone, TaskFailed:
		l.completed = append(l.completed, id)
	}
}

func (l *TaskList) updateTaskDesc(id, desc string) {
	task, ok := l.tasks[id]
	if !ok {
		return
	}

	if task.desc != desc {
		task.desc = desc
	}
}

func (l *TaskList) removeTask(id string) {
	task, ok := l.tasks[id]
	if !ok {
		return
	}

	l.removeFromCache(id, task.state)
	delete(l.tasks, id)
}

func (l *TaskList) Init() tea.Cmd { return nil }

func (l *TaskList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case AddTaskMsg:
		l.addTask(msg.ID, msg.Desc)
		l.recomputeEntries()
	case UpdateTaskStateMsg:
		cmd := l.updateTaskState(msg.ID, msg.State, msg.Progress)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		l.recomputeEntries()
	case UpdateTaskDescMsg:
		l.updateTaskDesc(msg.ID, msg.Desc)
		l.recomputeEntries()
	case RemoveTaskMsg:
		l.removeTask(msg.ID)
		l.recomputeEntries()
	case spinner.TickMsg:
		for _, id := range l.runningTasks {
			if task, ok := l.tasks[id]; ok {
				var cmd tea.Cmd
				task.spinner, cmd = task.spinner.Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		}
	case tapioca.ResizeMsg:
		l.impl.Update(msg)
		l.recomputeEntries()
	default:
		var cmd tea.Cmd
		l.impl, cmd = l.impl.UpdateInto(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) == 0 {
		return l, nil
	}
	l.recomputeEntries()
	return l, tea.Batch(cmds...)
}

func (l *TaskList) View() string { return l.impl.View() }

func (l *TaskList) recomputeEntries() {
	h := l.impl.Height()
	l.impl.Clear()

	rc := len(l.runningTasks)
	pc := len(l.pendingTasks)
	cc := len(l.completed)

	pHeight := min(1, pc) // at least one line for pending if there are any
	cHeight := min(1, cc) // at least one line for completed if there are any
	if h < 3 {
		pHeight = 0
		cHeight = 0
	}
	rHeight := min(rc, h-pHeight-cHeight)
	rest := h - rHeight - pHeight - cHeight
	for rest > 0 {
		update := false
		if pHeight < pc {
			pHeight++
			rest--
			update = true
		}
		if rest > 0 && pHeight < pc {
			pHeight++
			rest--
			update = true
		}
		if rest > 0 && cHeight < cc {
			cHeight++
			rest--
			update = true
		}
		if !update {
			break
		}
	}

	// get last n running tasks
	cList := l.completed[max(0, cc-cHeight):]
	for x := range cHeight {
		task, ok := l.tasks[cList[x]]
		if !ok {
			continue
		}
		l.impl.Append(task.render())
	}

	rList := l.runningTasks[:min(rc, rHeight)]
	for x := range rHeight {
		task, ok := l.tasks[rList[x]]
		if !ok {
			continue
		}
		l.impl.Append(task.render())
	}

	pList := l.pendingTasks[:min(pc, pHeight)]
	for x := range pHeight {
		task, ok := l.tasks[pList[x]]
		if !ok {
			continue
		}
		l.impl.Append(task.render())
	}
}
