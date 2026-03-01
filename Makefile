
proto:
	protoc --proto_path=./protobuf/proto --go_out=./protobuf/gen/auth --go_opt=paths=source_relative --go-grpc_out=./protobuf/gen/auth --go-grpc_opt=paths=source_relative protobuf/proto/auth.proto
	protoc --proto_path=./protobuf/proto --go_out=./protobuf/gen/goods --go_opt=paths=source_relative --go-grpc_out=./protobuf/gen/goods --go-grpc_opt=paths=source_relative protobuf/proto/goods.proto
