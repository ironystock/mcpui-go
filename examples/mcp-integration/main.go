// Example: mcp-integration
//
// This example demonstrates how to integrate mcpui-go with an MCP server.
// It shows patterns for returning UI resources from tool handlers and
// processing UI actions.
//
// Note: This is a demonstration of the integration patterns. For a complete
// working example, you would need the MCP Go SDK installed.
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ironystock/mcpui-go"
)

// Server represents an MCP server with UI capabilities
type Server struct {
	router *mcpui.Router
	state  *ServerState
}

// ServerState holds the server's current state
type ServerState struct {
	Connected bool
	Recording bool
	Streaming bool
	Volume    float64
}

func main() {
	// Create server with initial state
	server := &Server{
		router: mcpui.NewRouter(),
		state: &ServerState{
			Connected: true,
			Recording: false,
			Streaming: false,
			Volume:    0.75,
		},
	}

	// Setup UI action handlers
	server.setupUIHandlers()

	fmt.Println("=== MCP Integration Example ===")

	// Demonstrate tool response with UI resource
	fmt.Println("--- Tool Response with UI Resource ---")
	server.demonstrateToolResponse()

	// Demonstrate UI action handling
	fmt.Println("\n--- UI Action Handling ---")
	server.demonstrateActionHandling()

	// Demonstrate multiple resources
	fmt.Println("\n--- Multiple Resources ---")
	server.demonstrateMultipleResources()
}

func (s *Server) setupUIHandlers() {
	// Handle tool actions from UI
	s.router.HandleType(mcpui.ActionTypeTool, mcpui.WrapToolHandler(
		func(ctx context.Context, toolName string, params map[string]any) (any, error) {
			return s.executeTool(toolName, params)
		},
	))

	// Handle intent actions
	s.router.HandleType(mcpui.ActionTypeIntent, mcpui.WrapIntentHandler(
		func(ctx context.Context, intent string, params map[string]any) (any, error) {
			switch intent {
			case "refresh":
				return map[string]any{"refreshed": true, "state": s.state}, nil
			case "toggle_recording":
				s.state.Recording = !s.state.Recording
				return map[string]any{"recording": s.state.Recording}, nil
			default:
				return nil, fmt.Errorf("unknown intent: %s", intent)
			}
		},
	))

	// Handle dashboard-specific actions
	s.router.HandleResource("ui://dashboard/main",
		func(ctx context.Context, req *mcpui.UIActionRequest) (*mcpui.UIActionResult, error) {
			fmt.Println("  Dashboard-specific handler invoked")
			return &mcpui.UIActionResult{
				Response: map[string]any{
					"handler":     "dashboard",
					"resourceURI": req.ResourceURI,
				},
			}, nil
		},
	)
}

func (s *Server) executeTool(name string, params map[string]any) (any, error) {
	switch name {
	case "get_status":
		return s.state, nil
	case "start_recording":
		s.state.Recording = true
		return map[string]any{"recording": true}, nil
	case "stop_recording":
		s.state.Recording = false
		return map[string]any{"recording": false}, nil
	case "set_volume":
		if vol, ok := params["volume"].(float64); ok {
			s.state.Volume = vol
			return map[string]any{"volume": s.state.Volume}, nil
		}
		return nil, fmt.Errorf("invalid volume parameter")
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func (s *Server) demonstrateToolResponse() {
	// This shows how a tool handler would return a UI resource
	// In a real MCP server, this would be returned in CallToolResult.Content

	// Generate dynamic dashboard HTML based on state
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: sans-serif; padding: 20px; background: #f5f5f5; }
        .dashboard { background: white; border-radius: 8px; padding: 20px; max-width: 400px; margin: 0 auto; }
        .status { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid #eee; }
        .status:last-child { border-bottom: none; }
        .label { color: #666; }
        .value { font-weight: bold; }
        .value.active { color: #4caf50; }
        .value.inactive { color: #999; }
    </style>
</head>
<body>
    <div class="dashboard">
        <h2>System Status</h2>
        <div class="status">
            <span class="label">Connected</span>
            <span class="value %s">%v</span>
        </div>
        <div class="status">
            <span class="label">Recording</span>
            <span class="value %s">%v</span>
        </div>
        <div class="status">
            <span class="label">Streaming</span>
            <span class="value %s">%v</span>
        </div>
        <div class="status">
            <span class="label">Volume</span>
            <span class="value">%.0f%%</span>
        </div>
    </div>
</body>
</html>`,
		statusClass(s.state.Connected), s.state.Connected,
		statusClass(s.state.Recording), s.state.Recording,
		statusClass(s.state.Streaming), s.state.Streaming,
		s.state.Volume*100,
	)

	content := &mcpui.HTMLContent{HTML: html}
	rc, _ := mcpui.NewUIResourceContents("ui://dashboard/status", content)

	// This would be embedded in MCP CallToolResult.Content
	fmt.Println("Generated UI Resource:")
	fmt.Printf("  URI: %s\n", rc.URI)
	fmt.Printf("  MIMEType: %s\n", rc.MIMEType)
	fmt.Printf("  Content length: %d bytes\n", len(rc.Text))
}

func (s *Server) demonstrateActionHandling() {
	ctx := context.Background()

	// Simulate receiving a UI action
	action := &mcpui.UIAction{
		MessageID: "ui-msg-001",
		Type:      mcpui.ActionTypeTool,
		Payload:   json.RawMessage(`{"name":"start_recording","parameters":{}}`),
	}

	request := &mcpui.UIActionRequest{
		ResourceURI: "ui://dashboard/status",
		Action:      action,
	}

	fmt.Printf("Received action from %s\n", request.ResourceURI)
	fmt.Printf("  Type: %s\n", action.Type)

	// Route the action
	result, err := s.router.Dispatch(ctx, request)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
		return
	}

	// Send response
	var resp *mcpui.UIResponse
	if result.Error != nil {
		resp = mcpui.NewErrorResponse(action.MessageID, result.Error)
	} else {
		resp = mcpui.NewSuccessResponse(action.MessageID, result.Response)
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("Response:\n%s\n", string(data))
}

func (s *Server) demonstrateMultipleResources() {
	// Show how to return multiple UI resources

	resources := []struct {
		uri     string
		content mcpui.UIContent
	}{
		{
			uri: "ui://panel/status",
			content: &mcpui.HTMLContent{
				HTML: "<div class='panel'><h3>Status</h3><p>OK</p></div>",
			},
		},
		{
			uri: "ui://panel/controls",
			content: &mcpui.HTMLContent{
				HTML: "<div class='panel'><h3>Controls</h3><button>Start</button></div>",
			},
		},
		{
			uri: "ui://panel/metrics",
			content: &mcpui.HTMLContent{
				HTML: "<div class='panel'><h3>Metrics</h3><p>CPU: 45%</p></div>",
			},
		},
	}

	fmt.Println("Generated multiple UI resources:")
	for _, r := range resources {
		rc, _ := mcpui.NewUIResourceContents(r.uri, r.content)
		fmt.Printf("  - %s (%s)\n", rc.URI, rc.MIMEType)
	}

	fmt.Println("\nThese would be returned in CallToolResult.Content as an array")
}

func statusClass(active bool) string {
	if active {
		return "active"
	}
	return "inactive"
}
