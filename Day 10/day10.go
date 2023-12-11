package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)
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
var pipemap [][]rune
var p2PathPositions []Position
var start Position
func main() {
    file, err := os.Open("input.txt")
    if err != nil { log.Fatal(err) }
    scanner := bufio.NewScanner(file)
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
    p2PathPositions = append(p2PathPositions, Position {
        y: start.y,
        x: start.x,
    })
    var paths []Position
    for _, direction := range []rune { 'L', 'U', 'R', 'D' } {
        startPath := Position { x: start.x, y: start.y, steps: 0 }
        Move(&startPath, direction)
        if startPath.x != start.x || startPath.y != start.y {
            paths = append(paths, startPath)
            p2PathPositions = append(p2PathPositions, Position {
                y: startPath.y,
                x: startPath.x,
            })
        }
    }
    for idx, path := range paths {
        fmt.Printf("Path %v starting at %v (%v,%v)\n",
            idx, string(pipemap[path.y][path.x]), path.y, path.x)
    }
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
    containedCount := 0
    for _,path := range p2PathPositions {
        pipemap[path.y][path.x] = '@'
    }
    for row := range pipemap {
        if pipemap[row][0] == '@' { continue }
        pipemap[row][0] = '#'
    }
    for row := range pipemap {
        if pipemap[row][len(pipemap[row])-1] == '@' { continue }
        pipemap[row][len(pipemap[row])-1] = '#'
    }
    for col := range pipemap[0] {
        if pipemap[0][col] == '@' { continue }
        pipemap[0][col] = '#'
    }
    for col := range pipemap[len(pipemap) - 1] {
        if pipemap[len(pipemap)-1][col] == '@' { continue }
        pipemap[len(pipemap)-1][col] = '#'
    }
    changedCnt := 1
    for changedCnt != 0 {
        changedCnt = 0
        for row := 1; row < len(pipemap) -1; row++ {
            for col := 1; col < len(pipemap[row]) -1; col++ {
                if pipemap[row][col] != '@' && 
                    pipemap[row][col] != '#' && 
                    (pipemap[row][col -1] == '#' ||
                    pipemap[row][col + 1] == '#' ||
                    pipemap[row - 1][col] == '#' ||
                    pipemap[row + 1][col] == '#') {
                    changedCnt += 1
                    pipemap[row][col] = '#'
                }
            }
        }
    }
    for row := range pipemap {
        for col := range pipemap[row] {
            if pipemap[row][col] != '@' &&
                pipemap[row][col] != '#' {
                containedCount++
            }
        }
    }
    for rowIdx := range pipemap {
        for colIdx := range pipemap[rowIdx] {
            fmt.Print(string(pipemap[rowIdx][colIdx]))
        }
        fmt.Println()
    }
    fmt.Printf("Part 1: Total Steps: %v\n", paths[0].steps)
    fmt.Printf("Part 2: Area: %v\n", containedCount)
}
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

func AllSamePosition(paths []Position) bool {
    posX := paths[0].x
    posY := paths[0].y
    for _,path := range paths[1:] {
        if path.x != posX || path.y != posY { return false }
    }
    return true
}