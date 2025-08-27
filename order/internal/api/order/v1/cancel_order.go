package v1

import (
	"context"

	orderv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(
	ctx context.Context,
	params orderv1.CancelOrderParams,
) (orderv1.CancelOrderRes, error) {
	return &orderv1.InternalError{
		Code:    500,
		Message: "Internal server error.",
	}, nil
}
