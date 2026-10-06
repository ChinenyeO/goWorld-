package userservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// pgUniqueViolation is the PostgreSQL SQLSTATE for unique_violation.
	pgUniqueViolation = "23505"

	// usersEmailUniqueConstraint must match the constraint name in
	// db/migrations/20261005000000_create_users.sql.
	usersEmailUniqueConstraint = "users_email_key"
)

// PostgresUserRepository is a UserRepository backed by the PostgreSQL
// users table.
type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

// Compile-time check: *PostgresUserRepository implements UserRepository.
var _ UserRepository = (*PostgresUserRepository)(nil)

// NewPostgresUserRepository returns a repository that runs its queries on
// the given connection pool. The caller owns the pool and closes it.
func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		pool: pool,
	}
}

// CreateUser inserts a user and returns the UUID PostgreSQL generated for it.
func (r *PostgresUserRepository) CreateUser(
	ctx context.Context,
	name string,
	email string,
) (string, error) {
	const query = `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id`

	var id pgtype.UUID
	if err := r.pool.QueryRow(ctx, query, name, email).Scan(&id); err != nil {
		return "", translateInsertUserError(err)
	}

	return id.String(), nil
}

// translateInsertUserError converts PostgreSQL-specific insert failures into
// repository errors. A duplicate email becomes ErrUserAlreadyExists; any
// other error is wrapped unchanged so it remains inspectable.
func translateInsertUserError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) &&
		pgErr.Code == pgUniqueViolation &&
		pgErr.ConstraintName == usersEmailUniqueConstraint {
		return fmt.Errorf("insert user: %w", ErrUserAlreadyExists)
	}

	return fmt.Errorf("insert user: %w", err)
}
