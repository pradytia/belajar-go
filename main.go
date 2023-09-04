package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	dbDriver   = "mysql"
	dbUser     = "root"
	dbPassword = "pass"
	dbName     = "my_db"
)

//func dbConn() (db *sql.DB) {
//	dataSourceName := fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName)
//	db, err := sql.Open(dbDriver, dataSourceName)
//	if err != nil {
//		panic(err.Error())
//	}
//	return db
//}

var db *sql.DB

func init() {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName)
	db, err = sql.Open(dbDriver, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetItems(res http.ResponseWriter, req *http.Request) {
	//db := dbConn()
	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var array []Item

	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			log.Fatal(err)
		}
		array = append(array, item)
	}

	res.Header().Set("Content-Type", "application/json")
	errJson := json.NewEncoder(res).Encode(array)
	if errJson != nil {
		return
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", GetItems).Methods("GET")
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
