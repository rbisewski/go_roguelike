/*
 * File: ai.go
 *
 * Description: Class to handle monster AI.
 */

package main

import "fmt"
import "math"

//! Routine to handle the monster AI actions.
/*
 * @params     Game*    pointer to the game instance
 *
 * @returns    none
 */
func (g *Game) process_ai() {

    // Variable declaration
    var dy, dx int

    // Cycle thru all of the creatures present in the area level.
    for _, m := range g.Area.Creatures {

        // The player has no need for AI as the human controls it.
        if m == g.Player {
            continue
        }

        // Figure out the difference between the player coords and
        // given monster.
        ydist := g.Player.Y - m.Y
        xdist := g.Player.X - m.X

        // Safety check, if both the x and y distance are zero, do nothing.
        if ydist == 0 && xdist == 0 {
            continue
        }

        // Set the current viewing distance.
        distance := math.Sqrt(float64(xdist*xdist + ydist*ydist))

        // Safety check, if the distance is zero, do nothing.
        if distance == 0 {
            continue
        }

        // If the distance is too big between the player and the creature,
        // do nothing. This will stop the creature from attempting to move
        // needlessly.
        if (distance > 6) {
            continue
        }

        // Determine the derived distances.
        dx = Round(float64(int(xdist)/Round(distance)))
        dy = Round(float64(int(ydist)/Round(distance)))

        // If debug mode, then display this...
        DebugLog(g, fmt.Sprintf("dx, dy, dist = %d, %d, %g->%d | xdist: %d - ydist: %d    ",
                    dx,
                    dy,
                    distance,
                    Round(distance),
                    xdist,
                    ydist))

        // Tell the monster to move to the determined location.
        m.Move(dy, dx)
    }
}
