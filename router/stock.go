package router

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
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
		s.stock_name, 
		ARRAY_TO_STRING(ARRAY_AGG(s.price), ',')
	FROM 
		stock_data s
	INNER JOIN (
		SELECT 
			name
		FROM
			stocks s
		WHERE
			is_valid='t'
		) T
	ON 
		s.stock_name = T.name
	GROUP BY
		s.stock_name
	`)

	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Data Not Found", nil)
			return
		}
		util.GlobalErr(w, 500, "GET Data ERROR", err)
		return
	}

	var price string
	for rows.Next() {
		var stock util.AllStockData
		err := rows.Scan(&stock.Name, &price)

		strPrice := strings.Split(price, ",")
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
		buyInfo util.DealStock
		// 주식 가격
		price    int
		userCoin int
	)

	err := json.NewDecoder(r.Body).Decode(&buyInfo)
	if err != nil {
		util.GlobalErr(w, 400, "cannot read data", err)
		return
	}

	err = db.QueryRow(`SELECT stock_id FROM stocks WHERE 'is_valid'='t' AND name=$1`, &buyInfo.StockName).Err()
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
	`, buyInfo.StockName).Scan(&price)

	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Not Found Data", nil)
			return
		}
		util.GlobalErr(w, 500, "Get Data error", err)
		return
	}

	// 유저의 자본금 가져오기
	err = db.QueryRow(`SELECT coin FROM account WHERE user_id=$1;`, buyInfo.UserId).Scan(&userCoin)
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

	err = db.QueryRow(`
		SELECT user_id FROM user_stock WHERE user_id=$1 AND name=$2;
	`, buyInfo.UserId, buyInfo.StockName).Scan(&buyInfo.UserId)

	// 주식 수익률 = 현재 주가 / 이때까지 구매한 총액 * 100 - 100
	// 수익금 = 구매가 * 수익률 * 구매량
	if err == sql.ErrNoRows {
		_, err = db.Exec(`
			INSERT INTO
				user_stock
			VALUES($1, $2, $3, $4)
		`, buyInfo.UserId, buyInfo.StockName, price*buyInfo.Count, buyInfo.Count)

		if err != nil {
			util.GlobalErr(w, 500, "Update Data error", err)
			return
		}
	} else if err != nil {
		util.GlobalErr(w, 500, "cannot update or create", err)
		return
	} else {
		_, err = db.Exec(`
			UPDATE user_stock SET cost=cost + $3, count=count + $4 WHERE user_id=$1 AND name=$2
		`, buyInfo.UserId, buyInfo.StockName, price*buyInfo.Count, buyInfo.Count)

		if err != nil {
			util.GlobalErr(w, 500, "Update Data error", err)
			return
		}
	}

	userCoin -= (price * buyInfo.Count)
	_, err = db.Exec(`UPDATE account SET coin=$2 WHERE user_id=$1`, buyInfo.UserId, userCoin)
	if err != nil {
		util.GlobalErr(w, 500, "Update Data error", err)
		return
	}

	util.ResOk(w, 200, "update success")
}

func SellStock(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		sellInfo  util.DealStock
		stockInfo util.UserStock
		// 현재 주가
		price int
	)

	err := json.NewDecoder(r.Body).Decode(&sellInfo)
	if err != nil {
		util.GlobalErr(w, 400, "cannot read data", nil)
		return
	}
	err = db.QueryRow(`SELECT stock_id FROM stocks WHERE 'is_valid'='t' AND name=$1`, &sellInfo.StockName).Err()
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
	`, sellInfo.StockName).Scan(&price)

	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Not Found Data", nil)
			return
		}
		util.GlobalErr(w, 500, "Get Data error", err)
		return
	}

	// 유저의 보유 주식 정보 가져오기
	err = db.QueryRow(`
	SELECT 
		cost, 
		count 
	FROM 
		user_stock 
	WHERE 
		user_id=$1 
	AND 
		name=$2;
	`, sellInfo.UserId, sellInfo.StockName).
		Scan(&stockInfo.Cost, &stockInfo.Count)

	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Not Found Data", nil)
			return
		}
		util.GlobalErr(w, 500, "Get Data error", err)
		return
	}

	// 보유 주식량보다 많다면
	if stockInfo.Count < sellInfo.Count {
		util.GlobalErr(w, 400, "Not enough stock count", nil)
		return
	}
	_, err = db.Exec(`
			UPDATE user_stock SET cost=cost - $3, count=count - $4 WHERE user_id=$1 AND name=$2
		`, sellInfo.UserId, sellInfo.StockName, price*sellInfo.Count, sellInfo.Count)

	if err != nil {
		util.GlobalErr(w, 500, "Update Data error", err)
		return
	}

	_, err = db.Exec(`UPDATE account SET coin=coin+$2 WHERE user_id=$1`, sellInfo.UserId, price*sellInfo.Count)
	if err != nil {
		util.GlobalErr(w, 500, "Update Data error", err)
		return
	}

	y := math.Round(
		float64(
			(((price*stockInfo.Count)/stockInfo.Cost)*100-100)*100)) / 100
	util.ResOk(w, 200, util.SellStockRes{
		StockName: sellInfo.StockName,
		Yield:     y,
		Proceeds: int(
			math.Round(
				float64(stockInfo.Cost) * y * float64(stockInfo.Count))),
	})
}
