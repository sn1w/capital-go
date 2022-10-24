package bitflyer

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	cerror "github.com/sn1w/capital-go/error"
)

func TestBitFlyer_generateSign(t *testing.T) {
	type fields struct {
		endPoint  string
		apiSecret string
	}
	type args struct {
		method      string
		url         string
		requestBody string
		unixTime    uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "success",
			fields: fields{
				apiSecret: "hello",
			},
			args: args{
				method:      "GET",
				url:         "https://localhost/v1/me",
				requestBody: "",
				unixTime:    10000,
			},
			want: "bea433fb0f3611f5e5accfed6823eebe3b52abeca4373206e06653ae52131602",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BitFlyer{
				hc:        http.DefaultClient,
				endPoint:  tt.fields.endPoint,
				apiSecret: tt.fields.apiSecret,
			}
			if got := b.generateSign(tt.args.method, tt.args.url, tt.args.requestBody, tt.args.unixTime); got != tt.want {
				t.Errorf("BitFlyer.generateSign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitFlyer_GetAvaiableMarkets(t *testing.T) {
	type fields struct {
		endPoint        string
		apiKey          string
		apiSecret       string
		apiResponseCode int
		apiResponse     string
	}
	tests := []struct {
		name          string
		fields        fields
		want          GetMarketsResponse
		wantErr       bool
		expectedError error
	}{
		{
			name: "Success",
			fields: fields{
				endPoint:        "http://localhost",
				apiResponseCode: 200,
				apiResponse:     "[{\"product_code\": \"TST\", \"market_type\": \"Spot\", \"alias\": \"Testing\"}]",
			},
			want: GetMarketsResponse{
				MarketResponse{
					ProductCode: "TST",
					MarketType:  "Spot",
					Alias:       "Testing",
				},
			},
			wantErr: false,
		},
		{
			name: "Failed By 401",
			fields: fields{
				endPoint:        "http://localhost",
				apiResponseCode: 401,
				apiResponse:     "unauthorized",
			},
			wantErr:       true,
			expectedError: cerror.ErrUnAuthorized,
		},
		{
			name: "Failed By Unknown Status Code",
			fields: fields{
				endPoint:        "http://localhost",
				apiResponseCode: 429,
				apiResponse:     "I'm a teapot",
			},
			wantErr:       true,
			expectedError: cerror.ErrUnknown,
		},
	}
	for _, tt := range tests {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "http://localhost/v1/markets",
			httpmock.NewStringResponder(tt.fields.apiResponseCode, tt.fields.apiResponse),
		)

		t.Run(tt.name, func(t *testing.T) {
			b := BitFlyer{
				hc:        http.DefaultClient,
				endPoint:  tt.fields.endPoint,
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
			}
			got, err := b.GetAvaiableMarkets()
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyer.GetAvaiableMarkets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("BitFlyer.GetAvaiableMarkets() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitFlyer.GetAvaiableMarkets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitFlyer_GetBoard(t *testing.T) {
	type fields struct {
		endPoint        string
		apiKey          string
		apiSecret       string
		apiResponseCode int
		apiResponse     string
	}
	type args struct {
		productCode string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          *BoardResponse
		wantErr       bool
		expectedError error
	}{
		{
			name: "Success",
			fields: fields{
				endPoint:        "http://localhost",
				apiResponseCode: 200,
				apiResponse: `
					{
						"mid_price": 10502.345,
						"bids": [
							{ "size": 10.8, "price": 10280.3 },
							{ "size": 1.5, "price": 5472 }
						],
						"asks": [
							{ "size": 24, "price": 10601 },
							{ "size": 11.4, "price": 10482.5 }
						]
					}
				`,
			},
			args: args{
				productCode: "TST",
			},
			want: &BoardResponse{
				MidPrice: 10502.345,
				Asks: []PriceResponse{
					{
						Price: 10601,
						Size:  24,
					},
					{
						Price: 10482.5,
						Size:  11.4,
					},
				},
				Bids: []PriceResponse{
					{
						Price: 10280.3,
						Size:  10.8,
					},
					{
						Price: 5472,
						Size:  1.5,
					},
				},
			},
		},
		{
			name: "Unexpected Product Code",
			fields: fields{
				endPoint:        "http://localhost",
				apiResponseCode: 404,
				apiResponse:     "<html></html>",
			},
			args: args{
				productCode: "TST",
			},
			wantErr:       true,
			expectedError: cerror.ErrResourceNotFound,
		},
	}
	for _, tt := range tests {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "http://localhost/v1/board?product_code=TST",
			httpmock.NewStringResponder(tt.fields.apiResponseCode, tt.fields.apiResponse),
		)

		t.Run(tt.name, func(t *testing.T) {
			b := BitFlyer{
				hc:        http.DefaultClient,
				endPoint:  tt.fields.endPoint,
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
			}
			got, err := b.GetBoard(tt.args.productCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyer.GetBoard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("BitFlyer.GetBoard() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitFlyer.GetBoard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitFlyer_GetBalance(t *testing.T) {
	type fields struct {
		endPoint        string
		apiKey          string
		apiSecret       string
		apiResponseCode int
		apiResponse     string
	}
	tests := []struct {
		name          string
		fields        fields
		want          GetBalancesResponse
		wantErr       bool
		expectedError error
	}{
		{
			name: "Success",
			fields: fields{
				endPoint:        "http://localhost",
				apiResponseCode: 200,
				apiResponse: `
					[
						{"currency_code": "JPY", "amount": 10000.0, "available": 128.4}
					]
				`,
			},
			want: GetBalancesResponse{
				BalanceResponse{
					CurrencyCode: "JPY",
					Amount:       10000.0,
					Available:    128.4,
				},
			},
		},
		{
			name: "UnAuthorized",
			fields: fields{
				endPoint:        "http://localhost",
				apiResponseCode: 401,
				apiResponse:     "UnAuthorized",
			},
			wantErr:       true,
			expectedError: cerror.ErrUnAuthorized,
		},
	}
	for _, tt := range tests {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "http://localhost/v1/me/getbalance",
			httpmock.NewStringResponder(tt.fields.apiResponseCode, tt.fields.apiResponse),
		)

		t.Run(tt.name, func(t *testing.T) {
			b := BitFlyer{
				hc:        http.DefaultClient,
				endPoint:  tt.fields.endPoint,
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
			}
			got, err := b.GetBalance()
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyer.GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("BitFlyer.GetBalance() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitFlyer.GetBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}
