package main

import rl "github.com/gen2brain/raylib-go/raylib"

var tileRects = map[int]rl.Rectangle{
	TileGrass: {X: 0, Y: 416, Width: TileSize, Height: TileSize},
	TileTree:  {X: 0, Y: 800, Width: TileSize, Height: TileSize},
	TileRock:  {X: 32, Y: 576, Width: TileSize, Height: TileSize},
}

type Tile struct {
	Type int
}

func (t Tile) IsWalkable() bool {
	return t.Type == TileGrass
}

func (t Tile) IsGatherable() bool {
	return t.Type == TileTree || t.Type == TileWater || t.Type == TileRock
}

func (t Tile) Draw(texture rl.Texture2D, x, y int32) {

	switch t.Type {
	case TileRock,
		TileTree:
		src := tileRects[TileGrass] // `tileRects` is your atlas frame map
		dest := rl.Rectangle{
			X:      float32(x * TileSize),
			Y:      float32(y * TileSize),
			Width:  TileSize,
			Height: TileSize,
		}
		rl.DrawTexturePro(texture, src, dest, rl.Vector2{X: 0, Y: 0}, 0, rl.White)

		src = tileRects[t.Type] // `tileRects` is your atlas frame map
		dest = rl.Rectangle{
			X:      float32(x * TileSize),
			Y:      float32(y * TileSize),
			Width:  TileSize,
			Height: TileSize,
		}
		rl.DrawTexturePro(texture, src, dest, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
	default:
		src := tileRects[t.Type] // `tileRects` is your atlas frame map
		dest := rl.Rectangle{
			X:      float32(x * TileSize),
			Y:      float32(y * TileSize),
			Width:  TileSize,
			Height: TileSize,
		}
		rl.DrawTexturePro(texture, src, dest, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
	}

}
