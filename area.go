/*
 * File: area.go
 *
 * Description: Class to handle game Areas.
 */

package main

import "fmt"
import "math/rand"

// Structure to hold the points of the form (x,y)
type Coords struct {
    y int
    x int
}

// Structure to hold the initial tile and its properties
type Tile struct {
    Ch         rune
    BlockMove  bool
    BlockSight bool
}

// Structure to hold the entire location of a given area.
type Area struct {

    // Holds the current rune appearance.
    Tiles []Tile

    // Stores all of the creatures present in this level.
    Creatures  []*Creature

    // Contains all of the items in a level not held by creatures.
    Items []*Item

    // The length and width of the level.
    Height int
    Width  int

    // Whether or not the current area has been populated already.
    IsPopulatedWithCreatures bool
}

//! Generates an area and assigns a start location to the PC
/*
 * @param    int               height
 * @param    int               width
 *
 * @returns  Area* and (x,y)   A generated area array and (x,y) starting
 *                             points. 
 */
func NewArea(h, w int) (*Area, int, int) {

    // Input validation.
    if (h < 1 || w < 1) {
        DebugLog(&G, fmt.Sprintf("NewArea() --> invalid input"))
        return nil, 0, 0
    }

    // Variable declaration
    var ry, rx int
    var creatures []*Creature
    var items []*Item

    // Real x,y
    ry, rx = 0, 0

    // Assign memory for the creature and item arrays.
    creatures = make([]*Creature,0)
    items     = make([]*Item,0)

    // Number of iterations to use
    nIts := 4

    // Setup an iterator
    t := make([][]Tile, nIts)

    // Set the area padding
    SetPad(h, w)

    // Cycle thru all the elements of an iterator
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

//! Grab info about a given tile, specific what it is, whether it blocks,
//! and which creatures are present here.
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

    // Input validation, make sure this actually got an address to an
    // Area object.
    if a == nil {
        return ' ', false, nil, nil
    }

    // Variable declaration.
    var c *Creature = nil
    var items []*Item

    // Assign a chunk of memory for the array of Item objects.
    items = make([]*Item,0)

    // Read from the tiles array at the requested point to get the
    // specific Unicode rune value.
    ch = a.Tiles[int(x)+int(y)*a.Width].Ch

    // Read from the tiles array at a given point to determine whether
    // or not a tile is blocking.
    blocks = a.Tiles[int(x)+int(y)*a.Width].BlockMove

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

//! Randomly return a tile (e.g. # == wall and . == ground)
/*
 * @return    Tile    newly initialized tile object
 */
func placeRandomTile() Tile {

    // Make about 30% of the tiles walls (i.e. --> #)
    if rand.Intn(100) <= 30 {
        return Tile{'#', true, true}
    }

    // Otherwise return the ground tile.
    return Tile{'.', false, false}
}

//! Returns a random set of coordinates.
/*
 * @param      int       height
 * @param      int       width
 *
 * @returns    points    an (x,y) coord
 */
func selectRandomTile(h, w int) (int, int) {

    // Input validation.
    if (h < 1 || w < 1) {
        DebugLog(&G, fmt.Sprintf("selectRandomTile() --> invalid input"))
        return 0, 0
    }

    // Randomly generate a y-value.
    y := rand.Intn(h)

    // Randomly generate a x-value.
    x := rand.Intn(w)

    // Return the (x,y) coord.
    return y, x
}

//! With the tile given as argument make some more of those randomly around it.
/*
 * @param     int        y-coord
 * @param     int        x-coord
 * @param     int        w-coord
 * @param     *Tile[]    pointer to an array of tiles
 *
 * @return    none
 */
func explodeTile(y, x, w int, t *[]Tile) {

    // Input validation.
    if (t == nil) {
        DebugLog(&G, fmt.Sprintf("explodeTile() --> invalid input"))
        return
    }

    // Grab the tile currently in that location.
    originalTile := (*t)[x+y*w]

    // Cycle 5 times...
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
                ry = ry*-1
                rx = rx*-1
            }
        }()

        // If outside of the map borders, revert back to the original tile.
        if !mapBorders(y+ry, x+rx) || !mapBorders(y, x) {
            (*t)[(x+rx)+(y+ry)*w] = originalTile
        }

        // Adjust the coords accordingly.
        y = y+ry
        x = x+rx
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
    if (t == nil) {
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

    // Input validation, if no pointer to the tile array is given, do nothing.
    if (t == nil) {
        DebugLog(&G, fmt.Sprintf("floodFill() --> invalid input"))
        return
    }

    // Assign memory of size WorldHeight by WorldWidth.
    c := make([]Coords, WorldHeight*WorldWidth)

    // First element is the given coords.
    c[0] = Coords{y, x}

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

            // Then append it to the array.
            c = append(c[:i], c[i+1:]...)
        }
    }

    // Having completed filling the given tile array, go back.
    return
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

    // Input validation, was this given an array of coords and tiles?
    if (c == nil || t == nil) {
        DebugLog(&G, fmt.Sprintf("appendCoords() --> invalid input"))
        return
    }

    // Define the width
    w := WorldWidth

    // Array for each of the 3x3 chunks.
    three_by_three_chunks_y := []int{1, -1, -1,  1,  1, -1,  0,  0}
    three_by_three_chunks_x := []int{1, -1,  1, -1,  0,  0,  1, -1}

    // Get the derived (x,y) value
    for i := 0; i < 9; i++ {

        // Add the relevant chunk the derived x/y values.
        dy := y + three_by_three_chunks_y[i]
        dx := x + three_by_three_chunks_x[i]

        // Sanity check, make sure the values are still within the bounds.
        if !withinBounds(dy, dx) {
            continue;
        }

        // Further check, skip the non-ground tiles.
        if (*t)[dx+dy*w].Ch != '.' {
            continue;
        }

        // Append the values to the adjusted Coord array.
        *c = append(*c, Coords{dx, dy})
    }

    // Having appended the coords to the given array, go back.
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

    // Input validation, make sure this was given a tile array.
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

    // Otherwise return false.
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

    // Input validation, make sure this was given a tile array.
    if t == nil {
        DebugLog(&G, fmt.Sprintf("adjacentWalls() --> invalid input"))
        return 0
    }

    // Variable declaration
    counter := 0

    // Nearby wall layer? Add 2 then.
    if t[x+y*w].BlockMove {
        counter += 2
    }

    // Adjacent blocking tile, increment 1.
    if t[(x+1)+(y+1)*w].BlockMove {
        counter++
    }

    // Adjacent blocking tile, increment 1.
    if t[(x-1)+(y-1)*w].BlockMove {
        counter++
    }

    // Adjacent blocking tile, increment 1.
    if t[(x+1)+(y-1)*w].BlockMove {
        counter++
    }

    // Adjacent blocking tile, increment 1.
    if t[x+(y+1)*w].BlockMove {
        counter++
    }

    // Adjacent blocking tile, increment 1.
    if t[x+(y-1)*w].BlockMove {
        counter++
    }

    // Adjacent blocking tile, increment 1.
    if t[(x+1)+y*w].BlockMove {
        counter++
    }

    // Adjacent blocking tile, increment 1.
    if t[(x-1)+y*w].BlockMove {
        counter++
    }

    // Finally return the entire counter.
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

    // Sanity check, make sure this isn't a null pointer.
    if a == nil {
        DebugLog(&G, fmt.Sprintf("populateAreaWithCreatures() --> " +
                                 "invalid input"))
        return false
    }

    // Further check, make sure this area has not been previously populated.
    if a.IsPopulatedWithCreatures {
        DebugLog(&G, fmt.Sprintf("populateAreaWithCreatures() --> " +
                                 "area already previously populated..."))
        return false
    }

    // Define arrays to hold all the (x,y) points being used.
    var ys = make([]int,0)
    var xs = make([]int,0)

    // Variable to keep track of the current number.
    var coord_num uint = 0

    // Counters.
    var i uint = 0
    var j uint = 0

    // Flag to check if a coord (x,y) pair has already been used.
    var CoordIsAlreadyUtilized bool = false

    // Determine the number monsters to add to the area.
    //
    // Currently the formula is:
    //
    // Divide height and width each by 10, then multiply them both.
    //
    MaxNumberOfMonsters := uint(a.Height / 10) * uint(a.Width / 10)

    // Safety check, if less than zero, cap it at zero.
    if MaxNumberOfMonsters < 0 {
        MaxNumberOfMonsters = 0
    }

    // Continue to add monsters until the max has been reached.
    for i = 0; i < MaxNumberOfMonsters; i++ {

        // Grab a random x coord value.
        x_coord := getRandomNumBetweenZeroAndMax(a.Width)

        // Grab a random y coord value.
        y_coord := getRandomNumBetweenZeroAndMax(a.Height)

        // Cast them both to the derived coord pair.
        dy := y_coord
        dx := x_coord

        // Check if it already exists in the array holding the already
        // utilized points.
        for j = 0; j < coord_num; j++ {

            // Safety check, if the array is empty, skip this part.
            if len(xs) == 0 || len(ys) == 0 {
                break
            }

            // TODO: delete this once the loop works
            break

            // Set the flag to true if the (x,y) coord pair is already used.
            if xs[j] == dx && ys[j] == dy {
                CoordIsAlreadyUtilized = true
                break
            }
        }

        // If the (x,y) pair has already been used...
        if CoordIsAlreadyUtilized {

            // Decrement the value of i to a min of zero.
            if (i > 0) {
                i--
            }

            // Move on to the next instance to avoid spawning spawning
            // multiple monsters at the (x,y) pair.
            continue
        }

        // Since the (x,y) pair has yet to be used, attempt to grab details
        // about the tile at this (x,y) location.
        tile_rune, blocking, _, _ := a.GetTileInfo(dy, dx)

        // Safety check, make sure the tile isn't a wall or blocking tile.
        if tile_rune == '#' || blocking {

            // If it is a wall tile, decrement the value of i
            if (i > 0) {
                i--
            }

            // Move on to the next instance.
            continue
        }

        // Since this is a new point, go ahead and determine the number of
        // blocking tiles.
        nearby_wall_count := adjacentWalls(y_coord, x_coord, a.Width, a.Tiles)

        // If the number of blocking tiles is greater than 1...
        if nearby_wall_count > 1 {

            // If it is greater than 1, decrement the value of i
            if (i > 0) {
                i--
            }

            // Move on to the next instance
            continue
        }

        // As there are 1 or fewer blocking tiles nearby, it ought to be safe
        // to spawn a monster since this is a wide open area (which is to say,
        // very few walls or blockers).
        spawnCreatureToArray("dog", dx, dy, a)

        // Append the points to the relevant arrays.
        xs = append(xs, dx)
        ys = append(ys, dy)

        // Increment the coord_num counter.
        coord_num++
    }

    // Set the "IsPopulatedWithCreatures" flag to true since it now
    // contains various creatures / critters / monsters.
    a.IsPopulatedWithCreatures = true

    // Tell the developer this population function was successful.
    DebugLog(&G, fmt.Sprintf("populateAreaWithCreatures() --> " +
                             "creatures populated into area successfully"))

    // Since the monsters were added successfully, go ahead and return true.
    return true
}
