package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)
type Pattern struct {
    data [][]rune
    score int
}
var patterns []Pattern
func main() {
    file, err := os.Open("input.txt")
    if err != nil { log.Fatal(err) }
    scanner := bufio.NewScanner(file)
    current := Pattern {}
    lastLineWasData := false
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { 
            lastLineWasData = false
            patterns = append(patterns, current)
            current = Pattern {}
            continue
        }
        current.data = append(current.data, []rune(line))
        lastLineWasData = true
    }
    if lastLineWasData {
        patterns = append(patterns, current)
    }
    p1Total := 0
    p2Total := 0
    PATTERNLOOP:
    for i,pattern := range patterns {
        fmt.Printf("Checking pattern %v\n", i+1)
        for row := range pattern.data {
            fmt.Println(string(pattern.data[row]))
        }
        fmt.Println()
        FindMirrors(&pattern, -1, false)
        p1Total += pattern.score
        originalScore := pattern.score
        fmt.Println("Original:")
        for _,prow := range pattern.data {
            fmt.Println(string(prow))
        }
        for row := range pattern.data {
            for col := range pattern.data[row] {
                origChar := pattern.data[row][col]
                if origChar == '.' {
                    pattern.data[row][col] = '#'
                } else {
                    pattern.data[row][col] = '.'
                }
                FindMirrors(&pattern, originalScore, false)
                if pattern.score != originalScore {
                    fmt.Println("Successful Mutation:")
                    for _,prow := range pattern.data {
                        fmt.Println(string(prow))
                    }
                    fmt.Printf("Original score: %v, Mutation score: %v\n", originalScore,
                        pattern.score)
                    p2Total += pattern.score
                    continue PATTERNLOOP
                }
                pattern.data[row][col] = origChar
            }
        }

    }
    fmt.Printf("Part 1: Total: %v\n", p1Total)
    fmt.Printf("Part 2: Total: %v\n", p2Total)
}
func FindMirrors(pattern *Pattern, original int, debug bool) {
    SEARCHCOLUMNS:
    for col := 1; col < len(pattern.data[0]); col++ {
        for row := range pattern.data {
            left := string(pattern.data[row][0:col])
            right := string(pattern.data[row][col:])
            left = reverse(left)
            if len(left) < len(right) {
                right = right[0:len(left)]
            } else if len(right) < len(left) {
                left = left[0:len(right)]
            }
            if left != right || col == original {
                continue SEARCHCOLUMNS
            }
        }
        pattern.score = col
        return 
    }
    if debug {
        fmt.Println("No mirror columns found.  Checking rows.")
    }
    SEARCHROWS:
    for row := 1; row < len(pattern.data); row ++ {
        top := pattern.data[0:row]
        bottom := pattern.data[row:]
        if len(top) < len(bottom) {
            bottom = bottom[0:len(top)]
        }
        if len(bottom) < len(top) {
            top = top[len(top) - len(bottom):]
        }
        top = reverseRows(top)
        if debug {
            fmt.Printf("ROW: %v\n", row)
            fmt.Println("Top:           Bottom:")
            for tr := range top {
                fmt.Printf("%v",string(top[tr]))
                fmt.Printf("   %v\n",string(bottom[tr]))
            }
        }
        for row := range top {
            for col := range top[row] {
                if top[row][col] != bottom[row][col] || 
                    row * 100 == original {
                    continue SEARCHROWS
                }
            }
        }
        pattern.score = row * 100
        return
    }
    if debug {
        fmt.Println("WARNING: MIRROR NOT FOUND!")
    }
}
func reverseRows(input [][]rune) [][]rune {
    output := make([][]rune, len(input))
    for i, o := 0, len(input)-1; i < len(input); i, o = i+1, o-1 {
        output[o] = input[i]
    }
    return output
}
func reverse(input string) (result string) {
    for _,v := range input {
        result = string(v) + result
    } 
    return
}
