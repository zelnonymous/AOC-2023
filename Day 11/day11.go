package main
import (
    "bufio"
    "fmt"
    "log"
    "os"
    "math"
)
// Type to store information about a galaxy (mainly x, y cooridnates, but also
//added an id for troubleshooting 
type Galaxy struct {
    id int
    x int
    y int
}
// Type to store pairs of galaxies for path checking
type Pair[T any] struct {
    First T
    Second T
}
func main() { 
    file, err := os.Open("input.txt") 
    if err != nil { log.Fatal(err) } 
    scanner := bufio.NewScanner(file) 
    var image [][]rune 
    /* Parser: Another fairly simple one.  The input represents a telescope
    image.  Each line should have the same number of characters which are either
    '.' (representing space) or '#' representing a galaxy.*/
    for scanner.Scan() {
        line := scanner.Text()
        row := make([]rune, len(line))
        for col, c := range line {
            row[col] = c
        }
        image = append(image, row)
    }
    // For diagnostics, echo the original image
    fmt.Println("Original:")
    for row := range image {
        for col := range image[row] {
            fmt.Printf(string(image[row][col]))
        }
        fmt.Println()
    }
    // The expand function will apply expansion to space and return 
    // the full expansion as well as a slice of galaxies contained within
    // For part 1 only, echo this for troublshooting.
    expanded, galaxies := Expand(image, 2)
    fmt.Println("Expanded:")
    for row := range expanded {
        for col := range expanded[row] {
            fmt.Printf(string(expanded[row][col]))
        }
        fmt.Println()
    }
    // We only want to compare each pair of galaxies once, so lets
    // build that as a slice of pairs
    var pairs []Pair[Galaxy]
    for idx,galaxy := range galaxies {
        for idx2 := idx + 1; idx2 < len(galaxies); idx2++ {
            pairs = append(pairs, Pair[Galaxy]{ 
                First: galaxy, 
                Second: galaxies[idx2],
            })
        }
    }
    /* Part 1: given the expanded image, find the shortest distance between
    each pair of galaxies.  This should just be the sum of the difference in
    x and the difference in y.  Print each distance for troubleshooting and
    sum the results. */
    p1ShortestPaths := 0.00
    for _, set := range pairs {
        dist := math.Abs(float64(set.First.x - set.Second.x)) + 
            math.Abs(float64(set.First.y - set.Second.y))
        fmt.Printf("Distance between %v and %v: %v\n", 
            set.First.id, set.Second.id, dist)
        p1ShortestPaths += dist
    }
    /* Part 2: Same scenario, but instead of increasing empty space rows and
    columns by 2x, we now increase it by 1000000x.  The resulting map is
    obviously fairly large, but we don't even need it; We just need the 
    positions of the galaxies.  Discard the map and rebuild the pairs */
    _, galaxies = Expand(image, 1000000) 
    var p2pairs []Pair[Galaxy]
    for idx,galaxy := range galaxies {
        for idx2 := idx + 1; idx2 < len(galaxies); idx2++ {
            p2pairs = append(p2pairs, Pair[Galaxy]{ 
                First: galaxy, 
                Second: galaxies[idx2],
            })
        }
    }
    /* Same as part 1.  Calculate the distance between each distinct pair of
    galaxies and sum those paths. */
    p2ShortestPaths := 0.00
    for _, set := range p2pairs {
        dist := math.Abs(float64(set.First.x - set.Second.x)) + 
            math.Abs(float64(set.First.y - set.Second.y))
        fmt.Printf("Distance between %v and %v: %v\n", 
            set.First.id, set.Second.id, dist)
        p2ShortestPaths += dist
    }
    fmt.Printf("Part 1: Total of shortest distances: %f\n", p1ShortestPaths)
    fmt.Printf("Part 2: Total of shortest distances: %f\n", p2ShortestPaths)
}
/* Function to expand the universe. Excercise caution ;)
...But seriously, given the source telescope image and the number of copies
that should be created by a row or column that is all space characters ('.'),
produce a new image and a slice of the galaxies contained therein with their
coordinates. */
func Expand(image [][]rune, expansionCopies int) ([][]rune, []Galaxy) {
    var result [][]rune
    var copycols []int
    var galaxies []Galaxy
    gid := 1
    // The approach we took here won't work if the input is empty,
    // so we'll just fire that back with an empty slice of galaxies.
    if len(image) < 1 { return image, galaxies }
    /* First, we need to know where we need more columns in the output
    than we have in the input.  We'll build copycols as a list of 
    column indexes where the column is all space. */
    COLCHECK:
    for col := range image[0] {
        for row := range image {
            if image[row][col] != '.' {
                continue COLCHECK
            }
        }
        copycols = append(copycols, col) 
    }
    // Now that we know which columns will need duplicated, we'll walk the whole
    // image and build the output.
    for row := range image {
        var rowchars []rune
        allSpace := true
        for col := range image[row] {
            rowchars = append(rowchars, image[row][col])
            // This column is all space, so we need to copy it into the 
            // output based on the supplied nubmer of copies
            if Contains[int](copycols, col) {
                for copy := 1; copy < expansionCopies; copy ++ {
                    rowchars = append(rowchars, image[row][col])
                }
            }
            // We'll do our row detection inline here.  if a row is
            // all space ('.'), we'll want to duplicate it in the output
            // the supplied number of times
            if image[row][col] != '.' { allSpace = false }
            // If this is a galaxy, we'll want to assign it an ID and capture
            // its coordinates.
            if image[row][col] == '#' {
                galaxies = append(galaxies, Galaxy {
                    id: gid,
                    x: len(rowchars) - 1,
                    y: len(result),
                })
                gid++
            }
        }
        result = append(result, rowchars)
        if !allSpace { continue }
        for copy := 1; copy < expansionCopies; copy ++ {
            result = append(result, rowchars)
        }
    }
    return result, galaxies
}
// Helper function to see if a slice contains a supplied value
func Contains[T comparable](list []T, val T) bool {
    for _,lval := range list {
        if lval == val { return true }
    }
    return false
}
