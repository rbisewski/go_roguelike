/*
 * File: engine.go
 *
 * Description: Handles basic game engine aspects.
 */

package main

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/rbisewski/gocurses"
)

// log ... Holds the part of the window where in-game messages are shown.
type log struct {

	// Section of the window where messages are displayed
	pad *gocurses.Window

	// Line where the messages are added.
	line int

	// Line where the messages begin.
	dline int
}

// GamePad ... gameplay screen section where the level information is rendered
var GamePad *gocurses.Window

// debugWindow ... window for debug information
var debugWindow *gocurses.Window

// MessageLog ... in-game message log
var MessageLog log

// StatsWindow ... player stats section
var StatsWindow *gocurses.Window

// ConsoleHeight ... global to handle console height
var ConsoleHeight int

// ConsoleWidth ... global to handle console height
var ConsoleWidth int

// ScreenHeight ... gamepad height
var ScreenHeight int

// ScreenWidth ... gamepad width
var ScreenWidth int

// WorldHeight ... dungeon map height
var WorldHeight int

// WorldWidth ... dungeon map width
var WorldWidth int

// Init ... Function to initialize the game engine.
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

// InitColours ... Function to initialize the colours needed by gocurses.
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

// SetPad ... sets GamePad / WH-WW info to current area in the game object
/*
 * @param     int     height
 * @param     int     width
 *
 * @return    bool    whether or not the pad could be set
 */
func SetPad(h, w int) bool {

	if h < 1 || w < 1 {
		DebugLog(&G, fmt.Sprintf("SetPad() --> invalid input"))
		return false
	}

	// Initialize a new game pad based on the provided height / width.
	GamePad = gocurses.NewPad(h, w)
	if GamePad == nil {
		DebugLog(&G, fmt.Sprintf("SetPad() --> invalid input"))
		return false
	}

	WorldHeight = h
	WorldWidth = w

	return true
}

// End ... send the end() ncurse to this game.
/*
 * @return    none
 */
func End() {
	gocurses.End()
}

// Clear ... send the clear() ncurse to this game.
/*
 * @return    none
 */
func Clear() {
	gocurses.Clear()
}

// Draw ... function to draw a rune at a given (x,y) point.
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

// DrawColours ... draw a given ASCII character, with the defined colour.
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

// DrawMap ... given an Area object, attempt to draw a game level.
/*
 * @param     Area*    pointer to an Area object
 *
 * @return    bool     whether or not the map draw action succeeded.
 */
func DrawMap(a *Area) bool {

	if a == nil {
		DebugLog(&G, fmt.Sprintf("DrawMap() --> invalid input"))
		return false
	}

	// Cycle thru all of the elements via height...
	for y := 0; y < a.Height; y++ {

		// Cycle thru all of the elements via width...
		for x := 0; x < a.Width; x++ {

			// Draw the walls in a brownish / yellow colour.
			if a.Tiles[x+y*a.Width].Ch == '#' {
				DrawColours(y, x, a.Tiles[x+y*a.Width].Ch, 2)
				continue
			}

			// Else just take the character given and draw it onto the
			// gamepad viewscreen.
			Draw(y, x, a.Tiles[x+y*a.Width].Ch)
		}
	}
	return true
}

// RefreshPad ... function to redraw a given gamepad.
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

// Write ... send the necessary ASCII characters into console via gocurses.
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

// DebugLog ... function to write output messages to the debug viewscreen.
/*
 * @param     string    debug log output.
 *
 * @return    none
 */
func DebugLog(g *Game, s string) {

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

// log ... function to write data to the in-game log screen
/*
 * @param     string    line of data to log.
 *
 * @return    none
 */
func (l *log) log(s string) {

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
	// down to a maximum of 57, then end with an ellipse.
	if len(s) > 60 {

		tmp := ""
		for i := 0; i < 57; i++ {
			tmp += string(s[i])
		}

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

	// the log stores 99 lines, if it reaches 100, reset back to zero
	if l.line >= 100 {
		l.line = 0
		l.dline = 0
		return
	}

	l.line++
}

// UpdateStats ... Adjust the stats viewscreen of the player.
/*
 * @param    Creature*   pointer to creature object that defines the player
 *
 * @return   none
 */
func (p *Creature) UpdateStats() {

	if len(p.name) < 1 || p.species != "player" {
		return
	}

	// Print out the name of the player character.
	StatsWindow.Mvaddstr(1, 0, fmt.Sprintf("%s", p.name))

	// Print out the class of the character.
	StatsWindow.Mvaddstr(3, 0, fmt.Sprintf("%s", p.class.Name))

	// Format and write the HP row in the Stats viewscreen.
	//
	// NOTE: several whitespaces were added here to ensure ncurses properly
	//       wipes away and remaining ASCII data from long hitpoints, etc
	//
	StatsWindow.Mvaddstr(5, 0, fmt.Sprintf("HP: %d / %d    ", p.Hp, p.MaxHp))

	// Print out the four primary attributes; strength, intelligence,
	// agility, and wisdom.
	StatsWindow.Mvaddstr(7, 0, fmt.Sprintf("Strength:     %d ",
		p.Strength))
	StatsWindow.Mvaddstr(8, 0, fmt.Sprintf("Intelligence: %d ",
		p.Intelligence))
	StatsWindow.Mvaddstr(9, 0, fmt.Sprintf("Agility:      %d ",
		p.Agility))
	StatsWindow.Mvaddstr(10, 0, fmt.Sprintf("Wisdom:       %d ",
		p.Wisdom))

	// Refresh the screen.
	StatsWindow.NoutRefresh()
}

// GetInput ... Grab the keyboard input and then pass back a string.
/*
 * @return    string    Keyboard ASCII character input (Getch() = get character)
 */
func GetInput() string {
	gocurses.Doupdate()
	return string(gocurses.Getch())
}

// Confirm ... Display a message asking end-user for y/N confirmation.
/*
 * @param     string    message to display on-screen
 *
 * @return    bool      whether confirmed or denied
 */
func Confirm(msg string) bool {

	if len(msg) < 1 || len(msg) > 30 {
		return false
	}

	var GuiSize = len(msg) + 2
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
	Write(ScreenHeight/2, ScreenWidth/2, "| "+msg+" |")
	Write((ScreenHeight/2)+1, ScreenWidth/2, GuiLeftRight)
	Write((ScreenHeight/2)+2, ScreenWidth/2, GuiTopBottom)

	// Take a look at the keyboard input...
	key := GetInput()

	// End-user pressed Y/y? Go ahead and consider that as confirmation!
	if key == "Y" || key == "y" {
		return true
	}

	return false
}

// PickupGroundItem ... pickup an item from the list of ground items
/*
 * @param     Game*    pointer to the current game object
 * @param     string   the given key that was pressed
 *
 * @return    error    error message, if any
 */
func PickupGroundItem(g *Game, keyPressed string) error {

	if g == nil || len(keyPressed) < 1 {
		return fmt.Errorf("PickupGroundItem() --> invalid input")
	}

	var givenItem *Item

	// If there is less than 1 item, go back.
	if len(g.GroundItems) < 1 {
		return nil
	}

	// Do a switch to check if a key between 1-6 was pressed, and
	// grab that item.
	switch keyPressed {

	// Number 1
	case "31":

		// Safety check, ensure there is at least 1 item.
		if len(g.GroundItems) < 1 {
			return nil
		}

		// Grab the 1st item.
		givenItem = g.GroundItems[0]

	// Number 2
	case "32":

		// Safety check, ensure there is at least 2 items.
		if len(g.GroundItems) < 2 {
			return nil
		}

		// Add code to grab the 2nd item.
		givenItem = g.GroundItems[1]

	// Number 3
	case "33":

		// Safety check, ensure there is at least 3 items.
		if len(g.GroundItems) < 3 {
			return nil
		}

		// Add code to grab the 3rd item.
		givenItem = g.GroundItems[2]

	// Number 4
	case "34":

		// Safety check, ensure there is at least 4 items.
		if len(g.GroundItems) < 4 {
			return nil
		}

		// Add code to grab the 4th item.
		givenItem = g.GroundItems[3]

	// Number 5
	case "35":

		// Safety check, ensure there is at least 5 items.
		if len(g.GroundItems) < 5 {
			return nil
		}

		// Add code to grab the 5th item.
		givenItem = g.GroundItems[4]

	// Number 6
	case "36":

		// Safety check, ensure there is at least 6 items.
		if len(g.GroundItems) < 6 {
			return nil
		}

		// Add code to grab the 6th item.
		givenItem = g.GroundItems[5]
	}

	// If the item is nil, then skip this step.
	if givenItem == nil {
		return nil
	}

	// Set the current area of that item to nil.
	givenItem.area = nil

	// Attempt to add that item to the player's inventory.
	g.Player.inventory = append(g.Player.inventory, givenItem)

	// Delete the item from the current area.
	for index, item := range g.Area.Items {

		// if the item matches, remove it
		if item == givenItem {

			// Slick trick --> remove an element from an array while
			// preserving order.
			copy(g.Area.Items[index:], g.Area.Items[index+1:])
			g.Area.Items[len(g.Area.Items)-1] = nil
			g.Area.Items = g.Area.Items[:len(g.Area.Items)-1]
			break
		}
	}

	return nil
}

// DrawGroundItemsUI ... display the items currently present on the ground.
/*
 * @param     Game*    pointer to the current game object
 * @param     string   key pressed, as a string
 *
 * @return    none
 */
func DrawGroundItemsUI(g *Game, key string) {

	if g == nil || g.Player == nil || len(key) < 1 {
		return
	}

	var GuiHeight = 0
	var GuiWidth = 30
	var GuiTopBottom = "+"
	var GuiLeftRight = "|"
	var GuiLines = make([]string, 0)

	// Obtain the (x,y) coord of where the player character is currently
	// standing.
	locX := g.Player.X
	locY := g.Player.Y

	// Remake the ground items array.
	g.GroundItems = make([]*Item, 0)

	// Grab the list of items from the Area te player is currently in and
	// see if they are present in the same coord.
	for _, itm := range g.Area.Items {

		// If the item is at Player (x,y) position, add it to the list of
		// items present on the ground; i.e. itemsAtCurrentCoord
		if locX == itm.X && locY == itm.Y {
			g.GroundItems = append(g.GroundItems, itm)
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
		"| "+AlignAndSpaceString("Items on the Ground", "centre", 9)+" |")
	GuiLines = append(GuiLines, GuiLeftRight)
	GuiLines = append(GuiLines, GuiTopBottom)

	// Variable to store the current page.
	currentPage := 1
	numOfPages := int(len(g.GroundItems) / 6)
	itemPrintedCounter := 0

	// Generate a ncurses UI here based on the number of items on the
	// ground; display the items of the relevant inventory page.
	//
	// TODO: add logic to increment the current page
	//
	for i, itm := range g.GroundItems {

		// Skip elements of a forward or backward page.
		if i < ((currentPage - 1) * 6) {
			continue
		}

		// Append the item with spacing
		GuiLines = append(GuiLines, GuiLeftRight)
		GuiLines = append(GuiLines,
			"| "+AlignAndSpaceString(strconv.Itoa(i+1)+") "+
				itm.name, "right", GuiWidth-2)+" |")

		// Increment the current number of items printed
		itemPrintedCounter++

		// If there are 6 or more items, create a pagination to allow the
		// end-user to cycle thru all of the items on the ground
		if itemPrintedCounter >= 6 && numOfPages > 1 {

			// assemble the text for the 'Page x of y' label
			pageLabel := "Page " + strconv.Itoa(currentPage) + " of " +
				strconv.Itoa(numOfPages)

			// append it to the bottom of the page
			GuiLines = append(GuiLines, GuiLeftRight)
			GuiLines = append(GuiLines,
				"| "+AlignAndSpaceString(pageLabel, "right",
					GuiWidth-2)+" |")
			GuiLines = append(GuiLines, GuiLeftRight)

			// end the loop since this will only render 6
			break
		}
	}

	// If there are no items on the ground, display a small message
	// stating that there are no items here.
	if len(g.GroundItems) < 1 {
		GuiLines = append(GuiLines, GuiLeftRight)
		GuiLines = append(GuiLines,
			"| "+AlignAndSpaceString("No items are here.", "centre", 10)+" |")
		GuiLines = append(GuiLines, GuiLeftRight)
	}

	// Get the current number of lines and store it as the height of the UI.
	GuiHeight = len(GuiLines)

	// While the UI height is less than 17, keep appending |_| lines.
	for GuiHeight != 17 {
		GuiLines = append(GuiLines, GuiLeftRight)
		GuiHeight = len(GuiLines)
	}

	// Assemble the bottom portion of the ground items UI.
	GuiLines = append(GuiLines, GuiLeftRight)
	GuiLines = append(GuiLines, GuiTopBottom)

	// Using the calculated height, go ahead and determine the upper bounds
	// of the ground items interface, as it relates to the currently drawn
	// ncurses window.
	offset := int(GuiHeight/2) + 1

	// Safety check, this shouldn't happen but to safe-guard console offsets,
	// if the calculated height is less than one or the offset is zero, tell
	// the developer what happened and leave this function.
	if GuiHeight < 1 || offset == 0 {
		DebugLog(&G, "DrawGroundItemsUI() --> improper height and offset, "+
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

// DrawInventoryUI ... display the inventory the character currently
// has in their backpack.
/*
 * @param     Game*    pointer to the current game object
 * @param     string   key pressed, as a string
 *
 * @return    none
 */
func DrawInventoryUI(g *Game, key string) {

	if g == nil || g.Player == nil || len(key) < 1 {
		return
	}

	var GuiHeight = 0
	var GuiWidth = 30
	var GuiTopBottom = "+"
	var GuiLeftRight = "|"
	var GuiLines = make([]string, 0)

	// Obtain the (x,y) coord of where the player character is currently
	// standing.
	locX := g.Player.X
	locY := g.Player.Y

	// Remake the ground items array.
	g.GroundItems = make([]*Item, 0)

	// Grab the list of items from the Area te player is currently in and
	// see if they are present in the same coord. This will be useful for
	// when the player wants to drop an item.
	for _, itm := range g.Area.Items {

		// If the item is at Player (x,y) position, add it to the list of
		// items present on the ground; i.e. itemsAtCurrentCoord
		if locX == itm.X && locY == itm.Y {
			g.GroundItems = append(g.GroundItems, itm)
		}
	}

	// Assemble the various parts of the GUI.
	for i := 0; i < GuiWidth; i++ {
		GuiTopBottom += "-"
		GuiLeftRight += " "
	}
	GuiTopBottom += "+"
	GuiLeftRight += "|"

	// Assemble the inventory screen header.
	GuiLines = append(GuiLines, GuiTopBottom)
	GuiLines = append(GuiLines, GuiLeftRight)
	GuiLines = append(GuiLines,
		"| "+AlignAndSpaceString("Inventory", "centre", 19)+" |")
	GuiLines = append(GuiLines, GuiLeftRight)
	GuiLines = append(GuiLines, GuiTopBottom)

	// Variable to store the current page.
	currentPage := 1
	numOfPages := int(len(g.Player.inventory) / 6)
	itemPrintedCounter := 0

	// Generate a ncurses UI here based on the number of items on the
	// ground; display the items of the relevant inventory page.
	//
	// TODO: add logic to increment the current page
	//
	for i, itm := range g.Player.inventory {

		// Skip elements of a forward or backward page.
		if i < ((currentPage - 1) * 6) {
			continue
		}

		// Append the item with spacing
		GuiLines = append(GuiLines, GuiLeftRight)
		GuiLines = append(GuiLines,
			"| "+AlignAndSpaceString(strconv.Itoa(i+1)+") "+
				itm.name, "right", GuiWidth-2)+" |")

		// Increment the current number of items printed
		itemPrintedCounter++

		// If there are 6 or more items, create a pagination to allow the
		// end-user to cycle thru all of the items on the ground
		if itemPrintedCounter >= 6 && numOfPages > 1 {

			// assemble the text for the 'Page x of y' label
			pageLabel := "Page " + strconv.Itoa(currentPage) + " of " +
				strconv.Itoa(numOfPages)

			// append it to the bottom of the page
			GuiLines = append(GuiLines, GuiLeftRight)
			GuiLines = append(GuiLines,
				"| "+AlignAndSpaceString(pageLabel, "right",
					GuiWidth-2)+" |")
			GuiLines = append(GuiLines, GuiLeftRight)

			// end the loop since this will only render 7
			break
		}
	}

	// If the player is currently holding no items, go ahead and mention
	// that the backpack of the player is empty.
	if len(g.Player.inventory) < 1 {
		GuiLines = append(GuiLines, GuiLeftRight)
		GuiLines = append(GuiLines,
			"| "+AlignAndSpaceString("Backpack is empty.", "centre", 10)+" |")
		GuiLines = append(GuiLines, GuiLeftRight)
	}

	// Get the current number of lines and store it as the height of the UI.
	GuiHeight = len(GuiLines)

	// While the UI height is less than 17, keep appending |_| lines.
	for GuiHeight < 17 {
		GuiLines = append(GuiLines, GuiLeftRight)
		GuiHeight = len(GuiLines)
	}

	// Assemble the bottom portion of the ground items UI.
	GuiLines = append(GuiLines, GuiLeftRight)
	GuiLines = append(GuiLines, GuiTopBottom)

	// Using the calculated height, go ahead and determine the upper bounds
	// of the ground items interface, as it relates to the currently drawn
	// ncurses window.
	offset := int(GuiHeight/2) + 1

	// Safety check, this shouldn't happen but to safe-guard console offsets,
	// if the calculated height is less than one or the offset is zero, tell
	// the developer what happened and leave this function.
	if GuiHeight < 1 || offset == 0 {
		DebugLog(&G, "DrawInventoryUI() --> improper height and offset, "+
			"terminating function")
		return
	}

	// Write the ground item UI to the screen.
	for _, line := range GuiLines {
		Write((ScreenHeight/2)-offset, ScreenWidth/2, line)
		offset--
	}
}

// DrawEquipmentUI ... display the equipment the character currently is
// wearing and what items they are holding.
/*
 * @param     Game*    pointer to the current game object
 * @param     string   key pressed, as a string
 *
 * @return    none
 */
func DrawEquipmentUI(g *Game, key string) {

	if g == nil || g.Player == nil || len(key) < 1 {
		return
	}

	var GuiHeight = 0
	var GuiWidth = 30
	var GuiTopBottom = "+"
	var GuiLeftRight = "|"
	var GuiLines = make([]string, 0)
	var offset = 0
	var HeadItem = "nothing"
	var NeckItem = "nothing"
	var TorsoItem = "nothing"
	var RHandItem = "nothing"
	var LHandItem = "nothing"
	var PantsItem = "nothing"

	// If the player has equipment, go ahead a grab the name of the item.
	if g.Player.Head != nil {
		HeadItem = g.Player.Head.name
	}
	if g.Player.Neck != nil {
		NeckItem = g.Player.Neck.name
	}
	if g.Player.Torso != nil {
		TorsoItem = g.Player.Torso.name
	}
	if g.Player.RightHand != nil {
		RHandItem = g.Player.RightHand.name
	}
	if g.Player.LeftHand != nil {
		LHandItem = g.Player.LeftHand.name
	}
	if g.Player.Pants != nil {
		PantsItem = g.Player.Pants.name
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
		"| "+AlignAndSpaceString("Equipped Items", "centre", 14)+" |")
	GuiLines = append(GuiLines, GuiLeftRight)
	GuiLines = append(GuiLines, GuiTopBottom)
	GuiLines = append(GuiLines, GuiLeftRight)

	// Assemble the portion of the inventory screen for the head item.
	GuiLines = append(GuiLines,
		"| Head       --> "+AlignAndSpaceString(HeadItem, "right", 13)+" |")

	// Add a spacer.
	GuiLines = append(GuiLines, GuiLeftRight)

	// Assemble the portion of the inventory screen for the neck item.
	GuiLines = append(GuiLines,
		"| Neck       --> "+AlignAndSpaceString(NeckItem, "right", 13)+" |")

	// Add a spacer.
	GuiLines = append(GuiLines, GuiLeftRight)

	// Assemble the portion of the inventory screen for the torso item.
	GuiLines = append(GuiLines,
		"| Torso      --> "+AlignAndSpaceString(TorsoItem, "right", 13)+" |")

	// Add a spacer.
	GuiLines = append(GuiLines, GuiLeftRight)

	// Assemble the portion of the inventory screen for the right hand item.
	GuiLines = append(GuiLines,
		"| Right Hand --> "+AlignAndSpaceString(RHandItem, "right", 13)+" |")

	// Add a spacer.
	GuiLines = append(GuiLines, GuiLeftRight)

	// Assemble the portion of the inventory screen for the left hand item.
	GuiLines = append(GuiLines,
		"| Left Hand  --> "+AlignAndSpaceString(LHandItem, "right", 13)+" |")

	// Add a spacer.
	GuiLines = append(GuiLines, GuiLeftRight)

	// Assemble the portion of the inventory screen for the pants item.
	GuiLines = append(GuiLines,
		"| Pants      --> "+AlignAndSpaceString(PantsItem, "right", 13)+" |")

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
		DebugLog(&G, "ToggleInventoryUI() --> improper height and offset, "+
			"terminating function")
		return
	}

	// Write the character equipment/inventory screen.
	for _, line := range GuiLines {
		Write((ScreenHeight/2)-offset, ScreenWidth/2, line)
		offset--
	}
}

// SaveGame ... Handles a "save game to disk" event.
/*
 * @param     Game*    pointer to the current game instance.
 *
 * @return    none
 */
func (g *Game) SaveGame() {

	// Attempt to open the saved game.
	file, err := os.OpenFile("player.sav", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	// Attempt to prepare to encode the save game file.
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(g)
	if err != nil {
		panic(err)
	}
}

// LoadGame ... handles a "load game from disk" event.
/*
 * @param     Game*    pointer to the current game instance.
 *
 * @return    bool     whether or not the load was successful
 */
func (g *Game) LoadGame(filename string) bool {

	if filename == "" {
		return false
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		panic(err.Error())
	}

	// Prepare to decode the file in question.
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(g)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}

	return true
}
