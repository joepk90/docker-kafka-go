package utils

import (
	"github.com/google/uuid"
)

var (
	accountNSUUID = uuid.NewSHA1(uuid.UUID{}, []byte("customer"))
)

func hashAccountID(accountNumber string) string {
	return uuid.NewSHA1(accountNSUUID, []byte(accountNumber)).String()
}
