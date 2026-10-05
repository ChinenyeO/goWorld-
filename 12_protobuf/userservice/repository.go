package userservice

import (
	"context"
)

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
}
