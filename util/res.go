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
