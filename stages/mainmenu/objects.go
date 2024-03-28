package mainmenu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jsperandio/ebiten/core"
	"github.com/solarlune/resolv"
)

type startGameObject struct {
	*core.SceneObject
	name      string
	observers map[string]core.Observer
}

func NewStartGameObject() *startGameObject {
	w := float64(core.ScreenWidth / 4)
	h := 50.0
	posx := float64(core.ScreenWidth/2) - (w / 2)
	posy := float64(core.ScreenHeight/4) - (h / 2)

	return &startGameObject{
		name:      "StartGameButton",
		observers: make(map[string]core.Observer),
		SceneObject: &core.SceneObject{
			CoreImage: *core.NewCoreTextRectImage(ebiten.NewImage(int(w), int(h)), posx, posy, core.ClrGray, "START", core.ClrWhite, core.FontSizeMedium),
			Visible:   true,
			Space:     resolv.NewObject(posx, posy, w, h, "StartGame"),
		},
	}
}

func (so *startGameObject) Name() string {
	return so.name
}

func (so *startGameObject) OnClick() {
	so.Notify("Stage", so.name)
}

func (so *startGameObject) OnRightClick() {}

func (so *startGameObject) OnMiddleClick() {}

func (so *startGameObject) Register(id string, observer core.Observer) {
	so.observers[id] = observer
}

func (so *startGameObject) Deregister(id string) {
	delete(so.observers, id)
}

func (so *startGameObject) Notify(id string, data interface{}) {
	so.observers[id].OnNotify(data)
}

type exitGameObject struct {
	name string
	*core.SceneObject
	observers map[string]core.Observer
}

func NewExitGameObject() *exitGameObject {
	w := float64(core.ScreenWidth / 4)
	h := 50.0
	posx := float64(core.ScreenWidth/2) - (w / 2)
	posy := float64(core.ScreenHeight/2) - (h / 2)

	return &exitGameObject{
		name:      "ExitGameButton",
		observers: make(map[string]core.Observer),
		SceneObject: &core.SceneObject{
			CoreImage: *core.NewCoreTextRectImage(ebiten.NewImage(int(w), int(h)), posx, posy, core.ClrGray, "EXIT", core.ClrWhite, core.FontSizeMedium),
			Visible:   true,
			Space:     resolv.NewObject(posx, posy, w, h, "ExitGame"),
		},
	}
}

func (eo *exitGameObject) Name() string {
	return eo.name
}

func (eo *exitGameObject) OnClick() {
	eo.Notify("Stage", eo.name)
}

func (eo *exitGameObject) OnRightClick() {}

func (eo *exitGameObject) OnMiddleClick() {}

func (eo *exitGameObject) Register(id string, observer core.Observer) {
	eo.observers[id] = observer
}

func (eo *exitGameObject) Deregister(id string) {
	delete(eo.observers, id)
}

func (eo *exitGameObject) Notify(id string, data interface{}) {
	eo.observers[id].OnNotify(data)
}
