package stages

import (
	"github.com/jsperandio/ebiten/core"
	"github.com/jsperandio/ebiten/stages/mainmenu"
	"github.com/jsperandio/ebiten/stages/stage1"
)

var MossBgPng []byte

func InitStage(stageName string, player *core.Player) *core.Stage {
	switch stageName {
	case "Main Menu":
		return mainmenu.Init()
	case "Moss Cavern":
		stage1.MossBgPng = MossBgPng
		return stage1.InitStage1(player)
	}
	return nil
}

func GetStage(name string) *core.Stage {
	switch name {
	case "Main Menu":
		return mainmenu.Init()
	case "Moss Cavern":
		stage1.MossBgPng = MossBgPng
		return stage1.Init()
	}
	return nil
}
