package convert

import (
	"fmt"
	"math/big"
	"strconv"
)

func HexToDecimal(hex string) (*int64, error) {
	decimal, err := strconv.ParseInt(hex, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting from hex to decimal: %w", err)
	}

	return &decimal, nil
}

func HexToBigInt(hexValue string) (*big.Int, error) {
	bigIntValue := new(big.Int)
	bigIntValue.SetString(hexValue, 0)

	return bigIntValue, nil
}

func WeiToEth(wei int64) float64 {
	const weiInEth = 1e18

	return float64(wei) / weiInEth
}

func WeiToEthUsingBigInt(weiValue *big.Int) *big.Float {
	const weiInEth = 1e18

	return new(big.Float).Quo(new(big.Float).SetInt(weiValue), big.NewFloat(weiInEth))
}

func FormatAmount(amount float64, decimalPlaces int) string {
	// Create the format string dynamically based on the number of decimal places
	format := fmt.Sprintf("%%.%df", decimalPlaces)

	return fmt.Sprintf(format, amount)
}

func FormatAmountBigFloat(value *big.Float, decimalPlaces int) string {
	format := fmt.Sprintf("%%.%df", decimalPlaces)

	return fmt.Sprintf(format, value)
}
