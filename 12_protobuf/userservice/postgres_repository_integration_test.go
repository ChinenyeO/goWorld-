package userservice

// Integration tests for PostgresUserRepository against a real PostgreSQL.
//
// They only run when TEST_DATABASE_URL is set, and refuse to run unless it
// names a database ending in "_test" that is not the DATABASE_URL database.
// The schema must already be migrated:
//
//	TEST_DATABASE_URL='postgres://.../goworld_test?sslmode=disable' atlas migrate apply --env test
//	TEST_DATABASE_URL='postgres://.../goworld_test?sslmode=disable' go test ./12_protobuf/userservice -run Integration -v

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// integrationPool returns a pool for the explicitly configured test database,
// or skips the test if none is configured.
func integrationPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping PostgreSQL integration test")
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatalf("parse TEST_DATABASE_URL: %v", err)
	}

	// Safety checks: misconfiguration fails loudly instead of skipping.
	if !strings.HasSuffix(cfg.ConnConfig.Database, "_test") {
		t.Fatalf("refusing to run: TEST_DATABASE_URL database %q must end in _test", cfg.ConnConfig.Database)
	}
	if devURL := os.Getenv("DATABASE_URL"); devURL != "" {
		if devCfg, err := pgxpool.ParseConfig(devURL); err == nil &&
			devCfg.ConnConfig.Host == cfg.ConnConfig.Host &&
			devCfg.ConnConfig.Port == cfg.ConnConfig.Port &&
			devCfg.ConnConfig.Database == cfg.ConnConfig.Database {
			t.Fatal("refusing to run: TEST_DATABASE_URL points at the same database as DATABASE_URL")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("create test pool: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping test database: %v", err)
	}

	var hasUsers bool
	if err := pool.QueryRow(ctx, `SELECT to_regclass('public.users') IS NOT NULL`).Scan(&hasUsers); err != nil {
		t.Fatalf("check schema: %v", err)
	}
	if !hasUsers {
		t.Fatal("users table missing; run: atlas migrate apply --env test")
	}

	return pool
}

// uniqueEmail returns an email no other test run uses, and deletes any row
// with that email when the test finishes.
func uniqueEmail(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()

	email := fmt.Sprintf("it-%d@integration.test", time.Now().UnixNano())
	t.Cleanup(func() {
		_, err := pool.Exec(context.Background(), `DELETE FROM users WHERE email = $1`, email)
		if err != nil {
			t.Errorf("cleanup %s: %v", email, err)
		}
	})

	return email
}

func TestIntegration_CreateAndGetUser(t *testing.T) {
	pool := integrationPool(t)
	repo := NewPostgresUserRepository(pool)
	ctx := context.Background()
	email := uniqueEmail(t, pool)

	before := time.Now()
	id, err := repo.CreateUser(ctx, "Integration User", email)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if id == "" {
		t.Fatal("CreateUser returned an empty id")
	}

	// A: the row really exists in the table.
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM users WHERE email = $1`, email).Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row for %s, got %d", email, count)
	}

	// B: GetUser scans every column correctly.
	user, err := repo.GetUser(ctx, id)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if user.ID != id {
		t.Errorf("expected id %q, got %q", id, user.ID)
	}
	if user.Name != "Integration User" {
		t.Errorf("expected name %q, got %q", "Integration User", user.Name)
	}
	if user.Email != email {
		t.Errorf("expected email %q, got %q", email, user.Email)
	}
	// The database clock may differ slightly from ours, so allow some skew.
	const skew = time.Minute
	if user.CreatedAt.Before(before.Add(-skew)) || user.CreatedAt.After(time.Now().Add(skew)) {
		t.Errorf("expected created_at near %v, got %v", before, user.CreatedAt)
	}

	// E: the status column defaults to "active".
	var status string
	if err := pool.QueryRow(ctx, `SELECT status FROM users WHERE id = $1`, id).Scan(&status); err != nil {
		t.Fatalf("read status: %v", err)
	}
	if status != "active" {
		t.Errorf("expected default status %q, got %q", "active", status)
	}
}

func TestIntegration_GetUser_NotFound(t *testing.T) {
	pool := integrationPool(t)
	repo := NewPostgresUserRepository(pool)

	// C: a well-formed UUID that no row uses.
	_, err := repo.GetUser(context.Background(), "00000000-0000-4000-8000-000000000000")
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestIntegration_CreateUser_DuplicateEmail(t *testing.T) {
	pool := integrationPool(t)
	repo := NewPostgresUserRepository(pool)
	ctx := context.Background()
	email := uniqueEmail(t, pool)

	if _, err := repo.CreateUser(ctx, "First", email); err != nil {
		t.Fatalf("first CreateUser: %v", err)
	}

	// D: the real UNIQUE constraint fires and is translated.
	_, err := repo.CreateUser(ctx, "Second", email)
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}
