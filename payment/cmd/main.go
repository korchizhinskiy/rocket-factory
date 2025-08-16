package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"buf.build/go/protovalidate"
	"github.com/google/uuid"
	logging_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	paymentv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50051

type paymentService struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func (s *paymentService) PayOrder(
	_ context.Context,
	request *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	if !s.IsPaymentMethodValid(request) {
		return nil, status.Error(codes.InvalidArgument, "Invalid Payment method")
	}
	return &paymentv1.PayOrderResponse{
		TransactionUuid: uuid.NewString(),
	}, nil
}

func (s *paymentService) IsPaymentMethodValid(request *paymentv1.PayOrderRequest) bool {
	switch request.PaymentMethod {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
		paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
		paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return true
	default:
		return false
	}
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

	defer func() {
		if err := lis.Close(); err != nil {
			logger.Info(
				"Failed to close listener",
				slog.Any("error", err),
			)
		}
	}()
	validator, err := protovalidate.New()
	if err != nil {
		logger.Info(
			"Failed to create validators",
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
	service := &paymentService{}

	paymentv1.RegisterPaymentServiceServer(server, service)

	go func() {
		log.Printf(
			"gRPC server listening on %d\n",
			grpcPort,
		)
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
