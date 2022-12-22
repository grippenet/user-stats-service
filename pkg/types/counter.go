package types

// Counter metric
type Counter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
	Error string      `json:"error,omitempty"`
}

const COUNTER_TYPE_COUNT = "count"
const COUNTER_TYPE_MAP = "map"
