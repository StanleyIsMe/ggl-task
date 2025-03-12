package http

import (
	"fmt"
	"ggltask/internal/task/domain/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type TaskHandler struct {
	taskUsecase usecase.TaskUseCase
}

func NewTaskHandler(taskUsecase usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{taskUsecase: taskUsecase}
}

// @Summary Create task
// @Description Create a new task
// @Tags task
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Create task request"
// @Success 200 {object} CreateTaskResponse "Create task response"
// @Failure 400 {object} ErrorResponse "invalid request"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/v1/tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, InvalidRequestError())

		return
	}

	newTask, err := h.taskUsecase.CreateTask(ctx, usecase.CreateTaskParams{Name: req.Name})
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("task create error")

		c.JSON(UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, CreateTaskResponse{
		Task: newTask,
	})
}

// @Summary List tasks
// @Description List tasks
// @Tags task
// @Accept json
// @Produce json
// @Param request query ListTasksRequest true "List tasks request"
// @Success 200 {object} ListTasksResponse "List tasks response"
// @Failure 400 {object} ErrorResponse "invalid request"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/v1/tasks [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	ctx := c.Request.Context()

	var req ListTasksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, InvalidRequestError())

		return
	}

	result, err := h.taskUsecase.ListTasks(ctx, usecase.ListTasksParams{
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	})
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", req),
			"error":   err,
		}).Msg("task list error")

		c.JSON(UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, ListTasksResponse{
		Tasks: result.Tasks,
		Total: result.Total,
	})
}

// @Summary Update task
// @Description Update a task
// @Tags task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param request body UpdateTaskRequest true "Update task request"
// @Success 200 {object} UpdateTaskResponse "Update task response"
// @Failure 400 {object} ErrorResponse "invalid request"
// @Failure 404 {object} ErrorResponse "not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/v1/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, InvalidRequestError())

		return
	}

	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, InvalidRequestError())

		return
	}

	updateTaskParams := usecase.UpdateTaskParams{
		ID:     uint(idUint),
		Name:   req.Name,
		Status: req.Status,
	}
	updatedTask, err := h.taskUsecase.UpdateTask(ctx, updateTaskParams)
	if err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": fmt.Sprintf("%+v", updateTaskParams),
			"error":   err,
		}).Msg("task update error")

		c.JSON(UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, UpdateTaskResponse{
		Task: updatedTask,
	})
}

// @Summary Delete task
// @Description Delete a task
// @Tags task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} nil "empty result"
// @Failure 400 {object} ErrorResponse "invalid request"
// @Failure 404 {object} ErrorResponse "not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/v1/tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, InvalidRequestError())

		return
	}

	if err := h.taskUsecase.DeleteTask(ctx, uint(idUint)); err != nil {
		zerolog.Ctx(ctx).Error().Fields(map[string]any{
			"payload": id,
			"error":   err,
		}).Msg("task delete error")

		c.JSON(UseCaesErrorToErrorResp(err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
