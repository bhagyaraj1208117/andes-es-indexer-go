package logsevents

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"

	"github.com/bhagyaraj1208117/andes-core-go/core"
	"github.com/bhagyaraj1208117/andes-core-go/data/transaction"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/mock"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/process/elasticproc/tokeninfo"
	"github.com/stretchr/testify/require"
)

func TestDctPropertiesProcCreateRoleShouldWork(t *testing.T) {
	t.Parallel()

	dctPropProc := newDctPropertiesProcessor(&mock.PubkeyConverterMock{})

	event := &transaction.Event{
		Address:    []byte("addr"),
		Identifier: []byte(core.BuiltInFunctionSetDCTRole),
		Topics:     [][]byte{[]byte("MYTOKEN-abcd"), big.NewInt(0).Bytes(), big.NewInt(0).Bytes(), []byte(core.DCTRoleNFTCreate)},
	}

	tokenRolesAndProperties := tokeninfo.NewTokenRolesAndProperties()
	dctPropProc.processEvent(&argsProcessEvent{
		event:                   event,
		tokenRolesAndProperties: tokenRolesAndProperties,
	})

	expected := map[string][]*tokeninfo.RoleData{
		core.DCTRoleNFTCreate: {
			{
				Token:   "MYTOKEN-abcd",
				Set:     true,
				Address: "61646472",
			},
		},
	}
	require.Equal(t, expected, tokenRolesAndProperties.GetRoles())
}

func TestDctPropertiesProcTransferCreateRole(t *testing.T) {
	t.Parallel()

	dctPropProc := newDctPropertiesProcessor(&mock.PubkeyConverterMock{})

	event := &transaction.Event{
		Address:    []byte("addr"),
		Identifier: []byte(core.BuiltInFunctionDCTNFTCreateRoleTransfer),
		Topics:     [][]byte{[]byte("MYTOKEN-abcd"), big.NewInt(0).Bytes(), big.NewInt(0).Bytes(), []byte(strconv.FormatBool(true))},
	}

	tokenRolesAndProperties := tokeninfo.NewTokenRolesAndProperties()
	dctPropProc.processEvent(&argsProcessEvent{
		event:                   event,
		tokenRolesAndProperties: tokenRolesAndProperties,
	})

	expected := map[string][]*tokeninfo.RoleData{
		core.DCTRoleNFTCreate: {
			{
				Token:   "MYTOKEN-abcd",
				Set:     true,
				Address: "61646472",
			},
		},
	}
	require.Equal(t, expected, tokenRolesAndProperties.GetRoles())
}

func TestDctPropertiesProcUpgradeProperties(t *testing.T) {
	t.Parallel()

	dctPropProc := newDctPropertiesProcessor(&mock.PubkeyConverterMock{})

	event := &transaction.Event{
		Address:    []byte("addr"),
		Identifier: []byte(upgradePropertiesEvent),
		Topics:     [][]byte{[]byte("MYTOKEN-abcd"), big.NewInt(0).Bytes(), []byte("canMint"), []byte("true"), []byte("canBurn"), []byte("false")},
	}

	tokenRolesAndProperties := tokeninfo.NewTokenRolesAndProperties()
	dctPropProc.processEvent(&argsProcessEvent{
		event:                   event,
		tokenRolesAndProperties: tokenRolesAndProperties,
	})

	expected := []*tokeninfo.PropertiesData{
		{
			Token: "MYTOKEN-abcd",
			Properties: map[string]bool{
				"canMint": true,
				"canBurn": false,
			},
		},
	}
	require.Equal(t, expected, tokenRolesAndProperties.GetAllTokensWithProperties())
}

func TestCheckRolesBytes(t *testing.T) {
	t.Parallel()

	role1, _ := hex.DecodeString("01")
	role2, _ := hex.DecodeString("02")
	rolesBytes := [][]byte{role1, role2}
	require.False(t, checkRolesBytes(rolesBytes))

	role1 = []byte("DCTRoleNFTCreate")
	rolesBytes = [][]byte{role1}
	require.True(t, checkRolesBytes(rolesBytes))
}
