package main

import (
	"Ex2_Week3/cmd"
	"Ex2_Week3/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	err := db.DbConnection()
	if err!= nil {
        log.Fatal(err)
    }

	router := mux.NewRouter()
	cmd.SetupRoutes(router)
	
	srv := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	log.Println("Starting server on port 8080")
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("ListenAndServe error: %v", err)
    }
}