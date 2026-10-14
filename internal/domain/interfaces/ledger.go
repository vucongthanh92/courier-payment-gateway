package interfaces

import (
	"context"

	"github.com/vucongthanh92/courier-payment-gateway/internal/domain/entities"
	errHandler "github.com/vucongthanh92/go-base-utils/helper/http_error"
)

type LedgerQueryRepoI interface {
	GetWalletLiabilityAccount(context.Context, uint64) (*entities.LedgerAccount, *errHandler.ErrorBuilder)
	GetAccountByCode(context.Context, string) (*entities.LedgerAccount, *errHandler.ErrorBuilder)
}
type LedgerCommandRepoI interface {
	CreateJournalEntries(context.Context, *entities.LedgerJournal, []entities.LedgerEntry) *errHandler.ErrorBuilder
	CreateProviderTransaction(context.Context, *entities.ProviderTransaction) *errHandler.ErrorBuilder
}
