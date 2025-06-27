package main

type ItemSlot struct {
	Name  string
	Count int
	Type  string // e.g. "Weapon", "Shield", "Body"
}

type Inventory struct {
	slots [28]ItemSlot
}

func NewInventory() *Inventory {
	return &Inventory{}
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
