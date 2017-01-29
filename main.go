/*
 * File: main.go
 *
 * Description: Contains the main.go routine.
 */

package main

// Global variable declaration.
var G Game

//
// Main
//
func main() {

    // Let's get (gocurses) started!
    Init()
    defer End()

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
