package nerdweb_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/app-nerds/nerdweb/v2"
	"github.com/sirupsen/logrus"
)

func TestRealIP(t *testing.T) {
	type args struct {
		r *http.Request
	}

	withHeaders := make(http.Header)
	withHeaders["X-Forwarded-For"] = []string{"127.0.0.2"}

	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "Returns RemoteAddr when X-Forwarded-For is not present",
			want: "127.0.0.1",
			args: args{
				r: &http.Request{
					RemoteAddr: "127.0.0.1",
					Header:     http.Header{},
				},
			},
		},
		{
			name: "Returns X-Forwarded-For when present",
			want: "127.0.0.2",
			args: args{
				r: &http.Request{
					RemoteAddr: "127.0.0.1",
					Header:     withHeaders,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nerdweb.RealIP(tt.args.r)

			if got != tt.want {
				t.Errorf("want %s, got %s", tt.want, got)
			}
		})
	}
}

func TestValidateHTTPMethod(t *testing.T) {
	type args struct {
		r              *http.Request
		w              *httptest.ResponseRecorder
		expectedMethod string
	}

	logger := logrus.New().WithField("who", "testing")

	tests := []struct {
		name             string
		wantErr          bool
		wantErrorMessage string
		args             args
	}{
		{
			name:             "Returns no error when method is valid",
			wantErr:          false,
			wantErrorMessage: "",
			args: args{
				r: &http.Request{
					Method: "POST",
				},
				w:              nil,
				expectedMethod: "POST",
			},
		},
		{
			name:             "Returns an error and writes a message to the write when the method is invalid",
			wantErr:          true,
			wantErrorMessage: `{"message":"method not allowed"}`,
			args: args{
				r: &http.Request{
					Method: "GET",
				},
				w:              httptest.NewRecorder(),
				expectedMethod: "POST",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := nerdweb.ValidateHTTPMethod(tt.args.r, tt.args.w, tt.args.expectedMethod, logger)

			if !tt.wantErr && gotErr != nil {
				t.Errorf("did not expect an error!")
			}

			if tt.wantErr {
				if gotErr == nil {
					t.Errorf("wanted an error")
				}

				if strings.TrimSpace(tt.args.w.Body.String()) != strings.TrimSpace(tt.wantErrorMessage) {
					t.Errorf("wanted error '%s', got '%s'", strings.TrimSpace(tt.wantErrorMessage), strings.TrimSpace(tt.args.w.Body.String()))
				}
			}
		})
	}
}
