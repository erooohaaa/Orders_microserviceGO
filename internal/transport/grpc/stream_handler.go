package grpc

import (
	"errors"
	"time"

	"Orders/internal/domain"
	"Orders/internal/usecase"

	api "github.com/erooohaaa/orders-generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderStreamHandler struct {
	api.UnimplementedOrderServiceServer
	uc *usecase.OrderUseCase
}

func NewOrderStreamHandler(uc *usecase.OrderUseCase) *OrderStreamHandler {
	return &OrderStreamHandler{uc: uc}
}

func (h *OrderStreamHandler) SubscribeToOrderUpdates(req *api.OrderRequest, stream api.OrderService_SubscribeToOrderUpdatesServer) error {
	if req.OrderId == "" {
		return status.Error(codes.InvalidArgument, "order_id is required")
	}

	lastStatus := ""
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-ticker.C:
			order, err := h.uc.GetOrder(stream.Context(), req.OrderId)
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					return status.Error(codes.NotFound, "order not found")
				}
				return status.Error(codes.Internal, "failed to fetch order")
			}

			if order.Status == lastStatus {
				continue
			}

			if err := stream.Send(&api.OrderStatusUpdate{
				OrderId:   order.ID,
				Status:    order.Status,
				UpdatedAt: timestamppb.New(time.Now()),
			}); err != nil {
				return err
			}
			lastStatus = order.Status
		}
	}
}
