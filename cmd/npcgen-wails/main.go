package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/lackmus/npcgengo"
	uiwails "github.com/lackmus/npcgengo/ui/wails"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func main() {
	if err := syncSharedUICore(); err != nil {
		log.Fatal("failed to sync shared ui core:", err)
	}

	dataDir := flag.String("data-dir", "", "path to data directory (overrides NPCGEN_DATA env)")
	flag.Parse()

	npcGen, err := npcgengo.NewNPCGenWithDataDir(*dataDir)
	if err != nil {
		log.Fatal("failed to initialize NPCGen: ", err)
	}

	app := NewWailsAPI(npcGen.NPCListController)

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

func syncSharedUICore() error {
	sourcePath := filepath.Clean(filepath.Join("..", "..", "ui", "shared", "npc-ui-core.js"))
	targetPath := filepath.Clean(filepath.Join("..", "..", "ui", "wails", "dist", "npc-ui-core.js"))

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return err
	}

	return targetFile.Sync()
}
