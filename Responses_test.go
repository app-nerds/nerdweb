package nerdweb_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/app-nerds/nerdweb"
	"github.com/sirupsen/logrus"
)

func TestReadJSONBody(t *testing.T) {
	type sampleStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type args struct {
		r       *http.Request
		gotDest *sampleStruct
	}

	successStruct := &sampleStruct{Name: "Adam", Age: 10}
	successStructJSON, _ := json.Marshal(successStruct)
	successRequestBody := io.NopCloser(bytes.NewReader(successStructJSON))
	badRequestBody := io.NopCloser(strings.NewReader(`{"bad"`))

	tests := []struct {
		name     string
		wantErr  bool
		wantDest *sampleStruct
		args     args
	}{
		{
			name:     "Returns nil error, and populates the dest variable upon success",
			wantErr:  false,
			wantDest: successStruct,
			args: args{
				r: &http.Request{
					Body: successRequestBody,
				},
			},
		},
		{
			name:     "Returns an error when unmarshing fails",
			wantErr:  true,
			wantDest: nil,
			args: args{
				r: &http.Request{
					Body: badRequestBody,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := nerdweb.ReadJSONBody(tt.args.r, &tt.args.gotDest)

			if tt.wantErr && gotErr == nil {
				t.Errorf("wanted error")
			}

			if tt.wantDest != nil {
				if !reflect.DeepEqual(tt.args.gotDest, tt.wantDest) {
					t.Errorf("want:\n%#v\ngot:\n%#v", tt.wantDest, tt.args.gotDest)

				}
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	type sampleStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type args struct {
		w      *httptest.ResponseRecorder
		status int
		value  interface{}
	}

	logger := logrus.New().WithField("who", "testing")

	tests := []struct {
		name            string
		wantStatus      int
		wantContentType string
		want            string
		args            args
	}{
		{
			name:            "Writes JSON data to the response upon success",
			wantStatus:      http.StatusOK,
			wantContentType: "application/json",
			want:            `{"name":"Adam","age":10}`,
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
				value: sampleStruct{
					Name: "Adam",
					Age:  10,
				},
			},
		},
		{
			name:            "Writes an error message when there is a problem marshaling JSON data",
			wantStatus:      http.StatusInternalServerError,
			wantContentType: "application/json",
			want:            `{"message":"Error marshaling value for writing","suggestion":"See error log for more information"}`,
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusInternalServerError,
				value:  func() string { return "not in my house!" },
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nerdweb.WriteJSON(logger, tt.args.w, tt.args.status, tt.args.value)

			gotStatus := tt.args.w.Result().StatusCode
			gotContentType := tt.args.w.Header().Get("Content-Type")
			got := strings.TrimSpace(tt.args.w.Body.String())

			if gotStatus != tt.wantStatus {
				t.Errorf("wanted status %d, got %d", tt.wantStatus, gotStatus)
			}

			if gotContentType != tt.wantContentType {
				t.Errorf("wanted content type '%s', got '%s'", tt.wantContentType, gotContentType)
			}

			if got != strings.TrimSpace(tt.want) {
				t.Errorf("want: %s\ngot: %s", strings.TrimSpace(tt.want), got)
			}
		})
	}
}

func TestWriteString(t *testing.T) {
	w := httptest.NewRecorder()
	logger := logrus.New().WithField("who", "testing")

	nerdweb.WriteString(logger, w, http.StatusBadRequest, "this is a test")

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("wanted status %d, got %d", http.StatusBadRequest, w.Result().StatusCode)
	}

	if w.Header().Get("Content-Type") != "text/plain" {
		t.Errorf("wanted content type of text/plain, got %s", w.Header().Get("Content-Type"))
	}

	if w.Body.String() != "this is a test" {
		t.Errorf("wanted 'this is a test', got %s", w.Body.String())
	}
}
