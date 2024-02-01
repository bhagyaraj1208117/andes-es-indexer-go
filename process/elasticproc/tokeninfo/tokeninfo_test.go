package tokeninfo

import (
	"testing"

	"github.com/bhagyaraj1208117/andes-core-go/core"
	"github.com/stretchr/testify/require"
)

func TestTokenRolesAndPropertiesAddRole(t *testing.T) {
	t.Parallel()

	tokenRolesAndProp := NewTokenRolesAndProperties()

	tokenRolesAndProp.AddRole("MY-abcd", "addr-1", core.DCTRoleNFTBurn, true)
	tokenRolesAndProp.AddRole("MY-abcd", "addr-2", core.DCTRoleNFTBurn, true)

	expected := map[string][]*RoleData{
		core.DCTRoleNFTBurn: {
			{
				Token:   "MY-abcd",
				Address: "addr-1",
				Set:     true,
			},
			{
				Token:   "MY-abcd",
				Address: "addr-2",
				Set:     true,
			},
		},
	}
	require.Equal(t, expected, tokenRolesAndProp.GetRoles())
}

func TestTokenAndROlesPropertiesAddProperties(t *testing.T) {
	t.Parallel()

	tokenRolesAndProp := NewTokenRolesAndProperties()

	properties1 := map[string]bool{
		"canDo":   true,
		"canBurn": false,
	}
	properties2 := map[string]bool{
		"canDo":   false,
		"canBurn": false,
	}
	tokenRolesAndProp.AddProperties("MY-aaaa", properties1)
	tokenRolesAndProp.AddProperties("MY-aaaa", properties2)

	expected := []*PropertiesData{
		{
			Token:      "MY-aaaa",
			Properties: properties1,
		},
		{
			Token:      "MY-aaaa",
			Properties: properties2,
		},
	}
	require.Equal(t, expected, tokenRolesAndProp.GetAllTokensWithProperties())
}
