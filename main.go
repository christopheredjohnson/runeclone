package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	player        Player
	gameMap       *Map
	showInventory bool
	Recipes       = []Recipe{
		{
			Name: "Bronze Sword",
			Inputs: []ItemSlot{
				{Name: "Ore", Count: 2},
				{Name: "Logs", Count: 1},
			},
			Output: ItemSlot{Name: "Bronze Sword", Count: 1, Type: "Weapon"},
		},
	}
)

func Update() {
	clickedIndex, clickedUI := -1, false

	if showInventory {
		clickedIndex, clickedUI = player.CheckInventoryClick(10, 10)

		if clickedIndex >= 0 {
			item := player.Inventory.Get(clickedIndex)
			slot := mapItemTypeToSlot(item.Type)

			if slot != "" {
				swapped := player.Equipment.Unequip(slot)
				player.Inventory.Set(clickedIndex, ItemSlot{}) // remove old item
				player.Equipment.Equip(slot, item)

				if swapped.Name != "" {
					player.Inventory.Add(swapped)
				}
			}
		}

		if slot, eqClicked := player.CheckEquipmentClick(400, 10); eqClicked {
			clickedUI = true // <- mark UI was clicked
			item := player.Equipment.Unequip(slot)
			if item.Name != "" {
				player.Inventory.Add(item)
			}
		}
	}

	// Only click map if not interacting with UI
	if !clickedUI && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
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
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	gameMap.Draw()
	player.Draw()

	if showInventory {
		player.DrawInventory(10, 10)
		player.DrawEquipment(400, 10)

		hovered := player.GetHoveredInventoryIndex(10, 10)
		if hovered >= 0 {
			item := player.Inventory.Get(hovered)
			if item.Name != "" {
				mouse := rl.GetMousePosition()
				text := item.Name
				textWidth := rl.MeasureText(text, 16)
				padding := 4
				rect := rl.NewRectangle(mouse.X, mouse.Y-24, float32(textWidth+int32(padding)*2), 20)

				rl.DrawRectangleRec(rect, rl.Fade(rl.Black, 0.8))
				rl.DrawText(text, int32(mouse.X+float32(padding)), int32(mouse.Y-20), 16, rl.White)
			}
		}

		// After inventory tooltip
		hoveredEq := player.GetHoveredEquipmentSlot(400, 10)
		if hoveredEq != "" {
			item := player.Equipment.Slots[hoveredEq]
			if item.Name != "" {
				mouse := rl.GetMousePosition()
				text := item.Name
				textWidth := rl.MeasureText(text, 16)
				padding := 4
				rect := rl.NewRectangle(mouse.X, mouse.Y-24, float32(textWidth+int32(padding)*2), 20)

				rl.DrawRectangleRec(rect, rl.Fade(rl.Black, 0.8))
				rl.DrawText(text, int32(mouse.X+float32(padding)), int32(mouse.Y-20), 16, rl.White)
			}
		}

		drawCraftingUI(600, 10) // adjust x/y as needed
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
