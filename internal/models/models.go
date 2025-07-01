package models

import "time"

type Geolocation struct {
	Latitude  float64 `json:"latitude" example:"-6.200000"`
	Longitude float64 `json:"longitude" example:"106.816666"`
}

type Task struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" example:"Give medication"`
	Description string `json:"description" example:"Administer morning pills with water."`

	Completed          bool   `json:"completed"`
	NotCompletedReason string `json:"notCompletedReason,omitempty" example:"Client refused medication."`
}

type ClientContact struct {
	Email string `json:"email" example:"melisa@example.com"`
	Phone string `json:"phone" example:"+44 1232 212 3233"`
}

type Location struct {
	Address     string      `json:"address" example:"123 Main St"`
	Coordinates Geolocation `json:"coordinates"`
}

type Schedule struct {
	ID            string        `json:"id" example:"1"`
	ClientName    string        `json:"clientName" example:"Melisa Adam"`
	ServiceName   string        `json:"serviceName" example:"Casa Grande Apartment"`
	ShiftDate     string        `json:"shiftDate" example:"2025-01-15"`
	ShiftTime     string        `json:"shiftTime" example:"09:00 - 10:00"`
	AmOrPm        string        `json:"amOrPm" example:"AM"`        // "AM" or "PM"
	Status        string        `json:"status" example:"scheduled"` // "scheduled", "in-progress", "completed", "missed", "cancelled"
	Tasks         []Task        `json:"tasks"`
	ClientContact ClientContact `json:"clientContact"`
	ServiceNotes  string        `json:"serviceNotes,omitempty" example:"Client may be a bit groggy."`

	ClockInTime      *time.Time   `json:"clockInTime,omitempty"`
	ClockOutTime     *time.Time   `json:"clockOutTime,omitempty"`
	ClockInLocation  *Geolocation `json:"clockInLocation,omitempty"`
	ClockOutLocation *Geolocation `json:"clockOutLocation,omitempty"`
	Location         Location     `json:"location"`
}

type StartVisitRequest struct {
	Timestamp string      `json:"timestamp"`
	Location  Geolocation `json:"location"`
}

type EndVisitRequest struct {
	Timestamp string      `json:"timestamp"`
	Location  Geolocation `json:"location"`
}

type UpdateTaskRequest struct {
	Completed          bool   `json:"completed"`
	NotCompletedReason string `json:"notCompletedReason,omitempty"`
}

type AddTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
