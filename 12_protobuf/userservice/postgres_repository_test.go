package userservice

import (
	"context"
	"errors"
	"fmt"
	"testing"

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
