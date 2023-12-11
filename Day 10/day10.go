package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)
/* Type to hold a set of x and y coordinates and the exit direction for a node.
I also use one of these to track path positions for part 1, so I added a step
count as well */
type Position struct {
    x int
    y int
    steps int
    exitDirection rune
}
/* Map of output direction based on entry direction (L, U, R, D)
and type of pipe.  For example, if we enter an 'F' pipe from the right,
we will go down. */
var outdirection  = map[rune]map[rune]rune {
    'L' :  {
        '-': 'L',
        'F': 'D',
        'L': 'U',
    },
    'D' : {
        '|': 'D',
        'J': 'L',
        'L': 'R',
    },
    'R' : {
        'J': 'U',
        '7': 'D',
        '-': 'R',
    },
    'U' : {
        '|': 'U',
        '7': 'L',
        'F': 'R', 
    },
}
// Store the entire map
var pipemap [][]rune
// Added this for part 2.  Store each position on the path individually.
var p2PathPositions []Position
// Store the start position
var start Position
func main() {
    file, err := os.Open("input.txt")
    if err != nil { log.Fatal(err) }
    scanner := bufio.NewScanner(file)
    /* Parser: each line of input contains a series of characters representing
    various map symbols.  While we're storing these, we'll also look for 'S'
    which represents the starting pipe. */
    for scanner.Scan() {
        line := scanner.Text()
        row := make([]rune, len(line))
        for i, c := range line {
            row[i] = c            
            if c == 'S' {
                start = Position {
                    y: len(pipemap),
                    x: i,
                }
            }
        }
        pipemap = append(pipemap, row)
    }
    fmt.Printf("Starting from row %v col %v\n", start.y, start.x)
    // Added for part 2: make the start position part of the path list
    p2PathPositions = append(p2PathPositions, Position {
        y: start.y,
        x: start.x,
    })
    // Initially this was planned to support any number of paths coming
    // from start even though the directions state there will be exactly
    // two. Still, this will contain the start position for each path we're
    // going to take.
    var paths []Position
    // This is a somewhat hacky solution added for part 2.  We'll need to 
    // know which directions we take from start so that we can replace 'S'
    // with the appropriate pipe.
    var startDirections []rune
    // Check each direction from start to see if we are able to make a valid
    // move to it.  If so, we'll add it to our paths list.
    for _, direction := range []rune { 'L', 'U', 'R', 'D' } {
        startPath := Position { x: start.x, y: start.y, steps: 0 }
        Move(&startPath, direction)
        if startPath.x == start.x && startPath.y == start.y {
            continue
        }
        startDirections = append(startDirections, direction)
        paths = append(paths, startPath)
        p2PathPositions = append(p2PathPositions, Position {
            y: startPath.y,
            x: startPath.x,
        })
    }
    // Added for part 2.  Replace 'S' in the map with the appropriate pipe
    if len(startDirections) < 2 {
        log.Fatal("Unexpected number of paths")
    }
    if startDirections[0] == 'L' && startDirections[1] == 'U' {
        pipemap[start.y][start.x] = 'J'
    }
    if startDirections[0] == 'L' && startDirections[1] == 'R' {
        pipemap[start.y][start.x] = '-'
    }
    if startDirections[0] == 'L' && startDirections[1] == 'D' {
        pipemap[start.y][start.x] = '7'
    }
    if startDirections[0] == 'U' && startDirections[1] == 'R' {
        pipemap[start.y][start.x] = 'L'
    }
    if startDirections[0] == 'U' && startDirections[1] == 'D' {
        pipemap[start.y][start.x] = '|'
    }
    if startDirections[0] == 'R' && startDirections[1] == 'D' {
        pipemap[start.y][start.x] = 'F'
    }
    for idx, path := range paths {
        fmt.Printf("Path %v starting at %v (%v,%v)\n",
            idx, string(pipemap[path.y][path.x]), path.y, path.x)
    }
    /* Part 1: walk each path until they reach the same position (which will
    be the furthest position on the loop from the start
    */
    for !AllSamePosition(paths) {
        for idx := range paths {
            pstartX := paths[idx].x
            pstartY := paths[idx].y
            Move(&paths[idx], paths[idx].exitDirection)
            fmt.Printf("Path %v moved from %v (%v,%v) to %v (%v,%v)\n",
                idx, string(pipemap[pstartY][pstartX]), pstartY, pstartX,
                string(pipemap[paths[idx].y][paths[idx].x]), 
                paths[idx].y, paths[idx].x)
            p2PathPositions = append(p2PathPositions, Position {
                y: paths[idx].y,
                x: paths[idx].x,
            })
            if paths[idx].x == pstartX && paths[idx].y == pstartY {
                log.Fatal("Reached a dead end!")
            }
        }
    }
    /* Part 2: It took me a few tries to get this right, but I was settled on
    ray tracing as the solution from the start.  We need to find the total 
    number of points that are completely enclosed within the loop. Even if there
    is not a full space between two pipes, it can still constitute an opening.
    To account for this, I cast a ray in a single direction from each point on
    the map.  If I cross an odd number of vertical pipes (including '|', 'J',
    or 'L'), then that map position must be enclosed within the loop. */
    containedCount := 0
    for row := range pipemap {
        for col := range pipemap[row] {
            if P2PathContains(row, col) { continue }
            isEnclosed := false
            for x := col; x >= 0; x-- {
                if P2PathContains(row, x) && (
                    pipemap[row][x] == '|' ||
                    pipemap[row][x] == 'J' ||
                    pipemap[row][x] == 'L') {
                    isEnclosed = !isEnclosed
                }
            }
            if isEnclosed { containedCount++ }
        }
    }
    // Print out the map for diagnostics
    for rowIdx := range pipemap {
        for colIdx := range pipemap[rowIdx] {
            fmt.Print(string(pipemap[rowIdx][colIdx]))
        }
        fmt.Println()
    }
    fmt.Printf("Part 1: Total Steps: %v\n", paths[0].steps)
    fmt.Printf("Part 2: Area: %v\n", containedCount)
}
/* Function to move along the path in a given direction.  I use this
at the beginning to figure out which directions we move from start,
then I continue to call it for path moves in part 1 based on the "output
direction" from the previous pipe (which I also set within the Position 
object when a valid move is made)
If a move is invalid, then the position will be unaltered.  This also prints
a warning if a move would result in going out of bounds (which shouldn't happen)
*/
func Move(pos *Position, direction rune) {
    newY := pos.y
    newX := pos.x
    if direction == 'L' { newX-- }
    if direction == 'U' { newY-- }
    if direction == 'R' { newX++ }
    if direction == 'D' { newY++ }
    if newX < 0 || newY < 0 || 
        newX > len(pipemap[0]) || 
        newY > len(pipemap) {
        fmt.Printf("Direction %v to row %v, col %v out of bounds\n", 
            string(direction), newY, newX)
        return
    }
    outdir := outdirection[direction][pipemap[newY][newX]]
    if outdir != 0 {
        pos.x = newX
        pos.y = newY
        pos.exitDirection = outdir
        pos.steps++
    }
}
// Helper to see if all supplied paths are at the same position
func AllSamePosition(paths []Position) bool {
    posX := paths[0].x
    posY := paths[0].y
    for _,path := range paths[1:] {
        if path.x != posX || path.y != posY { return false }
    }
    return true
}
// For part two, we need to see if a given set of coordinates sits on the
// path that we recorded (p2PathPositions)
func P2PathContains(y, x int) bool {
    for _,path := range p2PathPositions {
        if path.x == x && path.y == y {
            return true
        }
    }
    return false
}
