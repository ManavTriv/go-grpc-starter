# go-grpc-starter

A very small gRPC service in Go, built to learn the fundamentals of gRPC and Protocol Buffers. It's a skeleton with a `UserService` containing `CreateUser` and `GetUser` (all in memory).

---

## What it does

A client calls a server over the network to create and fetch users, stored in a Go map.

```
Created user: id:"user-1"  name:"John Doe"  email:"john.doe@example.com"
Fetched user: id:"user-1"  name:"John Doe"  email:"john.doe@example.com"
```

---

## gRPC and protobuf

**gRPC** lets you call methods on another service over a network as if they were local function calls. It's mostly used for service-to-service communication inside a company, rather than public-facing APIs, where speed and type safety matter more than readability.

**Protocol Buffers (protobuf)** is the binary format and interface definition language gRPC is built on. You define your service's methods and data structures once in a `.proto` file, and a compiler generates matching client and server code from it.

---

## How it works here

**1. Define the contract** — [`proto/user.proto`](proto/user.proto) declares the service's methods and message shapes:

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

Contains no logic but only just the shape of the data and the method names. The numbers (`= 1`, `= 2`, `= 3`) aren't values but are field identifiers used in the binary encoding. Once a field number's used in a real system, it should never be reused.

**2. Generate code from it** — `protoc` reads `user.proto` and outputs [`gen/userpb/user.pb.go`](gen/userpb/user.pb.go) (message structs) and [`gen/userpb/user_grpc.pb.go`](gen/userpb/user_grpc.pb.go) (service interface and client stub). Nothing in `gen/` is hand-written. Regenerate it whenever the proto changes but don't touch it directly.

**3. Implement the server** — [`server/main.go`](server/main.go) has a `server` struct implementing the generated `UserServiceServer` interface:

```go
func (s *server) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user with id %s not found", req.Id)
	}
	return user, nil
}
```

It embeds `userpb.UnimplementedUserServiceServer`. This means if the proto gains a new method later and you haven't written it yet, your code still compilest. The new method just returns an "unimplemented" error until you get around to it, instead of breaking the build.

**4. Call it from a client** — [`client/main.go`](client/main.go) connects and calls those methods like they're local functions:

```go
newUser, err := client.CreateUser(ctx, &userpb.CreateUserRequest{
	Name:  "John Doe",
	Email: "john.doe@example.com",
})
```

What happens under the hood is that the request gets serialised to binary, sent over HTTP/2, then deserialised on the server, and routed to `CreateUser`, and the response makes the same trip back. This is all done by the generated code

**Difference between REST/JSON:** with REST, nothing guarantees at compile time that the client's request shape matches what the server expects. You only know when you run if something went wrong. Here both sides are generated from the same `.proto` file, so mismatches usually get caught when you build.

---

## A few things worth knowing

- **Field numbers are forever.** Delete a field, mark its number `reserved` so nobody accidentally reuses it later:
  ```proto
  reserved 3;
  reserved "email";
  ```
- **Fields are implicitly optional in proto3.** By default there's no way to tell "never set" apart from "set to empty/zero." Mark a field `optional` if that distinction actually matters (it turns the generated Go field into a pointer).
- **Everything's a pointer.** Message types are passed around as pointers (`*userpb.User`, `*userpb.CreateUserRequest`) rather than values, so you're not copying potentially large structs on every call.

---

## How this looks in a real system

Client and server are usually owned by different teams. One team runs the server which might be an authorisation service exposing `IsAuthorised(userID, resource, action)`. Everyone else across the company is a client of it, generating their own client code from the same shared `.proto` file, regardless of what language their service is written in.

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

## Regenerating code after changing `user.proto`

```bash
protoc --go_out=gen/userpb --go_opt=paths=source_relative \
  --go-grpc_out=gen/userpb --go-grpc_opt=paths=source_relative \
  proto/user.proto
```

---

## Running it

**Terminal 1 — server:**
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