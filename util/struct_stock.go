package util

type Stock struct {
	StockId int      `json:"stock_id"`
	Name    string   `json:"name"`
	Data    []string `json:"data"`
	Date    []string `json:"date"`
}
