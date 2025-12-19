# Responses

UI responses are messages sent from the MCP server back to the UI layer in response to actions.

## UIResponse

The primary response structure.

### Definition

```go
type UIResponse struct {
    MessageID string       `json:"messageId"`
    Type      ResponseType `json:"type"`
    Payload   any          `json:"payload,omitempty"`
    Error     *ErrorInfo   `json:"error,omitempty"`
}
```

### Fields

| Field | Description |
|-------|-------------|
| `MessageID` | Correlates with the action's MessageID |
| `Type` | Response type (received, success, error) |
| `Payload` | Response data (for success responses) |
| `Error` | Error information (for error responses) |

## Response Types

```go
const (
    ResponseTypeReceived ResponseType = "received"
    ResponseTypeSuccess  ResponseType = "success"
    ResponseTypeError    ResponseType = "error"
)
```

### ResponseTypeReceived

Acknowledges that the action was received and is being processed.

```go
resp := mcpui.NewReceivedResponse("msg-123")
// {
//   "messageId": "msg-123",
//   "type": "received"
// }
```

### ResponseTypeSuccess

Indicates successful completion with optional result data.

```go
resp := mcpui.NewSuccessResponse("msg-123", map[string]any{
    "status": "ok",
    "data":   result,
})
// {
//   "messageId": "msg-123",
//   "type": "success",
//   "payload": {"status": "ok", "data": ...}
// }
```

### ResponseTypeError

Indicates an error occurred during processing.

```go
resp := mcpui.NewErrorResponse("msg-123", errors.New("tool not found"))
// {
//   "messageId": "msg-123",
//   "type": "error",
//   "error": {"code": "error", "message": "tool not found"}
// }
```

## Response Builders

### NewReceivedResponse

Creates an acknowledgment response.

```go
func NewReceivedResponse(messageID string) *UIResponse
```

Example:
```go
// Immediately acknowledge receipt
ack := mcpui.NewReceivedResponse(action.MessageID)
sendToUI(ack)

// Then process asynchronously...
```

### NewSuccessResponse

Creates a success response with payload.

```go
func NewSuccessResponse(messageID string, payload any) *UIResponse
```

Example:
```go
result := map[string]any{
    "volume":  0.75,
    "muted":   false,
    "channel": "master",
}
resp := mcpui.NewSuccessResponse(action.MessageID, result)
```

### NewErrorResponse

Creates an error response from a Go error.

```go
func NewErrorResponse(messageID string, err error) *UIResponse
```

Example:
```go
err := errors.New("input 'Microphone' not found")
resp := mcpui.NewErrorResponse(action.MessageID, err)
```

### NewErrorResponseWithCode

Creates an error response with a custom error code.

```go
func NewErrorResponseWithCode(messageID, code, message string) *UIResponse
```

Example:
```go
resp := mcpui.NewErrorResponseWithCode(
    action.MessageID,
    "not_found",
    "The requested resource does not exist",
)
```

## ErrorInfo

Structured error information.

### Definition

```go
type ErrorInfo struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

### Common Error Codes

| Code | Description |
|------|-------------|
| `error` | Generic error (default) |
| `not_found` | Resource or tool not found |
| `invalid_request` | Malformed request |
| `unauthorized` | Permission denied |
| `timeout` | Operation timed out |
| `internal` | Internal server error |

## Response Flow

A typical action→response flow:

```go
func handleAction(ctx context.Context, action *mcpui.UIAction) {
    // 1. Send immediate acknowledgment
    ack := mcpui.NewReceivedResponse(action.MessageID)
    sendToUI(ack)

    // 2. Process the action
    result, err := processAction(ctx, action)

    // 3. Send final response
    var resp *mcpui.UIResponse
    if err != nil {
        resp = mcpui.NewErrorResponse(action.MessageID, err)
    } else {
        resp = mcpui.NewSuccessResponse(action.MessageID, result)
    }
    sendToUI(resp)
}
```

## JSON Serialization

Responses serialize cleanly to JSON:

```go
resp := mcpui.NewSuccessResponse("msg-123", map[string]any{
    "status": "ok",
})

data, err := json.Marshal(resp)
// {"messageId":"msg-123","type":"success","payload":{"status":"ok"}}
```

## Best Practices

1. **Always respond** - Every action should get at least one response
2. **Send acknowledgments** - For long operations, send `received` first
3. **Use meaningful error codes** - Help UI handle errors appropriately
4. **Include context in errors** - Make error messages actionable
5. **Preserve MessageID** - Always use the original action's MessageID
