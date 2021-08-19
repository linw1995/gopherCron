package scheduling

import (
	"reflect"
	"testing"
)

func TestKahn(t *testing.T) {
	t.Run("Common", func(t *testing.T) {
		got, err := Kahn(TasksToGraph(FlowExample.Tasks))
		if err != nil {
			t.Fatal(err)
		}
		if want := []string{"C", "E"}; !reflect.DeepEqual(got[3:], want) {
			t.Errorf("want %v got %v", want, got)
		}
		if want := []string{"A", "B", "D"}; !SameStringSet(got[:3], want) {
			t.Errorf("want %v got %v", want, got)
		}
	})

	t.Run("Bad", func(t *testing.T) {
		got, err := Kahn(TasksToGraph(BadFlowExample.Tasks))
		if err == nil {
			t.Fatal("Should fail, but got: ", got)
		}
	})
}

func TestKahnPlus(t *testing.T) {

	t.Run("Common", func(t *testing.T) {
		got, err := KahnPlus(TasksToGraph(FlowExample.Tasks), DummyQueueFactory)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(got)
	})

	t.Run("Bad", func(t *testing.T) {
		got, err := KahnPlus(TasksToGraph(BadFlowExample.Tasks), DummyQueueFactory)
		if err == nil {
			t.Fatal("Should fail, but got: ", got)
		}
	})

	flow := FlowExample

	t.Run("PreferInOrder", func(t *testing.T) {
		got, err := KahnPlus(TasksToGraph(flow.Tasks), PreferInOrder(flow.Tasks))
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"A", "B", "C", "D", "E"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v got %v", want, got)
		}
	})

	t.Run("PreferLongestPathViaDFS", func(t *testing.T) {
		got, err := KahnPlus(TasksToGraph(flow.Tasks), PreferLongestPath(flow.Tasks))
		if err != nil {
			t.Fatal(err)
		}
		if want := []string{"D", "C", "E"}; !reflect.DeepEqual(got[2:], want) {
			t.Errorf("want %v got %v", want, got)
		}
		if want := []string{"A", "B"}; !SameStringSet(got[:2], want) {
			t.Errorf("want %v got %v", want, got)
		}
	})

	t.Run("CombineTwoTypeOfWeight", func(t *testing.T) {
		got, err := KahnPlus(TasksToGraph(flow.Tasks), Prefer(flow.Tasks, LongestPathWeight, InOrderWeight))
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"A", "B", "C", "D", "E"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v got %v", want, got)
		}

	})
}
