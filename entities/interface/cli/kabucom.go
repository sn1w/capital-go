package cli

import (
	"fmt"

	"github.com/sn1w/capital-go/entities/usecases"
)

type KabucomCLI struct {
	usecase usecases.KabucomUseCase
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
