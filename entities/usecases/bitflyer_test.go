package usecases

import (
	"errors"
	"reflect"
	"testing"

	"github.com/sn1w/capital-go/entities/infrastructures/bitflyer"
	cerror "github.com/sn1w/capital-go/error"
)

type mockedBitFlyerClient struct {
	bitflyer.BitFlyer
	getMarket  func() (bitflyer.GetMarketsResponse, error)
	getBoard   func(pc string) (*bitflyer.BoardResponse, error)
	getBalance func() (bitflyer.GetBalancesResponse, error)
	sendOrder  func(req bitflyer.SendOrderRequest) (*bitflyer.OrderResponse, error)
}

func (m mockedBitFlyerClient) GetAvaiableMarkets() (bitflyer.GetMarketsResponse, error) {
	return m.getMarket()
}
func (m mockedBitFlyerClient) GetBoard(productCode string) (*bitflyer.BoardResponse, error) {
	return m.getBoard(productCode)
}
func (m mockedBitFlyerClient) GetBalance() (bitflyer.GetBalancesResponse, error) {
	return m.getBalance()
}
func (m mockedBitFlyerClient) SendOrder(req bitflyer.SendOrderRequest) (*bitflyer.OrderResponse, error) {
	return m.sendOrder(req)
}

func TestBitFlyerUseCase_ShowAvaiableMarkets(t *testing.T) {
	type fields struct {
		Client BitFlyerClient
	}
	tests := []struct {
		name        string
		fields      fields
		want        AvaiableMarkets
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success",
			fields: fields{
				Client: &mockedBitFlyerClient{
					getMarket: func() (bitflyer.GetMarketsResponse, error) {
						return bitflyer.GetMarketsResponse{
							bitflyer.MarketResponse{
								ProductCode: "TEST_PRODUCT",
								MarketType:  "Market",
								Alias:       "TST",
							},
						}, nil
					},
				},
			},
			want: []AvaiableMarket{
				{
					ProductCode: "TEST_PRODUCT", MarketType: "Market", Alias: "TST",
				},
			},
		},
		{
			name: "got error",
			fields: fields{
				Client: &mockedBitFlyerClient{
					getMarket: func() (bitflyer.GetMarketsResponse, error) {
						return nil, cerror.ErrUnAuthorized
					},
				},
			},
			wantErr:     true,
			expectedErr: cerror.ErrUnAuthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BitFlyerUseCase{
				Client: tt.fields.Client,
			}
			got, err := b.ShowAvaiableMarkets()
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyerUseCase.ShowAvaiableMarkets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && !errors.Is(err, tt.expectedErr) {
				t.Errorf("BitFlyerUseCase.ShowAvaiableMarkets() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitFlyerUseCase.ShowAvaiableMarkets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitFlyerUseCase_GetBoard(t *testing.T) {
	type fields struct {
		Client BitFlyerClient
	}
	type args struct {
		productCode string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        BoardInformation
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success",
			fields: fields{
				Client: mockedBitFlyerClient{
					getBoard: func(pc string) (*bitflyer.BoardResponse, error) {
						return &bitflyer.BoardResponse{
							MidPrice: 100.5,
							Bids:     []bitflyer.PriceResponse{{Price: 10.2, Size: 5.4}},
							Asks:     []bitflyer.PriceResponse{{Price: 9.8, Size: 6.2}},
						}, nil
					},
				},
			},
			want: BoardInformation{
				MidPrice: 100.5,
				Asks:     []BoardPrice{{Price: 9.8, Size: 6.2}},
				Bids:     []BoardPrice{{Price: 10.2, Size: 5.4}},
			},
		},
		{
			name: "error",
			fields: fields{
				Client: mockedBitFlyerClient{
					getBoard: func(pc string) (*bitflyer.BoardResponse, error) {
						return nil, cerror.ErrResourceNotFound
					},
				},
			},
			wantErr:     true,
			expectedErr: cerror.ErrResourceNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BitFlyerUseCase{
				Client: tt.fields.Client,
			}
			got, err := b.GetBoard(tt.args.productCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyerUseCase.GetBoard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && !errors.Is(err, tt.expectedErr) {
				t.Errorf("BitFlyerUseCase.GetBoard() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitFlyerUseCase.GetBoard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitFlyerUseCase_GetBalance(t *testing.T) {
	type fields struct {
		Client BitFlyerClient
	}
	tests := []struct {
		name        string
		fields      fields
		want        Balances
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success",
			fields: fields{
				Client: mockedBitFlyerClient{
					getBalance: func() (bitflyer.GetBalancesResponse, error) {
						return []bitflyer.BalanceResponse{
							{
								CurrencyCode: "JPY", Amount: 10000, Available: 39.48,
							},
						}, nil
					},
				},
			},
			want: []Balance{
				{
					CurrencyCode: "JPY",
					Amount:       10000,
					Available:    39.48,
				},
			},
		},
		{
			name: "error",
			fields: fields{
				Client: mockedBitFlyerClient{
					getBalance: func() (bitflyer.GetBalancesResponse, error) {
						return nil, cerror.ErrUnknown
					},
				},
			},
			wantErr:     true,
			expectedErr: cerror.ErrUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BitFlyerUseCase{
				Client: tt.fields.Client,
			}
			got, err := b.GetBalance()
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyerUseCase.GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && !errors.Is(err, tt.expectedErr) {
				t.Errorf("BitFlyerUseCase.GetBalance() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitFlyerUseCase.GetBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitFlyerUseCase_CreateOrder(t *testing.T) {
	type fields struct {
		Client BitFlyerClient
	}
	type args struct {
		req OrderCreate
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        *OrderInformation
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success",
			fields: fields{
				Client: mockedBitFlyerClient{
					sendOrder: func(req bitflyer.SendOrderRequest) (*bitflyer.OrderResponse, error) {
						if req.ProductCode == "" ||
							req.Size == 0 ||
							req.Price == 0 {
							return nil, cerror.ErrBadRequest
						}
						return &bitflyer.OrderResponse{ChildOrderAcceptanceId: "test_id"}, nil
					},
				},
			},
			args: args{
				req: OrderCreate{
					Price:       10242,
					Size:        10.5,
					ProductCode: "TEST_TOKEN",
				},
			},
			want: &OrderInformation{OrderAcceeptanceId: "test_id"},
		},
		{
			name: "error",
			fields: fields{
				Client: mockedBitFlyerClient{
					sendOrder: func(req bitflyer.SendOrderRequest) (*bitflyer.OrderResponse, error) {
						return nil, cerror.ErrUnknown
					},
				},
			},
			wantErr:     true,
			expectedErr: cerror.ErrUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BitFlyerUseCase{
				Client: tt.fields.Client,
			}
			got, err := b.CreateOrder(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyerUseCase.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && !errors.Is(err, tt.expectedErr) {
				t.Errorf("BitFlyerUseCase.CreateOrder() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitFlyerUseCase.CreateOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
