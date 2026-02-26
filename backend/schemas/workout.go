package schemas

import "time"

// --- inputs ---

type ListWorkoutsInput struct {
	UserID int64 `query:"userId" doc:"Filter workouts by user ID"`
}

type GetWorkoutInput struct {
	WorkoutID int64 `path:"workoutId" doc:"Workout ID"`
}

type CreateWorkoutInput struct {
	Body struct {
		UserID          int64  `json:"user_id,omitempty" doc:"Owner user ID (dev only; derived from auth token in production)"`
		Name            string `json:"name" minLength:"1" doc:"Workout name"`
		Description     string `json:"description,omitempty" doc:"Optional description"`
		DurationMinutes int    `json:"duration_minutes" minimum:"0" doc:"Duration in minutes"`
	}
}

type UpdateWorkoutInput struct {
	WorkoutID int64 `path:"workoutId" doc:"Workout ID"`
	Body      struct {
		Name            string `json:"name,omitempty" doc:"Workout name"`
		Description     string `json:"description,omitempty" doc:"Optional description"`
		DurationMinutes int    `json:"duration_minutes,omitempty" minimum:"0" doc:"Duration in minutes"`
	}
}

type DeleteWorkoutInput struct {
	WorkoutID int64 `path:"workoutId" doc:"Workout ID"`
}

// --- outputs / response bodies ---

type WorkoutResponse struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	DurationMinutes int       `json:"duration_minutes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type GetWorkoutOutput struct {
	Body *WorkoutResponse
}

type CreateWorkoutOutput struct {
	Body *WorkoutResponse
}

type UpdateWorkoutOutput struct {
	Body *WorkoutResponse
}

type ListWorkoutsOutput struct {
	Body []WorkoutResponse
}
