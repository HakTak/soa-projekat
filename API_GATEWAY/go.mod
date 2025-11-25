module PROJEKAT/API_GATEWAY

go 1.25.4

require (
	PROJEKAT/COMMON v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	google.golang.org/grpc v1.77.0
)

require (
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251111163417-95abcf5c77ba // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace PROJEKAT/COMMON => ../COMMON
