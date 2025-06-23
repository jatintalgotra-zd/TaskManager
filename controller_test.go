package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetAll(t *testing.T) {
	tm := NewTaskManager()
	h := NewHandler(tm)
	server := httptest.NewServer(http.HandlerFunc(h.GetAll))
	url := server.URL + "/api/task"

	req, err0 := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)

	if err0 != nil {
		t.Error(err0)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	err2 := res.Body.Close()
	if err2 != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %v but got %v", 200, res.StatusCode)
	}
}

//func TestHandler_PostTask(t *testing.T) {
//	tm := NewTaskManager()
//	h := NewHandler(tm)
//
//}
