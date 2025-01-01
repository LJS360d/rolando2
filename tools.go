//go:build tools

package tools

import (
	_ "github.com/air-verse/air"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
