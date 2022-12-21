package types

// Counter metric
type Counter struct {
	Name  string       `json:"name"`
	Value CounterValue `json:"value"`
	Error string       `json:"error,omitempty"`
}

type CounterValue interface {
	MarshalJSON() ([]byte, error)
}
