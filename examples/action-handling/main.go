// Example: action-handling
//
// This example demonstrates the complete action->response flow,
// including parsing payloads, processing actions, and sending responses.
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ironystock/mcpui-go"
)

// Simulated state
var state = struct {
	Volume float64
	Muted  bool
}{
	Volume: 0.75,
	Muted:  false,
}

func main() {
	ctx := context.Background()

	// Simulate receiving different actions and handling them
	fmt.Println("=== Action Handling Example ===")

	// Example 1: Tool action - toggle mute
	handleToolAction(ctx)

	// Example 2: Notify action
	handleNotifyAction(ctx)

	// Example 3: Full request/response cycle
	handleFullCycle(ctx)
}

func handleToolAction(ctx context.Context) {
	fmt.Println("--- Example 1: Tool Action (Toggle Mute) ---")

	// Simulate receiving an action from UI
	action := &mcpui.UIAction{
		MessageID: "msg-001",
		Type:      mcpui.ActionTypeTool,
		Payload:   json.RawMessage(`{"name":"toggle_mute","parameters":{"inputName":"Microphone"}}`),
	}

	// Parse the tool payload
	payload, err := action.ToolPayload()
	if err != nil {
		sendError(action.MessageID, err)
		return
	}

	fmt.Printf("Tool: %s\n", payload.ToolName)
	fmt.Printf("Parameters: %v\n", payload.Params)

	// Process the action
	state.Muted = !state.Muted
	result := map[string]any{
		"inputName": payload.Params["inputName"],
		"muted":     state.Muted,
	}

	// Send success response
	resp := mcpui.NewSuccessResponse(action.MessageID, result)
	printResponse(resp)
	fmt.Println()
}

func handleNotifyAction(ctx context.Context) {
	fmt.Println("--- Example 2: Notify Action ---")

	// Create a notify action
	action, err := mcpui.NewNotifyAction("Recording started", "info")
	if err != nil {
		fmt.Printf("Error creating action: %v\n", err)
		return
	}

	// Parse notify payload
	payload, err := action.NotifyPayload()
	if err != nil {
		fmt.Printf("Error parsing payload: %v\n", err)
		return
	}

	fmt.Printf("Notification: %s (level: %s)\n", payload.Message, payload.Level)
	fmt.Println()
}

func handleFullCycle(ctx context.Context) {
	fmt.Println("--- Example 3: Full Request/Response Cycle ---")

	// Simulate an action request
	request := &mcpui.UIActionRequest{
		ResourceURI: "ui://audio/mixer",
		Action: &mcpui.UIAction{
			MessageID: "msg-004",
			Type:      mcpui.ActionTypeTool,
			Payload:   json.RawMessage(`{"name":"set_volume","parameters":{"inputName":"Music","volume":0.5}}`),
		},
	}

	fmt.Printf("Received action from: %s\n", request.ResourceURI)

	// Step 1: Send acknowledgment
	ack := mcpui.NewReceivedResponse(request.Action.MessageID)
	fmt.Println("\nStep 1: Send acknowledgment")
	printResponse(ack)

	// Step 2: Process the action
	fmt.Println("\nStep 2: Process action...")
	payload, _ := request.Action.ToolPayload()

	var finalResp *mcpui.UIResponse

	switch payload.ToolName {
	case "set_volume":
		if vol, ok := payload.Params["volume"].(float64); ok {
			state.Volume = vol
			finalResp = mcpui.NewSuccessResponse(request.Action.MessageID, map[string]any{
				"inputName": payload.Params["inputName"],
				"volume":    state.Volume,
			})
		} else {
			finalResp = mcpui.NewErrorResponseWithCode(
				request.Action.MessageID,
				"invalid_parameter",
				"volume must be a number",
			)
		}
	default:
		finalResp = mcpui.NewErrorResponseWithCode(
			request.Action.MessageID,
			"unknown_tool",
			fmt.Sprintf("unknown tool: %s", payload.ToolName),
		)
	}

	// Step 3: Send final response
	fmt.Println("\nStep 3: Send final response")
	printResponse(finalResp)
}

func sendError(messageID string, err error) {
	resp := mcpui.NewErrorResponse(messageID, err)
	printResponse(resp)
}

func printResponse(resp *mcpui.UIResponse) {
	data, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(data))
}
