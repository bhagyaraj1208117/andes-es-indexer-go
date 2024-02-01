package mock

import (
	"github.com/bhagyaraj1208117/andes-core-go/core"
	"github.com/bhagyaraj1208117/andes-core-go/data/block"
)

// BlockContainerStub -
type BlockContainerStub struct {
	GetCalled func(headerType core.HeaderType) (block.EmptyBlockCreator, error)
}

// Get -
func (bcs *BlockContainerStub) Get(headerType core.HeaderType) (block.EmptyBlockCreator, error) {
	if bcs.GetCalled != nil {
		return bcs.GetCalled(headerType)
	}

	return nil, nil
}
