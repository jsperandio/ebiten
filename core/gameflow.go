package core

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameFlow struct {
	Player     *Player
	gameStages []*Stage
	Current    int
}

func NewGameFlow(p *Player) *GameFlow {
	return &GameFlow{
		Player:     p,
		gameStages: []*Stage{},
		Current:    0,
	}
}

func (gf *GameFlow) Update(tick int) {
	gf.CurrentStage().Update(tick)
}

func (gf *GameFlow) Draw(screen *ebiten.Image, tick int) {
	gf.CurrentStage().Draw(screen)
}

func (gf *GameFlow) AddStage(stage *Stage) {
	stage.Register("GameFlow", gf)
	gf.gameStages = append(gf.gameStages, stage)
}

func (gf *GameFlow) Next() {
	if gf.Current >= len(gf.gameStages) {
		return
	}
	gf.Current++
}

func (gf *GameFlow) CurrentStage() *Stage {
	return gf.gameStages[gf.Current]
}

func (gf *GameFlow) OnNotify(value interface{}) {
	switch value {
	case "StartGameButton":
		gf.Next()
	case "ExitGameButton":
		os.Exit(-1)
	case "PlayerDied":
		gf.GoToGameOver()
	case "ReturnToMenuButton":
		gf.GoToMainMenu()
	}
}

func (gf *GameFlow) GoToStage(stageName string) {
	for i, stage := range gf.gameStages {
		if stage.Name == stageName {
			gf.Current = i
			return
		}
	}
}

func (gf *GameFlow) GoToMainMenu() {
	gf.Current = 0
	// Reset player when returning to main menu
	if gf.Player != nil {
		gf.Player.LifePoints = 100
		gf.Player.Hunger = Full
		// Reset player position to initial spawn
		gf.Player.updatePlayerPosition(25, 125)
		gf.Player.velocityX = 0
		gf.Player.velocityY = 0
		gf.Player.FaceDirection = Left
		// Clear any draw events
		gf.Player.DrawEvents = []*DrawableEvent{}
		// Restore player image from sprite sheet
		gf.Player.RestoreImage()
		// Reset animation frame
		gf.Player.FrameOX = 0
		gf.Player.FrameOY = 0
		// Notify GUI to reset health bar
		gf.Player.Notify("GUI", 100)
	}
}

func (gf *GameFlow) GoToGameOver() {
	gf.GoToStage("Game Over")
}
