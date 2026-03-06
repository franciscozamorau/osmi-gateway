module github.com/franciscozamorau/osmi-gateway

go 1.24.0

require (
	github.com/franciscozamorau/osmi-protobuf v0.0.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0
	google.golang.org/grpc v1.79.1
)

require (
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260209200024-4cfbd4190f57 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260209200024-4cfbd4190f57 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/franciscozamorau/osmi-protobuf => ../osmi-protobuf
