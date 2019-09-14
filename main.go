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
	Version      = "0.0"

	G Game

	PlayerName  string
	PlayerClass *types.ClassTypeInfo

	// global debug mode
	DeveloperMode = false

	// global error message via menu
	MenuErrorMsg string

	// map for all of the creature types
	GlobalCreatureTypeInfoMap = make(map[string]types.CreatureTypeInfo)

	// check if the creature types has already been populated.
	GlobalCreatureTypeInfoMapIsPopulated = false

	// map for all of the item types.
	GlobalItemTypeInfoMap = make(map[string]types.ItemTypeInfo)

	// check if the item types has already been populated.
	GlobalItemTypeInfoMapIsPopulated = false

	// map of all of the class types.
	GlobalClassTypeInfoMap = make(map[string]types.ClassTypeInfo)

	// check if the class types has already been populated.
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
