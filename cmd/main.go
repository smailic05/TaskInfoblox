package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/smailic05/TaskInfoblox/internal/config"
	"github.com/smailic05/TaskInfoblox/internal/db"
	"github.com/smailic05/TaskInfoblox/internal/model"
	"github.com/smailic05/TaskInfoblox/internal/pb"
	"github.com/smailic05/TaskInfoblox/internal/service"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	dsn := "host=localhost user=postgres password=postgres dbname=backend port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal().Err(err).Msg("Database connection error")
	}
	err = database.AutoMigrate(&model.User{})
	if err != nil {
		logger.Panic().Err(err).Msg("Failed to migrate models")
	}
	repoDB := db.NewDB(database)
	cfg, err := config.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("Configuration error")
	}
	userHandler := service.New(repoDB)
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, userHandler)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logger.Fatal().Err(err).Msg("Listening gRPC error")
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	go func() {
		logger.Info().Msgf("GRPC server is listening on :%d", cfg.GRPCPort)
		err := grpcServer.Serve(lis)
		if err != nil && err != grpc.ErrServerStopped {
			logger.Fatal().Err(err).Msg("GRPC server error")
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", cfg.GRPCPort), opts)
	if err != nil {
		logger.Fatal().Err(err).Msg("Registering gRPC gateway endpoint error")
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	go func() {
		logger.Info().Msgf("GRPC gateway server is listening on :%d", cfg.Port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("GRPC gateway server error")
		}
	}()

	<-shutdown

	logger.Info().Msg("Shutdown signal received")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("GRPC gateway server shutdown error")
	}

	grpcServer.GracefulStop()

	logger.Info().Msg("Server stopped gracefully")
}
