package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	jwt "github.com/golang-jwt/jwt/v4"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

var mySigningKey = []byte("quakegodmode")

type Stock struct {
	Id      int             `json:"id"`
	Name    string          `json:"name"`
	Ticker  string          `json:"ticker"`
	Amount  int8            `json:"amount"`
	Buycost decimal.Decimal `json:"buycost"`
}

func checkAuth(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "close")
		defer r.Body.Close()

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("errror in header token parsing")
				}
				return mySigningKey, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				return
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "Not Authorizeddd")
		}
	})
}

func main() {
	fmt.Println("REST server here")
	r := mux.NewRouter()

	r.Handle("/stocks", checkAuth(getStocks))

	log.Fatal(http.ListenAndServe(":8081", r))

}

func getStocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var ArrStocks = myCurrenttocks()
	json.NewEncoder(w).Encode(ArrStocks)
}

func myCurrenttocks() []Stock {
	var allStocks []Stock
	db, err := sql.Open("mysql", "root:Super_22@tcp(127.0.0.1:3306)/mystocks")
	//when to use db?
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	res, err := db.Query("SELECT * FROM `usdstocks`")
	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		var stockrow Stock
		err = res.Scan(&stockrow.Id, &stockrow.Name, &stockrow.Ticker, &stockrow.Amount, &stockrow.Buycost)
		if err != nil {
			log.Fatal(err)
		}
		allStocks = append(allStocks, stockrow)
	}
	return allStocks
}
