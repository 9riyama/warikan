package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/warikan/api/domain/repository"
	handler "github.com/warikan/api/handler/rest"
	"github.com/warikan/api/usecase"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	port   string
	driver string
)

func init() {
	flag.StringVar(&port, "port", ":9292", "tcp host:port to connect")
	flag.StringVar(&driver, "driver", "postgres", "db driver: postgres")
	flag.Parse()
}

func main() {
	dsn := os.Getenv("WARIKAN_DB_DSN")
	sqlDB, err := sqlx.Connect(driver, dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer sqlDB.Close()

	sqlDB.SetMaxOpenConns(5)

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

	healthRepositry := repository.NewHealthRepository(sqlDB)
	healthUseCase := usecase.NewHealthUseCase(healthRepositry)
	healthHandler := handler.NewHealthHandler(healthUseCase)

	paymentRepository := repository.NewPaymentRepository(sqlDB)
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
		r.Get("/health", healthHandler.Health)
	})

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Listen and serve failed.%v\n", err)
			os.Exit(1)
		}
	}()
	log.Println("start warikan-api server.")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to srv.Shutdown  %v\n", err)
	}
	log.Println("shutdown warikan-api server.")
}

func allowedOrigins() []string {
	s := os.Getenv("WARIKAN_ALLOWED_ORIGINS")
	if s == "" {
		return []string{"*"}
	}
	return strings.Split(s, ",")
}
