package kabucom

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/sn1w/capital-go/entities/infrastructures/kabucom/autogen"
	cerror "github.com/sn1w/capital-go/error"
)

func ptr[T any](v T) *T {
	return &v
}

type MockedRoundTrip struct {
	RoundTripCb func(*http.Request) (*http.Response, error)
}

func (t *MockedRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripCb(req)
}

func NewMockedClient(res *http.Response) *autogen.ClientWithResponses {
	c, _ := autogen.NewClientWithResponses("localhost", autogen.WithHTTPClient(
		&http.Client{
			Transport: &MockedRoundTrip{
				RoundTripCb: func(r *http.Request) (*http.Response, error) {
					return res, nil
				},
			},
		},
	))
	return c
}

func TestKabucomClient_GetToken(t *testing.T) {
	type fields struct {
		client *autogen.ClientWithResponses
	}
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		{
			name: "Success",
			fields: fields{
				client: NewMockedClient(
					&http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Content-Type": []string{"application/json"},
						},
						Body: io.NopCloser(strings.NewReader(`
						{
							"ResultCode": 0,
							"Token": "Hello Token"
						}
						`)),
					},
				),
			},
			args: args{
				password: "HELP Wanted",
			},
			want: "Hello Token",
		},
		{
			name: "UnAuthorized",
			fields: fields{
				client: NewMockedClient(
					&http.Response{
						StatusCode: 401,
						Header: http.Header{
							"Content-Type": []string{"application/json"},
						},
						Body: io.NopCloser(strings.NewReader(`
						{
							"Code": 4001001,
							"Message": "内部エラー"
						}
						`)),
					},
				),
			},
			args: args{
				password: "HELP wanted",
			},
			wantErr: cerror.ErrUnAuthorized,
		},
		{
			name: "Unknown Response Format",
			fields: fields{
				client: NewMockedClient(
					&http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Content-Type": []string{"application/json"},
						},
						Body: io.NopCloser(strings.NewReader(`
						{
							"Message": "OK Computer"
						}
						`)),
					},
				),
			},
			args: args{
				password: "HELP wanted",
			},
			wantErr: cerror.ErrUnknownResponseFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := KabucomClient{
				client: tt.fields.client,
			}
			got, err := c.GetToken(tt.args.password)
			if err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("KabucomClient.GetToken() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("KabucomClient.GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKabucomClient_GetPosition(t *testing.T) {
	type fields struct {
		client *autogen.ClientWithResponses
	}
	type args struct {
		req GetPositionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []autogen.PositionsSuccess
		wantErr error
	}{
		{
			name: "Success",
			fields: fields{
				client: NewMockedClient(
					&http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Content-Type": []string{"application/json"},
						},
						Body: io.NopCloser(strings.NewReader(`
						[
							{
							"ExecutionID": "E20220404xxxxx",
							"AccountType": 4,
							"Symbol": "8306",
							"SymbolName": "三菱ＵＦＪフィナンシャル・グループ",
							"Exchange": 1,
							"ExchangeName": "東証プ",
							"ExecutionDay": 20220404,
							"Price": 704,
							"LeavesQty": 500,
							"HoldQty": 0,
							"Side": "1",
							"Expenses": 0,
							"Commission": 1620,
							"CommissionTax": 162,
							"ExpireDay": 20220404,
							"MarginTradeType": 1,
							"CurrentPrice": 414.5,
							"Valuation": 207250,
							"ProfitLoss": 144750,
							"ProfitLossRate": 41.12215909090909
							}
						]
						`)),
					},
				),
			},
			args: args{},
			want: []autogen.PositionsSuccess{
				{
					ExecutionID:     ptr("E20220404xxxxx"),
					AccountType:     ptr(int32(4)),
					Symbol:          ptr("8306"),
					SymbolName:      ptr("三菱ＵＦＪフィナンシャル・グループ"),
					Exchange:        ptr(int32(1)),
					ExchangeName:    ptr("東証プ"),
					ExecutionDay:    ptr(int32(20220404)),
					Price:           ptr(704.0),
					LeavesQty:       ptr(500.0),
					HoldQty:         ptr(0.0),
					Side:            ptr("1"),
					Expenses:        ptr(0.0),
					Commission:      ptr(1620.0),
					CommissionTax:   ptr(162.0),
					ExpireDay:       ptr(int32(20220404)),
					MarginTradeType: ptr(int32(1)),
					CurrentPrice:    ptr(414.5),
					Valuation:       ptr(207250.0),
					ProfitLoss:      ptr(144750.0),
					ProfitLossRate:  ptr(41.12215909090909),
				},
			},
		},
		{
			name: "Invalid Response",
			fields: fields{
				client: NewMockedClient(
					&http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Content-Type": []string{"application/json"},
						},
						Body: io.NopCloser(strings.NewReader(`{"reason":"error"}`)),
					}),
			},
			args:    args{},
			wantErr: cerror.ErrUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &KabucomClient{
				client: tt.fields.client,
			}
			got, err := c.GetPosition(tt.args.req)
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("KabucomClient.GetPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KabucomClient.GetPosition() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
