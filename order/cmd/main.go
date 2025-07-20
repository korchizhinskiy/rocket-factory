package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type WeatherHandler struct { }

func NewWeatherHandler() *WeatherHandler {
	return &WeatherHandler{}
}

// func (h *WeatherHandler) NewError(_ context.Context, err error) *weatherV1.GenericErrorStatusCode {
// 	return &weatherV1.GenericErrorStatusCode{
// 		StatusCode: http.StatusInternalServerError,
// 		Response: weatherV1.GenericError{
// 			Code:    weatherV1.NewOptInt(http.StatusInternalServerError),
// 			Message: weatherV1.NewOptString(err.Error()),
// 		},
// 	}
// }

func main() {
	weatherHandler := NewWeatherHandler()

	// weatherServer, err := weatherV1.NewServer(weatherHandler)
	// if err != nil {
	// 	log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	// }

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	// r.Use(customMiddleware.RequestLogger)

	// r.Mount("/", weatherServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
