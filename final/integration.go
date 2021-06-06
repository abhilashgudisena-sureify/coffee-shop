package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Coffee struct {
	ID    string `json:"id`
	Size  string `json:"size"`
	Name  string `json:"name"`
	Price string `json:"price`
}

var coffees []Coffee
var db, sql_err = sql.Open("mysql", "root:root@tcp(172.17.0.2:3306)/starbucks")

func fetchData() []Coffee {

	var temp []Coffee
	query := "SELECT * FROM coffees "
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var x Coffee

		err := rows.Scan(&x.ID, &x.Size, &x.Name, &x.Price)
		if err != nil {
			log.Fatal(err)
			fmt.Print("Sql error fetching")
		}
		temp = append(temp, x)
	}
	return temp
}

//Get all
func getCoffees(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	coffees = fetchData()
	json.NewEncoder(response).Encode(coffees)

}

//Get Coffee
func getCoffee(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(request)
	coffees = fetchData()
	for _, item := range coffees {
		if item.Name == params["name"] {
			json.NewEncoder(response).Encode(item)
			return
		}
	}

}

//Add Coffee
func addCoffee(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	var newCoffee Coffee
	json.NewDecoder(request.Body).Decode(&newCoffee)
	if newCoffee.ID != "0" {
		newCoffee.ID = strconv.Itoa(len(coffees) + 1)
	}
	result, err := db.Exec(`INSERT INTO coffees (id,name,size,price) VALUES(?,?,?,?)`, newCoffee.ID, newCoffee.Name, newCoffee.Size, newCoffee.Price)

	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	fmt.Println("inserted row:", id)
	// coffees = append(coffees, newCoffee)
	json.NewEncoder(response).Encode(newCoffee)
}

//Update Coffee
func updateCoffee(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(request)
	for i, item := range coffees {
		if item.Name == params["name"] {
			coffees = append(coffees[:i], coffees[i+1:]...)
			var newCoffee Coffee
			json.NewDecoder(request.Body).Decode(&newCoffee)
			newCoffee.Name = params["name"]
			coffees = append(coffees, newCoffee)
			json.NewEncoder(response).Encode(newCoffee)
			return
		}
	}

}

//Delete Coffee
func deleteCoffee(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	// params := mux.Vars(request)
	var newCoffee Coffee
	json.NewDecoder(request.Body).Decode(&newCoffee)
	fmt.Print("params:", newCoffee.ID)
	_, err := db.Exec(`DELETE FROM coffees WHERE id = ?`, newCoffee.ID) // check err

	if err != nil {
		log.Fatal(err)
	}

	// for i, item := range coffees {
	// 	if item.ID == params["id"] {
	// 		coffees = append(coffees[:i], coffees[i+1:]...)
	// 		break
	// 	}
	// }
	coffees = fetchData()
	json.NewEncoder(response).Encode(coffees)

}

func main() {

	// db, sql_err = sql.Open("mysql", "root:root@tcp(172.17.0.2:3306)/starbucks")
	if sql_err != nil {
		fmt.Print("yo")
		log.Fatal(sql_err)

	}
	if err := db.Ping(); err != nil {
		fmt.Print("jotaro")
		log.Fatal(err)
	}

	fmt.Print("yo")

	coffees = append(coffees, Coffee{ID: "1", Size: "Small", Name: "Latte", Price: "2.99"})
	coffees = append(coffees, Coffee{ID: "2", Size: "Large", Name: "Cappucino", Price: "1.99"})
	coffees = append(coffees, Coffee{ID: "3", Size: "Medium", Name: "Americano", Price: "5.99"})
	handler := mux.NewRouter()

	//endpoints
	handler.HandleFunc("/coffee", getCoffees).Methods("GET")

	handler.HandleFunc("/coffee/{name}", getCoffee).Methods("GET")

	handler.HandleFunc("/coffee", addCoffee).Methods("POST")

	handler.HandleFunc("/coffee/{name}", updateCoffee).Methods("PUT")

	handler.HandleFunc("/coffee", deleteCoffee).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8083", handler))
}
