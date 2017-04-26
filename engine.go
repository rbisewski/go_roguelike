/*
 * File: engine.go
 *
 * Description: Handles basic game engine aspects.
 */
package main

import (
    "encoding/gob"
    "fmt"
    "./gocurses"
    "math/rand"
    "os"
    "time"
    "strconv"
)

// Structure to hold the part of the window where in-game messages are shown.
type log struct {

    // Section of the window where messages are displayed
    pad   *gocurses.Window

    // Line where the messages are added.
    line  int

    // Line where the messages begin.
    dline int
}

// Gameplay screen section (i.e. where the level information is rendered.
var GamePad *gocurses.Window

// Debug window
var debugWindow *gocurses.Window

// In-game message log.
var MessageLog log

// Player stats section.
var StatsWindow *gocurses.Window

// Globals to handle console height and width.
var ConsoleHeight int
var ConsoleWidth int

// Globals to handle resolve the game-pad height / width.
var ScreenHeight int
var ScreenWidth int

// Globals to handle resolve the ASCII dungeon height / width.
var WorldHeight int
var WorldWidth int

//! Function to initialize the game engine.
/*
 * @return   none
 */
func Init() {

    // Setup the screen.
    gocurses.Initscr()

    // Send the terminal a 'break'
    gocurses.Cbreak()

    // Tell console that this isn't going to send data back to bash, etc.
    gocurses.Noecho()

    // Tell console this'll be using the keypad.
    gocurses.Stdscr.Keypad(true)

    // Nullify the curses set.
    gocurses.CursSet(0)

    // No colours? Then give up here via panic()
    if !gocurses.HasColors() {
        panic("Panic: Console does not use colours!")
    }

    // Since we otherwise have colours, go ahead and just run it.
    gocurses.StartColor()

    // Initialize the colours from the ncurses definitions.
    InitColours()

    // Figure out the limits of the provided console.
    ConsoleHeight, ConsoleWidth = gocurses.Getmaxyx()

    // Carve out a section for the gamepad viewscreen.
    ScreenHeight, ScreenWidth = Percent(85, ConsoleHeight),
                                Percent(70, ConsoleWidth)

    // Carve out another section for the stats viewscreen.
    StatsWindow = gocurses.NewWindow(ScreenHeight,
                                     ConsoleWidth-ScreenWidth,
                                     0,
                                     ScreenWidth+1)

    // Assign some space for the debug message section.
    debugWindow = gocurses.NewWindow(5,
                                     ConsoleWidth,
                                     ConsoleHeight-1,
                                     1)

    // Need in-game messages for those times when the player runs into the
    // wall or kills a monster, and the like...
    MessageLog.pad = gocurses.NewPad(100, ScreenWidth)

    // When the game starts, generate a seed from the nanosecond time.
    rand.Seed(time.Now().UnixNano())
}


//! Function to initialize the colours needed by gocurses.
/*
 * @return   none
 */
func InitColours() {

    // Initialize a red-black colour pair (for corpses, etc)
    gocurses.InitPair(1, gocurses.COLOR_RED, gocurses.COLOR_BLACK)

    // Initialize a yellow-black colour pair (for walls, etc)
    gocurses.InitPair(2, gocurses.COLOR_YELLOW, gocurses.COLOR_BLACK)

    // Initialize a magenta-black colour pair (for items, etc)
    gocurses.InitPair(3, gocurses.COLOR_MAGENTA, gocurses.COLOR_BLACK)
}

//Sets the GamePad and WH-WW info to the current area in the game object.
//!
/*
 * @param     int     height
 * @param     int     width
 *
 * @return    bool    whether or not the pad could be set
 */
func SetPad(h, w int) bool {

    // Input validation, make sure this is greater than 0.
    if h < 1 || w < 1 {
        DebugLog(&G, fmt.Sprintf("SetPad() --> invalid input"))
        return false
    }

    // Initialize a new game pad based on the provided height / width.
    GamePad = gocurses.NewPad(h, w)

    // Sanity check, make sure this actually contains a pointer to a
    // valid Pad object.
    if (GamePad == nil) {
        DebugLog(&G, fmt.Sprintf("SetPad() --> invalid input"))
        return false
    }

    // Define the global world height / width.
    WorldHeight = h
    WorldWidth  = w

    // Return true here since the game pad has been properly defined.
    return true
}

//! Send the end() ncurse to this game.
/*
 * @return    none
 */
func End() {
    gocurses.End()
}

//! Send the clear() ncurse to this game.
/*
 * @return    none
 */
func Clear() {
    gocurses.Clear()
}

//! Function to draw a rune at a given (x,y) point.
/*
 * @param     int      y-value
 * @param     int      x-value
 * @param     rune     ASCII character representation 
 *
 * @return    none
 */
func Draw(y, x int, ch rune) {

    // Draw the aforementioned character.
    GamePad.Mvaddch(int(y), int(x), ch)
}

//! Draw a given ASCII character, with the defined colour.
/*
 * @param     int     y-value
 * @param     int     x-value
 * @param     rune    ASCII character graphic
 * @param     int     colour value
 *
 * @return    none
 */
func DrawColours(y, x int, ch rune, col int) {

    // Apply a colour filter to the character drawing.
    GamePad.Attron(gocurses.ColorPair(col))

    // Add the character to the specific location.
    GamePad.Mvaddch(int(y), int(x), ch)

    // Revert the given filter back to the original console colours afterwards.
    GamePad.Attroff(gocurses.ColorPair(col))
}

//! Given an Area object, attempt to draw a game level.
/*
 * @param     Area*    pointer to an Area object
 *
 * @return    bool     whether or not the map draw action succeeded.
 */
func DrawMap(a *Area) bool {

    // Input validation, make sure this actually got an area object.
    if a == nil {
        DebugLog(&G, fmt.Sprintf("DrawMap() --> invalid input"))
        return false
    }

    // Cycle thru all of the elements via height...
    for y := 0; y < a.Height; y++ {

        // Cycle thru all of the elements via width...
        for x := 0; x < a.Width; x++ {

            // Draw the walls in a brownish / yellow colour.
            if (a.Tiles[x+y*a.Width].Ch == '#') {
                DrawColours(y, x, a.Tiles[x+y*a.Width].Ch, 2)
                continue
            }

            // Else just take the character given and draw it onto the
            // gamepad viewscreen.
            Draw(y, x, a.Tiles[x+y*a.Width].Ch)
        }
    }

    // With all of the characters drawn, this worked as intended.
    return true
}

//! Function to redraw a given gamepad.
/*
 * @param     int    y-value
 * @param     int    x-value
 *
 * @return    none
 */
func RefreshPad(y int, x int) {

    // Determine the relative Y from the overall screen height.
    fromY := Max(0, y-ScreenHeight/2)

    // Determine the relative X from the overall screen width.
    fromX := Max(0, x-ScreenWidth/2)

    // Align the y-value based on the overall height of the screen AND the
    // world. This is done for the purposes of giving a "camera-like".
    if bottomPoint := fromY + ScreenHeight; bottomPoint >= WorldHeight {
        fromY = (WorldHeight - ScreenHeight)
    }

    // Align the x-value based on the overall height of the screen AND the
    // world. This is done for the purposes of giving a "camera-like".
    if rightmostPoint := fromX + ScreenWidth; rightmostPoint >= WorldWidth {
        fromX = (WorldWidth - ScreenWidth)
    }

    // Refresh the output given to the stdout pointer.
    GamePad.PnoutRefresh(fromY, fromX, 0, 0, ScreenHeight-1, ScreenWidth-1)
}

//! Send the necessary ASCII characters into the console via gocurses.
/*
 * @param     int       y-value
 * @param     int       x-value
 * @param     string    ASCII character array (e.g. a "string")
 *
 * @return    none
 */
func Write(y int, x int, s string) {
    gocurses.Mvaddstr(y, x, s)
}

//! Function to write output messages to the debug viewscreen.
/*
 * @param     string    debug log output.
 *
 * @return    none
 */
func DebugLog(g *Game, s string) {

    // If debug mode is off (i.e. false) then do nothing here.
    if g == nil || !g.DebugMode {
        return
    }

    // Add some " " buffers to the character pad.
    debugWindow.Mvaddstr(0, 0, "                         ")

    // Add the given output message.
    debugWindow.Mvaddstr(0, 0, s)

    // Refresh
    debugWindow.NoutRefresh()
}

//! Function to write data to the in-game log screen
/*
 * @param     string    line of data to log.
 *
 * @return    none
 */
func (l *log) log(s string) {

    // Input validation.
    if len(s) < 1 {
        return
    }

    // If the string is less than 60 characters, add space buffers. This
    // is done for certain terminals that "refresh" by writing over existing
    // character buffers, thus preventing a mish-mash of message data.
    for i := len(s); i < 60; i++ {

        // Append a whitespace character to the end of the string.
        s += " "
    }

    // If the string is greater than 60 characters, go ahead and trim it
    // down to a maximum of 57 with an ellipse.
    if len(s) > 60 {

        // Define a temp string variable.
        tmp := ""

        // Grab the first 57 characters of the string, and append them to
        // the temp variable.
        for i := 0; i < 57; i++ {

            // Golang reads string addresses as bytes, so it needs to be
            // recast back to a string type after grabbing the [] address.
            tmp += string(s[i])
        }

        // Dump the concat'd string with ellipse from the tmp into the
        // original string variable.
        s = tmp + "..."
    }

    // Format and write the string.
    l.pad.Mvaddstr(l.line,
                   0,
                   fmt.Sprintf("%s", s))

    // Refresh the screen to account for the newly added log message.
    l.pad.PnoutRefresh(l.dline,
                       0,
                       ScreenHeight+1,
                       0,
                       ConsoleHeight-2,
                       ConsoleWidth)

    // Checks if we need to scroll the window
    if l.line >= ((ConsoleHeight - 2) - (ScreenHeight + 1)) {
        l.dline++
    }

    // Checks if we need to start over on the log.
    if l.line >= 100 {

        // Reset the start and current line back to zero.
        l.line = 0
        l.dline = 0

        // All done here.
        return
    }

    // If everything is good then just move on to the next line.
    l.line++
}

//! Adjust the stats viewscreen of the player.
/*
 * @param    Creature*   pointer to creature object that defines the player
 *
 * @return   none
 */
func (p *Creature) UpdateStats() {

    // Safety check, make sure this is actually the player and that it
    // has a name.
    if len(p.name) < 1 || p.species != "player" {

        // Since this failed due to not being the player, end this function.
        return
    }

    // Print out the name of the player character.
    StatsWindow.Mvaddstr(1, 0, fmt.Sprintf("%s", p.name))

    // Format and write the HP row in the Stats viewscreen.
    //
    // NOTE: several whitespaces were added here to ensure ncurses properly
    //       wipes away and remaining ASCII data from long hitpoints, etc
    //
    StatsWindow.Mvaddstr(3, 0, fmt.Sprintf("HP: %d / %d    ", p.Hp, p.MaxHp))

    // Refresh the screen.
    StatsWindow.NoutRefresh()
}

//! Grab the keyboard input and pass back a string.
/*
 * @return    string    Keyboard ASCII character input.
 */
func GetInput() string {

    // Update the current environment.
    gocurses.Doupdate()

    // Dump the keyboard input to a string and then pass it back.
    return string(gocurses.Getch())
}

//! Display a message asking end-user for y/N confirmation.
/*
 * @param     string    message to display on-screen
 *
 * @return    bool      whether confirmed or denied
 */
func Confirm(msg string) bool {

    // Input validation, make sure the string is between 1 to 30 characters.
    if len(msg) < 1 || len(msg) > 30 {
        return false
    }

    // Variable declaration.
    var GuiSize      = len(msg) + 2
    var GuiTopBottom = "+"
    var GuiLeftRight = "|"

    // Assemble the various parts of the GUI.
    for i := 0; i < GuiSize; i++ {
        GuiTopBottom += "-"
        GuiLeftRight += " "
    }
    GuiTopBottom += "+"
    GuiLeftRight += "|"

    // Write the confirmation message to the screen.
    Write((ScreenHeight/2)-2, ScreenWidth/2, GuiTopBottom)
    Write((ScreenHeight/2)-1, ScreenWidth/2, GuiLeftRight)
    Write(ScreenHeight/2, ScreenWidth/2, "| " + msg + " |")
    Write((ScreenHeight/2)+1, ScreenWidth/2, GuiLeftRight)
    Write((ScreenHeight/2)+2, ScreenWidth/2, GuiTopBottom)

    // Take a look at the keyboard input...
    key := GetInput()

    // End-user pressed Y/y? Go ahead and consider that as confirmation!
    if key == "Y" || key == "y" {
        return true
    }

    // Otherwise the end-user pressed some other key, so close the
    // ncurses confirmation UI.
    return false
}

//! Display the items currently present on the ground.
/*
 * @param     Game*    pointer to the current game object
 *
 * @return    none
 */
func DrawGroundItemsUI(g *Game) {

    // Input validation
    if g == nil || g.Player == nil {
        return
    }

    // Variable declaration
    var itemsAtCurrentCoord = make([]*Item,0)
    var GuiHeight    = 0
    var GuiWidth     = 30
    var GuiTopBottom = "+"
    var GuiLeftRight = "|"
    var GuiLines     = make([]string,0)

    // Obtain the (x,y) coord of where the player character is currently
    // standing.
    loc_x := g.Player.X
    loc_y := g.Player.Y

    // Grab the list of items from the Area te player is currently in and
    // see if they are present in the same coord.
    for _, itm := range g.Area.Items {

        // If the item is at Player (x,y) position, add it to the list of
        // items present on the ground; i.e. itemsAtCurrentCoord
        if loc_x == itm.X && loc_y == itm.Y {
            itemsAtCurrentCoord = append(itemsAtCurrentCoord, itm)
        }
    }

    // Assemble the various parts of the GUI.
    for i := 0; i < GuiWidth; i++ {
        GuiTopBottom += "-"
        GuiLeftRight += " "
    }
    GuiTopBottom += "+"
    GuiLeftRight += "|"

    // Assemble the top element of the inventory screen that displays the
    // words "Items on the Ground" surrounded by '-' characters.
    GuiLines = append(GuiLines, GuiTopBottom)
    GuiLines = append(GuiLines, GuiLeftRight)
    GuiLines = append(GuiLines,
      "| " + AlignAndSpaceString("Items on the Ground", "centre", 9) + " |")
    GuiLines = append(GuiLines, GuiLeftRight)
    GuiLines = append(GuiLines, GuiTopBottom)
    GuiLines = append(GuiLines, GuiLeftRight)

    // Generate a ncurses UI here based on the number of items on the
    // ground.
    for i, itm := range itemsAtCurrentCoord {

        // Append the item with spacing
        GuiLines = append(GuiLines, GuiLeftRight)
        GuiLines = append(GuiLines,
          "| " + AlignAndSpaceString(strconv.Itoa(i) + ") " +
          itm.name, "centre", 9) + " |")
        GuiLines = append(GuiLines, GuiLeftRight)

        // If there are 7 or more items, create a pagination to allow the
        // end-user to cycle thru all of the items on the ground
        if i >= 7 {

            //
            // TODO: implement the below pseudo-code
            //

            // determine the page number and total pages

            // append it to the bottom of the page

            // end the loop since this will only render 7
            break
        }
    }

    // If there are no items on the ground, display a small message
    // stating that there are no items here.
    if len(itemsAtCurrentCoord) < 1 {
        GuiLines = append(GuiLines, GuiLeftRight)
        GuiLines = append(GuiLines,
          "| " + AlignAndSpaceString("No items are here.", "centre", 10) + " |")
        GuiLines = append(GuiLines, GuiLeftRight)
    }

    // Assemble the bottom portion of the ground items UI.
    GuiLines = append(GuiLines, GuiLeftRight)
    GuiLines = append(GuiLines, GuiTopBottom)

    // Get the current number of lines and store it as the height of the UI.
    GuiHeight = len(GuiLines)

    // Using the calculated height, go ahead and determine the upper bounds
    // of the ground items interface, as it relates to the currently drawn
    // ncurses window.
    offset := int(GuiHeight / 2)

    // Safety check, this shouldn't happen but to safe-guard console offsets,
    // if the calculated height is less than one or the offset is zero, tell
    // the developer what happened and leave this function.
    if GuiHeight < 1 || offset == 0 {
        DebugLog(&G,"DrawGroundItemsUI() --> improper height and offset, " +
                    "terminating function")
        return
    }

    // Write the ground item UI to the screen.
    for _, line := range GuiLines {

        // Write the given line to the console output.
        Write((ScreenHeight/2)-offset, ScreenWidth/2, line)

        // Decrement the offset.
        offset--
    }

    // Leave this function, since this needs to redraw the UI if the player
    // picks up all of the items on the ground, etc.
    return
}

//! Display the equipment the character currently is wearing and what items
//! they are holding.
/*
 * @param     Game*    pointer to the current game object
 *
 * @return    none
 */
func DrawInventoryUI(g *Game) {

    // Input validation
    if g == nil || g.Player == nil {
        return
    }

    // Variable declaration.
    var GuiHeight    = 0
    var GuiWidth     = 30
    var GuiTopBottom = "+"
    var GuiLeftRight = "|"
    var GuiLines     = make([]string,0)
    var offset       = 0
    var HeadItem     = "nothing"
    var NeckItem     = "nothing"
    var TorsoItem    = "nothing"
    var RHandItem    = "nothing"
    var LHandItem    = "nothing"
    var PantsItem    = "nothing"

    // If the player has equipment, go ahead a grab the name of the item.
    if g.Player.Head != nil {
        HeadItem  = g.Player.Head.name
    }
    if g.Player.Neck != nil {
        NeckItem  = g.Player.Neck.name
    }
    if g.Player.Torso != nil {
        TorsoItem  = g.Player.Torso.name
    }
    if g.Player.RightHand != nil {
        RHandItem  = g.Player.RightHand.name
    }
    if g.Player.LeftHand != nil {
        LHandItem  = g.Player.LeftHand.name
    }
    if g.Player.Pants != nil {
        PantsItem  = g.Player.Pants.name
    }

    // Assemble the various parts of the GUI.
    for i := 0; i < GuiWidth; i++ {
        GuiTopBottom += "-"
        GuiLeftRight += " "
    }
    GuiTopBottom += "+"
    GuiLeftRight += "|"

    // Assemble the top element of the inventory screen that displays the
    // words "Equipped Items" surrounded by '-' characters.
    GuiLines = append(GuiLines, GuiTopBottom)
    GuiLines = append(GuiLines, GuiLeftRight)
    GuiLines = append(GuiLines,
      "| " + AlignAndSpaceString("Equipped Items", "centre", 14) + " |")
    GuiLines = append(GuiLines, GuiLeftRight)
    GuiLines = append(GuiLines, GuiTopBottom)
    GuiLines = append(GuiLines, GuiLeftRight)

    // Assemble the portion of the inventory screen for the head item.
    GuiLines = append(GuiLines,
      "| Head       --> " + AlignAndSpaceString(HeadItem, "right", 13) + " |")

    // Add a spacer.
    GuiLines = append(GuiLines, GuiLeftRight)

    // Assemble the portion of the inventory screen for the neck item.
    GuiLines = append(GuiLines,
      "| Neck       --> " + AlignAndSpaceString(NeckItem, "right", 13) + " |")

    // Add a spacer.
    GuiLines = append(GuiLines, GuiLeftRight)

    // Assemble the portion of the inventory screen for the torso item.
    GuiLines = append(GuiLines,
      "| Torso      --> " + AlignAndSpaceString(TorsoItem, "right", 13) + " |")

    // Add a spacer.
    GuiLines = append(GuiLines, GuiLeftRight)

    // Assemble the portion of the inventory screen for the right hand item.
    GuiLines = append(GuiLines,
      "| Right Hand --> " + AlignAndSpaceString(RHandItem, "right", 13) + " |")

    // Add a spacer.
    GuiLines = append(GuiLines, GuiLeftRight)

    // Assemble the portion of the inventory screen for the left hand item.
    GuiLines = append(GuiLines,
      "| Left Hand  --> " + AlignAndSpaceString(LHandItem, "right", 13) + " |")

    // Add a spacer.
    GuiLines = append(GuiLines, GuiLeftRight)

    // Assemble the portion of the inventory screen for the pants item.
    GuiLines = append(GuiLines,
      "| Pants      --> " + AlignAndSpaceString(PantsItem, "right", 13) + " |")

    // Assemble the bottom portion of the inventory.
    GuiLines = append(GuiLines, GuiLeftRight)
    GuiLines = append(GuiLines, GuiTopBottom)

    // Get the current number of lines and store it as the height of the UI.
    GuiHeight = len(GuiLines)

    // Using the calculated height, go ahead and determine the upper bounds
    // of the inventory interface, as it relates to the currently drawn
    // ncurses window.
    offset = int(GuiHeight / 2)

    // Safety check, this shouldn't happen but to safe-guard console offsets,
    // if the calculated height is less than one or the offset is zero, tell
    // the developer what happened and leave this function.
    if GuiHeight < 1 || offset == 0 {
        DebugLog(&G,"ToggleInventoryUI() --> improper height and offset, " +
                    "terminating function")
        return
    }

    // Write the character equipment/inventory screen.
    for _, line := range GuiLines {

        // Write the given line to the console output.
        Write((ScreenHeight/2)-offset, ScreenWidth/2, line)

        // Decrement the offset.
        offset--
    }

    // All done here, so then this can return.
    return
}

//! Handles a "save game to disk" event.
/*
 * @param     Game*    pointer to the current game instance.
 *
 * @return    none
 */
func (g *Game) SaveGame() {

    // Attempt to open the saved game.
    file, err := os.OpenFile("player.sav", os.O_WRONLY|os.O_CREATE, 0600)

    // Error? Print it to stdout and kill the program.
    if err != nil {
        panic(err)
    }

    // Take into account possible file issues via this.
    defer func() {
        if err := file.Close(); err != nil {
            panic(err)
        }
    }()

    // Attempt to prepare to encode the save game file.
    encoder := gob.NewEncoder(file)

    // Go ahead and encode now...
    err = encoder.Encode(g)

    // Error? Print it to stdout and kill the program.
    if err != nil {
        panic(err)
    }
}

//! Handles a "load game from disk" event.
/*
 * @param     Game*    pointer to the current game instance.
 *
 * @return    bool     whether or not the load was successful
 */
func (g *Game) LoadGame(filename string) bool {

    // Input validation, make sure the filename is valid.
    if len(filename) < 1 {
        return false
    }

    // Safety check, attempt to stat() the given filename.
    _, err := os.Stat(filename)

    // If an error occurred here, either the filename does not exist or
    // is inaccessible, simply return false to allow the game to continue
    // without terminating.
    if err != nil {
        return false
    }

    // Loading requires that this open an existing save game, so do that.
    // Probably only a read lock is needed here, so that's all that is used.
    file, err := os.OpenFile(filename, os.O_RDONLY, 0600)

    // If an error occurs at this point, terminate the program since
    // probably a memory or library error has occurred.
    if err != nil {
        panic(err)
    }

    // If closing the file causes unforeseen consequences, go ahead and
    // terminate the program.
    if file.Close() != nil {
        panic(err)
    }

    // Prepare to decode the file in question.
    decoder := gob.NewDecoder(file)

    // Go ahead and decode the loaded file.
    err = decoder.Decode(g)

    // Error occurred during decoding? Terminate the program via panic()
    if err != nil {
        panic(err)
    }

    // Otherwise everything loads as intended, so return true.
    return true
}
