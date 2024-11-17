package main

import (
	"fmt"
)

type Node struct {
	State    string
	Parent   *Node
	Action   string
	PathCost int
}

type Problem struct {
	InitialState Node                             // In(state)
	Actions      func(string) []string            // Actions(state) -> []string, Actions descriptions for each state
	Result       func(string, string) string      // Result(state, action) -> state, transition model
	GoalTest     func(string) bool                // GoalTest(state) -> bool, func that returns if a given state is a goal state
	Cost         func(string, string, string) int // Cost(stateA, action, stateB) -> int, cost function that returns the cost to reach stateB from stateA using action
}

type Solution struct {
	Actions []string // Actions to reach the goal state
}

func PopFifo[T any](arr *[]T) (T, error) {
	if len(*arr) == 0 {
		var zeroValue T
		return zeroValue, fmt.Errorf("Array is empty")
	}
	el := (*arr)[0]
	*arr = (*arr)[1:]
	return el, nil
}

func IsStateInList(state string, list []Node) bool {
	for _, node := range list {
		if node.State == state {
			return true
		}
	}
	return false
}

func ChildNode(problem Problem, parent Node, action string) Node {
	childState := problem.Result(parent.State, action)
	return Node{State: childState, Parent: &parent, Action: action, PathCost: problem.Cost(parent.State, action, childState)}
}

func SolutionPath(node Node) *Solution {
	actions := []string{}
	for node.Parent != nil {
		actions = append(actions, node.Action)
		node = *node.Parent
	}
	return &Solution{Actions: actions}
}

func TreeSearch(problem Problem) (*Solution, error) {
	// Initialize the frontier using the initial state of the problem
	frontier := []Node{problem.InitialState}
	actionsTaken := []string{}
	for len(frontier) > 0 {
		// Choose a leaf node and remove it from the frontier, we will choose the FIFO approach
		currentNode, _ := PopFifo(&frontier)

		// If the node contains a goal state, return the corresponding solution
		if problem.GoalTest(currentNode.State) {
			return &Solution{Actions: actionsTaken}, nil
		}

		// Expand the chosen node, adding the resulting nodes to the frontier
		for _, action := range problem.Actions(currentNode.State) {
			childNode := ChildNode(problem, currentNode, action)
			frontier = append(frontier, childNode)
			actionsTaken = append(actionsTaken, action)
		}
	}

	return nil, fmt.Errorf("No solution found")
}

func GraphSearch(problem Problem) (*Solution, error) {
	frontier := []Node{problem.InitialState}
	actionsTaken := []string{}
	exploredStates := map[string]bool{}
	for len(frontier) > 0 {
		currentNode, _ := PopFifo(&frontier)
		if problem.GoalTest(currentNode.State) {
			return &Solution{Actions: actionsTaken}, nil
		}

		if exploredStates[currentNode.State] {
			continue
		}

		exploredStates[currentNode.State] = true
		for _, action := range problem.Actions(currentNode.State) {
			childNode := ChildNode(problem, currentNode, action)
			frontier = append(frontier, childNode)
			actionsTaken = append(actionsTaken, action)
		}
	}

	return nil, fmt.Errorf("No solution found")
}

func BreadthFirstSearch(problem Problem) (*Solution, error) {
	node := problem.InitialState
	if problem.GoalTest(node.State) {
		return &Solution{}, nil
	}
	frontier := []Node{node}
	exploredStates := map[string]bool{}
	for len(frontier) > 0 {
		currentNode, _ := PopFifo(&frontier)
		exploredStates[currentNode.State] = true
		for _, action := range problem.Actions(currentNode.State) {
			childNode := ChildNode(problem, currentNode, action)
			if !exploredStates[childNode.State] && !IsStateInList(childNode.State, frontier) {
				if problem.GoalTest(childNode.State) {
					return SolutionPath(childNode), nil
				}
				frontier = append(frontier, childNode)
			}
		}
	}

	return nil, fmt.Errorf("No solution found")
}
