package usecases

import (
	"fmt"

	"github.com/sn1w/capital-go/entities/infrastructures/kabucom"
)

type KabucomUseCase struct {
	client KabucomClient
}

type KabucomClient interface {
	GetToken(pwd string) (string, error)
}

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
