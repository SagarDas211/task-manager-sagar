package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"task-manager/internal/handler"
	"task-manager/internal/repository/memory"
	"task-manager/internal/service"
)

func main() {

	// -------- Repository Layer --------
	taskRepo := memory.NewTaskRepository()

	// -------- Service Layer --------
	taskService := service.NewTaskService(taskRepo)

	// -------- Handler Layer --------
	taskHandler := handler.NewTaskHandler(taskService)

	// -------- Router Setup --------
	router := gin.Default()

	// Routes
	router.POST("/tasks", taskHandler.CreateTask)
	router.GET("/tasks/:id", taskHandler.GetTask)
	router.PUT("/tasks/:id", taskHandler.UpdateTask)
	router.DELETE("/tasks/:id", taskHandler.DeleteTask)
	router.GET("/tasks", taskHandler.ListTasks)

	log.Printf("listening on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
