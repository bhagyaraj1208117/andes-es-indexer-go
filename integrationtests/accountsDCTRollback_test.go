//go:build integrationtests

package integrationtests

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/bhagyaraj1208117/andes-core-go/core"
	"github.com/bhagyaraj1208117/andes-core-go/data/alteredAccount"
	dataBlock "github.com/bhagyaraj1208117/andes-core-go/data/block"
	"github.com/bhagyaraj1208117/andes-core-go/data/dct"
	"github.com/bhagyaraj1208117/andes-core-go/data/outport"
	"github.com/bhagyaraj1208117/andes-core-go/data/transaction"
	indexerdata "github.com/bhagyaraj1208117/andes-es-indexer-go/process/dataindexer"
	"github.com/stretchr/testify/require"
)

func TestAccountsDCTDeleteOnRollback(t *testing.T) {
	setLogLevelDebug()

	esClient, err := createESClient(esURL)
	require.Nil(t, err)

	dctToken := &dct.DCToken{
		Value:      big.NewInt(1000),
		Properties: []byte("3032"),
		TokenMetaData: &dct.MetaData{
			Creator: []byte("creator"),
		},
	}
	addr := "moa1sqy2ywvswp09ef7qwjhv8zwr9kzz3xas6y2ye5nuryaz0wcnfzzs7cfj8p"
	coreAlteredAccounts := map[string]*alteredAccount.AlteredAccount{
		addr: {
			Address: addr,
			Tokens: []*alteredAccount.AccountTokenData{
				{
					Identifier: "TOKEN-eeee",
					Nonce:      2,
					Balance:    "1000",
					MetaData: &alteredAccount.TokenMetaData{
						Creator: "creator",
					},
					Properties: "3032",
				},
			},
		},
	}

	esProc, err := CreateElasticProcessor(esClient)
	require.Nil(t, err)

	// CREATE SEMI-FUNGIBLE TOKEN
	dctDataBytes, _ := json.Marshal(dctToken)
	pool := &outport.TransactionPool{
		Logs: []*outport.LogData{
			{
				TxHash: hex.EncodeToString([]byte("h1")),
				Log: &transaction.Log{
					Events: []*transaction.Event{
						{
							Address:    decodeAddress(addr),
							Identifier: []byte(core.BuiltInFunctionDCTNFTCreate),
							Topics:     [][]byte{[]byte("TOKEN-eeee"), big.NewInt(2).Bytes(), big.NewInt(1).Bytes(), dctDataBytes},
						},
						nil,
					},
				},
			},
		},
	}

	body := &dataBlock.Body{}
	header := &dataBlock.Header{
		Round:     50,
		TimeStamp: 5040,
		ShardID:   2,
	}

	err = esProc.SaveTransactions(createOutportBlockWithHeader(body, header, pool, coreAlteredAccounts, testNumOfShards))
	require.Nil(t, err)

	ids := []string{fmt.Sprintf("%s-TOKEN-eeee-02", addr)}
	genericResponse := &GenericResponse{}
	err = esClient.DoMultiGet(context.Background(), ids, indexerdata.AccountsDCTIndex, true, genericResponse)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/accountsDCTRollback/account-after-create.json"), string(genericResponse.Docs[0].Source))

	// DO ROLLBACK
	err = esProc.RemoveAccountsDCT(5040, 2)
	require.Nil(t, err)

	err = esClient.DoMultiGet(context.Background(), ids, indexerdata.AccountsDCTIndex, true, genericResponse)
	require.Nil(t, err)
	require.False(t, genericResponse.Docs[0].Found)
}
