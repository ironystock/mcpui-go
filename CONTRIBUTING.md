# Contributing to mcpui-go

Thank you for your interest in contributing to mcpui-go!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/mcpui-go.git`
3. Create a branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Push and create a Pull Request

## Development Setup

### Prerequisites

- Go 1.25.5 or later
- Git

### Building

```bash
go build ./...
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run with verbose output
go test ./... -v
```

### Code Quality

```bash
# Format code
gofmt -w .

# Run static analysis
go vet ./...
```

## Code Style

- Follow standard Go conventions (gofmt, go vet)
- Use meaningful variable and function names
- Keep functions focused and concise
- Add comments for exported functions and types
- Include examples in documentation

### Example Documentation

For exported functions, include an example:

```go
// NewRouter creates a new action router.
//
// Example:
//
//	router := mcpui.NewRouter()
//	router.HandleType(mcpui.ActionTypeTool, toolHandler)
func NewRouter() *Router {
    // ...
}
```

## Testing Requirements

- All new features must have tests
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Test edge cases and error conditions

### Test Example

```go
func TestNewUIResourceContents(t *testing.T) {
    tests := []struct {
        name    string
        uri     string
        content UIContent
        wantErr bool
    }{
        {
            name: "valid HTML content",
            uri:  "ui://test/hello",
            content: &HTMLContent{HTML: "<div>Hello</div>"},
            wantErr: false,
        },
        {
            name:    "nil content",
            uri:     "ui://test/nil",
            content: nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewUIResourceContents(tt.uri, tt.content)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewUIResourceContents() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Pull Request Process

1. **Title**: Use a clear, descriptive title
   - Good: "Add validation for Remote DOM framework types"
   - Bad: "Fix bug"

2. **Description**: Explain what and why
   - What problem does this solve?
   - How does it solve it?
   - Are there any breaking changes?

3. **Tests**: Ensure all tests pass

4. **Review**: Address any feedback from reviewers

## Commit Messages

Follow conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Formatting, no code change
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance

Examples:
```
feat(router): add pattern matching for resource URIs

fix(content): handle empty HTML validation correctly

docs(readme): add integration example
```

## Reporting Issues

When reporting issues, please include:

1. Go version (`go version`)
2. Operating system
3. Steps to reproduce
4. Expected behavior
5. Actual behavior
6. Any relevant code or error messages

## Questions?

If you have questions, feel free to:

1. Open an issue with the "question" label
2. Check existing issues for similar questions
3. Review the documentation

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
