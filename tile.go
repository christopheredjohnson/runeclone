package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Tile struct {
	Type int
}

func (t Tile) IsWalkable() bool {
	return t.Type == TileGrass
}

func (t Tile) IsGatherable() bool {
	return t.Type == TileTree || t.Type == TileWater || t.Type == TileRock
}

func (t Tile) Color() rl.Color {
	switch t.Type {
	case TileGrass:
		return rl.Green
	case TileTree:
		return rl.DarkGreen
	case TileWater:
		return rl.Blue
	case TileRock:
		return rl.Gray
	default:
		return rl.Magenta // unknown
	}
}
