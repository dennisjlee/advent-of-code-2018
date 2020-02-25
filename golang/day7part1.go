package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Node struct {
	Name string
	EdgesOut map[*Node]struct{}
	EdgesIn map[*Node]struct{}
}

func newNode(name string) *Node {
	return &Node{name, make(map[*Node]struct{}), make(map[*Node]struct{})}
}

var keyExists = struct{}{}

func parseGraph(fileName string) map[string]*Node {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	} else {
		defer file.Close()
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	regex := regexp.MustCompile("Step ([A-Z]) must be finished before step ([A-Z]) can begin\\.")

	graph := make(map[string]*Node)
	for scanner.Scan() {
		matches := regex.FindStringSubmatch(scanner.Text())
		srcName := matches[1]
		destName := matches[2]
		//fmt.Println(srcName, "->", destName)
		src, ok := graph[srcName]
		if !ok {
			src = newNode(srcName)
			graph[srcName] = src
		}

		dest, ok := graph[destName]
		if !ok {
			dest = newNode(destName)
			graph[destName] = dest
		}
		src.EdgesOut[dest] = keyExists
		dest.EdgesIn[src] = keyExists
	}
	return graph
}

func insertSortedNode(nodes []*Node, node *Node) []*Node {
	i := sort.Search(len(nodes), func(i int) bool { return nodes[i].Name >= node.Name })
	nodes = append(nodes, nil)
	copy(nodes[i+1:], nodes[i:])
	nodes[i] = node
	return nodes
}

func main() {
	graph := parseGraph(os.Args[1])

	var readyNodes = make([]*Node, 0)
	for _, node := range graph {
		if len(node.EdgesIn) == 0 {
			readyNodes = append(readyNodes, node)
		}
	}
	sort.Slice(readyNodes, func(i, j int) bool { return readyNodes[i].Name < readyNodes[j].Name })
	processedNames := make([]string, len(graph))
	for len(readyNodes) > 0 {
		next := readyNodes[0]
		readyNodes = readyNodes[1:]
		processedNames = append(processedNames, next.Name)
		for neighbor := range next.EdgesOut {
			delete(neighbor.EdgesIn, next)
			if len(neighbor.EdgesIn) == 0 {
				readyNodes = insertSortedNode(readyNodes, neighbor)
			}
		}
	}

	fmt.Println(strings.Join(processedNames, ""))
}
