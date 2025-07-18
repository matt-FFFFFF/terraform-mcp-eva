package tool

import (
	"context"
	"errors"
	"fmt"

	"github.com/lonegunmanb/terraform-mcp-eva/pkg/azapi"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AzAPIResourceSchemaQueryParam struct {
	ResourceType string `json:"resource_type" jsonschema:"Azure resource type, for example: Microsoft.Compute/virtualMachines, combined with api_version to identify the resource schema, like: Microsoft.Compute/virtualMachines@2024-11-01"`
	ApiVersion   string `json:"api_version" jsonschema:"Azure resource api-version, for example: 2024-11-01, combined with resource_type to identify the resource schema, like: Microsoft.Compute/virtualMachines@2024-11-01"`
	Path         string `json:"path,omitempty" jsonschema:"JSON path to query the resource schema, for example: body.properties.osProfile.secrets.sourceVault.id, if not specified, the whole resource schema will be returned"`
}

func QueryAzAPIResourceSchema(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[AzAPIResourceSchemaQueryParam]) (*mcp.CallToolResultFor[any], error) {
	resourceType := params.Arguments.ResourceType
	apiVersion := params.Arguments.ApiVersion
	if resourceType == "" || apiVersion == "" {
		return nil, errors.New("`resource_type` and `api_version` are required parameters")
	}
	path := params.Arguments.Path
	schema, err := azapi.GetResourceSchema(resourceType, apiVersion, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource schema for %s@%s: %w", resourceType, apiVersion, err)
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: schema,
			},
		},
	}, nil
}
