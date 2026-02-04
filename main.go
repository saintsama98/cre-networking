package main

import (
	"fmt"
	"log"
	"os"

	"cre-networking-test/workflows"
)

// main is the entry point for the CRE workflow application
// In CRE, workflows are registered via cre.yaml configuration file
// This main function serves as a reference and can be used for local testing
func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "register":
		registerWorkflows()
	case "test":
		runTests()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("CRE Networking Test - Workflow Management")
	fmt.Println("=========================================")
	fmt.Println("Usage: go run main.go <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  register  - Register all workflows (for reference)")
	fmt.Println("  test      - Run local workflow tests")
	fmt.Println("  help      - Show this help message")
	fmt.Println()
	fmt.Println("Note: Workflows are registered via cre.yaml configuration file")
	fmt.Println("      when deploying to CRE environment.")
}

func registerWorkflows() {
	fmt.Println("Registering workflows...")

	registry := workflows.NewWorkflowRegistry()

	if err := registry.RegisterAllWorkflows(); err != nil {
		log.Fatalf("Failed to register workflows: %v", err)
	}

	fmt.Println("✓ Receiver workflow registered")
	fmt.Println("✓ Sender workflow registered")
	fmt.Println()
	fmt.Println("Workflows are ready for deployment via cre.yaml")
}

func runTests() {
	fmt.Println("Running local workflow tests...")
	fmt.Println("Note: For full integration tests, use: go test ./tests/...")
	fmt.Println()
	fmt.Println("Use the following commands for testing:")
	fmt.Println("  go test -v ./tests/...")
	fmt.Println("  go run ./examples/example_usage.go")
}
