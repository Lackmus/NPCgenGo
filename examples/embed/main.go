package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/loader"
	"github.com/lackmus/npcgengo/service"
)

func main() {
	// Example: embed NPCGen server inside another application without the
	// server package calling os.Exit or log.Fatal itself.

	creationLoader := loader.NewJSONNPCConfigLoader("data/creation_data")
	creationSupplier, err := service.NewNPCCreationSupplier(creationLoader)
	if err != nil {
		log.Fatalf("failed to init creation supplier: %v", err)
	}

	npcService, err := service.NewNPCService(context.Background(), loader.NewJSONNPCStorage("data/npc_database"))
	if err != nil {
		log.Printf("warning: NPCService initialized with partial data: %v", err)
	}
	nc := controller.NewNPCListController(creationSupplier, npcService, "default")
	srv := controller.NewServer(nc)

	// Start server non-fatally so the host app controls lifecycle.
	go func() {
		if err := srv.Start(":8080"); err != nil {
			log.Printf("embedded server stopped: %v", err)
		}
	}()

	log.Println("embedded NPC server started on :8080")

	// Wait for termination signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down host application")

	// attempt graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("error during server shutdown: %v", err)
	} else {
		log.Println("server shutdown complete")
	}
}
