package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var gatherSettings = map[int]struct {
	Label string
	Item  string
	Time  float32
}{
	TileTree:  {"Chopping...", "Logs", 2.0},
	TileRock:  {"Mining...", "Ore", 2.5},
	TileWater: {"Fishing...", "Fish", 3.0},
}

type Player struct {
	Pos           rl.Vector2
	Size          rl.Vector2
	Speed         float32
	Color         rl.Color
	Target        rl.Vector2
	Map           *Map
	Path          []Point
	Inventory     *Inventory
	Gathering     bool
	GatherTarget  Point
	GatherLabel   string
	GatherTimer   float32
	GatherItem    string
	PendingGather *Point
	Equipment     *Equipment
	Texture       rl.Texture2D
	Health        int
	MaxHealth     int
}

func NewPlayer(x, y float32, m *Map, texture rl.Texture2D, itemTexture rl.Texture2D) Player {
	return Player{
		Pos:       rl.NewVector2(x, y),
		Size:      rl.NewVector2(32, 32),
		Speed:     180,
		Color:     rl.Brown,
		Target:    rl.NewVector2(x, y),
		Map:       m,
		Inventory: NewInventory(itemTexture),
		Equipment: NewEquipment(),
		Texture:   texture,
		Health:    100,
		MaxHealth: 100,
	}
}

func (p *Player) MoveToTile(tileX, tileY int) {
	start := Point{
		X: int(p.Pos.X+p.Size.X/2) / TileSize,
		Y: int(p.Pos.Y+p.Size.Y/2) / TileSize,
	}
	goal := Point{tileX, tileY}
	path := FindPath(start, goal, p.Map)

	if len(path) == 0 {
		fmt.Println("No valid path to target:", goal)
		p.PendingGather = nil
		return
	}

	p.Path = path
}

func (p *Player) FindAdjacentWalkable(target Point) *Point {
	dirs := []Point{
		{X: 0, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0},
	}

	start := Point{
		X: int(p.Pos.X+p.Size.X/2) / TileSize,
		Y: int(p.Pos.Y+p.Size.Y/2) / TileSize,
	}

	var bestPath []Point
	var bestAdj *Point

	for _, d := range dirs {
		adj := Point{X: target.X + d.X, Y: target.Y + d.Y}
		tile := p.Map.GetTile(adj.X, adj.Y)

		if tile != nil && tile.IsWalkable() {
			path := FindPath(start, adj, p.Map)
			if len(path) > 0 && (bestPath == nil || len(path) < len(bestPath)) {
				bestPath = path
				bestAdj = &adj
			}
		}
	}

	return bestAdj
}

func (p *Player) Update(dt float32) {
	if len(p.Path) > 0 {
		next := p.Path[0]
		centerX := float32(next.X*TileSize + TileSize/2)
		centerY := float32(next.Y*TileSize + TileSize/2)
		target := rl.NewVector2(centerX-p.Size.X/2, centerY-p.Size.Y/2)

		dir := rl.Vector2Subtract(target, p.Pos)
		if rl.Vector2Length(dir) < 2 {
			p.Pos = target
			p.Path = p.Path[1:]
		} else {
			dir = rl.Vector2Normalize(dir)
			dir = rl.Vector2Scale(dir, p.Speed*dt)
			p.Pos = rl.Vector2Add(p.Pos, dir)
		}
	}

	if p.Gathering {
		p.GatherTimer -= dt
		if p.GatherTimer <= 0 {
			p.FinishGather()
		}
	}

	if p.PendingGather != nil && len(p.Path) == 0 {
		playerTileX := int(p.Pos.X+p.Size.X/2) / TileSize
		playerTileY := int(p.Pos.Y+p.Size.Y/2) / TileSize

		dx := p.PendingGather.X - playerTileX
		dy := p.PendingGather.Y - playerTileY

		fmt.Printf("Arrived near gather target? Player at (%d,%d), target (%d,%d)\n", playerTileX, playerTileY, p.PendingGather.X, p.PendingGather.Y)

		if abs(dx)+abs(dy) <= 1 {
			p.startGather(p.PendingGather.X, p.PendingGather.Y)
		} else {
			fmt.Println("Target not adjacent after walking — skipping gather")
		}
		p.PendingGather = nil
	}
}

func (p *Player) Draw() {
	for _, step := range p.Path {
		center := rl.NewVector2(float32(step.X*TileSize+TileSize/2), float32(step.Y*TileSize+TileSize/2))
		rl.DrawCircleV(center, 2, rl.Red)
	}

	source := rl.Rectangle{
		X:      0,
		Y:      32,
		Width:  TileSize,
		Height: TileSize,
	}

	rl.DrawTextureRec(p.Texture, source, p.Pos, rl.White)
}

func (p *Player) DrawInventory(x, y int) {
	slots := p.Inventory.Slots()
	boxSize := 40
	columns := 7

	for i, slot := range slots {
		cx := x + (i%columns)*(boxSize+4)
		cy := y + (i/columns)*(boxSize+4)
		rect := rl.NewRectangle(float32(cx), float32(cy), float32(boxSize), float32(boxSize))

		rl.DrawRectangleRec(rect, rl.LightGray)
		rl.DrawRectangleLinesEx(rect, 1, rl.DarkGray)

		if slot.Name != "" {
			rl.DrawText(slot.Name[:1], int32(cx+4), int32(cy+2), 20, rl.Black)

			rl.DrawTextureRec(player.Inventory.ItemsTexture, slot.FrameRect, rl.NewVector2(float32(cx), float32(cy)), rl.White)
			rl.DrawText(fmt.Sprintf("%d", slot.Count), int32(cx+4), int32(cy+20), 16, rl.DarkBlue)
		}
	}
}

func (p *Player) TryGatherAt(tileX, tileY int) {
	tile := p.Map.GetTile(tileX, tileY)
	if tile == nil || !tile.IsGatherable() {
		fmt.Println("Tile not gatherable")
		return
	}

	playerTileX := int(p.Pos.X+p.Size.X/2) / TileSize
	playerTileY := int(p.Pos.Y+p.Size.Y/2) / TileSize

	dx := tileX - playerTileX
	dy := tileY - playerTileY

	if abs(dx)+abs(dy) <= 1 {
		fmt.Println("Gathering immediately at", tileX, tileY)
		p.startGather(tileX, tileY)
	} else {
		adj := p.FindAdjacentWalkable(Point{tileX, tileY})
		if adj == nil {
			fmt.Println("No adjacent walkable tile to gather target")
			return
		}
		p.MoveToTile(adj.X, adj.Y)
		p.PendingGather = &Point{tileX, tileY}
	}
}

func (p *Player) startGather(tileX, tileY int) {
	p.Gathering = true
	p.GatherTarget = Point{tileX, tileY}

	tile := p.Map.GetTile(tileX, tileY)

	if settings, ok := gatherSettings[tile.Type]; ok {
		p.Gathering = true
		p.GatherTarget = Point{tileX, tileY}
		p.GatherLabel = settings.Label
		p.GatherItem = settings.Item
		p.GatherTimer = settings.Time
	} else {
		fmt.Println("Invalid gather target")
	}
}

func (p *Player) FinishGather() {
	tile := p.Map.GetTile(p.GatherTarget.X, p.GatherTarget.Y)
	if tile == nil {
		p.Gathering = false
		return
	}

	switch tile.Type {
	case TileTree, TileRock:
		tile.Type = TileGrass
	}

	// Add item to inventory
	p.Inventory.Add(ItemSlot{
		Name:  p.GatherItem,
		Count: 1,
		Type:  inferItemType(p.GatherItem),
	})

	p.Gathering = false
	p.GatherLabel = ""
}

func (p *Player) DrawEquipment(x, y int) {
	boxSize := 40
	padding := 4

	slotOrder := []EquipmentSlot{
		SlotHead,
		SlotBody,
		SlotLegs,
		SlotWeapon,
		SlotShield,
	}

	for i, slot := range slotOrder {
		cx := x
		cy := y + i*(boxSize+padding)
		rect := rl.NewRectangle(float32(cx), float32(cy), float32(boxSize), float32(boxSize))

		// Draw slot box
		rl.DrawRectangleRec(rect, rl.LightGray)
		rl.DrawRectangleLinesEx(rect, 1, rl.DarkGray)

		// Draw slot label
		rl.DrawText(string(slot), int32(cx+boxSize+6), int32(cy+12), 16, rl.Black)

		item := p.Equipment.Slots[slot]
		if item.Name != "" {

			rl.DrawTextureRec(player.Inventory.ItemsTexture, item.FrameRect, rl.NewVector2(float32(cx), float32(cy)), rl.White)
			rl.DrawText(item.Name[:1], int32(cx+4), int32(cy+2), 20, rl.Black)
		}
	}
}

func (p *Player) CheckInventoryClick(x, y int) (clickedIndex int, clickedUI bool) {
	slots := p.Inventory.Slots()
	boxSize := 40
	columns := 7
	mouse := rl.GetMousePosition()

	clickedIndex = -1

	for i := range slots {
		cx := x + (i%columns)*(boxSize+4)
		cy := y + (i/columns)*(boxSize+4)
		rect := rl.NewRectangle(float32(cx), float32(cy), float32(boxSize), float32(boxSize))

		if rl.CheckCollisionPointRec(mouse, rect) {
			clickedUI = true
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				clickedIndex = i
			}
		}
	}

	return clickedIndex, clickedUI
}

func (p *Player) GetHoveredInventoryIndex(x, y int) int {
	slots := p.Inventory.Slots()
	boxSize := 40
	columns := 7
	mouse := rl.GetMousePosition()

	for i := range slots {
		cx := x + (i%columns)*(boxSize+4)
		cy := y + (i/columns)*(boxSize+4)
		rect := rl.NewRectangle(float32(cx), float32(cy), float32(boxSize), float32(boxSize))

		if rl.CheckCollisionPointRec(mouse, rect) {
			return i
		}
	}

	return -1
}

func (p *Player) CheckEquipmentClick(x, y int) (EquipmentSlot, bool) {
	boxSize := 40
	padding := 4

	slotOrder := []EquipmentSlot{
		SlotHead,
		SlotBody,
		SlotLegs,
		SlotWeapon,
		SlotShield,
	}

	mouse := rl.GetMousePosition()

	for i, slot := range slotOrder {
		cx := x
		cy := y + i*(boxSize+padding)
		rect := rl.NewRectangle(float32(cx), float32(cy), float32(boxSize), float32(boxSize))

		if rl.CheckCollisionPointRec(mouse, rect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			return slot, true
		}
	}

	return "", false
}

func (p *Player) GetHoveredEquipmentSlot(x, y int) EquipmentSlot {
	boxSize := 40
	padding := 4

	slotOrder := []EquipmentSlot{
		SlotHead,
		SlotBody,
		SlotLegs,
		SlotWeapon,
		SlotShield,
	}

	mouse := rl.GetMousePosition()

	for i, slot := range slotOrder {
		cx := x
		cy := y + i*(boxSize+padding)
		rect := rl.NewRectangle(float32(cx), float32(cy), float32(boxSize), float32(boxSize))
		if rl.CheckCollisionPointRec(mouse, rect) {
			return slot
		}
	}

	return ""
}

func (p *Player) TryCraft(recipe Recipe) bool {
	if p.Inventory.HasItems(recipe.Inputs) {
		p.Inventory.ConsumeItems(recipe.Inputs)
		p.Inventory.Add(recipe.Output)
		return true
	}
	return false
}
