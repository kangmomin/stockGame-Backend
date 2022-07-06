package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"stockServer/util"

	"github.com/julienschmidt/httprouter"
)

func GetUserInfo(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	var (
		userInfo   util.UserInfo
		userStocks []util.UserStock
	)

	userId := p.ByName("userId")

	err := db.QueryRow(`
	SELECT 
		name, coin, bank, tax, gamble_ticket 
	FROM 
		account 
	WHERE 
		user_id=$1;
	`, userId).Scan(&userInfo.Name, &userInfo.Coin, &userInfo.Bank, &userInfo.Tax, &userInfo.GambleTicket)
	if err != nil {
		if err == sql.ErrNoRows {
			util.ResOk(w, 404, "User Not Found")
			return
		}
		util.GlobalErr(w, 500, "cannot get info", err)
		return
	}

	stockData, err := db.Query(`SELECT name, cost, count FROM user_stock WHERE user_id=$1;`, userId)
	if err != nil {
		util.GlobalErr(w, 500, "cannot get info", err)
		return
	}

	for stockData.Next() {
		var userStock util.UserStock
		err := stockData.Scan(&userStock.Name, &userStock.Cost, userStock.Count)
		if err != nil {
			continue
		}

		userStocks = append(userStocks, userStock)
	}

	totalVal := struct {
		UserInfo util.UserInfo    `json:"user_info"`
		Stocks   []util.UserStock `json:"user_stock"`
	}{
		UserInfo: userInfo,
		Stocks:   userStocks,
	}

	util.ResOk(w, 200, totalVal)
}

func SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		signUpData util.SignUp
	)

	err := json.NewDecoder(r.Body).Decode(&signUpData)
	if err != nil {
		util.GlobalErr(w, 400, "Wrong Request", err)
		return
	}

	if !util.IsValidUserId(signUpData.Token, signUpData.UserId) {
		util.GlobalErr(w, 400, "Not Found User", nil)
		return
	}

	err = db.QueryRow(`SELECT user_id FROM account WHERE userId=$1`, &signUpData.UserId).Err()
	if err != nil {
		util.GlobalErr(w, 400, "already sign-up account", nil)
		return
	}

	_, err = db.Exec(`INSERT INTO account(user_id, name) VALUES ($1, $2)`, &signUpData.UserId, &signUpData.Name)
	if err != nil {
		util.GlobalErr(w, 400, "cannot create", err)
		return
	}

	util.ResOk(w, http.StatusCreated, nil)
}
