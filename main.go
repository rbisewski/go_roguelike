/*
 * File: main.go
 *
 * Description: Contains the main.go routine.
 */

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rbisewski/go_roguelike/types"
)

var (
	printVersion = false

	// Version ... stores the version of the software
	Version = "0.0"

	// G ... global game object
	G Game

	// PlayerName ... stores the name of the player character
	PlayerName string

	// PlayerClass ... stores the name of the player class
	PlayerClass *types.ClassTypeInfo

	// DeveloperMode ... global debug mode
	DeveloperMode = false

	// MenuErrorMsg ... global error message via menu
	MenuErrorMsg string

	// GlobalCreatureTypeInfoMap ... map for all of the creature types
	GlobalCreatureTypeInfoMap = make(map[string]types.CreatureTypeInfo)

	// GlobalCreatureTypeInfoMapIsPopulated ... check if the creature types has already been populated.
	GlobalCreatureTypeInfoMapIsPopulated = false

	// GlobalItemTypeInfoMap ... map for all of the item types.
	GlobalItemTypeInfoMap = make(map[string]types.ItemTypeInfo)

	// GlobalItemTypeInfoMapIsPopulated ... check if the item types has already been populated.
	GlobalItemTypeInfoMapIsPopulated = false

	// GlobalClassTypeInfoMap ... map of all of the class types.
	GlobalClassTypeInfoMap = make(map[string]types.ClassTypeInfo)

	// GlobalClassTypeInfoMapIsPopulated ... check if the class types has already been populated.
	GlobalClassTypeInfoMapIsPopulated = false
)

func init() {
	flag.BoolVar(&printVersion, "version", false,
		"Print the current version of this program and exit.")
}

func main() {

	flag.Parse()

	if printVersion {
		fmt.Println("go-roguelike v" + Version)
		os.Exit(0)
	}

	Init()
	defer End()

	// setup creature types, item types, and player class types
	GlobalCreatureTypeInfoMapIsPopulated = types.GenCreatureTypes(GlobalCreatureTypeInfoMap)
	GlobalItemTypeInfoMapIsPopulated = types.GenItemTypes(GlobalItemTypeInfoMap)
	GlobalClassTypeInfoMapIsPopulated = types.GenClassTypes(GlobalClassTypeInfoMap)

	G.state = "menu"

	G.DebugMode = DeveloperMode

	// infinite loop which is present as long as the game is running
	for !G.state.Quiting() {

		// In the menu?
		if G.state.Menuing() {
			G.state = G.Menu()
			continue
		}

		// handlers for screen output and keyboard input
		G.Output()
		G.Input()
	}
}
