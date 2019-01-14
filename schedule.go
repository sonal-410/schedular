package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type scheduler struct {
	TaskID    int   `json:"taskId"`
	TimeStamp int64 `json:"timestamp"`
	Flag      bool  `json:"flag"`
}

var db *sql.DB

func schedule(w http.ResponseWriter, r *http.Request) {
	var q scheduler
	_ = json.NewDecoder(r.Body).Decode(&q)
	fmt.Println(q)
	stmtIns, er := db.Prepare("insert into scheduler values(?,?,?)")
	if er != nil {
		panic(er.Error())
	} else {
		fmt.Println("entity add")
	}
	_, err := stmtIns.Exec(q.TaskID, q.TimeStamp, q.Flag)
	if err != nil {
		panic(err.Error())
	}

}

func main() {
	r := mux.NewRouter()
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/todo")
	defer db.Close()
	r.HandleFunc("/schedule/{id:[0-9]+}", schedule).Methods("GET")
	http.ListenAndServe(":8000", r)
}
