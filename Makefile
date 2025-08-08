PROTO_SRC := proto/auth.proto
OUT_DIR := pkg/pb/auth

generate-proto:
	protoc \
	  --proto_path=proto \
	  --go_out=paths=source_relative:$(OUT_DIR) \
	  --go-grpc_out=paths=source_relative:$(OUT_DIR) \
	  $(PROTO_SRC)
