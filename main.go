package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"buf.build/go/protovalidate"
	"github.com/susana-garcia/go-crud/config"
	"github.com/susana-garcia/go-crud/server"
	"github.com/susana-garcia/go-crud/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

//go:generate ./scripts/generate-pb.sh

func main() {
	// load configuration from environment variables
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	logger.Info("starting server on", "host", cfg.Server.Host, "port", cfg.Server.Port)

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("unable to listen to tcp port", "port", cfg.Server.Port, "error", err)
		os.Exit(1)
	}
	defer func() {
		logger.Info("closing tcp connection")
		_ = listener.Close()
	}()

	logger.Info("connecting to database", "name", cfg.Name)

	db := config.OpenConnection(cfg.Database)

	// run DB migration
	logger.Info("running database migration for blogs table")
	err = db.AutoMigrate(&service.Blog{})
	if err != nil {
		logger.Error("error running auto migrate", "err", err)
		os.Exit(1)
	}
	logger.Info("database migration completed successfully")

	bService := service.New(db, logger)
	server := server.New(bService, logger)

	// create protovalidate validator
	validator, err := protovalidate.New()
	if err != nil {
		logger.Error("failed to create validator", "error", err)
		os.Exit(1)
	}

	// create gRPC server with validation interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
			// validate request using protovalidate if it's a protobuf message
			if msg, ok := req.(proto.Message); ok {
				if err := validator.Validate(msg); err != nil {
					return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
				}
			}
			return handler(ctx, req)
		}),
	)
	server.Register(s)

	// enable server reflection so tools like grpcurl can discover services without a proto file
	reflection.Register(s)

	logger.Info(fmt.Sprintf("server listening on %s", address))

	err = s.Serve(listener)
	if err != nil {
		logger.Error("unable to start server", "error", err)
	}
}
