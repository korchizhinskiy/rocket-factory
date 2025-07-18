package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"buf.build/go/protovalidate"
	logging_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50052

type inventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer

	parts map[string]*inventoryv1.Part
	mu    sync.RWMutex
}

func (inv *inventoryService) ListPart(
	_ context.Context,
	request *inventoryv1.ListPartRequest,
) (*inventoryv1.ListPartResponse, error) {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	uuidFilter := partUUIDFilter{}
	tagFilter := partTagFilter{}
	nameFilter := partNameFilter{}
	categoryFilter := partCategoryFilter{}
	manufacturerCountriesFilter := partManufacturerCountriesFilter{}

	partsMap := uuidFilter.filter(inv.parts, request.GetFilter().Uuids)
	partsMap = tagFilter.filter(partsMap, request.GetFilter().Tags)
	partsMap = nameFilter.filter(partsMap, request.GetFilter().Names)
	partsMap = categoryFilter.filter(partsMap, request.GetFilter().Categories)
	partsMap = manufacturerCountriesFilter.filter(partsMap, request.GetFilter().ManufactorerCountries)

	partSlice := make([]*inventoryv1.Part, len(partsMap))

	idx := 0
	for _, v := range partsMap {
		partSlice[idx] = v
		idx++
	}

	return &inventoryv1.ListPartResponse{Parts: partSlice}, nil
}

func main() {
	logger := GetLogger()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Info("Failed to listen", slog.Any("error", err))
		return
	}

	defer func() {
		if err := lis.Close(); err != nil {
			logger.Info("Failed to close listener", slog.Any("error", err))
		}
	}()

	validator, err := protovalidate.New()
	if err != nil {
		logger.Info("Failed to create Validators", slog.Any("error", err))
		return
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			protovalidate_middleware.UnaryServerInterceptor(validator),
			logging_middleware.UnaryServerInterceptor(
				InterceptorLogger(logger),
				[]logging_middleware.Option{
					logging_middleware.WithLogOnEvents(
						logging_middleware.StartCall,
						logging_middleware.FinishCall,
					),
				}...),
		),
	)

	reflection.Register(server)
	service := &inventoryService{parts: generateParts()}
	inventoryv1.RegisterInventoryServiceServer(server, service)

	go func() {
		logger.Info("gRPC server listening", slog.Any("port", grpcPort))
		err := server.Serve(lis)
		if err != nil {
			logger.Info("Failed to serve", slog.Any("error", err))
			return
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down gRPC server...")
	server.GracefulStop()
	logger.Info("Server stopped")
}

func GetLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func InterceptorLogger(l *slog.Logger) logging_middleware.Logger {
	return logging_middleware.LoggerFunc(
		func(ctx context.Context, lvl logging_middleware.Level, msg string, fields ...any) {
			l.Log(ctx, slog.Level(lvl), msg, fields...)
		},
	)
}
