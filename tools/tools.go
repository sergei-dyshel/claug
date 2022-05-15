//go:build tools
// +build tools

package tools

import (
	_ "github.com/a-h/generate/cmd/schema-generate"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
