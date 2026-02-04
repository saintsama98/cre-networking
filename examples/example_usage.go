package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"cre-networking-test/config"
	"cre-networking-test/workflows"
	"github.com/smartcontractkit/cre-sdk-go/cre/testutils"
)

// getTestRuntime creates a test runtime for examples
// In production, the runtime comes from CRE
func getTestRuntime() *testutils.TestRuntime {
	// Create a minimal test to get a real testing.T
	t := &testing.T{}
	return testutils.NewRuntime(t, nil)
}


// ExampleReceiverWorkflowUsage demonstrates how to use the receiver workflow
func ExampleReceiverWorkflowUsage() {
	fmt.Println("=== Receiver Workflow Example ===")
	
	// Create a test runtime (in production, this comes from CRE)
	runtime := getTestRuntime()
	
	receiverConfig := workflows.ReceiverConfig{}

	// Example input
	input := workflows.CreateTestReceiverInput(
		"Hello from example!",
		"example_sender",
		map[string]interface{}{
			"user_id":    "user123",
			"action":     "test_communication",
			"metadata":   map[string]string{"env": "development"},
		},
	)

	// Execute the receiver workflow
	output, err := workflows.ReceiverWorkflowHandler(receiverConfig, input, runtime)
	if err != nil {
		log.Fatalf("Failed to execute receiver workflow: %v", err)
	}

	// Print results
	outputJSON, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println("Receiver Workflow Output:")
	fmt.Println(string(outputJSON))
	fmt.Println()
}

// ExampleSenderWorkflowUsage demonstrates how to use the sender workflow
func ExampleSenderWorkflowUsage() {
	fmt.Println("=== Sender Workflow Example ===")
	
	cfg := config.LoadConfig()
	runtime := getTestRuntime()

	// Create sender workflow with receiver URL
	receiverURL := cfg.ReceiverWorkflowURL
	if receiverURL == "" {
		receiverURL = "https://01.gateway.zone-a.cre.chain.link/workflows/receiver-workflow-id"
		fmt.Println("Note: Using placeholder URL. Set RECEIVER_WORKFLOW_URL for actual communication.")
	}

	senderConfig := workflows.SenderConfig{
		ReceiverWorkflowURL: receiverURL,
	}

	// Example input
	input := workflows.CreateTestSenderInput(
		"Test communication between workflows",
		"",
		map[string]interface{}{
			"source":      "sender_workflow",
			"destination": "receiver_workflow",
			"payload": map[string]interface{}{
				"data": "example payload",
			},
		},
	)

	// Execute the sender workflow
	output, err := workflows.SenderWorkflowHandler(senderConfig, input, runtime)
	if err != nil {
		log.Printf("Failed to execute sender workflow (expected if receiver URL not accessible): %v", err)
		return
	}

	// Print results
	outputJSON, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println("Sender Workflow Output:")
	fmt.Println(string(outputJSON))
	fmt.Println()
}

// ExampleFullCommunicationFlow demonstrates the complete communication flow
func ExampleFullCommunicationFlow() {
	fmt.Println("=== Full Communication Flow Example ===")
	
	runtime := getTestRuntime()

	// Step 1: Receiver workflow receives a message
	fmt.Println("Step 1: Receiver Workflow")
	receiverConfig := workflows.ReceiverConfig{}
	receiverInput := workflows.CreateTestReceiverInput(
		"Initial message",
		"external_system",
		map[string]interface{}{
			"step": 1,
		},
	)

	receiverOutput, err := workflows.ReceiverWorkflowHandler(receiverConfig, receiverInput, runtime)
	if err != nil {
		log.Fatalf("Receiver workflow failed: %v", err)
	}

	fmt.Printf("Receiver Status: %s\n", receiverOutput.Status)
	fmt.Printf("Receiver Response: %s\n", receiverOutput.Response)

	// Step 2: Sender workflow sends a message (simulated)
	fmt.Println("\nStep 2: Sender Workflow")
	cfg := config.LoadConfig()
	receiverURL := cfg.ReceiverWorkflowURL
	if receiverURL == "" {
		fmt.Println("Note: Set RECEIVER_WORKFLOW_URL to test actual HTTP communication")
		receiverURL = "http://localhost:8080/receiver" // Placeholder
	}

	senderConfig := workflows.SenderConfig{
		ReceiverWorkflowURL: receiverURL,
	}
	senderInput := workflows.CreateTestSenderInput(
		"Follow-up message from sender",
		"",
		map[string]interface{}{
			"step":           2,
			"previous_status": receiverOutput.Status,
		},
	)

	senderOutput, err := workflows.SenderWorkflowHandler(senderConfig, senderInput, runtime)
	if err != nil {
		fmt.Printf("Sender workflow error (expected if receiver URL not accessible): %v\n", err)
	} else {
		fmt.Printf("Sender Status: %s\n", senderOutput.Status)
		fmt.Printf("Response Status: %s\n", senderOutput.ResponseStatus)
	}
	fmt.Println()
}

func main() {
	fmt.Println("CRE Workflow Communication Examples")
	fmt.Println("===================================\n")

	// Run examples
	ExampleReceiverWorkflowUsage()
	ExampleSenderWorkflowUsage()
	ExampleFullCommunicationFlow()

	fmt.Println("=== Workflow Registry Example ===")
	registry := workflows.NewWorkflowRegistry()
	if err := registry.RegisterAllWorkflows(); err != nil {
		log.Fatalf("Failed to register workflows: %v", err)
	}
	fmt.Println("✓ All workflows registered successfully")
	fmt.Println()
	fmt.Println("Note: In production, workflows are registered via cre.yaml configuration")
	fmt.Println("      when deploying to CRE environment.")
}
