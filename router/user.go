package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"stockServer/util"

	"github.com/julienschmidt/httprouter"
)

func GetUserInfo(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	var userInfo util.UserInfo
	userId := p.ByName("userId")
	err := db.QueryRow(`SELECT * FROM account WHERE discord_id=$1;`, userId).
		Scan(&userInfo.DiscordId, &userInfo.Name, &userInfo.Coin, &userInfo.Bank, &userInfo.Tax, &userInfo.GambleTicket)
	if err != nil {
		if err == sql.ErrNoRows {
			util.ResOk(w, 404, "User Not Found")
			return
		}
		util.GlobalErr(w, 500, "cannot get info", err)
		return
	}
	data, err := json.Marshal(userInfo)

	if err != nil {
		util.GlobalErr(w, 500, "idk", err)
	}

	util.ResOk(w, 200, string(data))
}

func GetMyInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		userInfo   util.UserInfo
		userStocks []util.UserStock

		// body data
		user struct {
			Id string `json:"user_id"`
		}
	)

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		util.GlobalErr(w, 400, "raw type error", nil)
		return
	}

	err = db.QueryRow(`SELECT * FROM account WHERE user_id=$1;`, user.Id).Scan(&userInfo.DiscordId, &userInfo.Name, &userInfo.Coin, &userInfo.Bank, &userInfo.Bank)
	if err != nil {
		if err == sql.ErrNoRows {
			util.ResOk(w, 404, "User Not Found")
			return
		}
		util.GlobalErr(w, 500, "cannot get info", err)
		return
	}

	stockData, err := db.Query(`SELECT name, cost, count FROM user_stock WHERE user_id=$1;`, user.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			util.ResOk(w, 404, "User Not Found")
			return
		}
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

	util.ResOk(w, 200, userStocks)
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

	_, err = db.Exec(`INSERT INTO account(user_id, name) VALUES ($1, $2)`, &signUpData.UserId, &signUpData.Name)
	if err != nil {
		util.GlobalErr(w, 400, "cannot create", err)
		return
	}

	util.ResOk(w, http.StatusCreated, nil)
}
