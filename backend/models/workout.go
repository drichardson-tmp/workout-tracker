package models

type Workout struct {
	BaseModel
	UserID          int64 `gorm:"not null;index"`
	User            User
	Name            string `gorm:"not null"`
	Description     string
	DurationMinutes int
}

//
