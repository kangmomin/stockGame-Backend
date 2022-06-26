package util

type Stock struct {
	StockId int      `json:"stock_id"`
	Name    string   `json:"name"`
	Data    []string `json:"data"`
	Date    []string `json:"date"`
}

type BuyStock struct {
	StockName string `json:"stock_name"`
	UserId    string `json:"user_id"`
	Count     int    `json:"count"`
}
