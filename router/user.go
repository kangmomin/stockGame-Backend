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
		Scan(&userInfo.DiscordId, &userInfo.Name, &userInfo.Coin, &userInfo.Bank, &userInfo.Tax, &userInfo.GambleTicket, &userInfo.Stock)
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
