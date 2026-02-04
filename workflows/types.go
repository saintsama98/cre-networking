package workflows

// ReceiverInput represents the input payload for the receiver workflow
// This is shared between sender and receiver workflows
type ReceiverInput struct {
	Message   string                 `json:"message"`
	SenderID  string                 `json:"sender_id,omitempty"`
	Timestamp int64                  `json:"timestamp,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}
