// Example: router
//
// This example demonstrates the Router for handling different action types
// and resource-specific handlers.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ironystock/mcpui-go"
)

func main() {
	// Create a new router
	router := mcpui.NewRouter()

	// Register a handler for tool actions
	router.HandleType(mcpui.ActionTypeTool, mcpui.WrapToolHandler(
		func(ctx context.Context, toolName string, params map[string]any) (any, error) {
			fmt.Printf("Tool called: %s with params: %v\n", toolName, params)
			return map[string]any{
				"tool":   toolName,
				"status": "executed",
			}, nil
		},
	))

	// Register a handler for prompt actions
	router.HandleType(mcpui.ActionTypePrompt, mcpui.WrapPromptHandler(
		func(ctx context.Context, prompt string) (any, error) {
			fmt.Printf("Prompt received: %s\n", prompt)
			return map[string]any{
				"prompt": prompt,
				"status": "processed",
			}, nil
		},
	))

	// Register a resource-specific handler (overrides type handlers)
	router.HandleResource("ui://audio/mixer",
		func(ctx context.Context, req *mcpui.UIActionRequest) (*mcpui.UIActionResult, error) {
			fmt.Println("Audio mixer specific handler called")
			return &mcpui.UIActionResult{
				Response: map[string]string{"handler": "audio-mixer-specific"},
			}, nil
		},
	)

	// Set a default handler for unmatched actions
	router.SetDefault(func(ctx context.Context, req *mcpui.UIActionRequest) (*mcpui.UIActionResult, error) {
		fmt.Printf("Default handler called for action type: %s\n", req.Action.Type)
		return &mcpui.UIActionResult{
			Response: map[string]string{"handler": "default"},
		}, nil
	})

	// Test routing different actions
	ctx := context.Background()

	// Test 1: Tool action
	fmt.Println("\n--- Test 1: Tool Action ---")
	result1, err := router.Dispatch(ctx, &mcpui.UIActionRequest{
		ResourceURI: "ui://dashboard/main",
		Action: &mcpui.UIAction{
			MessageID: "msg-1",
			Type:      mcpui.ActionTypeTool,
			Payload:   json.RawMessage(`{"name":"get_status","parameters":{}}`),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	printResult("Tool action", result1)

	// Test 2: Prompt action
	fmt.Println("\n--- Test 2: Prompt Action ---")
	result2, err := router.Dispatch(ctx, &mcpui.UIActionRequest{
		ResourceURI: "ui://dashboard/main",
		Action: &mcpui.UIAction{
			MessageID: "msg-2",
			Type:      mcpui.ActionTypePrompt,
			Payload:   json.RawMessage(`{"prompt":"What is the current status?"}`),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	printResult("Prompt action", result2)

	// Test 3: Resource-specific handler
	fmt.Println("\n--- Test 3: Resource-Specific Handler ---")
	result3, err := router.Dispatch(ctx, &mcpui.UIActionRequest{
		ResourceURI: "ui://audio/mixer",
		Action: &mcpui.UIAction{
			MessageID: "msg-3",
			Type:      mcpui.ActionTypeTool,
			Payload:   json.RawMessage(`{"name":"toggle_mute","parameters":{}}`),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	printResult("Resource-specific", result3)

	// Test 4: Default handler (intent action not registered)
	fmt.Println("\n--- Test 4: Default Handler ---")
	result4, err := router.Dispatch(ctx, &mcpui.UIActionRequest{
		ResourceURI: "ui://unknown/resource",
		Action: &mcpui.UIAction{
			MessageID: "msg-4",
			Type:      mcpui.ActionTypeIntent,
			Payload:   json.RawMessage(`{"intent":"unknown"}`),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	printResult("Default handler", result4)
}

func printResult(name string, result *mcpui.UIActionResult) {
	data, _ := json.MarshalIndent(result.Response, "", "  ")
	fmt.Printf("Result from %s:\n%s\n", name, string(data))
}
