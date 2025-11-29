package core

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	MainPlayer *Player
	GameFlow   *GameFlow
	GUI        *GUI

	tick        int
	embedAssets map[string][]byte
}

// NewGame returns a new game with the given values.
// The values are copied to the new game.
// So you can safely modify the given values.
// If you want to modify the given values, create a new game.
func NewGame(p *Player, gf *GameFlow, gui *GUI, am map[string][]byte) *Game {
	return &Game{
		MainPlayer:  p,
		GameFlow:    gf,
		GUI:         gui,
		tick:        0,
		embedAssets: am,
	}
}

// NewGameFromAssets returns a new game with empty values.
func NewEmptyGame() *Game {
	return &Game{}
}

// NewDefaultGame returns a new game with default values.
// The default values are:
//   - MainPlayer: NewPlayer(25, 125)
//   - GameFlow:   NewGameFlow(MainPlayer)
//   - GUI:        NewGUI(MainPlayer)
//   - tick:       0
//   - embedAssets: nil
func NewDefaultGame(embedAssets map[string][]byte) *Game {
	game := NewEmptyGame()
	game.embedAssets = embedAssets
	game.MainPlayer = NewPlayer(embedAssets[PlayerSheetAsset], 25, 125, &game.tick)
	game.GameFlow = NewGameFlow(game.MainPlayer)
	game.GUI = NewGUI(game.MainPlayer)
	return game
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	g.tick++

	g.GameFlow.Update(g.tick)

	if g.GameFlow.Current > 0 {
		g.MainPlayer.Update(g.GameFlow.CurrentStage().Gravity)
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	g.GameFlow.Draw(screen, g.tick)

	// Only draw GUI on actual game stages (not main menu or game over)
	currentStage := g.GameFlow.CurrentStage()
	if currentStage != nil && currentStage.Name != "Main Menu" && currentStage.Name != "Game Over" {
		g.GUI.Draw(screen, g.tick)
	}

	if DEBUG {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

			cx, cy := ebiten.CursorPosition()
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: %d Y: %d", cx, cy), cx, cy)
			obj := g.GameFlow.CurrentStage().StageSpace.CheckCells(cx, cy, 1, 1)
			if obj != nil {
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Object at: %s", strings.Join(obj.Tags(), ",")), cx, cy+10)
			}
		}
	}
}

// AddStage adds a stage to the game.
// The stage is added to the game flow by given order.
func (g *Game) AddStage(s *Stage, hasPlayer bool) {
	if hasPlayer {
		s.SetPlayer(g.MainPlayer)
	}
	g.GameFlow.AddStage(s)
}

// AddAsset adds an asset to the game.
// The asset is added to the game embed assets.
// The name is the key of the asset.
func (g *Game) AddAsset(name string, data []byte) {
	g.embedAssets[name] = data
}
