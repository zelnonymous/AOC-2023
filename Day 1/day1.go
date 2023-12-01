package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "regexp"
    "strconv"
)

func main() {
    p2exprs := make(map[string]*regexp.Regexp)
    p2exprs["one1one"], _ = regexp.Compile("one")
    p2exprs["two2two"], _ = regexp.Compile("two")
    p2exprs["three3three"], _ = regexp.Compile("three")
    p2exprs["four4four"], _ = regexp.Compile("four")
    p2exprs["five5five"], _ = regexp.Compile("five")
    p2exprs["six6six"], _ = regexp.Compile("six")
    p2exprs["seven7seven"], _ = regexp.Compile("seven")
    p2exprs["eight8eight"], _ = regexp.Compile("eight")
    p2exprs["nine9nine"], _ = regexp.Compile("nine")
    file, err := os.Open("input.txt")
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(file)
    p1total := 0
    p2total := 0
    for scanner.Scan() {
        line := scanner.Text()
        val, err := GetLineValue(line)
        if err != nil {
            log.Fatal(err)
            continue
        }
        p1total += val
        p2line := line
        for r, expr := range p2exprs {
            p2line = expr.ReplaceAllString(p2line, r)
        }
        modval, err := GetLineValue(p2line)
        if err != nil {
            log.Fatal(err)
            continue
        }
        p2total += modval
        fmt.Printf("New Value: %v, Running Total: %v\n", modval, p2total)
    }
    fmt.Printf("Part 1 Total: %v\n", p1total)
    fmt.Printf("Part 2 Total: %v\n\n", p2total)
    if err := scanner.Err(); err != nil {
        log.Fatal((err))
    }
    file.Close()
}
func GetLineValue(line string) (int, error) {
    numexpr, _ := regexp.Compile("[0-9]")
    nummatches := numexpr.FindAllString(line, -1)
    flvalue := fmt.Sprintf(
        "%v%v", 
        nummatches[0], 
        nummatches[len(nummatches) -1],
    )
    return strconv.Atoi(flvalue)
}
