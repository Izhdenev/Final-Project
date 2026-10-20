package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Izhdenev/Final-Project/pkg/api"
	"github.com/Izhdenev/Final-Project/pkg/db"
)

const (
	defaultPort   = 7540
	defaultDBFile = "tests/scheduler.db"
	defaultWebDir = "web"
)

func main() {
	port := defaultPort
	dbfile := defaultDBFile
	webDir := defaultWebDir

	if v := os.Getenv("TODO_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}
	if v := os.Getenv("TODO_DBFILE"); v != "" {
		dbfile = v
	}
	if v := os.Getenv("TODO_WEBDIR"); v != "" {
		webDir = v
	}

	password := os.Getenv("TODO_PASSWORD")

	if err := db.Init(dbfile); err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close DB: %v", err)
		}
	}()

	api.Init(password)
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: nil, // DefaultServeMux
	}

	go func() {
		log.Printf("Starting server on %s (db=%s, web=%s)", srv.Addr, dbfile, webDir)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
