// Example: basic
//
// This example demonstrates the simplest use case: creating HTML content
// and serializing it as an MCP resource.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ironystock/mcpui-go"
)

func main() {
	// Create HTML content for a simple greeting card
	content := &mcpui.HTMLContent{
		HTML: `<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            margin: 0;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .card {
            background: white;
            border-radius: 12px;
            padding: 40px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            text-align: center;
        }
        h1 {
            color: #333;
            margin: 0 0 10px 0;
        }
        p {
            color: #666;
            margin: 0;
        }
    </style>
</head>
<body>
    <div class="card">
        <h1>Hello, World!</h1>
        <p>This is rendered in a sandboxed iframe via MCP-UI.</p>
    </div>
</body>
</html>`,
	}

	// Create resource contents for MCP response
	rc, err := mcpui.NewUIResourceContents("ui://greeting/hello", content)
	if err != nil {
		log.Fatal(err)
	}

	// Print the JSON output (this would be returned in an MCP tool response)
	data, err := json.MarshalIndent(rc, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UI Resource Contents:")
	fmt.Println(string(data))
}
