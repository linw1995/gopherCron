package scheduling

import (
	"reflect"
	"sort"
)

func dfs(graph map[string][]string, node string, longestPathMap map[string]int) (depth int) {
	longestPathMap[node] = 0
	for _, outcome := range graph[node] {
		if _, existed := longestPathMap[outcome]; !existed {
			dfs(graph, outcome, longestPathMap)
		}
		if outcomeDepth := longestPathMap[outcome]; outcomeDepth > depth {
			depth = outcomeDepth
		}
	}
	depth++
	longestPathMap[node] = depth
	return
}

// LongestPathViaDFS produce longest path length of nodes in DAG.
func LongestPathViaDFS(graph map[string][]string) (longestPathMap map[string]int) {
	longestPathMap = make(map[string]int)
	for node := range graph {
		if _, existed := longestPathMap[node]; existed {
			continue
		}
		dfs(graph, node, longestPathMap)
	}
	return longestPathMap
}

func SameStringSet(a []string, b []string) bool {
	copy := func(src []string) []string {
		dst := make([]string, len(src))
		copy(dst, src)
		return dst
	}

	a = copy(a)
	sort.Strings(a)
	b = copy(b)
	sort.Strings(b)
	return reflect.DeepEqual(a, b)
}
