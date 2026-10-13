package main

import (
	"strings"
	"testing"
	"time"
)

const testDatabaseURL = "postgres://goworld:goworld@localhost:5432/goworld?sslmode=disable"

// setConfigEnv sets every config variable for the duration of the test, so
// values from the developer's shell can't leak in. Variables missing from env
// are set to "", which loadConfig treats as unset.
func setConfigEnv(t *testing.T, env map[string]string) {
	t.Helper()

	for _, name := range []string{
		"DATABASE_URL",
		"DB_MAX_CONNS",
		"DB_MIN_CONNS",
		"DB_MAX_CONN_LIFETIME",
		"DB_MAX_CONN_IDLE_TIME",
	} {
		t.Setenv(name, env[name])
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	setConfigEnv(t, map[string]string{"DATABASE_URL": testDatabaseURL})

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := Config{
		DatabaseURL:     testDatabaseURL,
		MaxConns:        10,
		MinConns:        0,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}
	if cfg != want {
		t.Errorf("expected %+v, got %+v", want, cfg)
	}
}

func TestLoadConfig_Overrides(t *testing.T) {
	setConfigEnv(t, map[string]string{
		"DATABASE_URL":          testDatabaseURL,
		"DB_MAX_CONNS":          "25",
		"DB_MIN_CONNS":          "2",
		"DB_MAX_CONN_LIFETIME":  "45m",
		"DB_MAX_CONN_IDLE_TIME": "5m",
	})

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := Config{
		DatabaseURL:     testDatabaseURL,
		MaxConns:        25,
		MinConns:        2,
		MaxConnLifetime: 45 * time.Minute,
		MaxConnIdleTime: 5 * time.Minute,
	}
	if cfg != want {
		t.Errorf("expected %+v, got %+v", want, cfg)
	}
}

func TestLoadConfig_Invalid(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		// Every substring must appear in the error message.
		wantErr []string
	}{
		{
			name:    "missing DATABASE_URL",
			env:     map[string]string{},
			wantErr: []string{"DATABASE_URL is required"},
		},
		{
			name:    "malformed DB_MAX_CONNS",
			env:     map[string]string{"DATABASE_URL": testDatabaseURL, "DB_MAX_CONNS": "ten"},
			wantErr: []string{`DB_MAX_CONNS: "ten" is not a valid 32-bit integer`},
		},
		{
			name:    "DB_MAX_CONNS overflows int32",
			env:     map[string]string{"DATABASE_URL": testDatabaseURL, "DB_MAX_CONNS": "99999999999"},
			wantErr: []string{`DB_MAX_CONNS: "99999999999" is not a valid 32-bit integer`},
		},
		{
			name:    "malformed DB_MIN_CONNS",
			env:     map[string]string{"DATABASE_URL": testDatabaseURL, "DB_MIN_CONNS": "1.5"},
			wantErr: []string{`DB_MIN_CONNS: "1.5" is not a valid 32-bit integer`},
		},
		{
			name:    "malformed DB_MAX_CONN_LIFETIME",
			env:     map[string]string{"DATABASE_URL": testDatabaseURL, "DB_MAX_CONN_LIFETIME": "1 hour"},
			wantErr: []string{`DB_MAX_CONN_LIFETIME: "1 hour" is not a valid duration`},
		},
		{
			name:    "DB_MAX_CONN_IDLE_TIME without unit",
			env:     map[string]string{"DATABASE_URL": testDatabaseURL, "DB_MAX_CONN_IDLE_TIME": "30"},
			wantErr: []string{`DB_MAX_CONN_IDLE_TIME: "30" is not a valid duration`},
		},
		{
			name:    "DB_MIN_CONNS exceeds DB_MAX_CONNS",
			env:     map[string]string{"DATABASE_URL": testDatabaseURL, "DB_MAX_CONNS": "5", "DB_MIN_CONNS": "6"},
			wantErr: []string{"DB_MIN_CONNS (6) must not exceed DB_MAX_CONNS (5)"},
		},
		{
			name:    "DB_MIN_CONNS exceeds default DB_MAX_CONNS",
			env:     map[string]string{"DATABASE_URL": testDatabaseURL, "DB_MIN_CONNS": "11"},
			wantErr: []string{"DB_MIN_CONNS (11) must not exceed DB_MAX_CONNS (10)"},
		},
		{
			name:    "DB_MAX_CONNS zero",
			env:     map[string]string{"DATABASE_URL": testDatabaseURL, "DB_MAX_CONNS": "0"},
			wantErr: []string{"DB_MAX_CONNS must be at least 1, got 0"},
		},
		{
			name:    "DB_MIN_CONNS negative",
			env:     map[string]string{"DATABASE_URL": testDatabaseURL, "DB_MIN_CONNS": "-1"},
			wantErr: []string{"DB_MIN_CONNS must not be negative, got -1"},
		},
		{
			name:    "DB_MAX_CONN_LIFETIME zero",
			env:     map[string]string{"DATABASE_URL": testDatabaseURL, "DB_MAX_CONN_LIFETIME": "0s"},
			wantErr: []string{"DB_MAX_CONN_LIFETIME must be positive, got 0s"},
		},
		{
			name:    "DB_MAX_CONN_IDLE_TIME negative",
			env:     map[string]string{"DATABASE_URL": testDatabaseURL, "DB_MAX_CONN_IDLE_TIME": "-5m"},
			wantErr: []string{"DB_MAX_CONN_IDLE_TIME must be positive, got -5m0s"},
		},
		{
			name: "several problems reported together",
			env: map[string]string{
				"DB_MAX_CONNS":         "ten",
				"DB_MAX_CONN_LIFETIME": "forever",
			},
			wantErr: []string{
				"DATABASE_URL is required",
				`DB_MAX_CONNS: "ten"`,
				`DB_MAX_CONN_LIFETIME: "forever"`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setConfigEnv(t, test.env)

			_, err := loadConfig()
			if err == nil {
				t.Fatal("expected an error, got nil")
			}

			for _, want := range test.wantErr {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("expected error to contain %q, got %q", want, err.Error())
				}
			}
		})
	}
}

func TestConfigPoolConfig_AppliesSettings(t *testing.T) {
	cfg := Config{
		// pool_max_conns in the URL must be overridden by MaxConns.
		DatabaseURL:     testDatabaseURL + "&pool_max_conns=3",
		MaxConns:        25,
		MinConns:        2,
		MaxConnLifetime: 45 * time.Minute,
		MaxConnIdleTime: 5 * time.Minute,
	}

	poolCfg, err := cfg.poolConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if poolCfg.MaxConns != 25 {
		t.Errorf("expected MaxConns 25, got %d", poolCfg.MaxConns)
	}
	if poolCfg.MinConns != 2 {
		t.Errorf("expected MinConns 2, got %d", poolCfg.MinConns)
	}
	if poolCfg.MaxConnLifetime != 45*time.Minute {
		t.Errorf("expected MaxConnLifetime 45m, got %s", poolCfg.MaxConnLifetime)
	}
	if poolCfg.MaxConnIdleTime != 5*time.Minute {
		t.Errorf("expected MaxConnIdleTime 5m, got %s", poolCfg.MaxConnIdleTime)
	}

	if poolCfg.ConnConfig.Host != "localhost" || poolCfg.ConnConfig.Database != "goworld" {
		t.Errorf("expected connection to localhost/goworld, got %s/%s",
			poolCfg.ConnConfig.Host, poolCfg.ConnConfig.Database)
	}
}

func TestConfigPoolConfig_InvalidURL(t *testing.T) {
	cfg := Config{
		DatabaseURL:     "postgres://goworld:s3cret@localhost:notaport/goworld",
		MaxConns:        10,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}

	_, err := cfg.poolConfig()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "parse DATABASE_URL") {
		t.Errorf("expected error to mention DATABASE_URL, got %q", err.Error())
	}
	if strings.Contains(err.Error(), "s3cret") {
		t.Errorf("expected password to be redacted, got %q", err.Error())
	}
}
