package utils

import (
	"crypto/rand"
	"errors"
	"fmt"

	baseUtils "github.com/vucongthanh92/go-base-utils/helper/utils"
)

const paymentCodeAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type SnowflakeConfig = baseUtils.SnowflakeConfig

func InitSnowflake(cfg SnowflakeConfig) {
	baseUtils.InitSnowflake(cfg)
}

func NewSnowflakeID() (uint64, error) {
	return baseUtils.NewSnowflakeID()
}

// NewSePayPaymentCode generates a Test mode-compatible payment code in the
// format COUR[A-Za-z0-9]{6,8}. The current format uses eight suffix characters.
func NewSePayPaymentCode(prefix string, suffixLength int) (string, error) {
	if prefix == "" || suffixLength < 6 || suffixLength > 8 {
		return "", errors.New("invalid SePay payment-code format")
	}

	bytes := make([]byte, suffixLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate SePay payment code: %w", err)
	}
	for index, value := range bytes {
		bytes[index] = paymentCodeAlphabet[int(value)%len(paymentCodeAlphabet)]
	}

	return prefix + string(bytes), nil
}
