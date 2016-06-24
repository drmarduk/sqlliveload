package main

import (
	"database/sql"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func main() {
	var err error

	db, err = sql.Open("sqlite3", "foo.db")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/query", QueryHandler)
	http.Handle("/", http.FileServer(http.Dir("./html/")))
	http.ListenAndServe(":8888", nil)
}

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	for _, v := range r.Form {

		print(v)
	}
	q := r.FormValue("query")
	print(q + "\n")
	if q == "" {
		return
	}

	result, err := queryDb(q)
	if err != nil {
		return
	}

	for _, x := range result {
		w.Write([]byte(x + "\n"))
	}
}

func queryDb(q string) ([]string, error) {

	query := "select name from files where instr(name, ?) limit 0, 2000;"
	rows, err := db.Query(query, q)
	if err != nil {
		panic(err)
	}

	var result []string

	for rows.Next() {
		var x string
		err = rows.Scan(&x)
		if err != nil {
			print("err while scanning: " + q + "\n")
			return result, err
		}
		result = append(result, x)
	}

	return result, nil
}
