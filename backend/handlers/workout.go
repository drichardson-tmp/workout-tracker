package handlers

import (
	"context"
	"errors"
	"net/http"

	"workout-tracker/backend/middleware"
	"workout-tracker/backend/models"
	"workout-tracker/backend/schemas"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type WorkoutHandler struct {
	db *gorm.DB
}

func NewWorkoutHandler(db *gorm.DB) *WorkoutHandler {
	return &WorkoutHandler{db: db}
}

// resolveUserID returns the local user.ID for the current request.
// When Zitadel auth is active it finds or auto-creates a local User from the
// token claims. In dev/test mode (no auth context) it returns 0 so callers
// can fall back to a user_id supplied in the request body.
func (h *WorkoutHandler) resolveUserID(ctx context.Context) (int64, error) {
	info := middleware.GetUserInfo(ctx)
	if info == nil {
		return 0, nil
	}

	var user models.User
	err := h.db.Where("zitadel_id = ?", info.ZitadelID).First(&user).Error
	if err == nil {
		return user.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	// First request from this Zitadel user — provision a local record.
	email := info.Email
	if email == "" {
		email = info.ZitadelID + "@zitadel.local"
	}
	name := info.Username
	if name == "" {
		name = info.ZitadelID
	}
	user = models.User{ZitadelID: &info.ZitadelID, Email: email, Name: name}
	if err := h.db.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (h *WorkoutHandler) ListWorkouts(ctx context.Context, input *schemas.ListWorkoutsInput) (*schemas.ListWorkoutsOutput, error) {
	var workouts []models.Workout
	q := h.db

	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to resolve user")
	}
	if userID != 0 {
		// Auth active: scope results to the current user.
		q = q.Where("user_id = ?", userID)
	} else if input.UserID != 0 {
		// Dev/test fallback: honour the query-param filter.
		q = q.Where("user_id = ?", input.UserID)
	}

	if err := q.Find(&workouts).Error; err != nil {
		return nil, huma.Error500InternalServerError("failed to fetch workouts")
	}
	out := &schemas.ListWorkoutsOutput{Body: make([]schemas.WorkoutResponse, len(workouts))}
	for i, w := range workouts {
		out.Body[i] = workoutToResponse(w)
	}
	return out, nil
}

func (h *WorkoutHandler) GetWorkout(ctx context.Context, input *schemas.GetWorkoutInput) (*schemas.GetWorkoutOutput, error) {
	var workout models.Workout
	if err := h.db.First(&workout, input.WorkoutID).Error; err != nil {
		return nil, huma.NewError(http.StatusNotFound, "workout not found")
	}
	r := workoutToResponse(workout)
	return &schemas.GetWorkoutOutput{Body: &r}, nil
}

func (h *WorkoutHandler) CreateWorkout(ctx context.Context, input *schemas.CreateWorkoutInput) (*schemas.CreateWorkoutOutput, error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to resolve user")
	}
	if userID == 0 {
		// Dev/test fallback: accept user_id from the request body.
		if input.Body.UserID == 0 {
			return nil, huma.NewError(http.StatusUnauthorized, "authentication required")
		}
		userID = input.Body.UserID
	}

	workout := models.Workout{
		UserID:          userID,
		Name:            input.Body.Name,
		Description:     input.Body.Description,
		DurationMinutes: input.Body.DurationMinutes,
	}
	if err := h.db.Create(&workout).Error; err != nil {
		return nil, huma.Error500InternalServerError("failed to create workout")
	}
	r := workoutToResponse(workout)
	return &schemas.CreateWorkoutOutput{Body: &r}, nil
}

func (h *WorkoutHandler) UpdateWorkout(ctx context.Context, input *schemas.UpdateWorkoutInput) (*schemas.UpdateWorkoutOutput, error) {
	var workout models.Workout
	if err := h.db.First(&workout, input.WorkoutID).Error; err != nil {
		return nil, huma.NewError(http.StatusNotFound, "workout not found")
	}
	if input.Body.Name != "" {
		workout.Name = input.Body.Name
	}
	if input.Body.Description != "" {
		workout.Description = input.Body.Description
	}
	if input.Body.DurationMinutes != 0 {
		workout.DurationMinutes = input.Body.DurationMinutes
	}
	if err := h.db.Save(&workout).Error; err != nil {
		return nil, huma.Error500InternalServerError("failed to update workout")
	}
	r := workoutToResponse(workout)
	return &schemas.UpdateWorkoutOutput{Body: &r}, nil
}

func (h *WorkoutHandler) DeleteWorkout(ctx context.Context, input *schemas.DeleteWorkoutInput) (*struct{}, error) {
	if err := h.db.Delete(&models.Workout{}, input.WorkoutID).Error; err != nil {
		return nil, huma.Error500InternalServerError("failed to delete workout")
	}
	return nil, nil
}

func workoutToResponse(w models.Workout) schemas.WorkoutResponse {
	return schemas.WorkoutResponse{
		ID:              w.ID,
		UserID:          w.UserID,
		Name:            w.Name,
		Description:     w.Description,
		DurationMinutes: w.DurationMinutes,
		CreatedAt:       w.CreatedAt,
		UpdatedAt:       w.UpdatedAt,
	}
}
