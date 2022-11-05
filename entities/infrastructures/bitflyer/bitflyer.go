package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sn1w/capital-go/config"
	cerror "github.com/sn1w/capital-go/error"
)

type BitFlyer struct {
	hc        *http.Client
	endPoint  string
	apiKey    string
	apiSecret string
}

func (b *BitFlyer) currentTimestamp() uint64 {
	return uint64(time.Now().Unix())
}

// NewBitFlyer returns a Default BitFlyer client.
func NewBitFlyer(cfg config.Config) *BitFlyer {
	return &BitFlyer{
		hc:        http.DefaultClient,
		endPoint:  "https://api.bitflyer.com",
		apiKey:    cfg.BitFlyerApiKey,
		apiSecret: cfg.BitFlyerApiSecret,
	}
}

// ChildOrderType represents Order Type used in SendChildOrder.
// https://lightning.bitflyer.com/docs?lang=ja&_gl=1*1rx0t7g*_ga*MjI5Nzg4NDM1LjE2NjQ1OTA4MTU.*_ga_3VYMQNCVSM*MTY2NzQ0MzUyMy4xNi4wLjE2Njc0NDM1MjMuNjAuMC4w#%E6%96%B0%E8%A6%8F%E6%B3%A8%E6%96%87%E3%82%92%E5%87%BA%E3%81%99
type ChildOrderType string

const (
	ChildOrderTypeLimit  ChildOrderType = "LIMIT"
	ChildOrderTypeMarket ChildOrderType = "MARKET"
)

// ChildOrderSide represents Side Types used in SendChildOrder.
// https://lightning.bitflyer.com/docs?lang=ja&_gl=1*1rx0t7g*_ga*MjI5Nzg4NDM1LjE2NjQ1OTA4MTU.*_ga_3VYMQNCVSM*MTY2NzQ0MzUyMy4xNi4wLjE2Njc0NDM1MjMuNjAuMC4w#%E6%96%B0%E8%A6%8F%E6%B3%A8%E6%96%87%E3%82%92%E5%87%BA%E3%81%99
type ChildOrderSide string

const (
	SideBuy  ChildOrderSide = "BUY"
	SideSell ChildOrderSide = "SELL"
)
const MiniuteToExpireDefault = 43200

// TimeInForceType represents TimeInForce param used in SendChildOrder.
// https://lightning.bitflyer.com/docs?lang=ja&_gl=1*1rx0t7g*_ga*MjI5Nzg4NDM1LjE2NjQ1OTA4MTU.*_ga_3VYMQNCVSM*MTY2NzQ0MzUyMy4xNi4wLjE2Njc0NDM1MjMuNjAuMC4w#%E6%96%B0%E8%A6%8F%E6%B3%A8%E6%96%87%E3%82%92%E5%87%BA%E3%81%99
type TimeInForceType string

const (
	TimeInForceGTC TimeInForceType = "GTC"
	TimeInForceIOC TimeInForceType = "IOC"
	TimeInForceFOK TimeInForceType = "FOK"
)

type SendOrderRequest struct {
	ProductCode    string          `json:"product_code"`
	ChildOrderType ChildOrderType  `json:"child_order_type"`
	Side           ChildOrderSide  `json:"side"`
	Price          float64         `json:"price"`
	Size           float64         `json:"size"`
	MinuteToExpire int             `json:"minute_to_expire"`
	TimeInForce    TimeInForceType `json:"time_in_force"`
}

type OrderResponse struct {
	ChildOrderAcceptanceId string `json:"child_order_acceptance_id"`
}

type MarketResponse struct {
	ProductCode string `json:"product_code"`
	MarketType  string `json:"market_type"`
	Alias       string `json:"alias"`
}

type GetMarketsResponse = []MarketResponse

type PriceResponse struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type PriceResponses = []PriceResponse

type BoardResponse struct {
	MidPrice float64        `json:"mid_price"`
	Bids     PriceResponses `json:"bids"`
	Asks     PriceResponses `json:"asks"`
}

type BalanceResponse struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

type GetBalancesResponse = []BalanceResponse

func request[REQ any, RES any](b *BitFlyer, method string, url string, body *REQ, useSecret bool) (*RES, error) {
	path := b.endPoint + url

	var requestBody []byte
	var err error

	if body != nil {
		requestBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshall request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(requestBody))

	if useSecret {
		timestamp := b.currentTimestamp()
		req.Header.Add("ACCESS-KEY", b.apiKey)
		req.Header.Add("ACCESS-TIMESTAMP", fmt.Sprint(timestamp))
		req.Header.Add("ACCESS-SIGN", b.generateSign(method, url, string(requestBody), timestamp))
	}

	if method == "POST" {
		req.Header.Add("Content-Type", "application/json")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to make request to %s, %w", url, err)
	}

	resp, err := b.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("got error from url %s: %w", url, err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error from read %s's  response: %w", url, err)
	}

	if resp.StatusCode >= 400 {
		rootError := cerror.ErrUnknown
		switch resp.StatusCode {
		case 400:
			rootError = cerror.ErrBadRequest
		case 401:
			rootError = cerror.ErrUnAuthorized
		case 404:
			rootError = cerror.ErrResourceNotFound
		}
		return nil, fmt.Errorf("%w: unexpected response %d. reason = %s", rootError, resp.StatusCode, string(res))
	}

	var result RES

	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall %s's response: %w", url, err)
	}

	return &result, nil
}

func getRequest[T any](b *BitFlyer, url string, useSecret bool) (*T, error) {
	return request[any, T](b, "GET", url, nil, useSecret)
}

// generateSign makes a token used to call HTTP Private API.
// https://lightning.bitflyer.com/docs?lang=ja#%E8%AA%8D%E8%A8%BC
func (b *BitFlyer) generateSign(method string, url string, requestBody string, unixTime uint64) string {
	seed := fmt.Sprintf("%d%s%s%s", unixTime, method, url, requestBody)
	mac := hmac.New(sha256.New, []byte(b.apiSecret))
	mac.Write([]byte(seed))

	return hex.EncodeToString(mac.Sum(nil))
}

// GetAvaiableMarkets represents an API call to `GET /v1/markets`.
//
// https://lightning.bitflyer.com/docs?lang=ja#%E3%83%9E%E3%83%BC%E3%82%B1%E3%83%83%E3%83%88%E3%81%AE%E4%B8%80%E8%A6%A7
func (b *BitFlyer) GetAvaiableMarkets() (GetMarketsResponse, error) {
	response, err := getRequest[GetMarketsResponse](b, "/v1/markets", false)
	if err != nil {
		return nil, err
	}

	return *response, nil
}

// GetBoard represents an API call to `GET /v1/board`.
//
// https://lightning.bitflyer.com/docs?lang=ja#%E6%9D%BF%E6%83%85%E5%A0%B1
func (b *BitFlyer) GetBoard(productCode string) (*BoardResponse, error) {
	url := fmt.Sprintf("/v1/board?product_code=%s", productCode)
	response, err := getRequest[BoardResponse](b, url, false)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetBalance represents an API call to `GET /v1/me/balance`.
//
// https://lightning.bitflyer.com/docs?lang=ja#%E8%B3%87%E7%94%A3%E6%AE%8B%E9%AB%98%E3%82%92%E5%8F%96%E5%BE%97
func (b *BitFlyer) GetBalance() (GetBalancesResponse, error) {
	response, err := getRequest[GetBalancesResponse](b, "/v1/me/getbalance", true)
	if err != nil {
		return nil, err
	}

	return *response, nil
}

// SendOrder represents an API call to `POST /v1/me/sendchildorder`.
//
// https://lightning.bitflyer.com/docs?lang=ja#%E6%96%B0%E8%A6%8F%E6%B3%A8%E6%96%87%E3%82%92%E5%87%BA%E3%81%99
func (b *BitFlyer) SendOrder(req SendOrderRequest) (*OrderResponse, error) {
	response, err := request[SendOrderRequest, OrderResponse](b, "POST", "/v1/me/sendchildorder", &req, true)
	if err != nil {
		return nil, err
	}

	return response, nil
}
