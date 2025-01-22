package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"video-api/db"
	"video-api/models"

	"github.com/gorilla/mux"
)

func main() {
	db.Init()
	err := migration()
	if err != nil {
		fmt.Println("Migration failed", err)
		return
	}
	router := mux.NewRouter()
	authenticatedRouter := router.NewRoute().Subrouter()
	authenticatedRouter.Use(authMiddleware)

	addRoutes(router, authenticatedRouter)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		fmt.Println("Starting server at port 8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-stop

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server exiting")
}

func migration() error {
	err := db.Migrate(&models.Video{})
	return err
}
