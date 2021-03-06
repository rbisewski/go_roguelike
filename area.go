/*
 * File: area.go
 *
 * Description: Class to handle game Areas.
 */

package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

// Coords ... Structure to hold the points of the form (x,y)
type Coords struct {

	// Quick comparison value in the form "x:y"
	id string

	// Coord (x,y) pair
	y int
	x int
}

// Tile ... Structure to hold the initial tile and its properties
type Tile struct {
	Ch         rune
	BlockMove  bool
	BlockSight bool
}

// Area ... Structure to hold the entire location of a given area.
type Area struct {

	// Holds the current rune appearance.
	Tiles []Tile

	// Stores all of the creatures present in this level.
	Creatures []*Creature

	// Contains all of the items in a level not held by creatures.
	Items []*Item

	// The length and width of the level.
	Height int
	Width  int

	// Whether or not the current area has been populated already.
	IsPopulatedWithCreatures bool
}

// NewArea ... Generates an area and assigns a start location to the PC
/*
 * @param    int               height
 * @param    int               width
 *
 * @returns  Area* and (x,y)   A generated area array and (x,y) starting
 *                             points.
 */
func NewArea(h, w int) (*Area, int, int) {

	if h < 1 || w < 1 {
		DebugLog(&G, fmt.Sprintf("NewArea() --> invalid input"))
		return nil, 0, 0
	}

	var ry, rx int
	var creatures []*Creature
	var items []*Item

	// Real x,y
	ry, rx = 0, 0

	creatures = make([]*Creature, 0)
	items = make([]*Item, 0)

	// Number of iterations to use
	nIts := 4

	// Setup an iterator
	t := make([][]Tile, nIts)

	// Set the area padding
	SetPad(h, w)

	for it := 0; it < nIts; it++ {

		// Assign a tile to that given location
		t[it] = make([]Tile, w*h)

		// For all of the y-coords (height)...
		for y := 0; y < h; y++ {

			// For all of the x-coords (width)...
			for x := 0; x < w; x++ {

				// On first iteration, place random tiles.
				if it == 0 {
					t[it][x+y*w] = placeRandomTile()
					continue
				}

				// Otherwise check for wall placement.
				if mapBorders(y, x) || adjacentWalls(y, x, w, t[it-1]) >= 4 {
					t[it][x+y*w] = Tile{'#', true, true}
					continue
				}

				// Or draw ground.
				t[it][x+y*w] = Tile{'.', false, false}

				// If we are at last Iteration
				if it == nIts-1 {

					//set the spawn coords
					ry = y
					rx = x
				}
			}
		}
	}

	// Return the completed area-object plus start coords.
	return &Area{t[nIts-1], creatures, items, h, w, false}, ry, rx
}

// GetTileInfo ... Grab info about a given tile, specific what it is,
// whether it blocks, and which creatures are present here.
/*
 * @param         int    y-value
 * @param         int    x-value
 *
 * @returns      rune    unicode rune value
 *               bool    whether or not the tile is "blocking"
 *           Creature*   pointer to a monster, if any is present
 */
func (a *Area) GetTileInfo(y, x int) (ch rune,
	blocks bool,
	hasCreature *Creature,
	hasItems []*Item) {

	if a == nil {
		return ' ', false, nil, nil
	}

	var c *Creature
	var items []*Item
	var RequestedTile = x + y*a.Width

	// Safety check, make sure the requested tile is a sane value.
	if RequestedTile < 0 || RequestedTile > len(a.Tiles) {
		return ' ', false, nil, nil
	}

	items = make([]*Item, 0)

	// Read from the tiles array at the requested point to get the
	// specific Unicode rune value.
	ch = a.Tiles[RequestedTile].Ch

	// Read from the tiles array at a given point to determine whether
	// or not a tile is blocking.
	blocks = a.Tiles[RequestedTile].BlockMove

	// Cycle thru every monster in the array...
	for _, m := range a.Creatures {

		// If the Monster is alive and at (x,y) point.
		if m.Hp > 0 && m.X == x && m.Y == y {

			// Grab a reference to that creature and break.
			c = m
			break
		}
	}

	// Cycle thru every item in the array...
	for _, itm := range a.Items {

		// If an item is located at (x,y) point...
		if itm.X == x && itm.Y == y {

			// Append it to the list of items.
			items = append(items, itm)
		}
	}

	// Return the rune, whether this is a blocking tile, and reference to
	// the creature (if any) and a reference to items laying on the ground
	// at the particular spot (if any).
	return ch, blocks, c, items
}

// placeRandomTile ... randomly return a tile (e.g. # == wall and . == ground)
/*
 * @return    Tile    newly initialized tile object
 */
func placeRandomTile() Tile {

	// Make about 30% of the tiles walls (i.e. --> #)
	if rand.Intn(100) <= 30 {
		return Tile{'#', true, true}
	}

	return Tile{'.', false, false}
}

// selectRandomTile ... Returns a random set of coordinates.
/*
 * @param      int       height
 * @param      int       width
 *
 * @returns    points    an (x,y) coord
 */
func selectRandomTile(h, w int) (int, int) {

	if h < 1 || w < 1 {
		DebugLog(&G, fmt.Sprintf("selectRandomTile() --> invalid input"))
		return 0, 0
	}

	// Randomly generate a y-value and an x-value
	y := rand.Intn(h)
	x := rand.Intn(w)

	return y, x
}

// explodeTile ... With the tile given as argument make some new tiles
// randomly around it.
/*
 * @param     int        y-coord
 * @param     int        x-coord
 * @param     int        w-coord
 * @param     *Tile[]    pointer to an array of tiles
 *
 * @return    none
 */
func explodeTile(y, x, w int, t *[]Tile) {

	if t == nil {
		DebugLog(&G, fmt.Sprintf("explodeTile() --> invalid input"))
		return
	}

	// Grab the tile currently in that location.
	originalTile := (*t)[x+y*w]

	for it := 0; it < 5; it++ {

		// Randomly generate some small integers.
		ry := rand.Intn(2)
		rx := rand.Intn(2)

		// If heads then go back 1 for y-coord.
		if TossCoin() {
			ry *= -1
		}

		// If heads then go back 1 for x-coord.
		if TossCoin() {
			rx *= -1
		}

		// If anything happens (e.g. a panic) then go ahead and just
		// decrement the real x,y values.
		defer func() {
			if r := recover(); r != nil {
				ry = ry * -1
				rx = rx * -1
			}
		}()

		// If outside of the map borders, revert back to the original tile.
		if !mapBorders(y+ry, x+rx) || !mapBorders(y, x) {
			(*t)[(x+rx)+(y+ry)*w] = originalTile
		}

		// Adjust the coords accordingly.
		y = y + ry
		x = x + rx
	}
}

//! Searches for the first walkable tile on the map.
/*
 * @param     *Tile[]    pointer to an array of Tile objects
 *
 * @return    int        y-coord
 *            int        x-coord
 */
func firstGroundTile(t *[]Tile) (int, int) {

	// If no pointer to a tile array is given, default to (0,0)
	if t == nil {
		return 0, 0
	}

	// For every unit of height...
	for y := 0; y < WorldHeight; y++ {

		// For every unit of width...
		for x := 0; x < WorldWidth; x++ {

			// Not "blocking"? Then return that...
			if !(*t)[x+y*WorldWidth].BlockMove {
				return y, x
			}
		}
	}

	// Default is (0,0)
	return 0, 0
}

//! En-masse fill a section of the area map with tiles.
/*
 * @param      int       y-value
 * @param      int       x-value
 * @param      *Tile[]   pointer to an array of tiles.
 *
 * @returns    none
 */
func floodFill(y, x int, t *[]Tile) {

	if t == nil {
		DebugLog(&G, fmt.Sprintf("floodFill() --> invalid input"))
		return
	}

	c := make([]Coords, WorldHeight*WorldWidth)

	// First element is the given coords.
	c[0] = Coords{strconv.Itoa(x) + ":" + strconv.Itoa(y),
		y,
		x}

	// If all else fails... end here.
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	// Iterate 9 times (i.e. a 3x3 chuck)
	for it := 0; it < 9; it++ {

		// Cycle thru the world map chunks...
		for i, coord := range c {

			// Sanity check, make sure the X-coord is reasonable (i.e. on
			// the "visible" worldmap).
			if coord.x < 0 || coord.x >= WorldWidth {
				continue
			}

			// Sanity check, make sure the Y-coord is reasonable (i.e. on
			// the "visible" worldmap).
			if coord.y < 0 || coord.y >= WorldHeight {
				continue
			}

			// If tile is "blocking" then skip.
			if (*t)[coord.x+coord.y*WorldWidth].BlockMove {
				continue
			}

			// Since everything seems fine, go ahead and generate a tile.
			(*t)[coord.x+coord.y*WorldWidth] = Tile{'1', false, false}

			// Attach the coords.
			appendCoords(coord.y, coord.x, &c, t)
			c = append(c[:i], c[i+1:]...)
		}
	}
}

//! Function to append (x,y) to a set of coords / tiles
/*
 * @param     int         y-value
 * @param     int         x-value
 * @param     *Coords[]   array of coords
 * @param     *Tile[]     array of tiles
 *
 * @return    none
 */
func appendCoords(y, x int, c *[]Coords, t *[]Tile) {

	if c == nil || t == nil {
		DebugLog(&G, fmt.Sprintf("appendCoords() --> invalid input"))
		return
	}

	w := WorldWidth

	// Array for each of the 3x3 chunks.
	threeByThreeChunksY := []int{1, -1, -1, 1, 1, -1, 0, 0}
	threeByThreeChunksX := []int{1, -1, 1, -1, 0, 0, 1, -1}

	for i := 0; i < 9; i++ {

		// Add the relevant chunk the derived x/y values.
		dy := y + threeByThreeChunksY[i]
		dx := x + threeByThreeChunksX[i]

		// Sanity check, make sure the values are still within the bounds.
		if !withinBounds(dy, dx) {
			continue
		}

		// Further check, skip the non-ground tiles.
		if (*t)[dx+dy*w].Ch != '.' {
			continue
		}

		// Append the values to the adjusted Coord array.
		*c = append(*c, Coords{strconv.Itoa(dx) + ":" + strconv.Itoa(dy),
			dx,
			dy})
	}

	return
}

//! Determine whether or not the coords are not out of range.
/*
 * @param     int     y-value
 * @param     int     x-value
 *
 * @return    bool    whether or not a given point is in bounds.
 */
func withinBounds(y, x int) bool {
	return y > 0 && y < WorldHeight && x > 0 && x < WorldWidth
}

//! Returns true if any of the adjacent tiles block move
/*
 * @param     int       y-value
 * @param     int       x-value
 * @param     int       width
 * @param     Tile[]    array of tiles
 *
 * @return    bool      whether or not a given tile has nearby walls
 */
func anyAdjacentWalls(y, x, w int, t []Tile) bool {

	if t == nil {
		DebugLog(&G, fmt.Sprintf("anyAdjacentWalls() --> invalid input"))
		return false
	}

	// If "blocking" then this is, for all purposes, a "wall" tile, so
	// go ahead and return true here.
	if t[x+y*w].BlockMove ||
		t[(x+1)+(y+1)*w].BlockMove ||
		t[(x-1)+(y-1)*w].BlockMove ||
		t[(x+1)+(y-1)*w].BlockMove ||
		t[(x-1)+(y+1)*w].BlockMove ||
		t[x+(y+1)*w].BlockMove ||
		t[x+(y-1)*w].BlockMove ||
		t[(x+1)+y*w].BlockMove ||
		t[(x-1)+y*w].BlockMove {
		return true
	}

	return false
}

//! Function to help take into account nearby tiles for level generation.
/*
 * @param    int       y-value
 * @param    int       x-value
 * @param    int       width
 * @param    Tile[]    array of tiles
 *
 * @return   int       count of nearby blocking tiles
 */
func adjacentWalls(y, x, w int, t []Tile) int {

	if t == nil {
		DebugLog(&G, fmt.Sprintf("adjacentWalls() --> invalid input"))
		return 0
	}

	counter := 0

	// Nearby wall layer? Add 2 then.
	if t[x+y*w].BlockMove {
		counter += 2
	}

	// Adjacent blocking tile in a given direction? Then increment 1...
	if t[(x+1)+(y+1)*w].BlockMove {
		counter++
	}
	if t[(x-1)+(y-1)*w].BlockMove {
		counter++
	}
	if t[(x+1)+(y-1)*w].BlockMove {
		counter++
	}
	if t[x+(y+1)*w].BlockMove {
		counter++
	}
	if t[x+(y-1)*w].BlockMove {
		counter++
	}
	if t[(x+1)+y*w].BlockMove {
		counter++
	}
	if t[(x-1)+y*w].BlockMove {
		counter++
	}

	return counter
}

//! Function to check if a given (x,y) is in the bounds of the world map.
/*
 * @param     int    y-coord
 * @param     int    x-coord
 *
 * @return    bool   whether in or out of bounds
 */
func mapBorders(y, x int) bool {
	return y == 0 || y == WorldHeight-1 || x == 0 || x == WorldWidth-1
}

//! Populate an area with creatures / critters / monsters; this only works
//! if the level has yet to be populated.
/*
 * @returns    bool    whether or not the tile is "blocking"
 */
func (a *Area) populateAreaWithCreatures() bool {

	if a == nil {
		DebugLog(&G, fmt.Sprintf("populateAreaWithCreatures() --> "+
			"invalid input"))
		return false
	}

	if a.IsPopulatedWithCreatures {
		DebugLog(&G, fmt.Sprintf("populateAreaWithCreatures() --> "+
			"area already previously populated..."))
		return false
	}

	var CoordsArray = make([]Coords, 0)
	var coordNum uint
	var i uint
	var CoordIsAlreadyUtilized = false

	// Determine the number monsters to add to the area.
	//
	// Currently the formula is:
	//
	// Divide height and width each by 10, then multiply them both.
	//
	MaxNumberOfMonsters := uint(a.Height/10) * uint(a.Width/10)

	// Safety check, if less than zero, cap it at zero.
	if MaxNumberOfMonsters < 0 {
		MaxNumberOfMonsters = 0
	}

	// Continue to add monsters until the max has been reached.
	for i = 0; i < MaxNumberOfMonsters; i++ {

		// Set the already utilized flag back to false.
		CoordIsAlreadyUtilized = false

		// Grab a random x coord value.
		dx := getRandomNumBetweenZeroAndMax(a.Width)

		// Grab a random y coord value.
		dy := getRandomNumBetweenZeroAndMax(a.Height)

		// Assemble a Coords object from the above info.
		CurrentCoordPair := Coords{strconv.Itoa(dx) + ":" + strconv.Itoa(dy),
			dy,
			dx}

		// Check if it already exists in the array holding the already
		// utilized points.
		for _, coord := range CoordsArray {

			// Safety check, if the id is empty, then skip this part.
			if len(coord.id) == 0 {
				break
			}

			// Set the flag to true if the (x,y) coord pair is already used.
			if CurrentCoordPair.id == coord.id {
				CoordIsAlreadyUtilized = true
				break
			}
		}

		// If the (x,y) pair has already been used...
		if CoordIsAlreadyUtilized {

			// Decrement the value of i to a min of zero.
			if i > 0 {
				i--
			}

			// Move on to the next instance to avoid spawning multiple
			// monsters at the (x,y) pair.
			continue
		}

		// Since the (x,y) pair has yet to be used, attempt to grab details
		// about the tile at this (x,y) location.
		tileRune, blocking, _, _ := a.GetTileInfo(dy, dx)

		// Safety check, make sure the tile isn't a wall or blocking tile.
		if tileRune == '#' || blocking {

			// If it is a wall tile, decrement the value of i
			if i > 0 {
				i--
			}

			// Move on to the next instance.
			continue
		}

		// Since this is a new point, go ahead and determine the number of
		// blocking tiles.
		nearbyWallCount := adjacentWalls(dy, dx, a.Width, a.Tiles)

		// If the number of blocking tiles is greater than 1...
		if nearbyWallCount > 1 {

			// If it is greater than 1, decrement the value of i
			if i > 0 {
				i--
			}

			// Move on to the next instance
			continue
		}

		//
		// As there are 1 or fewer blocking tiles nearby, it ought to be safe
		// to spawn a monster since this is a wide open area (which is to say,
		// very few walls or blockers).
		//
		// The code in this below block will randomly spawn creatures into
		// a given area. Perhaps in the future consider making it area type
		// dependant; e.g. spawn goblins and bats in cave areas
		//
		// Firstly, do a quick safety check to ensure the creature types are
		// actually populated correctly into the global array.
		//
		if GlobalCreatureTypeInfoMapIsPopulated {

			// Determine the current number of types
			numOfTypes := len(GlobalCreatureTypeInfoMap)

			// Attempt to grab a number between 0 and numOfTypes
			chosenTypeNum := getRandomNumBetweenZeroAndMax(numOfTypes - 1)

			// Grab a creature type stored at the address specified by
			// chosenTypeNum.
			var chosenCreatureType string
			var i int
			for k := range GlobalCreatureTypeInfoMap {

				// If the counter matches the requested element, break
				if i == chosenTypeNum {
					chosenCreatureType = k
					break
				}

				// Else increment the counter
				i++
			}

			// Safety check, ensure that the chosenCreatureType isn't empty.
			if chosenCreatureType == "" {
				break
			}

			// Attempt to spawn a creature of that type
			wasSuccessful := spawnCreatureToArray(chosenCreatureType, dx, dy, a)
			if !wasSuccessful {
				DebugLog(&G, fmt.Sprintf("populateAreaWithCreatures() --> "+
					"Unable to spawn chosen creature into the area!"))
				break
			}
		}

		CoordsArray = append(CoordsArray, CurrentCoordPair)
		coordNum++
	}

	// Set the "IsPopulatedWithCreatures" flag to true since it now
	// contains various creatures / critters / monsters.
	a.IsPopulatedWithCreatures = true

	DebugLog(&G, fmt.Sprintf("populateAreaWithCreatures() --> "+
		"creatures populated into area successfully"))

	return true
}
