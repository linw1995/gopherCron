package scheduling

import (
	"fmt"
)

func New(tasks []Task, weightFactories ...WeightFactory) (s State, err error) {
	graph := TasksToGraph(tasks)
	TopologicalSortingResult, err := KahnPlus(graph, Prefer(tasks, weightFactories...))
	if err != nil {
		return
	}
	states := make(map[string]TaskState)
	for _, task := range tasks {
		states[task.ID] = SCHEDULING
	}
	s.Tasks = tasks
	s.TaskStates = states
	s.IGraph = InverseGraph(graph)
	s.SchedulingOrder = TopologicalSortingResult
	return
}

type State struct {
	Tasks           []Task               `yaml:"tasks" toml:"tasks" json:"tasks"`
	TaskStates      map[string]TaskState `yaml:"task_states" toml:"task_states" json:"task_states"`
	SchedulingOrder []string             `yaml:"order" toml:"order" json:"order"` // Topological Sorting Result
	IGraph          map[string][]string  `yaml:"i_graph" toml:"i_graph" json:"i_graph"`
}

type TaskState string

const (
	SCHEDULING TaskState = "scheduling"
	RUNNING    TaskState = "running"
	DONE       TaskState = "done"
)

func (s *State) Schedule() (readyIDs []string) {
	for _, taskID := range s.SchedulingOrder {
		if s.TaskStates[taskID] != SCHEDULING {
			continue
		}
		ready := true
		for _, income := range s.IGraph[taskID] {
			if s.TaskStates[income] != DONE {
				ready = false
				break
			}
		}
		if ready {
			readyIDs = append(readyIDs, taskID)
		}
	}
	return
}

// Run flags states of tasks of given IDs Running
func (s *State) Run(taskIDs ...string) error {
	for _, taskID := range taskIDs {
		if status, ok := s.TaskStates[taskID]; !ok {
			return fmt.Errorf("%v is not found", taskID)
		} else if status != SCHEDULING {
			return fmt.Errorf("can not run task %v in %v", taskID, status)
		}
		s.TaskStates[taskID] = RUNNING
	}
	return nil
}

// Finish flags states of tasks of given IDs DONE
func (s *State) Finish(taskIDs ...string) error {
	for _, taskID := range taskIDs {
		if status, ok := s.TaskStates[taskID]; !ok {
			return fmt.Errorf("%v is not found", taskID)
		} else if status == DONE {
			return fmt.Errorf("can not finish task %v twice", taskID)
		}
		s.TaskStates[taskID] = DONE
	}
	return nil
}
