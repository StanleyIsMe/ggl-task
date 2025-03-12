package http

import "ggltask/internal/task"

type CreateTaskRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

type UpdateTaskRequest struct {
	Name   string          `json:"name" binding:"required,max=50"`
	Status task.TaskStatus `json:"status" binding:"oneof=0 1"`
}

type ListTasksRequest struct {
	PageIndex int `form:"page_index,default=1" binding:"required,gte=1"`
	PageSize  int `form:"page_size,default=10" binding:"required,gte=1,lte=100"`
}
