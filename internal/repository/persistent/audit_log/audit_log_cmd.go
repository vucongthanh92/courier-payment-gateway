package auditlog

import (
	"context"

	"github.com/vucongthanh92/courier-payment-gateway/database"
	"github.com/vucongthanh92/courier-payment-gateway/internal/domain/entities"
	"github.com/vucongthanh92/courier-payment-gateway/internal/domain/interfaces"
	errHandler "github.com/vucongthanh92/go-base-utils/helper/http_error"
	"github.com/vucongthanh92/go-base-utils/helper/transaction"
	"gorm.io/gorm"
)

type auditLogCmdRepository struct {
	writeDB *gorm.DB
}

func InitAuditLogCmdRepository(db *database.GormWriteDb) interfaces.AuditLogCommandRepoI {
	return &auditLogCmdRepository{
		writeDB: *db,
	}
}

func (r *auditLogCmdRepository) Create(ctx context.Context, auditLog *entities.AuditLog) *errHandler.ErrorBuilder {
	if err := transaction.RunnerFromCtx(ctx, r.writeDB).Create(auditLog).Error; err != nil {
		return errHandler.InitErrorBuilder(ctx).ValidateError(err)
	}

	return nil
}
