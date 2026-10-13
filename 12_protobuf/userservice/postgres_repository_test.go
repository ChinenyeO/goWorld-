package userservice

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestTranslateInsertUserError(t *testing.T) {
	emailTaken := &pgconn.PgError{Code: "23505", ConstraintName: "users_email_key"}
	otherUnique := &pgconn.PgError{Code: "23505", ConstraintName: "users_pkey"}
	notNull := &pgconn.PgError{Code: "23502", ColumnName: "name"}

	tests := []struct {
		name            string
		err             error
		wantExists      bool
		wantPreservedAs error // checked with errors.Is when wantExists is false
	}{
		{"email unique violation", emailTaken, true, nil},
		{"wrapped email unique violation", fmt.Errorf("x: %w", emailTaken), true, nil},
		{"other unique constraint", otherUnique, false, otherUnique},
		{"other postgres error", notNull, false, notNull},
		{"non-postgres error", context.Canceled, false, context.Canceled},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := translateInsertUserError(test.err)

			if errors.Is(got, ErrUserAlreadyExists) != test.wantExists {
				t.Fatalf("errors.Is(%v, ErrUserAlreadyExists) = %v, want %v",
					got, !test.wantExists, test.wantExists)
			}

			if !test.wantExists && !errors.Is(got, test.wantPreservedAs) {
				t.Errorf("expected original error %v to be preserved in %v", test.wantPreservedAs, got)
			}
		})
	}
}

func TestTranslateGetUserError(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "57014"} // query_canceled

	tests := []struct {
		name            string
		err             error
		wantNotFound    bool
		wantPreservedAs error // checked with errors.Is when wantNotFound is false
	}{
		{"no rows", pgx.ErrNoRows, true, nil},
		{"wrapped no rows", fmt.Errorf("x: %w", pgx.ErrNoRows), true, nil},
		{"postgres error", pgErr, false, pgErr},
		{"non-postgres error", context.Canceled, false, context.Canceled},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := translateGetUserError(test.err)

			if errors.Is(got, ErrUserNotFound) != test.wantNotFound {
				t.Fatalf("errors.Is(%v, ErrUserNotFound) = %v, want %v",
					got, !test.wantNotFound, test.wantNotFound)
			}

			if test.wantNotFound && errors.Is(got, pgx.ErrNoRows) {
				t.Errorf("expected pgx.ErrNoRows not to leak through %v", got)
			}

			if !test.wantNotFound && !errors.Is(got, test.wantPreservedAs) {
				t.Errorf("expected original error %v to be preserved in %v", test.wantPreservedAs, got)
			}
		})
	}
}

// Malformed IDs are rejected before the pool is used, so a repository without
// a database is enough for this test.
func TestPostgresGetUser_MalformedID(t *testing.T) {
	repo := &PostgresUserRepository{}

	tests := []struct {
		name string
		id   string
	}{
		{"not a uuid", "not-a-uuid"},
		{"too short", "1234"},
		{"right length, not hex", "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := repo.GetUser(context.Background(), test.id)

			if !errors.Is(err, ErrInvalidUserID) {
				t.Errorf("expected ErrInvalidUserID, got %v", err)
			}
		})
	}
}
