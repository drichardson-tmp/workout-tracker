package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"workout-tracker/backend/middleware"
	"workout-tracker/backend/models"
	"workout-tracker/backend/roles"
	"workout-tracker/backend/schemas"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) ListUsers(ctx context.Context, input *schemas.ListUsersInput) (*schemas.ListUsersOutput, error) {
	if authCtx := middleware.GetAuth(ctx); authCtx != nil && !authCtx.IsGrantedRole(roles.Admin) {
		return nil, huma.NewError(http.StatusForbidden, "admin role required")
	}

	var users []models.User
	q := h.db
	if input.Email != "" {
		q = q.Where("email = ?", input.Email)
	}
	if err := q.Find(&users).Error; err != nil {
		return nil, huma.Error500InternalServerError("failed to fetch users")
	}
	out := &schemas.ListUsersOutput{Body: make([]schemas.UserResponse, len(users))}
	for i, u := range users {
		out.Body[i] = userToResponse(u)
	}
	return out, nil
}

func (h *UserHandler) GetUser(ctx context.Context, input *schemas.GetUserInput) (*schemas.GetUserOutput, error) {
	var user models.User
	if err := h.db.First(&user, input.UserID).Error; err != nil {
		return nil, huma.NewError(http.StatusNotFound, "user not found")
	}
	r := userToResponse(user)
	return &schemas.GetUserOutput{Body: &r}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, input *schemas.CreateUserInput) (*schemas.CreateUserOutput, error) {
	user := models.User{
		Email:        input.Body.Email,
		Name:         input.Body.Name,
		PasswordHash: input.Body.Password, // TODO: hash before storing
	}
	if err := h.db.Create(&user).Error; err != nil {
		return nil, huma.Error500InternalServerError("failed to create user")
	}
	r := userToResponse(user)
	return &schemas.CreateUserOutput{Body: &r}, nil
}

func userToResponse(u models.User) schemas.UserResponse {
	return schemas.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
