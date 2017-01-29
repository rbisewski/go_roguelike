/*
 * File: ai.go
 *
 * Description: Class to handle monster AI.
 */

package main

import "math"
import "fmt"

//! Routine to handle the monster AI actions.
/*
 * @params     Game*    pointer to the game instance
 *
 * @returns    none
 */
func (g *Game) processAi() {

    // Variable declaration
    var dy, dx Coord

    // Cycle thru all of the creatures present in the area level.
    for _, m := range g.Area.Creatures {

        // If not the PC then do this (i.e. monsters).
        if m != g.Player {

            // Figure out the difference between the player coords and
            // given monster.
            ydist := g.Player.Y - m.Y
            xdist := g.Player.X - m.X

            // Set the current viewing distance.
            distance := math.Sqrt(float64(xdist*xdist + ydist*ydist))

            // Determine the derived distances.
            dx = Coord(Round(float64(int(xdist)/Round(distance))))
                dy = Coord(Round(float64(int(ydist)/Round(distance))))

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
}
