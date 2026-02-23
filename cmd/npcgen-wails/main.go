package main

import (
	"flag"
	"log"
	"os"

	"github.com/lackmus/npcgengo"
	uiwails "github.com/lackmus/npcgengo/ui/wails"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func main() {
	dataDir := flag.String("data-dir", "", "path to data directory (overrides NPCGEN_DATA env)")
	flag.Parse()

	if *dataDir == "" {
		if v := os.Getenv("NPCGEN_DATA"); v != "" {
			*dataDir = v
		}
	}

	npcGen, err := npcgengo.NewNPCGenWithDataDir(*dataDir)
	if err != nil {
		log.Fatal("failed to initialize NPCGen: ", err)
	}

	app := NewApp(npcGen.NPCListController)

	err = wails.Run(&options.App{
		Title:  "NPCGen",
		Width:  1200,
		Height: 820,
		AssetServer: &assetserver.Options{
			Assets: uiwails.Assets(),
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
