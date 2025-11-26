gen:
	@protoc \
		--proto_path=protobuf "protobuf/tour.proto" \
		--go_out=protobuf --go_opt=paths=source_relative \
		--go-grpc_out=protobuf --go-grpc_opt=paths=source_relative