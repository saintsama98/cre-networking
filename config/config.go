package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	ReceiverWorkflowURL string
	GatewayURL          string
	PrivateKey          string
	AuthorizedKeys      []string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		ReceiverWorkflowURL: getEnv("RECEIVER_WORKFLOW_URL", ""),
		GatewayURL:          getEnv("GATEWAY_URL", "https://01.gateway.zone-a.cre.chain.link"),
		PrivateKey:          getEnv("PRIVATE_KEY", ""),
		AuthorizedKeys:      getEnvSlice("AUTHORIZED_KEYS", []string{}),
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvSlice retrieves a comma-separated environment variable as a slice
func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple split by comma - in production, handle more complex cases
		result := []string{}
		for _, v := range splitString(value, ",") {
			if trimmed := trimSpace(v); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultValue
}

// Helper functions
func splitString(s, sep string) []string {
	result := []string{}
	current := ""
	for _, char := range s {
		if string(char) == sep {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

