// Example: remote-dom
//
// This example demonstrates creating Remote DOM content for dynamic,
// framework-based UIs using React.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ironystock/mcpui-go"
)

func main() {
	// Create Remote DOM content with React
	content := &mcpui.RemoteDOMContent{
		Script: `
// Counter component using React hooks
const [count, setCount] = React.useState(0);
const [theme, setTheme] = React.useState('light');

const styles = {
    container: {
        fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
        padding: '20px',
        backgroundColor: theme === 'light' ? '#f5f5f5' : '#1a1a2e',
        color: theme === 'light' ? '#333' : '#eee',
        minHeight: '100vh',
        transition: 'all 0.3s ease',
    },
    card: {
        backgroundColor: theme === 'light' ? 'white' : '#16213e',
        borderRadius: '12px',
        padding: '30px',
        maxWidth: '400px',
        margin: '0 auto',
        boxShadow: theme === 'light'
            ? '0 4px 20px rgba(0,0,0,0.1)'
            : '0 4px 20px rgba(0,0,0,0.4)',
    },
    count: {
        fontSize: '48px',
        fontWeight: 'bold',
        textAlign: 'center',
        margin: '20px 0',
        color: theme === 'light' ? '#667eea' : '#a8d8ea',
    },
    buttonRow: {
        display: 'flex',
        gap: '10px',
        justifyContent: 'center',
        marginBottom: '20px',
    },
    button: {
        padding: '12px 24px',
        fontSize: '16px',
        border: 'none',
        borderRadius: '8px',
        cursor: 'pointer',
        backgroundColor: '#667eea',
        color: 'white',
        transition: 'transform 0.1s, background-color 0.2s',
    },
    themeButton: {
        padding: '8px 16px',
        fontSize: '14px',
        border: '1px solid',
        borderColor: theme === 'light' ? '#ddd' : '#444',
        borderRadius: '6px',
        cursor: 'pointer',
        backgroundColor: 'transparent',
        color: theme === 'light' ? '#666' : '#aaa',
    },
};

return React.createElement('div', { style: styles.container },
    React.createElement('div', { style: styles.card },
        React.createElement('h2', { style: { textAlign: 'center', margin: '0 0 10px 0' } },
            'React Counter'
        ),
        React.createElement('p', { style: { textAlign: 'center', color: '#888', margin: '0 0 20px 0' } },
            'Rendered via Remote DOM'
        ),
        React.createElement('div', { style: styles.count }, count),
        React.createElement('div', { style: styles.buttonRow },
            React.createElement('button', {
                style: styles.button,
                onClick: () => setCount(c => c - 1),
            }, '-'),
            React.createElement('button', {
                style: styles.button,
                onClick: () => setCount(0),
            }, 'Reset'),
            React.createElement('button', {
                style: styles.button,
                onClick: () => setCount(c => c + 1),
            }, '+')
        ),
        React.createElement('div', { style: { textAlign: 'center' } },
            React.createElement('button', {
                style: styles.themeButton,
                onClick: () => setTheme(t => t === 'light' ? 'dark' : 'light'),
            }, theme === 'light' ? 'Dark Mode' : 'Light Mode')
        )
    )
);
`,
		Framework: mcpui.FrameworkReact,
	}

	// Create resource contents
	rc, err := mcpui.NewUIResourceContents("ui://demo/counter", content)
	if err != nil {
		log.Fatal(err)
	}

	// Print the result
	fmt.Println("Remote DOM Resource Contents:")
	fmt.Printf("URI: %s\n", rc.URI)
	fmt.Printf("MIMEType: %s\n", rc.MIMEType)
	fmt.Println("\nParsed content:")

	// Parse and pretty-print the text content
	var textContent map[string]any
	if err := json.Unmarshal([]byte(rc.Text), &textContent); err != nil {
		log.Fatal(err)
	}

	data, _ := json.MarshalIndent(textContent, "", "  ")
	fmt.Println(string(data))

	// Also show Web Components framework
	fmt.Println("\n--- Web Components Framework Example ---")

	wcContent := &mcpui.RemoteDOMContent{
		Script: `
class MyGreeting extends HTMLElement {
    connectedCallback() {
        this.innerHTML = '<h1>Hello from Web Components!</h1>';
    }
}
customElements.define('my-greeting', MyGreeting);
`,
		Framework: mcpui.FrameworkWebComponents,
	}
	fmt.Printf("WebComponents framework: %s\n", wcContent.Framework)
}
