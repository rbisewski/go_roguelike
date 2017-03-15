/*
 * File: game.go
 *
 * Description: Handles the game and its relevant data.
 */
package main

import (
    "fmt"
)

// Attributes for the `Game` structure.
type GameState string

// Structure used hold the game object.
type Game struct {

    // Determines whether or not to display debug messages.
    DebugMode bool

    // Current status of the game.
    state GameState

    // Pointer to the player-character object.
    Player *Creature

    // Pointer to area array.
    Area   *Area
}

//! Function to initialize the game.
/*
 * @param     Game*    pointer to a game object
 *
 * @return    none
 */
func (g *Game) Init() {

    // Variable declaration
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

    // Generate an area map
    g.Area, y, x = NewArea(240, 250)

    // Safety check, if the player name is blank, default to anonymous.
    if len(PlayerName) == 0 {
        PlayerName = "Anonymous"
    }

    // The player-character will be represented by an @ symbol.
    g.Player = NewCreature(PlayerName,
                           "player",
                           y,
                           x,
                           '@',
                           g.Area,
                           make([]*Item,0),
                           30,
                           30,
                           10,
                           5)

    // Attach the player-character creature to the map.
    g.Area.Creatures = append(g.Area.Creatures, g.Player)

    // Pass along the area, and populate the world with a number of monsters.
    g.Area.populateAreaWithCreatures()
}

//! Determines whether or not 
/*
 * @param      Gamestate    current state of the game.
 *
 * @returns    bool         whether or not the game is still in the menu.
 */
func (s GameState) Menuing() bool {
    return s == "menu"
}

//! Write out the in-game menu.
/*
 * @param     *Game       pointer to an instance of this game
 *
 * @returns   GameState   Current state of this game (e.g. "menu", "playing")
 */
func (g *Game) Menu() GameState {

    // Print out the title.
    Write(Percent(25, ConsoleHeight), ConsoleWidth/2,
      "GoRogue - A Rogue-like written in golang.")

    // Print out the options, which currently are as follows:
    //
    // * Start a new game.
    // * Load the last game.
    // * Quit the current game.
    //
    Write(Percent(25, ConsoleHeight)+2, ConsoleWidth/2,
      "Press 'N' to start a new game.")
    Write(Percent(25, ConsoleHeight)+3, ConsoleWidth/2,
      "Press 'L' to load a previous game.")
    Write(Percent(25, ConsoleHeight)+4, ConsoleWidth/2,
      "Press 'Q' to quit.")

    // Print out the most recent menu error message, if any.
    Write(Percent(25, ConsoleHeight)+6, ConsoleWidth/2,
      MenuErrorMsg)

    // Grab the current keyboard input.
    key := GetInput()

    // If the N key was pressed, then this needs to start a new game.
    if key == "N" || key == "n" {

        // Wipe away the current menu screen once the player has elected to start
        // a new game.
        Clear()

        // If any partly loaded data is present, clear it away.
        PlayerName = ""

        // Endless loop that is designed to allow the player character to enter
        // the name of their character by typing via the keyboard.
        for true {

            // Tell the end user to enter the name of their character.
            Write(Percent(25, ConsoleHeight), ConsoleWidth/2,
              "Enter the name of your character:")

            // Write the current PlayerName to the below console, which is
            // in location 'ConsoleHeight+2' so that it appears two lines below.
            Write(Percent(25, ConsoleHeight)+2, ConsoleWidth/2,
              PlayerName)

            // Grab the current keyboard input.
            key = GetInput()

            // If the enter key was pressed and the character name has at least
            // one or more characters.
            if WasEnterPressed(key) && len(PlayerName) > 0 {
                break

            // Else if it was a valid a-zA-Z key then...
            } else if IsAlphaCharacter(key) && len(PlayerName) < 13 {

                // ... append it to the name string.
                PlayerName += key

            // Else if it was a backspace or delete key && length of PlayerName
            // is greater than 0 then...
            } else if IsDeleteOrBackspace(key) && len(PlayerName) > 0 {

                // ... remove the last string from PlayerName.
                PlayerName = string(PlayerName[:len(PlayerName)-1])
            }

            // Wipe away the old screen, so that it can be reprinted during
            // the next cycle.
            Clear()
        }

        // Wipe away the old screen, in the event that some part of the previous
        // enter your character name interface happens to remain.
        Clear()

        // Initialize the game.
        g.Init()

        // Set the current GameState to "playing".
        return "playing"
    }

    // Pressed the "L" key? Then attempt to load a game...
    if key == "L" || key == "l" {

        // Setup the game environment
        g.Init()

        // Attempt to load the previous game.
        if !g.LoadGame("player.save") {

            // Tell the end user that the load was unsuccessful.
            MenuErrorMsg = "No recent save game was detected. Please start " +
              "a new game."

            // Tell the developer that the load function has failed.
            DebugLog(g, fmt.Sprintf("Menu() --> unable to load previous game"))

            // Return back to the global game menu.
            return "menu"
        }

        // Give the main game pad a height and width.
        SetPad(g.Area.Height, g.Area.Width)

        // Draw the recorded map.
        DrawMap(g.Area)

        // Wipe away the menu screen.
        Clear()

        // Set the current GameState to "playing".
        return "playing"
    }

    // If the end user pressed the Q key, attempt to exit the game.
    if key == "Q" || key == "q" {
        return "quit"
    }

    // Otherwise if none of the above keys were pressed, then default to
    // remaining in the initial game menu.
    return "menu"
}

//! Handle the event of a PC death (by monsters or the like).
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

//! Determine if the current game is in a state of quitting.
/*
 * @param    GameState    current GameState
 *
 * @return   bool         whether or not the game is quitting.
 */
func (s GameState) Quiting() bool {
    return s == "quit"
}

//! Generate the game screen output.
/*
 * @param    Game*    pointer to the current game instance
 *
 * @return   none
 */
func (g *Game) Output() {

    // Draw the given area map.
    DrawMap(g.Area)

    // Cycle thru all of the item present in the current area. If an item is
    // at the given (x,y) coords, then go ahead and draw it on the map.
    for index, item := range g.Area.Items {

        // Sanity check, make sure this actually got a valid item.
        if (item == nil) {

            // Otherwise tell the developer something odd was appended here.
            DebugLog(g, fmt.Sprintf("Output() --> invalid item at index [%i]",
              index))

            // Move on to the next item.
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

        // Sanity check, make sure this actually got a valid monster.
        if (m == nil) {

            // Otherwise tell the developer something odd was appended here.
            DebugLog(g, fmt.Sprintf("Output() --> null monster at index [%i]",
              index))

            // Move on to the next monster.
            continue
        }

        // Draw a monster at its current coords.
        Draw(m.Y, m.X, m.ch)
    }

    // Refresh the tile the PC is currently on.
    RefreshPad(int(g.Player.Y), int(g.Player.X))

    // Update the current PC stats.
    g.Player.UpdateStats()
}

//! Keyboard input parser.
/*
 * @param     Game*   pointer to the current game instance
 *
 * @return    none
 */
func (g *Game) Input() {

    // Grab the key command pressed.
    key := GetInput()

    // If debug enabled, show the key that was pressed.
    DebugLog(g, fmt.Sprintf("Key pressed --> %x", key))

    // Convert the key pressed to a hex string value.
    key_as_string := fmt.Sprintf("%x", key)

    // For a given key...
    switch key_as_string {

    // Numpad 8 --> Move player north
    case "38":
        g.Player.Move(-1, 0)
        g.process_ai()

    // Numpad 9 --> Move player north-east
    case "39":
        g.Player.Move(-1, 1)
        g.process_ai()

    // Numpad 6 --> Move player east
    case "36":
        g.Player.Move(0, 1)
        g.process_ai()

    // Numpad 3 --> Move player south-west
    case "33":
        g.Player.Move(1, 1)
        g.process_ai()

    // Numpad 2 --> Move player south
    case "32":
        g.Player.Move(1, 0)
        g.process_ai()

    // Numpad 1 --> Move player south-east
    case "31":
        g.Player.Move(1, -1)
        g.process_ai()

    // Numpad 4 --> Move player west
    case "34":
        g.Player.Move(0, -1)
        g.process_ai()

    // Numpad 7 --> Move player north-west
    case "37":
        g.Player.Move(-1, -1)
        g.process_ai()

    // Down Arrow --> Move player south
    case "c482":
        g.Player.Move(1, 0)
        g.process_ai()

    // Up Arrow --> Move player north
    case "c483":
        g.Player.Move(-1, 0)
        g.process_ai()

    // Left Arrow --> Move player west
    case "c484":
        g.Player.Move(0, -1)
        g.process_ai()

    // Right Arrow --> Move player east
    case "c485":
        g.Player.Move(0, 1)
        g.process_ai()

    // I || i --> Open inventory
    case "49":
    case "69":
        //g.OpenInventory()

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
