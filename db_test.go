package main

import "testing"

func TestAddTask(t *testing.T) {
	tm := NewTaskManager()
	id := tm.AddTask("test task")

	task, err := tm.ListTaskByID(id)
	if err != nil {
		t.Error("failed: ", err)
		return
	}

	if task.Desc != "test task" {
		t.Errorf("Expected: test task, got %s", task.Desc)
		return
	}

	// clean
	err2 := tm.DeleteTaskByID(id)
	if err2 != nil {
		t.Error("failed: ", err2)
	}
}

func TestListTasks(t *testing.T) {
	tm := NewTaskManager()
	id1 := tm.AddTask("task 1")
	id2 := tm.AddTask("task 2")

	tasks := tm.ListTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected: 2 tasks, got %d", len(tasks))
		return
	}

	if tasks[0].Desc != "task 1" || tasks[1].Desc != "task 2" {
		t.Errorf("Expected: task 1, task 2 got %s, %s", tasks[0].Desc, tasks[1].Desc)
		return
	}

	// clean
	err := tm.DeleteTaskByID(id1)
	if err != nil {
		t.Error("failed: ", err)
		return
	}

	err2 := tm.DeleteTaskByID(id2)
	if err2 != nil {
		t.Error("failed: ", err2)
	}
}

func TestListTaskByID(t *testing.T) {
	tm := NewTaskManager()
	id1 := tm.AddTask("task 1")

	task1, err := tm.ListTaskByID(id1)
	if err != nil {
		t.Error("failed: ", err)
		return
	}

	if task1.Desc != "task 1" {
		t.Errorf("Expected: task 1, got %s", task1.Desc)
		return
	}

	// clean
	err2 := tm.DeleteTaskByID(id1)
	if err2 != nil {
		t.Error("failed: ", err2)
		return
	}
}

func TestDeleteTaskByID(t *testing.T) {
	tm := NewTaskManager()
	id1 := tm.AddTask("task 1")

	err := tm.DeleteTaskByID(id1)
	if err != nil {
		t.Error("failed: ", err)
		return
	}

	tasks := tm.ListTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected: 0 tasks, got %d", len(tasks))
		return
	}
}

func TestCompleteTask(t *testing.T) {
	tm := NewTaskManager()
	id1 := tm.AddTask("task 1")

	err := tm.CompleteTask(id1)
	if err != nil {
		t.Error("failed: ", err)
	}

	tasks := tm.ListTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected: 1 task, got %d", len(tasks))
	}

	if tasks[0].Desc != "task 1" {
		t.Errorf("Expected: task 1, got %s", tasks[0].Desc)
	}

	if !tasks[0].Status {
		t.Errorf("Expected: true, got %v", tasks[0].Status)
	}

	// clean
	err2 := tm.DeleteTaskByID(id1)
	if err2 != nil {
		t.Error("failed: ", err2)
	}
}
