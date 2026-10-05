package userservice

import (
	"context"
)

// ---------------------------------------------------------------------------
// HINT: a type implementing an interface
//
//	type MemoryProductStore struct{}
//
//	var _ ProductStore = (*MemoryProductStore)(nil)
//
//	func (m *MemoryProductStore) SaveProduct(
//	    ctx context.Context,
//	    name string,
//	) (string, error) {
//	    // TODO
//	}
// ---------------------------------------------------------------------------

// InMemoryUserRepository is a UserRepository that keeps nothing in a real
// database. For now it only needs to hand back a generated user ID.
type InMemoryUserRepository struct{}

// Compile-time check: this fails to compile until *InMemoryUserRepository
// has every method UserRepository requires, with matching signatures.
var _ UserRepository = (*InMemoryUserRepository)(nil)

// TODO: Implement CreateUser on *InMemoryUserRepository.
//
// Its signature must match the CreateUser method in the UserRepository
// interface (see repository.go). It should return a generated user ID and
// an error.

func (m *InMemoryUserRepository) CreateUser(
	ctx context.Context,
	name string,
	email string,
) (string, error) {
	return "user-memory-123", nil
}
