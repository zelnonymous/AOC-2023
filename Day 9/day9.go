package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)
// Expression to split by any amount of whitespace
var exprWS = regexp.MustCompile(`\s+`)
// Type to store nodes with the left and right node
// from which their value was derived
type Node struct {
	value int
	left  *Node
	right *Node
}
// A function delegate for testing some condition on
// a node and returning true/ false
type NodeTest func(node Node) bool

// Store input data.  Each "history" will be a slice of nodes
// that represent the end of the factor tree (all zeroes), and
// those histories will each be an entry in the outer slice.
var histories [][]Node

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
    /* Parser: Each line of input is a "history" with a series of
    numbers separated by  spaces.  I start by initializing a "Node" for each
    number. I built the "factor trees" inside of the parser.  The actual history
    set that gets added to histories is the end of the tree (all zeroes as 
    mentioned above), but with references to the left and right parent for that
    entry from which its value is derived. */
	for scanner.Scan() {
		line := scanner.Text()
		var history []Node
		for _, strHistory := range exprWS.Split(line, -1) {
			intHistory, err := strconv.Atoi(strHistory)
			if err != nil {
				fmt.Printf("Warning: %v is not a number\n", strHistory)
				continue
			}
			history = append(history, Node{value: intHistory})
		}
        // Loop until every node in the current layer has a value of zero. This
        // indicates we are at the end of the tree (root nodes).
		for !All(history, 
            func(node Node) bool {
			    return node.value == 0
		    },
        ) {
			history = GetFactors(history)
		}
        histories = append(histories, history)
	}
    /* Part 1: We don't actually need to modify the tree to get the answer.
    instead, we start with a value of 0 (the inserted child) and starting with
    the last root node, walk up the "right" nodes for each layer of the tree
    and increment the current value with the sum of the previous value and the 
    value of the right node.  The grand total is the sum of the individual
    totals. */
    totalHistValue := 0
    for _, history := range histories {
        value := 0
        stepLastNode := &history[len(history) -1]
        for stepLastNode.right != nil {
            value += stepLastNode.right.value
            stepLastNode = stepLastNode.right
        }
        totalHistValue += value
    }
    /* Part 2: This is exactly the same as part 1, but instead of starting with
    the last root node, we start with the first root node and walk up the tree
    using the "left" nodes.  The new value is the value of each left node minus
    the last value. */
    p2TotalHistValue := 0
    for _, history := range histories {
        value := 0
        stepLastNode := &history[0]
        for stepLastNode.left != nil {
            value = stepLastNode.left.value - value
            stepLastNode = stepLastNode.left
        }
        p2TotalHistValue += value
    }
    fmt.Printf("Part 1: Total of New History Values: %v\n", totalHistValue)
    fmt.Printf("Part 2: Total of New History Values: %v\n", p2TotalHistValue)
}
/* Given a set of source nodes, compose a new 'layer' of nodes that
refer back to the left and right parents from the source */ 
func GetFactors(source []Node) []Node {
	var result []Node
	for idx := range source {
        // If we are on the last node in the list, there won't be a child because
        // there's no "right" node for it to derive its value from.
		if idx == len(source)-1 {
			continue
		}
        // The value of the new node is the difference between the value of the
        // node at the current index and the value of the node at the next index.
		result = append(result, Node{
			left:  &source[idx],
			right: &source[idx+1],
			value: source[idx+1].value - source[idx].value,
		})
	}
	return result
}
/* Given a set of nodes and test delegate, this will return true only if ALL nodes
in the set pass the supplied test */
func All(source []Node, test NodeTest) bool {
	for _, node := range source {
		if !test(node) {
			return false
		}
	}
	return true
}
