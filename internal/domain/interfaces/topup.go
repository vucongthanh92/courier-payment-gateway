package interfaces

import (
	"context"

	"github.com/vucongthanh92/courier-payment-gateway/internal/domain/entities"
	"github.com/vucongthanh92/courier-payment-gateway/internal/domain/models"
	errHandler "github.com/vucongthanh92/go-base-utils/helper/http_error"
)

type TopUpQueryRepoI interface {
	GetForUpdateByPaymentCode(context.Context, string, string) (*entities.TopUpIntent, *errHandler.ErrorBuilder)
}

type TopUpCommandRepoI interface {
	Create(context.Context, *entities.TopUpIntent) *errHandler.ErrorBuilder
	MarkSucceeded(context.Context, uint64, string) *errHandler.ErrorBuilder
}

type TopUpServiceI interface {
	CreateTopUp(context.Context, models.CreateTopUpRequest) (models.CheckoutInstruction, bool, *errHandler.ErrorBuilder)
}
