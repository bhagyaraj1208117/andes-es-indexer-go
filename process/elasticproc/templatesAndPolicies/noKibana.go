package templatesAndPolicies

import (
	"bytes"

	indexer "github.com/bhagyaraj1208117/andes-es-indexer-go/process/dataindexer"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/templates/noKibana"
)

type templatesAndPolicyReaderNoKibana struct{}

// NewTemplatesAndPolicyReaderNoKibana will create a new instance of templatesAndPolicyReaderNoKibana
func NewTemplatesAndPolicyReaderNoKibana() *templatesAndPolicyReaderNoKibana {
	return new(templatesAndPolicyReaderNoKibana)
}

// GetElasticTemplatesAndPolicies will return templates and policies
func (tr *templatesAndPolicyReaderNoKibana) GetElasticTemplatesAndPolicies() (map[string]*bytes.Buffer, map[string]*bytes.Buffer, error) {
	indexPolicies := make(map[string]*bytes.Buffer)
	indexTemplates := make(map[string]*bytes.Buffer)

	indexTemplates["opendistro"] = noKibana.OpenDistro.ToBuffer()
	indexTemplates[indexer.TransactionsIndex] = noKibana.Transactions.ToBuffer()
	indexTemplates[indexer.BlockIndex] = noKibana.Blocks.ToBuffer()
	indexTemplates[indexer.MiniblocksIndex] = noKibana.Miniblocks.ToBuffer()
	indexTemplates[indexer.RatingIndex] = noKibana.Rating.ToBuffer()
	indexTemplates[indexer.RoundsIndex] = noKibana.Rounds.ToBuffer()
	indexTemplates[indexer.ValidatorsIndex] = noKibana.Validators.ToBuffer()
	indexTemplates[indexer.AccountsIndex] = noKibana.Accounts.ToBuffer()
	indexTemplates[indexer.AccountsHistoryIndex] = noKibana.AccountsHistory.ToBuffer()
	indexTemplates[indexer.AccountsDCTIndex] = noKibana.AccountsDCT.ToBuffer()
	indexTemplates[indexer.AccountsDCTHistoryIndex] = noKibana.AccountsDCTHistory.ToBuffer()
	indexTemplates[indexer.EpochInfoIndex] = noKibana.EpochInfo.ToBuffer()
	indexTemplates[indexer.ReceiptsIndex] = noKibana.Receipts.ToBuffer()
	indexTemplates[indexer.ScResultsIndex] = noKibana.SCResults.ToBuffer()
	indexTemplates[indexer.SCDeploysIndex] = noKibana.SCDeploys.ToBuffer()
	indexTemplates[indexer.TokensIndex] = noKibana.Tokens.ToBuffer()
	indexTemplates[indexer.TagsIndex] = noKibana.Tags.ToBuffer()
	indexTemplates[indexer.LogsIndex] = noKibana.Logs.ToBuffer()
	indexTemplates[indexer.DelegatorsIndex] = noKibana.Delegators.ToBuffer()
	indexTemplates[indexer.OperationsIndex] = noKibana.Operations.ToBuffer()
	indexTemplates[indexer.DCTsIndex] = noKibana.DCTs.ToBuffer()

	return indexTemplates, indexPolicies, nil
}
