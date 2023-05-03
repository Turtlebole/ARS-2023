package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	routerChan := mux.NewRouter()
	routerChan.StrictSlash(true)

	server := configServer{
		data:      map[string]*Config{},
		groupData: map[string]*Group{},
	}
	routerChan.HandleFunc("/config/", server.createConfigHandler).Methods("POST")
	routerChan.HandleFunc("/configs/", server.getAllHandler).Methods("GET")
	routerChan.HandleFunc("/config/{id}/", server.getConfigHandler).Methods("GET")
	routerChan.HandleFunc("/config/{id}/", server.delConfigHandler).Methods("DELETE")
	routerChan.HandleFunc("/group/", server.createGroupHandler).Methods("POST")
	routerChan.HandleFunc("/groups/", server.getAllHandler).Methods("GET")
	routerChan.HandleFunc("/group/{id}/", server.getGroupHandler).Methods("GET")
	routerChan.HandleFunc("/group/{groupId}/{id}/", server.addGroupConfig).Methods("PUT")

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
