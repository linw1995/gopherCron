package scheduling

import (
	"encoding/json"
	"io"
)

func FlowFromJSONData(data []byte) (f Flow, err error) {
	err = json.Unmarshal(data, &f)
	return
}

func FlowFromJSON(rd io.Reader) (f Flow, err error) {
	err = json.NewDecoder(rd).Decode(&f)
	return
}

func FlowToJSONData(f Flow) (data []byte, err error) {
	return json.Marshal(f)
}

func FlowToJSON(wd io.Writer, f Flow) (err error) {
	return json.NewEncoder(wd).Encode(f)
}

func StateFromJSONData(data []byte) (s State, err error) {
	err = json.Unmarshal(data, &s)
	return
}

func StateFromJSON(rd io.Reader) (s State, err error) {
	err = json.NewDecoder(rd).Decode(&s)
	return
}

func StateToJSONData(s State) (data []byte, err error) {
	return json.Marshal(s)
}

func StateToJSON(wd io.Writer, s State) (err error) {
	return json.NewEncoder(wd).Encode(s)
}
