package main

import (
	_ "embed"
	"flag"
	_ "image/png"

	log "github.com/sirupsen/logrus"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jsperandio/ebiten/core"
	"github.com/jsperandio/ebiten/stages"
)

// Because native embed can't handle files from another package, we have to load in main.go
// and embed the files manually.
// Maybe we can find a way to do this in future.
var (

	//go:embed slime-Sheet.png
	MainPlayerSheet []byte
	//go:embed mossbg.png
	MossBgPng []byte

	MainGame *core.Game
)

func initLogs() {
	inDebugMode := flag.Bool("d", false, "debug mode")
	flag.Parse()

	core.DEBUG = *inDebugMode

	log.SetLevel(core.GetLogLevel())
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
}

// Pass away the sprites for desired package
func initSprites() {
	core.LoadFonts()
	stages.MossBgPng = MossBgPng
}

// Initialize all needed structures/objects for the game
// is called once at the beginning of the game
func init() {
	initLogs()
	initSprites()

	MainGame = core.NewDefaultGame(map[string][]byte{
		core.PlayerSheetAsset: MainPlayerSheet,
	})
	MainGame.AddStage(stages.GetStage("Main Menu"), false)
	MainGame.AddStage(stages.GetStage("Moss Cavern"), true)
	MainGame.AddStage(stages.GetStage("Game Over"), false)
	
	// Register player with GameFlow to receive death notifications
	MainGame.MainPlayer.Register("GameFlow", MainGame.GameFlow)
}

func main() {
	ebiten.SetWindowSize(core.ScreenWidth*2, core.ScreenHeight*2)
	ebiten.SetWindowTitle("Slime Life")

	if err := ebiten.RunGame(MainGame); err != nil {
		log.Fatal(err)
	}
}
