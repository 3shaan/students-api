package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/3shaan/students-api/internals/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Eshan!"))
	})

	// setup server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	fmt.Println("server is running on", cfg.Address)

	serverErr := server.ListenAndServe()

	if serverErr != nil {
		log.Fatalf("Server running error, %s", serverErr.Error())
	}

}
