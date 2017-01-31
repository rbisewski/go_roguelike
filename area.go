/*
 * File: area.go
 *
 * Description: Class to handle game Areas.
 */

package main

import "math/rand"
import "fmt"

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
    Tiles []Tile
    Creatures  []*Creature
    Items []*Creature

    Height int
    Width  int
}

//! Generates an area and assigns a start location to the PC
/*
 * @param    int               height
 * @param    int               width
 *
 * @returns  Area* and (x,y)   A generated area array and (x,y) starting
 *                             points. 
 */
func NewArea(h, w int) (*Area, Coord, Coord) {

    // Variable declaration
    var ry, rx Coord

    // Number of iterations to use
    nIts := 4

    // Real x,y
    ry, rx = 0, 0

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
                        ry = Coord(y)
                        rx = Coord(x)
                    }
                }
            }

            // Iteration complete
            Write(51, 2, fmt.Sprint("Done"))
        }

    // Return the completed area-object plus start coords.
    return &Area{t[nIts-1], nil, nil, h, w}, ry, rx
}

//! Grab info about a given tile, specific what it is, whether it blocks,
//! and which creatures are present here.
/*
 * @param    Coord    y-value
 * @param    Coord    x-value
 *
 * @returns      rune     unicode rune value
 *               bool     whether or not the tile is "blocking"
 *           Creature*    pointer to a monster, if any is present
 */
func (a *Area) GetTileInfo(y, x Coord) (ch rune,
                                         blocks bool,
                                         hasCreature *Creature) {

    // Input validation, make sure this actually got an address to an
    // Area object.
    if a == nil {
        return
    }

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

            // State that a monster has been found.
            hasCreature = m

            // All done here.
            return
        }
    }

    // Go back.
    return
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

    // Randomly generate a y-value.
    y := rand.Intn(h)

    // Randomly generate a x-value.
    x := rand.Intn(w)

    // Return the (x,y) coord.
    return y, x
}

//! With the tile given as argument make some more of those randomly around it.
/*
 * @param    int        y-coord
 * @param    int        x-coord
 * @param    int        w-coord
 * @param    *Tile[]    pointer to an array of tiles
 */
func explodeTile(y, x, w int, t *[]Tile) {

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
                ry, rx = ry*-1, rx*-1
            }
        }()

        // If outside of the map borders, revert back to the original tile.
        if !mapBorders(y+ry, x+rx) || !mapBorders(y, x) {
            (*t)[(x+rx)+(y+ry)*w] = originalTile
        }

        // Adjust the coords accordingly.
        y, x = y+ry, x+rx
    }
}

// Searches for the first walkable tile on the map.
func firstGroundTile(t *[]Tile) (int, int) {

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

// Returns true if any of the adjacent tiles block move
func anyAdjacentWalls(y, x, w int, t []Tile) bool {

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
 * @return 
 */
func adjacentWalls(y, x, w int, t []Tile) int {

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
