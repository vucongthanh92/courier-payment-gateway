package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/vucongthanh92/courier-payment-gateway/helper/constants"
	"github.com/vucongthanh92/courier-payment-gateway/internal/domain/models"
	cacheRepo "github.com/vucongthanh92/courier-payment-gateway/internal/repository/external/redis"
	user_grpc "github.com/vucongthanh92/courier-payment-gateway/internal/repository/external/user_grpc"
	errHandler "github.com/vucongthanh92/go-base-utils/helper/http_error"
)

// decodePublicKeyPEM decodes a PEM-encoded RSA public key.
func decodePublicKeyPEM(publicPEM string) (interface{}, *errHandler.ErrorBuilder) {
	pub, errJWT := jwt.ParseRSAPublicKeyFromPEM([]byte(publicPEM))
	if errJWT != nil {
		return nil, errHandler.InitErrorBuilder(context.Background()).
			SetStatus(http.StatusInternalServerError).
			SetError(models.ErrorDTO{Code: "invalid_public_key", Message: "Invalid public key"})
	}
	return pub, nil
}

// ResolvePublicKey fetches the public key by kid, first checking the cache, then falling back to the JWK client.
func ResolvePublicKey(ctx context.Context, cache cacheRepo.JWKCacheRepo, client user_grpc.UserGrpcClient, kid string) (any, *errHandler.ErrorBuilder) {
	if kid != "" && cache != nil {
		if cached, err := cache.GetByKid(ctx, kid); err == nil && cached != nil && cached.PublicPEM != "" {
			return decodePublicKeyPEM(cached.PublicPEM)
		}
	}

	respKid, respPublicPEM, respAlg, errBuilder := client.GetPublicKey(ctx, kid)
	if errBuilder != nil {
		return nil, errBuilder
	}

	pub, errBuilder := decodePublicKeyPEM(respPublicPEM)
	if errBuilder != nil {
		return nil, errBuilder
	}

	if cache != nil && respKid != "" {
		_ = cache.SetByKid(ctx, cacheRepo.JWKCacheEntry{Kid: respKid, PublicPEM: respPublicPEM, Alg: respAlg}, constants.Time_Cache_15_minutes)
	}
	return pub, nil
}
