package azapi

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetApiVersions(t *testing.T) {
	resourceType := "Microsoft.Compute/virtualMachines"
	versions, err := GetApiVersions(resourceType)
	require.NoError(t, err)
	assert.Contains(t, versions, "2024-11-01", "Expected API version 2024-11-01 to be in the list of versions for %s", resourceType)
}
