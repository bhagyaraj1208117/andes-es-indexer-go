package logsevents

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/bhagyaraj1208117/andes-core-go/core"
	"github.com/bhagyaraj1208117/andes-core-go/data/outport"
	"github.com/bhagyaraj1208117/andes-core-go/data/transaction"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/data"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/mock"
	elasticIndexer "github.com/bhagyaraj1208117/andes-es-indexer-go/process/dataindexer"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/process/elasticproc/converters"
	"github.com/stretchr/testify/require"
)

func createMockArgs() ArgsLogsAndEventsProcessor {
	balanceConverter, _ := converters.NewBalanceConverter(10)
	return ArgsLogsAndEventsProcessor{
		PubKeyConverter:  &mock.PubkeyConverterMock{},
		Marshalizer:      &mock.MarshalizerMock{},
		BalanceConverter: balanceConverter,
		Hasher:           &mock.HasherMock{},
	}
}

func TestNewLogsAndEventsProcessor(t *testing.T) {
	t.Parallel()

	args := createMockArgs()
	args.PubKeyConverter = nil
	_, err := NewLogsAndEventsProcessor(args)
	require.Equal(t, elasticIndexer.ErrNilPubkeyConverter, err)

	args = createMockArgs()
	args.Marshalizer = nil
	_, err = NewLogsAndEventsProcessor(args)
	require.Equal(t, elasticIndexer.ErrNilMarshalizer, err)

	args = createMockArgs()
	args.BalanceConverter = nil
	_, err = NewLogsAndEventsProcessor(args)
	require.Equal(t, elasticIndexer.ErrNilBalanceConverter, err)

	args = createMockArgs()
	args.Hasher = nil
	_, err = NewLogsAndEventsProcessor(args)
	require.Equal(t, elasticIndexer.ErrNilHasher, err)

	args = createMockArgs()
	proc, err := NewLogsAndEventsProcessor(args)
	require.NotNil(t, proc)
	require.Nil(t, err)
}

func TestLogsAndEventsProcessor_ExtractDataFromLogsAndPutInAltered(t *testing.T) {
	t.Parallel()

	logsAndEvents := []*outport.LogData{
		nil,
		{
			TxHash: hex.EncodeToString([]byte("h3")),
			Log: &transaction.Log{
				Events: []*transaction.Event{
					{
						Address:    []byte("addr"),
						Identifier: []byte(core.SCDeployIdentifier),
						Topics:     [][]byte{[]byte("addr1"), []byte("addr2"), []byte("codeHash")},
					},
				},
			},
		},
		{
			TxHash: hex.EncodeToString([]byte("h1")),
			Log: &transaction.Log{
				Address: []byte("address"),
				Events: []*transaction.Event{
					{
						Address:    []byte("addr"),
						Identifier: []byte(core.BuiltInFunctionDCTNFTTransfer),
						Topics:     [][]byte{[]byte("my-token"), big.NewInt(0).SetUint64(1).Bytes(), big.NewInt(100).Bytes(), []byte("receiver")},
					},
				},
			},
		},
		{
			TxHash: hex.EncodeToString([]byte("h2")),
			Log: &transaction.Log{
				Events: []*transaction.Event{
					{
						Address:    []byte("addr"),
						Identifier: []byte(core.BuiltInFunctionDCTTransfer),
						Topics:     [][]byte{[]byte("dct"), big.NewInt(0).Bytes(), big.NewInt(0).SetUint64(100).Bytes(), []byte("receiver")},
					},
					nil,
				},
			},
		},
		{
			TxHash: hex.EncodeToString([]byte("h4")),
			Log: &transaction.Log{
				Events: []*transaction.Event{
					{
						Address:    []byte("addr"),
						Identifier: []byte(issueSemiFungibleDCTFunc),
						Topics:     [][]byte{[]byte("SEMI-abcd"), []byte("semi-token"), []byte("SEMI"), []byte(core.SemiFungibleDCT)},
					},
					nil,
				},
			},
		},
		{
			TxHash: hex.EncodeToString([]byte("h5")),
			Log: &transaction.Log{
				Address: []byte("contract"),
				Events: []*transaction.Event{
					{
						Address:    []byte("addr"),
						Identifier: []byte(delegateFunc),
						Topics:     [][]byte{big.NewInt(1000).Bytes(), big.NewInt(1000000000).Bytes(), big.NewInt(10).Bytes(), big.NewInt(1000000000).Bytes()},
					},
				},
			},
		},
		{
			TxHash: hex.EncodeToString([]byte("h6")),
			Log: &transaction.Log{
				Address: []byte("contract-second"),
				Events: []*transaction.Event{
					{
						Address:    []byte("addr"),
						Identifier: []byte(delegateFunc),
						Topics:     [][]byte{big.NewInt(1000).Bytes(), big.NewInt(1000000000).Bytes(), big.NewInt(10).Bytes(), big.NewInt(1000000000).Bytes()},
					},
				},
			},
		},
	}

	res := &data.PreparedResults{
		Transactions: []*data.Transaction{
			{
				Hash: "6831",
			},
		},
		ScResults: []*data.ScResult{
			{
				Hash: "6832",
			},
		},
	}

	args := createMockArgs()
	balanceConverter, _ := converters.NewBalanceConverter(10)
	args.BalanceConverter = balanceConverter
	proc, _ := NewLogsAndEventsProcessor(args)

	resLogs := proc.ExtractDataFromLogs(logsAndEvents, res, 1000, core.MetachainShardId, 3)
	require.NotNil(t, resLogs.Tokens)
	require.True(t, res.Transactions[0].HasOperations)
	require.True(t, res.ScResults[0].HasOperations)
	require.True(t, res.Transactions[0].HasLogs)
	require.True(t, res.ScResults[0].HasLogs)

	require.Equal(t, &data.ScDeployInfo{
		TxHash:       "6833",
		Creator:      "6164647232",
		CurrentOwner: "6164647232",
		Timestamp:    uint64(1000),
		CodeHash:     []byte("codeHash"),
	}, resLogs.ScDeploys["6164647231"])

	require.Equal(t, &data.TokenInfo{
		Name:         "semi-token",
		Ticker:       "SEMI",
		Token:        "SEMI-abcd",
		Type:         core.SemiFungibleDCT,
		Timestamp:    1000,
		Issuer:       "61646472",
		CurrentOwner: "61646472",
		OwnersHistory: []*data.OwnerData{
			{
				Address:   "61646472",
				Timestamp: 1000,
			},
		},
		Properties: &data.TokenProperties{},
	}, resLogs.TokensInfo[0])

	require.Equal(t, &data.Delegator{
		Address:        "61646472",
		Contract:       "636f6e7472616374",
		ActiveStakeNum: 0.1,
		ActiveStake:    "1000000000",
		Timestamp:      time.Duration(1000),
	}, resLogs.Delegators["61646472636f6e7472616374"])
	require.Equal(t, &data.Delegator{
		Address:        "61646472",
		Contract:       "636f6e74726163742d7365636f6e64",
		ActiveStakeNum: 0.1,
		ActiveStake:    "1000000000",
		Timestamp:      time.Duration(1000),
	}, resLogs.Delegators["61646472636f6e74726163742d7365636f6e64"])
}

func TestLogsAndEventsProcessor_PrepareLogsForDB(t *testing.T) {
	t.Parallel()

	logsAndEvents := []*outport.LogData{
		nil,
		{
			TxHash: hex.EncodeToString([]byte("txHash")),
			Log: &transaction.Log{
				Address: []byte("address"),
				Events: []*transaction.Event{
					{
						Address:        []byte("addr"),
						Identifier:     []byte(core.BuiltInFunctionDCTNFTTransfer),
						Topics:         [][]byte{[]byte("my-token"), big.NewInt(0).SetUint64(1).Bytes(), []byte("receiver")},
						AdditionalData: [][]byte{[]byte("something")},
					},
				},
			},
		},
	}

	args := createMockArgs()
	proc, _ := NewLogsAndEventsProcessor(args)

	_ = proc.ExtractDataFromLogs(nil, &data.PreparedResults{ScResults: []*data.ScResult{
		{
			Hash:           "747848617368",
			OriginalTxHash: "orignalHash",
		},
	}}, 1234, 0, 3)

	logsDB := proc.PrepareLogsForDB(logsAndEvents, 1234)
	require.Equal(t, &data.Logs{
		ID:             "747848617368",
		Address:        "61646472657373",
		OriginalTxHash: "orignalHash",
		Timestamp:      time.Duration(1234),
		Events: []*data.Event{
			{
				Address:        "61646472",
				Identifier:     core.BuiltInFunctionDCTNFTTransfer,
				Topics:         [][]byte{[]byte("my-token"), big.NewInt(0).SetUint64(1).Bytes(), []byte("receiver")},
				AdditionalData: [][]byte{[]byte("something")},
			},
		},
	}, logsDB[0])
}

func TestLogsAndEventsProcessor_ExtractDataFromLogsNFTBurn(t *testing.T) {
	t.Parallel()

	logsAndEventsSlice := make([]*outport.LogData, 1)
	logsAndEventsSlice[0] = &outport.LogData{
		TxHash: "h1",
		Log: &transaction.Log{
			Address: []byte("address"),
			Events: []*transaction.Event{
				{
					Address:    []byte("addr"),
					Identifier: []byte(core.BuiltInFunctionDCTNFTBurn),
					Topics:     [][]byte{[]byte("MY-NFT"), big.NewInt(2).Bytes(), big.NewInt(1).Bytes()},
				},
			},
		},
	}

	res := &data.PreparedResults{
		Transactions: []*data.Transaction{
			{
				Hash: "6831",
			},
		},
		ScResults: []*data.ScResult{
			{
				Hash: "6832",
			},
		},
	}

	args := createMockArgs()
	balanceConverter, _ := converters.NewBalanceConverter(10)
	args.BalanceConverter = balanceConverter
	proc, _ := NewLogsAndEventsProcessor(args)

	resLogs := proc.ExtractDataFromLogs(logsAndEventsSlice, res, 1000, 2, 3)
	require.Equal(t, 1, resLogs.TokensSupply.Len())

	tokensSupply := resLogs.TokensSupply.GetAll()
	require.Equal(t, "MY-NFT", tokensSupply[0].Token)
	require.Equal(t, "MY-NFT-02", tokensSupply[0].Identifier)
}
