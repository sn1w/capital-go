package kabucom

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/sn1w/capital-go/entities/infrastructures/kabucom/autogen"
	cerror "github.com/sn1w/capital-go/error"
)

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
