package http

import "ggltask/internal/task/domain/entities"

type ListTasksResponse struct {
	Tasks []*entities.Task `json:"tasks"`
	Total int              `json:"total"`
}

type CreateTaskResponse struct {
	Task *entities.Task `json:"task"`
}

type UpdateTaskResponse struct {
	Task *entities.Task `json:"task"`
}
