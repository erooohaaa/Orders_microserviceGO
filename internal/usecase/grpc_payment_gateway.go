package usecase

import (
	"context"

	api "github.com/erooohaaa/orders-generated"
	"google.golang.org/grpc"
)

type GRPCPaymentGateway struct {
	client api.PaymentServiceClient
}

func NewGRPCPaymentGateway(conn *grpc.ClientConn) *GRPCPaymentGateway {
	return &GRPCPaymentGateway{client: api.NewPaymentServiceClient(conn)}
}

func (g *GRPCPaymentGateway) Authorize(ctx context.Context, orderID string, amount int64) (string, string, error) {
	resp, err := g.client.ProcessPayment(ctx, &api.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	})
	if err != nil {
		return "", "", err
	}
	return resp.TransactionId, resp.Status, nil
}
