package v1

import (
	"context"
	"github.com/huynguyen-quoc/go/investment/pb/v1"
)

type InvestmentService struct {
}

func (i InvestmentService) OpenInvestment(ctx context.Context, request *v1.InvestmentRequest) (*v1.InvestmentResponse, error) {
	panic("implement me")
}
