/*
 * File: main.go
 *
 * Description: Contains the main.go routine.
 */

package main

// Import the types since the game needs to populate the creature type
// for the purpose of spawning creatures later on.
import "./types"

// Global variable declaration.
var G              Game
var DeveloperMode  bool   = false
var PlayerName     string = ""
var MenuErrorMsg   string = ""

// Global variable to hold all of the creature types.
var GlobalCreatureTypeInfoMap = make(map[string]types.CreatureTypeInfo)

// Global variable to check if the map has already been populated.
var GlobalCreatureTypeInfoMapIsPopulated = false

//
// Main
//
func main() {

    // Let's get (gocurses) started!
    Init()
    defer End()

    // Populate the various creature types into the game.
    GlobalCreatureTypeInfoMapIsPopulated = types.GenCreatureTypes(GlobalCreatureTypeInfoMap)

    // The default state shall be to set the menu.
    G.state = "menu"

    // Set the debug mode flag.
    G.DebugMode = false

    // As long as we're not quting, then do this...
    for !G.state.Quiting() {

        // In the menu?
	if G.state.Menuing() {

            // The state remains on menu then!
	    G.state = G.Menu()
	    continue
	}

        // Handle output.
	G.Output()

        // Handle input.
	G.Input()
    }
}
