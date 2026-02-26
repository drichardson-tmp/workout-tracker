// Package roles defines Zitadel role constants shared across the backend.
// Keeping them in a leaf package avoids import cycles between middleware and handlers.
package roles

const (
	Admin = "admin"
	User  = "user"
)
