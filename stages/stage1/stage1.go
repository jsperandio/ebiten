package stage1

import (
	"github.com/jsperandio/ebiten/core"
	"github.com/solarlune/resolv"
)

var MossBgPng []byte

func Init() *core.Stage {
	Stage1 := core.NewStage("Moss Cavern", getScene(), getSpace(), 0.1)
	Stage1.AddEvents(getEvents())
	Stage1.AddObjects(getStageObjects())

	return Stage1
}

func InitStage1(player *core.Player) *core.Stage {
	s1 := Init()
	s1.SetPlayer(player)

	return s1
}

func getScene() *core.Scene {
	return core.NewScene(MossBgPng, 0, 0)
}

func getSpace() *resolv.Space {
	return resolv.NewSpace(core.ScreenWidth, core.ScreenHeight, 1, 1)
}

func getEvents() []core.Event {
	return []core.Event{
		LifeDecreaseEvent(),
		HungerIncreaseEvent(),
	}
}

func getStageObjects() []core.StageObject {
	return []core.StageObject{
		NewFloorSceneObject(),
		NewCeilSceneObject(),
	}
}
