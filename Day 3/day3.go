package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
    "strconv"
)

    
var lines []string
var nums [][]string
var symbols [][]string
var exprDigit, _ = regexp.Compile(`\d`)

func main() {
    file, err := os.Open("input.txt")
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(file)
    maxY := 0 
    maxX := 0 
    // For this one, I had to read the entire file into memory
    // or we would have been seeking the stream all over the place.
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) > maxX {
            maxX = len(line) + 1
        }
        lines = append(lines, line)
        maxY++
    }
    /*
    Parsing the input: I wasn't sure what was going to be coming
    for part 2, but I figured it would be best to store the positions
    of all items of interest (numbers and symbols)
    For numbers, we just need to track the starting x position
    since we are storing the number as a string and can get
    the ending index from its length.  This does mean we are technically
    allocating more memory than needed since the other decimal positions
    will be empty otherwise, but it does give us the capability to track
    a number starting at any index.
    Symbols are only 1 character wide, so that slice is just storing the
    y and x coordinates of each symbol.
    For both lists, the outside index is the row number (y position) and
    the inside index is the column number (x position) 
    */
    nums = make([][]string, maxY)
    symbols = make([][]string, maxY)
    for y := range nums {
        nums[y] = make([]string, maxX)
        symbols[y] = make([]string, maxX)
    }
    for y, line := range lines {
        numStart := -1
        numMatch := false
        numCurrent := ""
        for x, ch := range strings.Split(line, "") {
            if ch == "." {
                if !numMatch { continue }
                
                fmt.Printf(
                    "Found number %v at row %v starting at column %v\n",
                    numCurrent, y, numStart,
                )
                
                nums[y][numStart] = numCurrent 
                numCurrent = ""
                numMatch = false
                numStart = -1
                continue
            }
            if exprDigit.MatchString(ch) {
                numMatch = true
                numCurrent += ch
                if (numStart == -1) {
                    numStart = x
                }
                continue
            }
            fmt.Printf(
                "Found symbol %v at row %v column %v\n",
                ch, y, x,
            )
            symbols[y][x] = ch
            if !numMatch { continue }
            nums[y][numStart] = numCurrent 
            numCurrent = ""
            numMatch = false
            numStart = -1
        }
        if !numMatch { continue }
        nums[y][numStart] = numCurrent
    }
    /*
    Part 1: find the sum of part numbers.  A number is a part if it is
    adjacent to any symbol (including diagonals).  I first find the starting
    and ending rows to check based on the row prior to to the number's row
    (unless its in the first row) and the row after the number's row
    (unless its in the last row) as well as the starting and ending columns
    to check based on the column prior to the number's starting column
    (unless its in the first column) and the column after the number's ending
    column (unless it is in the last column).  I check each position in the
    parts array to see if there is a part in range.  As soon as I find a part,
    I increment the total accordingly and move on to the next number.  I could
    technically skip the positions in the parts slice where I know the number
    itself is, but it would not make it significantly faster.
    */
    total := 0
    for y, row := range nums {
        for numStart, num := range row  {
            if num == "" { continue }
            startcol := 0
            if numStart > 0 { startcol = numStart - 1 }
            endcol := numStart + len(num) 
            if endcol > maxX { endcol = maxX }
            startrow := 0
            if y > 0 { startrow = y - 1 }
            endrow := y + 1
            if endrow > maxY - 1 { endrow = maxY - 1 }
             
            fmt.Printf(
                "Number %v on row %v starts at col %v and ends at col %v\n",
                num, y, numStart, numStart + len(num))
            fmt.Printf(
                "We will check for symbols between %v,%v and %v,%v (inclusive)\n",
                startrow, startcol, endrow, endcol)
            
            symbol := ""
            for row := startrow; row <= endrow; row++ {
                for col := startcol; col <= endcol; col++ {
                    fmt.Printf("Symbol at row %v col %v: %v\n", row, col,
                        symbols[row][col])
                    if symbols[row][col] == "" { continue }
                    symbol = symbols[row][col]
                    break;
                }
                if symbol != "" { break }
            }
            if symbol == "" { continue }
            fmt.Printf("Found part number %v\n", num)
            partnum, err := strconv.Atoi(num)
            if err != nil {
                log.Fatal(err)
            }
            total += partnum 
        }
    }
    // Work in progress for day 2...
    /*
    for y, row := range symbols {
        for x, sym := range row {
            if sym != "*" { continue }
            startcol := 0
            if x > 0 { startcol = x - 1 }
            endcol := x + 1
            if y > maxX { endcol = maxX }
            startrow := 0
            if y > 0 { startrow = y -1 }
            endrow := y + 1
            if endrow > maxY -1 { endrow = maxY - 1 }
            fmt.Printf(
                "Searching for numbers between %v,%v and %v,%v (inclusive)\n",
                startrow, startcol, endrow, endcol)
            for row := startrow; row <= endrow; row++ {
                for col := startcol; col <= endcol; col++ {
                    
                }
            }
        }
    }
    */

    fmt.Printf("Part 1: Part Numbers Total: %v", total)
} 
