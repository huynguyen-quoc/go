package echo

import (
	"context"

	"github.com/huynguyen-quoc/go/common/server/endpoint"
	"github.com/huynguyen-quoc/go/common/server/metadata"
	"github.com/huynguyen-quoc/go/echo/dto"
)

// Endpoints ...
type Endpoints struct {
	EchoEndpoint endpoint.Endpoint
}

type ServiceHandler struct {
	Service Service
}

func (h ServiceHandler) MakeEchoEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (metadata.Metadata, interface{}, error) {
		userLoginReq := request.(*dto.EchoRequest)

		value, err := h.Service.Echo(ctx, userLoginReq.Value)

		return metadata.Metadata{}, &dto.EchoResponse{Value: value}, err
	}
}
