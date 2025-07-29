package tool

import (
	"context"
	"fmt"

	"github.com/matt-FFFFFF/tfpluginschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	blockTypeResource  = "resource"
	blockTypeData      = "data"
	blockTypeEphemeral = "ephemeral"
)

type FineGrainedSchemaQueryParam struct {
	BlockType         string `json:"block_type" jsonschema:"Terraform block type, possible values: resource, data, ephemeral"`
	ProviderName      string `json:"provider_name" jsonschema:"The name of the provider: azapi, azurerm, etc. This is the first segment of the block label, e.g. for azurerm_virtual_machine it's azurerm."`
	ProviderNamespace string `json:"provider_namespace" jsonschema:"The namespace of the provider, e.g. Azure, hashicorp, etc. This can be found in the terraform.required_providers block."`
	ProviderVersion   string `json:"provider_version" jsonschema:"The version of the provider, e.g. 2.5.0. Can be obtained by running 'terraform providers' command."`
	BlockLabel        string `json:"block_label" jsonschema:"The first label of the block, e.g. azurerm_virtual_machine"`
}

var validCategories = map[string]struct{}{
	blockTypeResource:  {},
	blockTypeData:      {},
	blockTypeEphemeral: {},
}

func QueryResourceSchema(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[FineGrainedSchemaQueryParam]) (*mcp.CallToolResultFor[any], error) {
	if _, ok := validCategories[params.Arguments.BlockType]; !ok {
		return nil, fmt.Errorf("invalid category: %s", params.Arguments.BlockType)
	}

	server, ok := ctx.Value(tfpluginschema.ContextKey{}).(*tfpluginschema.Server)
	if !ok {
		return nil, fmt.Errorf("failed to get tfpluginschema server from context")
	}

	req := tfpluginschema.Request{
		Namespace: params.Arguments.ProviderNamespace,
		Version:   params.Arguments.ProviderVersion,
		Name:      params.Arguments.ProviderName,
	}
	if err := server.Get(req); err != nil {
		return nil, fmt.Errorf("failed to get provider %s: %w", req.Name, err)
	}

	var err error
	var returnData []byte
	switch params.Arguments.BlockType {
	case blockTypeResource:
		returnData, err = server.GetResourceSchema(req, params.Arguments.BlockLabel)
	case blockTypeData:
		returnData, err = server.GetDataSourceSchema(req, params.Arguments.BlockLabel)
	case blockTypeEphemeral:
		returnData, err = server.GetEphemeralResourceSchema(req, params.Arguments.BlockLabel)
	}

	if err != nil || len(returnData) == 0 {
		return nil, fmt.Errorf("failed to get schema for %s %s: %w", params.Arguments.BlockType, params.Arguments.BlockLabel, err)
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(returnData),
				Annotations: &mcp.Annotations{
					Audience: []mcp.Role{
						"assistant",
					},
				},
			},
		},
	}, nil
}
