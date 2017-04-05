/*
 * File: types/creature_types.go
 *
 * Description: Hold type information about a creature.
 *
 * TODO: fix this so that the creature types are populated on game start
 */

package types

// Global variable to hold all of the creature types.
var GlobalCreatureTypeInfoMap = make(map[string]CreatureTypeInfo)

// Global variable to check if the map has already been populated.
var GlobalCreatureTypeInfoMapIsPopulated = false

// Structure to hold creature information
type CreatureTypeInfo struct {

    // Holds the name of the given creature.
    name string

    // Holds the type of the given creature.
    species string

    // Appearance of the creature.
    ch rune

    // Health, max health, attack, and defence
    Hp    int
    MaxHp int
    Att   int
    Def   int
}

// 
// Specific creature information
//

// Dog
//GlobalCreatureTypeInfoMap["dog"]
//  = &CreatureTypeInfo{"dog", "canine", "d", 20, 30, 5, 0}
