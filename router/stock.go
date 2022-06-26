package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"stockServer/util"

	"github.com/julienschmidt/httprouter"
	pg "github.com/lib/pq"
)

func AllStockList(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	var stocks []util.Stock
	rows, err := db.Query(`SELECT stock_id, name, data, date FROM stocks WHERE expire_t IS NULL`)
	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Data Not Found", nil)
			return
		}
		util.GlobalErr(w, 500, "GET Data ERROR", err)
		return
	}

	for rows.Next() {
		var stock util.Stock
		err := rows.Scan(&stock.StockId, &stock.Name, pg.Array(&stock.Data), pg.Array(&stock.Date))
		if err != nil {
			util.GlobalErr(w, 500, "GET Data ERROR", err)
			return
		}

		stocks = append(stocks, stock)
	}

	util.ResOk(w, 200, stocks)
}

func BuyStock(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		buyInfo   util.DealStock
		stockInfo util.Stock
		price     int
		userCoin  int

		// 주식 원가
		cost int
	)

	err := json.NewDecoder(r.Body).Decode(&buyInfo)
	if err != nil {
		util.GlobalErr(w, 400, "cannot read data", nil)
		return
	}

	// 가장 최신 주가
	err = db.QueryRow(`SELECT stock_id, data[array_upper(data, 1)] FROM stocks WHERE stock_name=$1 AND expire_t IS NULL;`, buyInfo.StockName).
		Scan(stockInfo.StockId, price)
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
		UPDATE user_stock SET count=$1, cost=$2 WHERE user_id=$3 NOT EXIST (
			INSERT INTO user_stock VALUES ($3, $4, $2, $1)
		);`, buyInfo.Count, cost, buyInfo.UserId, buyInfo.StockName)

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
			SELECT s.stock_id, s.data[array_upper(data, 1)], a.count FROM stocks s 
			INNER JOIN user_stock u ON a.user_id=$2 AND a.name=$1
			WHERE s.stock_name=$1 AND expire_t IS NULL;`, sellInfo.StockName, sellInfo.UserId).
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

	_, err = db.Exec(`UPDATE account SET coin=(
		SELECT coin FROM account WHERE discord_id=$1
	) + $2`, sellInfo.UserId, sellInfo.Count*price)

	if err != nil {
		if err == sql.ErrNoRows {
			util.GlobalErr(w, 404, "Not Found User", nil)
			return
		}
		util.GlobalErr(w, 400, "update error", err)
		return
	}
}
