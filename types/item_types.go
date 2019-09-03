/*
 * File: types/item_types.go
 *
 * Description: Hold type information about items.
 */

package types

// Structure to hold creature information
type ItemTypeInfo struct {

	// Holds the name of the given item.
	Name string

	// Holds the type of the given item.
	Category string

	// Appearance of the item.
	Ch rune

	// Determines whether an item can be equipped by a creature.
	Can_equip bool

	// Variable that holds whether or not an item is broken.
	Is_broken bool

	// The current condition of an item. If the max durability and the
	// current are equal, then the item is in perfect condition.
	Durability_current int
	Durability_maximum int

	// Price of the item in gold coins.
	Price_to_purchase int
	Price_to_sell     int

	// Weight of the item in grams.
	Weight int

	// Variables that determine how much an item increases (if positive)
	// or decreases (if negative) the attack / defence of a creature.
	Attack_increase  int
	Defence_increase int
}

//! Function to populate details about various items types
/*
 * @return    none
 */
func GenItemTypes(itype map[string]ItemTypeInfo) bool {

	// Input validation
	if itype == nil {
		return false
	}

	//
	// Dagger
	//
	itype["dagger"] = ItemTypeInfo{"Dagger", "blade", '%', true, false, 5,
		5, 10, 5, 10000, 1, 0}

	//
	// Sword
	//
	itype["sword"] = ItemTypeInfo{"Sword", "blade", '%', true, false, 10,
		10, 10, 5, 10000, 2, 0}

	//
	// Mace
	//
	itype["mace"] = ItemTypeInfo{"Mace", "blunt", '%', true, false, 8,
		8, 11, 3, 8000, 2, 0}

	//
	// Buckler
	//
	itype["Buckler"] = ItemTypeInfo{"Buckler", "shield", '%', true, false, 11,
		11, 20, 10, 20000, 0, 1}

	//
	// Helm
	//
	itype["Helm"] = ItemTypeInfo{"Helm", "helmet", '%', true, false, 10,
		10, 25, 8, 15000, 0, 1}

	//
	// Amulet of Defence
	//
	itype["amulet_of_defence"] = ItemTypeInfo{"Amulet of Defence", "necklace",
		'%', true, false, 20, 20, 50, 25, 5000, 0, 1}

	//
	// Leather Armour
	//
	itype["leather_armour"] = ItemTypeInfo{"Leather Armour", "armour", '%',
		true, false, 15, 15, 40, 20, 75000, 0, 1}

	//
	// Greaves
	//
	itype["greaves"] = ItemTypeInfo{"Greaves", "pants", '%', true, false,
		15, 15, 20, 10, 20000, 0, 2}

	// All of the items have been populates successfully, so return true.
	return true
}
