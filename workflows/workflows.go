package workflows

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// WorkflowRegistry manages all registered workflows
type WorkflowRegistry struct {
	receiverWorkflow cre.Workflow[ReceiverConfig]
	senderWorkflow   cre.Workflow[SenderConfig]
}

// NewWorkflowRegistry creates a new workflow registry with all workflows initialized
func NewWorkflowRegistry() *WorkflowRegistry {
	return &WorkflowRegistry{
		receiverWorkflow: NewReceiverWorkflow(),
		senderWorkflow:   NewSenderWorkflow(),
	}
}

// GetReceiverWorkflow returns the receiver workflow
func (r *WorkflowRegistry) GetReceiverWorkflow() cre.Workflow[ReceiverConfig] {
	return r.receiverWorkflow
}

// GetSenderWorkflow returns the sender workflow
func (r *WorkflowRegistry) GetSenderWorkflow() cre.Workflow[SenderConfig] {
	return r.senderWorkflow
}

// RegisterAllWorkflows registers all workflows in the registry
// In CRE, workflows are registered via cre.yaml configuration
// This function is a helper for programmatic access
func (r *WorkflowRegistry) RegisterAllWorkflows() error {
	// Workflows are automatically registered when deployed via cre.yaml
	// This function can be used for validation or testing
	if len(r.receiverWorkflow) == 0 {
		return fmt.Errorf("receiver workflow not initialized")
	}
	if len(r.senderWorkflow) == 0 {
		return fmt.Errorf("sender workflow not initialized")
	}
	return nil
}

// LogWorkflowCommunication logs workflow communication events
func LogWorkflowCommunication(runtime cre.Runtime, event string, details map[string]interface{}) {
	runtime.Logger().Info("Workflow Communication", "event", event)
	for k, v := range details {
		runtime.Logger().Info("Communication detail", "key", k, "value", v)
	}
}

// ValidateWorkflowInput validates workflow input data
func ValidateWorkflowInput(input interface{}) error {
	// Basic validation - can be extended
	data, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("input cannot be empty")
	}

	return nil
}

// CreateTestReceiverInput creates a test receiver input for testing
func CreateTestReceiverInput(message string, senderID string, data map[string]interface{}) ReceiverInput {
	return ReceiverInput{
		Message:   message,
		SenderID:  senderID,
		Timestamp: 0, // Will be set by workflow
		Data:      data,
	}
}

// CreateTestSenderInput creates a test sender input for testing
func CreateTestSenderInput(message string, targetURL string, data map[string]interface{}) SenderInput {
	return SenderInput{
		Message:    message,
		TargetURL:  targetURL,
		Data:       data,
		RetryCount: 0,
	}
}
