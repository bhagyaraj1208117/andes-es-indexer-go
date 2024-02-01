package logsevents

import (
	"testing"
	"time"

	"github.com/bhagyaraj1208117/andes-core-go/core"
	"github.com/bhagyaraj1208117/andes-core-go/data/transaction"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/data"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/mock"
	"github.com/stretchr/testify/require"
)

func TestIssueDCTProcessor(t *testing.T) {
	t.Parallel()

	dctIssueProc := newDCTIssueProcessor(&mock.PubkeyConverterMock{})

	event := &transaction.Event{
		Address:    []byte("addr"),
		Identifier: []byte(issueNonFungibleDCTFunc),
		Topics:     [][]byte{[]byte("MYTOKEN-abcd"), []byte("my-token"), []byte("MYTOKEN"), []byte(core.NonFungibleDCT)},
	}
	args := &argsProcessEvent{
		timestamp:   1234,
		event:       event,
		selfShardID: core.MetachainShardId,
	}

	res := dctIssueProc.processEvent(args)

	require.Equal(t, &data.TokenInfo{
		Token:        "MYTOKEN-abcd",
		Name:         "my-token",
		Ticker:       "MYTOKEN",
		Timestamp:    time.Duration(1234),
		Type:         core.NonFungibleDCT,
		Issuer:       "61646472",
		CurrentOwner: "61646472",
		OwnersHistory: []*data.OwnerData{
			{
				Address:   "61646472",
				Timestamp: time.Duration(1234),
			},
		},
		Properties: &data.TokenProperties{},
	}, res.tokenInfo)
}

func TestIssueDCTProcessor_TransferOwnership(t *testing.T) {
	t.Parallel()

	dctIssueProc := newDCTIssueProcessor(&mock.PubkeyConverterMock{})

	event := &transaction.Event{
		Address:    []byte("addr"),
		Identifier: []byte(transferOwnershipFunc),
		Topics:     [][]byte{[]byte("MYTOKEN-abcd"), []byte("my-token"), []byte("MYTOKEN"), []byte(core.NonFungibleDCT), []byte("newOwner")},
	}
	args := &argsProcessEvent{
		timestamp:   1234,
		event:       event,
		selfShardID: core.MetachainShardId,
	}

	res := dctIssueProc.processEvent(args)

	require.Equal(t, &data.TokenInfo{
		Token:        "MYTOKEN-abcd",
		Name:         "my-token",
		Ticker:       "MYTOKEN",
		Timestamp:    time.Duration(1234),
		Type:         core.NonFungibleDCT,
		Issuer:       "61646472",
		CurrentOwner: "6e65774f776e6572",
		OwnersHistory: []*data.OwnerData{
			{
				Address:   "6e65774f776e6572",
				Timestamp: time.Duration(1234),
			},
		},
		TransferOwnership: true,
		Properties:        &data.TokenProperties{},
	}, res.tokenInfo)
}

func TestIssueDCTProcessor_EventWithShardID0ShouldBeIgnored(t *testing.T) {
	t.Parallel()

	dctIssueProc := newDCTIssueProcessor(&mock.PubkeyConverterMock{})

	event := &transaction.Event{
		Address:    []byte("addr"),
		Identifier: []byte(transferOwnershipFunc),
		Topics:     [][]byte{[]byte("MYTOKEN-abcd"), []byte("my-token"), []byte("MYTOKEN"), []byte(core.NonFungibleDCT), []byte("newOwner")},
	}
	args := &argsProcessEvent{
		timestamp:   1234,
		event:       event,
		selfShardID: 0,
	}

	res := dctIssueProc.processEvent(args)
	require.False(t, res.processed)
}
