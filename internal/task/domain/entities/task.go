package entities

import (
	"ggltask/internal/task"
	"time"
)

type Task struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name"`
	Status    task.TaskStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
