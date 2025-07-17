package azapi

import (
	"fmt"
	"github.com/ms-henglu/go-azure-types/types"
)

func GetApiVersions(resourceType string) ([]string, error) {
	versions := types.DefaultAzureSchemaLoader().ListApiVersions(resourceType)
	if len(versions) == 0 {
		return nil, fmt.Errorf("no API versions found for resource type %s", resourceType)
	}
	return versions, nil
}
