package main

import (
	"bytes"
	"encoding/json"
	"evv-logger-backend/internal/models"
	"evv-logger-backend/internal/router"
	"evv-logger-backend/internal/store"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	_ "time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTest() (*fiber.App, *store.Store) {
	dataStore := store.NewStore()
	dataStore.SetupInitialData()

	app := fiber.New()
	router.SetupRoutes(app, dataStore)

	return app, dataStore
}

func TestScheduleHandlers(t *testing.T) {
	app, _ := setupTest()

	t.Run("Get All Schedules - Success and Sorted", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/schedules", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var schedules []models.Schedule
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &schedules)

		assert.Len(t, schedules, 6)

		assert.Equal(t, "1", schedules[0].ID)
		assert.Equal(t, "3", schedules[1].ID)
		assert.Equal(t, "2", schedules[2].ID)
		assert.Equal(t, "4", schedules[3].ID)
	})

	t.Run("Get Today's Schedules - Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/schedules/today", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var schedules []models.Schedule
		json.NewDecoder(resp.Body).Decode(&schedules)

		assert.Len(t, schedules, 6)
	})

	t.Run("Get Schedule By ID - Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/schedules/2", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var schedule models.Schedule
		json.NewDecoder(resp.Body).Decode(&schedule)
		assert.Equal(t, "John Doe", schedule.ClientName)
	})

	t.Run("Get Schedule By ID - Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/schedules/999", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestTaskHandlers(t *testing.T) {
	app, dataStore := setupTest()

	t.Run("Update Task - Success", func(t *testing.T) {
		updateBody := `{"completed": true}`
		req := httptest.NewRequest("PUT", "/api/tasks/1/update", bytes.NewBufferString(updateBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updatedTask models.Task
		json.NewDecoder(resp.Body).Decode(&updatedTask)
		assert.True(t, updatedTask.Completed)

		assert.True(t, dataStore.Schedules["1"].Tasks[0].Completed)
	})

	t.Run("Update Task - Not Completed with Reason", func(t *testing.T) {
		updateBody := `{"completed": false, "notCompletedReason": "Client refused"}`
		req := httptest.NewRequest("PUT", "/api/tasks/2/update", bytes.NewBufferString(updateBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updatedTask models.Task
		json.NewDecoder(resp.Body).Decode(&updatedTask)
		assert.False(t, updatedTask.Completed)
		assert.Equal(t, "Client refused", updatedTask.NotCompletedReason)

		assert.Equal(t, "Client refused", dataStore.Schedules["1"].Tasks[1].NotCompletedReason)
	})

	t.Run("Update Task - Not Found", func(t *testing.T) {
		updateBody := `{"completed": true}`
		req := httptest.NewRequest("PUT", "/api/tasks/999/update", bytes.NewBufferString(updateBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Add Task to Schedule - Success", func(t *testing.T) {
		assert.Len(t, dataStore.Schedules["2"].Tasks, 2)

		taskBody := `{"name": "New Task", "description": "A new test task"}`
		req := httptest.NewRequest("POST", "/api/schedules/2/tasks", bytes.NewBufferString(taskBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Len(t, dataStore.Schedules["2"].Tasks, 3)
		lastTask := dataStore.Schedules["2"].Tasks[2]
		assert.Equal(t, "New Task", lastTask.Name)
		assert.Equal(t, "A new test task", lastTask.Description)
	})
}

func TestVisitFlow(t *testing.T) {
	app, dataStore := setupTest()

	t.Run("Start and End Visit Flow", func(t *testing.T) {
		assert.Equal(t, "scheduled", dataStore.Schedules["2"].Status)

		startBody := `{"location": {"latitude": 10.0, "longitude": 20.0}}`
		startReq := httptest.NewRequest("POST", "/api/schedules/2/start", bytes.NewBufferString(startBody))
		startReq.Header.Set("Content-Type", "application/json")
		startResp, _ := app.Test(startReq)
		assert.Equal(t, http.StatusOK, startResp.StatusCode)

		schedule := dataStore.Schedules["2"]
		assert.Equal(t, "in_progress", schedule.Status)
		assert.NotNil(t, schedule.ClockInTime)
		assert.Equal(t, 10.0, schedule.ClockInLocation.Latitude)

		endBody := `{"location": {"latitude": 10.1, "longitude": 20.1}}`
		endReq := httptest.NewRequest("POST", "/api/schedules/2/end", bytes.NewBufferString(endBody))
		endReq.Header.Set("Content-Type", "application/json")
		endResp, _ := app.Test(endReq)
		assert.Equal(t, http.StatusOK, endResp.StatusCode)

		assert.Equal(t, "completed", schedule.Status)
		assert.NotNil(t, schedule.ClockOutTime)
		assert.Equal(t, 10.1, schedule.ClockOutLocation.Latitude)
	})

	t.Run("Clock-in and Cancel Flow", func(t *testing.T) {
		assert.Equal(t, "scheduled", dataStore.Schedules["4"].Status)

		clockInReq := httptest.NewRequest("GET", "/api/schedules/4/clock-in", nil)
		clockInResp, _ := app.Test(clockInReq)
		assert.Equal(t, http.StatusOK, clockInResp.StatusCode)

		schedule := dataStore.Schedules["4"]
		assert.Equal(t, "in_progress", schedule.Status)
		assert.NotNil(t, schedule.ClockInTime)

		cancelReq := httptest.NewRequest("POST", "/api/schedules/4/cancel-clock-in", nil)
		cancelResp, _ := app.Test(cancelReq)
		assert.Equal(t, http.StatusOK, cancelResp.StatusCode)

		assert.Equal(t, "scheduled", schedule.Status)
		assert.Nil(t, schedule.ClockInTime)
		assert.Nil(t, schedule.ClockInLocation)
	})
}

func TestAdminHandlers(t *testing.T) {
	app, dataStore := setupTest()

	t.Run("Reset Store", func(t *testing.T) {
		originalStatus := dataStore.Schedules["1"].Status
		dataStore.Schedules["1"].Status = "testing_change"
		assert.NotEqual(t, originalStatus, dataStore.Schedules["1"].Status)

		req := httptest.NewRequest("POST", "/api/reset", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, originalStatus, dataStore.Schedules["1"].Status)
	})
}
