// Command 12_protobuf runs UserService over HTTP using ConnectRPC.
//
// It requires DATABASE_URL to point at PostgreSQL; see Config in config.go for
// the optional DB_* pool settings. Run it from the repository root with:
//
//	DATABASE_URL='postgres://goworld:goworld@localhost:5432/goworld?sslmode=disable' go run ./12_protobuf
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"example/goWorld/12_protobuf/gen/user/v1/userv1connect"
	"example/goWorld/12_protobuf/userservice"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ---------------------------------------------------------------------------
// HINT 1:
// A standard Go HTTP mux can be created with:
//
//	mux := http.NewServeMux()
//
// HINT 2:
// A generated ConnectRPC handler constructor conceptually returns:
//
//	path, handler := someconnect.NewSomeServiceHandler(service)
//
// HINT 3:
// A handler can be registered with:
//
//	mux.Handle(path, handler)
//
// HINT 4:
// A basic Go HTTP server can be started with:
//
//	http.ListenAndServe(":9000", mux)
//
// For this exercise, use port 8080.
// ---------------------------------------------------------------------------

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	poolCfg, err := cfg.poolConfig()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	// Startup gets 5 seconds to reach PostgreSQL; requests use their own contexts.
	startupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// One pool for the whole process, shared by every request.
	pool, err := pgxpool.NewWithConfig(startupCtx, poolCfg)
	if err != nil {
		log.Fatalf("create database pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(startupCtx); err != nil {
		// log.Fatalf exits without running deferred calls, so close explicitly.
		pool.Close()
		log.Fatalf("ping database: %v", err)
	}

	log.Printf("database pool: max_conns=%d min_conns=%d max_conn_lifetime=%s max_conn_idle_time=%s",
		poolCfg.MaxConns, poolCfg.MinConns, poolCfg.MaxConnLifetime, poolCfg.MaxConnIdleTime)

	// TODO 1: Create your UserService (from the userservice package).
	//         Remember that CreateUser has a pointer receiver (*UserService).
	repo := userservice.NewPostgresUserRepository(pool)

	service := userservice.NewUserService(repo)

	// TODO 2: Call userv1connect.NewUserServiceHandler with your service.
	//         It gives you back a path and an http.Handler.

	path, handler := userv1connect.NewUserServiceHandler(service)

	// TODO 3: Create an http.ServeMux.

	mux := http.NewServeMux()

	// TODO 4: Register the generated path and handler on the mux.

	mux.Handle(path, handler)

	// TODO 5: Start an HTTP server on :8080 using the mux.
	//         ListenAndServe returns an error; consider what to do with it.

	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		fmt.Println("something failed:", err)
	}
}
