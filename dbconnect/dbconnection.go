package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Coffee struct {
	ID    string `json:"id`
	Size  string `json:"size"`
	Name  string `json:"name"`
	Price string `json:"price`
}

func main() {
	db, sql_err := sql.Open("mysql", "root:root@tcp(172.17.0.2:3306)/starbucks")
	if sql_err != nil {
		fmt.Print("yo")
		log.Fatal(sql_err)

	}
	if err := db.Ping(); err != nil {
		fmt.Print("jotaro")
		log.Fatal(err)
	}
	fmt.Println("Sql Worked")
	// var temp Coffee
	query := `select * from coffees`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var temp []Coffee

	for rows.Next() {
		var x Coffee

		err := rows.Scan(&x.ID, &x.Size, &x.Name, &x.Price)
		if err != nil {
			log.Fatal(err)
			fmt.Print("Sql error fetching")
		}
		temp = append(temp, x)
	}

	fmt.Print(temp)
}
