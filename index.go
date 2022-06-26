package main

import (
	"net/http"
	"stockServer/router"

	"github.com/julienschmidt/httprouter"
)

func main() {
	app := httprouter.New()

	// user
	app.GET("/user/info/:userId", router.GetUserInfo)

	// stock
	app.GET("/stock/all", router.AllStockList)
	app.POST("/stock/buy", router.BuyStock)
	app.POST("/stock/sell", router.SellStock)

	http.ListenAndServe(":8080", app)
}
