package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sahil-dharwani/students-api/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//setup any custom logger if any
	//database setup
	//setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to students api"))
	})

	//setup server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("server Started...", slog.String("info", cfg.Address))

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server %s", err.Error())
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shut down server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown successfully")
}
