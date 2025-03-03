package handlers_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ndreyserg/ushort/internal/app/auth"
	"github.com/ndreyserg/ushort/internal/app/handlers"
	"github.com/ndreyserg/ushort/internal/app/queue"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

type config struct {
	ServerAddr string
	BaseURL    string
}

func Example() {
	cfg := config{
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost:8080",
	}

	st := storage.NewMemoryStorage()
	s := auth.NewJWTSession("1234")
	q := queue.NewQueue(st)
	go func() {
		http.ListenAndServe(cfg.ServerAddr, handlers.MakeRouter(st, cfg.BaseURL, s, q))
	}()

	time.Sleep(1 * time.Second)

	//Отправим запрос на сохранение ссылки
	resp, err := http.Post("http://localhost:8080/", "", bytes.NewReader([]byte("http://ya.ru")))

	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		panic(err)
	}

	respBody, _ := io.ReadAll(resp.Body)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// Перейдем по пройденной ссылке и проверим статус и заголовок
	resp, err = client.Get(string(respBody))
	fmt.Println(resp.Status)
	fmt.Println(resp.Header.Get("Location"))

	// Output:
	// 307 Temporary Redirect
	// http://ya.ru
}
