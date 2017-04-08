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
    ct["dog"] = CreatureTypeInfo{"dog", "canine", 'd', 20, 30, 5, 0}

    //
    // Wolf
    //
    ct["wolf"] = CreatureTypeInfo{"wolf", "canine", 'w', 25, 35, 7, 0}

    // Set the populated flag to true.
    return true
}

