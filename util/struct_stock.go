package util

import "time"

type Stock struct {
	StockId int       `json:"stock_id"`
	Name    string    `json:"name"`
	ExpireT time.Time `json:"expire_t"`
	IsValid bool      `json:"isValid"`
}

type StockData struct {
	DataId    int       `json:"data_id"`
	Price     int       `json:"price"`
	UpdatedAt time.Time `json:"update_at"`
	Name      string    `json:"name"`
}

type DealStock struct {
	StockName string `json:"stock_name"`
	UserId    string `json:"user_id"`
	Count     int    `json:"count"`
}

type UserStock struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Cost   int    `json:"cost"`
	Count  int    `json:"count"`
}
