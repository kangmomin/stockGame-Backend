package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ResOk(w http.ResponseWriter, code int, data any) {
	w.WriteHeader(code)
	resData, _ := json.Marshal(res{
		Data: data,
		Err:  false,
	})

	fmt.Fprint(w, string(resData))
}

func GlobalErr(w http.ResponseWriter, code int, data string, err error) {
	w.WriteHeader(code)
	if err != nil {
		log.Println(err)
	}

	resData, _ := json.Marshal(res{
		Data: data,
		Err:  true,
	})

	fmt.Fprint(w, string(resData))
}

type SellStockRes struct {
	StockName string `json:"stock_name"`

	// 수익률
	Yield float64 `json:"yield"`

	// 수익금
	Proceeds int `json:"proceeds"`

	// // 자본금
	// Coin int `json:"coin"`

	// // 남은 주식 개수
	// Count int `json:"count"`
}
