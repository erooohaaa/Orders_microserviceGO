package grpc

import (
	"Orders/internal/usecase"
	"github.com/erooohaaa/orders-generated/api"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type OrderStreamHandler struct {
	api.UnimplementedOrderServiceServer
	uc *usecase.OrderUseCase
}

func NewOrderStreamHandler(uc *usecase.OrderUseCase) *OrderStreamHandler {
	return &OrderStreamHandler{uc: uc}
}

func (h *OrderStreamHandler) SubscribeToOrderUpdates(req *api.OrderRequest, stream api.OrderService_SubscribeToOrderUpdatesServer) error {
	lastStatus := ""

	for {
		// Проверяем статус в базе каждые 2 секунды (простая реализация стриминга)
		order, err := h.uc.GetOrder(stream.Context(), req.OrderId)
		if err != nil {
			return err
		}

		// Шлем обновление только если статус реально изменился
		if order.Status != lastStatus {
			if err := stream.Send(&api.OrderStatusUpdate{
				OrderId:   order.ID,
				Status:    order.Status,
				UpdatedAt: timestamppb.New(time.Now()),
			}); err != nil {
				return err
			}
			lastStatus = order.Status
		}

		time.Sleep(2 * time.Second)
	}
}
