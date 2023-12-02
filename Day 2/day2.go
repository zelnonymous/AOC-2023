package main

import (
	"bufio"
    "fmt"
	"log"
	"os"
	"strings"
    "strconv"
)
type pull struct {
    cubes map[string]int
}

func main() {
    games := make(map[string][]pull)
    file, err := os.Open("input.txt")
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        gameid, pulls := ParseGame(line) 
        games[gameid] = pulls
    }

}
func ParseGame(line string) (string, []pull) {
    gameinfo := strings.Split(line, ":")
    if len(gameinfo) < 2 {
        log.Fatal("Game with no data detected!")
    }
    var result []pull
    pullstrings := strings.Split(gameinfo[1], ";")
    for _, pullstring := range pullstrings {
        cubes := make(map[string]int)
        cubeinfo := strings.Split(pullstring, ",")
        for _, info := range cubeinfo {
            details := strings.Split(info, " ")
            if len(details) < 2 {
                fmt.Println("Warning: Found cube pull without color or count:")
                fmt.Println(details)
                continue
            }
            qty, err := strconv.Atoi(details[1])
            if err != nil {
                fmt.Println("Warning: Cube quantity is not a number")
                fmt.Println(details)
                continue
            }
            color := strings.TrimSpace(details[0])
            cubes[color] = qty
        }
        result = append(result, pull {
            cubes: cubes,
        })
    }
    return gameinfo[0], result
}
