package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ndreyserg/ushort/internal/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, reqBody string, path string) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(reqBody))
	require.NoError(t, err)
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageMock := mocks.NewMockStorage(ctrl)
	storageMock.EXPECT().Check(gomock.Any()).Return(nil)
	storageMock.EXPECT().Get(gomock.Any(), gomock.Eq("unknown_key")).Return("", errors.New(""))
	storageMock.EXPECT().Get(gomock.Any(), gomock.Eq("existed_key")).Return("https://ya.ru", nil)
	storageMock.EXPECT().Set(gomock.Any(), gomock.Eq("http://practicum.yndex.ru")).Return("new_short_link", nil).Times(2)

	type want struct {
		statusCode int
		body       string
	}
	const baseURL = "http://localhost:8080"
	ts := httptest.NewServer(MakeRouter(storageMock, baseURL))
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	tests := []struct {
		name    string
		request string
		body    string
		method  string
		want    want
	}{
		{
			name:    "unknown method",
			request: "/",
			body:    "",
			method:  http.MethodPut,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "method not allowed",
			},
		},
		{
			name:    "empty key",
			request: "/",
			body:    "",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "method not allowed",
			},
		},
		{
			name:    "unknown key",
			request: "/unknown_key",
			body:    "",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "key not found",
			},
		},
		{
			name:    "existed key",
			request: "/existed_key",
			body:    "",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				body:       "<a href=\"https://ya.ru\">Temporary Redirect</a>.",
			},
		},
		{
			name:    "post empty link",
			request: "/",
			body:    "",
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "empty request body",
			},
		},
		{
			name:    "post link",
			request: "",
			body:    "http://practicum.yndex.ru",
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusCreated,
				body:       fmt.Sprintf("%s/new_short_link", baseURL),
			},
		},
		{
			name:    "post json link",
			request: "/api/shorten",
			body:    `{"url" :"http://practicum.yndex.ru"}`,
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusCreated,
				body:       fmt.Sprintf(`{"result":"%s/new_short_link"}`, baseURL),
			},
		},
		{
			name:    "ping DB",
			request: "/ping",
			body:    "",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusOK,
				body:       "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			resp, body := testRequest(t, ts, test.method, test.body, test.request)
			defer resp.Body.Close()
			if assert.Equal(
				t,
				test.want.statusCode,
				resp.StatusCode,
				"expected status code %d got %d",
				test.want.statusCode, resp.StatusCode,
			) {
				assert.Equal(
					t,
					test.want.body,
					strings.Trim(body, "\n"),
					"expected body \"%s\" got  \"%s\"",
					test.want.body,
					body,
				)
			}

		})
	}
}
