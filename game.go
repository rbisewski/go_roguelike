/*
 * File: game.go
 *
 * Description: Handles the game and its relevant data.
 */

package main

import (
	"fmt"
	"strconv"
)

// GameState ... Attributes for the `Game` structure.
type GameState string

// Game ... Structure used hold the game object.
type Game struct {

	// Determines whether or not to display debug messages.
	DebugMode bool

	// Current status of the game.
	state GameState

	// Pointer to the player-character object.
	Player *Creature

	// Pointer to area array.
	Area *Area

	// List of items on the ground at a give coord
	GroundItems []*Item
}

// Init ... Function to initialize the game.
/*
 * @param     Game*    pointer to a game object
 *
 * @return    none
 */
func (g *Game) Init() {

	var y int
	var x int

	// Grab the global setting and assign it to the current instance of
	// this game, specifically this will enable / disable debugging
	// functionality and messages.
	//
	// This setting can be adjusted in the main.go included in this project.
	//
	g.DebugMode = DeveloperMode

	// Set the game state.
	g.state = "menu"

	// Initially the player is not picking up items from thr ground.
	g.GroundItems = make([]*Item, 0)

	// Generate an area map
	g.Area, y, x = NewArea(240, 250)

	// Safety check, if the player name is blank, default to anonymous.
	if len(PlayerName) == 0 {
		PlayerName = "Anonymous"
	}

	// Safety check, if the player class is blank, default to warrior.
	if PlayerClass == nil {
		defaultClass := GlobalClassTypeInfoMap["1"]
		PlayerClass = &defaultClass
	}

	// The player-character will be represented by an @ symbol.
	g.Player = NewCreatureWithEquipment(PlayerName, "player", y, x, '@',
		g.Area, make([]*Item, 0), 30, 30, 10, 5, PlayerClass, 10, 10, 10,
		10, 10, 0)

	// Attach the player-character creature to the map.
	g.Area.Creatures = append(g.Area.Creatures, g.Player)

	// Pass along the area, and populate the world with a number of monsters.
	g.Area.populateAreaWithCreatures()
}

// Menuing ... Determines if in-menu.
func (s GameState) Menuing() bool {
	return s == "menu"
}

// Menu ... Write out the in-game menu.
/*
 * @param     *Game       pointer to an instance of this game
 *
 * @returns   GameState   Current state of this game (e.g. "menu", "playing")
 */
func (g *Game) Menu() GameState {

	var state GameState = ""

	Write(Percent(25, ConsoleHeight), ConsoleWidth/2, "GoRogue - A Rogue-like written in golang.")
	Write(Percent(25, ConsoleHeight)+2, ConsoleWidth/2, "Press 'N' to start a new game.")
	Write(Percent(25, ConsoleHeight)+3, ConsoleWidth/2, "Press 'L' to load a previous game.")
	Write(Percent(25, ConsoleHeight)+4, ConsoleWidth/2, "Press 'Q' to quit.")

	// Print out the most recent menu error message, if any.
	Write(Percent(25, ConsoleHeight)+6, ConsoleWidth/2, MenuErrorMsg)

	key := GetInput()
	switch key {
	case "n":
		fallthrough
	case "N":
		// If any partly loaded data is present, clear it away.
		Clear()
		PlayerName = ""

		// Endless loop that is designed to allow the player character to enter
		// the name of their character by typing via the keyboard.
		for true {

			Write(Percent(25, ConsoleHeight), ConsoleWidth/2, "Enter the name of your character:")
			Write(Percent(25, ConsoleHeight)+2, ConsoleWidth/2, PlayerName)

			key = GetInput()

			if WasEnterPressed(key) && len(PlayerName) > 0 {
				// if Enter was pressed and name is at least 1
				// then assume end-user is done typing their name
				break

			} else if IsAlphaCharacter(key) && len(PlayerName) < 13 {
				PlayerName += key

			} else if IsDeleteOrBackspace(key) && len(PlayerName) > 0 {
				PlayerName = string(PlayerName[:len(PlayerName)-1])
			}

			// Wipe away the old screen, so that it can be reprinted during
			// the next cycle.
			Clear()
		}

		// Endless loop that is designed to allow the player character to enter
		// the class of their character by typing via the keyboard.
		classCounter := 1
		for true {

			Write(Percent(25, ConsoleHeight), ConsoleWidth/2, "The name of your character is:     ")
			Write(Percent(25, ConsoleHeight)+2, ConsoleWidth/2, PlayerName)
			Write(Percent(25, ConsoleHeight)+5, ConsoleWidth/2, "Now select a class:")

			for range GlobalClassTypeInfoMap {

				// Obtain the given class
				strref := strconv.Itoa(classCounter)
				givenClass := GlobalClassTypeInfoMap[strref]

				// If unknown, move to the next element.
				if givenClass.Name == "Unknown" || len(givenClass.Name) < 1 {
					continue
				}

				// Print out the given class options.
				Write(Percent(25, ConsoleHeight)+6+classCounter, ConsoleWidth/2, strref+") "+givenClass.Name+"   ")

				classCounter++
			}

			// if the player selected a class, print it out so that there
			// is feedback for the player to see
			if PlayerClass != nil {
				Write(Percent(25, ConsoleHeight)+8+classCounter, ConsoleWidth/2, "You have selected... "+PlayerClass.Name+"     ")
				Write(Percent(25, ConsoleHeight)+10+classCounter, ConsoleWidth/2, "Press [Enter] to begin the game.")
			}

			key = GetInput()

			// If a number has been pressed...
			if IsNumeric(key) {

				// convert the keystroke to an uint64 representation of the
				// original hexidecimal typed by the keyboard
				num, err := ConvertKeyToNumeric(key)
				if err != nil {
					num = 1
				}

				// attempt to grab the selected class using the above
				n := strconv.FormatUint(num, 10)
				selectedClass := GlobalClassTypeInfoMap[n]
				PlayerClass = &selectedClass
			}

			// If the enter key was pressed and the character class has
			// been selected by the player.
			if WasEnterPressed(key) && PlayerClass != nil {
				break
			}

			classCounter = 1
			Clear()
		}

		Clear()

		g.Init()
		state = "playing"

	case "l":
		fallthrough
	case "L":

		g.Init()

		// Attempt to load the previous game.
		if !g.LoadGame("player.sav") {

			MenuErrorMsg = "No recent save game was detected. Please start a new game."
			DebugLog(g, fmt.Sprintf("Menu() --> unable to load previous game"))

			state = "menu"
			break
		}

		// Give the main game pad a height and width.
		SetPad(g.Area.Height, g.Area.Width)

		// Draw the recorded map.
		DrawMap(g.Area)

		// Wipe away the menu screen.
		Clear()
		state = "playing"

	case "q":
		fallthrough
	case "Q":
		state = "quit"

	default:
		state = "menu"
	}

	return state
}

// Death ... handle the event of a PC death (by monsters or the like).
/*
 * @param      Game    current game instance
 *
 * @returns    none
 */
func (g *Game) Death() {

	// Wipe away the game screen.
	Clear()

	// Print out helpful death messages.
	Write(Percent(25, ConsoleHeight),
		ConsoleWidth/2,
		"Death overcomes you...")
	Write(Percent(25, ConsoleHeight)+1,
		ConsoleWidth/2,
		"Banished from the realm of the living for all time.")

	// Grab the present keyboard input.
	GetInput()

	// Set the current game state to "quit"
	g.state = "quit"
}

// Quiting ... determine if the current game is in a state of quitting.
func (s GameState) Quiting() bool {
	return s == "quit"
}

// Output ... generate the game screen output.
func (g *Game) Output() {

	DrawMap(g.Area)

	// Cycle thru all of the item present in the current area. If an item is
	// at the given (x,y) coords, then go ahead and draw it on the map.
	for index, item := range g.Area.Items {

		if item == nil {
			numStr := strconv.Itoa(index)
			DebugLog(g, "Output() --> invalid or null item at index ["+numStr+"]")
			continue
		}

		// If the item is in fact a dead creature, colour it red.
		if item.category == "corpse" {

			// Draw the item char rune with the red colour.
			DrawColours(item.Y, item.X, item.ch, 1)

			// Move on to the next item.
			continue
		}

		// Otherwise this is just a plain ol' item, so assign colours for a
		// striking magenta-black colour.
		DrawColours(item.Y, item.X, item.ch, 3)
	}

	// Cycle thru every monster present in the given area.
	for index, m := range g.Area.Creatures {

		if m == nil {
			numStr := strconv.Itoa(index)
			DebugLog(g, "Output() --> null monster at index ["+numStr+"]")
			continue
		}

		// Draw a monster at its current coords.
		Draw(m.Y, m.X, m.ch)
	}

	// Refresh the tile the PC is currently on.
	RefreshPad(int(g.Player.Y), int(g.Player.X))
	g.Player.UpdateStats()
}

// Input ... Keyboard input parser.
/*
 * @param     Game*   pointer to the current game instance
 *
 * @return    none
 */
func (g *Game) Input() {

	key := GetInput()
	DebugLog(g, fmt.Sprintf("Key pressed --> %x", key))

	// Convert the key pressed to a hex string value.
	keyAsString := fmt.Sprintf("%x", key)

	// If the player presses the ESC key in the inventory or ground items
	// screen, then switch back to playing mode.
	if (g.state == "equipment" || g.state == "inventory" ||
		g.state == "ground_items") && keyAsString == "1b" {

		// Set the state back to playing.
		g.state = "playing"
		return

		// Equipment screen is open and the player presses the left arrow.
	} else if g.state == "equipment" && keyAsString == "c484" {

		// Draw and populate the inventory ncurses UI.
		g.state = "equipment"
		DrawEquipmentUI(g, keyAsString)
		return

		// Equipment screen is open and the player presses the left arrow.
	} else if g.state == "equipment" && keyAsString == "c485" {

		// Draw and populate the inventory ncurses UI.
		g.state = "inventory"
		DrawInventoryUI(g, keyAsString)
		return

		// Inventory screen is open and the player presses the left arrow.
	} else if g.state == "inventory" && keyAsString == "c484" {

		// Draw and populate the inventory ncurses UI.
		g.state = "equipment"
		DrawEquipmentUI(g, keyAsString)
		return

		// Inventory screen is open and the player presses the right arrow.
	} else if g.state == "inventory" && keyAsString == "c485" {

		// Draw the UI and populate the global list of ground items.
		g.state = "ground_items"
		DrawGroundItemsUI(g, keyAsString)

		// Leave here since this needs to continue showing the ground
		// items UI to the player.
		return

		// Ground items screen is open and the player presses the left arrow.
	} else if g.state == "ground_items" && keyAsString == "c484" {

		// Draw and populate the inventory ncurses UI.
		g.state = "inventory"
		DrawInventoryUI(g, keyAsString)
		return

		// Ground items screen is open and the player presses the right arrow.
	} else if g.state == "ground_items" && keyAsString == "c485" {

		// Draw the UI and populate the global list of ground items.
		DrawGroundItemsUI(g, keyAsString)

		// Leave here since this needs to continue showing the ground
		// items UI to the player.
		return

		// If the player character inventory is open, and the key being pressed
		// is not "e" then do nothing.
	} else if g.state == "equipment" && keyAsString != "65" {

		// Draw and populate the inventory ncurses UI.
		DrawEquipmentUI(g, keyAsString)
		return

		// If the player character inventory is open, and the key being pressed
		// is not "i" then do nothing.
	} else if g.state == "inventory" && keyAsString != "69" {

		// Draw and populate the inventory ncurses UI.
		DrawInventoryUI(g, keyAsString)
		return

		// If the ground items UI is open, and the key being pressed
		// is not "g" then do nothing.
	} else if g.state == "ground_items" && keyAsString != "67" {

		// Draw the UI and populate the global list of ground items.
		DrawGroundItemsUI(g, keyAsString)

		// Do a check to see if a player presses the key 1-7 then attempt
		// to add that item to the player's inventory.
		err := PickupGroundItem(g, keyAsString)

		// If there was an error, print it out.
		if err != nil {
			fmt.Printf(err.Error())
			return
		}

		// Draw the UI and populate the global list of ground items.
		DrawGroundItemsUI(g, keyAsString)

		// Leave here since this needs to continue showing the ground
		// items UI to the player.
		return
	}

	// For a given key...
	switch keyAsString {

	// Numpad 8 --> Move player north
	case "38":
		g.Player.Move(-1, 0)
		g.processAI()

	// Numpad 9 --> Move player north-east
	case "39":
		g.Player.Move(-1, 1)
		g.processAI()

	// Numpad 6 --> Move player east
	case "36":
		g.Player.Move(0, 1)
		g.processAI()

	// Numpad 3 --> Move player south-west
	case "33":
		g.Player.Move(1, 1)
		g.processAI()

	// Numpad 2 --> Move player south
	case "32":
		g.Player.Move(1, 0)
		g.processAI()

	// Numpad 1 --> Move player south-east
	case "31":
		g.Player.Move(1, -1)
		g.processAI()

	// Numpad 4 --> Move player west
	case "34":
		g.Player.Move(0, -1)
		g.processAI()

	// Numpad 7 --> Move player north-west
	case "37":
		g.Player.Move(-1, -1)
		g.processAI()

	// Down Arrow --> Move player south
	case "c482":
		g.Player.Move(1, 0)
		g.processAI()

	// Up Arrow --> Move player north
	case "c483":
		g.Player.Move(-1, 0)
		g.processAI()

	// Left Arrow --> Move player west
	case "c484":
		g.Player.Move(0, -1)
		g.processAI()

	// Right Arrow --> Move player east
	case "c485":
		g.Player.Move(0, 1)
		g.processAI()

	// e --> Open / close the equipment screen.
	case "65":

		// If the equipment screen is not yet open.
		if g.state != "equipment" {

			// Enable the equipment state and draw the UI
			g.state = "equipment"
			DrawEquipmentUI(g, keyAsString)
			return
		}

		// Otherwise the equipment screen is open, so flip the state.
		g.state = "playing"

	// g --> Open / close the grab-item-from-ground interface.
	case "67":

		// If the ground items UI screen is not yet open.
		if g.state != "ground_items" {

			// Enable the inventory state and draw the UI
			g.state = "ground_items"
			DrawGroundItemsUI(g, keyAsString)
			return
		}

		// Otherwise the ground items UI is open, so flip the state.
		g.state = "playing"

	// i --> Open / close the inventory of the player character.
	case "69":

		// If the inventory is not yet open.
		if g.state != "inventory" {

			// Enable the inventory state and draw the UI
			g.state = "inventory"
			DrawInventoryUI(g, keyAsString)
			return
		}

		// Otherwise the inventory is open, so flip the state.
		g.state = "playing"

	// S --> Save game
	case "53":
		if Confirm("Save and Quit? Y/N") {
			g.SaveGame()
			MessageLog.log("Game Saved")
			g.state = "quit"
		}

	// Q --> Quit game
	case "51":
		if Confirm("Quit Without Saving? Y/N") {
			g.state = "quit"
		}
	}
}
