/*
 * File: curses.go
 *
 * Description: Contains a number of useful golang wrappers for the ncurses
 *              library. Ergo, this file relies on Cgo to properly access the
 *              char drawing and colouring functionality.
 */

package gocurses

// #cgo LDFLAGS: -lncurses
// #include <stdlib.h>
// #include <ncurses.h>
// int wrapper_getmaxx(WINDOW* win) { return getmaxx(win); }
// int wrapper_getmaxy(WINDOW* win) { return getmaxy(win); }
// void wrapper_wclrtoeol(WINDOW* win) { wclrtoeol(win); }
// void wrapper_wattrset(WINDOW* win, int attr) { wattrset(win, attr); }
// int wrapper_color_pair(int i) { return COLOR_PAIR(i); }
// int wrapper_set_escdelay(int i) { return set_escdelay(i); }
//
import "C"
import "unsafe"
import "fmt"

// Curses window type.
type Window struct {
    cwin *C.WINDOW
}

// Standard window.
var Stdscr *Window = &Window{cwin: C.stdscr}

// Initializes curses.
// This function should be called before using the package.
func Initscr() *Window {
    Stdscr.cwin = C.initscr()
    // set the ESCDELAY to 1 millisecond
    C.wrapper_set_escdelay(1)
    return Stdscr
}

// Raw input. No buffering.
// CTRL+Z and CTRL+C passed to the application.
func Raw() {
    C.raw()
}

// No buffering.
func Cbreak() {
    C.cbreak()
}

// Enable character echoing while reading.
func Echo() {
    C.echo()
}

// Disables character echoing while reading.
func Noecho() {
    C.noecho()
}

// Hides cursor if 0, visible if 1, very visible if 2
func CursSet(i int) {
  C.curs_set(C.int(i))
}

// Starts color capabilities, check with HasColors if terminal has the capability.
func StartColor() {
    C.start_color()
}

// Checks if the terminal supports colors.
func HasColors() bool {
    return bool(C.has_colors())
}

func InitPair(pair, fg, bg int) {
    C.init_pair(C.short(pair), C.short(fg), C.short(bg))
}

//! Determine the intensity of an 8bit terminal colour.
/*
 * @param      int    id of the new colour to generate
 * @param      int    red value of the colour
 * @param      int    green value of the colour
 * @param      int    blue value of the colour
 *
 * @returns    int    If success: colour intensity (as int)
 *                    If failure: -1
 */
func ColorContent(color_id, red, green, blue int) int {

    // Input validation
    if color_id < 0 || red < 0 || green < 0 || blue < 0 ||
      color_id > 1000 || red > 1000 || green > 1000 || blue > 1000 {

        // Return a nil value since something bad has happened.
        return -1
    }

    // Convert the values in question to C-style shorts.
    red_as_short   := C.short(red)
    green_as_short := C.short(green)
    blue_as_short  := C.short(blue)

    // Sanity check, make sure that golang was able to convert them properly.
    if red_as_short < 0 || green_as_short < 0 || blue_as_short < 0 {

        // Return a nil value since something bad has happened.
        return -1
    }

    // Otherwise call ncurses color_content() function and return that value.
    return int(C.color_content(C.short(color_id),
                               &red_as_short,
                               &green_as_short,
                               &blue_as_short))
}

//! Initialize a new 8bit terminal colour.
/*
 * @param      int    id of the new colour to generate
 * @param      int    red value of the colour
 * @param      int    green value of the colour
 * @param      int    blue value of the colour
 *
 * @returns    int    If success: generated colour (as int)
 *                    If failure: -1
 */
func InitColor(color_id, red, green, blue int) int {

    // Input validation
    if color_id < 0 || red < 0 || green < 0 || blue < 0 ||
      color_id > 1000 || red > 1000 || green > 1000 || blue > 1000 {

        // Return a nil value since something bad has happened.
        return -1
    }

    // Convert the values in question to C-style shorts.
    red_as_short   := C.short(red)
    green_as_short := C.short(green)
    blue_as_short  := C.short(blue)

    // Sanity check, make sure that golang was able to convert them properly.
    if red_as_short < 0 || green_as_short < 0 || blue_as_short < 0 {

        // Return a nil value since something bad has happened.
        return -1
    }

    // Otherwise call ncurses color_content() function and return that value.
    return int(C.init_color(C.short(color_id),
                               red_as_short,
                               green_as_short,
                               blue_as_short))
}


// Enable reading of function keys.
func (window *Window) Keypad(on bool) {
    C.keypad(window.cwin, C.bool(on))
}

// Get char from the standard in.
func (window *Window) Getch() int {
    return int(C.wgetch(window.cwin))
}

// Get char from the standard in.
func Getch() int {
    return int(C.getch())
}

// Enable attribute
func Attron(attr int) {
    C.attron(C.int(attr))
}

// Disable attribute
func Attroff(attr int) {
    C.attroff(C.int(attr))
}

// Set attribute
func Attrset(attr int) {
    C.attrset(C.int(attr))
}

func (window *Window) Attron(attr int) {
    C.wattron(window.cwin, C.int(attr))
}

func (window *Window) Attroff(attr int) {
    C.wattroff(window.cwin, C.int(attr))
}

func (window *Window) Attrset(attr int) {
    C.wrapper_wattrset(window.cwin, C.int(attr))
}

func ColorPair(pair int) int {
    return int(C.wrapper_color_pair(C.int(pair)))
}

// Refresh screen.
func Refresh() {
    C.refresh()
}

// Refresh given window.
func (window *Window) Refresh() {
    C.wrefresh(window.cwin)
}

// Refresh given window, for pads.
func (window *Window) PRefresh(pminrow,pmincol,sminrow,smincol, smaxrow, smaxcol int) {
    C.prefresh(window.cwin,C.int(pminrow),C.int(pmincol),C.int(sminrow),C.int(smincol),C.int(smaxrow),C.int(smaxcol))
}

//! Allows for multiple updates with more efficiency than refresh alone.
//!
//! Normally drawing functionality might have to call refresh multiple times.
//! In the case of Noutrefresh, only updates to the virtual screen and then
//! the actual screen can be updated by calling Doupdate, which checks all
//! pending changes.
/*
 * @param     Window*    pointer to a Cwindow
 *
 * @return    none
 */
func (window *Window) NoutRefresh() {
    C.wnoutrefresh(window.cwin)
}

// Same as NoutRefresh, but for pads.
func (window *Window) PnoutRefresh(pminrow,pmincol,sminrow,smincol,smaxrow,smaxcol int) {
    C.pnoutrefresh(window.cwin,C.int(pminrow),C.int(pmincol),C.int(sminrow),C.int(smincol),C.int(smaxrow),C.int(smaxcol))
}

// Compares the virtual screen to the physical screen and does the actual update.
func Doupdate() {
    C.doupdate()
}

// Finalizes curses.
func End() {
    C.endwin()
}

// Create new window.
func NewWindow(height, width, starty, startx int) *Window {
    w := new(Window)
    w.cwin = C.newwin(C.int(height), C.int(width),
        C.int(starty), C.int(startx))
    return w
}

// Create new pad.
func NewPad(nlines int, ncols int) *Window {
    w := new(Window)
    w.cwin = C.newpad(C.int(nlines), C.int(ncols))
    return w
}

// Set box lines.
func (window *Window) Box(v, h int) {
    C.box(window.cwin, C.chtype(v), C.chtype(h))
}

//! Function to set the border characters of a given window.
/*
 * @param    int    character for the left side of the window
 * @param    int    right side of the window
 * @param    int    top side of the window
 * @param    int    bottom side of the window
 * @param    int    top left corner of the window
 * @param    int    top right corner of the window
 * @param    int    bottom left corner of the window
 * @param    int    bottom right corner of the window
 *
 * @return   none
 */
func (window *Window) Border(ls, rs, ts, bs, tl, tr, bl, br int) {

    // Define the window borders using the given characters.
    C.wborder(window.cwin,
              C.chtype(ls),
              C.chtype(rs),
              C.chtype(ts),
              C.chtype(bs),
              C.chtype(tl),
              C.chtype(tr),
              C.chtype(bl),
              C.chtype(br))
}

// Delete current window.
func (window *Window) Del() {
    C.delwin(window.cwin)
}

// Get windows sizes.
func (window *Window) Getmaxyx() (row, col int) {
    row = int(C.wrapper_getmaxy(window.cwin))
    col = int(C.wrapper_getmaxx(window.cwin))
    return row, col
}

func (window *Window) Setscrreg(top, bot int) {
    C.wsetscrreg(window.cwin, C.int(top), C.int(bot))
}

func Addstr(str ...interface{}) {
    res := (*C.char)(C.CString(fmt.Sprint(str...)))
    defer C.free(unsafe.Pointer(res))
    C.addstr(res)
}

func Mvaddstr(y, x int, str ...interface{}) {
    res := (*C.char)(C.CString(fmt.Sprint(str...)))
    defer C.free(unsafe.Pointer(res))
    C.mvaddstr(C.int(y), C.int(x), res)
}

func Addch(ch int) {
	C.addch(C.chtype(ch))
}

func Mvaddch(y, x int, ch rune) {
	C.mvaddch(C.int(y), C.int(x), C.chtype(ch))
}

func (window *Window) Addstr(str ...interface{}) {
    res := (*C.char)(C.CString(fmt.Sprint(str...)))
    defer C.free(unsafe.Pointer(res))
    C.waddstr(window.cwin, res)
}

func (window *Window) Mvaddstr(y, x int, str ...interface{}) {
    res := (*C.char)(C.CString(fmt.Sprint(str...)))
    defer C.free(unsafe.Pointer(res))
    C.mvwaddstr(window.cwin, C.int(y), C.int(x), res)
}

func (window *Window) Addch(ch int) {
    C.waddch(window.cwin, C.chtype(ch))
}

func (window *Window) Mvaddch(y, x int, ch rune) {
    C.mvwaddch(window.cwin, C.int(y), C.int(x), C.chtype(ch))
}

// Hardware insert/delete feature.
func (window *Window) Idlok(bf bool) {
    C.idlok(window.cwin, C.bool(bf))
}

// Enable window scrolling.
func (window *Window) Scrollok(bf bool) {
    C.scrollok(window.cwin, C.bool(bf))
}

// Scroll given window.
func (window *Window) Scroll() {
    C.scroll(window.cwin)
}

// Get terminal size.
func Getmaxyx() (row, col int) {
    row = int(C.LINES)
    col = int(C.COLS)
    return row, col
}

// Erases content from cursor to end of line inclusive.
func (window *Window) Clrtoeol() {
    C.wrapper_wclrtoeol(window.cwin)
}

// Clears the console.
func Clear() int {
  return int(C.clear())
}
