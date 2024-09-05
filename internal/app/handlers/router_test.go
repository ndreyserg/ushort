package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeStorage struct {
	state map[string]string
}

func (fk fakeStorage) Get(key string) (string, error) {
	res, ok := fk.state[key]
	if ok {
		return res, nil
	}
	return "", errors.New("")
}

func (fk fakeStorage) Set(val string) string {
	return "linkID"
}

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
	type want struct {
		statusCode int
		body       string
	}
	const urlAddr = "localhost:8080"
	ts := httptest.NewServer(MakeRouter(fakeStorage{state: map[string]string{"asdasd": "https://ya.ru"}}, urlAddr))
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
			request: "/dddd",
			body:    "",
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "key not found",
			},
		},
		{
			name:    "existed key",
			request: "/asdasd",
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
				body:       fmt.Sprintf("http://%s/linkID", urlAddr),
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
					test.body,
					test.want.body,
				)
			}

		})
	}
}
