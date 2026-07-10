# go-grpc-starter

```
protoc --go_out=gen/userpb --go_opt=paths=source_relative \
  --go-grpc_out=gen/userpb --go-grpc_opt=paths=source_relative \
  proto/user.proto
```