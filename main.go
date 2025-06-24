package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrInvalid = errors.New("invalid task id")

type TaskManager struct {
	tasks  []Task
	nextID func() int
}

func NewTaskManager() *TaskManager {
	return &TaskManager{tasks: make([]Task, 0), nextID: GenID()}
}

// Task struct
// added fields for description and status.
type Task struct {
	Desc   string `json:"desc"`
	Status bool   `json:"status"`
}

// GenID () func() int
// returns function to get next id.
func GenID() func() int {
	id := 0

	return func() int {
		id++
		return id
	}
}

// AddTask (description string, nextID func() int, mp map[int]*Task)
// adds new Task by generating id.
func (t *TaskManager) AddTask(description string) {
	id := t.nextID()
	t1 := Task{Desc: description}
	t.tasks = append(t.tasks, t1)

	fmt.Println("Task added:", id, "-", description)
}

// ListTasks () []string
// prints all pending tasks.
func (t *TaskManager) ListTasks() []Task {
	return t.tasks
}

func (t *TaskManager) ListTaskByID(id int) (Task, error) {
	var err error

	idx := id - 1
	if idx < 0 || idx >= len(t.tasks) {
		err = ErrInvalid
		return Task{}, err
	}

	if t.tasks[idx].Status {
		return Task{}, ErrInvalid
	}

	return t.tasks[idx], nil
}

func (t *TaskManager) DeleteTaskByID(id int) error {
	idx := id - 1
	if idx < 0 || idx >= len(t.tasks) {
		fmt.Println("invalid task id for Delete method")
		return ErrInvalid
	}

	t.tasks[idx] = t.tasks[len(t.tasks)-1]
	t.tasks = t.tasks[:len(t.tasks)-1]

	return nil
}

// CompleteTask (int)
// marks Task complete by id.
func (t *TaskManager) CompleteTask(id int) {
	idx := id - 1
	if idx >= 0 && idx < len(t.tasks) {
		fmt.Println("Marking task", id, "as completed...")

		(t.tasks)[idx].Status = true
	} else {
		fmt.Println("Invalid task ID for Complete method")
	}
}

func main() {
	taskManager := NewTaskManager()
	handler := NewHandler(taskManager)
	http.HandleFunc("GET /api/task", handler.GetAll)
	http.HandleFunc("GET /api/task/{id}", handler.GetByID)
	http.HandleFunc("POST /api/task", handler.PostTask)
	http.HandleFunc("PUT /api/task/{id}", handler.PutTask)
	http.HandleFunc("DELETE /api/task/{id}", handler.DeleteTask)

	svr := http.Server{
		Addr:              ":8000",
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: 20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	err := svr.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server")
	}
}
