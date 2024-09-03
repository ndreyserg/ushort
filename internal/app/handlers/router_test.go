package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type state map[string]string

type fakeStorage struct {
	state state
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

type want struct {
	statusCode int
	body       string
}

func TestGet(t *testing.T) {
	const host = "http://localhost:8080/"
	tests := []struct {
		name    string
		request string
		body    string
		method  string
		state   state
		want    want
	}{
		{
			name:    "unknown method",
			request: host,
			body:    "",
			state:   state{"asdasd": "https://ya.ru"},
			method:  http.MethodPut,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "",
			},
		},
		{
			name:    "empty key",
			request: host,
			body:    "",
			state:   state{"asdasd": "https://ya.ru"},
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "empty key",
			},
		},
		{
			name:    "unknown key",
			request: host + "dddd",
			body:    "",
			state:   state{"asdasd": "https://ya.ru"},
			method:  http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "key not found",
			},
		},
		{
			name:    "empty storage",
			request: host + "dddd",
			body:    "",
			method:  http.MethodGet,
			state:   state{},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "key not found",
			},
		},
		{
			name:    "existed key",
			request: host + "asdasd",
			body:    "",
			method:  http.MethodGet,
			state:   state{"asdasd": "https://ya.ru"},
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				body:       "<a href=\"https://ya.ru\">Temporary Redirect</a>.",
			},
		},
		{
			name:    "post empty link",
			request: host,
			body:    "",
			method:  http.MethodPost,
			state:   state{},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "empty request body",
			},
		},
		{
			name:    "post link",
			request: host,
			body:    "http://ya.ru",
			method:  http.MethodPost,
			state:   state{},
			want: want{
				statusCode: http.StatusCreated,
				body:       host + "linkID",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fk := fakeStorage{
				state: test.state,
			}
			request := httptest.NewRequest(test.method, test.request, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			MakeRouter(fk)(w, request)
			res := w.Result()
			if assert.Equal(
				t,
				test.want.statusCode,
				res.StatusCode,
				"expected status code %d got %d",
				test.want.statusCode, res.StatusCode,
			) {
				assert.Equal(
					t,
					test.want.body,
					strings.Trim(w.Body.String(), "\n"),
					"expected body \"%s\" got  \"%s\"",
					test.body,
					test.want.body,
				)
			}

		})
	}
}
