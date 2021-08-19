package scheduling

import (
	"reflect"
	"testing"
)

func runTwice(do func()) {
	for i := 0; i < 2; i++ {
		do()
	}
}

func TestScheduleFail(t *testing.T) {
	_, err := New(BadFlowExample.Tasks, LongestPathWeight, InOrderWeight)
	if err == nil {
		t.Fatalf("Schedule on bad tasks flow, it should fail")
	}
}

func TestScheduleRunAndFinish(t *testing.T) {
	state, err := New(FlowExample.Tasks, LongestPathWeight, InOrderWeight)
	if err != nil {
		t.Fatal(err)
	}

	var (
		readyIDs []string
		wants    = [][]string{
			{"A", "B", "D"},
			{"C"},
			{"E"},
			nil,
		}
	)

	for _, want := range wants {
		runTwice(func() {
			readyIDs = state.Schedule()
			if !reflect.DeepEqual(readyIDs, want) {
				t.Errorf("got %v want %v", readyIDs, want)
			}
		})
		if readyIDs != nil {
			if err := state.Run(readyIDs...); err != nil {
				t.Errorf("state.Run on %v err %s", readyIDs, err)
			}
			if err := state.Finish(readyIDs...); err != nil {
				t.Errorf("state.Finish on %v err %s", readyIDs, err)
			}
		}
	}
}

func TestSchedulePartial(t *testing.T) {
	for _, param := range []struct {
		Name            string
		WeightFactories []WeightFactory
		Pick            []string
		Wants           [][]string
	}{
		{
			Name:            "Prefer_LongestPath_InOrder",
			WeightFactories: []WeightFactory{LongestPathWeight, InOrderWeight},
			Pick: []string{
				"A",
				"B",
				"C",
				"D",
				"E",
			},
			Wants: [][]string{
				{"A", "B", "D"},
				{"B", "D"},
				{"C", "D"},
				{"D"},
				{"E"},
				nil,
			},
		},
		{
			Name:            "Prefer_LongestPath",
			WeightFactories: []WeightFactory{LongestPathWeight},
			Pick: []string{
				"A",
				"B",
				"D",
				"C",
				"E",
			},
			Wants: [][]string{
				{"A", "B", "D"},
				{"B", "D"},
				{"C", "D"},
				{"C"},
				{"E"},
				nil,
			},
		},
		{
			Name:            "Prefer_nothing",
			WeightFactories: nil,
			Pick: []string{
				"A",
				"B",
				"C",
				"D",
				"E",
			},
			Wants: [][]string{
				{"A", "B", "D"},
				{"B", "D"},
				{"C", "D"},
				{"D"},
				{"E"},
				nil,
			},
		},
	} {
		t.Run(param.Name, func(t *testing.T) {
			state, err := New(FlowExample.Tasks, param.WeightFactories...)
			if err != nil {
				t.Fatal(err)
			}
			var (
				readyIDs []string
			)
			for idx, want := range param.Wants {
				runTwice(func() {
					readyIDs = state.Schedule()
					if !SameStringSet(readyIDs, want) {
						t.Errorf("got %v want %v", readyIDs, want)
					}
				})
				if readyIDs != nil {
					picked := param.Pick[idx]
					if err := state.Run(picked); err != nil {
						t.Errorf("state.Run on %v err %s", picked, err)
					}
					if err := state.Finish(picked); err != nil {
						t.Errorf("state.Finish on %v err %s", picked, err)
					}
				}
			}
		})
	}
}

func TestScheduleRunFail(t *testing.T) {
	state, err := New(FlowExample.Tasks, LongestPathWeight, InOrderWeight)
	if err != nil {
		t.Fatal(err)
	}

	if err := state.Run("A"); err != nil {
		t.Errorf("state.Run should succeed")
	}

	if err := state.Run("A"); err == nil {
		t.Errorf("state.Run on same id twice should fail")
	}

	state.TaskStates["A"] = DONE
	if err := state.Run("A"); err == nil {
		t.Errorf("state.Run on done id should fail")
	}

	if err := state.Run("NotFound"); err == nil {
		t.Errorf("state.Run on not found id should fail")
	}
}

func TestScheduleFinishFail(t *testing.T) {
	state, err := New(FlowExample.Tasks, LongestPathWeight, InOrderWeight)
	if err != nil {
		t.Fatal(err)
	}

	if err := state.Finish("A"); err != nil {
		t.Errorf("state.Finish should succeed")
	}

	state.TaskStates["A"] = RUNNING
	if err := state.Finish("A"); err != nil {
		t.Errorf("state.Finish should succeed")
	}

	state.TaskStates["A"] = DONE
	if err := state.Finish("A"); err == nil {
		t.Errorf("state.Finish on same id twice should fail")
	}

	if err := state.Finish("NotFound"); err == nil {
		t.Errorf("state.Finish on not found id should fail")
	}
}
