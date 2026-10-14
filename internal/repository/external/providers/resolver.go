package providers

import (
	"context"
	"net/http"
	"strings"

	"github.com/vucongthanh92/courier-payment-gateway/internal/domain/interfaces"
	"github.com/vucongthanh92/courier-payment-gateway/internal/domain/models"
	"github.com/vucongthanh92/courier-payment-gateway/internal/repository/external/providers/sepay"
	errHandler "github.com/vucongthanh92/go-base-utils/helper/http_error"
)

type resolver struct {
	providers map[string]interfaces.PaymentGateway
}

func InitPaymentGatewayResolver(sePayProvider *sepay.Provider) interfaces.PaymentGatewayResolverI {
	return &resolver{
		providers: map[string]interfaces.PaymentGateway{
			sePayProvider.Name(): sePayProvider,
		},
	}
}

func (r *resolver) Resolve(ctx context.Context, providerName string) (interfaces.PaymentGateway, *errHandler.ErrorBuilder) {
	providerName = strings.ToLower(strings.TrimSpace(providerName))
	provider, exists := r.providers[providerName]
	if exists {
		return provider, nil
	}

	return nil, errHandler.InitErrorBuilder(ctx).
		SetStatus(http.StatusUnprocessableEntity).
		SetError(models.ErrorDTO{
			Code:    "provider_not_available",
			Field:   "provider-name",
			Message: "the requested payment provider is not available",
		})
}
