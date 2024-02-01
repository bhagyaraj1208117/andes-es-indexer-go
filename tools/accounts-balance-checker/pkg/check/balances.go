package check

type balancesDCT map[string]map[string]string

func newBalancesDCT() balancesDCT {
	return make(map[string]map[string]string)
}

func (be balancesDCT) extractBalancesFromResponse(responseAccounts *ResponseAccounts) {
	for _, hit := range responseAccounts.Hits.Hits {
		tokenIdentifier := hit.Source.TokenIdentifier
		if hit.Source.TokenIdentifier == "" {
			tokenIdentifier = hit.Source.TokenName
		}

		be.add(hit.Source.Address, tokenIdentifier, hit.Source.Balance)
	}
}

func (be balancesDCT) add(address, tokenIdentifier, value string) {
	_, ok := be[address]
	if !ok {
		be[address] = map[string]string{}
	}

	be[address][tokenIdentifier] = value
}

func (be balancesDCT) getBalancesForAddress(address string) map[string]string {
	return be[address]
}

func (be balancesDCT) getBalance(address, tokenIdentifier string) string {
	return be[address][tokenIdentifier]
}
