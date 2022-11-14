package usecases

import (
	"fmt"

	"github.com/sn1w/capital-go/entities/infrastructures/kabucom"
	"github.com/sn1w/capital-go/entities/infrastructures/kabucom/autogen"
)

func valueOrDefault[T int32 | float64 | string](value *T) T {
	var val T
	if value == nil {
		return val
	}

	return *value
}

type KabucomUseCase struct {
	client KabucomClient
}

type KabucomClient interface {
	GetToken(pwd string) (string, error)
	GetPosition(req kabucom.GetPositionRequest) ([]autogen.PositionsSuccess, error)
}

type KabucomBalanceRequest struct {
	Product string
}

type KabucomBalance struct {
	AccountType     int
	Commision       float64
	CommisionTax    float64
	CurrentPrice    float64
	Exchange        int
	ExchangeName    string
	ExecutionDay    int
	ExecutionID     string
	Expenses        float64
	ExpireDay       int
	HoldQty         float64
	LeavesQty       float64
	MarginTradeTyoe int
	Price           float64
	ProfitLoss      float64
	ProfitLossRate  float64
	SecurityType    int
	Side            string
	Symbol          string
	SymbolName      string
	Valuation       float64
}

type KabucomBalances = []KabucomBalance

func NewKabucomUseCase(client KabucomClient) KabucomUseCase {
	return KabucomUseCase{
		client: client,
	}
}

var _ KabucomClient = &kabucom.KabucomClient{}

func (k *KabucomUseCase) DoAuthorize(pwd string) (string, error) {
	result, err := k.client.GetToken(pwd)

	if err != nil {
		return "", fmt.Errorf("failed to authorize: %w", err)
	}

	return result, nil
}

func (k *KabucomUseCase) GetBalance(token string, req KabucomBalanceRequest) (KabucomBalances, error) {
	result, err := k.client.GetPosition(kabucom.GetPositionRequest{
		Product: req.Product,
		APIKey:  token,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	var response KabucomBalances

	for _, v := range result {
		response = append(response, KabucomBalance{
			AccountType:     int(*v.AccountType),
			Commision:       valueOrDefault(v.Commission),
			CommisionTax:    valueOrDefault(v.CommissionTax),
			CurrentPrice:    valueOrDefault(v.CurrentPrice),
			Exchange:        int(valueOrDefault(v.Exchange)),
			ExchangeName:    valueOrDefault(v.ExchangeName),
			ExecutionDay:    int(valueOrDefault(v.ExecutionDay)),
			ExecutionID:     valueOrDefault(v.ExecutionID),
			Expenses:        valueOrDefault(v.Expenses),
			ExpireDay:       int(valueOrDefault(v.ExpireDay)),
			HoldQty:         valueOrDefault(v.HoldQty),
			LeavesQty:       valueOrDefault(v.LeavesQty),
			MarginTradeTyoe: int(valueOrDefault(v.MarginTradeType)),
			Price:           valueOrDefault(v.Price),
			ProfitLoss:      valueOrDefault(v.ProfitLoss),
			ProfitLossRate:  valueOrDefault(v.ProfitLossRate),
			SecurityType:    int(valueOrDefault(v.SecurityType)),
			Side:            valueOrDefault(v.Side),
			Symbol:          valueOrDefault(v.Symbol),
			SymbolName:      valueOrDefault(v.SymbolName),
			Valuation:       valueOrDefault(v.Valuation),
		})
	}

	return response, nil
}
