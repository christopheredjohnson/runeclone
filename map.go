package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Map struct {
	Width  int
	Height int
	Tiles  [][]Tile
}

func NewMap(width, height int) *Map {
	tiles := make([][]Tile, height)
	for y := range tiles {
		tiles[y] = make([]Tile, width)
		for x := range tiles[y] {
			tiles[y][x] = Tile{Type: TileGrass} // default to grass
		}
	}
	return &Map{
		Width:  width,
		Height: height,
		Tiles:  tiles,
	}
}

func (m *Map) GetTile(x, y int) *Tile {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return nil
	}
	return &m.Tiles[y][x]
}

func (m *Map) SetTile(x, y int, tileType int) {
	if tile := m.GetTile(x, y); tile != nil {
		tile.Type = tileType
	}
}

func (m *Map) Draw() {
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			tile := m.Tiles[y][x]
			pos := rl.NewRectangle(float32(x*TileSize), float32(y*TileSize), TileSize, TileSize)
			rl.DrawRectangleRec(pos, tile.Color())
			rl.DrawRectangleLinesEx(pos, 1, rl.Black)
		}
	}
}

func (m *Map) Generate(treeChance, rockChance, waterChance float64) {
	rand.Seed(time.Now().UnixNano())

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			// Skip spawn zone
			if x >= SpawnX && x < SpawnX+SpawnWidth &&
				y >= SpawnY && y < SpawnY+SpawnHeight {
				m.Tiles[y][x] = Tile{Type: TileGrass}
				continue
			}

			// Randomly assign tile type
			r := rand.Float64()
			switch {
			case r < treeChance:
				m.Tiles[y][x] = Tile{Type: TileTree}
			case r < treeChance+rockChance:
				m.Tiles[y][x] = Tile{Type: TileRock}
			case r < treeChance+rockChance+waterChance:
				m.Tiles[y][x] = Tile{Type: TileWater}
			default:
				m.Tiles[y][x] = Tile{Type: TileGrass}
			}
		}
	}
}
