package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	_ "os"
	"path/filepath"
	"strconv"
	_ "strconv"

	"github.com/Izhdenev/Final-Project/pkg/api"
	"github.com/Izhdenev/Final-Project/pkg/db"
	"github.com/Izhdenev/Final-Project/tests"
)

func main() {
	dbfile := filepath.Join("tests", tests.DBFile)
	port := tests.Port
	webDir := tests.WebDir

	if v := os.Getenv("TODO_DBFILE"); v != "" {
		dbfile = v
	}
	if v := os.Getenv("TODO_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}

	if err := db.Init(dbfile); err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting file server on %s, serving from %q", addr, webDir)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
