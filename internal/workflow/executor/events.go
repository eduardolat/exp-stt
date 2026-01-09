package executor

import "time"

// ExecutionStatus represents the status of a node during execution.
type ExecutionStatus string

const (
	StatusPending   ExecutionStatus = "pending"
	StatusRunning   ExecutionStatus = "running"
	StatusCompleted ExecutionStatus = "completed"
	StatusError     ExecutionStatus = "error"
	StatusSkipped   ExecutionStatus = "skipped"
)

// ExecutionEvent is emitted during workflow execution for observability.
// These events can be streamed to the frontend via SSE for real-time updates.
type ExecutionEvent struct {
	ExecutionID   string          `json:"execution_id"`
	NodeID        string          `json:"node_id"`
	Status        ExecutionStatus `json:"status"`
	OutputPreview string          `json:"output_preview,omitempty"`
	ErrorMessage  string          `json:"error_message,omitempty"`
	Timestamp     time.Time       `json:"timestamp"`
}

// NewExecutionEvent creates a new execution event.
func NewExecutionEvent(executionID, nodeID string, status ExecutionStatus) ExecutionEvent {
	return ExecutionEvent{
		ExecutionID: executionID,
		NodeID:      nodeID,
		Status:      status,
		Timestamp:   time.Now(),
	}
}

// WithOutputPreview adds an output preview to the event.
func (e ExecutionEvent) WithOutputPreview(preview string) ExecutionEvent {
	e.OutputPreview = preview
	return e
}

// WithError adds an error message to the event.
func (e ExecutionEvent) WithError(err error) ExecutionEvent {
	if err != nil {
		e.ErrorMessage = err.Error()
	}
	return e
}
