# Kalinga Backend

MCP (Model Context Protocol) server for the Kalinga financial intelligence assistant.

## Overview

This Go backend implements a JSON-RPC 2.0 server that exposes financial tools for the Kalinga frontend to interact with via Gemini function calling.

## Available Tools

1. **check_status** - Verifies connection health and returns server status
2. **connect_momo_account** - Connects mobile money accounts with phone number and provider
3. **build_financial_profile** - Builds loan readiness and credit scores from transaction history
4. **verify_claim** - Checks refund/payment claims against suspicious transaction patterns

## Setup

### Prerequisites
- Go 1.27.1 or higher

### Installation

1. Navigate to the backend directory:
```bash
cd /Users/chawanangwa/Desktop/My Projects/Kalinga/backend/kalinga-backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the server:
```bash
go build ./cmd/server
```

## Running the Server

### Development
```bash
go run cmd/server/main.go
```

### Production
```bash
./server
```

The server will start on `http://localhost:8080/mcp`

## API Endpoints

### POST /mcp

JSON-RPC 2.0 endpoint for tool calls.

**Example Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "check_status",
    "arguments": {}
  }
}
```

**Example Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "status": "ok",
    "service": "kalinga",
    "version": "0.1.0",
    "health": "healthy"
  }
}
```

## Tool Parameters

### check_status
- No parameters required

### connect_momo_account
- `phone_number` (string): The phone number to connect
- `provider` (string): The mobile money provider (e.g., MTN, Airtel)

### build_financial_profile
- `user_id` (string): The user identifier

### verify_claim
- `claim_id` (string): The claim identifier to verify
- `amount` (number): The claim amount
- `merchant` (string): The merchant name

## Development

The server is implemented as a simple HTTP handler that processes JSON-RPC requests and routes them to the appropriate tool handlers. Each tool returns JSON responses that can be consumed by the frontend via the Gemini API.
