// Package userservice contains your implementation of the UserService
// defined in proto/user/v1/user.proto.
//
// This is YOUR code. The generated code lives in gen/ and must not be edited.
package userservice

import (
	"context"
	"errors"
	"log"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "example/goWorld/12_protobuf/gen/user/v1"
	"example/goWorld/12_protobuf/gen/user/v1/userv1connect"
)

// ---------------------------------------------------------------------------
// DEPENDENCY INJECTION HINT:
//
//	type OrderService struct {
//	    store OrderStore
//	}
//
//	func NewOrderService(store OrderStore) *OrderService {
//	    return &OrderService{
//	        store: store,
//	    }
//	}
// ---------------------------------------------------------------------------

// UserService implements userv1connect.UserServiceHandler.
type UserService struct {
	// TODO 1: Add a field named repo with type UserRepository.
	repo UserRepository
}

// TODO 2: Create a constructor named NewUserService that accepts a
//         UserRepository.

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}

}

//
// TODO 3: Have NewUserService return a *UserService initialized with that
//         repository.

// Compile-time check: this line fails to compile until *UserService has
// every method the generated UserServiceHandler interface requires, with
// exactly the right signatures. If the compiler complains here, compare your
// method with the interface in gen/user/v1/userv1connect/user.connect.go.
var _ userv1connect.UserServiceHandler = (*UserService)(nil)

// ---------------------------------------------------------------------------
// HINT 1:
// A method can look like:
//
//   func (s *ProductService) CreateProduct(
//       ctx context.Context,
//       req *connect.Request[productv1.CreateProductRequest],
//   ) (*connect.Response[productv1.CreateProductResponse], error) {
//       // implementation
//   }
//
// HINT 2:
// The protobuf message inside a Connect request is available through:
//
//   req.Msg
//
// Generated protobuf getters can then be called on that message.
//
// HINT 3:
// A protobuf response object can be created with a struct literal, for example:
//
//   product := &productv1.CreateProductResponse{
//       Id: "product-123",
//   }
//
// HINT 4:
// ConnectRPC provides connect.NewResponse(...) for wrapping a protobuf response.
// ---------------------------------------------------------------------------

// TODO: Implement CreateUser below.
//
// It must match the generated interface:
//
//   CreateUser(context.Context, *connect.Request[userv1.CreateUserRequest]) (*connect.Response[userv1.CreateUserResponse], error)
//
// In this file, the generated message package is imported as "userv1".

func (s *UserService) CreateUser(
	ctx context.Context,
	req *connect.Request[userv1.CreateUserRequest],
) (*connect.Response[userv1.CreateUserResponse], error) {
	if strings.TrimSpace(req.Msg.GetName()) == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name is required"))
	}

	if strings.TrimSpace(req.Msg.GetEmail()) == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email is required"))
	}

	id, err := s.repo.CreateUser(
		ctx,
		req.Msg.GetName(),
		req.Msg.GetEmail(),
	)

	if errors.Is(err, ErrUserAlreadyExists) {
		return nil, connect.NewError(connect.CodeAlreadyExists, ErrUserAlreadyExists)
	}

	if err != nil {
		return nil, unexpectedError("CreateUser", err)
	}

	response := &userv1.CreateUserResponse{
		Id: id,
	}

	return connect.NewResponse(response), nil

}

// GetUser returns the user with the requested ID.
func (s *UserService) GetUser(
	ctx context.Context,
	req *connect.Request[userv1.GetUserRequest],
) (*connect.Response[userv1.GetUserResponse], error) {
	if strings.TrimSpace(req.Msg.GetId()) == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is required"))
	}

	user, err := s.repo.GetUser(ctx, req.Msg.GetId())

	if errors.Is(err, ErrInvalidUserID) {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrInvalidUserID)
	}

	if errors.Is(err, ErrUserNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, ErrUserNotFound)
	}

	if err != nil {
		return nil, unexpectedError("GetUser", err)
	}

	response := &userv1.GetUserResponse{
		User: toProtoUser(user),
	}

	return connect.NewResponse(response), nil
}

// errInternal is the only detail clients receive for unexpected failures.
var errInternal = errors.New("internal error")

// unexpectedError converts a repository error that has no domain meaning into
// a safe Connect error. The underlying error can contain database details
// (SQL state, constraint or table names, connection info), so it is logged
// server-side and never returned to the client.
//
// Cancellations and deadlines are the caller's doing, not server faults, so
// they keep their own codes and are not logged.
func unexpectedError(procedure string, err error) error {
	switch {
	case errors.Is(err, context.Canceled):
		return connect.NewError(connect.CodeCanceled, errors.New("request canceled"))
	case errors.Is(err, context.DeadlineExceeded):
		return connect.NewError(connect.CodeDeadlineExceeded, errors.New("request deadline exceeded"))
	}

	log.Printf("userservice: %s: unexpected repository error: %v", procedure, err)
	return connect.NewError(connect.CodeInternal, errInternal)
}

// toProtoUser converts the domain User into its protobuf API representation.
func toProtoUser(user User) *userv1.User {
	return &userv1.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
