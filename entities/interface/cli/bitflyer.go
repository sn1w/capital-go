package cli

import (
	"fmt"
	"sort"

	"github.com/sn1w/capital-go/entities/usecases"
)

// BitFlyerCLI executes a BitFlyer's useCase in a form suitable for CLI input/output.
type BitFlyerCLI struct {
	useCase usecases.BitFlyerUseCase
}

func NewBitFlyerCli(usecase usecases.BitFlyerUseCase) BitFlyerCLI {
	return BitFlyerCLI{useCase: usecase}
}

type CreateOrderArgument struct {
	Size        float64
	Price       float64
	ProductCode string
	Buy         bool
}

func (c *BitFlyerCLI) GetAvaiableMarkets() (string, error) {
	res, err := c.useCase.ShowAvaiableMarkets()

	if err != nil {
		return "", err
	}

	output := "Product Code, Alias, Market Type\n"

	for _, v := range res {
		output += fmt.Sprintf("%s, %s, %s\n", v.ProductCode, v.Alias, v.MarketType)
	}

	return output, nil
}

func (c *BitFlyerCLI) GetBoard(productCode string) (string, error) {
	res, err := c.useCase.GetBoard(productCode)

	if err != nil {
		return "", err
	}

	sort.Slice(res.Asks, func(i, j int) bool { return res.Asks[i].Price > res.Asks[j].Price })

	output := fmt.Sprintf("mid_price: %f\n", res.MidPrice)
	output += "\nAsk\n===========\n"
	askLen := len(res.Asks)

	for _, v := range res.Asks[(askLen - 10):] {
		output += fmt.Sprintf("Price: %f, Size: %f\n", v.Price, v.Size)
	}
	output += "\nBid\n===========\n"
	for _, v := range res.Bids[:10] {
		output += fmt.Sprintf("Price: %f, Size: %f\n", v.Price, v.Size)
	}

	return output, nil
}

func (c *BitFlyerCLI) GetBalance() (string, error) {
	res, err := c.useCase.GetBalance()
	if err != nil {
		return "", err
	}

	output := ""
	for _, v := range res {
		output += fmt.Sprintf("%s, %f, %f\n", v.CurrencyCode, v.Amount, v.Available)
	}

	return output, nil
}

func (c *BitFlyerCLI) CreateOrder(arg CreateOrderArgument) (string, error) {
	orderReq := usecases.OrderCreate{
		Size:        arg.Size,
		Price:       arg.Price,
		Buy:         arg.Buy,
		ProductCode: arg.ProductCode,
	}

	res, err := c.useCase.CreateOrder(orderReq)
	if err != nil {
		return "", err
	}

	output := res.OrderAcceeptanceId

	return output, nil
}
