# CRE Networking Test

This project demonstrates inter-workflow communication in Chainlink's Cloud Runtime Environment (CRE) using HTTP triggers and the CRE Go SDK. It includes two workflows that communicate with each other via HTTP, following industry best practices.

## Overview

The project consists of:

1. **Receiver Workflow** - Accepts HTTP triggers and processes incoming data
2. **Sender Workflow** - Sends HTTP POST requests to other workflows
3. **Test Suite** - Integration tests for workflow communication
4. **Examples** - Usage examples demonstrating workflow communication

## Architecture

```
┌─────────────────┐         HTTP POST          ┌─────────────────┐
│  Sender         │ ──────────────────────────> │  Receiver       │
│  Workflow       │                             │  Workflow       │
│                 │ <────────────────────────── │                 │
│  (HTTP Client)  │      JSON Response         │  (HTTP Trigger) │
└─────────────────┘                             └─────────────────┘
```

### Communication Flow

1. **Sender Workflow** uses `http.SendRequest` from CRE SDK to send HTTP POST requests
2. **Receiver Workflow** accepts HTTP triggers via CRE's HTTP trigger mechanism
3. Communication uses JSON-RPC format with JWT authentication (in production)
4. All requests are cryptographically signed and validated

## Prerequisites

- Go 1.21 or higher
- CRE CLI (for deployment and simulation)
- Access to CRE environment (for deployed testing)
- Private key for signing requests (for production)

## Installation

1. Clone the repository:
```bash
cd cre-networking-test
```

2. Install dependencies:
```bash
make deps
# or manually:
go mod download
go mod tidy
```

3. Set up environment variables (copy from `.env.example`):
```bash
cp .env.example .env
# Edit .env with your configuration
```

## Configuration

### Environment Variables

Create a `.env` file with the following variables:

```bash
# Receiver Workflow URL (after deployment)
RECEIVER_WORKFLOW_URL=https://01.gateway.zone-a.cre.chain.link/workflows/{workflow-id}

# CRE Gateway URL
GATEWAY_URL=https://01.gateway.zone-a.cre.chain.link

# Private Key for signing requests (hex format: 0x...)
PRIVATE_KEY=0x...

# Authorized Keys (comma-separated EVM addresses)
AUTHORIZED_KEYS=0xAddress1,0xAddress2
```

### CRE Configuration

The `cre.yaml` file contains workflow deployment configuration. Update it with:
- Authorized keys for HTTP triggers
- Workflow entrypoints
- Network settings

## Project Structure

```
cre-networking-test/
├── workflows/
│   ├── receiver.go      # Receiver workflow implementation
│   ├── sender.go         # Sender workflow implementation
│   └── workflows.go      # Workflow registry and helpers
├── config/
│   └── config.go         # Configuration management
├── tests/
│   └── integration_test.go  # Integration tests
├── examples/
│   └── example_usage.go     # Usage examples
├── go.mod                 # Go module definition
├── cre.yaml              # CRE deployment configuration
├── Makefile              # Build and test commands
└── README.md             # This file
```

## Usage

### Running Examples

```bash
make run-examples
# or
go run ./examples/example_usage.go
```

### Running Tests

```bash
make test
# or
go test -v ./tests/...
```

### Building

```bash
make build
```

### Local Simulation

For local testing without deployment:

```bash
# Simulate receiver workflow
cre workflow simulate \
  --workflow receiver-workflow \
  --input '{"message":"test message","sender_id":"test"}'

# Simulate sender workflow
cre workflow simulate \
  --workflow sender-workflow \
  --input '{"message":"test","target_url":"http://localhost:8080/receiver"}'
```

## Workflow Details

### Receiver Workflow

The receiver workflow accepts HTTP triggers and processes incoming data:

**Input:**
```json
{
  "message": "Hello from sender",
  "sender_id": "sender_workflow",
  "timestamp": 1234567890,
  "data": {
    "key": "value"
  }
}
```

**Output:**
```json
{
  "status": "success",
  "received_at": 1234567890,
  "processed": true,
  "response": "Successfully received and processed: Hello from sender",
  "echoed_data": {
    "key": "value",
    "processed_by": "receiver_workflow",
    "original_message": "Hello from sender"
  }
}
```

### Sender Workflow

The sender workflow sends HTTP POST requests to other workflows:

**Input:**
```json
{
  "message": "Test communication",
  "target_url": "https://gateway.cre.chain.link/workflows/receiver-id",
  "data": {
    "test": true
  }
}
```

**Output:**
```json
{
  "status": "success",
  "sent_at": 1234567890,
  "response_status": "200",
  "response_data": {
    "status": "success",
    "processed": true
  },
  "request_id": "req_1234567890"
}
```

## HTTP Communication Format

CRE workflows communicate using:

1. **HTTP POST requests** to CRE gateway endpoints
2. **JSON-RPC format** for request bodies
3. **JWT authentication** with cryptographic signatures
4. **EVM address-based authorization** via `authorizedKeys`

### Request Format

```json
{
  "jsonrpc": "2.0",
  "method": "workflow_execute",
  "params": {
    "workflow_id": "your-workflow-id",
    "input": {
      "message": "your data"
    }
  },
  "id": 1
}
```

### Authentication

- Requests must be signed with a private key
- Corresponding EVM address must be in `authorizedKeys`
- JWT tokens are automatically generated by CRE tools

## Testing Communication

### Unit Testing

Run the test suite:
```bash
go test -v ./tests/...
```

### Integration Testing

For full integration testing:

1. Deploy both workflows to CRE
2. Get the receiver workflow URL
3. Update `.env` with the receiver URL
4. Trigger the sender workflow with the receiver URL

### Using CRE HTTP Trigger Tool

For testing deployed workflows, use the `cre-http-trigger` tool:

```bash
# Install (requires Bun)
bun install -g cre-http-trigger

# Configure .env with PRIVATE_KEY and GATEWAY_URL
# Run the tool
cre-http-trigger
```

This starts a local proxy server on port 2000 that handles JWT generation automatically.

## Best Practices

1. **Single Execution**: Use `CacheSettings` in HTTP requests to prevent duplicate executions across DON nodes
2. **Error Handling**: Always handle errors and return appropriate status codes
3. **Logging**: Use `sdk.Log()` for debugging and monitoring
4. **Input Validation**: Validate all inputs before processing
5. **Security**: Never expose private keys in code or logs
6. **Idempotency**: Design workflows to be idempotent when possible

## Industry Standards

This implementation follows:

- ✅ **CRE SDK v1.0.0** standards
- ✅ **HTTP trigger** best practices
- ✅ **Error handling** patterns
- ✅ **Structured logging**
- ✅ **Configuration management**
- ✅ **Test coverage**
- ✅ **Documentation** standards

## Troubleshooting

### Common Issues

1. **"target URL not specified"**
   - Set `RECEIVER_WORKFLOW_URL` in `.env`
   - Or provide `target_url` in sender workflow input

2. **Authentication failures**
   - Verify `PRIVATE_KEY` is set correctly
   - Check `authorizedKeys` in workflow configuration
   - Ensure EVM address matches the private key

3. **Connection errors**
   - Verify `GATEWAY_URL` is correct
   - Check network connectivity
   - Ensure workflow is deployed and active

## References

- [CRE Go SDK Documentation](https://pkg.go.dev/github.com/smartcontractkit/cre-sdk-go/sdk)
- [CRE HTTP Trigger Guide](https://docs.chain.link/cre/guides/workflow/using-triggers/http-trigger/overview-go)
- [CRE HTTP Client Guide](https://docs.chain.link/cre/guides/workflow/using-http-client/post-request-go)
- [CRE SDK Reference](https://docs.chain.link/cre/reference/sdk/overview-go)

## License

This project is for testing and demonstration purposes.

## Contributing

When contributing:
1. Follow Go best practices
2. Add tests for new features
3. Update documentation
4. Ensure code passes linting

## Support

For CRE-specific issues, refer to:
- [CRE Documentation](https://docs.chain.link/cre)
- [CRE GitHub Repository](https://github.com/smartcontractkit/cre-sdk-go)

