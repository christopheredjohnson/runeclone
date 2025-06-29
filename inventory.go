package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ItemSlot struct {
	Name      string
	Count     int
	Type      string // e.g. "Weapon", "Shield", "Body"
	FrameRect rl.Rectangle
}

type Inventory struct {
	slots        [28]ItemSlot
	ItemsTexture rl.Texture2D
}

func NewInventory(texture rl.Texture2D) *Inventory {
	return &Inventory{
		ItemsTexture: texture,
	}
}

func (inv *Inventory) Add(slot ItemSlot) {
	for i := 0; i < len(inv.slots); i++ {
		s := &inv.slots[i]
		if s.Name == slot.Name || s.Name == "" {
			if s.Name == "" {
				*s = slot
			} else {
				s.Count += slot.Count
			}
			return
		}
	}
}

func (inv *Inventory) AddByName(name string, count int, itemType string) {
	inv.Add(ItemSlot{Name: name, Count: count, Type: itemType})
}

func (inv *Inventory) Get(index int) ItemSlot {
	if index < 0 || index >= len(inv.slots) {
		return ItemSlot{}
	}
	return inv.slots[index]
}

func (inv *Inventory) Set(index int, slot ItemSlot) {
	if index < 0 || index >= len(inv.slots) {
		return
	}
	inv.slots[index] = slot
}

func (inv *Inventory) Slots() [28]ItemSlot {
	return inv.slots
}

func inferItemType(name string) string {
	switch name {
	case "Logs":
		return "Material"
	case "Ore":
		return "Material"
	case "Fish":
		return "Food"
	}
	return "Misc"
}

func (inv *Inventory) HasItems(requirements []ItemSlot) bool {
	for _, req := range requirements {
		count := 0
		for _, slot := range inv.slots {
			if slot.Name == req.Name {
				count += slot.Count
			}
		}
		if count < req.Count {
			return false
		}
	}
	return true
}

func (inv *Inventory) ConsumeItems(requirements []ItemSlot) {
	for _, req := range requirements {
		remaining := req.Count
		for i := range inv.slots {
			slot := &inv.slots[i]
			if slot.Name == req.Name {
				if slot.Count > remaining {
					slot.Count -= remaining
					break
				} else {
					remaining -= slot.Count
					slot.Name = ""
					slot.Count = 0
				}
			}
		}
	}
}

func drawCraftingUI(x, y int) {
	boxWidth := 180
	boxHeight := 24
	padding := 6
	mouse := rl.GetMousePosition()

	for i, recipe := range Recipes {
		rect := rl.NewRectangle(float32(x), float32(y+i*(boxHeight+padding)), float32(boxWidth), float32(boxHeight))

		// Check if player can craft
		canCraft := player.Inventory.HasItems(recipe.Inputs)

		bg := rl.Gray
		if canCraft {
			bg = rl.LightGray
		}
		if rl.CheckCollisionPointRec(mouse, rect) {

			bg = rl.DarkGray
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) && canCraft {
				player.TryCraft(recipe)
			}

			tooltip := ""
			for _, input := range recipe.Inputs {
				tooltip += fmt.Sprintf("%dx %s\n", input.Count, input.Name)
			}
			rl.DrawText(tooltip, int32(mouse.X+8), int32(mouse.Y+8), 16, rl.DarkBlue)
		}

		rl.DrawRectangleRec(rect, bg)
		rl.DrawRectangleLinesEx(rect, 1, rl.Black)
		rl.DrawText(recipe.Name, int32(rect.X+6), int32(rect.Y+4), 16, rl.Black)
	}
}
