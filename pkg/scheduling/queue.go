package scheduling

import (
	"container/heap"
)

func PreferInOrder(tasks []Task) QueueFactory {
	return Prefer(tasks, InOrderWeight)
}

func PreferLongestPath(tasks []Task) QueueFactory {
	return Prefer(tasks, LongestPathWeight)
}

func Prefer(tasks []Task, weightFactories ...WeightFactory) QueueFactory {
	return func() Queue {
		lenTasks := len(tasks)
		var (
			weight map[string]int
		)
		switch len(weightFactories) {
		case 0:
			return &DummyQueue{}
		case 1:
			weight = weightFactories[0](tasks)
		default:
			weight = make(map[string]int)
			for _, weightFactory := range weightFactories {
				newWeight := weightFactory(tasks)
				for _, task := range tasks {
					weight[task.ID] = weight[task.ID]*lenTasks + newWeight[task.ID]
				}
			}
		}
		return &SortByWeightQueue{&SortByWeight{Weight: weight}}
	}
}

type WeightFactory func([]Task) map[string]int

func LongestPathWeight(tasks []Task) map[string]int {
	graph := TasksToGraph(tasks)
	weight := LongestPathViaDFS(graph)
	return weight
}

func InOrderWeight(tasks []Task) map[string]int {
	rv := make(map[string]int)
	lenTasks := len(tasks)
	for no, t := range tasks {
		rv[t.ID] = lenTasks - no
	}
	return rv
}

type QueueFactory func() Queue

var (
	_ Queue = (*DummyQueue)(nil)
	_ Queue = (*SortByWeightQueue)(nil)
)

type Queue interface {
	Put(string)
	Get() (string, bool)
	Len() int
}

func DummyQueueFactory() Queue {
	return &DummyQueue{}
}

type DummyQueue struct {
	nodes []string
}

func (q *DummyQueue) Put(node string) {
	q.nodes = append(q.nodes, node)
}

func (q *DummyQueue) Get() (string, bool) {
	if lenNodes := len(q.nodes); lenNodes > 0 {
		rv := q.nodes[lenNodes-1]
		q.nodes = q.nodes[:lenNodes-1]
		return rv, true
	}
	return "", false
}

func (q *DummyQueue) Len() int {
	return len(q.nodes)
}

type SortByWeightQueue struct {
	heap.Interface
}

func (s *SortByWeightQueue) Put(node string) {
	heap.Push(s, node)
}

func (s *SortByWeightQueue) Get() (string, bool) {
	if s.Len() <= 0 {
		return "", false
	}
	node := heap.Pop(s).(string)
	return node, true
}

var (
	_ heap.Interface = (*SortByWeight)(nil)
)

type SortByWeight struct {
	Weight map[string]int
	Seq    []string
}

func (s *SortByWeight) Len() int {
	return len(s.Seq)
}

func (s *SortByWeight) Less(i, j int) bool {
	// reverse order to make value larger first
	return s.Weight[s.Seq[i]] > s.Weight[s.Seq[j]]
}

func (s *SortByWeight) Swap(i, j int) {
	s.Seq[i], s.Seq[j] = s.Seq[j], s.Seq[i]
}

func (s *SortByWeight) Push(v interface{}) {
	s.Seq = append(s.Seq, v.(string))
}

func (s *SortByWeight) Pop() interface{} {
	l := s.Len() - 1
	rv := s.Seq[l]
	s.Seq = s.Seq[0:l]
	return rv
}
