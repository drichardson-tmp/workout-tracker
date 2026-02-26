package models

type User struct {
	BaseModel
	ZitadelID    *string `gorm:"uniqueIndex"`           // nil for manually-created users
	Email        string  `gorm:"uniqueIndex;not null"`
	Name         string  `gorm:"not null"`
	PasswordHash string  // empty for Zitadel-provisioned users
}
