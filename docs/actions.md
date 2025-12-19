# Actions

UI actions represent user interactions sent from the UI layer to the MCP server.

## UIAction

The primary action structure containing interaction details.

### Definition

```go
type UIAction struct {
    MessageID string          `json:"messageId"`
    Type      ActionType      `json:"type"`
    Payload   json.RawMessage `json:"payload"`
}
```

### Fields

| Field | Description |
|-------|-------------|
| `MessageID` | Unique identifier for request/response correlation |
| `Type` | Action type (tool, prompt, resource, custom) |
| `Payload` | Type-specific payload data |

## Action Types

```go
const (
    ActionTypeTool     ActionType = "tool"
    ActionTypePrompt   ActionType = "prompt"
    ActionTypeResource ActionType = "resource"
    ActionTypeCustom   ActionType = "custom"
)
```

### ActionTypeTool

Requests execution of an MCP tool.

```go
type ToolActionPayload struct {
    Name       string         `json:"name"`
    Parameters map[string]any `json:"parameters,omitempty"`
}
```

Example:
```go
action := &mcpui.UIAction{
    MessageID: "msg-123",
    Type:      mcpui.ActionTypeTool,
    Payload:   json.RawMessage(`{"name":"get_status","parameters":{}}`),
}
```

### ActionTypePrompt

Requests execution of an MCP prompt.

```go
type PromptActionPayload struct {
    Name      string         `json:"name"`
    Arguments map[string]any `json:"arguments,omitempty"`
}
```

Example:
```go
action := &mcpui.UIAction{
    MessageID: "msg-456",
    Type:      mcpui.ActionTypePrompt,
    Payload:   json.RawMessage(`{"name":"health-check","arguments":{}}`),
}
```

### ActionTypeResource

Requests read of an MCP resource.

```go
type ResourceActionPayload struct {
    URI string `json:"uri"`
}
```

Example:
```go
action := &mcpui.UIAction{
    MessageID: "msg-789",
    Type:      mcpui.ActionTypeResource,
    Payload:   json.RawMessage(`{"uri":"obs://scene/Gaming"}`),
}
```

### ActionTypeCustom

Application-specific custom actions.

```go
type CustomActionPayload struct {
    Action string         `json:"action"`
    Data   map[string]any `json:"data,omitempty"`
}
```

Example:
```go
action := &mcpui.UIAction{
    MessageID: "msg-abc",
    Type:      mcpui.ActionTypeCustom,
    Payload:   json.RawMessage(`{"action":"refresh","data":{"force":true}}`),
}
```

## UIActionRequest

The complete request containing action and context.

### Definition

```go
type UIActionRequest struct {
    SourceURI string    // URI of the resource that sent the action
    Action    *UIAction // The action to process
}
```

### Example

```go
request := &mcpui.UIActionRequest{
    SourceURI: "ui://audio/mixer",
    Action: &mcpui.UIAction{
        MessageID: "msg-123",
        Type:      mcpui.ActionTypeTool,
        Payload:   json.RawMessage(`{"name":"toggle_input_mute","parameters":{"inputName":"Mic"}}`),
    },
}
```

## Parsing Payloads

### ParseToolPayload

```go
func ParseToolPayload(action *UIAction) (*ToolActionPayload, error)
```

Example:
```go
payload, err := mcpui.ParseToolPayload(action)
if err != nil {
    return nil, err
}
fmt.Printf("Tool: %s, Params: %v\n", payload.Name, payload.Parameters)
```

### ParsePromptPayload

```go
func ParsePromptPayload(action *UIAction) (*PromptActionPayload, error)
```

Example:
```go
payload, err := mcpui.ParsePromptPayload(action)
if err != nil {
    return nil, err
}
fmt.Printf("Prompt: %s, Args: %v\n", payload.Name, payload.Arguments)
```

### ParseResourcePayload

```go
func ParseResourcePayload(action *UIAction) (*ResourceActionPayload, error)
```

Example:
```go
payload, err := mcpui.ParseResourcePayload(action)
if err != nil {
    return nil, err
}
fmt.Printf("Resource URI: %s\n", payload.URI)
```

### ParseCustomPayload

```go
func ParseCustomPayload(action *UIAction) (*CustomActionPayload, error)
```

Example:
```go
payload, err := mcpui.ParseCustomPayload(action)
if err != nil {
    return nil, err
}
fmt.Printf("Custom action: %s, Data: %v\n", payload.Action, payload.Data)
```

## UIActionResult

The result of processing an action.

### Definition

```go
type UIActionResult struct {
    Response any   // Success response data
    Error    error // Error if action failed
}
```

### Example

```go
// Success result
result := &mcpui.UIActionResult{
    Response: map[string]string{"status": "ok"},
}

// Error result
result := &mcpui.UIActionResult{
    Error: errors.New("tool not found"),
}
```

## Validation

Use `ValidateAction` to check actions before processing:

```go
func ValidateAction(action *UIAction) error
```

### Validation Rules

| Field | Rules |
|-------|-------|
| `MessageID` | Must not be empty |
| `Type` | Must be a valid action type |
| `Payload` | Must be valid JSON (can be empty object) |

### Example

```go
action := &mcpui.UIAction{
    MessageID: "",
    Type:      mcpui.ActionTypeTool,
    Payload:   json.RawMessage(`{}`),
}

err := mcpui.ValidateAction(action)
// Error: MessageID cannot be empty
```

## Best Practices

1. **Always validate actions** - Check with `ValidateAction` before processing
2. **Use typed payloads** - Parse to specific payload types for type safety
3. **Preserve MessageID** - Use the same MessageID in responses for correlation
4. **Handle unknown types** - Gracefully handle unrecognized action types
5. **Log actions** - Record actions for debugging and auditing
