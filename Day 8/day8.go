package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)
// Type to store a node from the map
type Node struct {
    name string
    left *Node
    right *Node
}
// Expression for the instruction line from the input (only contains a series
// of instructions as "L" for left and "R" for right
var exprInst = regexp.MustCompile("^[RL]+$")
// Expression for node lines from the input.  These consist of a node name
// (several characters) followed by a space, an equal sign, another space,
// a left paren, the name of the node to the left, a comma, another space,
// the name of the node to the right, and a right paren.
var exprNode = regexp.MustCompile(`^([A-Z0-9]+)\s*=\s*\(([A-Z0-9]+),\s*([A-Z0-9]+)\)$`)
// Expression for start nodes in part 2 (start node names start with "A")
var exprStart = regexp.MustCompile(`A$`)
// Expression for end nodes in part 2 (end nodes end with "Z")
var exprEnd = regexp.MustCompile(`Z$`)
// A map to store nodes by name
var nodes = make(map[string]*Node)
// A signature type for a function capable of testing a node for some
// condition and return true or false
type NodeTest func(node *Node) bool
// A list of instruction characters
var instructions []string
func main() {
    file, err := os.Open("input.txt")
    if err != nil { log.Fatal(err) }
    scanner := bufio.NewScanner(file)
    /* Parser: As mentioned above, the interesting lines from the input are
    an instruction line (there should only be one of these) and node lines.
    For a node line, we extract three sets of characters: The name of the node,
    the name of the node to its left, and the name of the node to its right.
    If any of these three notes are not already in the map ("nodes") we add 
    them.  The Node type acts like a binary tree; Node.left is a reference to
    the left node and Node.right is a reference to the right node.  The last
    thing we do is assign those references and dump the node object to output
    for debugging. */
    for scanner.Scan() {
        line := scanner.Text()
        instMatch := exprInst.FindString(line) 
        nodeMatch := exprNode.FindStringSubmatch(line)
        if instMatch != "" {
            fmt.Printf("Found instructions: %v\n", line)
            instructions = strings.Split(line, "")
            continue
        }
        if nodeMatch != nil { 
            fmt.Printf("Found Node: %v\n", nodeMatch)
            if len(nodeMatch) <  4 {
                log.Fatal("Unexpected node string in input.")
            }
            _, found := nodes[nodeMatch[1]]
            if !found { 
                nodes[nodeMatch[1]] = &Node {name: nodeMatch[1]}
            }
            _, found = nodes[nodeMatch[2]]
            if !found { 
                nodes[nodeMatch[2]] = &Node {name: nodeMatch[2]} 
            }
            _, found = nodes[nodeMatch[3]]
            if !found { 
                nodes[nodeMatch[3]] = &Node {name: nodeMatch[3]} 
            }
            srcNode := nodes[nodeMatch[1]]
            srcNode.right = nodes[nodeMatch[3]]
            srcNode.left = nodes[nodeMatch[2]]
            fmt.Println(srcNode)
        }
    }
    /* Part 1: Starting from node "AAA", we follow the instructions sequentially
    until we reach node "ZZZ".  The GetStepCount function will return the number
    of steps required based on the test function we provide (for this part, our
    test for end node is node.name == "ZZZ." */
    startNode := nodes["AAA"]
    p1Steps := 0
    if startNode != nil {
        p1Steps = GetStepCount(startNode, 
            func(node *Node) bool { return node.name == "ZZZ"})
    }

    /* Part 2: This was a little more tricky.  We will be traversing multiple
    paths concurrently.  Any node that starts with an "A" is a starting node
    and any node that ends with a "Z" is an ending node.  We need to find the
    number of stpes for ALL paths to finish on an end node (one that ends with
    "Z"). Originally, I was trying to do this with loops (at least O(n^2)),
    but it was far to slow, so I figured I could find the number of steps
    required to reach the end for each path like we did in the first part, then
    find the least common multiple of those step counts to get the answer (which
    was much faster) */
    var p2Steps []int
    for _, node := range nodes {
        if exprStart.FindString(node.name) == "" {
            continue
        }
        nodeSteps := GetStepCount(
            node, 
            func(node *Node) bool { 
                return exprEnd.FindString(node.name) != ""
            })
        fmt.Printf("Starting at %v, needed %v steps\n", node.name,
            nodeSteps)
        p2Steps = append(p2Steps, nodeSteps)
    }
    p2Total := LCM(p2Steps)
    fmt.Printf("Part 1: Step Count: %v\n", p1Steps)
    fmt.Printf("Part 2: Step Count: %v\n", p2Total)
}
/* I refactored this after completing part 1.  This function takes a reference
to a starting node and a test function to identify an ending node.  It walks
the instruction list and executes each move until the test passes.  Per the
puzzle, if the end of the instructions are reached before we find the end node,
we loop back to the beginning of the instructions and continue.
Along the way, we track the number of steps we've taken and return the result
at the end. */
func GetStepCount (startNode *Node, test NodeTest) int {
    stepCount := 0
    currentNode := startNode
    for !test(currentNode) {
        for _, inst := range instructions {
            if inst == "L" {
                fmt.Printf("Moving from %v to %v\n",
                    currentNode.name, currentNode.left.name) 
                currentNode = currentNode.left
            }
            if inst == "R" {
                fmt.Printf("Moving from %v to %v\n",
                    currentNode.name, currentNode.right.name)
                currentNode = currentNode.right
            }
            stepCount++
            if test(currentNode) { break }
        }
    }
    return stepCount
}
// Given a slice of ints, find the least common multiple
// by iterating the list and finding the LCM of the current entry
// and the result of the prior entry (first entry is just the entry itself)
func LCM (vals []int) int {
    previous := 1
    for _, val := range vals {
       previous = val * previous / GCD(val, previous) 
    }
    return previous
}
// Greatest common divisor of two numbers
func GCD(a, b int) int {
    for b !=0 {
        t:= b
        b = a % b
        a = t
    }
    return a
}
