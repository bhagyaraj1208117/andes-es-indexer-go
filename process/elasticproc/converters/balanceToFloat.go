package converters

import (
	"errors"
	"math"
	"math/big"

	"github.com/bhagyaraj1208117/andes-core-go/core"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/data"
	indexer "github.com/bhagyaraj1208117/andes-es-indexer-go/process/dataindexer"
)

const (
	numDecimalsInFloatBalance    = 10
	numDecimalsInFloatBalanceDCT = 18
)

var (
	errValueTooBig        = errors.New("provided value is too big")
	errCastStringToBigInt = errors.New("cannot convert string to big value")
)

var zero = big.NewInt(0)

type balanceConverter struct {
	dividerForDenomination float64
	balancePrecision       float64
	balancePrecisionDCT    float64
}

// NewBalanceConverter will create a new instance of balance converter
func NewBalanceConverter(denomination int) (*balanceConverter, error) {
	if denomination < 0 {
		return nil, indexer.ErrNegativeDenominationValue
	}

	return &balanceConverter{
		balancePrecision:       math.Pow(10, float64(numDecimalsInFloatBalance)),
		balancePrecisionDCT:    math.Pow(10, float64(numDecimalsInFloatBalanceDCT)),
		dividerForDenomination: math.Pow(10, float64(denomination)),
	}, nil
}

// ComputeBalanceAsFloat will compute balance as float
func (bc *balanceConverter) ComputeBalanceAsFloat(balance *big.Int) (float64, error) {
	return bc.computeBalanceAsFloat(balance, bc.balancePrecision)
}

// ConvertBigValueToFloat will convert big value to float
func (bc *balanceConverter) ConvertBigValueToFloat(balance *big.Int) (float64, error) {
	return bc.computeBalanceAsFloat(balance, bc.balancePrecisionDCT)
}

// ComputeSliceOfStringsAsFloat will compute the provided slice of string values in float values
func (bc *balanceConverter) ComputeSliceOfStringsAsFloat(values []string) ([]float64, error) {
	floatValues := make([]float64, 0, len(values))

	for _, value := range values {
		valueBig, ok := big.NewInt(0).SetString(value, 10)
		if !ok {
			return nil, errCastStringToBigInt
		}

		valueNum, err := bc.ConvertBigValueToFloat(valueBig)
		if err != nil {
			return nil, err
		}

		floatValues = append(floatValues, valueNum)
	}

	return floatValues, nil
}

func (bc *balanceConverter) computeBalanceAsFloat(balance *big.Int, balancePrecision float64) (float64, error) {
	if balance == nil || balance.Cmp(zero) == 0 {
		return 0, nil
	}
	if len(balance.Bytes()) > data.MaxDCTValueLength {
		return 0, errValueTooBig
	}

	balanceBigFloat := big.NewFloat(0).SetInt(balance)
	balanceFloat64, _ := balanceBigFloat.Float64()

	bal := balanceFloat64 / bc.dividerForDenomination

	balanceFloatWithDecimals := math.Round(bal*balancePrecision) / balancePrecision

	value := core.MaxFloat64(balanceFloatWithDecimals, 0)
	if math.IsInf(value, +1) || math.IsInf(value, -1) {
		return 0, errValueTooBig
	}

	return value, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (bc *balanceConverter) IsInterfaceNil() bool {
	return bc == nil
}

// BigIntToString will convert a big.Int to string
func BigIntToString(value *big.Int) string {
	if value == nil {
		return "0"
	}

	return value.String()
}
