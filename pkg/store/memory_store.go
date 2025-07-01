package store

import (
	"log"
	"sync"
	"time"

	"github.com/IkoAfianando/mini_evv_logger_go/pkg/models"
)

type Store struct {
	mu        sync.Mutex
	Schedules map[string]*models.Schedule
	Tasks     map[int]*models.Task
}

func NewStore() *Store {
	return &Store{
		Schedules: make(map[string]*models.Schedule),
		Tasks:     make(map[int]*models.Task),
	}
}

func (s *Store) SetupInitialData() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Schedules = make(map[string]*models.Schedule)
	s.Tasks = make(map[int]*models.Task)

	initialSchedules := []*models.Schedule{
		{
			ID:          "1",
			ClientName:  "Melisa Adam",
			ServiceName: "Casa Grande Apartment",
			ShiftDate:   time.Now().Format("2006-01-02"),
			ShiftTime:   "00:00 - 6:00",
			AmOrPm:      "AM",
			Status:      "scheduled",
			Tasks: []models.Task{
				{ID: 1, Name: "Give medication", Description: "Administer morning pills with water."},
				{ID: 2, Name: "Assist with bathing", Description: "Ensure safety during shower."},
			},
			ClientContact: models.ClientContact{Email: "melisa@example.com", Phone: "+44 1232 212 3233"},
			ServiceNotes:  "Client may be a bit groggy in the morning. Speak clearly and be patient.",
			Location: models.Location{
				Address: "123 Main St, Springfield, IL",
				Coordinates: models.Geolocation{
					Latitude:  40.712776,
					Longitude: -74.005974,
				},
			},
		},
		{
			ID:          "2",
			ClientName:  "John Doe",
			ServiceName: "Senior Living Center",
			ShiftDate:   time.Now().Format("2006-01-02"),
			ShiftTime:   "06:00 - 12:00",
			AmOrPm:      "AM",
			Status:      "scheduled",
			Tasks: []models.Task{
				{ID: 3, Name: "Prepare lunch", Description: "Low-sodium, soft food diet."},
				{ID: 4, Name: "Light housekeeping", Description: "Tidy up living room and kitchen."},
			},
			ClientContact: models.ClientContact{Email: "john.doe@example.com", Phone: "+1 555 123 4567"},
			ServiceNotes:  "John enjoys listening to classical music during his lunch.",
			Location: models.Location{
				Address: "456 Oak Ave, Springfield, IL",
				Coordinates: models.Geolocation{
					Latitude:  40.712776,
					Longitude: -74.005974,
				},
			},
		},
		{
			ID:          "3",
			ClientName:  "Jane Smith",
			ServiceName: "Private Residence",
			ShiftDate:   time.Now().Format("2006-01-02"),
			ShiftTime:   "2:00 - 3:00",
			AmOrPm:      "AM",
			Status:      "completed",
			Tasks: []models.Task{
				{ID: 5, Name: "Physical therapy exercises", Description: "Follow the chart from Dr. Evans.", Completed: false, NotCompletedReason: "Client was too tired."},
				{ID: 6, Name: "Check vitals", Description: "Measure blood pressure and heart rate.", Completed: true},
			},
			ClientContact: models.ClientContact{Email: "jane.s@example.com", Phone: "+1 555 987 6543"},
			ServiceNotes:  "Client was in good spirits and completed all exercises without issue.",
			Location: models.Location{
				Address: "789 Pine Rd, Springfield, IL",
				Coordinates: models.Geolocation{
					Latitude:  40.712776,
					Longitude: -74.005974,
				},
			},
		},
		{
			ID:          "4",
			ClientName:  "Alice Johnson",
			ServiceName: "Community Health Center",
			ShiftDate:   time.Now().Format("2006-01-02"),
			ShiftTime:   "00:00 - 06:00",
			AmOrPm:      "PM",
			Status:      "scheduled",
			Tasks: []models.Task{
				{ID: 7, Name: "Administer insulin", Description: "Check blood sugar before administering."},
				{ID: 8, Name: "Assist with mobility", Description: "Help client walk to the therapy room."},
			},
			ClientContact: models.ClientContact{Email: "alice@example.com", Phone: "+1 555 321 6543"},
			ServiceNotes:  "Alice is diabetic and requires regular monitoring. Ensure she has her glucose meter.",
			Location: models.Location{
				Address: "321 Maple St, Springfield, IL",
				Coordinates: models.Geolocation{
					Latitude:  40.712776,
					Longitude: -74.005974,
				},
			},
		},
		{
			ID:          "5",
			ClientName:  "Bob Brown",
			ServiceName: "Assisted Living Facility",
			ShiftDate:   time.Now().Format("2006-01-02"),
			ShiftTime:   "06:00 - 11:59",
			AmOrPm:      "PM",
			Status:      "scheduled",
			Tasks: []models.Task{
				{ID: 9, Name: "Monitor heart rate", Description: "Use the portable ECG machine."},
				{ID: 10, Name: "Provide companionship", Description: "Spend time reading and chatting."},
			},
			ClientContact: models.ClientContact{Email: "bob@example.com", Phone: "+1 555 456 7890"},
			ServiceNotes:  "Bob enjoys reading mystery novels. Bring a book to read together.",
			Location: models.Location{
				Address: "654 Cedar Blvd, Springfield, IL",
				Coordinates: models.Geolocation{
					Latitude:  40.712776,
					Longitude: -74.005974,
				},
			},
		},
		{
			ID:          "6",
			ClientName:  "Charlie Green",
			ServiceName: "Home Care Services",
			ShiftDate:   time.Now().Format("2006-01-02"),
			ShiftTime:   "2:00 - 3:00",
			AmOrPm:      "PM",
			Status:      "missed",
			Tasks: []models.Task{
				{ID: 11, Name: "Check medication schedule", Description: "Ensure all medications are taken as prescribed.", Completed: false, NotCompletedReason: "Client was not home."},
				{ID: 12, Name: "Assist with meal prep", Description: "Prepare a light snack for the client.", Completed: false, NotCompletedReason: "Client refused meal."},
			},
			ClientContact: models.ClientContact{Email: "charlie@example.com", Phone: "+1 555 789 1234"},
			ServiceNotes:  "Charlie was not home during the scheduled visit. Attempted to call but no answer.",
			Location: models.Location{
				Address: "987 Birch St, Springfield, IL",
				Coordinates: models.Geolocation{
					Latitude:  40.712776,
					Longitude: -74.005974,
				},
			},
		},
	}

	for _, schedule := range initialSchedules {
		s.Schedules[schedule.ID] = schedule
		for i := range schedule.Tasks {
			task := schedule.Tasks[i]
			s.Tasks[task.ID] = &task
		}
	}
	log.Println("In-memory data store initialized.")
}
