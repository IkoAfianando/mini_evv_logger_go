package handler

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"evv-logger-backend/internal/models"
	"evv-logger-backend/internal/store"
)

type TaskHandler struct {
	store *store.Store
}

func NewTaskHandler(st *store.Store) *TaskHandler {
	return &TaskHandler{store: st}
}

// UpdateTask handles updating the status of a specific task.
// @Summary      Update a task status
// @Description  Updates the status of a specific task to "completed" or "not_completed".
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        taskId   path      int  true  "Task ID"
// @Param        update   body models.UpdateTaskRequest true "Task Status Update"
// @Success      200  {object}  models.Task
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/tasks/{taskId}/update [put]
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	taskIDStr := c.Params("taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	var req models.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	var updatedTask *models.Task
	var found bool

	for _, schedule := range h.store.Schedules {
		for i, task := range schedule.Tasks {
			if task.ID == taskID {
				schedule.Tasks[i].Completed = req.Completed
				if req.NotCompletedReason != "" {
					schedule.Tasks[i].NotCompletedReason = req.NotCompletedReason
				} else {
					schedule.Tasks[i].NotCompletedReason = ""
				}

				updatedTask = &schedule.Tasks[i]
				found = true
				log.Printf("Task %d in Schedule %s updated: completed=%v, reason=%s", taskID, schedule.ID, schedule.Tasks[i].Completed, schedule.Tasks[i].NotCompletedReason)
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found in any schedule"})
	}

	return c.JSON(updatedTask)
}
