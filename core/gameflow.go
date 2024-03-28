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
	}
}
