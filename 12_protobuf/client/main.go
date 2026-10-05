// Command client calls UserService.CreateUser on a running server.
//
// Start the server first (from the repository root):
//
//	go run ./12_protobuf
//
// Then, in another terminal:
//
//	go run ./12_protobuf/client
package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	userv1 "example/goWorld/12_protobuf/gen/user/v1"
	"example/goWorld/12_protobuf/gen/user/v1/userv1connect"
)

// ---------------------------------------------------------------------------
// HINT 1:
// A generated client can conceptually be created like:
//
//	client := productconnect.NewProductServiceClient(
//	    http.DefaultClient,
//	    "http://localhost:9000",
//	)
//
// HINT 2:
// A protobuf message can be constructed like:
//
//	msg := &productv1.CreateProductRequest{
//	    Name: "Keyboard",
//	}
//
// HINT 3:
// A ConnectRPC request can be wrapped with:
//
//	req := connect.NewRequest(msg)
//
// HINT 4:
// An RPC returning a response and error can be called like:
//
//	response, err := client.SomeMethod(context.Background(), req)
//
// HINT 5:
// The protobuf response message is inside:
//
//	response.Msg
// ---------------------------------------------------------------------------

func main() {
	// TODO 1: Create a UserService client with userv1connect.NewUserServiceClient.
	//         Your server is running at http://localhost:8080.

	client := userv1connect.NewUserServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	// TODO 2: Create a CreateUserRequest (from the userv1 package) with:
	//           name  = "Chinenye"
	//           email = "chinenye@example.com"

	msg := &userv1.CreateUserRequest{
		Name:  "Chinenye",
		Email: "chinenye@example.com",
	}

	// TODO 3: Wrap the message with connect.NewRequest.

	req := connect.NewRequest(msg)

	// TODO 4: Call client.CreateUser using context.Background().
	response, err := client.CreateUser(context.Background(), req)

	// TODO 5: Handle the returned error.
	if err != nil {
		fmt.Println("CreateUser failed", err)
		return
	}

	// TODO 6: Read the user ID from the response and print it with fmt.

	fmt.Println("Created user:", response.Msg.GetId())
}
