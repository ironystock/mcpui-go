# mcpui-go Documentation

This directory contains detailed documentation for the mcpui-go SDK.

## Contents

| Document | Description |
|----------|-------------|
| [content-types.md](content-types.md) | UI content types (HTML, URL, Remote DOM) |
| [resources.md](resources.md) | UI resources and resource contents |
| [actions.md](actions.md) | Action types, payloads, and parsing |
| [responses.md](responses.md) | Response builders and message types |
| [handlers.md](handlers.md) | Action handlers and routing |
| [integration.md](integration.md) | MCP server integration guide |

## Quick Reference

### Package Structure

```
mcpui-go/
├── content.go      # HTMLContent, URLContent, RemoteDOMContent
├── resource.go     # UIResource, UIResourceContents
├── action.go       # UIAction, action types, payloads
├── response.go     # UIResponse builders
├── handler.go      # UIActionHandler, Router
└── doc.go          # Package documentation
```

### Import

```go
import "github.com/ironystock/mcpui-go"
```

### Key Types

| Type | Purpose |
|------|---------|
| `UIContent` | Interface for content types (HTML, URL, Remote DOM) |
| `UIResource` | A complete UI resource with URI, name, and content |
| `UIResourceContents` | MCP-compatible resource contents |
| `UIAction` | User interaction from UI layer |
| `UIResponse` | Response message to UI layer |
| `UIActionHandler` | Function to handle UI actions |
| `Router` | Routes actions to handlers |

## See Also

- [README.md](../README.md) - Package overview and quick start
- [examples/](../examples/) - Runnable example programs
- [Go Reference](https://pkg.go.dev/github.com/ironystock/mcpui-go@v0.1.0) - API documentation
