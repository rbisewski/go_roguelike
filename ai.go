/*
 * File: ai.go
 *
 * Description: Class to handle monster AI.
 */

package main

import (
	"fmt"
	"math"
)

//! Routine to handle the monster AI actions.
/*
 * @params     Game*    pointer to the game instance
 *
 * @returns    none
 */
func (g *Game) processAI() {

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
		// have the monster do nothing half of the time.
		if distance > 6 && TossCoin() {
			continue
		}

		// The other half of the time, pick a random location nearby from
		// one of the 8 cells nearby to move the creature.
		if distance > 6 {

			// Variable declaration.
			dx := 0
			dy := 0

			// Pick a number between 0-8 and add 1, so as to make it similar
			// in concept to a QWERTY keypad setup:
			//
			// 7 8 9
			// 4 5 6
			// 1 2 3
			//
			// Where 7 is equal to north-east, 2 is equal to south, and etc.
			//
			KeypadLocation := getRandomNumBetweenZeroAndMax(8) + 1

			// Adjust the X if the creature moved to the west.
			if KeypadLocation == 3 || KeypadLocation == 6 || KeypadLocation == 9 {
				dx = -1

				// Adjust the X if the creature moved to the east.
			} else if KeypadLocation == 7 || KeypadLocation == 4 || KeypadLocation == 1 {
				dx = 1
			}

			// Adjust the Y if the creature moved to the north.
			if KeypadLocation == 7 || KeypadLocation == 8 || KeypadLocation == 9 {
				dy = -1

				// Adjust the Y if the creature moved to the south.
			} else if KeypadLocation == 1 || KeypadLocation == 2 || KeypadLocation == 3 {
				dy = 1
			}

			// Attempt to move the creature to that location.
			m.Move(dy, dx)

			// This creature's movement action is now completed.
			continue
		}

		// Determine the derived distances.
		dx = Round(float64(int(xdist) / Round(distance)))
		dy = Round(float64(int(ydist) / Round(distance)))

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
