package v1

import (
	"context"
	"net/http"

	orderv1 "github.com/korchizhinskiy/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) NewError(
	_ context.Context,
	err error,
) *orderv1.GenericErrorStatusCode {
	return &orderv1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderv1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}
