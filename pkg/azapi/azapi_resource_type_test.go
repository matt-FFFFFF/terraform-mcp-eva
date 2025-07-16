package azapi

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestGetAzAPIType_WithoutJsonPath(t *testing.T) {
	schema, err := GetResourceJsonSchema("Microsoft.Resources/resourcegroups", "2024-07-01", "")
	require.NoError(t, err)
	assert.NotEmpty(t, schema)
	v := make(map[string]any)
	require.NoError(t, json.Unmarshal([]byte(schema), &v))
	require.Contains(t, v, "attributes")
	attributes, ok := v["attributes"].(map[string]any)
	require.True(t, ok)
	require.Contains(t, attributes, "location")
	location, ok := attributes["location"].(map[string]any)
	require.True(t, ok)
	require.Contains(t, location, "type")
	assert.Equal(t, "string", location["type"])
}

func TestGetAzAPIType_WithJsonPath(t *testing.T) {
	cases := []struct {
		desc         string
		resourceType string
		apiVersion   string
		path         string
		expectedType string
	}{
		{
			resourceType: "Microsoft.MachineLearningServices/workspaces",
			apiVersion:   "2025-06-01",
			path:         "body.properties.encryption.identity.userAssignedIdentity",
			expectedType: `String`,
		},
		{
			resourceType: "Microsoft.MachineLearningServices/workspaces",
			apiVersion:   "2025-06-01",
			path:         "body.properties.publicNetworkAccess",
			expectedType: `String`,
		},
		{
			desc:         "path in array",
			resourceType: "Microsoft.Compute/virtualMachines",
			apiVersion:   "2024-11-01",
			path:         "body.properties.osProfile.secrets.sourceVault.id",
			expectedType: `String`,
		},
		{
			desc:         "object",
			resourceType: "Microsoft.Compute/virtualMachines",
			apiVersion:   "2024-11-01",
			path:         "body.properties.osProfile.secrets.sourceVault",
			expectedType: `ObjectWithOptionalAttrs(map[string]Type{"id":String}, []string{"id"})`,
		},
	}
	for _, c := range cases {
		caseName := c.desc
		if c.desc == "" {
			caseName = strings.Join([]string{c.resourceType, c.apiVersion, c.path}, "-")
		}
		t.Run(caseName, func(t *testing.T) {
			schema, err := GetResourceJsonSchema(c.resourceType, c.apiVersion, c.path)
			require.NoError(t, err)
			assert.Equal(t, c.expectedType, schema)
		})
	}
}
