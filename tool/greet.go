package tool

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GreetingInput struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type GreetingOutput struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input GreetingInput) (
	*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Hi " + input.Name},
		},
	}, GreetingOutput{Greeting: "Hi " + input.Name}, nil
}
