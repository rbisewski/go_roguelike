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
    Y  Coord
    X  Coord
    ch rune
    area *Area
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
 * @param     Coord    y-value
 * @param     Coord    x-value
 * @param     rune     ASCII character graphic for the Creature
 * @param     Area*    object storing the level area 
 *
 * @return    Creature*    pointer to a new mob.
 */
func NewCreature(y Coord, x Coord, ch rune, area *Area) *Creature {

    // Return an address to a newly allocated Creature object.
    return &Creature{y, x, ch, area, nil}
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
 * @param     Coord   y-value
 * @param     Coord   x-value
 * @param     rune    ASCII character graphic
 * @param     Area*   pointer to an Area object
 * @param     int     current hit points
 * @param     int     maximum hit points
 * @param     int     attack
 * @param     int     defence
 *
 * @return    Creature*    pointer to a Creature w/ Stats
 */
func NewCreatureWithStats(y Coord, x Coord, ch rune, area *Area, hp, max, att, def int) *Creature {
    return &Creature{y, x, ch, area, newStats(hp, max, att, def)}
}

//! Function to move the mob to a new (x,y) location.
/*
 * @param     Coord    y-value
 * @param     Coord    x-value
 *
 * @return    none
 */
func (m *Creature) Move(y, x Coord) {

    // If there is either a monster or a non-blocking tile, then do this...
    if blocks, hasCreature := m.area.IsBlocking(m.Y+y, m.X+x); !blocks {

        // If the chosen square has no other monster present, then move there.
        if hasCreature == nil || m == hasCreature {
            m.Y += y
            m.X += x
            return
        }

        // Run the attack function.
        m.attack(hasCreature)

        // Give the end-user a clue what is going on.
        MessageLog.log(fmt.Sprintf("Monster HP: %d", hasCreature.Hp))

        // End here.
        return
    }

    // If the player character, then print this message instead.
    if m.ch == '@' {
        MessageLog.log("The wall is solid and damp, and you cannot move past.")
    }
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
