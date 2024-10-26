package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type BatchRequestItem struct {
	ID       string `json:"correlation_id"`
	Original string `json:"original_url"`
}

type BatchRequest []BatchRequestItem

type BatchResultItem struct {
	ID    string `json:"correlation_id"`
	Short string `json:"short_url"`
}

type BatchResult []BatchResultItem
