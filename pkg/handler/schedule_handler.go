package handler

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/IkoAfianando/mini_evv_logger_go/pkg/models"
	"github.com/IkoAfianando/mini_evv_logger_go/pkg/store"
)

type ScheduleHandler struct {
	store *store.Store
}

func NewScheduleHandler(st *store.Store) *ScheduleHandler {
	return &ScheduleHandler{store: st}
}

// ResetStore handles resetting the in-memory data store to its initial state.
// @Summary      Reset data store
// @Description  Resets the in-memory data to the initial set of schedules and tasks, useful for testing.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/reset [post]
func (h *ScheduleHandler) ResetStore(c *fiber.Ctx) error {
	log.Println("Received request to reset data store.")
	h.store.SetupInitialData()
	log.Println("In-memory data store has been reset.")
	return c.JSON(fiber.Map{"message": "Data store has been reset to initial state"})
}

// GetSchedules handles fetching all caregiver schedules, sorted by date and time.
// @Summary      Get all schedules
// @Description  Fetches a list of all caregiver schedules, sorted chronologically
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Schedule
// @Router       /api/schedules [get]
func (h *ScheduleHandler) GetSchedules(c *fiber.Ctx) error {
	schedulesList := make([]*models.Schedule, 0, len(h.store.Schedules))
	for _, schedule := range h.store.Schedules {
		schedulesList = append(schedulesList, schedule)
	}

	sort.Slice(schedulesList, func(i, j int) bool {
		timeStrI := schedulesList[i].ShiftDate + " " + strings.Split(schedulesList[i].ShiftTime, " - ")[0] + " " + schedulesList[i].AmOrPm
		timeStrJ := schedulesList[j].ShiftDate + " " + strings.Split(schedulesList[j].ShiftTime, " - ")[0] + " " + schedulesList[j].AmOrPm

		layout := "2006-01-02 3:04 PM"

		timeI, errI := time.Parse(layout, timeStrI)
		if errI != nil {
			log.Printf("Error parsing time for schedule %s: %v", schedulesList[i].ID, errI)
			return false
		}

		timeJ, errJ := time.Parse(layout, timeStrJ)
		if errJ != nil {
			log.Printf("Error parsing time for schedule %s: %v", schedulesList[j].ID, errJ)
			return true
		}
		return timeI.Before(timeJ)
	})

	return c.JSON(schedulesList)
}

// GetTodaySchedules handles fetching all schedules for the current date.
// @Summary      Get today's schedules
// @Description  Fetches all schedules scheduled for the current date
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Schedule
// @Router       /api/schedules/today [get]
func (h *ScheduleHandler) GetTodaySchedules(c *fiber.Ctx) error {
	today := time.Now().Format("2006-01-02")
	todaySchedules := make([]*models.Schedule, 0)

	for _, schedule := range h.store.Schedules {
		if schedule.ShiftDate == today {
			todaySchedules = append(todaySchedules, schedule)
		}
	}
	return c.JSON(todaySchedules)
}

// GetScheduleByID handles fetching a single schedule by its ID.
// @Summary      Get schedule by ID
// @Description  Fetches the details of a single schedule using its ID
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Schedule ID"
// @Success      200  {object}  models.Schedule
// @Failure      404  {object}  map[string]string
// @Router       /api/schedules/{id} [get]
func (h *ScheduleHandler) GetScheduleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	schedule, ok := h.store.Schedules[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Schedule with ID %s not found", id),
		})
	}
	return c.JSON(schedule)
}

// StartVisit handles the start of a visit.
// @Summary      Start a visit
// @Description  Marks a scheduled visit as "in_progress" and records the start time and location.
// @Tags         Visits
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Schedule ID"
// @Param        location body models.StartVisitRequest true "Start Location"
// @Success      200  {object}  models.Schedule
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/schedules/{id}/start [post]
func (h *ScheduleHandler) StartVisit(c *fiber.Ctx) error {
	id := c.Params("id")
	schedule, ok := h.store.Schedules[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
	}

	var req models.StartVisitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	now := time.Now()
	schedule.Status = "in_progress"
	schedule.ClockInTime = &now
	schedule.ClockInLocation = &models.Geolocation{
		Latitude:  req.Location.Latitude,
		Longitude: req.Location.Longitude,
	}

	log.Printf("Started visit for schedule ID %s at %v", id, now)
	return c.JSON(schedule)
}

// EndVisit handles the end of a visit.
// @Summary      End a visit
// @Description  Marks an in-progress visit as "completed" and records the end time and location.
// @Tags         Visits
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Schedule ID"
// @Param        location body models.EndVisitRequest true "End Location"
// @Success      200  {object}  models.Schedule
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/schedules/{id}/end [post]
func (h *ScheduleHandler) EndVisit(c *fiber.Ctx) error {
	id := c.Params("id")
	schedule, ok := h.store.Schedules[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
	}

	var req models.EndVisitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	now := time.Now()
	schedule.Status = "completed"
	schedule.ClockOutTime = &now
	schedule.ClockOutLocation = &models.Geolocation{
		Latitude:  req.Location.Latitude,
		Longitude: req.Location.Longitude,
	}

	log.Printf("Ended visit for schedule ID %s at %v", id, now)
	return c.JSON(schedule)
}

// ClockIn handles clocking in for a schedule.
// @Summary      Clock in for a schedule
// @Description  Records the clock-in time and location for a schedule.
// @Tags         Visits
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Schedule ID"
// @Success      200  {object}  models.Schedule
// @Failure      404  {object}  map[string]string
// @Router       /api/schedules/{id}/clock-in [get]
func (h *ScheduleHandler) ClockIn(c *fiber.Ctx) error {
	id := c.Params("id")
	schedule, ok := h.store.Schedules[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
	}

	if schedule.ClockInTime != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Already clocked in"})
	}

	now := time.Now()
	schedule.ClockInTime = &now
	schedule.ClockInLocation = &models.Geolocation{
		Latitude:  0, // Placeholder, should be replaced with actual location
		Longitude: 0, // Placeholder, should be replaced with actual location
	}
	schedule.Status = "in_progress"

	log.Printf("Clocked in for schedule ID %s at %v", id, now)
	return c.JSON(schedule)
}

// CancelClockIn handles cancellation of a previously recorded clock-in.
// @Summary      Cancel clock-in
// @Description  Cancels the clock-in by clearing time and location, and sets status back to "scheduled"
// @Tags         Visits
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Schedule ID"
// @Success      200  {object}  models.Schedule
// @Failure      404  {object}  map[string]string
// @Router       /api/schedules/{id}/cancel-clock-in [post]
func (h *ScheduleHandler) CancelClockIn(c *fiber.Ctx) error {
	id := c.Params("id")
	schedule, ok := h.store.Schedules[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
	}

	schedule.ClockInTime = nil
	schedule.ClockInLocation = nil
	schedule.Status = "scheduled"

	log.Printf("Cancelled clock-in for schedule ID %s", id)
	return c.JSON(schedule)
}

// AddTaskToSchedule adds a new task to a schedule.
// @Summary      Add a task to schedule
// @Description  Adds a new task with name and description to the given schedule
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Schedule ID"
// @Param        task body models.AddTaskRequest true "Task to add"
// @Success      200  {object}  models.Schedule
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/schedules/{id}/tasks [post]
func (h *ScheduleHandler) AddTaskToSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	schedule, ok := h.store.Schedules[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
	}

	var req models.AddTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	newTaskID := len(schedule.Tasks) + 1
	newTask := models.Task{
		ID:          newTaskID,
		Name:        req.Name,
		Description: req.Description,
	}

	schedule.Tasks = append(schedule.Tasks, newTask)

	log.Printf("Added task to schedule ID %s: %+v", id, newTask)
	return c.JSON(schedule)
}
