# mcpui-go Examples

This directory contains runnable examples demonstrating the mcpui-go SDK.

## Examples

| Example | Description |
|---------|-------------|
| [basic](basic/) | Minimal HTML resource creation |
| [router](router/) | Action routing with multiple handlers |
| [remote-dom](remote-dom/) | Remote DOM content with React framework |
| [action-handling](action-handling/) | Complete action→response flow |
| [mcp-integration](mcp-integration/) | Integration with MCP server |

## Running Examples

Each example can be run with:

```bash
cd examples/<name>
go run main.go
```

## Example Descriptions

### basic

Shows the simplest possible use case: creating HTML content and serializing it for MCP.

```bash
cd examples/basic
go run main.go
```

### router

Demonstrates the Router for handling different action types and resource-specific handlers.

```bash
cd examples/router
go run main.go
```

### remote-dom

Shows how to create Remote DOM content for dynamic, framework-based UIs.

```bash
cd examples/remote-dom
go run main.go
```

### action-handling

Complete example of handling UI actions and sending responses.

```bash
cd examples/action-handling
go run main.go
```

### mcp-integration

Shows how to integrate mcpui-go with an MCP server.

```bash
cd examples/mcp-integration
go run main.go
```
