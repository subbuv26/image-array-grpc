/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"log"
	"net"

	"github.com/google/uuid"

	pb "github.com/subbuv26/image-array-grpc/proto/image"

	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

const (
	port = ":50051"
)

type imageStore struct {
	images map[string]*pb.Image
}

var store imageStore

// server is used to implement helloworld.GreeterServer.
type imageServer struct {
	pb.UnimplementedImageServiceServer
}

func (is *imageServer) ListImages(
	ctx context.Context, in *pb.ListImagesRequest,
) (*pb.ListImagesResponse, error) {
	log.Printf("ListImages with: %v", in)

	maxImgs := len(store.images)
	if maxImgs > int(in.GetMaxImages()) {
		maxImgs = int(in.GetMaxImages())
	}

	images := []*pb.Image{}
	count := 0
	for _, image := range store.images {
		images = append(images, image)
		count++
		if count == maxImgs {
			break
		}
	}

	return &pb.ListImagesResponse{
		Images: images,
	}, nil
}

func (is *imageServer) GetImage(
	ctx context.Context, in *pb.GetImageRequest,
) (*pb.Image, error) {
	log.Printf("GetImage with: %v", in)

	if image, ok := store.images[in.GetId()]; ok {
		return image, nil
	}
	return nil, grpc.Errorf(grpcCodes.NotFound, "Image not Found")
}

func (is *imageServer) CreateImage(
	ctx context.Context, in *pb.CreateImageRequest,
) (*pb.CreateImageResponse, error) {
	log.Printf("CreateImage with: %v", in)
	id := uuid.NewString()
	store.images[id] = in.GetImage()

	return &pb.CreateImageResponse{Id: id}, nil
}

func (is *imageServer) UpdateImage(
	ctx context.Context, in *pb.UpdateImageRequest,
) (*pb.StatusResponse, error) {
	log.Printf("UpdateImage with: %v", in)

	store.images[uuid.NewString()] = in.GetImage()

	return &pb.StatusResponse{Success: true}, nil
}

func (is *imageServer) DeleteImage(
	ctx context.Context, in *pb.DeleteImageRequest,
) (*pb.StatusResponse, error) {
	log.Printf("DeleteImage with: %v", in)

	delete(store.images, in.GetId())
	return &pb.StatusResponse{Success: true}, nil
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	store = imageStore{
		images: make(map[string]*pb.Image),
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcAuth.UnaryServerInterceptor(AuthFunc),
		)))
	pb.RegisterImageServiceServer(s, &imageServer{})
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// AuthFunc is a middleware (interceptor) that extracts token from header
func AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}
	log.Printf("Token: %v", token)

	return ctx, nil
}
