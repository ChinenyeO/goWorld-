package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Development defaults, used when the corresponding variable is unset or empty.
const (
	defaultMaxConns        = 10
	defaultMinConns        = 0
	defaultMaxConnLifetime = time.Hour
	defaultMaxConnIdleTime = 30 * time.Minute
)

// Config is the service configuration, read from environment variables.
//
//	DATABASE_URL           required; never logged because it contains credentials
//	DB_MAX_CONNS           maximum open connections (default 10)
//	DB_MIN_CONNS           connections kept open when idle (default 0)
//	DB_MAX_CONN_LIFETIME   recycle connections older than this (default 1h)
//	DB_MAX_CONN_IDLE_TIME  close connections idle longer than this (default 30m)
//
// Durations use Go syntax, e.g. 30s, 5m, 1h30m. The DB_* variables always
// take precedence over pool_* parameters in DATABASE_URL.
type Config struct {
	DatabaseURL     string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// loadConfig reads Config from the environment. A variable that is set but
// invalid is an error, never silently replaced by its default. All problems
// are reported together.
func loadConfig() (Config, error) {
	cfg := Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		MaxConns:        defaultMaxConns,
		MinConns:        defaultMinConns,
		MaxConnLifetime: defaultMaxConnLifetime,
		MaxConnIdleTime: defaultMaxConnIdleTime,
	}

	var errs []error
	if cfg.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}

	parseErrs := errors.Join(
		envInt32("DB_MAX_CONNS", &cfg.MaxConns),
		envInt32("DB_MIN_CONNS", &cfg.MinConns),
		envDuration("DB_MAX_CONN_LIFETIME", &cfg.MaxConnLifetime),
		envDuration("DB_MAX_CONN_IDLE_TIME", &cfg.MaxConnIdleTime),
	)
	if parseErrs != nil {
		// Range checks would be misleading for values that did not parse.
		return Config{}, errors.Join(append(errs, parseErrs)...)
	}

	if cfg.MaxConns < 1 {
		errs = append(errs, fmt.Errorf("DB_MAX_CONNS must be at least 1, got %d", cfg.MaxConns))
	}
	if cfg.MinConns < 0 {
		errs = append(errs, fmt.Errorf("DB_MIN_CONNS must not be negative, got %d", cfg.MinConns))
	}
	if cfg.MinConns > cfg.MaxConns {
		errs = append(errs, fmt.Errorf("DB_MIN_CONNS (%d) must not exceed DB_MAX_CONNS (%d)", cfg.MinConns, cfg.MaxConns))
	}
	if cfg.MaxConnLifetime <= 0 {
		errs = append(errs, fmt.Errorf("DB_MAX_CONN_LIFETIME must be positive, got %s", cfg.MaxConnLifetime))
	}
	if cfg.MaxConnIdleTime <= 0 {
		errs = append(errs, fmt.Errorf("DB_MAX_CONN_IDLE_TIME must be positive, got %s", cfg.MaxConnIdleTime))
	}

	if err := errors.Join(errs...); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// poolConfig builds a pgxpool.Config from DATABASE_URL and applies the pool
// settings on top of it. It does not connect to the database.
func (c Config) poolConfig() (*pgxpool.Config, error) {
	poolCfg, err := pgxpool.ParseConfig(c.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}

	poolCfg.MaxConns = c.MaxConns
	poolCfg.MinConns = c.MinConns
	poolCfg.MaxConnLifetime = c.MaxConnLifetime
	poolCfg.MaxConnIdleTime = c.MaxConnIdleTime

	return poolCfg, nil
}

// envInt32 overwrites *dst with the named variable if it is set and non-empty.
func envInt32(name string, dst *int32) error {
	s := os.Getenv(name)
	if s == "" {
		return nil
	}

	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return fmt.Errorf("%s: %q is not a valid 32-bit integer", name, s)
	}

	*dst = int32(n)
	return nil
}

// envDuration overwrites *dst with the named variable if it is set and non-empty.
func envDuration(name string, dst *time.Duration) error {
	s := os.Getenv(name)
	if s == "" {
		return nil
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("%s: %q is not a valid duration (examples: 30s, 5m, 1h)", name, s)
	}

	*dst = d
	return nil
}
