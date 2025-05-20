package main

import (
	"context"
	Config "firstproject/internal/config"
	"firstproject/internal/http/student"
	// "firstproject/internal/storage"
	"firstproject/internal/storage/sqlite"

	// "fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := Config.MustLoad()
	storage,err:=sqlite.New(cfg)
	if err!=nil{
		log.Fatal(err)
	}
	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
 

	// Create a channel to listen for termination signals
	done := make(chan os.Signal, 1)

	// Setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /api/students", student.New(storage))
    //  router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	// router.HandleFunc("GET /api/students", student.GetList(storage))
	// Setup HTTP server
	server := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	// Notify `done` channel if interrupt or termination signal is caught
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in a separate goroutine
	go func() {
		slog.Info("Starting server", slog.String("address", cfg.HTTPServer.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait here until a signal is received
	<-done

	slog.Info("shutting down server...")

	// Gracefully shutdown the server with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
		return
	}

	slog.Info("server shutdown successfully")
}
