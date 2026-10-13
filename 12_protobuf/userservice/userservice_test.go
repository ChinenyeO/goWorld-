package userservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"

	userv1 "example/goWorld/12_protobuf/gen/user/v1"
)

// fakeUserRepository is a test double for UserRepository. Tests configure
// what it returns and inspect what it received.
type fakeUserRepository struct {
	// Configured by the test.
	id  string
	err error

	// Recorded by CreateUser.
	calls    int
	gotName  string
	gotEmail string

	// Configured by the test for GetUser.
	user   User
	getErr error

	// Recorded by GetUser.
	getCalls int
	gotID    string
}

var _ UserRepository = (*fakeUserRepository)(nil)

func (f *fakeUserRepository) CreateUser(
	ctx context.Context,
	name string,
	email string,
) (string, error) {
	f.calls++
	f.gotName = name
	f.gotEmail = email
	return f.id, f.err
}

func (f *fakeUserRepository) GetUser(ctx context.Context, id string) (User, error) {
	f.getCalls++
	f.gotID = id
	return f.user, f.getErr
}

func TestCreateUser_InvalidInput(t *testing.T) {
	tests := []struct {
		name        string
		reqName     string
		reqEmail    string
		wantMessage string
	}{
		{"empty name", "", "chinenye@example.com", "name is required"},
		{"whitespace-only name", "   ", "chinenye@example.com", "name is required"},
		{"empty email", "Chinenye", "", "email is required"},
		{"whitespace-only email", "Chinenye", "   ", "email is required"},
		{"both empty", "", "", "name is required"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &fakeUserRepository{id: "test-id-1"}
			service := NewUserService(repo)

			req := connect.NewRequest(&userv1.CreateUserRequest{
				Name:  test.reqName,
				Email: test.reqEmail,
			})

			resp, err := service.CreateUser(context.Background(), req)

			if resp != nil {
				t.Errorf("expected nil response, got %v", resp)
			}

			if got := connect.CodeOf(err); got != connect.CodeInvalidArgument {
				t.Errorf("expected code %v, got %v (err: %v)", connect.CodeInvalidArgument, got, err)
			}

			var connectErr *connect.Error
			if !errors.As(err, &connectErr) {
				t.Fatalf("expected a *connect.Error, got %T (%v)", err, err)
			}
			if connectErr.Message() != test.wantMessage {
				t.Errorf("expected message %q, got %q", test.wantMessage, connectErr.Message())
			}

			if repo.calls != 0 {
				t.Errorf("expected repository not to be called, got %d calls", repo.calls)
			}
		})
	}
}

func TestCreateUser_Valid(t *testing.T) {
	repo := &fakeUserRepository{id: "test-id-1"}
	service := NewUserService(repo)

	req := connect.NewRequest(&userv1.CreateUserRequest{
		Name:  "Chinenye",
		Email: "chinenye@example.com",
	})

	resp, err := service.CreateUser(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.calls != 1 {
		t.Errorf("expected repository to be called once, got %d calls", repo.calls)
	}

	if repo.gotName != "Chinenye" {
		t.Errorf("expected repository to receive name %q, got %q", "Chinenye", repo.gotName)
	}

	if repo.gotEmail != "chinenye@example.com" {
		t.Errorf("expected repository to receive email %q, got %q", "chinenye@example.com", repo.gotEmail)
	}

	if resp.Msg.GetId() != "test-id-1" {
		t.Errorf("expected id %q, got %q", "test-id-1", resp.Msg.GetId())
	}
}

func TestCreateUser_UnexpectedRepositoryError(t *testing.T) {
	for _, test := range unexpectedErrorCases() {
		t.Run(test.name, func(t *testing.T) {
			logs := captureLog(t)
			repo := &fakeUserRepository{err: test.repoErr}
			service := NewUserService(repo)

			req := connect.NewRequest(&userv1.CreateUserRequest{
				Name:  "Chinenye",
				Email: "chinenye@example.com",
			})

			resp, err := service.CreateUser(context.Background(), req)

			if repo.calls != 1 {
				t.Errorf("expected repository to be called once, got %d calls", repo.calls)
			}

			if resp != nil {
				t.Errorf("expected nil response, got %v", resp)
			}

			test.check(t, err, logs.String())
		})
	}
}

func TestCreateUser_AlreadyExists(t *testing.T) {
	// Wrapped on purpose: the service must use errors.Is, not ==.
	repo := &fakeUserRepository{err: fmt.Errorf("insert user: %w", ErrUserAlreadyExists)}
	service := NewUserService(repo)

	req := connect.NewRequest(&userv1.CreateUserRequest{
		Name:  "Chinenye",
		Email: "chinenye@example.com",
	})

	resp, err := service.CreateUser(context.Background(), req)

	if repo.calls != 1 {
		t.Errorf("expected repository to be called once, got %d calls", repo.calls)
	}

	if resp != nil {
		t.Errorf("expected nil response, got %v", resp)
	}

	if got := connect.CodeOf(err); got != connect.CodeAlreadyExists {
		t.Errorf("expected code %v, got %v (err: %v)", connect.CodeAlreadyExists, got, err)
	}

	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected a *connect.Error, got %T (%v)", err, err)
	}
	if connectErr.Message() != "user already exists" {
		t.Errorf("expected message %q, got %q", "user already exists", connectErr.Message())
	}
}

// requireConnectError checks that err is a *connect.Error with the given code
// and message.
func requireConnectError(t *testing.T, err error, wantCode connect.Code, wantMessage string) {
	t.Helper()

	if got := connect.CodeOf(err); got != wantCode {
		t.Errorf("expected code %v, got %v (err: %v)", wantCode, got, err)
	}

	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected a *connect.Error, got %T (%v)", err, err)
	}
	if connectErr.Message() != wantMessage {
		t.Errorf("expected message %q, got %q", wantMessage, connectErr.Message())
	}
}

func TestGetUser_InvalidInput(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{"empty id", ""},
		{"whitespace-only id", "   "},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &fakeUserRepository{}
			service := NewUserService(repo)

			req := connect.NewRequest(&userv1.GetUserRequest{Id: test.id})

			resp, err := service.GetUser(context.Background(), req)

			if resp != nil {
				t.Errorf("expected nil response, got %v", resp)
			}

			requireConnectError(t, err, connect.CodeInvalidArgument, "id is required")

			if repo.getCalls != 0 {
				t.Errorf("expected repository not to be called, got %d calls", repo.getCalls)
			}
		})
	}
}

func TestGetUser_Found(t *testing.T) {
	createdAt := time.Date(2026, 10, 7, 9, 30, 0, 0, time.UTC)
	repo := &fakeUserRepository{
		user: User{
			ID:        "0b6f8f9e-5d4c-4b1a-9f3e-2a1b0c9d8e7f",
			Name:      "Chinenye",
			Email:     "chinenye@example.com",
			CreatedAt: createdAt,
		},
	}
	service := NewUserService(repo)

	req := connect.NewRequest(&userv1.GetUserRequest{Id: "0b6f8f9e-5d4c-4b1a-9f3e-2a1b0c9d8e7f"})

	resp, err := service.GetUser(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.getCalls != 1 {
		t.Errorf("expected repository to be called once, got %d calls", repo.getCalls)
	}

	if repo.gotID != "0b6f8f9e-5d4c-4b1a-9f3e-2a1b0c9d8e7f" {
		t.Errorf("expected repository to receive id %q, got %q", "0b6f8f9e-5d4c-4b1a-9f3e-2a1b0c9d8e7f", repo.gotID)
	}

	user := resp.Msg.GetUser()

	if user.GetId() != "0b6f8f9e-5d4c-4b1a-9f3e-2a1b0c9d8e7f" {
		t.Errorf("expected id %q, got %q", "0b6f8f9e-5d4c-4b1a-9f3e-2a1b0c9d8e7f", user.GetId())
	}

	if user.GetName() != "Chinenye" {
		t.Errorf("expected name %q, got %q", "Chinenye", user.GetName())
	}

	if user.GetEmail() != "chinenye@example.com" {
		t.Errorf("expected email %q, got %q", "chinenye@example.com", user.GetEmail())
	}

	if got := user.GetCreatedAt().AsTime(); !got.Equal(createdAt) {
		t.Errorf("expected created_at %v, got %v", createdAt, got)
	}
}

func TestGetUser_InvalidID(t *testing.T) {
	// Wrapped on purpose: the service must use errors.Is, not ==.
	repo := &fakeUserRepository{getErr: fmt.Errorf("get user: %w", ErrInvalidUserID)}
	service := NewUserService(repo)

	req := connect.NewRequest(&userv1.GetUserRequest{Id: "not-a-uuid"})

	resp, err := service.GetUser(context.Background(), req)

	if repo.getCalls != 1 {
		t.Errorf("expected repository to be called once, got %d calls", repo.getCalls)
	}

	if resp != nil {
		t.Errorf("expected nil response, got %v", resp)
	}

	requireConnectError(t, err, connect.CodeInvalidArgument, "invalid user id")
}

func TestGetUser_NotFound(t *testing.T) {
	// Wrapped on purpose: the service must use errors.Is, not ==.
	repo := &fakeUserRepository{getErr: fmt.Errorf("get user: %w", ErrUserNotFound)}
	service := NewUserService(repo)

	req := connect.NewRequest(&userv1.GetUserRequest{Id: "0b6f8f9e-5d4c-4b1a-9f3e-2a1b0c9d8e7f"})

	resp, err := service.GetUser(context.Background(), req)

	if repo.getCalls != 1 {
		t.Errorf("expected repository to be called once, got %d calls", repo.getCalls)
	}

	if resp != nil {
		t.Errorf("expected nil response, got %v", resp)
	}

	requireConnectError(t, err, connect.CodeNotFound, "user not found")
}

func TestGetUser_UnexpectedRepositoryError(t *testing.T) {
	for _, test := range unexpectedErrorCases() {
		t.Run(test.name, func(t *testing.T) {
			logs := captureLog(t)
			repo := &fakeUserRepository{getErr: test.repoErr}
			service := NewUserService(repo)

			req := connect.NewRequest(&userv1.GetUserRequest{Id: "0b6f8f9e-5d4c-4b1a-9f3e-2a1b0c9d8e7f"})

			resp, err := service.GetUser(context.Background(), req)

			if repo.getCalls != 1 {
				t.Errorf("expected repository to be called once, got %d calls", repo.getCalls)
			}

			if resp != nil {
				t.Errorf("expected nil response, got %v", resp)
			}

			test.check(t, err, logs.String())
		})
	}
}

// dbErrorDetail mimics what a real PostgreSQL error message exposes.
const dbErrorDetail = `ERROR: relation "users" does not exist (SQLSTATE 42P01)`

type unexpectedErrorCase struct {
	name    string
	repoErr error
	check   func(t *testing.T, err error, logs string)
}

// unexpectedErrorCases are repository errors with no domain meaning, shared by
// the CreateUser and GetUser tests.
func unexpectedErrorCases() []unexpectedErrorCase {
	return []unexpectedErrorCase{
		{
			name:    "database error becomes internal",
			repoErr: fmt.Errorf("insert user: %w", errors.New(dbErrorDetail)),
			check: func(t *testing.T, err error, logs string) {
				requireConnectError(t, err, connect.CodeInternal, "internal error")

				// Nothing from the underlying error may reach the client.
				if strings.Contains(err.Error(), "SQLSTATE") || strings.Contains(err.Error(), "relation") {
					t.Errorf("expected database details to be hidden from client, got %q", err.Error())
				}

				// ...but it must be logged for debugging.
				if !strings.Contains(logs, dbErrorDetail) {
					t.Errorf("expected underlying error to be logged, got logs %q", logs)
				}
			},
		},
		{
			name:    "canceled context stays canceled",
			repoErr: fmt.Errorf("query: %w", context.Canceled),
			check: func(t *testing.T, err error, logs string) {
				requireConnectError(t, err, connect.CodeCanceled, "request canceled")
				if logs != "" {
					t.Errorf("expected cancellation not to be logged, got %q", logs)
				}
			},
		},
		{
			name:    "deadline stays deadline exceeded",
			repoErr: fmt.Errorf("query: %w", context.DeadlineExceeded),
			check: func(t *testing.T, err error, logs string) {
				requireConnectError(t, err, connect.CodeDeadlineExceeded, "request deadline exceeded")
				if logs != "" {
					t.Errorf("expected deadline not to be logged, got %q", logs)
				}
			},
		},
	}
}

// captureLog redirects the standard logger for the duration of the test.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()

	var buf bytes.Buffer
	original := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(original) })

	return &buf
}
