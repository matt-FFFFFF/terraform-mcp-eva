package pkg

import (
	"github.com/lonegunmanb/terraform-mcp-eva/pkg/prompt"
	"github.com/lonegunmanb/terraform-mcp-eva/pkg/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterMcpServer(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: p(false),
			IdempotentHint:  true,
			OpenWorldHint:   p(false),
			ReadOnlyHint:    true,
		},
		Description: "Query fine grained AzAPI resource body schema by `resource type`, `api_version` and optional `path`. The returned type is a Go type string, which can be used in Go code to represent the resource's `body` attribute. If you're querying corresponds to the AzAPI provider and the `body` attribute, this tool should have higher priority",
		Name:        "query_azapi_resource_body",
	}, tool.QueryAzAPIResourceSchema)

	mcp.AddTool(s, &mcp.Tool{
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: p(false),
			IdempotentHint:  true,
			OpenWorldHint:   p(false),
			ReadOnlyHint:    true,
		},
		Description: "Query Azure API versions by `resource type`, e.g. `Microsoft.Compute/virtualMachines`. The returned value is a list of API versions for the specified resource type, split by comma.",
		Name:        "list_azapi_api_versions",
	}, tool.QueryAzAPIVersions)

	mcp.AddTool(s, &mcp.Tool{
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: p(false),
			IdempotentHint:  true,
			OpenWorldHint:   p(false),
			ReadOnlyHint:    true,
		},
		Description: "Query fine grained AzAPI resource description by `resource type`, `api_version` and optional `path`. The returned value is either description of the property, or json object representing the object, the key is property name the value is the description of the property. Via description you can learn whether a property is id, readonly or writeonly, and possible values. If you're querying AzAPI provider and the `body` attribute, this tool should have higher priority",
		Name:        "query_azapi_resource_document",
	}, tool.QueryAzAPIDescriptionSchema)

	mcp.AddTool(s, &mcp.Tool{
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: p(false),
			IdempotentHint:  true,
			OpenWorldHint:   p(true),
			ReadOnlyHint:    true,
			Title:           "Query Terraform Provider Schema",
		},
		Description: "Query Terraform provider schemas by name. Supports resource, ephemeral and data blocks. MUST supply provider name, e.g. azurerm, provider version, e.g. 2.5.0, and the first block label. MUST get provider version from `terraform providers`, The returned value is a JSON string representing the resource schema, including attribute descriptions. If you're querying schema information about specified attribute or nested block schema, this tool should have higher priority.",
		Name:        "query_terraform_provider_schema",
	}, tool.QueryResourceSchema)
	prompt.AddSolveAvmIssuePrompt(s)
}

func p[T any](input T) *T {
	return &input
}
