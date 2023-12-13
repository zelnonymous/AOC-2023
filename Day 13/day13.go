package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)
type Pattern struct {
    data [][]rune
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
    SEARCHPATTERNS:
    for i,pattern := range patterns {
        fmt.Printf("Checking pattern %v\n", i+1)
        for row := range pattern.data {
            fmt.Println(string(pattern.data[row]))
        }
        fmt.Println()
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
                if left != right {
                    continue SEARCHCOLUMNS
                }
            }
            p1Total += col
            fmt.Printf("Found mirror at column %v, adding %v to total\n\n", 
                col, col)
            continue SEARCHPATTERNS
        }
        fmt.Println("No mirror columns found.  Checking rows.")
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
            fmt.Printf("ROW: %v\n", row)
            fmt.Println("Top:           Bottom:")
            for tr := range top {
                fmt.Printf("%v",string(top[tr]))
                fmt.Printf("   %v\n",string(bottom[tr]))
            }
            for row := range top {
                for col := range top[row] {
                    if top[row][col] != bottom[row][col] {
                        continue SEARCHROWS
                    }
                }
            }
            rowval := row * 100
            fmt.Printf("Found mirror at row %v, adding %v to total\n\n", 
                row, rowval)
            p1Total += rowval 
            continue SEARCHPATTERNS
        }
        log.Panic("NO MIRROR FOUND!")
    }
    fmt.Printf("Part 1: Total: %v\n", p1Total)

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
