package transactions

import (
	"github.com/bhagyaraj1208117/andes-core-go/data/outport"
	datafield "github.com/bhagyaraj1208117/andes-vm-common-go/parsers/dataField"
)

// DataFieldParser defines what a data field parser should be able to do
type DataFieldParser interface {
	Parse(dataField []byte, sender, receiver []byte, numOfShards uint32) *datafield.ResponseParseData
}

type feeInfoHandler interface {
	GetFeeInfo() *outport.FeeInfo
}
