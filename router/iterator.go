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
		T1.stock_name,
		T1.price
	FROM 
		stock_data T1
		INNER JOIN (
			SELECT 
				stock_name, 
				max(data_id) max_id
			FROM
				stock_data
			INNER JOIN 
				stocks
			ON 
				stock_data.stock_name = stocks.name
			WHERE 
				is_valid='t'
			GROUP BY 
				stock_name 
			) T2
		ON 
			T1.data_id = T2.max_id 
		AND 
			T1.stock_name = T2.stock_name;
		`)
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		var (
			newStock util.StockData
		)
		if i != 1 {
			query += ","
		}
		query += `($` + strconv.Itoa(i) + `, $` + strconv.Itoa(i+1) + `)`
		varPrice := getNewStock()
		err = rows.Scan(&newStock.Name, &newStock.Price)
		if err != nil {
			log.Println(err)
			return
		}

		newStock.Price += varPrice
		if newStock.Price < 0 {
			newStock.Price = 0
		}
		newStocks = append(newStocks, newStock.Name, newStock.Price)
		i += 2
	}

	if len(query) < 5 {
		return
	}

	_, err = db.Exec(`
		INSERT INTO
			stock_data(stock_name, price) 
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
