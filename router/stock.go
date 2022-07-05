package router

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"stockServer/util"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func AllStockList(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	var stocks []util.AllStockData
	rows, err := db.Query(`
	SELECT 
		stock_name, 
		ARRAY_TO_STRING(ARRAY_AGG(price), ',')
	FROM 
		stock_data
	WHERE 
		stock_name=(
			SELECT name FROM stocks WHERE 'isVaild'='t'
		)
	GROUP BY
		stock_name
	`)

	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Data Not Found", nil)
			return
		}
		util.GlobalErr(w, 500, "GET Data ERROR", err)
		return
	}

	for rows.Next() {
		var stock util.AllStockData
		err := rows.Scan(&stock.Name, &stock.Price)

		strPrice := strings.Split(stock.Price, ",")
		for _, i := range strPrice {
			p, err := strconv.Atoi(i)
			if err != nil {
				log.Println(err)
				return
			}

			stock.PriceList = append(stock.PriceList, p)
		}

		if err != nil {
			util.GlobalErr(w, 500, "Parse Data ERROR", err)
			return
		}

		stocks = append(stocks, stock)
	}

	util.ResOk(w, 200, stocks)
}

func BuyStock(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		buyInfo  util.DealStock
		price    int
		userCoin int

		// 주식 원가
		cost int
	)

	err := json.NewDecoder(r.Body).Decode(&buyInfo)
	if err != nil {
		util.GlobalErr(w, 400, "cannot read data", nil)
		return
	}

	err = db.QueryRow(`SELECT stock_id FROM stocks WHERE 'isValid'='t' AND stock_name=$1`, &buyInfo.StockName).Err()
	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 400, "Not found valid stocks", nil)
			return
		}

		util.GlobalErr(w, 500, "Can Not found valid stocks", err)
		return
	}

	// 가장 최신 주가
	err = db.QueryRow(`
	SELECT 
		price
	FROM
		stock_data 
	WHERE
		stock_name=$1 
	AND 
		data_id=(
			SELECT MAX(data_id) FROM stock_data WHERE stock_name=$1
		)
	`, buyInfo.StockName).Scan(price)

	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Not Found Data", nil)
			return
		}
		util.GlobalErr(w, 500, "Get Data error", err)
		return
	}

	// 유저의 자본금 가져오기
	err = db.QueryRow(`SELECT coin FROM account WHERE discord_id=$1;`, buyInfo.UserId).Scan(&userCoin)
	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Not Found Data", nil)
			return
		}
		util.GlobalErr(w, 500, "Get Data error", err)
		return
	}

	// 주식을 살 돈이 없다면
	if userCoin < (price * buyInfo.Count) {
		util.GlobalErr(w, 400, "Not enough coin", nil)
		return
	}

	_, err = db.Exec(`
	UPDATE 
		user_stock 
	SET 
		count=$1, cost=$2 
	WHERE 
		user_id=$3 NOT EXIST (
			INSERT INTO
				user_stock
			VALUES
				($3, $4, $2, $1)
		);
	`, buyInfo.Count, cost, buyInfo.UserId, buyInfo.StockName)

	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Not Found Data", nil)
			return
		}
		util.GlobalErr(w, 500, "Update Data error", err)
		return
	}

	util.ResOk(w, 200, "update success")
}

func SellStock(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		sellInfo   util.DealStock
		stockInfo  util.Stock
		price      int
		stockCount int
	)

	err := json.NewDecoder(r.Body).Decode(&sellInfo)
	if err != nil {
		util.GlobalErr(w, 400, "cannot read data", nil)
		return
	}

	// 가장 최신 주가
	err = db.QueryRow(`
	SELECT 
		s.stock_id, s.data[array_upper(data, 1)], a.count 
	FROM 
		stocks s 
	INNER JOIN 
		user_stock a
	ON
		a.user_id=$2 
	AND 
		a.name=$1
	WHERE
		s.stock_name=$1 
	AND
		expire_t IS NULL;
	`, sellInfo.StockName, sellInfo.UserId).
		Scan(stockInfo.StockId, price, stockCount)
	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Not Found Data", nil)
			return
		}
		util.GlobalErr(w, 500, "Get Data error", err)
		return
	}

	if stockCount < sellInfo.Count {
		util.GlobalErr(w, 400, "over stock's count", nil)
		return
	}

	_, err = db.Exec(`
	UPDATE 
		account 
	SET 
		coin=(
			SELECT 
				coin 
			FROM 
				account 
			WHERE 
				discord_id=$1
		) + $2
	`, sellInfo.UserId, sellInfo.Count*price)

	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Not Found User", nil)
			return
		}
		util.GlobalErr(w, 400, "update error", err)
		return
	}
}
