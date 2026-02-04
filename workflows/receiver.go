package workflows

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// ReceiverConfig represents the configuration for the receiver workflow
type ReceiverConfig struct {
	// Configuration can be extended as needed
}

// ReceiverOutput represents the output from the receiver workflow
type ReceiverOutput struct {
	Status     string                 `json:"status"`
	ReceivedAt int64                  `json:"received_at"`
	Processed  bool                   `json:"processed"`
	Response   string                 `json:"response"`
	EchoedData map[string]interface{} `json:"echoed_data,omitempty"`
}

// executeReceiverWorkflow executes the receiver workflow logic
func executeReceiverWorkflow(config ReceiverConfig, input ReceiverInput, runtime cre.Runtime) (ReceiverOutput, error) {
	// Log the received input
	runtime.Logger().Info("Receiver workflow received message", "message", input.Message)
	if input.SenderID != "" {
		runtime.Logger().Info("Message from sender", "sender_id", input.SenderID)
	}

	// Process the input data
	processed := true
	response := fmt.Sprintf("Successfully received and processed: %s", input.Message)

	// Echo back the data with additional metadata
	echoedData := make(map[string]interface{})
	if input.Data != nil {
		for k, v := range input.Data {
			echoedData[k] = v
		}
	}
	echoedData["processed_by"] = "receiver_workflow"
	echoedData["original_message"] = input.Message
	echoedData["received_timestamp"] = runtime.Now().Unix()

	output := ReceiverOutput{
		Status:     "success",
		ReceivedAt: runtime.Now().Unix(),
		Processed:  processed,
		Response:   response,
		EchoedData: echoedData,
	}

	runtime.Logger().Info("Receiver workflow completed", "status", output.Status)
	return output, nil
}

// ReceiverWorkflowHandler handles HTTP trigger requests for the receiver workflow
func ReceiverWorkflowHandler(config ReceiverConfig, input ReceiverInput, runtime cre.Runtime) (ReceiverOutput, error) {
	return executeReceiverWorkflow(config, input, runtime)
}

// NewReceiverWorkflow creates a new receiver workflow using CRE SDK patterns
func NewReceiverWorkflow() cre.Workflow[ReceiverConfig] {
	// Create HTTP trigger for the receiver workflow
	httpTrigger := http.Trigger(&http.Config{
		AuthorizedKeys: []*http.AuthorizedKey{
			// Add authorized keys in production
			// Example:
			// {
			// 	Type:      http.KeyType_KEY_TYPE_ECDSA_EVM,
			// 	PublicKey: "0xYourEVMAddress",
			// },
		},
	})

	// Define the workflow execution handler
	handler := cre.Handler(
		httpTrigger,
		func(config ReceiverConfig, runtime cre.Runtime, payload *http.Payload) (ReceiverOutput, error) {
			// Parse the input from the HTTP trigger payload
			var input ReceiverInput
			if err := json.Unmarshal(payload.Input, &input); err != nil {
				// If unmarshaling fails, try to extract message from raw body
				input = ReceiverInput{
					Message: string(payload.Input),
					Data:    make(map[string]interface{}),
				}
				runtime.Logger().Warn("Failed to parse JSON input, using raw body", "error", err)
			}

			// Set timestamp if not provided
			if input.Timestamp == 0 {
				input.Timestamp = runtime.Now().Unix()
			}

			return ReceiverWorkflowHandler(config, input, runtime)
		},
	)

	return cre.Workflow[ReceiverConfig]{handler}
}

// Helper function to create receiver input from JSON
func ParseReceiverInput(jsonData []byte) (ReceiverInput, error) {
	var input ReceiverInput
	if err := json.Unmarshal(jsonData, &input); err != nil {
		return ReceiverInput{}, fmt.Errorf("failed to parse receiver input: %w", err)
	}
	return input, nil
}
