/*
 * File: types/creature_types.go
 *
 * Description: Hold type information about a creature.
 */

package types

// Structure to hold creature information
type CreatureTypeInfo struct {

	// Holds the name of the given creature.
	Name string

	// Holds the type of the given creature.
	Species string

	// Appearance of the creature.
	Ch rune

	// Health, max health, attack, and defence
	Hp    int
	MaxHp int
	Att   int
	Def   int

	// Creature class
	Class *ClassTypeInfo

	// Four primary creature attributes
	Strength     uint
	Intelligence uint
	Agility      uint
	Wisdom       uint

	// The number of steps required to heal by 1 point of health.
	Healrate uint

	// The number of steps currently walked by the creature in question.
	Healcounter uint
}

//! Function to populate details about various creature types
/*
 * @return    none
 */
func GenCreatureTypes(ct map[string]CreatureTypeInfo) bool {

	// Input validation
	if ct == nil {
		return false
	}

	//
	// Dog
	//
	ct["dog"] = CreatureTypeInfo{"dog", "canine", 'd', 20, 20, 5, 0, nil,
		20, 10, 10, 10, 10, 0}

	//
	// Wolf
	//
	ct["wolf"] = CreatureTypeInfo{"wolf", "canine", 'w', 25, 25, 7, 0, nil,
		20, 10, 10, 10, 10, 0}

	//
	// Snake
	//
	ct["snake"] = CreatureTypeInfo{"snake", "reptile", 's', 18, 18, 10, 1,
		nil, 20, 10, 10, 10, 10, 0}

	//
	// Snake
	//
	ct["spider"] = CreatureTypeInfo{"spider", "arthropod", 'x', 8, 8, 2, 2,
		nil, 20, 10, 10, 10, 10, 0}

	//
	// Goblin
	//
	ct["goblin"] = CreatureTypeInfo{"goblin", "humanoid", 'g', 22, 22, 4,
		2, nil, 20, 10, 10, 10, 10, 0}

	//
	// Orc
	//
	ct["orc"] = CreatureTypeInfo{"orc", "humanoid", 'o', 40, 40, 12, 5, nil,
		20, 10, 10, 10, 10, 0}

	// Set the populated flag to true.
	return true
}
