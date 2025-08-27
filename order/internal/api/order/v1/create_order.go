package v1

import (
	"context"

	orderv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(
	ctx context.Context,
	req *orderv1.OrderCreateRequest,
) (orderv1.CreateOrderRes, error) {
	return &orderv1.InternalError{
		Code:    500,
		Message: "Internal server error.",
	}, nil
}
