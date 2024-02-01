package accounts

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/bhagyaraj1208117/andes-core-go/core"
	"github.com/bhagyaraj1208117/andes-core-go/core/check"
	"github.com/bhagyaraj1208117/andes-core-go/data/alteredAccount"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/data"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/process/dataindexer"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/process/elasticproc/converters"
	logger "github.com/bhagyaraj1208117/andes-logger-go"
)

var log = logger.GetOrCreate("indexer/process/accounts")

// accountsProcessor is a structure responsible for processing accounts
type accountsProcessor struct {
	addressPubkeyConverter core.PubkeyConverter
	balanceConverter       dataindexer.BalanceConverter
}

// NewAccountsProcessor will create a new instance of accounts processor
func NewAccountsProcessor(
	addressPubkeyConverter core.PubkeyConverter,
	balanceConverter dataindexer.BalanceConverter,
) (*accountsProcessor, error) {
	if check.IfNil(addressPubkeyConverter) {
		return nil, dataindexer.ErrNilPubkeyConverter
	}
	if check.IfNil(balanceConverter) {
		return nil, dataindexer.ErrNilBalanceConverter
	}

	return &accountsProcessor{
		addressPubkeyConverter: addressPubkeyConverter,
		balanceConverter:       balanceConverter,
	}, nil
}

// GetAccounts will get accounts for regular operations and dct operations
func (ap *accountsProcessor) GetAccounts(coreAlteredAccounts map[string]*alteredAccount.AlteredAccount) ([]*data.Account, []*data.AccountDCT) {
	regularAccountsToIndex := make([]*data.Account, 0)
	accountsToIndexDCT := make([]*data.AccountDCT, 0)

	for _, alteredAccount := range coreAlteredAccounts {
		regularAccounts, dctAccounts := splitAlteredAccounts(alteredAccount)

		regularAccountsToIndex = append(regularAccountsToIndex, regularAccounts...)
		accountsToIndexDCT = append(accountsToIndexDCT, dctAccounts...)
	}

	return regularAccountsToIndex, accountsToIndexDCT
}

func splitAlteredAccounts(
	account *alteredAccount.AlteredAccount,
) ([]*data.Account, []*data.AccountDCT) {
	regularAccountsToIndex := make([]*data.Account, 0)
	accountsToIndexDCT := make([]*data.AccountDCT, 0)

	isSender, balanceChanged := false, false
	if account.AdditionalData != nil {
		isSender, balanceChanged = account.AdditionalData.IsSender, account.AdditionalData.BalanceChanged
	} else {
		log.Debug("accountsProcessor.splitAlteredAccounts - nil additional data")
	}

	//if the balance of the DCT receiver is 0 the receiver is a new account most probably, and we should index it
	ignoreAddress := !balanceChanged && notZeroBalance(account.Balance) && !isSender
	if !ignoreAddress {
		regularAccountsToIndex = append(regularAccountsToIndex, &data.Account{
			UserAccount: account,
			IsSender:    isSender,
		})
	}

	for _, info := range account.Tokens {
		accountDCT := &data.AccountDCT{
			Account:         account,
			TokenIdentifier: info.Identifier,
			NFTNonce:        info.Nonce,
			IsSender:        isSender,
		}
		if info.AdditionalData != nil {
			accountDCT.IsNFTCreate = info.AdditionalData.IsNFTCreate
		}

		accountsToIndexDCT = append(accountsToIndexDCT, accountDCT)

	}

	return regularAccountsToIndex, accountsToIndexDCT
}

func notZeroBalance(balance string) bool {
	return len(balance) > 0 && balance != "0"
}

// PrepareRegularAccountsMap will prepare a map of regular accounts
func (ap *accountsProcessor) PrepareRegularAccountsMap(timestamp uint64, accounts []*data.Account, shardID uint32) map[string]*data.AccountInfo {
	accountsMap := make(map[string]*data.AccountInfo)
	for _, userAccount := range accounts {
		address := userAccount.UserAccount.Address
		addressBytes, err := ap.addressPubkeyConverter.Decode(address)
		if err != nil {
			log.Warn("accountsProcessor.PrepareRegularAccountsMap: cannot decode address", "address", address, "error", err)
			continue
		}
		balance, ok := big.NewInt(0).SetString(userAccount.UserAccount.Balance, 10)
		if !ok {
			log.Warn("accountsProcessor.PrepareRegularAccountsMap: cannot cast account's balance to big int", "value", userAccount.UserAccount.Balance)
			continue
		}

		balanceAsFloat, err := ap.balanceConverter.ComputeBalanceAsFloat(balance)
		if err != nil {
			log.Warn("accountsProcessor.PrepareRegularAccountsMap: cannot compute balance as num",
				"balance", balance, "address", address, "error", err)
		}

		acc := &data.AccountInfo{
			Address:         address,
			Nonce:           userAccount.UserAccount.Nonce,
			Balance:         converters.BigIntToString(balance),
			BalanceNum:      balanceAsFloat,
			IsSender:        userAccount.IsSender,
			IsSmartContract: core.IsSmartContractAddress(addressBytes),
			Timestamp:       time.Duration(timestamp),
			ShardID:         shardID,
		}

		ap.addAdditionalDataInAccount(userAccount.UserAccount.AdditionalData, acc)

		accountsMap[address] = acc
	}

	return accountsMap
}

func (ap *accountsProcessor) addAdditionalDataInAccount(additionalData *alteredAccount.AdditionalAccountData, account *data.AccountInfo) {
	if additionalData == nil {
		return
	}

	account.UserName = additionalData.UserName
	account.CurrentOwner = additionalData.CurrentOwner
	account.RootHash = additionalData.RootHash
	account.CodeHash = additionalData.CodeHash
	account.CodeMetadata = additionalData.CodeMetadata

	ap.addDeveloperRewardsInAccount(additionalData, account)
}

func (ap *accountsProcessor) addDeveloperRewardsInAccount(additionalData *alteredAccount.AdditionalAccountData, account *data.AccountInfo) {
	if additionalData.DeveloperRewards == "" {
		return
	}

	developerRewardsBig, ok := big.NewInt(0).SetString(additionalData.DeveloperRewards, 10)
	if !ok {
		log.Warn("ap.addDeveloperRewardsInAccountInfo cannot convert developer rewards in number", "address", account.Address)
		return
	}

	account.DeveloperRewards = additionalData.DeveloperRewards

	developerRewardsNum, err := ap.balanceConverter.ComputeBalanceAsFloat(developerRewardsBig)
	if err != nil {
		log.Warn("accountsProcessor.addDeveloperRewardsInAccount: cannot compute developer rewards as num",
			"developer rewards", developerRewardsBig, "error", err)
	}

	account.DeveloperRewardsNum = developerRewardsNum
}

// PrepareAccountsMapDCT will prepare a map of accounts with DCT tokens
func (ap *accountsProcessor) PrepareAccountsMapDCT(
	timestamp uint64,
	accounts []*data.AccountDCT,
	tagsCount data.CountTags,
	shardID uint32,
) (map[string]*data.AccountInfo, data.TokensHandler) {
	tokensData := data.NewTokensInfo()
	accountsDCTMap := make(map[string]*data.AccountInfo)
	for _, accountDCT := range accounts {
		address := accountDCT.Account.Address
		addressBytes, err := ap.addressPubkeyConverter.Decode(address)
		if err != nil {
			log.Warn("accountsProcessor.PrepareAccountsMapDCT: cannot decode address", "address", address, "error", err)
			continue
		}
		balance, properties, tokenMetaData, err := ap.getDCTInfo(accountDCT)
		if err != nil {
			log.Warn("accountsProcessor.PrepareAccountsMapDCT: cannot get dct info from account",
				"address", address,
				"error", err.Error())
			continue
		}

		if tokenMetaData != nil && accountDCT.IsNFTCreate {
			tagsCount.ParseTags(tokenMetaData.Tags)
		}

		tokenIdentifier := converters.ComputeTokenIdentifier(accountDCT.TokenIdentifier, accountDCT.NFTNonce)
		balanceNum, err := ap.balanceConverter.ConvertBigValueToFloat(balance)
		if err != nil {
			log.Warn("accountsProcessor.PrepareAccountsMapDCT: cannot compute dct balance as num",
				"balance", balance, "address", address, "error", err, "token", tokenIdentifier)
		}

		acc := &data.AccountInfo{
			Address:         address,
			TokenName:       accountDCT.TokenIdentifier,
			TokenIdentifier: tokenIdentifier,
			TokenNonce:      accountDCT.NFTNonce,
			Balance:         balance.String(),
			BalanceNum:      balanceNum,
			Properties:      properties,
			Frozen:          isFrozen(properties),
			IsSender:        accountDCT.IsSender,
			IsSmartContract: core.IsSmartContractAddress(addressBytes),
			Data:            tokenMetaData,
			Timestamp:       time.Duration(timestamp),
			ShardID:         shardID,
		}

		if acc.TokenNonce == 0 {
			acc.Type = core.FungibleDCT
		}

		keyInMap := fmt.Sprintf("%s-%s-%d", acc.Address, acc.TokenName, accountDCT.NFTNonce)
		accountsDCTMap[keyInMap] = acc

		if acc.Balance == "0" || acc.Balance == "" {
			continue
		}

		tokensData.Add(&data.TokenInfo{
			Token:      accountDCT.TokenIdentifier,
			Identifier: tokenIdentifier,
		})
	}

	return accountsDCTMap, tokensData
}

// PrepareAccountsHistory will prepare a map of accounts history balance from a map of accounts
func (ap *accountsProcessor) PrepareAccountsHistory(
	timestamp uint64,
	accounts map[string]*data.AccountInfo,
	shardID uint32,
) map[string]*data.AccountBalanceHistory {
	accountsMap := make(map[string]*data.AccountBalanceHistory)
	for _, userAccount := range accounts {
		acc := &data.AccountBalanceHistory{
			Address:         userAccount.Address,
			Balance:         userAccount.Balance,
			Timestamp:       time.Duration(timestamp),
			Token:           userAccount.TokenName,
			TokenNonce:      userAccount.TokenNonce,
			IsSender:        userAccount.IsSender,
			IsSmartContract: userAccount.IsSmartContract,
			Identifier:      converters.ComputeTokenIdentifier(userAccount.TokenName, userAccount.TokenNonce),
			ShardID:         shardID,
		}
		keyInMap := fmt.Sprintf("%s-%s-%d", acc.Address, acc.Token, acc.TokenNonce)
		accountsMap[keyInMap] = acc
	}

	return accountsMap
}

func (ap *accountsProcessor) getDCTInfo(accountDCT *data.AccountDCT) (*big.Int, string, *data.TokenMetaData, error) {
	if accountDCT.TokenIdentifier == "" {
		return big.NewInt(0), "", nil, nil
	}
	if accountDCT.IsNFTOperation && accountDCT.NFTNonce == 0 {
		return big.NewInt(0), "", nil, nil
	}

	accountTokenData := &alteredAccount.AccountTokenData{}
	for _, tokenData := range accountDCT.Account.Tokens {
		if tokenData.Identifier == accountDCT.TokenIdentifier && tokenData.Nonce == accountDCT.NFTNonce {
			accountTokenData = tokenData
		}
	}

	value, _ := big.NewInt(0).SetString(accountTokenData.Balance, 10)
	if value == nil {
		return big.NewInt(0), "", nil, nil
	}

	tokenMetaData := converters.PrepareTokenMetaData(accountTokenData.MetaData)

	return value, accountTokenData.Properties, tokenMetaData, nil
}

// PutTokenMedataDataInTokens will put the TokenMedata in provided tokens data
func (ap *accountsProcessor) PutTokenMedataDataInTokens(tokensData []*data.TokenInfo, coreAlteredAccounts map[string]*alteredAccount.AlteredAccount) {
	for _, tokenData := range tokensData {
		if tokenData.Data != nil || tokenData.Nonce == 0 {
			continue
		}

		metadata, errLoad := ap.loadMetadataForToken(tokenData, coreAlteredAccounts)
		if errLoad != nil {
			log.Warn("accountsProcessor.PutTokenMedataDataInTokens: cannot load token metadata",
				"token identifier ", tokenData.Identifier,
				"error", errLoad.Error())

			continue
		}

		tokenData.Data = converters.PrepareTokenMetaData(metadata)
	}
}

func (ap *accountsProcessor) loadMetadataForToken(
	tokenData *data.TokenInfo,
	coreAlteredAccounts map[string]*alteredAccount.AlteredAccount,
) (*alteredAccount.TokenMetaData, error) {
	for _, account := range coreAlteredAccounts {
		for _, token := range account.Tokens {
			if tokenData.Token == token.Identifier && tokenData.Nonce == token.Nonce {
				return token.MetaData, nil
			}
		}
	}

	return nil, fmt.Errorf("%w for identifier %s and nonce %d", errTokenNotFound, tokenData.Identifier, tokenData.Nonce)
}

func isFrozen(properties string) bool {
	decoded, err := hex.DecodeString(properties)
	if err != nil {
		log.Debug("isFrozen() cannot decode token properties", "error", err)
		return false
	}
	if len(decoded) == 0 {
		return false
	}

	return (decoded[0] & 1) != 0
}
