package workflows

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"google.golang.org/protobuf/types/known/durationpb"
)

// SenderConfig represents the configuration for the sender workflow
type SenderConfig struct {
	ReceiverWorkflowURL string `json:"receiver_workflow_url,omitempty"`
}

// SenderInput represents the input for the sender workflow
type SenderInput struct {
	Message    string                 `json:"message"`
	TargetURL  string                 `json:"target_url,omitempty"` // Override default receiver URL
	Data       map[string]interface{} `json:"data,omitempty"`
	RetryCount int                    `json:"retry_count,omitempty"`
}

// SenderOutput represents the output from the sender workflow
type SenderOutput struct {
	Status         string                 `json:"status"`
	SentAt         int64                  `json:"sent_at"`
	ResponseStatus string                 `json:"response_status,omitempty"`
	ResponseData   map[string]interface{} `json:"response_data,omitempty"`
	Error          string                 `json:"error,omitempty"`
	RequestID      string                 `json:"request_id,omitempty"`
}

// executeSenderWorkflow executes the sender workflow logic
func executeSenderWorkflow(config SenderConfig, input SenderInput, runtime cre.Runtime, httpClient *http.Client) (SenderOutput, error) {
	runtime.Logger().Info("Sender workflow initiating request", "message", input.Message)

	// Determine target URL
	targetURL := config.ReceiverWorkflowURL
	if input.TargetURL != "" {
		targetURL = input.TargetURL
	}

	if targetURL == "" {
		return SenderOutput{
			Status: "error",
			SentAt: runtime.Now().Unix(),
			Error:  "target URL not specified",
		}, fmt.Errorf("target URL not specified")
	}

	// Prepare the request payload
	requestPayload := ReceiverInput{
		Message:   input.Message,
		SenderID:  "sender_workflow",
		Timestamp: runtime.Now().Unix(),
		Data:      input.Data,
	}

	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return SenderOutput{
			Status: "error",
			SentAt: runtime.Now().Unix(),
			Error:  fmt.Sprintf("failed to marshal payload: %v", err),
		}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Configure HTTP request
	requestConfig := &http.Request{
		Method: "POST",
		Url:    targetURL,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: payloadBytes,
		// Cache settings ensure single execution across DON nodes
		CacheSettings: &http.CacheSettings{
			Store:  true,
			MaxAge: durationpb.New(5 * time.Minute),
		},
	}

	// Send the HTTP request using CRE SDK
	runtime.Logger().Info("Sending POST request", "url", targetURL)

	// Use SendRequest with proper callback pattern
	responsePromise := http.SendRequest(
		requestConfig,
		runtime,
		httpClient,
		func(req *http.Request, logger *slog.Logger, sendRequester *http.SendRequester) (*http.Response, error) {
			promise := sendRequester.SendRequest(req)
			return promise.Await()
		},
		cre.ConsensusIdenticalAggregation[*http.Response](),
	)

	response, err := responsePromise.Await()
	if err != nil {
		runtime.Logger().Error("Request failed", "error", err)
		return SenderOutput{
			Status: "error",
			SentAt: runtime.Now().Unix(),
			Error:  fmt.Sprintf("HTTP request failed: %v", err),
		}, fmt.Errorf("HTTP request failed: %w", err)
	}

	// Parse response
	var responseData map[string]interface{}
	if len(response.Body) > 0 {
		if err := json.Unmarshal(response.Body, &responseData); err != nil {
			runtime.Logger().Warn("Failed to parse response body", "error", err)
			// Continue even if parsing fails
			responseData = map[string]interface{}{
				"raw_body": string(response.Body),
			}
		}
	}

	runtime.Logger().Info("Request completed", "status_code", response.StatusCode)

	output := SenderOutput{
		Status:         "success",
		SentAt:         runtime.Now().Unix(),
		ResponseStatus: fmt.Sprintf("%d", response.StatusCode),
		ResponseData:   responseData,
		RequestID:      fmt.Sprintf("req_%d", runtime.Now().Unix()),
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		runtime.Logger().Info("Successfully communicated with receiver workflow")
	} else {
		output.Status = "warning"
		output.Error = fmt.Sprintf("Received non-2xx status code: %d", response.StatusCode)
	}

	return output, nil
}

// SenderWorkflowHandler handles execution of the sender workflow
func SenderWorkflowHandler(config SenderConfig, input SenderInput, runtime cre.Runtime) (SenderOutput, error) {
	// Create HTTP client
	httpClient := &http.Client{}

	return executeSenderWorkflow(config, input, runtime, httpClient)
}

// NewSenderWorkflow creates a new sender workflow using CRE SDK patterns
func NewSenderWorkflow() cre.Workflow[SenderConfig] {
	// Create HTTP trigger for the sender workflow
	httpTrigger := http.Trigger(&http.Config{
		AuthorizedKeys: []*http.AuthorizedKey{
			// Add authorized keys in production
		},
	})

	// Define the workflow execution handler
	handler := cre.Handler(
		httpTrigger,
		func(config SenderConfig, runtime cre.Runtime, payload *http.Payload) (SenderOutput, error) {
			// Parse the input from the HTTP trigger payload
			var input SenderInput
			if err := json.Unmarshal(payload.Input, &input); err != nil {
				return SenderOutput{
					Status: "error",
					Error:  fmt.Sprintf("failed to parse input: %v", err),
				}, fmt.Errorf("failed to parse input: %w", err)
			}

			return SenderWorkflowHandler(config, input, runtime)
		},
	)

	return cre.Workflow[SenderConfig]{handler}
}
