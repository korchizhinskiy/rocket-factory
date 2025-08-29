package v1

import (
	"context"
	"errors"

	"github.com/korchizhinskiy/rocket-factory/order/internal/api/converter/response"
	"github.com/korchizhinskiy/rocket-factory/order/internal/model"
	orderv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByID(
	ctx context.Context,
	params orderv1.GetOrderByIDParams,
) (orderv1.GetOrderByIDRes, error) {
	order, err := a.orderService.GetOrderByID(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderv1.NotFoundError{
				Code:    404,
				Message: "Order was not found.",
			}, nil
		}
		return &orderv1.InternalError{
			Code:    500,
			Message: "Internal server error.",
		}, nil
	}

	return response.ConvertModelToResponse(order), nil
}
