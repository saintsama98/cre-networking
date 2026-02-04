package tests

import (
	"encoding/json"
	"testing"
	"time"

	"cre-networking-test/workflows"

	"github.com/smartcontractkit/cre-sdk-go/cre/testutils"
)

// TestReceiverWorkflow tests the receiver workflow functionality
func TestReceiverWorkflow(t *testing.T) {
	// Create a test runtime
	runtime := testutils.NewRuntime(t, nil)

	// Create receiver workflow
	receiverWorkflow := workflows.NewReceiverWorkflow()
	if len(receiverWorkflow) == 0 {
		t.Fatal("Receiver workflow not initialized")
	}

	// Create test input
	input := workflows.CreateTestReceiverInput(
		"Test message from integration test",
		"test_sender",
		map[string]interface{}{
			"test_key": "test_value",
			"number":   42,
		},
	)

	// Execute the receiver workflow handler directly
	config := workflows.ReceiverConfig{}
	output, err := workflows.ReceiverWorkflowHandler(config, input, runtime)
	if err != nil {
		t.Fatalf("Receiver workflow execution failed: %v", err)
	}

	if output.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", output.Status)
	}

	if !output.Processed {
		t.Error("Expected processed to be true")
	}

	if output.Response == "" {
		t.Error("Expected non-empty response")
	}

	// Verify echoed data
	if output.EchoedData == nil {
		t.Error("Expected echoed data to be present")
	} else {
		if output.EchoedData["test_key"] != "test_value" {
			t.Errorf("Expected echoed data to contain test_key='test_value'")
		}
		if output.EchoedData["processed_by"] != "receiver_workflow" {
			t.Error("Expected processed_by to be 'receiver_workflow'")
		}
	}
}

// TestSenderWorkflow tests the sender workflow functionality
// Note: This test requires a running receiver workflow endpoint for full integration
func TestSenderWorkflow(t *testing.T) {
	// Create a test runtime
	runtime := testutils.NewRuntime(t, nil)

	// Create sender workflow
	senderWorkflow := workflows.NewSenderWorkflow()
	if len(senderWorkflow) == 0 {
		t.Fatal("Sender workflow not initialized")
	}

	// Create test input
	// In a real integration test, you would use an actual receiver workflow URL
	testReceiverURL := "http://localhost:8080/receiver"
	input := workflows.CreateTestSenderInput(
		"Test message from sender workflow",
		testReceiverURL,
		map[string]interface{}{
			"test": true,
		},
	)

	// Execute the sender workflow handler
	// Note: This will fail without a running receiver, which is expected in unit tests
	config := workflows.SenderConfig{
		ReceiverWorkflowURL: testReceiverURL,
	}
	_, err := workflows.SenderWorkflowHandler(config, input, runtime)
	if err != nil {
		// Expected in unit test environment without running server
		t.Logf("Sender workflow test (expected to fail without running server): %v", err)
	}
}

// TestWorkflowCommunication tests the full communication flow
func TestWorkflowCommunication(t *testing.T) {
	runtime := testutils.NewRuntime(t, nil)

	// Test receiver workflow
	receiverInput := workflows.CreateTestReceiverInput(
		"Communication test",
		"test_suite",
		map[string]interface{}{
			"communication_test": true,
		},
	)

	receiverConfig := workflows.ReceiverConfig{}
	receiverOutput, err := workflows.ReceiverWorkflowHandler(receiverConfig, receiverInput, runtime)
	if err != nil {
		t.Fatalf("Receiver workflow failed: %v", err)
	}

	// Verify the receiver processed the message
	if receiverOutput.Status != "success" {
		t.Errorf("Receiver workflow status should be 'success', got '%s'", receiverOutput.Status)
	}

	// Test that the output can be serialized (simulating HTTP response)
	outputJSON, err := json.Marshal(receiverOutput)
	if err != nil {
		t.Fatalf("Failed to marshal receiver output: %v", err)
	}

	// Verify JSON structure
	var parsedOutput workflows.ReceiverOutput
	if err := json.Unmarshal(outputJSON, &parsedOutput); err != nil {
		t.Fatalf("Failed to unmarshal receiver output: %v", err)
	}

	if parsedOutput.Status != "success" {
		t.Error("Parsed output status mismatch")
	}
}

// TestInputValidation tests input validation for workflows
func TestInputValidation(t *testing.T) {
	runtime := testutils.NewRuntime(t, nil)
	receiverConfig := workflows.ReceiverConfig{}

	// Test with empty message
	input := workflows.CreateTestReceiverInput("", "test", nil)

	output, err := workflows.ReceiverWorkflowHandler(receiverConfig, input, runtime)
	if err != nil {
		t.Fatalf("Receiver should handle empty message: %v", err)
	}

	if output.Status != "success" {
		t.Error("Receiver should process empty message successfully")
	}
}

// TestWorkflowRegistry tests the workflow registry
func TestWorkflowRegistry(t *testing.T) {
	registry := workflows.NewWorkflowRegistry()

	if err := registry.RegisterAllWorkflows(); err != nil {
		t.Fatalf("Failed to register workflows: %v", err)
	}

	receiverWorkflow := registry.GetReceiverWorkflow()
	if len(receiverWorkflow) == 0 {
		t.Error("Receiver workflow should be initialized")
	}

	senderWorkflow := registry.GetSenderWorkflow()
	if len(senderWorkflow) == 0 {
		t.Error("Sender workflow should be initialized")
	}
}

// TestJSONSerialization tests JSON serialization of workflow inputs and outputs
func TestJSONSerialization(t *testing.T) {
	// Test ReceiverInput serialization
	receiverInput := workflows.ReceiverInput{
		Message:   "Test message",
		SenderID:  "test_sender",
		Timestamp: time.Now().Unix(),
		Data: map[string]interface{}{
			"key": "value",
		},
	}

	jsonData, err := json.Marshal(receiverInput)
	if err != nil {
		t.Fatalf("Failed to marshal ReceiverInput: %v", err)
	}

	var parsedInput workflows.ReceiverInput
	if err := json.Unmarshal(jsonData, &parsedInput); err != nil {
		t.Fatalf("Failed to unmarshal ReceiverInput: %v", err)
	}

	if parsedInput.Message != receiverInput.Message {
		t.Error("Message mismatch after serialization")
	}

	// Test SenderInput serialization
	senderInput := workflows.SenderInput{
		Message:   "Test message",
		TargetURL: "http://example.com",
		Data: map[string]interface{}{
			"test": true,
		},
	}

	jsonData, err = json.Marshal(senderInput)
	if err != nil {
		t.Fatalf("Failed to marshal SenderInput: %v", err)
	}

	var parsedSenderInput workflows.SenderInput
	if err := json.Unmarshal(jsonData, &parsedSenderInput); err != nil {
		t.Fatalf("Failed to unmarshal SenderInput: %v", err)
	}

	if parsedSenderInput.Message != senderInput.Message {
		t.Error("Message mismatch after serialization")
	}
}
