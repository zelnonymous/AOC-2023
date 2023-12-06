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
type Race struct {
    time int
    distance int
}
var p1races []Race
var expr = regexp.MustCompile(`^(Time:|Distance:)\s*((\d+\s*)+)$`)
var digits = regexp.MustCompile(`\d+`)

func main() {
    file, err  := os.Open("input.txt")
    if err != nil { log.Fatal(err) }
    scanner := bufio.NewScanner(file)
    /* Parser: This was modified to accomodate part 2.  The input is expected 
    to be two lines, the first marked 'Time:' with a series of numbers, and the
    second marked 'Distance:' with a series of numbers.  For part 1, each of
    these pairs indicates a duration and distance record for a single race. We
    capture each number in a slice.  We expect the time slice and the distance
    slice to be the same length, and the same index in each should be the pair
    of numbers for a race.
    For part 2, the premise changes a bit.  There is only one big race and the
    spaces between the numbers is to be ignored.  We'll just allocate a couple
    of strings and use concatentation on the string form of the numbers to 
    elminate the spaces and get the time and distance for this big race as 
    strings. */
    var times []int
    var distances []int
    p2timestr := ""
    p2distancestr := ""
    for scanner.Scan() {
        line := scanner.Text()
        match := expr.FindStringSubmatch(line)
        if match == nil || len(match) < 3 { 
            log.Fatal("Unexpected input.")
        }
        strvals := digits.FindAllString(match[2], -1) 
        for _, val := range strvals {
            nval, err := strconv.Atoi(val)
            if err != nil { log.Fatal(err) }
            if (strings.Contains(match[1], "Time")) {
                times = append(times, nval)            
                p2timestr += val
            } else {
                distances = append(distances, nval)
                p2distancestr += val
            }
        }
    }
    // Part 1: make sure the length of the time and distance slices matches
    // (numbers must be in pairs)
    if len(times) != len(distances) {
        log.Fatal("Count mismatch between times and distances.")
    }
    // Part 1: we could probably just index into each of these lists,
    // but I had already defined a nice type for it, so we'll just reorganize
    // them into a single slice of Race
    for idx := range times {
        p1races = append(p1races, Race {
            time: times[idx],
            distance: distances[idx],
        })
    }
    /* Part 1: Calculate the total margin of error.  We can hold the button
    for any duration between 0 and the full duration of the race.  For each 
    millisecond the button is held, we increase our speed by 1mm.  We can
    calculate the distance we'll travel by holding the button for each possible
    duration with holdms * (race.time - holdms).  In each case where the result
    would exceed the record (race.distance), we count that as a "winOpt". The
    total margin of error is the product of the "winOpts" from each race. */
    totalMargin := 1
    for idx, race := range p1races {
        winOpts := 0
        fmt.Printf("Race %v: Time %v, Distance %v\n", idx + 1,
            race.time, race.distance)
        for holdms := 0; holdms <= race.time; holdms++ {
            distance := holdms * (race.time - holdms)
            if distance > race.distance {
                fmt.Printf("Winner! Hold for %v.  Distance: %v\n", 
                    holdms, distance)
                winOpts++
            }
            // This optimization was required for part 2, but I added it to
            // part 1 as well since it would be in a shared function if I
            // bothered to refactor
            // If we have found winning options but are no longer winning,
            // holding the button for a longer period of time is gauranteed
            // to not produce a win (values on a graph would be parabolic).
            // we're safe to break out of the loop here to avoid wasting cylces.
            if distance <= race.distance && winOpts > 0 { break }
        }
        fmt.Printf("Winnable: %v\n\n", winOpts)
        totalMargin *= winOpts
    }
    /* Part 2: Same scenario, except that it's one big race and we just need to
    calculate our total "winOpts".  We'll first need to conver the strings we
    captured and concatenated earlier into numbers, then apply the same loop as
    earlier to see which ones will result in a win. */
    p2time, err := strconv.Atoi(p2timestr)
    if err != nil { log.Fatal(err) }
    p2distance, err := strconv.Atoi(p2distancestr)
    if err != nil { log.Fatal(err) }
    fmt.Printf("Part 2 Race was %v time and %v distance\n", p2time, p2distance)
    p2opts := 0
    for holdms := 0; holdms <= p2distance; holdms++ {
        distance := holdms * (p2time - holdms)
        if distance > p2distance { p2opts++ }
        if distance <= p2distance && p2opts > 0 { break }
    }

    fmt.Printf("Part 1: Total margin: %v\n", totalMargin)
    fmt.Printf("Part 2: Available options: %v\n", p2opts)
}
