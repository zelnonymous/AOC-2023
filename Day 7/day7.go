package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)
// Type to store a hand of cards and its bid
type Hand struct {
    cards []string
    bid int
}
// For part 1, card trump value is in order with aces high.
var P1CardVals = map[string]int {
    "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, 
    "9": 9, "T": 10, "J": 11, "Q": 12, "K": 13, "A": 14,
}
// For part 2, J changes from jack to joker and becomes wild.
// As a result, it gets the lowest possible trump value.
var P2CardVals = map[string]int {
    "J": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, 
    "8": 8, "9": 9, "T": 10, "Q": 11, "K": 12, "A": 13,
}
// Slice of hands from input
var hands []Hand

func main() {
    file,err := os.Open("input.txt")
    if err != nil { log.Fatal(err) }
    scanner := bufio.NewScanner(file)
    /* Parser: This one is very simple.  Each line of the input is a series 
    of letters and numbers that each represent a card, then space as a 
    separator, and finally a number representing a bid for that hand. */
    for scanner.Scan() {
        line := scanner.Text()
        handParts := strings.Split(line, " ")
        bid, err := strconv.Atoi(handParts[1])
        if err != nil { log.Fatal(err) }
        hands = append(hands, Hand {
            cards: strings.Split(handParts[0], ""),
            bid: bid,
        })
    }
    /* For part 1, we need to sort our list of hands based on a set of rules.
    first, a hand gets a "HandTypeValue" based on typical poker rules.  Scored
    highest to lowest:
    - Five of a kind
    - Four of a kind
    - Full house
    - Three of a kind
    - Two pair
    - Single pair
    - High card
    If two hands are of the same type, we compare card-by-card left-to-right
    until one card trumps the other (based on the P1CardVals map)*/
    slices.SortFunc(hands, func(a, b Hand) int {
        aVal := GetHandTypeValue(a.cards, false)
        bVal := GetHandTypeValue(b.cards, false)
        if aVal != bVal { return cmp.Compare(aVal, bVal) }
        for idx := range a.cards {
            if a.cards[idx] == b.cards[idx] { continue }
            return cmp.Compare(P1CardVals[a.cards[idx]], 
                P1CardVals[b.cards[idx]])
        }
        log.Fatalf("Hands could not be compared: %v %v", 
            strings.Join(a.cards, ""),
            strings.Join(b.cards, ""))
        return 0
    })
    /* Once we have the hands in the correct order, we just iterate over them
    and multiply the rank (index in the list plus one) and its bid.  The sum of
    these values is the answer.
    I added an extra call to GetHandTypeValue here which isn't necessary to 
    solve, but was helpful for debugging. */
    p1TotalWinnings := 0 
    for idx, hnd := range hands {
        score := GetHandTypeValue(hnd.cards, false)
        fmt.Printf("Hand %v has hand type score %v and rank %v\n", 
            strings.Join(hnd.cards, ""), score, idx + 1)
        p1TotalWinnings += (idx + 1) * hnd.bid
    }
    fmt.Println() 
    p2TotalWinnings := 0
    /* Part 2: Again, we sort the list based on a set of rules, but
    'J' is now a wild joker, so the formulation of GetHandTypeValue is going
    to be different.  I added a second parameter to indicate whether jokers are
    wild.  Also, for the second scoring rule (trump card), J's now have the 
    lowest possible value, so I use the modified P2CardVals map. */
    slices.SortFunc(hands, func(a, b Hand) int {
        aVal := GetHandTypeValue(a.cards, true)
        bVal := GetHandTypeValue(b.cards, true)
        if aVal != bVal { return cmp.Compare(aVal, bVal) }
        for idx := range a.cards {
            if a.cards[idx] == b.cards[idx] { continue }
            return cmp.Compare(P2CardVals[a.cards[idx]], 
                P2CardVals[b.cards[idx]])
        }
        log.Fatalf("Hands could not be compared: %v %v", 
            strings.Join(a.cards, ""),
            strings.Join(b.cards, ""))
        return 0
    })
    for idx, hnd := range hands {
        score := GetHandTypeValue(hnd.cards, true)
        fmt.Printf("Hand %v has hand type score %v and rank %v\n", 
            strings.Join(hnd.cards, ""), score, idx + 1)
        p2TotalWinnings += (idx + 1) * hnd.bid
    }
    fmt.Printf("Part 1: Total Winnings: %v\n", p1TotalWinnings)
    fmt.Printf("Part 2: Total Winnings: %v\n", p2TotalWinnings)
}
/*
For a given hand of 5 cards, get a value representing the type score for hand.
the types of hands are described above.  I assign a value of 6 for five of a
kind, 5 for four of a kind, 4 for full house, 3 for three of a kind... down to
0 for high card only.
I added a flag here for part 2 to change the behavior if we are playing with 
jokers.
*/
func GetHandTypeValue(cards []string, jokers bool) int {
    if len(cards) != 5 { 
        log.Fatal("INVALID HAND: " + strings.Join(cards, "")) 
    }
    // The first step is to iterate the hand and, for each card value, get
    // the count of that card in the hand
    crdCnts := make(map[string]int)
    for _, card := range cards {
        crdCnts[card]++
    }
    // If we are playing with jokers, we'll need to know how many we have
    jokerCnt := crdCnts["J"]
    // If we have 5 "J", we have five of a kind regardless of whether they are 
    // considered jacks or jokers
    if jokerCnt == 5 { return 6 }
    // Any other natural 5 of a kind (unaffected by jokers)
    if len(CardsByCount(5, crdCnts, jokers)) == 1 { return 6 }
    if len(CardsByCount(4, crdCnts, jokers)) == 1 { 
        // If we are playing with jokers, 
        // four of a kind + 1 joker = 5 of a kind
        if jokers && jokerCnt == 1 { return 6 }
        // If we are not playing with jokers or we don't have any in hand,
        // this is a four of a kind.
        return 5 
    }
    if len(CardsByCount(3, crdCnts, jokers)) == 1 {
        // If we are playing with jokers, 
        // three of a kind + 2 jokers = 5 of a kind
        if jokers && jokerCnt == 2 { return 6 }
        // If we are playing with jokers, 
        // three of a kind + joker = 4 of a kind
        if jokers && jokerCnt == 1 { return 5 }
        // If we have 3 and 2 of a kind it's a
        // natural full house
        if len(CardsByCount(2, crdCnts, jokers)) == 1 {
            return 4
        }
        // If none of the above apply, then this is a three of a kind.
        return 3
    }
    // If we made it here, we are interested in pairs which we could either
    // have one or two of.  
    pairs := CardsByCount(2, crdCnts, jokers)
    // Start with the case of a single pair
    if len(pairs) == 1 {
        // The only way to form a full house with a single pair is
        // 3 jokers, but that's also a 5 of a kind which is better
        // Pair + 3 jokers = 5 of a kind
        if jokers && jokerCnt == 3 { return 6 }
        // Pair + 2 jokers = 4 of a kind
        if jokers && jokerCnt == 2 { return 5 }
        // Pair + 1 joker = 3 of a kind
        if jokers && jokerCnt == 1 { return 3 }
        // This is just a single pair
        return 1
    }
    // We have two natural pairs
    if len(pairs) == 2 {
        // Two pairs plus a joker = full house
        if jokers && jokerCnt == 1 { return 4 }
        // Otherwise, this is scores as two pair
        return 2
    }
    // High card plus 4 jokers = 5 of a kind
    if jokers && jokerCnt == 4 { return 6 }
    // High card plus 3 jokers = 4 of a kind
    if jokers && jokerCnt == 3 { return 5 }
    // High card plus 2 jokers = 3 of a kind
    if jokers && jokerCnt == 2 { return 3 }
    // High card plus 1 joker is a pair
    if jokers && jokerCnt == 1 { return 1 }
    return 0
}
// Get a list of cards that we have exactly 'count' of from the map of cards
// to counts.
func CardsByCount(count int, 
    crdCnts map[string]int,
    ignoreJokers bool) []string {
    var cardsWithCount []string 
    for crd, crdCnt := range crdCnts {
        // If we are playing with jokers, ignore them.  We'll
        // count and apply them separately.
        if ignoreJokers && crd == "J" { continue }
        if crdCnt == count {
            cardsWithCount = append(cardsWithCount, crd)
        }
    }
    return cardsWithCount
}
