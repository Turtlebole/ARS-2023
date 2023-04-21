package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/gorilla/mux"
)

type Service struct {
	Data map[string]*Config
}

func main() {
	fmt.Println("Hello World!")
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	routerChan := mux.NewRouter()
	routerChan.StrictSlash(true)

	server := configServer{
		data: map[string]*Config{},
	}
	routerChan.HandleFunc("/config/", server.createPostHandler).Methods("POST")
	routerChan.HandleFunc("/configs/", server.getAllHandler).Methods("GET")
	routerChan.HandleFunc("/config/{id}/", server.getPostHandler).Methods("GET")
	routerChan.HandleFunc("/config/{id}/", server.delPostHandler).Methods("DELETE")

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: routerChan}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
