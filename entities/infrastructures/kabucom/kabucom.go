package kabucom

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sn1w/capital-go/config"
	"github.com/sn1w/capital-go/entities/infrastructures/kabucom/autogen"
	cerror "github.com/sn1w/capital-go/error"
)

type KabucomClient struct {
	client *autogen.ClientWithResponses
}

type GetPositionRequest struct {
	APIKey  string
	Product string
}

func NewKabucomClient(cfg config.Config) *KabucomClient {
	c, err := autogen.NewClientWithResponses(cfg.KabucomAPIHost,
		autogen.WithHTTPClient(http.DefaultClient))
	if err != nil {
		panic(err)
	}
	return &KabucomClient{
		client: c,
	}
}

func (c *KabucomClient) GetToken(password string) (string, error) {
	ctx := context.Background()
	res, err := c.client.TokenPostWithResponse(ctx, autogen.RequestToken{
		APIPassword: password,
	})

	if err != nil {
		return "", err
	}
	defer res.HTTPResponse.Body.Close()

	if res.StatusCode() == 401 {
		return "", fmt.Errorf("unexpected error %w, body = %s", cerror.ErrUnAuthorized, res.Body)
	}

	if res.JSON200 == nil || res.JSON200.Token == nil {
		return "", fmt.Errorf("unexpected error %w, body = %s", cerror.ErrUnknownResponseFormat, res.Body)
	}

	return *res.JSON200.Token, nil
}

func (c *KabucomClient) GetPosition(req GetPositionRequest) ([]autogen.PositionsSuccess, error) {
	ctx := context.Background()
	res, err := c.client.PositionsGetWithResponse(ctx, &autogen.PositionsGetParams{
		Product: nil,
		Symbol:  nil,
		Side:    nil,
		Addinfo: nil,
		XAPIKEY: req.APIKey,
	})

	if err != nil {
		return nil, fmt.Errorf("unexpected error %w, reason = %s", cerror.ErrUnknown, err.Error())
	}
	defer res.HTTPResponse.Body.Close()

	if res.StatusCode() == 401 {
		return nil, fmt.Errorf("unexpected error %w, body = %s", cerror.ErrUnAuthorized, res.Body)
	}

	if res.JSON200 == nil {
		return nil, fmt.Errorf("unexpected error %w, body = %s", cerror.ErrUnknownResponseFormat, res.Body)
	}

	return *res.JSON200, nil
}
