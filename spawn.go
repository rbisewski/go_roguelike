/*
 * File: spawn.go
 *
 * Description: File to handle spawning new monsters into a given area.
 */

package main

import "fmt"

//! Function to spawn a creature in a given area.
/*
 * @param     string    name of the creature to add
 * @param     int       x-coord as int
 * @param     int       y-coord as int
 * @param     Area*     pointer to the intended area
 *
 * @return    bool      whether or not the creature was added
 */
func spawnCreatureToArray(name string, x int, y int, a *Area) bool {

    // Input validation, make sure this got a valid string, coords, and area.
    if len(name) < 1 || x < 0 || y < 0 || a == nil {
        DebugLog(&G, fmt.Sprintf("spawnCreatureToArray() --> invalid input"))
        return false
    }

    // Safety check, if the creature type info has not yet been populated,
    // then leave this function.
    if !GlobalCreatureTypeInfoMapIsPopulated {
        return false
    }

    // Further safety check, make sure that creature actually exists as a
    // valid type in the global creature type array.
    _, IsCreatureTypeDefined := GlobalCreatureTypeInfoMap[name]

    // If the given creature name is not present, leave and return false.
    if !IsCreatureTypeDefined {

        // When debug mode is enabled, also log a message about the improper
        // creature name string given.
        DebugLog(&G, fmt.Sprintf("spawnCreatureToArray() --> improper " +
          "monster string given: %s", name))

        // Send back a false since the creature was *not* spawned.
        return false
    }

    // Grab the creature's name, species, rune-graphic, health, max-health,
    // attack, and defence attributes from the global creature type map.
    SpawnedCreatureName    := GlobalCreatureTypeInfoMap[name].Name
    SpawnedCreatureSpecies := GlobalCreatureTypeInfoMap[name].Species
    SpawnedCreatureGfx     := GlobalCreatureTypeInfoMap[name].Ch
    SpawnedCreatureHp      := GlobalCreatureTypeInfoMap[name].Hp
    SpawnedCreatureMaxHp   := GlobalCreatureTypeInfoMap[name].MaxHp
    SpawnedCreatureAttack  := GlobalCreatureTypeInfoMap[name].Att
    SpawnedCreatureDefence := GlobalCreatureTypeInfoMap[name].Def

    // Append it to the array.
    a.Creatures = append(a.Creatures, NewCreature(SpawnedCreatureName,
      SpawnedCreatureSpecies, y, x, SpawnedCreatureGfx, a, nil,
      SpawnedCreatureHp, SpawnedCreatureMaxHp, SpawnedCreatureAttack,
      SpawnedCreatureDefence))

    // With the monster successfully added, consider this complete.
    return true
}
