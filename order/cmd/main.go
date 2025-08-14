package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type SwaggerDoc struct {
	URL  string
	Name string
}

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderv1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderv1.OrderDto),
	}
}

type OrderHandler struct {
	storage   *OrderStorage
	invClient inventoryv1.InventoryServiceClient
	payClient paymentv1.PaymentServiceClient
}

func (o *OrderHandler) GetOrderByID(
	ctx context.Context,
	params orderv1.GetOrderByIDParams,
) (orderv1.GetOrderByIDRes, error) {
	order, ok := o.storage.orders[params.OrderUUID.String()]
	if !ok {
		return &orderv1.NotFoundError{
			Code:    404,
			Message: "Order was not found.",
		}, nil
	}

	return order, nil
}

func (o *OrderHandler) CancelOrder(
	ctx context.Context,
	params orderv1.CancelOrderParams,
) (orderv1.CancelOrderRes, error) {
	order, ok := o.storage.orders[params.OrderUUID.String()]
	if !ok {
		return &orderv1.NotFoundError{
			Code:    404,
			Message: "Order was not found.",
		}, nil
	}
	if order.Status.Value == orderv1.OrderStatusCANCELLED {
		return &orderv1.ConflictError{
			Code:    409,
			Message: "Заказ уже отменен",
		}, nil
	}
	if order.Status.Value == orderv1.OrderStatusPAID {
		return &orderv1.ConflictError{
			Code:    409,
			Message: "Заказ уже оплачен и не может быть отменен",
		}, nil
	}
	(&order.Status).SetTo(orderv1.OrderStatusCANCELLED)
	return (orderv1.CancelOrderRes)(nil), nil
}

func (o *OrderHandler) PayOrder(
	ctx context.Context,
	req *orderv1.OrderPayRequest,
	params orderv1.PayOrderParams,
) (orderv1.PayOrderRes, error) {
	order, ok := o.storage.orders[params.OrderUUID.String()]
	if !ok {
		return &orderv1.NotFoundError{
			Code:    404,
			Message: "Order was not found.",
		}, nil
	}

	resp, err := o.payClient.PayOrder(
		ctx,
		&paymentv1.PayOrderRequest{
			OrderUuid: params.OrderUUID.String(),
			UserUuid:  req.UserUUID.String(),
			PaymentMethod: ConvertPaymentMethodStrToProto(
				req.PaymentMethod,
			),
		},
	)
	if err != nil {
		return &orderv1.InternalError{
			Code:    500,
			Message: "Internal Error",
		}, nil
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return &orderv1.InternalError{
			Code:    500,
			Message: "Internal Error",
		}, nil
	}
	(&order.Status).SetTo(orderv1.OrderStatusPAID)
	(&order.PatmentMethod).SetTo(req.PaymentMethod)

	return &orderv1.OrderPayResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

func (o *OrderHandler) CreateOrder(
	ctx context.Context,
	req *orderv1.OrderCreateRequest,
) (orderv1.CreateOrderRes, error) {
	o.storage.mu.RLock()
	defer func() {
		o.storage.mu.RUnlock()
	}()

	resp, err := o.invClient.ListPart(
		ctx,
		&inventoryv1.ListPartRequest{
			Filter: &inventoryv1.PartsFilter{
				Uuids: slices.Collect(
					func(yield func(string) bool) {
						for _, u := range req.PartUuids {
							yield(u.String())
						}
					},
				),
			},
		},
	)
	if err != nil {
		return &orderv1.OrderCreateResponse{
			OrderUUID:  uuid.UUID{},
			TotalPrice: 1000,
		}, nil
	}
	// TODO: Сделать проверку на вхождение всех элементов
	if len(resp.Parts) != len(req.PartUuids) {
		return &orderv1.NotFoundError{
			Code:    404,
			Message: "Some of part was not fount in inventory",
		}, nil
	}
	orderUUID, err := uuid.NewUUID()
	if err != nil {
		return &orderv1.InternalError{
			Code:    500,
			Message: "Internal Error",
		}, nil
	}
	transactionUUID, err := uuid.NewUUID()
	if err != nil {
		return &orderv1.InternalError{
			Code:    500,
			Message: "Internal Error",
		}, nil
	}
	order := orderv1.OrderDto{
		OrderUUID:  orderv1.NewOptUUID(orderUUID),
		UserUUID:   orderv1.NewOptUUID(req.UserUUID),
		PartUuids:  req.PartUuids,
		TotalPrice: orderv1.NewOptFloat64(1000),
		TransactionUUID: orderv1.NewOptUUID(
			transactionUUID,
		),
		PatmentMethod: orderv1.OptPaymentMethod{},
		Status: orderv1.NewOptOrderStatus(
			orderv1.OrderStatusPENDINGPAYMENT,
		),
	}

	var totalPrice float64
	for _, part := range resp.Parts {
		totalPrice += part.Price
	}

	o.storage.orders[orderUUID.String()] = &order
	return &orderv1.OrderCreateResponse{
		OrderUUID:  order.OrderUUID.Value,
		TotalPrice: totalPrice,
	}, nil
}

func (o *OrderHandler) NewError(
	_ context.Context,
	err error,
) *orderv1.GenericErrorStatusCode {
	return &orderv1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderv1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

func NewOrderHandler(
	storage *OrderStorage,
	invClient inventoryv1.InventoryServiceClient,
	payClient paymentv1.PaymentServiceClient,
) *OrderHandler {
	return &OrderHandler{
		storage:   storage,
		invClient: invClient,
		payClient: payClient,
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	fileServer := http.FileServer(http.Dir("shared/api"))

	invConn, err := grpc.NewClient(
		net.JoinHostPort("localhost", "50052"),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		log.Printf("Ошибка подключения к сервису Inventory")
		return
	}
	payConn, err := grpc.NewClient(
		net.JoinHostPort("localhost", "50051"),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		log.Printf("Ошибка подключения к сервису Payment")
		return
	}
	invClient := inventoryv1.NewInventoryServiceClient(
		invConn,
	)
	payClient := paymentv1.NewPaymentServiceClient(payConn)
	storage := NewOrderStorage()
	orderHandler := NewOrderHandler(
		storage,
		invClient,
		payClient,
	)
	orderServer, err := orderv1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("Ошибка создания сервера")
	}
	r.Mount("/", orderServer)
	r.Handle(
		"/docs/*",
		http.StripPrefix("/docs/", fileServer),
	)
	r.Get("/docs/swagger", SwaggerUIHandler(
		"https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css",
		"https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js",
		"https://shorturl.at/YUNWT",
		"Rocket Factory API Docs",
		"/docs/order/v1/order.openapi.yaml",
	))

	srv := &http.Server{
		Addr: net.JoinHostPort(
			"localhost",
			httpPort,
		),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	go func() {
		log.Printf(
			"🚀 HTTP-сервер запущен на http://localhost:%s\n",
			httpPort,
		)
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Printf(
				"❌ Ошибка запуска сервера: %v\n",
				err,
			)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf(
			"❌ Ошибка при остановке сервера: %v\n",
			err,
		)
	} else {
		log.Println("✅ Сервер остановлен")
	}
}
