package scheduling

import (
	"fmt"
)

func TasksToGraph(tasks []Task) (rv map[string][]string) {
	rv = make(map[string][]string)
	for _, task := range tasks {
		rv[task.ID] = task.Deps
	}
	return InverseGraph(rv)
}

func InverseGraph(graph map[string][]string) (igraph map[string][]string) {
	igraph = make(map[string][]string)
	for node, outcomes := range graph {
		for _, outcome := range outcomes {
			igraph[outcome] = append(igraph[outcome], node)
		}
		if _, existed := igraph[node]; !existed {
			igraph[node] = make([]string, 0)
		}
	}
	return igraph
}

func Kahn(graph map[string][]string) (L []string, err error) {
	var S []string
	igraph := InverseGraph(graph)
	inCountMap := make(map[string]int)
	for node, incomes := range igraph {
		inCountMap[node] = len(incomes)
		if inCountMap[node] == 0 {
			S = append(S, node)
		}
	}

	for {
		if len(S) == 0 {
			break
		}
		node := S[0]
		S = S[1:]
		L = append(L, node)

		for _, outcome := range graph[node] {
			inCountMap[outcome]--
			if inCountMap[outcome] == 0 {
				S = append(S, outcome)
			}
		}
	}

	for node, inCount := range inCountMap {
		if inCount != 0 {
			err = fmt.Errorf("Invalid DAG node %v", node)
			return
		}
	}

	return

}

func KahnPlus(graph map[string][]string, factory QueueFactory) (selected []string, err error) {
	queue := factory()
	igraph := InverseGraph(graph)
	inCountMap := make(map[string]int)
	for node, incomes := range igraph {
		inCountMap[node] = len(incomes)
		if inCountMap[node] == 0 {
			queue.Put(node)
		}
	}

	for {
		node, ok := queue.Get()
		if !ok {
			break
		}
		selected = append(selected, node)

		for _, outcome := range graph[node] {
			inCountMap[outcome]--
			if inCountMap[outcome] == 0 {
				queue.Put(outcome)
			}
		}
	}

	for node, inCount := range inCountMap {
		if inCount != 0 {
			err = fmt.Errorf("Invalid DAG node %v", node)
			return
		}
	}

	return
}
