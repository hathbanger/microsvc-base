package microsvc

import (
	"context"
	"errors"

	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

func (s service) Foo(
	ctx context.Context,
	req models.FooRequest,
) (models.FooResponse, error) {
	var (
		response models.FooResponse
		err      error
	)

	if req.Str == "" {
		err = errors.New("no string was passed")
		s.logger.Log("err", "boo")
		return response, err
	}

	product := req.Str + "bar"
	response.Res = product
	s.logger.Log("response", response.Res)

	return response, err
}
