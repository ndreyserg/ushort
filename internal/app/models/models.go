// Package models с основными моделями
package models

// Request - структура запроса на сохранение ссылки
type Request struct {
	URL string `json:"url"`
}

// Response - структура ответ на запрос сохранения ссылки
type Response struct {
	Result string `json:"result"`
}

// BatchRequestItem - структура элемента batch запроса на сохранение ссылки
type BatchRequestItem struct {
	ID       string `json:"correlation_id"`
	Original string `json:"original_url"`
}

// BatchRequest - структура batch запроса на сохранение ссылки
type BatchRequest []BatchRequestItem

// BatchResultItem - структура элемента ответа на batch запрос на сохранение ссылки
type BatchResultItem struct {
	ID    string `json:"correlation_id"`
	Short string `json:"short_url"`
}

// BatchResult - структура ответа на batch запрос на сохранение ссылки
type BatchResult []BatchResultItem
