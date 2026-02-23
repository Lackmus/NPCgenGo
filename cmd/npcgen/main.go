package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lackmus/npcgengo"
	"github.com/lackmus/npcgengo/cmd/npcgen/handlers"
)

func main() {
	dataDir := flag.String("data-dir", "", "path to data directory (overrides NPCGEN_DATA env)")
	flag.Parse()

	// env var fallback
	if *dataDir == "" {
		if v := os.Getenv("NPCGEN_DATA"); v != "" {
			*dataDir = v
		}
	}

	//npcGen, err := npcgengo.NewNPCGenWithDataDir(*dataDir)
	npcGen, err := npcgengo.NewNPCGen()
	if err != nil {
		log.Fatal("failed to initialize NPCGen:", err)
	}

	srv := handlers.NewServer(npcGen.NPCListController)

	go func() {
		if err := srv.Start(":8080"); err != nil {
			log.Printf("embedded server stopped: %v", err)
		}
	}()

	log.Println("embedded NPC server started on :8080")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down host application")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("error during server shutdown: %v", err)
	} else {
		log.Println("server shutdown complete")
	}
}
