package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"buf.build/go/protovalidate"
	logging_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

type inventoryService struct {
	inventoryv1.InventoryServiceServer

	mu sync.RWMutex
}

func (s *inventoryService) ListPart(_ context.Context, request *inventoryv1.ListPartRequest) (*inventoryv1.ListPartResponse, error) {
	parts := []*inventoryv1.Part{
		{
			Uuid:          "da8bf2b2-7278-4ff8-bc99-864315109354",
			Name:          "Hello",
			Price:         100,
			StockQuantity: 100,
			Category:      inventoryv1.PartCategory_PART_CATEGORY_ENGINE,
			Tags:          []string{},
		},
	}
	log.Print(parts[0].Uuid)
	return &inventoryv1.ListPartResponse{Parts: parts}, nil
}

func main() {
	logger := GetLogger()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Info("Failed to listen: %v\n", err)
		return
	}

	defer func() {
		if err := lis.Close(); err != nil {
			logger.Info("Failed to close listener: %v\n", err)
		}
	}()
	validator, _ := protovalidate.New()
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			protovalidate_middleware.UnaryServerInterceptor(validator),
			logging_middleware.UnaryServerInterceptor(InterceptorLogger(logger), []logging_middleware.Option{
				logging_middleware.WithLogOnEvents(logging_middleware.StartCall, logging_middleware.FinishCall),
			}...),
		),
	)

	reflection.Register(server)
	service := &inventoryService{}

	inventoryv1.RegisterInventoryServiceServer(server, service)

	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)
		err := server.Serve(lis)
		if err != nil {
			log.Printf("Failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	server.GracefulStop()
	log.Println("Server stopped")

}
func GetLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func InterceptorLogger(l *slog.Logger) logging_middleware.Logger {
	return logging_middleware.LoggerFunc(func(ctx context.Context, lvl logging_middleware.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
