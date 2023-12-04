package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
    "regexp"
    "strconv"
)
type card struct {
    myNums []int
    winNums []int
}
var cards = make(map[int]card)
var p2Scores = make(map[int]int)
var exprNum = regexp.MustCompile(`\d+`)

func main() {
    file, err := os.Open("input.txt")
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(file)
    /* 
    Parsing - Each line in the input represents a scratchoff card.
    The line starts with a card id followed by a colon, then a whitespace
    separated list of winning numbers and player numbers separated by a pipe.
    */
    for scanner.Scan() {
        line := scanner.Text()
        lineparts := strings.Split(line, ":")
        if len(lineparts) < 2 {
            log.Fatal("Invalid card found.")
        }
        cardnum := exprNum.FindString(lineparts[0])
        id, err := strconv.Atoi(cardnum)
        if err != nil {
            log.Fatal(err)
        }
        crd := card { } 
        numpart := strings.Split(lineparts[1], "|")
        winNumStr := exprNum.FindAllString(numpart[1], -1)
        myNumStr := exprNum.FindAllString(numpart[0], -1)
        for _, num := range myNumStr {
            n, err := strconv.Atoi(num)
            if err != nil { log.Fatal(err) }
            crd.myNums = append(crd.myNums, n)
        }
        for _, num := range winNumStr {
            n, err := strconv.Atoi((num))
            if err != nil { log.Fatal(err) }
            crd.winNums = append(crd.winNums, n)
        }
        cards[id] = crd
    }
    /*
    Part 1: points for each card is based on the number of player numbers
    that match winning numbers.  The card is worth 1 point for the initial 
    match and value doubles for each subsequent match. The answer is the sum of
    all card scores.
    */
    /*
    Revisited for Part 2: We need the total number of matches for each card,
    so we might as well get that here while we're iterating over the cards.
    */
    p1TotalPoints := 0
    for id, crd := range cards {
        p1val := 0
        p2val := 0
        for _, mynum := range crd.myNums {
            matches := false
            for _, winnum := range crd.winNums {
                if mynum == winnum {
                    matches = true
                    break
                }
            }
            if matches {
                p2val++
                if p1val == 0 {
                    p1val = 1
                    continue
                }
                p1val = p1val * 2
            }
        }
        p1TotalPoints += p1val
        p2Scores[id] = p2val
    }
    /*
    Part 2: for each matching number on a card, the player wins the cards that
    follow it in succession (eg. if card 2 has 3 matches, the player wins a copy
    of card 3, card 4, and card 5; if card 3 then has one match, their original
    card 3 wins a copy of card 4 and their copy of card 3 from card 2 also
    wins a copy of card 4). The answer is the total number of cards the player
    posseses after scoring. Recursion seemed like the simplest solution here,
    getting the won copies, then the one copies of each of those, etc. until
    a card that wins nothing is reached. We track the count of copies produced
    (including the original card itself) for each card in succession and sum
    those for the solution.
    */
    p2Cnt := 0
    for id := range p2Scores {
        copycnt := TotalCards(id)
        p2Cnt += copycnt
    }
    fmt.Printf("Part 1: Total Points: %v\n", p1TotalPoints)
    fmt.Printf("Part 2: Card count: %v\n", p2Cnt)
}
// Recusive function to get the total count of "won" cards
// for Part 2
func TotalCards(id int) int {
    cnt := 1 
    score := p2Scores[id]
    if score < 1 { return cnt }
    for copyid := id + 1; copyid <= id + score; copyid++ {
        // Diagnostics
        //fmt.Printf("card %v wins a copy of card %v\n", id, copyid)
        cnt += TotalCards(copyid)
    }
    return cnt
}


