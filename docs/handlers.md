# Handlers

The mcpui-go SDK provides a flexible system for routing and handling UI actions.

## UIActionHandler

The function signature for action handlers.

### Definition

```go
type UIActionHandler func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error)
```

### Example

```go
handler := func(ctx context.Context, req *mcpui.UIActionRequest) (*mcpui.UIActionResult, error) {
    // Process the action
    result := processAction(req.Action)

    return &mcpui.UIActionResult{
        Response: result,
    }, nil
}
```

## Router

Routes actions to appropriate handlers based on type or resource URI.

### Definition

```go
type Router struct {
    typeHandlers     map[ActionType]UIActionHandler
    resourceHandlers map[string]UIActionHandler
    defaultHandler   UIActionHandler
}
```

### Creating a Router

```go
router := mcpui.NewRouter()
```

### Registering Handlers

#### By Action Type

```go
router.HandleType(mcpui.ActionTypeTool, toolHandler)
router.HandleType(mcpui.ActionTypePrompt, promptHandler)
router.HandleType(mcpui.ActionTypeResource, resourceHandler)
router.HandleType(mcpui.ActionTypeCustom, customHandler)
```

#### By Resource URI

```go
router.HandleResource("ui://audio/mixer", audioMixerHandler)
router.HandleResource("ui://scene/preview", scenePreviewHandler)
router.HandleResource("ui://dashboard/main", dashboardHandler)
```

#### Default Handler

```go
router.SetDefault(func(ctx context.Context, req *mcpui.UIActionRequest) (*mcpui.UIActionResult, error) {
    return nil, errors.New("unhandled action")
})
```

### Routing Actions

```go
result, err := router.Route(ctx, &mcpui.UIActionRequest{
    SourceURI: "ui://audio/mixer",
    Action:    action,
})
```

### Routing Priority

1. **Resource handlers** - Exact URI match takes precedence
2. **Type handlers** - If no resource handler matches
3. **Default handler** - If no type handler matches
4. **Error** - If no default handler is set

## Typed Handler Wrappers

Convenience wrappers for type-specific handlers.

### WrapToolHandler

Wraps a function that handles tool actions.

```go
func WrapToolHandler(fn func(ctx context.Context, toolName string, params map[string]any) (any, error)) UIActionHandler
```

Example:
```go
router.HandleType(mcpui.ActionTypeTool, mcpui.WrapToolHandler(
    func(ctx context.Context, toolName string, params map[string]any) (any, error) {
        switch toolName {
        case "get_status":
            return getStatus()
        case "toggle_mute":
            inputName := params["inputName"].(string)
            return toggleMute(inputName)
        default:
            return nil, fmt.Errorf("unknown tool: %s", toolName)
        }
    },
))
```

### WrapPromptHandler

Wraps a function that handles prompt actions.

```go
func WrapPromptHandler(fn func(ctx context.Context, promptName string, args map[string]any) (any, error)) UIActionHandler
```

Example:
```go
router.HandleType(mcpui.ActionTypePrompt, mcpui.WrapPromptHandler(
    func(ctx context.Context, promptName string, args map[string]any) (any, error) {
        switch promptName {
        case "health-check":
            return runHealthCheck()
        case "audio-check":
            return runAudioCheck()
        default:
            return nil, fmt.Errorf("unknown prompt: %s", promptName)
        }
    },
))
```

### WrapResourceHandler

Wraps a function that handles resource actions.

```go
func WrapResourceHandler(fn func(ctx context.Context, uri string) (any, error)) UIActionHandler
```

Example:
```go
router.HandleType(mcpui.ActionTypeResource, mcpui.WrapResourceHandler(
    func(ctx context.Context, uri string) (any, error) {
        return readResource(uri)
    },
))
```

### WrapCustomHandler

Wraps a function that handles custom actions.

```go
func WrapCustomHandler(fn func(ctx context.Context, action string, data map[string]any) (any, error)) UIActionHandler
```

Example:
```go
router.HandleType(mcpui.ActionTypeCustom, mcpui.WrapCustomHandler(
    func(ctx context.Context, action string, data map[string]any) (any, error) {
        switch action {
        case "refresh":
            return handleRefresh(data)
        case "export":
            return handleExport(data)
        default:
            return nil, fmt.Errorf("unknown custom action: %s", action)
        }
    },
))
```

## Complete Example

```go
func setupRouter() *mcpui.Router {
    router := mcpui.NewRouter()

    // Tool handler
    router.HandleType(mcpui.ActionTypeTool, mcpui.WrapToolHandler(
        func(ctx context.Context, toolName string, params map[string]any) (any, error) {
            return executeTool(toolName, params)
        },
    ))

    // Prompt handler
    router.HandleType(mcpui.ActionTypePrompt, mcpui.WrapPromptHandler(
        func(ctx context.Context, promptName string, args map[string]any) (any, error) {
            return executePrompt(promptName, args)
        },
    ))

    // Resource-specific handlers
    router.HandleResource("ui://audio/mixer",
        func(ctx context.Context, req *mcpui.UIActionRequest) (*mcpui.UIActionResult, error) {
            return handleAudioMixer(ctx, req)
        },
    )

    // Default handler for unmatched actions
    router.SetDefault(func(ctx context.Context, req *mcpui.UIActionRequest) (*mcpui.UIActionResult, error) {
        return &mcpui.UIActionResult{
            Error: fmt.Errorf("no handler for action type %s", req.Action.Type),
        }, nil
    })

    return router
}
```

## Best Practices

1. **Use typed wrappers** - They handle payload parsing and provide type safety
2. **Register resource handlers for complex UIs** - Override type handlers for specific resources
3. **Always set a default handler** - Gracefully handle unknown actions
4. **Keep handlers focused** - Each handler should do one thing well
5. **Use context for cancellation** - Respect context cancellation in long-running handlers
