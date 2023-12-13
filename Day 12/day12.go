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
type Record struct {
    springs []rune
    grpDamaged []int
    matchCnt int
    matchExpr *regexp.Regexp
}
var p1Records []Record
var p2Records []Record
func main () {
    file, err := os.Open("example.txt")
    if err != nil { log.Fatal(err) }
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        recordParts := strings.Split(line, " ")
        if len(recordParts) != 2 {
            log.Fatal("Unexpected formatting in records")
        }
        lineRecord := Record { matchCnt: 0 }
        p2Record := Record { matchCnt: 0 }
        for _,spring := range recordParts[0] {
            lineRecord.springs = append(lineRecord.springs, spring)
        }
        p2Springs := strings.Repeat(string(lineRecord.springs) + "?", 5)
        p2Record.springs = []rune(p2Springs)
        strGrps := strings.Split(recordParts[1], ",")
        matchExp := `^\.*` 
        for idx,strGrp := range strGrps {
            grp, err := strconv.Atoi(strGrp)
            if err != nil { 
                log.Fatal("Unexpected formatting in records")
            }
            matchExp += fmt.Sprintf(`#{%v}`, grp)
            if idx < len(strGrps) - 1 {
                matchExp += `\.+`
            }
            lineRecord.grpDamaged = append(lineRecord.grpDamaged, grp)
        }
        matchExp +=`\.*$`
        lineRecord.matchExpr = regexp.MustCompile(matchExp)
        p1Records = append(p1Records, lineRecord)
        p2MatchExp := `^\.*`
        for iter := 0; iter < 5; iter++ {
            for _,grp := range lineRecord.grpDamaged {
                p2Record.grpDamaged = append(p2Record.grpDamaged, grp)
                p2MatchExp += fmt.Sprintf(`#{%v}`, grp)
            }
            if iter == 4 { p2MatchExp += `\.+` }
        }
        p2MatchExp += `.*$`
        p2Record.matchExpr = regexp.MustCompile(p2MatchExp) 
        p2Records = append(p2Records, p2Record)
    }
    p1TotalArrangements := 0
    for _, rec := range p1Records {
        GenerateArrangements(&rec)
        p1TotalArrangements += rec.matchCnt
    }
    p2TotalArrangements := 0
    for _, rec := range p2Records {
        GenerateArrangements(&rec)
        p2TotalArrangements += rec.matchCnt
    }
    fmt.Printf("Part 1: Total Posssible Arrangements: %v\n", 
        p1TotalArrangements)

}
func GenerateArrangements(rec *Record) {
    var qIdx []int
    for idx,spring := range rec.springs {
        if spring == '?' { 
            qIdx = append(qIdx, idx) 
        }
    }
    first := strings.ReplaceAll(string(rec.springs), "?", ".")
    if rec.matchExpr.MatchString(first) {
        rec.matchCnt++
    }
    GenAllArrangements(first, qIdx, rec)
}
func GenAllArrangements(source string, qIdx []int, rec *Record) {
    for mIdx,idx := range qIdx {
        next := []rune(source)
        next[idx] = '#'
        fmt.Printf("Testing arrangement %v with expression %v on source line %v\n",
            string(next), rec.matchExpr.String(), string(rec.springs))
        if rec.matchExpr.MatchString(string(next)) {
            rec.matchCnt++
            fmt.Printf("Arrangement %v matched expression %v on source line %v. ",
                string(next), rec.matchExpr.String(), string(rec.springs))
            fmt.Printf("Match count for this line is now %v.\n",
                rec.matchCnt)

        }
        if mIdx > len(qIdx) { continue }
        var nextQIdx []int = qIdx[mIdx + 1:]
        GenAllArrangements(string(next), nextQIdx, rec)
    }
}
