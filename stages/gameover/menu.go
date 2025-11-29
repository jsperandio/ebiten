package gameover

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jsperandio/ebiten/core"
	"github.com/solarlune/resolv"
)

func Init() *core.Stage {
	gameOver := core.NewStage("Game Over", getScene(), getSpace(), 0)
	gameOver.AddObjects(getStageObjects())

	return gameOver
}

func getStageObjects() []core.StageObject {
	return []core.StageObject{
		NewReturnToMenuObject(),
	}
}

func getSpace() *resolv.Space {
	return resolv.NewSpace(core.ScreenWidth, core.ScreenHeight, 1, 1)
}

func getScene() *core.Scene {
	return core.NewSceneFromImage(ebiten.NewImage(core.ScreenWidth, core.ScreenHeight), 0, 0)
}
