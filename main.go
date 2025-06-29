package main

import (
	"fmt"
	"math/rand"

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
	enemies          []Enemy
	inCombat         bool
	currentEnemy     *Enemy
	playerTurn       bool
	combatTimer      float32
	combatInterval   = float32(1.0) // seconds per turn
	lootMessage      string
	lootMessageTimer float32
)

func Update() {

	if lootMessageTimer > 0 {
		lootMessageTimer -= rl.GetFrameTime()
	}

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
			if rl.Vector2Distance(player.Pos, enemies[i].Pos) < float32(TileSize) {
				inCombat = true
				currentEnemy = &enemies[i]
				playerTurn = true
				combatTimer = combatInterval
				break
			}
		}
	}

	if inCombat && currentEnemy != nil {
		combatTimer -= rl.GetFrameTime()
		if combatTimer <= 0 {
			if playerTurn {
				currentEnemy.Health -= 10
				fmt.Println("Player hits", currentEnemy.Name, "for 10 damage")

				if currentEnemy.Health <= 0 {
					fmt.Println(currentEnemy.Name, "is defeated!")

					for _, loot := range currentEnemy.LootTable {
						if rand.Float32() <= loot.Chance {
							player.Inventory.Add(loot.Item)
							fmt.Printf("Looted: %s x%d\n", loot.Item.Name, loot.Item.Count)
							lootMessage = fmt.Sprintf("You looted %s x%d", loot.Item.Name, loot.Item.Count)
							lootMessageTimer = 2.0
						}
					}

					inCombat = false
					currentEnemy = nil
					return
				}
			} else {
				player.Health -= 5
				fmt.Println(currentEnemy.Name, "hits Player for 5 damage")

				if player.Health <= 0 {
					fmt.Println("You died!")
					inCombat = false
					currentEnemy = nil
					return
				}
			}
			playerTurn = !playerTurn
			combatTimer = combatInterval
		}
	}

	if inCombat && currentEnemy != nil && rl.Vector2Distance(player.Pos, currentEnemy.Pos) > float32(TileSize*2) {
		fmt.Println("You escaped combat.")
		inCombat = false
		currentEnemy = nil
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
		status := "Fighting " + currentEnemy.Name
		turn := "Turn: Player"
		if !playerTurn {
			turn = "Turn: " + currentEnemy.Name
		}
		rl.DrawText(status, 10, ScreenHeight-60, 20, rl.Red)
		rl.DrawText(turn, 10, ScreenHeight-40, 20, rl.DarkGray)
	}

	rl.DrawText(fmt.Sprintf("Player HP: %d", player.Health), 10, 10, 20, rl.Black)

	if lootMessageTimer > 0 {
		rl.DrawText(lootMessage, 10, ScreenHeight-90, 20, rl.DarkGreen)
	}

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
		LootTable: []LootEntry{
			{
				Item: ItemSlot{
					Name:      "Club",
					Count:     1,
					Type:      "Weapon",
					FrameRect: rl.NewRectangle(0, 256, TileSize, TileSize),
				},
				Chance: 0.3, // 30% chance
			},
			{
				Item: ItemSlot{
					Name:  "Coins",
					Count: 5,
					Type:  "Misc",
				},
				Chance: 0.7, // 70% chance
			},
		},
	})

	for !rl.WindowShouldClose() {
		Update()
		Draw(tilemap)
	}

	rl.CloseWindow()
}
