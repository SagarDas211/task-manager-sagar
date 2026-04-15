package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"task-manager/internal/repository/memory"
	"task-manager/internal/service"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	repo := memory.NewTaskRepository()
	svc := service.NewTaskService(repo)
	handler := NewTaskHandler(svc)

	r := gin.Default()

	r.POST("/tasks", handler.CreateTask)
	r.GET("/tasks/:id", handler.GetTask)
	r.PUT("/tasks/:id", handler.UpdateTask)
	r.DELETE("/tasks/:id", handler.DeleteTask)
	r.GET("/tasks", handler.ListTasks)

	return r
}

func TestCreateTask(t *testing.T) {
	router := setupRouter()

	body := map[string]interface{}{
		"title":    "Test Task",
		"due_date": "2030-01-01",
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestGetTask(t *testing.T) {
	router := setupRouter()

	// Create first
	body := `{"title":"Test","due_date":"2030-01-01"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	id := resp["id"].(string)

	// Now GET
	req, _ = http.NewRequest("GET", "/tasks/"+id, nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateTask(t *testing.T) {
	router := setupRouter()

	// Create
	body := `{"title":"Test","due_date":"2030-01-01"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	id := resp["id"].(string)

	// Update
	updateBody := `{"title":"Updated"}`
	req, _ = http.NewRequest("PUT", "/tasks/"+id, bytes.NewBuffer([]byte(updateBody)))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	router := setupRouter()

	// Create
	body := `{"title":"Test","due_date":"2030-01-01"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	id := resp["id"].(string)

	// Delete
	req, _ = http.NewRequest("DELETE", "/tasks/"+id, nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestListTasks(t *testing.T) {
	router := setupRouter()

	// Create multiple
	for i := 0; i < 2; i++ {
		body := `{"title":"Task","due_date":"2030-01-01"}`
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// List
	req, _ := http.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var tasks []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &tasks)

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}
