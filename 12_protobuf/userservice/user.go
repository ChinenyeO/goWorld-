package userservice

import "time"

// User is the application's view of a stored user. It is independent of both
// the protobuf API types and the database driver.
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}
