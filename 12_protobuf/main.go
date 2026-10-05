// Command 12_protobuf runs UserService over HTTP using ConnectRPC.
//
// Run it from the repository root with:
//
//	go run ./12_protobuf
package main

import (
	"net/http"

	"example/goWorld/12_protobuf/gen/user/v1/userv1connect"
	"example/goWorld/12_protobuf/userservice"
	"fmt"
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
	// TODO 1: Create your UserService (from the userservice package).
	//         Remember that CreateUser has a pointer receiver (*UserService).
	repo := &userservice.InMemoryUserRepository{}

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

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		fmt.Println("something failed:", err)
	}
}
