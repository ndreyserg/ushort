package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ndreyserg/ushort/internal/app/mocks"
	"github.com/ndreyserg/ushort/internal/app/models"
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

func TestRouterPostJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	d := initDeps(t, ctrl)
	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("", errors.New(""))

	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("andrey", nil)
	d.storage.EXPECT().Set(gomock.Any(), gomock.Eq("new_url"), gomock.Eq("andrey")).Return("", errors.New(""))

	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("andrey", nil)
	d.storage.EXPECT().Set(gomock.Any(), gomock.Eq("new_url"), gomock.Eq("andrey")).Return("", storage.ErrConflict)

	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("andrey", nil)
	d.storage.EXPECT().Set(gomock.Any(), gomock.Eq("new_url"), gomock.Eq("andrey")).Return("new_short_link", nil)

	ts := httptest.NewServer(MakeRouter(d.storage, baseURL, d.sessison, d.queue))
	tests := []tCase{
		{
			name:           "uncorrect json",
			request:        "/api/shorten",
			body:           `{'ee':`,
			method:         http.MethodPost,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "empty url",
			request:        "/api/shorten",
			body:           `{"url":""}`,
			method:         http.MethodPost,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "session error",
			request:        "/api/shorten",
			body:           `{"url":"new_url"}`,
			method:         http.MethodPost,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "storage error",
			request:        "/api/shorten",
			body:           `{"url":"new_url"}`,
			method:         http.MethodPost,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "storage conflict",
			request:        "/api/shorten",
			body:           `{"url":"new_url"}`,
			method:         http.MethodPost,
			wantStatusCode: http.StatusConflict,
		},

		{
			name:            "storage conflict",
			request:         "/api/shorten",
			body:            `{"url":"new_url"}`,
			method:          http.MethodPost,
			wantStatusCode:  http.StatusCreated,
			hasResponseBody: true,
			responseBody:    fmt.Sprintf(`{"result":"%s/new_short_link"}`, baseURL),
		},
	}

	runTest(t, tests, ts)
}

func TestRouterPostBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	d := initDeps(t, ctrl)

	req := models.BatchRequest{
		models.BatchRequestItem{ID: "1", Original: "original1"},
		models.BatchRequestItem{ID: "2", Original: "original2"},
	}

	res := models.BatchResult{
		models.BatchResultItem{ID: "1", Short: "short1"},
		models.BatchResultItem{ID: "2", Short: "short2"},
	}

	reqB, _ := json.Marshal(req)

	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("", errors.New(""))

	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("andrey", nil)
	d.storage.EXPECT().SetBatch(gomock.Any(), gomock.Eq(req), gomock.Eq("andrey")).Return(nil, errors.New(""))

	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("andrey", nil)
	d.storage.EXPECT().SetBatch(gomock.Any(), gomock.Eq(req), gomock.Eq("andrey")).Return(nil, storage.ErrConflict)

	d.sessison.EXPECT().Open(gomock.All(), gomock.Any()).Return("andrey", nil)
	d.storage.EXPECT().SetBatch(gomock.Any(), gomock.Eq(req), gomock.Eq("andrey")).Return(res, nil)

	ts := httptest.NewServer(MakeRouter(d.storage, baseURL, d.sessison, d.queue))
	tests := []tCase{
		{
			name:           "uncorrect json",
			request:        "/api/shorten/batch",
			body:           `{'ee':`,
			method:         http.MethodPost,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "empty request",
			request:        "/api/shorten/batch",
			body:           `[]`,
			method:         http.MethodPost,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "session error",
			request:        "/api/shorten/batch",
			body:           string(reqB),
			method:         http.MethodPost,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "storage error",
			request:        "/api/shorten/batch",
			body:           string(reqB),
			method:         http.MethodPost,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "storage conflict",
			request:        "/api/shorten/batch",
			body:           string(reqB),
			method:         http.MethodPost,
			wantStatusCode: http.StatusCreated,
		},

		{
			name:            "success",
			request:         "/api/shorten/batch",
			body:            string(reqB),
			method:          http.MethodPost,
			wantStatusCode:  http.StatusCreated,
			hasResponseBody: true,
			responseBody: fmt.Sprintf(
				`[{"correlation_id":"1","short_url":"%s/short1"},{"correlation_id":"2","short_url":"%s/short2"}]`,
				baseURL, baseURL,
			),
		},
	}
	runTest(t, tests, ts)
}

func TestRouterPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	d := initDeps(t, ctrl)

	d.storage.EXPECT().Check(gomock.Any()).Return(errors.New(""))
	d.storage.EXPECT().Check(gomock.Any()).Return(nil)

	ts := httptest.NewServer(MakeRouter(d.storage, baseURL, d.sessison, d.queue))
	tests := []tCase{
		{
			name:           "error",
			request:        "/ping",
			body:           "",
			method:         http.MethodGet,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "success",
			request:        "/ping",
			body:           "",
			method:         http.MethodGet,
			wantStatusCode: http.StatusOK,
		},
	}
	runTest(t, tests, ts)
}

func TestRouterGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	d := initDeps(t, ctrl)
	d.storage.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Return("", errors.New(""))
	d.storage.EXPECT().Get(gomock.Any(), gomock.Eq("2")).Return("", storage.ErrIsGone)

	d.storage.EXPECT().Get(gomock.Any(), gomock.Eq("3")).Return("original url", nil)
	ts := httptest.NewServer(MakeRouter(d.storage, baseURL, d.sessison, d.queue))

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	tests := []tCase{
		{
			name:           "empty",
			request:        "/",
			body:           "",
			method:         http.MethodGet,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "storage error",
			request:        "/1",
			body:           "",
			method:         http.MethodGet,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "url is gone",
			request:        "/2",
			body:           "",
			method:         http.MethodGet,
			wantStatusCode: http.StatusGone,
		},
		{
			name:           "success",
			request:        "/3",
			body:           "",
			method:         http.MethodGet,
			wantStatusCode: http.StatusTemporaryRedirect,
		},
	}
	runTest(t, tests, ts)
}

func TestRouterGetUserUrls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	d := initDeps(t, ctrl)

	d.sessison.EXPECT().Open(gomock.Any(), gomock.Any()).Return("andrey", nil)
	d.storage.EXPECT().GetUserUrls(gomock.Any(), gomock.Eq("andrey")).Return(nil, errors.New(""))

	d.sessison.EXPECT().Open(gomock.Any(), gomock.Any()).Return("andrey", nil)
	d.storage.EXPECT().GetUserUrls(gomock.Any(), gomock.Eq("andrey")).Return([]storage.StorageItem{}, nil)

	d.sessison.EXPECT().Open(gomock.Any(), gomock.Any()).Return("andrey", nil)
	d.storage.EXPECT().GetUserUrls(gomock.Any(), gomock.Eq("andrey")).Return([]storage.StorageItem{
		{Original: "orig", Short: "short"},
	}, nil)

	ts := httptest.NewServer(MakeRouter(d.storage, baseURL, d.sessison, d.queue))
	tests := []tCase{

		{
			name:           "storage error",
			request:        "/api/user/urls",
			body:           "",
			method:         http.MethodGet,
			wantStatusCode: http.StatusInternalServerError,
		},

		{
			name:           "empty",
			request:        "/api/user/urls",
			body:           "",
			method:         http.MethodGet,
			wantStatusCode: http.StatusNoContent,
		},

		{
			name:            "success",
			request:         "/api/user/urls",
			body:            "",
			method:          http.MethodGet,
			wantStatusCode:  http.StatusOK,
			hasResponseBody: true,
			responseBody:    fmt.Sprintf(`[{"short_url":"%s/short","original_url":"orig"}]`, baseURL),
		},
	}
	runTest(t, tests, ts)
}

func TestRouterDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	d := initDeps(t, ctrl)

	d.sessison.EXPECT().Open(gomock.Any(), gomock.Any()).Return("", errors.New(""))

	d.sessison.EXPECT().Open(gomock.Any(), gomock.Any()).Return("andrey", nil)
	d.queue.EXPECT().AddTask(gomock.Eq([]string{"short"}), gomock.Eq("andrey"))

	ts := httptest.NewServer(MakeRouter(d.storage, baseURL, d.sessison, d.queue))
	tests := []tCase{
		{
			name:           "bad request",
			request:        "/api/user/urls",
			body:           `["short]`,
			method:         http.MethodDelete,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "session error",
			request:        "/api/user/urls",
			body:           `["short"]`,
			method:         http.MethodDelete,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "success",
			request:        "/api/user/urls",
			body:           `["short"]`,
			method:         http.MethodDelete,
			wantStatusCode: http.StatusAccepted,
		},
	}
	runTest(t, tests, ts)
}
