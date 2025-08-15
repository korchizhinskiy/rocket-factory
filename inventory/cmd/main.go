package main

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"buf.build/go/protovalidate"
	logging_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50052

type inventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer

	parts map[string]*inventoryv1.Part
	mu    sync.RWMutex
}

func (inv *inventoryService) GetPart(
	_ context.Context,
	request *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	part, ok := inv.parts[request.Uuid]

	if !ok {
		return nil, status.Error(
			codes.NotFound,
			"part not found",
		)
	}
	return &inventoryv1.GetPartResponse{Part: part}, nil
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

	parts := inv.parts

	if filter := request.GetFilter(); filter != nil {
		parts = uuidFilter.filter(
			parts,
			request.GetFilter().Uuids,
		)
		parts = tagFilter.filter(
			parts,
			request.GetFilter().Tags,
		)
		parts = nameFilter.filter(
			parts,
			request.GetFilter().Names,
		)
		parts = categoryFilter.filter(
			parts,
			request.GetFilter().Categories,
		)
		parts = manufacturerCountriesFilter.filter(
			parts,
			request.GetFilter().ManufactorerCountries,
		)

	}
	var partList []*inventoryv1.Part
	maps.Values(parts)(func(p *inventoryv1.Part) bool {
		partList = append(partList, p)
		return true
	})
	return &inventoryv1.ListPartResponse{Parts: partList}, nil
}

func main() {
	logger := GetLogger()
	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", grpcPort),
	)
	if err != nil {
		logger.Info(
			"Failed to listen",
			slog.Any("error", err),
		)
		return
	}

	validator, err := protovalidate.New()
	if err != nil {
		logger.Info(
			"Failed to create Validators",
			slog.Any("error", err),
		)
		return
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			protovalidate_middleware.UnaryServerInterceptor(
				validator,
			),
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
	inventoryv1.RegisterInventoryServiceServer(
		server,
		service,
	)

	go func() {
		logger.Info(
			"gRPC server listening",
			slog.Any("port", grpcPort),
		)
		err := server.Serve(lis)
		if err != nil {
			logger.Info(
				"Failed to serve",
				slog.Any("error", err),
			)
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
	return slog.New(
		slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
}

func InterceptorLogger(
	l *slog.Logger,
) logging_middleware.Logger {
	return logging_middleware.LoggerFunc(
		func(ctx context.Context, lvl logging_middleware.Level, msg string, fields ...any) {
			l.Log(ctx, slog.Level(lvl), msg, fields...)
		},
	)
}
