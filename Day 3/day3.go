package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
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
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) > maxX {
            maxX = len(line)
        }
        lines = append(lines, line)
        maxY++
    }
    nums = make([][]string, maxY)
    symbols = make([][]string, maxY)
    for y := range nums {
        nums[y] = make([]string, maxX)
        symbols[y] = make([]string, maxX)
    }
    for y, line := range lines {
        numStart := 0
        numMatch := false
        numCurrent := ""
        for x, ch := range strings.Split(line, "") {
            if ch == "." {
                if !numMatch { continue }
                /*
                fmt.Printf(
                    "Found number %v at row %v starting at column %v\n",
                    numCurrent, y, numStart,
                )
                */
                nums[y][numStart] = numCurrent 
                numCurrent = ""
                numMatch = false
                continue
            }
            if exprDigit.MatchString(ch) {
                numMatch = true
                numCurrent += ch
                numStart = x
                continue
            }
            /*
            fmt.Printf(
                "Found symbol %v at row %v column %v\n",
                ch, y, x,
            )
            */
            symbols[y][x] = ch
            if !numMatch { continue }
            nums[y][numStart] = numCurrent 
            numCurrent = ""
            numMatch = false
        }
    }
    total := 0
    for y, row := range nums {
        for numStart, num := range row  {
            if num == "" { continue }
            startcol := 0
            if numStart > 0 { startcol = numStart - 1 }
            endcol := startcol + len(num)
            if endcol > maxX - 1 { endcol = maxX - 1 }
            startrow := 0
            if y > 0 { startrow = y - 1 }
            endrow := y + 1
            if endrow > maxY - 1 { endrow = maxY - 1 }
            /* 
            fmt.Printf(
                "Number %v on row %v starts at col %v and ends at col %v\n",
                num, y, numStart, startcol + len(num))
            fmt.Printf(
                "We will check for symbols between %v,%v and %v,%v (inclusive)\n",
                startrow, startcol, endrow, endcol)
            */
            symbol := ""
            for row := startrow; row <= endrow; row++ {
                for col := startcol; col <= endcol; col++ {
                    if symbols[col][row] == "" { continue }
                    symbol = symbols[col][row]
                    break;
                }
                if symbol != "" { break }
            }
            if symbol == "" { continue }
            partnum, err := strconv.Atoi(num)
            if err != nil {
                log.Fatal(err)
            }
            total += partnum 
        }
    }
    fmt.Printf("Part 1: Part Numbers Total: %v", total)
} 
