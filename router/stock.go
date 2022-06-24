package router

import (
	"database/sql"
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
