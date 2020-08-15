package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"

	handler "github.com/warikan/api/handler/rest"
	"github.com/warikan/api/infra"
	"github.com/warikan/api/usecase"
	"github.com/warikan/db"
	"github.com/warikan/log"
)

var (
	port           = ":8080"
	maxconn        = 5
	configFilePath = "_config/config.yaml"
	version        = "unknown"
)

func main() {
	flag.StringVar(&port, "port", port, "tcp host:port to connect")
	flag.StringVar(&configFilePath, "configFilePath", configFilePath, "config filePath")
	flag.IntVar(&maxconn, "maxconn", maxconn, "max db connection")

	log.Init()
	// nolint:errcheck
	defer log.Logger.Sync()
	if err := db.Init(maxconn, configFilePath); err != nil {
		log.Logger.Error("failed to initialize db", zap.Error(err))
		os.Exit(1)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(6, "gzip"))
	r.Use(middleware.StripSlashes)

	cors := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(cors.Handler)

	healthRepository := infra.NewPingPersistencePostgres(db.Pool)
	healthUseCase := usecase.NewHealthUseCase(healthRepository)
	healthHandler := handler.NewHealthHandler(healthUseCase, version)

	paymentRepository := infra.NewPaymentsRepository(db.Pool)
	paymentUsecase := usecase.NewPaymentUseCase(paymentRepository)
	paymentsHandler := handler.NewPaymentsHandler(paymentUsecase)

	r.Route("/warikan/v1", func(r chi.Router) {
		r.Route("/users/{user_id}/payments", func(r chi.Router) {
			r.Get("/", paymentsHandler.GetData)
			r.Post("/", paymentsHandler.CreateData)
			r.Patch("/{payment_id}", paymentsHandler.UpdateData)
			r.Delete("/{payment_id}", paymentsHandler.DeleteData)
			r.Get("/monthly_cost", paymentsHandler.FetchDate)
		})
		r.Get("/health", healthHandler.Check)
	})

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Logger.Error("Listen and serve failed.", zap.Error(err))
			os.Exit(1)
		}
	}()
	log.Logger.Info("start warikan-api server", zap.String("version", version))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Logger.Error("err", zap.Error(err))
	}
	log.Logger.Info("shutdown warikan-api server", zap.String("version", version))
}

func allowedOrigins() []string {
	s := os.Getenv("WARIKAN_ALLOWED_ORIGINS")
	if s == "" {
		return []string{"*"}
	}
	return strings.Split(s, ",")
}
