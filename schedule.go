package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

type scheduler struct {
	TaskId    int   `json:"taskId"`
	TimeStamp int64 `json:"timestamp"`
	Flag      bool  `json:"flag"`
}

var db *sql.DB

func schedule(w http.ResponseWriter, r *http.Request){
	
}

func main() {
	r := mux.NewRouter()
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/todo")
	defer db.Close()
	r.HandleFunc("/schedule/{id:[0-9]+}", schedule)
	http.ListenAndServe(":8000", r)
}
