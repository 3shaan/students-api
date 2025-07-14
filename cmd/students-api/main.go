package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/3shaan/students-api/internals/config"
	"github.com/3shaan/students-api/internals/https/handlers/students"
	"github.com/3shaan/students-api/internals/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	// database config
	storage, databaseInitErr := sqlite.New(cfg)
	if databaseInitErr != nil {
		log.Panic("Database initialzed failed", databaseInitErr.Error())
	}

	slog.Info("Database Initialize success")

	//setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", students.New(storage))
	router.HandleFunc("GET /api/students", students.GetAll(storage))
	router.HandleFunc("GET /api/students/{id}", students.GetStudentById(storage))
	router.HandleFunc("DELETE /api/students/{id}", students.DeleteStudentById(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	fmt.Println("server is running on", cfg.Address)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		serverErr := server.ListenAndServe()

		if serverErr != nil {
			log.Fatalf("Server running error, %s", serverErr.Error())
		}

	}()

	<-done

	slog.Info("Shutting down the sever...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdownErr := server.Shutdown(ctx)
	if shutdownErr != nil {
		slog.Error("Failed to shutdown", slog.String("error", shutdownErr.Error()))
	}

}
