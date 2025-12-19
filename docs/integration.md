# MCP Server Integration

This guide shows how to integrate mcpui-go with an MCP server built using the official [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk).

## Overview

The integration has two main aspects:

1. **Returning UI resources** - Include UI content in MCP tool responses
2. **Handling UI actions** - Process user interactions from UI components

## Returning UI Resources

### Basic Tool Response

Include UI content in a tool's response using embedded resources:

```go
import (
    "github.com/modelcontextprotocol/go-sdk/mcp"
    "github.com/ironystock/mcpui-go"
)

func handleDashboardTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Create HTML content for the dashboard
    content := &mcpui.HTMLContent{
        HTML: `<div style="padding: 20px;">
            <h1>System Dashboard</h1>
            <p>Status: All systems operational</p>
        </div>`,
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

### Dynamic Content Generation

Generate UI content based on server state:

```go
func handleStatusTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Get current status
    status := getSystemStatus()

    // Generate dynamic HTML
    html := fmt.Sprintf(`
        <div style="font-family: sans-serif; padding: 20px;">
            <h2>System Status</h2>
            <ul>
                <li>Connected: %v</li>
                <li>Recording: %v</li>
                <li>Streaming: %v</li>
            </ul>
            <p>Last updated: %s</p>
        </div>
    `, status.Connected, status.Recording, status.Streaming, time.Now().Format(time.RFC3339))

    content := &mcpui.HTMLContent{HTML: html}
    rc, _ := mcpui.NewUIResourceContents("ui://status/overview", content)

    return &mcp.CallToolResult{
        Content: []mcp.Content{
            {Type: "resource", Resource: &mcp.ResourceContents{
                URI: rc.URI, MimeType: rc.MimeType, Text: rc.Text,
            }},
        },
    }, nil
}
```

### Multiple Resources

Return multiple UI resources in a single response:

```go
func handleMultiPanelTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Create multiple panels
    statusRC, _ := mcpui.NewUIResourceContents("ui://panel/status",
        &mcpui.HTMLContent{HTML: generateStatusPanel()})

    metricsRC, _ := mcpui.NewUIResourceContents("ui://panel/metrics",
        &mcpui.HTMLContent{HTML: generateMetricsPanel()})

    return &mcp.CallToolResult{
        Content: []mcp.Content{
            {Type: "resource", Resource: &mcp.ResourceContents{
                URI: statusRC.URI, MimeType: statusRC.MimeType, Text: statusRC.Text,
            }},
            {Type: "resource", Resource: &mcp.ResourceContents{
                URI: metricsRC.URI, MimeType: metricsRC.MimeType, Text: metricsRC.Text,
            }},
        },
    }, nil
}
```

## Handling UI Actions

### Setting Up Action Handling

UI actions come through a dedicated channel or callback. Set up a router to handle them:

```go
type Server struct {
    mcpServer *mcp.Server
    uiRouter  *mcpui.Router
}

func NewServer() *Server {
    s := &Server{
        uiRouter: mcpui.NewRouter(),
    }

    // Register UI action handlers
    s.setupUIHandlers()

    return s
}

func (s *Server) setupUIHandlers() {
    // Handle tool actions
    s.uiRouter.HandleType(mcpui.ActionTypeTool, mcpui.WrapToolHandler(
        func(ctx context.Context, toolName string, params map[string]any) (any, error) {
            // Execute the MCP tool
            return s.executeTool(ctx, toolName, params)
        },
    ))

    // Handle prompt actions
    s.uiRouter.HandleType(mcpui.ActionTypePrompt, mcpui.WrapPromptHandler(
        func(ctx context.Context, promptName string, args map[string]any) (any, error) {
            // Execute the MCP prompt
            return s.executePrompt(ctx, promptName, args)
        },
    ))

    // Handle resource actions
    s.uiRouter.HandleType(mcpui.ActionTypeResource, mcpui.WrapResourceHandler(
        func(ctx context.Context, uri string) (any, error) {
            // Read the MCP resource
            return s.readResource(ctx, uri)
        },
    ))
}
```

### Processing Actions

When an action is received, route it and send the response:

```go
func (s *Server) handleUIAction(ctx context.Context, sourceURI string, action *mcpui.UIAction) {
    // Create request
    req := &mcpui.UIActionRequest{
        SourceURI: sourceURI,
        Action:    action,
    }

    // Send acknowledgment
    ack := mcpui.NewReceivedResponse(action.MessageID)
    s.sendUIResponse(ack)

    // Route to handler
    result, err := s.uiRouter.Route(ctx, req)

    // Send final response
    var resp *mcpui.UIResponse
    if err != nil {
        resp = mcpui.NewErrorResponse(action.MessageID, err)
    } else if result.Error != nil {
        resp = mcpui.NewErrorResponse(action.MessageID, result.Error)
    } else {
        resp = mcpui.NewSuccessResponse(action.MessageID, result.Response)
    }
    s.sendUIResponse(resp)
}
```

## Complete Integration Example

Here's a complete example showing MCP server integration:

```go
package main

import (
    "context"
    "encoding/json"
    "log"

    "github.com/modelcontextprotocol/go-sdk/mcp"
    "github.com/modelcontextprotocol/go-sdk/server"
    "github.com/ironystock/mcpui-go"
)

type MyServer struct {
    *server.MCPServer
    uiRouter *mcpui.Router
}

func main() {
    s := &MyServer{
        uiRouter: mcpui.NewRouter(),
    }

    // Setup MCP server
    s.MCPServer = server.NewMCPServer("my-server", "1.0.0",
        server.WithToolHandler(s.handleTool),
    )

    // Setup UI handlers
    s.setupUIHandlers()

    // Register tools
    s.MCPServer.AddTool(mcp.Tool{
        Name:        "show_dashboard",
        Description: "Display the main dashboard",
    })

    // Run server
    if err := s.MCPServer.ServeStdio(); err != nil {
        log.Fatal(err)
    }
}

func (s *MyServer) handleTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    switch req.Params.Name {
    case "show_dashboard":
        return s.showDashboard(ctx, req)
    default:
        return nil, fmt.Errorf("unknown tool: %s", req.Params.Name)
    }
}

func (s *MyServer) showDashboard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    content := &mcpui.HTMLContent{
        HTML: `<div id="dashboard">
            <h1>Dashboard</h1>
            <button onclick="sendAction('refresh')">Refresh</button>
        </div>`,
    }

    rc, _ := mcpui.NewUIResourceContents("ui://dashboard/main", content)

    return &mcp.CallToolResult{
        Content: []mcp.Content{
            {Type: "resource", Resource: &mcp.ResourceContents{
                URI: rc.URI, MimeType: rc.MimeType, Text: rc.Text,
            }},
        },
    }, nil
}

func (s *MyServer) setupUIHandlers() {
    s.uiRouter.HandleType(mcpui.ActionTypeTool, mcpui.WrapToolHandler(
        func(ctx context.Context, toolName string, params map[string]any) (any, error) {
            req := mcp.CallToolRequest{
                Params: mcp.CallToolRequestParams{
                    Name:      toolName,
                    Arguments: params,
                },
            }
            return s.handleTool(ctx, req)
        },
    ))

    s.uiRouter.HandleType(mcpui.ActionTypeCustom, mcpui.WrapCustomHandler(
        func(ctx context.Context, action string, data map[string]any) (any, error) {
            switch action {
            case "refresh":
                return map[string]string{"status": "refreshed"}, nil
            default:
                return nil, fmt.Errorf("unknown action: %s", action)
            }
        },
    ))
}
```

## Best Practices

1. **Separate concerns** - Keep UI resource generation separate from business logic
2. **Validate content** - Use `ValidateContent` before creating resources
3. **Handle errors gracefully** - Always send error responses for failed actions
4. **Use acknowledgments** - Send `received` responses for long-running operations
5. **Keep HTML self-contained** - Include all styles inline for HTMLContent
6. **Cache templates** - Pre-compile HTML templates for performance
7. **Test handlers** - Unit test action handlers independently
