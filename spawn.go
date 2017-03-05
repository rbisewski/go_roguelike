/*
 * File: spawn.go
 *
 * Description: File to handle spawning new monsters into a given area.
 */

package main

import "fmt"

//! Function to spawn a creature in a given area.
/*
 * @param     string         name of the creature to add
 * @param     Coord          x-coord as int
 * @param     Coord          y-coord as int
 * @param     Area*          pointer to the intended area 
 *
 * @return    bool           whether or not the creature was added
 */
func spawnCreatureToArray(name string, x Coord, y Coord, a *Area) bool {

    // Input validation, make sure this got a valid string, coords, and area.
    if len(name) < 1 || x < 0 || y < 0 || a == nil {
        DebugLog(&G, fmt.Sprintf("spawnCreatureToArray() --> invalid input"))
        return false
    }

    // If a "dog" creature was requested...
    if name == "dog" {

        // Append it to the array.
        a.Creatures = append(a.Creatures, NewCreatureWithStats("dog",
          "canine", y, x, 'd', a, nil, 20, 30, 5, 0))

        // With the monster successfully added, consider this complete.
        return true
    }

    // Otherwise an invalid string was given, so tell the dev this failed...
    DebugLog(&G, fmt.Sprintf("spawnCreatureToArray() --> improper monster " +
                             "string given: %s", name))

    // Finally give back a false here since this attempt failed.
    return false
}