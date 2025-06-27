package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	player        Player
	gameMap       *Map
	showInventory bool
)

func Update() {
	// Handle input
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mouse := rl.GetMousePosition()
		tileX := int(mouse.X) / TileSize
		tileY := int(mouse.Y) / TileSize

		tile := gameMap.GetTile(tileX, tileY)
		if tile != nil && tile.IsGatherable() {
			player.TryGatherAt(tileX, tileY)
		} else {
			player.MoveToTile(tileX, tileY)
		}
	}

	if rl.IsKeyPressed(rl.KeyB) {
		showInventory = !showInventory
	}

	player.Update(rl.GetFrameTime())
}

func Draw() {
	// Draw
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	gameMap.Draw()
	player.Draw()

	if showInventory {
		player.DrawInventory(10, 10)
		player.DrawEquipment(400, 10)
	}

	if player.Gathering {
		rl.DrawText(player.GatherLabel, 10, ScreenHeight-30, 20, rl.Black)
	}

	rl.EndDrawing()
}

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Encapsulated Player Example")
	rl.SetTargetFPS(60)

	gameMap = NewMap(20, 15)
	gameMap.Generate(0.1, 0.05, 0.05)

	player = NewPlayer(
		float32((SpawnX+SpawnWidth/2)*TileSize),
		float32((SpawnY+SpawnHeight/2)*TileSize),
		gameMap,
	)

	for !rl.WindowShouldClose() {
		Update()

		Draw()
	}

	rl.CloseWindow()
}
