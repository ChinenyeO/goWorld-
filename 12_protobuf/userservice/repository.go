package userservice

import (
	"context"
	"errors"
)

// ErrUserAlreadyExists is returned by UserRepository implementations when a
// user with the given email already exists. Callers should check for it with
// errors.Is, since implementations may wrap it.
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrUserNotFound is returned by UserRepository implementations when no user
// has the given ID.
var ErrUserNotFound = errors.New("user not found")

// ErrInvalidUserID is returned by UserRepository implementations when an ID is
// not in a format the repository can identify a user by.
var ErrInvalidUserID = errors.New("invalid user id")

// ---------------------------------------------------------------------------
// HINT:
//
//	type ProductRepository interface {
//	    SaveProduct(ctx context.Context, name string, price int) (string, error)
//	}
// ---------------------------------------------------------------------------

// UserRepository describes how users are stored. UserService will depend on
// this interface rather than on a specific database.
type UserRepository interface {
	// TODO: Declare a CreateUser method that:
	//   1. accepts a context.Context
	//   2. accepts name and email strings
	//   3. returns a string user ID and an error
	CreateUser(ctx context.Context, name string, email string) (string, error)

	// GetUser returns the user with the given ID, ErrUserNotFound if there is
	// none, or ErrInvalidUserID if id is malformed.
	GetUser(ctx context.Context, id string) (User, error)
}
