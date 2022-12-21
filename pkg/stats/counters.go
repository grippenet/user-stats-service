package stats

import (
	"encoding/json"
)

type SimpleCounter struct {
	Count int64
}

func (c SimpleCounter) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{"type": "count", "value": c.Count}
	return json.Marshal(m)
}

type MapCounter struct {
	Counts map[string]int64
}

func (c MapCounter) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{"type": "map", "value": c.Counts}
	return json.Marshal(m)
}
