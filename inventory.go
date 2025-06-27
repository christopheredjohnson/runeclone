package main

type ItemSlot struct {
	Name  string
	Count int
}

type Inventory struct {
	slots [28]ItemSlot
}

func NewInventory() *Inventory {
	return &Inventory{}
}

// Adds item to first matching or empty slot
func (inv *Inventory) Add(item string, amount int) {
	for i := 0; i < len(inv.slots); i++ {
		slot := &inv.slots[i]
		if slot.Name == item || slot.Name == "" {
			slot.Name = item
			slot.Count += amount
			return
		}
	}
}

// Get slots (for display)
func (inv *Inventory) Slots() [28]ItemSlot {
	return inv.slots
}
