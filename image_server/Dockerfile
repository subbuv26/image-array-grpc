FROM golang:1.18 as builder
WORKDIR /workspace/image_server

COPY main.go .
COPY go.mod .
COPY go.sum .
COPY vendor vendor

# RUN go mod tidy
# RUN go mod download


# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o image_server main.go 
RUN CGO_ENABLED=0 go build -a -o image_server main.go 


# RUN CGO_ENABLED=0 go get github.com/grpc-ecosystem/grpc-health-probe
# RUN CGO_ENABLED=0 go get github.com/aws-samples/grpc-examples/examples/helloworld/image_server

FROM scratch
COPY --from=builder /workspace/image_server/image_server image_server
EXPOSE 50051
ENTRYPOINT ["/image_server"]
