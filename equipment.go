package main

type EquipmentSlot string

const (
	SlotHead   EquipmentSlot = "Head"
	SlotBody   EquipmentSlot = "Body"
	SlotLegs   EquipmentSlot = "Legs"
	SlotWeapon EquipmentSlot = "Weapon"
	SlotShield EquipmentSlot = "Shield"
)

type Equipment struct {
	Slots map[EquipmentSlot]ItemSlot
}

func NewEquipment() *Equipment {
	return &Equipment{
		Slots: map[EquipmentSlot]ItemSlot{
			SlotHead:   {},
			SlotBody:   {},
			SlotLegs:   {},
			SlotWeapon: {},
			SlotShield: {},
		},
	}
}

func (e *Equipment) Equip(slot EquipmentSlot, item ItemSlot) ItemSlot {
	e.Slots[slot] = ItemSlot{
		Name:  item.Name,
		Type:  item.Type,
		Count: 1,
	}
	// Return the remaining stack (if any)
	if item.Count > 1 {
		return ItemSlot{
			Name:  item.Name,
			Type:  item.Type,
			Count: item.Count - 1,
		}
	}
	return ItemSlot{}
}

func (e *Equipment) Unequip(slot EquipmentSlot) ItemSlot {
	item := e.Slots[slot]
	e.Slots[slot] = ItemSlot{}
	return item
}

func mapItemTypeToSlot(itemType string) EquipmentSlot {
	switch itemType {
	case "Weapon":
		return SlotWeapon
	case "Shield":
		return SlotShield
	case "Head":
		return SlotHead
	case "Body":
		return SlotBody
	case "Legs":
		return SlotLegs
	}
	return ""
}
