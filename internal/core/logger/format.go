package logger

type OperationLog struct {
	ID           string `json:"id"`
	Operation    string `json:"ops"`
	ResponseTime string `json:"rt"`
	IngoingData  any    `json:"in,omitempty"`
	OutgoingData any    `json:"out,omitempty"`
	MetaData     any    `json:"meta,omitempty"`
	Error        any    `json:"error,omitempty"`
}
