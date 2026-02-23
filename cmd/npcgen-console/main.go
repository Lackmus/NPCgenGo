package main

import (
	"flag"
	"log"
	"os"

	"github.com/lackmus/npcgengo"
)

func main() {
	dataDir := flag.String("data-dir", "", "path to data directory (overrides NPCGEN_DATA env)")
	locationID := flag.String("location-id", "default", "location id for console NPC view")
	flag.Parse()

	if *dataDir == "" {
		if v := os.Getenv("NPCGEN_DATA"); v != "" {
			*dataDir = v
		}
	}

	npcGen, err := npcgengo.NewNPCGenWithDataDir(*dataDir)
	if err != nil {
		log.Fatal("failed to initialize NPCGen:", err)
	}

	npcGen.InitNPCListView(*locationID)
}
