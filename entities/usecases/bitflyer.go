package usecases

import (
	"fmt"

	"github.com/sn1w/capital-go/entities/infrastructures/bitflyer"
)

type AvaiableMarkets = []AvaiableMarket

type AvaiableMarket struct {
	ProductCode string
	Alias       string
	MarketType  string
}

type BoardPrice struct {
	Price float64
	Size  float64
}

type BoardPrices = []BoardPrice

type BoardInformation struct {
	MidPrice float64
	Bids     BoardPrices
	Asks     BoardPrices
}

type Balance struct {
	CurrencyCode string
	Amount       float64
	Available    float64
}

type Balances = []Balance

type BitFlyerUseCase struct {
	Client BitFlyerClient
}

type BitFlyerClient interface {
	GetAvaiableMarkets() (bitflyer.GetMarketsResponse, error)
	GetBoard(productCode string) (*bitflyer.BoardResponse, error)
	GetBalance() (bitflyer.GetBalancesResponse, error)
}

var _ BitFlyerClient = &bitflyer.BitFlyer{}

func (b *BitFlyerUseCase) ShowAvaiableMarkets() (AvaiableMarkets, error) {
	result, err := b.Client.GetAvaiableMarkets()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch markets: %w", err)
	}

	responses := make(AvaiableMarkets, 0, len(result))

	for _, v := range result {
		responses = append(responses, AvaiableMarket{
			ProductCode: v.ProductCode,
			Alias:       v.Alias,
			MarketType:  v.MarketType,
		})
	}

	return responses, nil
}

func (b *BitFlyerUseCase) GetBoard(productCode string) (BoardInformation, error) {
	result, err := b.Client.GetBoard(productCode)

	response := BoardInformation{}

	if err != nil {
		return response, fmt.Errorf("failed to fetch board: %w", err)
	}

	response.MidPrice = result.MidPrice

	for _, v := range result.Asks {
		response.Asks = append(response.Asks, BoardPrice{
			Price: v.Price,
			Size:  v.Size,
		})
	}

	for _, v := range result.Bids {
		response.Bids = append(response.Bids, BoardPrice{
			Price: v.Price,
			Size:  v.Size,
		})
	}

	return response, nil
}

func (b *BitFlyerUseCase) GetBalance() (Balances, error) {
	result, err := b.Client.GetBalance()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch balance: %w", err)
	}

	response := Balances{}

	for _, v := range result {
		response = append(response, Balance{
			Amount:       v.Amount,
			CurrencyCode: v.CurrencyCode,
			Available:    v.Available,
		})
	}

	return response, nil
}
