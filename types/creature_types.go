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
    // TODO: several of these monster attributes are accidently healing
    //       the protagonist; also consider properly balancing the
    //       attributes of some of these monster to make combat more
    //       challenging
    //

    //
    // Dog
    //
    ct["dog"] = CreatureTypeInfo{"dog", "canine", 'd', 20, 20, 5, 0}

    //
    // Wolf
    //
    ct["wolf"] = CreatureTypeInfo{"wolf", "canine", 'w', 25, 25, 7, 0}

    //
    // Snake
    //
    ct["snake"] = CreatureTypeInfo{"snake", "reptile", 's', 18, 18, 10, 1}

    //
    // Snake
    //
    ct["spider"] = CreatureTypeInfo{"spider", "arthropod", 'x', 8, 8, 2, 2}

    //
    // Goblin
    //
    ct["goblin"] = CreatureTypeInfo{"goblin", "humanoid", 'g', 22, 22, 4, 2}

    //
    // Orc
    //
    ct["orc"] = CreatureTypeInfo{"orc", "humanoid", 'o', 40, 40, 12, 5}

    // Set the populated flag to true.
    return true
}

