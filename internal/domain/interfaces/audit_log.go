package interfaces

import (
	"context"

	"github.com/vucongthanh92/courier-payment-gateway/internal/domain/entities"
	"github.com/vucongthanh92/courier-payment-gateway/internal/domain/models"
	errHandler "github.com/vucongthanh92/go-base-utils/helper/http_error"
)

type AuditLogCommandRepoI interface {
	Create(context.Context, *entities.AuditLog) *errHandler.ErrorBuilder
}

type AuditLogServiceI interface {
	Write(context.Context, models.AuditLogInput) *errHandler.ErrorBuilder
}
