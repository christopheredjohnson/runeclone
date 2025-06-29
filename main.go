package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	player        Player
	gameMap       *Map
	showInventory bool
	Recipes       = []Recipe{
		{
			Name: "Sword",
			Inputs: []ItemSlot{
				{Name: "Ore", Count: 2},
				{Name: "Logs", Count: 1},
			},
			Output: ItemSlot{Name: "Sword", Count: 1, Type: "Weapon"},
		},
	}
	enemies      []Enemy
	inCombat     bool
	currentEnemy *Enemy
	playerTurn   bool
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

				remaining := player.Equipment.Equip(slot, item)
				player.Inventory.Set(clickedIndex, remaining) // either empty or rest of stack

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

	if !inCombat {
		for i := range enemies {
			dist := rl.Vector2Distance(player.Pos, enemies[i].Pos)
			if dist < float32(TileSize) {
				currentEnemy = &enemies[i]
				inCombat = true
				playerTurn = true
				break
			}
		}
	}

	if inCombat && currentEnemy != nil {
		if playerTurn {
			if rl.IsKeyPressed(rl.KeyF) { // Player chooses to attack
				currentEnemy.Health -= 10
				fmt.Println("Player attacks:", currentEnemy.Name)

				if currentEnemy.Health <= 0 {
					fmt.Println("Enemy defeated!")
					inCombat = false
					currentEnemy = nil
				} else {
					playerTurn = false
				}
			}
		} else {
			// Enemy turn
			player.Health -= 5
			fmt.Println("Enemy attacks player!")

			if player.Health <= 0 {
				fmt.Println("You died!")
				// Optional: reset game or show death screen
			}
			playerTurn = true
		}
	}
	player.Update(rl.GetFrameTime())

	alive := enemies[:0]
	for _, e := range enemies {
		if e.Health > 0 {
			alive = append(alive, e)
		}
	}
	enemies = alive
}

func Draw(tilemap rl.Texture2D) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	gameMap.Draw(tilemap)
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

	for _, enemy := range enemies {
		enemy.Draw()
	}

	if inCombat && currentEnemy != nil {
		msg := "Combat with " + currentEnemy.Name
		turn := "Your turn: Press F to attack"
		if !playerTurn {
			turn = "Enemy is attacking..."
		}

		rl.DrawText(msg, 10, ScreenHeight-60, 20, rl.Red)
		rl.DrawText(turn, 10, ScreenHeight-40, 20, rl.DarkGray)
	}

	rl.DrawText(fmt.Sprintf("HP: %d/%d", player.Health, player.MaxHealth), 10, 10, 20, rl.Black)

	rl.EndDrawing()
}

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Encapsulated Player Example")
	rl.SetTargetFPS(60)

	tilemap := rl.LoadTexture("assets/tiles.png")
	defer rl.UnloadTexture(tilemap)

	characterTilemap := rl.LoadTexture("assets/rogues.png")
	defer rl.UnloadTexture(characterTilemap)

	itemTexture := rl.LoadTexture("assets/items.png")
	defer rl.UnloadTexture(itemTexture)

	enemyTex := rl.LoadTexture("assets/monsters.png")
	defer rl.UnloadTexture(enemyTex)

	gameMap = NewMap(20, 15)
	gameMap.Generate(0.1, 0.05, 0.05)

	player = NewPlayer(
		float32((SpawnX+SpawnWidth/2)*TileSize),
		float32((SpawnY+SpawnHeight/2)*TileSize),
		gameMap,
		characterTilemap,
		itemTexture,
	)

	enemies = append(enemies, Enemy{
		Pos:       rl.NewVector2(100, 100),
		Health:    50,
		MaxHealth: 50,
		Name:      "Slime",
		Texture:   enemyTex,
		Frame:     rl.NewRectangle(0, 64, TileSize, TileSize),
	})

	for !rl.WindowShouldClose() {
		Update()
		Draw(tilemap)
	}

	rl.CloseWindow()
}
