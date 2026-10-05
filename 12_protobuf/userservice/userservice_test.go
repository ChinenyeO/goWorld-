package userservice

import (
	"context"
	"errors"
	"testing"

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

func TestCreateUser_RepositoryError(t *testing.T) {
	sentinelErr := errors.New("storage down")
	repo := &fakeUserRepository{err: sentinelErr}
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

	if !errors.Is(err, sentinelErr) {
		t.Errorf("expected error to wrap %v, got %v", sentinelErr, err)
	}
}
