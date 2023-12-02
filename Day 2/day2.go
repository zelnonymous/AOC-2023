package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)
type pull struct {
    cubes map[string]int
}

func main() {
    p1limits := map[string]int{
        "red": 12,
        "green": 13,
        "blue": 14,
    }
    p1possible := 0
    games := make(map[int][]pull)
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
    ordered_gameids := make([]int, 0)
    for gameid := range games {
        ordered_gameids = append(ordered_gameids, gameid)
    }
    sort.Ints(ordered_gameids)
    for _, gameid := range ordered_gameids {
        pulls := games[gameid]
        possible := true
        for _, pullinfo := range pulls {
            for color, qty := range pullinfo.cubes {
                if qty > p1limits[color] {
                    fmt.Printf(
                        "In game %v, pulled %v %v but limit was %v\n",
                        gameid,
                        qty,
                        color,
                        p1limits[color],
                    )
                    possible = false
                    break
                }
            }
            if possible == false {
                break
            }
        }
        if possible {
            fmt.Printf(
                "Game %v: Possible, Running Total: %v\n", 
                gameid, 
                p1possible)
            p1possible += gameid
        } 
    }
    fmt.Printf("Part 1: Sum of possible game IDs: %v\n", p1possible)

}
func ParseGame(line string) (int, []pull) {
    gameinfo := strings.Split(line, ":")
    if len(gameinfo) < 2 {
        log.Fatal("Game with no data detected!")
    }
    idstr := strings.TrimLeft(gameinfo[0], "Game ")
    id, err := strconv.Atoi(idstr)
    if err != nil {
        log.Fatal("Game with non-numeric ID detected!")
    }
    var result []pull
    pullstrings := strings.Split(gameinfo[1], ";")
    for _, pullstring := range pullstrings {
        result = append(result, ParsePull(pullstring))
    }
    return id, result
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
    info = strings.TrimSpace(info)
    details := strings.Split(info, " ")
    if len(details) < 2 {
        return "", 0, errors.New("Found cube pull without color or count")
    }
    color := strings.TrimSpace(details[1])
    qty, err := strconv.Atoi(details[0])
    if err != nil {
        return color, 0, errors.New("Warning: Cube quantity is not a number")
    }
    return color, qty, nil
}
