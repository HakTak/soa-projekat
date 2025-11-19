gen:
	@protoc \
		--proto_path=protobuf "protobuf/stakeholders.proto" \
		--go_out=common/genproto --go_opt=paths=source_relative \
		--go-grpc_out=common/genproto --go-grpc_opt=paths=source_relative