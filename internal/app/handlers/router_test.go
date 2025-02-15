package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ndreyserg/ushort/internal/app/mocks"
	"github.com/ndreyserg/ushort/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

type deps struct {
	storage  *mocks.MockStorage
	sessison *mocks.MockSession
	queue    *mocks.MockQueue
}

type want struct {
	statusCode int
	body       string
}

type tCase struct {
	name            string
	request         string
	body            string
	method          string
	wantStatusCode  int
	hasResponseBody bool
	responseBody    string
}

func initDeps(t *testing.T, ctrl *gomock.Controller) *deps {
	t.Helper()
	return &deps{
		storage:  mocks.NewMockStorage(ctrl),
		sessison: mocks.NewMockSession(ctrl),
		queue:    mocks.NewMockQueue(ctrl),
	}
}

func runTest(t *testing.T, tests []tCase, ts *httptest.Server) {
	t.Helper()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(
				test.method,
				ts.URL+test.request,
				strings.NewReader(test.body),
			)
			require.NoError(t, err)
			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer func() {
				_ = resp.Body.Close()
			}()
			assert.Equal(
				t,
				test.wantStatusCode,
				resp.StatusCode,
				"expected status code %d got %d",
				test.wantStatusCode,
				resp.StatusCode,
			)
			if test.hasResponseBody {
				respBody, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(
					t,
					test.responseBody,
					strings.Trim(string(respBody), "\n"),
					"expected body \"%s\" got  \"%s\"",
					test.responseBody,
					respBody,
				)
			}
		})
	}
}

func TestRouterPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	d := initDeps(t, ctrl)
	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("", errors.New(""))

	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("ndrey", nil)
	d.storage.EXPECT().Set(gomock.Any(), "new_link", "ndrey").Return("", storage.ErrConflict)

	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("ndrey", nil)
	d.storage.EXPECT().Set(gomock.Any(), "new_link", "ndrey").Return("", errors.New(""))

	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("ndrey", nil)
	d.storage.EXPECT().Set(gomock.Any(), "new_link", "ndrey").Return("new_short_link", nil)

	ts := httptest.NewServer(MakeRouter(d.storage, baseURL, d.sessison, d.queue))
	tests := []tCase{
		{
			name:            "empty body",
			request:         "/",
			body:            "",
			method:          http.MethodPost,
			wantStatusCode:  http.StatusBadRequest,
			hasResponseBody: true,
			responseBody:    "empty request body",
		},
		{
			name:           "session err",
			request:        "/",
			body:           "new_link",
			method:         http.MethodPost,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "conflict",
			request:        "/",
			body:           "new_link",
			method:         http.MethodPost,
			wantStatusCode: http.StatusConflict,
		},
		{
			name:           "storage error",
			request:        "/",
			body:           "new_link",
			method:         http.MethodPost,
			wantStatusCode: http.StatusInternalServerError,
		},

		{
			name:            "storage error",
			request:         "/",
			body:            "new_link",
			method:          http.MethodPost,
			wantStatusCode:  http.StatusCreated,
			hasResponseBody: true,
			responseBody:    baseURL + "/new_short_link",
		},
	}
	runTest(t, tests, ts)
}

// func testRequest(t *testing.T, ts *httptest.Server, method, reqBody string, path string) (*http.Response, string) {

// 	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(reqBody))
// 	require.NoError(t, err)
// 	resp, err := ts.Client().Do(req)
// 	require.NoError(t, err)
// 	defer resp.Body.Close()

// 	respBody, err := io.ReadAll(resp.Body)
// 	require.NoError(t, err)

// 	return resp, string(respBody)
// }

// func TestRouter(t *testing.T) {

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	incBatch := models.BatchRequest{
// 		models.BatchRequestItem{ID: "1", Original: "original1"},
// 		models.BatchRequestItem{ID: "2", Original: "original2"},
// 	}
// 	resBatch := models.BatchResult{
// 		models.BatchResultItem{ID: "1", Short: "short1"},
// 		models.BatchResultItem{ID: "2", Short: "short2"},
// 	}

// 	storageMock := mocks.NewMockStorage(ctrl)
// 	storageMock.EXPECT().Check(gomock.Any()).Return(nil)
// 	storageMock.EXPECT().Get(gomock.Any(), gomock.Eq("unknown_key")).Return("", errors.New(""))
// 	storageMock.EXPECT().Get(gomock.Any(), gomock.Eq("existed_key")).Return("https://ya.ru", nil)
// 	storageMock.EXPECT().Set(gomock.Any(), gomock.Eq("http://practicum.yndex.ru"), gomock.Any()).Return("new_short_link", nil).Times(2)
// 	storageMock.EXPECT().Set(gomock.Any(), gomock.Eq("conflict"), gomock.Any()).Return("old_short_link", storage.ErrConflict).Times(2)
// 	storageMock.EXPECT().SetBatch(gomock.Any(), gomock.Eq(incBatch), gomock.Any()).Return(resBatch, nil)

// 	sessionMock := mocks.NewMockSession(ctrl)
// 	sessionMock.EXPECT().Open(gomock.Any(), gomock.Any()).AnyTimes()

// 	queueMock := mocks.NewMockQueue(ctrl)

// 	type want struct {
// 		statusCode int
// 		body       string
// 	}
// 	const baseURL = "http://localhost:8080"
// 	ts := httptest.NewServer(MakeRouter(storageMock, baseURL, sessionMock, queueMock))
// 	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
// 		return http.ErrUseLastResponse
// 	}

// 	tests := []struct {
// 		name    string
// 		request string
// 		body    string
// 		method  string
// 		want    want
// 	}{
// 		{
// 			name:    "unknown method",
// 			request: "/",
// 			body:    "",
// 			method:  http.MethodPut,
// 			want: want{
// 				statusCode: http.StatusBadRequest,
// 				body:       "method not allowed",
// 			},
// 		},
// 		{
// 			name:    "empty key",
// 			request: "/",
// 			body:    "",
// 			method:  http.MethodGet,
// 			want: want{
// 				statusCode: http.StatusBadRequest,
// 				body:       "method not allowed",
// 			},
// 		},
// 		{
// 			name:    "unknown key",
// 			request: "/unknown_key",
// 			body:    "",
// 			method:  http.MethodGet,
// 			want: want{
// 				statusCode: http.StatusBadRequest,
// 				body:       "key not found",
// 			},
// 		},
// 		{
// 			name:    "existed key",
// 			request: "/existed_key",
// 			body:    "",
// 			method:  http.MethodGet,
// 			want: want{
// 				statusCode: http.StatusTemporaryRedirect,
// 				body:       "<a href=\"https://ya.ru\">Temporary Redirect</a>.",
// 			},
// 		},
// 		{
// 			name:    "post empty link",
// 			request: "/",
// 			body:    "",
// 			method:  http.MethodPost,
// 			want: want{
// 				statusCode: http.StatusBadRequest,
// 				body:       "empty request body",
// 			},
// 		},
// 		{
// 			name:    "post link",
// 			request: "",
// 			body:    "http://practicum.yndex.ru",
// 			method:  http.MethodPost,
// 			want: want{
// 				statusCode: http.StatusCreated,
// 				body:       fmt.Sprintf("%s/new_short_link", baseURL),
// 			},
// 		},
// 		{
// 			name:    "post conflict link",
// 			request: "",
// 			body:    "conflict",
// 			method:  http.MethodPost,
// 			want: want{
// 				statusCode: http.StatusConflict,
// 				body:       fmt.Sprintf("%s/old_short_link", baseURL),
// 			},
// 		},
// 		{
// 			name:    "post json link",
// 			request: "/api/shorten",
// 			body:    `{"url" :"http://practicum.yndex.ru"}`,
// 			method:  http.MethodPost,
// 			want: want{
// 				statusCode: http.StatusCreated,
// 				body:       fmt.Sprintf(`{"result":"%s/new_short_link"}`, baseURL),
// 			},
// 		},
// 		{
// 			name:    "post json conflict link",
// 			request: "/api/shorten",
// 			body:    `{"url" :"conflict"}`,
// 			method:  http.MethodPost,
// 			want: want{
// 				statusCode: http.StatusConflict,
// 				body:       fmt.Sprintf(`{"result":"%s/old_short_link"}`, baseURL),
// 			},
// 		},
// 		{
// 			name:    "ping DB",
// 			request: "/ping",
// 			body:    "",
// 			method:  http.MethodGet,
// 			want: want{
// 				statusCode: http.StatusOK,
// 				body:       "",
// 			},
// 		},
// 		{
// 			name:    "post empty batch",
// 			request: "/api/shorten/batch",
// 			body:    `[]`,
// 			method:  http.MethodPost,
// 			want: want{
// 				statusCode: http.StatusBadRequest,
// 				body:       "",
// 			},
// 		},
// 		{
// 			name:    "post batch",
// 			request: "/api/shorten/batch",
// 			body:    `[{"correlation_id": "1","original_url": "original1"}, {"correlation_id": "2","original_url": "original2"}]`,
// 			method:  http.MethodPost,
// 			want: want{
// 				statusCode: http.StatusCreated,
// 				body:       `[{"correlation_id":"1","short_url":"http://localhost:8080/short1"},{"correlation_id":"2","short_url":"http://localhost:8080/short2"}]`,
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {

// 			resp, body := testRequest(t, ts, test.method, test.body, test.request)
// 			defer resp.Body.Close()
// 			if assert.Equal(
// 				t,
// 				test.want.statusCode,
// 				resp.StatusCode,
// 				"expected status code %d got %d",
// 				test.want.statusCode, resp.StatusCode,
// 			) {
// 				assert.Equal(
// 					t,
// 					test.want.body,
// 					strings.Trim(body, "\n"),
// 					"expected body \"%s\" got  \"%s\"",
// 					test.want.body,
// 					body,
// 				)
// 			}

// 		})
// 	}
// }
