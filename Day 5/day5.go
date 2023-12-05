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
type MapRange struct {
    srcStart int
    destStart int
    rangeLength int
}
type MapEntry struct {
    destType string
    mapRanges []MapRange
}
type P2SeedRange struct {
    start int
    len int
}
var maps = make(map[string]*MapEntry)
var p1seeds []int
var p2seeds []P2SeedRange

func main() {
    var exprs = make(map[string]*regexp.Regexp)
    exprs["seed"] = regexp.MustCompile(`^seeds:\s*((\d+\s*)+)`)
    exprs["mapName"] = regexp.MustCompile(`^([a-z\-]*) map:`)
    exprs["mapVals"] = regexp.MustCompile(`^((\d+\s*)+)`)
	file, err := os.Open("input.txt")
    if err != nil { log.Fatal(err) }
    scanner := bufio.NewScanner(file)
    /* Parser: This was modified to accomodate part 2.  We're interested
    in input lines that match one of three expressions: seeds (which should
    just be the first line, but we'll expression match it anyway), map names
    (which has a dash separated source-to-destination format followed by
    the word 'map' and a colon), and map values (series of data points
    pertaining to the most recent map name line, formatted as a source range 
    start, a destination range start, and a lenght that applies to both ranges)
    Originally, I would build the entire range as a map of int to int, but that
    overextended memory for the real puzzle input, so I am instead just tracking
    the starts and ranges.
    For part 2, we also need to handle the seed line a little differently.
    We track pairs of numbers, where the first is the seed range start and the
    second is the seed range length.
    */
    curSrc := ""
    for scanner.Scan() {
        line := scanner.Text()
        for exprName, expr := range exprs {
            match := expr.FindStringSubmatch(line)
            if match == nil { continue }
            if exprName == "seed" {
                currentRange := P2SeedRange {}
                for idx, seed := range strings.Split(match[1], " ") {
                    seednum, err := strconv.Atoi(seed)
                    if err != nil { log.Fatal(err) }
                    p1seeds = append(p1seeds, seednum)
                    // Part 2
                    if idx % 2 == 0 {
                        currentRange.start = seednum
                    } else {
                        currentRange.len = seednum
                        p2seeds = append(p2seeds, currentRange)
                        currentRange = P2SeedRange {}
                    }
                }
                continue
            }
            if exprName == "mapName" {
                mapName := strings.Split(match[1], "-")
                if len(mapName) < 3 { 
                    fmt.Printf("Unexpected map name %v\n", match[1])
                    panic("ERROR")
                }
                curSrc = mapName[0]
                maps[curSrc] = &MapEntry { destType: mapName[2], }
                continue
            }
            if exprName == "mapVals" {
                vals := strings.Split(match[1], " ")
                if len(vals) < 3 {
                    fmt.Printf("Unexpected map value format %v\n", match[1])
                    panic("ERROR")
                }
                destStart, err := strconv.Atoi(vals[0])
                if err != nil { continue }
                srcStart, err := strconv.Atoi(vals[1]) 
                if err != nil { continue }
                rangeLen, err := strconv.Atoi(vals[2])
                if err != nil { continue }
                maps[curSrc].mapRanges = append(
                    maps[curSrc].mapRanges, 
                    MapRange { 
                        srcStart: srcStart, 
                        destStart: destStart, 
                        rangeLength: rangeLen,
                    },
                )
            }
        } 
    } 
    fmt.Println("The following maps were discovered:")
    for src, mapentry := range maps {
        fmt.Printf("Source: %v, Dest: %v\n", src, mapentry.destType)
        for _, rng := range mapentry.mapRanges {
            fmt.Printf("Source Start: %v ", rng.srcStart)
            fmt.Printf("Dest Start: %v ", rng.destStart)
            fmt.Printf("Range: %v ", rng.rangeLength)
            fmt.Println()
        }
        fmt.Println()
    }
    fmt.Println("Part 2 Seed Ranges:")
    for _, seedRng := range p2seeds {
        fmt.Printf("Start: %v, Length: %v\n", seedRng.start, seedRng.len)
    }
    fmt.Println()
    /* Part 1: For each seed, we traverse each map sequentially until we
    reach the location map. We want to find the lowest available location number
    for the seeds that we have. */
    minLocation := 0
    for _, seed := range p1seeds {
        locval := GetSeedLocation(seed) 
        if minLocation == 0 || locval < minLocation {
            minLocation = locval
        }
    }
    /* Part 2: We need to do the same thing as part one, but now for every
    seed in each seed range.  I iterate over the ranges and then get the 
    location value for each seed in the range.  This did get me the correct
    answer, but it had a long and grindy run time.  In retrospect, if I 
    considered the number of seeds involved, I likely would have taken a 
    different approach.  I would have sorted the location numbers smallest to
    largest and walked them to find the first one that tied to a valid seed
    (within any range, based on the difference between start and start + length)
    */
    minLocation2 := 0 
    for _, seedRng := range p2seeds {
        fmt.Printf("Searching range starting at %v with length %v\n",
            seedRng.start, seedRng.len)
        for currentSeed := seedRng.start; 
            currentSeed < seedRng.start + seedRng.len; 
            currentSeed++ {
            locval := GetSeedLocation(currentSeed)
            if minLocation2 == 0 || locval < minLocation2 {
                minLocation2 = locval
            }
        }
    }
    fmt.Printf("Part 1: Closest Location: %v\n", minLocation)
    fmt.Printf("Part 2: Closest Location: %v\n", minLocation2)
}
// Function to walk the maps to get the location value based on a seed value
func GetSeedLocation(seed int) int {
    srcType := "seed"
    srcVal := seed
    for srcType != "location" {
        destType := maps[srcType].destType
        destVal := 0
        for _, rng := range maps[srcType].mapRanges {
            if srcVal < rng.srcStart || 
            srcVal > rng.srcStart + rng.rangeLength - 1 {
                continue
            }
            destVal = rng.destStart + (srcVal - rng.srcStart)
        }
        if destVal == 0 { destVal = srcVal }
        /* 
        //Debug logging 
        fmt.Printf(
            "%v %v mapped to %v %v\n", 
            srcType, srcVal,
            destType, destVal,
            )
        */
        srcType = destType 
        srcVal = destVal
    }
    return srcVal
}
