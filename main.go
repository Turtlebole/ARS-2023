// API-Chan
//
//	Title: API-Chan
//
//	Schemes: http
//	Version: 0.0.1
//	BasePath: /
//
//	Produces:
//	  - application/json
//
// swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ps "github.com/Turtlebole/ARS-2023/poststore"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	routerChan := mux.NewRouter()
	routerChan.StrictSlash(true)

	store, err := ps.New()
	if err != nil {
		log.Fatal(err)
	}
	server := postServer{
		store: store,
	}
	routerChan.HandleFunc("/config/", count(server.createConfigHandler)).Methods("POST")
	routerChan.HandleFunc("/configs/", count(server.getAllHandler)).Methods("GET")
	routerChan.HandleFunc("/config/{id}/", count(server.getConfigHandler)).Methods("GET")
	routerChan.HandleFunc("/config/{id}/", count(server.delConfigHandler)).Methods("DELETE")

	routerChan.HandleFunc("/group/", count(server.createGroupHandler)).Methods("POST")
	routerChan.HandleFunc("/groups/", count(server.getAllGroupsHandler)).Methods("GET")
	routerChan.HandleFunc("/group/{id}/", count(server.getGroupHandler)).Methods("GET")
	routerChan.HandleFunc("/group/{id}/config/{id}", count(server.addGroupConfigHandler)).Methods("PUT")
	routerChan.HandleFunc("/group/{id}/", count(server.delGroupHandler)).Methods("DELETE")

	routerChan.HandleFunc("/swagger.yaml", server.swaggerHandler).Methods("GET")
	routerChan.Path("/metrics").Handler(metricsHandler())
	// SwaggerUI
	optionsDevelopers := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	developerDocumentationHandler := middleware.SwaggerUI(optionsDevelopers, nil)
	routerChan.Handle("/docs", developerDocumentationHandler)
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
