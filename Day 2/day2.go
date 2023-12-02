package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
        result = append(result, ParsePull(pullstring))
    }
    return gameinfo[0], result
}
func ParsePull(pullstring string) pull {
    cubes := make(map[string]int)
    cubeinfo := strings.Split(pullstring, ",")
    for _, info := range cubeinfo {
        color, qty, err := ParsePullCubes(info)
        if err != nil {
            fmt.Println(err)
            continue
        }
        cubes[color] = qty
    }
    return pull {
        cubes: cubes,
    }
}
func ParsePullCubes(info string) (string, int, error) {
    details := strings.Split(info, " ")
    if len(details) < 2 {
        return "", 0, errors.New("Found cube pull without color or count")
    }
    color := strings.TrimSpace(details[0])
    qty, err := strconv.Atoi(details[1])
    if err != nil {
        return color, 0, errors.New("Warning: Cube quantity is not a number")
    }
    return color, qty, nil
}
