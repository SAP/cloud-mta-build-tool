package mta

import (
	"fmt"

	"github.com/deckarep/golang-set"
	"github.com/pkg/errors"
)

type graphNode struct {
	module string
	deps   mapset.Set
	index  int
}

func new(module *string, deps mapset.Set, index int) *graphNode {
	return &graphNode{module: *module, deps: deps, index: index}
}

type Graph map[string]*graphNode

func (mta MTA) GetModulesOrder() ([]string, error) {
	var graph Graph = make(map[string]*graphNode)
	for index, module := range mta.Modules {
		deps := mapset.NewSet()
		if module.BuildParams.Requires != nil {
			for _, req := range module.BuildParams.Requires {
				deps.Add(req.Name)
			}
		}
		graph[module.Name] = new(&module.Name, deps, index)
	}
	return resolveGraph(graph, mta)
}

// Resolves the dependency Graph
func resolveGraph(graph Graph, mta MTA) ([]string, error) {
	overleft := graph

	// Iteratively find and remove nodes from the Graph which have no dependencies.
	// If at some point there are still nodes in the Graph and we cannot find
	// nodes without dependencies, that means we have a circular dependency
	var resolved []string
	for len(overleft) != 0 {
		// Get all nodes from the Graph which have no dependencies
		readyNodesSet := mapset.NewSet()
		readyModulesSet := mapset.NewSet()
		for _, node := range overleft {
			if node.deps.Cardinality() == 0 {
				readyNodesSet.Add(node)
				readyModulesSet.Add(node.module)
			}
		}

		// If there aren't any ready nodes, then we have a circular dependency
		if readyNodesSet.Cardinality() == 0 {
			module1 := ""
			module2 := ""
			index := 0
			for _, node := range overleft {
				if index == 0 {
					module1 = node.module
					index = 1
				} else {
					module2 = node.module
					break
				}
			}

			return nil, errors.New(fmt.Sprintf("Circular dependency found. Check modules %v and %v", module1, module2))
		}

		// Remove the ready nodes and add them to the resolved Graph
		readyModulesIndexes := mapset.NewSet()
		for node := range readyNodesSet.Iter() {
			delete(overleft, node.(*graphNode).module)
			readyModulesIndexes.Add(node.(*graphNode).index)
		}

		for index, module := range mta.Modules {
			if readyModulesIndexes.Contains(index) {
				resolved = append(resolved, module.Name)
			}
		}

		// remove the ready nodes from the remaining node dependencies as well
		for _, node := range overleft {
			node.deps = node.deps.Difference(readyModulesSet)
		}
	}

	return resolved, nil
}
