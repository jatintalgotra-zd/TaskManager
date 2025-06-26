package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var ErrInvalid = errors.New("invalid task id")

type TaskManager struct {
	db *sql.DB
}

func NewTaskManager() *TaskManager {
	db, err := sql.Open("mysql", "root:root123@tcp(localhost:3306)/test_db")
	if err != nil {
		fmt.Println("Failed to connect to mysql")
		return nil
	}

	return &TaskManager{db: db}
}

// Task struct
// added fields for description and status.
type Task struct {
	ID     int    `json:"id"`
	Desc   string `json:"desc"`
	Status bool   `json:"status"`
}

// AddTask (description string, nextID func() int, mp map[int]*Task)
// adds new Task by generating id.
func (t *TaskManager) AddTask(description string) int {
	result, err := t.db.Exec("INSERT INTO tasks (description, status) VALUES ( ?, ?)", description, false)
	if err != nil {
		fmt.Println("Failed to insert task: ", err)
		return 0
	}

	id, err2 := result.LastInsertId()
	if err2 != nil {
		fmt.Println("Failed to retrieve last id: ", err2)
		return 0
	}

	fmt.Println("Task added:", id, "-", description)

	return int(id)
}

// ListTasks () []string
// prints all pending tasks.
func (t *TaskManager) ListTasks() []Task {
	taskList := make([]Task, 0)

	rows, err := t.db.Query("SELECT id, description, status FROM tasks")
	if err != nil {
		fmt.Println("Failed to list tasks: ", err)
		return nil
	}

	if rows.Err() != nil {
		fmt.Println("Failed to list tasks: ", rows.Err())
	}

	for rows.Next() {
		var task Task

		err2 := rows.Scan(&task.ID, &task.Desc, &task.Status)
		if err2 != nil {
			fmt.Println("Failed to list tasks: ", err2)
			continue
		}

		taskList = append(taskList, task)
	}

	return taskList
}

func (t *TaskManager) ListTaskByID(id int) (Task, error) {
	var task Task

	row := t.db.QueryRow("SELECT id, description, status FROM tasks WHERE id = ?", id)

	err := row.Scan(&task.ID, &task.Desc, &task.Status)
	if err != nil {
		fmt.Println("Failed to list tasks: ", err)

		return Task{}, err
	}

	return task, nil
}

func (t *TaskManager) DeleteTaskByID(id int) error {
	result, err := t.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		fmt.Println("Failed to delete task: ", err)
		return err
	}

	affected, err2 := result.RowsAffected()
	if err2 != nil {
		fmt.Println("Failed to delete task: ", err2)
		return err2
	}

	if affected == 0 {
		return ErrInvalid
	}

	return nil
}

// CompleteTask (int)
// marks Task complete by id.
func (t *TaskManager) CompleteTask(id int) error {
	res, err := t.db.Exec("UPDATE tasks SET status = true WHERE id = ? AND status = false", id)
	if err != nil {
		fmt.Println("Failed to complete task: ", err)
		return err
	}

	affected, err2 := res.RowsAffected()
	if err2 != nil {
		fmt.Println("Failed to complete task: ", err2)
		return err2
	}

	if affected == 0 {
		return ErrInvalid
	}

	return nil
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
