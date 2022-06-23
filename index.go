package main

import (
	"net/http"
	"stockServer/router"

	"github.com/julienschmidt/httprouter"
)

func main() {
	app := httprouter.New()

	app.GET("/user/info/:userId", router.GetUserInfo)

	http.ListenAndServe(":8080", app)
}
