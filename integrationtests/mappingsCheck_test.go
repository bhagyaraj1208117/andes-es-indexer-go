//go:build integrationtests

package integrationtests

import (
	"testing"

	"github.com/bhagyaraj1208117/andes-es-indexer-go/process/dataindexer"
	"github.com/stretchr/testify/require"
)

func TestMappingsOfDCTsIndex(t *testing.T) {
	setLogLevelDebug()

	esClient, err := createESClient(esURL)
	require.Nil(t, err)

	_, err = CreateElasticProcessor(esClient)
	require.Nil(t, err)

	mappings, err := getIndexMappings(dataindexer.DCTsIndex)
	require.Nil(t, err)
	require.JSONEq(t, readExpectedResult("./testdata/mappings/dcts.json"), mappings)
}
