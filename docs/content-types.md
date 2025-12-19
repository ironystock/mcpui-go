# Content Types

The mcpui-go SDK supports three types of UI content, each implementing the `UIContent` interface.

## UIContent Interface

```go
type UIContent interface {
    MarshalJSON() ([]byte, error)
    ContentType() string
}
```

All content types implement this interface, providing JSON serialization and MIME type information.

## HTMLContent

Inline HTML content rendered via iframe `srcdoc`. Best for simple, self-contained UIs.

### Definition

```go
type HTMLContent struct {
    HTML string // The HTML content to render
}
```

### MIME Type

`text/html`

### Example

```go
content := &mcpui.HTMLContent{
    HTML: `<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: sans-serif; padding: 20px; }
        .status { color: green; }
    </style>
</head>
<body>
    <h1>Status Dashboard</h1>
    <p class="status">All systems operational</p>
</body>
</html>`,
}

// Serialize to JSON
data, err := content.MarshalJSON()
// Result: {"mimeType":"text/html","text":"<!DOCTYPE html>..."}
```

### Use Cases

- Simple status displays
- Forms and input collection
- Static information cards
- Self-contained widgets

### Considerations

- No external dependencies (CSS/JS must be inline)
- Limited to what fits in a string
- No live updates without full content replacement

## URLContent

External URL content rendered via iframe `src`. Best for existing web applications.

### Definition

```go
type URLContent struct {
    URL string // The URL to load in the iframe
}
```

### MIME Type

`text/uri-list`

### Example

```go
content := &mcpui.URLContent{
    URL: "https://example.com/dashboard?theme=dark",
}

// Serialize to JSON
data, err := content.MarshalJSON()
// Result: {"mimeType":"text/uri-list","text":"https://example.com/dashboard?theme=dark"}
```

### Use Cases

- Existing web applications
- Complex UIs with many dependencies
- Third-party integrations
- Live dashboards

### Considerations

- Requires network access from client
- Subject to CORS and CSP policies
- External server must be available
- Full page load on each update

## RemoteDOMContent

Script-based UI using remote DOM rendering. Best for dynamic, framework-based UIs.

### Definition

```go
type RemoteDOMContent struct {
    Script    string    // JavaScript code to execute
    Framework Framework // UI framework (React, Vue, etc.)
}
```

### MIME Type

`application/vnd.mcp-ui.remote-dom`

### Frameworks

```go
const (
    FrameworkReact   Framework = "react"
    FrameworkVue     Framework = "vue"
    FrameworkSvelte  Framework = "svelte"
    FrameworkVanilla Framework = "vanilla"
)
```

### Example

```go
content := &mcpui.RemoteDOMContent{
    Script: `
        const [count, setCount] = React.useState(0);
        return React.createElement('div', null,
            React.createElement('p', null, 'Count: ' + count),
            React.createElement('button',
                { onClick: () => setCount(count + 1) },
                'Increment'
            )
        );
    `,
    Framework: mcpui.FrameworkReact,
}

// Serialize to JSON
data, err := content.MarshalJSON()
// Result: {"mimeType":"application/vnd.mcp-ui.remote-dom","text":"{\"script\":\"...\",\"framework\":\"react\"}"}
```

### Use Cases

- Interactive components
- Real-time updates
- Complex state management
- Framework-specific patterns

### Considerations

- Requires client framework support
- More complex than HTML content
- Framework version compatibility
- Larger payload size

## Content Validation

Use `ValidateContent` to check content before use:

```go
func ValidateContent(content UIContent) error
```

### Validation Rules

| Content Type | Rules |
|--------------|-------|
| `HTMLContent` | HTML must not be empty |
| `URLContent` | URL must not be empty, must be valid URL |
| `RemoteDOMContent` | Script must not be empty, Framework must be valid |

### Example

```go
content := &mcpui.HTMLContent{HTML: ""}
err := mcpui.ValidateContent(content)
// Error: HTML content cannot be empty
```

## Choosing a Content Type

| Need | Recommended Type |
|------|------------------|
| Simple static display | `HTMLContent` |
| Existing web app | `URLContent` |
| Interactive components | `RemoteDOMContent` |
| No external dependencies | `HTMLContent` |
| Complex state management | `RemoteDOMContent` |
| Third-party integration | `URLContent` |
