package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

var errTest = errors.New("test error")

func TestHandler_GetAll(t *testing.T) {
	tm := NewTaskManager()
	tm.AddTask("task 1")
	tm.AddTask("task 2")

	h := NewHandler(tm)

	// testcase 1
	req1 := httptest.NewRequest(http.MethodPut, "/api/task", http.NoBody)
	res1 := httptest.NewRecorder()

	h.GetAll(res1, req1)

	if res1.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d but got %d", http.StatusMethodNotAllowed, res1.Code)
	}

	// testcase 2
	req2 := httptest.NewRequest(http.MethodGet, "/api/task", http.NoBody)
	res2 := httptest.NewRecorder()

	h.GetAll(res2, req2)

	if res2.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, res2.Code)
	}

	var tasks []Task

	err := json.Unmarshal(res2.Body.Bytes(), &tasks)
	if err != nil {
		t.Error("failed unmarshall")
	}

	if len(tasks) != 2 {
		t.Errorf("Expected length %d but got %d", 2, len(tasks))
	}

	if tasks[0].Desc != "task 1" || tasks[1].Desc != "task 2" {
		t.Errorf("Expected: task 1, task 2 got %s, %s", tasks[0].Desc, tasks[1].Desc)
	}

	// testcase 3
	req3 := httptest.NewRequest(http.MethodGet, "/api/task", http.NoBody)
	res3 := errWriter{0}

	h.GetAll(&res3, req3)

	if res3.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d but got %d", http.StatusInternalServerError, res3.Code)
	}
}

type errWriter struct {
	Code int
}

func (*errWriter) Header() http.Header {
	return nil
}
func (*errWriter) Write(_ []byte) (int, error) {
	return 0, errTest
}
func (e *errWriter) WriteHeader(statusCode int) {
	e.Code = statusCode
}

func TestHandler_GetByID(t *testing.T) {
	tm := NewTaskManager()
	tm.AddTask("task 1")
	tm.AddTask("task 2")

	h := NewHandler(tm)

	// testcase 1
	req1 := httptest.NewRequest(http.MethodPut, "/api/task/1", http.NoBody)
	res1 := httptest.NewRecorder()

	h.GetByID(res1, req1)

	if res1.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d but got %d", http.StatusMethodNotAllowed, res1.Code)
	}

	// testcase 2
	req2 := httptest.NewRequest(http.MethodGet, "/api/task/1", http.NoBody)
	res2 := httptest.NewRecorder()

	req2.SetPathValue("id", "1")

	h.GetByID(res2, req2)

	if res2.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, res2.Code)
	}

	var task1 Task

	err := json.Unmarshal(res2.Body.Bytes(), &task1)
	if err != nil {
		t.Error("failed unmarshal")
	}

	if task1.Desc != "task 1" {
		t.Errorf("Expected: task 1 got %s", task1.Desc)
	}

	// testcase 3
	tm.ListTasks()[0].Status = true
	req3 := httptest.NewRequest(http.MethodGet, "/api/task/1", http.NoBody)
	res3 := httptest.NewRecorder()

	req3.SetPathValue("id", "1")

	h.GetByID(res3, req3)

	if res3.Code != http.StatusBadRequest {
		t.Errorf("Expected %d got %d", http.StatusBadRequest, res3.Code)
	}

	// testcase 4
	req4 := httptest.NewRequest(http.MethodGet, "/api/task/1", http.NoBody)
	res4 := httptest.NewRecorder()

	req4.SetPathValue("id", "abc")

	h.GetByID(res4, req4)

	if res4.Code != http.StatusBadRequest {
		t.Errorf("Expected %d got %d", http.StatusBadRequest, res4.Code)
	}

	// testcase 5
	req5 := httptest.NewRequest(http.MethodGet, "/api/task/1", http.NoBody)
	res5 := httptest.NewRecorder()

	req5.SetPathValue("id", "10")

	h.GetByID(res5, req5)

	if res5.Code != http.StatusBadRequest {
		t.Errorf("Expected %d got %d", http.StatusBadRequest, res5.Code)
	}

	// testcase 6
	req6 := httptest.NewRequest(http.MethodGet, "/api/task/1", http.NoBody)
	res6 := errWriter{0}

	req6.SetPathValue("id", "2")

	h.GetByID(&res6, req6)

	if res6.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d but got %d", http.StatusInternalServerError, res6.Code)
	}
}

func TestHandler_PostTask(t *testing.T) {
	tm := NewTaskManager()
	h := NewHandler(tm)

	// testcase 1
	req1 := httptest.NewRequest(http.MethodGet, "/api/task", http.NoBody)
	res1 := httptest.NewRecorder()

	h.PostTask(res1, req1)

	if res1.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d but got %d", http.StatusMethodNotAllowed, res1.Code)
	}

	// testcase 2
	newTask := Task{Desc: "new task", Status: false}

	reqBytes, err := json.Marshal(newTask)
	if err != nil {
		t.Error("failed marshal")
	}

	reader := bytes.NewReader(reqBytes)

	req2 := httptest.NewRequest(http.MethodPost, "/api/task", reader)
	res2 := httptest.NewRecorder()

	h.PostTask(res2, req2)

	if res2.Code != http.StatusCreated {
		t.Errorf("Expected %d got %d", http.StatusCreated, res2.Code)
	}

	// testcase 3
	newTask2 := []struct {
		id     int64
		name   string
		height float64
	}{
		{
			id:     1,
			name:   "jack",
			height: 5.7,
		},
		{
			id:     2,
			name:   "jack2",
			height: 6.7,
		},
	}

	reqBytes2, err2 := json.Marshal(newTask2)
	if err2 != nil {
		t.Error("failed marshal")
	}

	reader2 := bytes.NewReader(reqBytes2)

	req3 := httptest.NewRequest(http.MethodPost, "/api/task", reader2)
	res3 := httptest.NewRecorder()

	h.PostTask(res3, req3)

	if res3.Code != http.StatusBadRequest {
		t.Errorf("Expected %d got %d", http.StatusBadRequest, res2.Code)
	}

	// testcase 4
	req4 := httptest.NewRequest(http.MethodPost, "/api/task", errReader(0))
	res4 := httptest.NewRecorder()

	h.PostTask(res4, req4)

	if res4.Code != http.StatusBadRequest {
		t.Errorf("Expected %d got %d", http.StatusBadRequest, res4.Code)
	}
}

type errReader int

func (errReader) Read(_ []byte) (n int, err error) {
	return 0, errTest
}

func TestHandler_PutTask(t *testing.T) {
	tm := NewTaskManager()
	tm.AddTask("task 1")
	tm.AddTask("task 2")
	h := NewHandler(tm)

	// testcase 1
	req1 := httptest.NewRequest(http.MethodGet, "/api/task/1", http.NoBody)
	res1 := httptest.NewRecorder()

	h.PutTask(res1, req1)

	if res1.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d but got %d", http.StatusMethodNotAllowed, res1.Code)
	}

	// testcase 2
	req2 := httptest.NewRequest(http.MethodPut, "/api/task/1", http.NoBody)
	res2 := httptest.NewRecorder()

	req2.SetPathValue("id", "abc")

	h.PutTask(res2, req2)

	if res2.Code != http.StatusBadRequest {
		t.Errorf("Expected %d got %d", http.StatusBadRequest, res2.Code)
	}

	// testcase 3
	req3 := httptest.NewRequest(http.MethodPut, "/api/task/1", http.NoBody)
	res3 := httptest.NewRecorder()

	req3.SetPathValue("id", "1")

	h.PutTask(res3, req3)

	if res3.Code != http.StatusOK {
		t.Errorf("Expected %d got %d", http.StatusOK, res3.Code)
	}

	if !tm.ListTasks()[0].Status {
		t.Errorf("Expected status true got %v", tm.ListTasks()[1].Status)
	}

	// testcase 4
	req4 := httptest.NewRequest(http.MethodPut, "/api/task/1", http.NoBody)
	res4 := httptest.NewRecorder()

	req4.SetPathValue("id", "3")

	h.PutTask(res4, req4)

	if res4.Code != http.StatusBadRequest {
		t.Errorf("Expected %d got %d", http.StatusBadRequest, res4.Code)
	}
}

func TestHandler_DeleteTask(t *testing.T) {
	tm := NewTaskManager()
	tm.AddTask("task 1")
	tm.AddTask("task 2")
	h := NewHandler(tm)

	// testcase 1
	req1 := httptest.NewRequest(http.MethodGet, "/api/task/1", http.NoBody)
	res1 := httptest.NewRecorder()

	h.DeleteTask(res1, req1)

	if res1.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d but got %d", http.StatusMethodNotAllowed, res1.Code)
	}

	// testcase 2
	req2 := httptest.NewRequest(http.MethodDelete, "/api/task/1", http.NoBody)
	res2 := httptest.NewRecorder()

	req2.SetPathValue("id", "abc")

	h.DeleteTask(res2, req2)

	if res2.Code != http.StatusBadRequest {
		t.Errorf("Expected %d got %d", http.StatusBadRequest, res2.Code)
	}

	// testcase 3
	req3 := httptest.NewRequest(http.MethodDelete, "/api/task/1", http.NoBody)
	res3 := httptest.NewRecorder()

	req3.SetPathValue("id", "3")

	h.DeleteTask(res3, req3)

	if res3.Code != http.StatusBadRequest {
		t.Errorf("Expected %d got %d", http.StatusBadRequest, res3.Code)
	}

	// testcase 4
	req4 := httptest.NewRequest(http.MethodDelete, "/api/task/1", http.NoBody)
	res4 := httptest.NewRecorder()

	req4.SetPathValue("id", "1")

	h.DeleteTask(res4, req4)

	if res4.Code != http.StatusOK {
		t.Errorf("Expected %d got %d", http.StatusOK, res4.Code)
	}

	if tm.ListTasks()[0].Desc == "task 1" {
		t.Errorf("Expected after delete %s got %v", "task 2", tm.ListTasks()[0].Desc)
	}
}
