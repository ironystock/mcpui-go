# mcpui-go

[![Go Reference](https://pkg.go.dev/badge/github.com/ironystock/mcpui-go@v0.1.0.svg)](https://pkg.go.dev/github.com/ironystock/mcpui-go@v0.1.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/ironystock/mcpui-go)](https://goreportcard.com/report/github.com/ironystock/mcpui-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitMCP](https://img.shields.io/endpoint?url=https://gitmcp.io/badge/ironystock/mcpui-go)](https://gitmcp.io/ironystock/mcpui-go)
[![Twitter Follow](https://img.shields.io/twitter/follow/ironystock?style=social)](https://twitter.com/ironystock)

A Go SDK for the MCP-UI protocol, enabling MCP servers to return interactive UI resources that clients can render in sandboxed iframes.

## Overview

MCP-UI extends the Model Context Protocol (MCP) to support rich, interactive user interfaces. This SDK provides the server-side implementation for Go-based MCP servers, allowing them to:

- Generate UI resources with HTML, URLs, or Remote DOM content
- Handle user interactions from UI components
- Send responses back to the UI layer
- Route actions to appropriate handlers

## Installation

```bash
go get github.com/ironystock/mcpui-go@v0.1.0
```

## Quick Start

### Creating a Simple HTML Resource

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/ironystock/mcpui-go"
)

func main() {
    // Create HTML content
    content := &mcpui.HTMLContent{
        HTML: `<div style="padding: 20px; font-family: sans-serif;">
            <h1>Hello, World!</h1>
            <p>This is rendered in a sandboxed iframe.</p>
        </div>`,
    }

    // Create resource contents for MCP response
    rc, _ := mcpui.NewUIResourceContents("ui://greeting/hello", content)

    data, _ := json.MarshalIndent(rc, "", "  ")
    fmt.Println(string(data))
}
```

### Handling UI Actions

```go
package main

import (
    "context"
    "fmt"

    "github.com/ironystock/mcpui-go"
)

func main() {
    router := mcpui.NewRouter()

    // Handle tool actions
    router.HandleType(mcpui.ActionTypeTool, mcpui.WrapToolHandler(
        func(ctx context.Context, toolName string, params map[string]any) (any, error) {
            fmt.Printf("Tool called: %s\n", toolName)
            return map[string]string{"status": "executed"}, nil
        },
    ))

    // Handle specific resource
    router.HandleResource("ui://dashboard/main",
        func(ctx context.Context, req *mcpui.UIActionRequest) (*mcpui.UIActionResult, error) {
            return &mcpui.UIActionResult{Response: "dashboard handled"}, nil
        },
    )

    fmt.Println("Router configured")
}
```

## Content Types

The SDK supports three types of UI content:

| Type | Description | Use Case |
|------|-------------|----------|
| `HTMLContent` | Inline HTML rendered via iframe srcdoc | Simple, self-contained UIs |
| `URLContent` | External URL rendered via iframe src | Existing web apps, complex UIs |
| `RemoteDOMContent` | Script-based UI using remote DOM | Dynamic, framework-based UIs |

### HTMLContent

Best for simple, self-contained UIs:

```go
content := &mcpui.HTMLContent{
    HTML: `<div>
        <h1>Status Dashboard</h1>
        <p>All systems operational</p>
    </div>`,
}
```

### URLContent

Best for existing web applications:

```go
content := &mcpui.URLContent{
    URL: "https://example.com/dashboard",
}
```

### RemoteDOMContent

Best for dynamic, framework-based UIs:

```go
content := &mcpui.RemoteDOMContent{
    Script:    "React.createElement('div', null, 'Hello from React!')",
    Framework: mcpui.FrameworkReact,
}
```

## Action Handling

When users interact with UI resources, the client sends `UIAction` messages. The SDK provides a flexible routing system:

```go
router := mcpui.NewRouter()

// Handle by action type
router.HandleType(mcpui.ActionTypeTool, toolHandler)
router.HandleType(mcpui.ActionTypePrompt, promptHandler)
router.HandleType(mcpui.ActionTypeResource, resourceHandler)

// Handle by specific resource URI
router.HandleResource("ui://audio/mixer", audioHandler)
router.HandleResource("ui://scene/preview", sceneHandler)

// Route an action
result, err := router.Route(ctx, &mcpui.UIActionRequest{
    SourceURI: "ui://audio/mixer",
    Action:    action,
})
```

## Response Messages

Send responses back to the UI:

```go
// Acknowledge receipt
ack := mcpui.NewReceivedResponse(action.MessageID)

// Success response
resp := mcpui.NewSuccessResponse(action.MessageID, result)

// Error response
errResp := mcpui.NewErrorResponse(action.MessageID, err)
```

## Documentation

- [Content Types](docs/content-types.md) - HTMLContent, URLContent, RemoteDOMContent
- [Resources](docs/resources.md) - UIResource, UIResourceContents, validation
- [Actions](docs/actions.md) - UIAction types, payloads, parsing
- [Responses](docs/responses.md) - UIResponse builders, success/error handling
- [Handlers](docs/handlers.md) - UIActionHandler, Router, typed wrappers
- [Integration](docs/integration.md) - How to integrate with MCP servers

## Examples

See the [examples](examples/) directory for complete, runnable examples:

- [basic](examples/basic/) - Minimal HTML resource example
- [router](examples/router/) - Router with multiple handlers
- [remote-dom](examples/remote-dom/) - Remote DOM with React framework
- [action-handling](examples/action-handling/) - Complete action→response flow
- [mcp-integration](examples/mcp-integration/) - Integration with MCP server

## License

MIT License - see [LICENSE](LICENSE) for details.

## Related

- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - Official MCP Go SDK
- [MCP-UI Specification](https://mcpui.dev/) - Protocol documentation
- [agentic-obs](https://github.com/ironystock/agentic-obs) - MCP server using this SDK
