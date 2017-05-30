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

// Global variable to check if the creature types has already been populated.
var GlobalCreatureTypeInfoMapIsPopulated = false

// Global variable to hold all of the item types.
var GlobalItemTypeInfoMap = make(map[string]types.ItemTypeInfo)

// Global variable to check if the item types has already been populated.
var GlobalItemTypeInfoMapIsPopulated = false

// Global variable to hold all of the class types.
var GlobalClassTypeInfoMap = make(map[string]types.ClassTypeInfo)

// Global variable to check if the class types has already been populated.
var GlobalClassTypeInfoMapIsPopulated = false

//
// Main
//
func main() {

    // Let's get (gocurses) started!
    Init()
    defer End()

    // Populate the various creature types into the game.
    GlobalCreatureTypeInfoMapIsPopulated = types.GenCreatureTypes(GlobalCreatureTypeInfoMap)

    // Populate the various item types into the game.
    GlobalItemTypeInfoMapIsPopulated = types.GenItemTypes(GlobalItemTypeInfoMap)

    // Populate the various class types into the game.
    GlobalClassTypeInfoMapIsPopulated = types.GenClassTypes(GlobalClassTypeInfoMap)

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
