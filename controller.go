package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	tm *TaskManager
}

func NewHandler(taskManager *TaskManager) *Handler {
	return &Handler{tm: taskManager}
}

func (h *Handler) GetAll(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		pendingTasks := h.tm.ListTasks()

		pBytes, err := json.Marshal(pendingTasks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err2 := w.Write(pBytes)
		if err2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := req.PathValue("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err2 := h.tm.ListTaskByID(int(idInt))
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tBytes, err3 := json.Marshal(task)
	if err3 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err4 := w.Write(tBytes)
	if err4 != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) PostTask(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		tBytes, err := io.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var task Task

		err2 := json.Unmarshal(tBytes, &task)
		if err2 != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.tm.AddTask(task.Desc)
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) PutTask(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := req.PathValue("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err2 := h.tm.CompleteTask(int(idInt))
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := req.PathValue("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err2 := h.tm.DeleteTaskByID(int(idInt))
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
