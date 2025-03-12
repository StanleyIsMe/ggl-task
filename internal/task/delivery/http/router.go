package http

import (
	"ggltask/internal/task/domain/usecase"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(router *gin.Engine, taskUsecase usecase.TaskUseCase) {
	taskHandler := NewTaskHandler(taskUsecase)

	v1 := router.Group("/api/v1")
	v1.POST("/tasks", taskHandler.CreateTask)
	v1.GET("/tasks", taskHandler.ListTasks)
	v1.PUT("/tasks/:id", taskHandler.UpdateTask)
	v1.DELETE("/tasks/:id", taskHandler.DeleteTask)
}
