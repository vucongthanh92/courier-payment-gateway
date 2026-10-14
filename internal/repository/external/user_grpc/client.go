package user_grpc

import (
	"context"
	"fmt"

	"github.com/vucongthanh92/courier-payment-gateway/config"
	jwkpb "github.com/vucongthanh92/courier/shared/grpc/user-service/jwk/gen"
	errhandler "github.com/vucongthanh92/go-base-utils/helper/http_error"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGrpcClient interface {
	GetPublicKey(context.Context, string) (string, string, string, *errhandler.ErrorBuilder)
}
type grpcClient struct{ jwkSvc jwkpb.JWKServiceClient }

func NewGrpcClient(cfg *config.AppConfig) UserGrpcClient {
	conn, err := grpc.NewClient(cfg.Client.UserService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	return &grpcClient{jwkSvc: jwkpb.NewJWKServiceClient(conn)}
}
func (c *grpcClient) GetPublicKey(ctx context.Context, kid string) (string, string, string, *errhandler.ErrorBuilder) {
	resp, err := c.jwkSvc.GetPublicKey(ctx, &jwkpb.GetPublicKeyRequest{Kid: kid})
	if err != nil {
		return "", "", "", errhandler.InitErrorBuilder(ctx).SetStatus(500).SetLogError(fmt.Errorf("get public key: %w", err))
	}
	return resp.Kid, resp.PublicPem, resp.Alg, nil
}
