# go-grpc-starter

A very small gRPC service in Go, built to learn the fundamentals of gRPC and Protocol Buffers. It's a skeleton with an in-memory `UserService` with `CreateUser` and `GetUser` methods.

---

## What it does

A client calls a server over the network to create and fetch users (stored in an in-memory Go map)

```
Created user: id:"user-1"  name:"John Doe"  email:"john.doe@example.com"
Fetched user: id:"user-1"  name:"John Doe"  email:"john.doe@example.com"
```

---

## gRPC? protobuf?

**gRPC** is a framework for calling methods on another service over a network, as if they were local function calls. It's commonly used for service-to-service communication inside a large enterprise (rather than public-facing APIs), where speed and type-safety matter the most.

**Protocol Buffers (protobuf)** is the binary data format and interface definition language gRPC is built on. You define your service's methods and data structures once, in a `.proto` file, and a compiler generates matching client and server code from it (not only for Go).

---

## How gRPC works in this context

**1. Define a contract** — [`proto/user.proto`](proto/user.proto) declares the service's methods and message shapes:

```proto
service UserService {
  rpc GetUser(GetUserRequest) returns (User);
  rpc CreateUser(CreateUserRequest) returns (User);
}

message User {
  string id = 1;
  string name = 2;
  string email = 3;
}
```

This file doesnt have any logic, it only defines the shape of the data and the names of the methods. Note the field numbers (`= 1`, `= 2`, `= 3`) arent the actual values, but rather they are identifiers used in the binary encoding. Once used in a real system they should never be reused (they should be retired), even if the field is later removed.

**2. Generate code from it** — `protoc` reads `user.proto` and generates [`gen/userpb/user.pb.go`](gen/userpb/user.pb.go) (message structs) and [`gen/userpb/user_grpc.pb.go`](gen/userpb/user_grpc.pb.go) (service interface + client stub). Nothing in `gen/` is hand-written. Regenerate it any time the proto changes, don't edit it directly.

**3. Implement the server** — [`server/main.go`](server/main.go) defines a `server` struct that implements the generated `UserServiceServer` interface:

```go
func (s *server) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user with id %s not found", req.Id)
	}
	return user, nil
}
```

It embeds `userpb.UnimplementedUserServiceServer`, which provides the default implementations for every method in the interface. This means if the proto file gains a new method later, existing code still compiles. The new method just returns an "unimplemented" error until you actually write it.

**4. Call it from a client** — [`client/main.go`](client/main.go) connects and calls those same methods as if they were local functions:

```go
newUser, err := client.CreateUser(ctx, &userpb.CreateUserRequest{
	Name:  "John Doe",
	Email: "john.doe@example.com",
})
```

The request is serialised to binary, sent over HTTP/2 (TCP), deserialised on the server, routed to `CreateUser`, and the response makes the same journey back. All of which is handled by the generated code.

**Why not REST/JSON here:** with REST there's no compile time guarantee the client's request shape matches what the server expects. This means that mismatches surface at runtime. Here, both sides are generated from the same `.proto` file, so mismatches are usually caught at build time instead.

---

## Important notes

- **Field numbers are permanent.** If you remove a field, mark it `reserved` so its never accidentally reused:
  ```proto
  reserved 3;
  reserved "email";
  ```
- **Fields are implicitly optional in proto3**, and by default there's no way to tell "never set" apart from "set to empty/zero." Mark a field `optional` if that distinction matters (it changes the generated Go field into a pointer).
- **Pointers:** every generated message type is used as a pointer (`*userpb.User`, `*userpb.CreateUserRequest`) rather than a value, to avoid copying potentially large structs on every call.

---

## Real usecase 

Client and server are usually separate services owned by separate teams. One team owns and runs the server (e.g. an authorisation service exposing `IsAuthorised(userID, resource, action)`). Other teams across the company are clients of it, each generating their own client code from the same shared `.proto` file, regardless of what language their own service is written in.

---

## Prereqs

- [Go](https://go.dev/dl/) 1.21+
- [`protoc`](https://grpc.io/docs/protoc-installation/)
- Go plugins for protoc:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

---

## Setup

```bash
git clone <this-repo>
cd go-grpc-starter
go mod tidy
```

---

## Generating/regnerating code after changes to `user.proto`

```bash
protoc --go_out=gen/userpb --go_opt=paths=source_relative \
  --go-grpc_out=gen/userpb --go-grpc_opt=paths=source_relative \
  proto/user.proto
```

---

## How to run

**Terminal 1 - server:**
```bash
go run server/main.go
```
```
Server listening on :50051
```

**Terminal 2 — client:**
```bash
go run client/main.go
```
```
Created user: id:"user-1"  name:"John Doe"  email:"john.doe@example.com"
Fetched user: id:"user-1"  name:"John Doe"  email:"john.doe@example.com"
```

---

## Project structure

```
go-grpc-starter/
├── proto/
│   └── user.proto          # the service contract
├── gen/userpb/
│   ├── user.pb.go          # generated message structs (don't edit)
│   └── user_grpc.pb.go     # generated service interface + client stub (don't edit)
├── server/
│   └── main.go             # server implementation
├── client/
│   └── main.go             # client that calls the server
├── go.mod
└── README.md
```

---

## Why i made this

A weekend project to understand gRPC and Go fundamentals :)