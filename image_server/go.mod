module github.com/subbuv26/image-array-grpc/image_server

replace github.com/subbuv26/image-array-grpc/proto/image => ../proto/image

go 1.18

require (
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/subbuv26/image-array-grpc/proto/image v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.53.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.29.0 // indirect
)
