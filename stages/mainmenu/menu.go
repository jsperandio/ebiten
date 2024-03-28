package mainmenu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jsperandio/ebiten/core"
	"github.com/solarlune/resolv"
)

func Init() *core.Stage {
	MainMenu := core.NewStage("Main Menu", getScene(), getSpace(), 0)
	MainMenu.AddObjects(getStageObjects())

	return MainMenu
}

func getStageObjects() []core.StageObject {
	return []core.StageObject{
		NewStartGameObject(),
		NewExitGameObject(),
	}
}

func getSpace() *resolv.Space {
	return resolv.NewSpace(core.ScreenWidth, core.ScreenHeight, 1, 1)
}

func getScene() *core.Scene {
	return core.NewSceneFromImage(ebiten.NewImage(core.ScreenWidth, core.ScreenHeight), 0, 0)
}
