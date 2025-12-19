# Resources

UI resources represent renderable content that can be returned from MCP servers to clients.

## UIResource

A complete UI resource with metadata and content.

### Definition

```go
type UIResource struct {
    URI         string    // Resource URI (e.g., "ui://dashboard/main")
    Name        string    // Human-readable name
    Description string    // Optional description
    Content     UIContent // The content to render
}
```

### Example

```go
resource := &mcpui.UIResource{
    URI:         "ui://status/overview",
    Name:        "Status Overview",
    Description: "Real-time system status dashboard",
    Content: &mcpui.HTMLContent{
        HTML: "<div>All systems operational</div>",
    },
}
```

### URI Format

UI resource URIs follow the pattern: `ui://{category}/{name}`

| Component | Description | Example |
|-----------|-------------|---------|
| Scheme | Always `ui://` | `ui://` |
| Category | Logical grouping | `dashboard`, `audio`, `scene` |
| Name | Specific resource | `main`, `mixer`, `preview` |

Examples:
- `ui://dashboard/main` - Main dashboard
- `ui://audio/mixer` - Audio mixer panel
- `ui://scene/preview` - Scene preview window
- `ui://settings/config` - Configuration panel

## UIResourceContents

MCP-compatible resource contents for returning in tool responses.

### Definition

```go
type UIResourceContents struct {
    URI      string `json:"uri"`
    MimeType string `json:"mimeType"`
    Text     string `json:"text,omitempty"`
    Blob     string `json:"blob,omitempty"`
}
```

### Creating Resource Contents

Use `NewUIResourceContents` to create properly formatted contents:

```go
func NewUIResourceContents(uri string, content UIContent) (*UIResourceContents, error)
```

### Example

```go
content := &mcpui.HTMLContent{
    HTML: "<div>Hello, World!</div>",
}

rc, err := mcpui.NewUIResourceContents("ui://greeting/hello", content)
if err != nil {
    log.Fatal(err)
}

// rc.URI = "ui://greeting/hello"
// rc.MimeType = "text/html"
// rc.Text = "<div>Hello, World!</div>"
```

### MCP Integration

Resource contents are designed to be returned as part of MCP tool responses:

```go
func handleTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Create UI content
    content := &mcpui.HTMLContent{
        HTML: generateDashboard(),
    }

    // Create resource contents
    rc, err := mcpui.NewUIResourceContents("ui://dashboard/main", content)
    if err != nil {
        return nil, err
    }

    // Return as embedded resource
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            {
                Type: "resource",
                Resource: &mcp.ResourceContents{
                    URI:      rc.URI,
                    MimeType: rc.MimeType,
                    Text:     rc.Text,
                },
            },
        },
    }, nil
}
```

## BlobContent

For binary content like images, use the `Blob` field instead of `Text`:

### Definition

```go
type BlobContent struct {
    Data     []byte // Binary data
    MimeType string // MIME type (e.g., "image/png")
}
```

### Example

```go
content := &mcpui.BlobContent{
    Data:     imageBytes,
    MimeType: "image/png",
}

// Creates base64-encoded blob in resource contents
rc, err := mcpui.NewUIResourceContents("ui://image/preview", content)
// rc.Blob contains base64-encoded data
```

## Validation

Use `ValidateResource` to check resources before use:

```go
func ValidateResource(resource *UIResource) error
```

### Validation Rules

| Field | Rules |
|-------|-------|
| URI | Must not be empty, must start with `ui://` |
| Name | Must not be empty |
| Content | Must not be nil, must pass content validation |

### Example

```go
resource := &mcpui.UIResource{
    URI:  "invalid-uri",
    Name: "Test",
    Content: &mcpui.HTMLContent{HTML: "<div>Test</div>"},
}

err := mcpui.ValidateResource(resource)
// Error: URI must start with "ui://"
```

## Resource Lists

When returning multiple resources, use a slice:

```go
resources := []*mcpui.UIResource{
    {
        URI:     "ui://dashboard/status",
        Name:    "Status",
        Content: &mcpui.HTMLContent{HTML: statusHTML},
    },
    {
        URI:     "ui://dashboard/metrics",
        Name:    "Metrics",
        Content: &mcpui.HTMLContent{HTML: metricsHTML},
    },
}
```

## Best Practices

1. **Use meaningful URIs** - URIs should describe the resource's purpose
2. **Include descriptions** - Help users understand what the resource shows
3. **Validate before sending** - Check resources with `ValidateResource`
4. **Keep URIs consistent** - Use the same URI for the same logical resource
5. **Handle errors gracefully** - Always check `NewUIResourceContents` errors
