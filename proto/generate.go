package proto

//go:generate protoc --proto_path=. --go_out=paths=source_relative:../pkg/pb/auth --go-grpc_out=paths=source_relative:../pkg/pb/auth auth.proto
