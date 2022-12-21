package types

// Counter metric
type Counter struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
	Error string `json:"error,omitempty"`
}
