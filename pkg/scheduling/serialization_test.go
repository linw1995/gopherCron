package scheduling

import (
	"bytes"
	"reflect"
	"testing"
)

var (
	FlowExample    Flow
	BadFlowExample Flow
)

func init() {
	var err error
	FlowExample, err = FlowFromJSONData([]byte(`
	{
		"tasks": [
			{"id": "A"},
			{"id": "B"},
			{"id": "C", "deps": ["A", "B"]},
			{"id": "D"},
			{"id": "E", "deps": ["C", "D"]}
		]
	}
	`))
	if err != nil {
		panic(err)
	}
	BadFlowExample, err = FlowFromJSONData([]byte(`
	{
		"tasks": [
			{"id": "A", "deps": ["C"]},
			{"id": "B", "deps": ["A"]},
			{"id": "C", "deps": ["B"]}
		]
	}
	`))
	if err != nil {
		panic(err)
	}
}

func TestFlow(t *testing.T) {
	data, err := FlowToJSONData(FlowExample)
	if err != nil {
		t.Fatal(err)
	}
	nf, err := FlowFromJSONData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(nf, FlowExample) {
		t.Errorf("deserialize fail")
	}

	wd := bytes.NewBuffer(nil)
	err = FlowToJSON(wd, FlowExample)
	if err != nil {
		t.Fatal(err)
	}
	nf, err = FlowFromJSON(wd)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(nf, FlowExample) {
		t.Errorf("deserialize fail")
	}
}

func TestState(t *testing.T) {
	s, err := New(FlowExample.Tasks, LongestPathWeight, InOrderWeight)
	if err != nil {
		t.Fatal(err)
	}

	data, err := StateToJSONData(s)
	if err != nil {
		t.Fatal(err)
	}
	ns, err := StateFromJSONData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ns, s) {
		t.Errorf("deserialize fail")
	}

	wd := bytes.NewBuffer(nil)
	err = StateToJSON(wd, s)
	if err != nil {
		t.Fatal(err)
	}
	ns, err = StateFromJSON(wd)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ns, s) {
		t.Errorf("deserialize fail")
	}
}
