package main

import rl "github.com/gen2brain/raylib-go/raylib"

type LootEntry struct {
	Item   ItemSlot
	Chance float32 // 0.0 to 1.0
}

type Enemy struct {
	Pos       rl.Vector2
	Texture   rl.Texture2D
	Frame     rl.Rectangle
	Health    int
	MaxHealth int
	Name      string
	LootTable []LootEntry
}

func (e *Enemy) Draw() {
	rl.DrawTextureRec(e.Texture, e.Frame, e.Pos, rl.White)

	// Optional health bar
	barWidth := TileSize
	rl.DrawRectangle(int32(e.Pos.X), int32(e.Pos.Y)-6, int32(barWidth), 4, rl.Red)
	rl.DrawRectangle(int32(e.Pos.X), int32(e.Pos.Y)-6, int32(float32(barWidth)*float32(e.Health)/float32(e.MaxHealth)), 4, rl.Green)
}
