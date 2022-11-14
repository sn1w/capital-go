package cli

import (
	"fmt"

	"github.com/sn1w/capital-go/entities/usecases"
)

type KabucomCLI struct {
	usecase usecases.KabucomUseCase
}

type Parameters struct {
	APIKey string
}

func NewKabucomCli(usecase usecases.KabucomUseCase) KabucomCLI {
	return KabucomCLI{usecase: usecase}
}

func (c *KabucomCLI) Authorization(pwd string) (string, error) {
	res, err := c.usecase.DoAuthorize(pwd)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("token: %s\n", res), nil
}

func (c *KabucomCLI) GetBalance(p Parameters) (string, error) {
	res, err := c.usecase.GetBalance(p.APIKey, usecases.KabucomBalanceRequest{})

	if err != nil {
		return "", err
	}

	output := "Type,Symbol,SymbolName,Exchange,ExName,Qty,Price,CurrentPrice,Value,ProfitLoss,ProfitLossRate\n"

	for _, v := range res {
		output += fmt.Sprintf("%d,%s,%s,%d,%s,%.0f,%.1f,%.1f,%.0f,%.0f,%.1f\n",
			v.AccountType, v.Symbol, v.SymbolName, v.Exchange, v.ExchangeName, v.LeavesQty, v.Price, v.CurrentPrice, v.Valuation, v.ProfitLoss, v.ProfitLossRate,
		)
	}

	return output, nil
}
