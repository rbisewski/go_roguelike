/*
 * File: creature.go
 *
 * Description: Class to handle player character and monsters.
 *
 * Future feature, add support for unicode char mobs like... 🦂
 */

package main

import "fmt"

// Coord is just an integer
type Coord int

// Structure to hold characters / monsters
type Creature struct {

    // Holds the name of the given creature.
    name string

    // Holds the type of the given creature.
    species string

    // Store the current (x,y) coord of the creature.
    Y  Coord
    X  Coord

    // Appearance of the creature.
    ch rune

    // Pointer to the level the creature is on.
    area *Area

    // Pointer to the stats attributes.
    *stats
}

// Structure to hold character attributes
type stats struct {
    Hp    int
    MaxHp int
    Att   int
    Def   int
}

//! Creature Constructor function.
/*
 * @param     string   creature name
 * @param     string   creature species (i.e type)
 * @param     Coord    y-value
 * @param     Coord    x-value
 * @param     rune     ASCII character graphic for the Creature
 * @param     Area*    object storing the level area 
 *
 * @return    Creature*    pointer to a new mob.
 */
func NewCreature(name string,
                 species string,
                 y Coord,
                 x Coord,
                 ch rune,
                 area *Area) *Creature {

    // Return an address to a newly allocated Creature object.
    return &Creature{name, species, y, x, ch, area, nil}
}

//! Monster Stats Constructor
/*
 * @param    int       current hit points
 * @param    int       maximum hit points
 * @param    int       attack
 * @param    int       defence
 *
 * @return   Stats*    pointer to a newly allocated Stats object.
 */
func newStats(hp, max, att, def int) *stats {
    return &stats{hp, max, att, def}
}

//! Creature w/ Stats Constructor
/*
 * @param     string    creature name
 * @param     string    creature species (i.e type)
 * @param     Coord     y-value
 * @param     Coord     x-value
 * @param     rune      ASCII character graphic
 * @param     Area*     pointer to an Area object
 * @param     int       current hit points
 * @param     int       maximum hit points
 * @param     int       attack
 * @param     int       defence
 *
 * @return    Creature*    pointer to a Creature w/ Stats
 */
func NewCreatureWithStats(name string,
                          species string,
                          y Coord,
                          x Coord,
                          ch rune,
                          area *Area,
                          hp int,
                          max int,
                          att int,
                          def int) *Creature {

    // Assign memory for a creature object and return the address.
    return &Creature{name, species, y, x, ch, area,
      newStats(hp, max, att, def)}
}

//! Function to move the mob to a new (x,y) location.
/*
 * @param     Coord    y-value
 * @param     Coord    x-value
 *
 * @return    none
 */
func (m *Creature) Move(y, x Coord) {

    // Input validation, make sure the (x,y) coords are reasonable.
    if x > 32767 || y > 32767 || x < -32767 || y < -32767 {

        // If debug mode is on, tell the end user what happened here.
        DebugLog(&G,"Error: Invalid (x,y) coord detected.")

        // Leave the function here.
        return
    }

    // Further input validation, make sure we actually got a creature here.
    if m == nil {

        // If debug mode is on, tell the end user what happened here.
        DebugLog(&G,"Error: Null pointer or invalid creature detected.")

        // Leave the function here.
        return
    }

    // Since this has a creature and it appears to have valid coords, then
    // go ahead and test it again the tile the creature in question wishes
    // to move to.
    tile_rune, blocks, hasCreature := m.area.GetTileInfo(m.Y+y, m.X+x)

    // Sanity check, make sure this actually got a tile rune.
    if tile_rune == 0 {

        // If debug mode is on, tell the end user what happened here.
        DebugLog(&G,"Error: Null or invalid Unicode rune.")

        // Leave the function here.
        return
    }

    // If the player attempts to move to a blocking tile, and it is a wall,
    // go ahead and print a short message and then leave function.
    if blocks && m.species == "player" && tile_rune == '#' {
        MessageLog.log("The wall is solid and damp, and you cannot move past.")
        return
    }

    // Catch-all message for when the player moves into a blocking tile.
    if blocks && m.species == "player" {
        MessageLog.log("Something here is blocking, and you cannot move past.")
        return
    }

    // If some other creature attempts to move, simply return here since
    // there is no need to print a message, except in debug mode.
    if blocks && m.species != "player" {

        // If debug mode, tell the developer where the creature has moved to.
        DebugLog(&G, fmt.Sprintf(
                 "The %s attempted to move to location (%d,%d), but it " +
                 "was blocked.",
                 m.name,
                 m.Y+y,
                 m.X+x))

        // Leave the function here.
        return
    }

    // If the tile is non-blocking, but a creature is here, go ahead and
    // switch to combat mode via the attack() function.
    if hasCreature != nil && m != hasCreature  {

        // If debug mode, tell the developer which creature is being attacked.
        DebugLog(&G, fmt.Sprintf(
                 "The %s is attacking %s at location (%d,%d).",
                 m.name,
                 hasCreature.name,
                 m.Y+y,
                 m.X+x))

        // Call the attack() function.
        m.attack(hasCreature)

        // Leave the function here.
        return
    }

    // If debug mode, tell the developer where the creature has moved to.
    DebugLog(&G, fmt.Sprintf(
             "The %s moved to location (%d,%d).",
             m.name,
             m.Y+y,
             m.X+x))

    // Since the tile is non-blocking, and no creature is present, then
    // go ahead and move there.
    m.Y += y
    m.X += x

    // Finally, leave the function since the move is finished.
    return
}

//! Function to handle what occurs if a monster attacks.
/*
 * @param     Creature*    defending creature / PC
 *
 * @return    none
 */
func (attacker *Creature) attack(defender *Creature) {

    // Adjust the defender's HP based on the damage dealt.
    defender.Hp -= attacker.Att - defender.Def

    // If the HP of the defender falls below zero, then he/she/it dies.
    if defender.Hp <= 0 {
        defender.die()
    }

    // Give the end-user a clue what is going on.
    //if m.ch == '@' {
        MessageLog.log(fmt.Sprintf("The %s has %d HP left.",
                                   attacker.name,
                                   attacker.Hp))
    //}
}

//! Function to handle what occurs when a monster dies.
/*
 * @param     Creature*    monster who is currently dying
 *
 * @return    none
 */
func (m *Creature) die() {

    // Change the character to a '%', eventually items will be added.
    m.ch = '%'

    // Adjust the array of monsters to account for the newly dead monster.
    for i, mm := range m.area.Creatures {

        // Uponing finding the monster who died...
        if m == mm {

            // Account for all of the monsters.
            m.area.Creatures = append(m.area.Creatures[:i], m.area.Creatures[i+1:]...)

            // As well as all of the monster's items.
            m.area.Items = append(m.area.Items, m)
        }
    }

    // If the "monster" who died is the player, then call that routine.
    if m == G.Player {
        G.Death()
    }
}
