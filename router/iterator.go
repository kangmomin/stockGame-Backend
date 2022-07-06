package router

import (
	"log"
	"math/rand"
	"stockServer/util"
	"strconv"
	"time"
)

func UpdateStock() {
	var (
		// update query
		query string

		// only name[string] & price[int]
		newStocks []any

		// iterated count for update query
		i int
	)

	i = 1

	rows, err := db.Query(`
		SELECT 
			stock_name, price 
		FROM 
			stock_data 
		WHERE 
			stock_name=(
				SELECT name FROM stocks WHERE is_valid='t'
			);
		`)
	if err != nil {
		log.Println(err)
		return
	}

	if !rows.Next() {
		return
	}

	for rows.Next() {
		var (
			newStock util.StockData
		)
		query += `($,` + strconv.Itoa(i) + ` $` + strconv.Itoa(i+1) + `)`
		varPrice := getNewStock()
		err = rows.Scan(&newStock.Name, &newStock.Price)
		if err != nil {
			log.Println(err)
			return
		}

		newStock.Price += varPrice
		newStocks = append(newStocks, newStock.Name, newStock.Price)
		i += 2
	}

	_, err = db.Exec(`
		INSERT INTO
			stock_data(price, stock_name) 
		VALUES`+query, newStocks...)
	if err != nil {
		log.Println(err)
		return
	}
}

func getNewStock() int {
	s := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(s)

	varPrice := r.Intn(101)
	if varPrice > 80 {
		varPrice = r.Intn(101)
	}

	if r.Intn(2) == 0 {
		varPrice *= -1
	}

	return varPrice
}
