/*
 * File: main.go
 *
 * Description: Contains the main.go routine.
 */

package main

// Import the types since the game needs to populate the creature type
// for the purpose of spawning creatures later on.
import (
	"flag"
	"fmt"
	"os"

	"github.com/rbisewski/go_roguelike/types"
)

//
// Global variable declaration.
//
var (
	// Whether or not to print the current version of the program
	printVersion = false
	Version      = "0.0"

	G Game

	PlayerName  string
	PlayerClass *types.ClassTypeInfo

	// DeveloperMode ... global debug mode
	DeveloperMode = false

	// MenuErrorMsg ... global error message via menu
	MenuErrorMsg string

	// GlobalCreatureTypeInfoMap ... Global variable to hold all of the creature
	// types.
	GlobalCreatureTypeInfoMap = make(map[string]types.CreatureTypeInfo)

	// GlobalCreatureTypeInfoMapIsPopulated ... Global variable to check if the
	// creature types has already been populated.
	GlobalCreatureTypeInfoMapIsPopulated = false

	// GlobalItemTypeInfoMap ... global variable to hold all of the item types.
	GlobalItemTypeInfoMap = make(map[string]types.ItemTypeInfo)

	// GlobalItemTypeInfoMapIsPopulated ... Global variable to check if the item
	// types has already been populated.
	GlobalItemTypeInfoMapIsPopulated = false

	// GlobalClassTypeInfoMap ... global variable to hold all of the class types.
	GlobalClassTypeInfoMap = make(map[string]types.ClassTypeInfo)

	// GlobalClassTypeInfoMapIsPopulated ... Global variable to check if the
	// class types has already been populated.
	GlobalClassTypeInfoMapIsPopulated = false
)

func init() {

	// Version mode flag
	flag.BoolVar(&printVersion, "version", false,
		"Print the current version of this program and exit.")
}

//
// Main
//
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
