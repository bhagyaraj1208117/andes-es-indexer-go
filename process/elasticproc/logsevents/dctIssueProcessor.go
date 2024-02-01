package logsevents

import (
	"math/big"
	"time"

	"github.com/bhagyaraj1208117/andes-core-go/core"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/data"
)

const (
	numIssueLogTopics = 4

	issueFungibleDCTFunc     = "issue"
	issueSemiFungibleDCTFunc = "issueSemiFungible"
	issueNonFungibleDCTFunc  = "issueNonFungible"
	registerMetaDCTFunc      = "registerMetaDCT"
	changeSFTToMetaDCTFunc   = "changeSFTToMetaDCT"
	transferOwnershipFunc    = "transferOwnership"
	registerAndSetRolesFunc  = "registerAndSetAllRoles"
)

type dctIssueProcessor struct {
	pubkeyConverter            core.PubkeyConverter
	issueOperationsIdentifiers map[string]struct{}
}

func newDCTIssueProcessor(pubkeyConverter core.PubkeyConverter) *dctIssueProcessor {
	return &dctIssueProcessor{
		pubkeyConverter: pubkeyConverter,
		issueOperationsIdentifiers: map[string]struct{}{
			issueFungibleDCTFunc:     {},
			issueSemiFungibleDCTFunc: {},
			issueNonFungibleDCTFunc:  {},
			registerMetaDCTFunc:      {},
			changeSFTToMetaDCTFunc:   {},
			transferOwnershipFunc:    {},
			registerAndSetRolesFunc:  {},
		},
	}
}

func (eip *dctIssueProcessor) processEvent(args *argsProcessEvent) argOutputProcessEvent {
	if args.selfShardID != core.MetachainShardId {
		return argOutputProcessEvent{}
	}

	identifierStr := string(args.event.GetIdentifier())
	_, ok := eip.issueOperationsIdentifiers[identifierStr]
	if !ok {
		return argOutputProcessEvent{}
	}

	topics := args.event.GetTopics()
	if len(topics) < numIssueLogTopics {
		return argOutputProcessEvent{
			processed: true,
		}
	}

	// topics slice contains:
	// topics[0] -- token identifier
	// topics[1] -- token name
	// topics[2] -- token ticker
	// topics[3] -- token type
	// topics[4] -- num decimals / new owner address in case of transferOwnershipFunc
	if len(topics[0]) == 0 {
		return argOutputProcessEvent{
			processed: true,
		}
	}

	numDecimals := uint64(0)
	if len(topics) == numIssueLogTopics+1 && identifierStr != transferOwnershipFunc {
		numDecimals = big.NewInt(0).SetBytes(topics[4]).Uint64()
	}

	encodedAddr := eip.pubkeyConverter.SilentEncode(args.event.GetAddress(), log)

	tokenInfo := &data.TokenInfo{
		Token:        string(topics[0]),
		Name:         string(topics[1]),
		Ticker:       string(topics[2]),
		Type:         string(topics[3]),
		NumDecimals:  numDecimals,
		Issuer:       encodedAddr,
		CurrentOwner: encodedAddr,
		Timestamp:    time.Duration(args.timestamp),
		OwnersHistory: []*data.OwnerData{
			{
				Address:   encodedAddr,
				Timestamp: time.Duration(args.timestamp),
			},
		},
		Properties: &data.TokenProperties{},
	}

	if identifierStr == transferOwnershipFunc && len(topics) >= numIssueLogTopics+1 {
		newOwner := eip.pubkeyConverter.SilentEncode(topics[4], log)
		tokenInfo.TransferOwnership = true
		tokenInfo.CurrentOwner = newOwner
		tokenInfo.OwnersHistory[0].Address = newOwner
	}

	return argOutputProcessEvent{
		tokenInfo: tokenInfo,
		processed: true,
	}
}
