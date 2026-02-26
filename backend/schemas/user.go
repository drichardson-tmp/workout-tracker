package schemas

import "time"

// --- inputs ---

type ListUsersInput struct {
	Email string `query:"email" doc:"Filter by email address"`
}

type CreateUserInput struct {
	Body struct {
		Email    string `json:"email" format:"email" doc:"User email address"`
		Name     string `json:"name" minLength:"1" doc:"Display name"`
		Password string `json:"password" minLength:"8" doc:"Password (min 8 characters)"`
	}
}

type GetUserInput struct {
	UserID int64 `path:"userId" doc:"User ID"`
}

// --- outputs / response bodies ---

type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetUserOutput struct {
	Body *UserResponse
}

type CreateUserOutput struct {
	Body *UserResponse
}

type ListUsersOutput struct {
	Body []UserResponse
}
