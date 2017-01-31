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

    // Pointer to the current states of given creatures.
    //Creatures   []*Creature

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
    var y Coord
    var x Coord

    // The default setting is for debug messages to be hidden.
    g.DebugMode = false

    // Set the game state.
    g.state = "menu"

    // Generate an area map
    g.Area, y, x = NewArea(240, 250)

    // The PC will be represented by an @ symbol.
    g.Player = NewCreatureWithStats("player",
                                    "player",
                                    y,
                                    x,
                                    '@',
                                    g.Area,
                                    30,
                                    30,
                                    10,
                                    5)

    // Attach the player to the map.
    g.Area.Creatures = append(g.Area.Creatures,
                             g.Player)

    // Right now all this does is add a single monster.
    g.Area.Creatures = append(g.Area.Creatures,
                              NewCreatureWithStats("dog",
                                                   "canine",
                                                   y-10,
                                                   x-1,
                                                   'd',
                                                   g.Area,
                                                   20,
                                                   30,
                                                   15,
                                                   0))
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
    Write(Percent(25, ConsoleHeight),
          ConsoleWidth/2,
          "GoRogue - A Rogue-like written in golang.")

    // Print out the options.
    Write(Percent(25, ConsoleHeight)+1,
          ConsoleWidth/2,
          "Press any key to continue, press 'L' to load")

    // Grab the current keyboard input.
    key := GetInput()


    // Pressed the "L" key? Then attempt to load a game...
    if key == "L" {

        // Setup the game environment
        g.Init()

        // Attempt to load the previous game.
        g.LoadGame()

        // Give the main game pad a height and width.
        SetPad(g.Area.Height, g.Area.Width)

        // Draw the recorded map.
        DrawMap(g.Area)

        // Wipe away the menu screen.
        Clear()

        // Set the current GameState to "playing".
        return "playing"
    }

    // Wipe away the current menu screen once the player has elected to start
    // a new game or load a game.
    Clear()

    // Initialize the game.
    g.Init()

    // Set the current GameState to "playing".
    return "playing"
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

    // Cycle thru a given list of items in an area.
    for _, i := range g.Area.Items {

        // Assign colours for a given item.
        if i.Hp <= 0 {
            DrawColours(i.Y, i.X, i.ch, 1)
            continue
        }
    }

    // For every monster present in this area.
    for _, m := range g.Area.Creatures {

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

    // Numpad 9 --> Move player north-east
    case "39":
        g.Player.Move(-1, 1)

    // Numpad 6 --> Move player east
    case "36":
        g.Player.Move(0, 1)

    // Numpad 3 --> Move player south-west
    case "33":
        g.Player.Move(1, 1)

    // Numpad 2 --> Move player south
    case "32":
        g.Player.Move(1, 0)

    // Numpad 1 --> Move player south-east
    case "31":
        g.Player.Move(1, -1)

    // Numpad 4 --> Move player west
    case "34":
        g.Player.Move(0, -1)

    // Numpad 7 --> Move player north-west
    case "37":
        g.Player.Move(-1, -1)

    // Down Arrow --> Move player south
    case "c482":
        g.Player.Move(1, 0)

    // Up Arrow --> Move player north
    case "c483":
        g.Player.Move(-1, 0)

    // Left Arrow --> Move player west
    case "c484":
        g.Player.Move(0, -1)

    // Right Arrow --> Move player east
    case "c485":
        g.Player.Move(0, 1)

    // S --> Save game
    case "53":
        if Confirm("Save and Quit? Y/N") {
            g.SaveGame()
            MessageLog.log("Game Saved")
            g.state = "quit"
        }

    // Q --> Quit game
    case "51":
        g.state = "quit"
    }

    // Handle what all of the other actors (e.g. monsters) are current doing.
    g.processAi()
}
