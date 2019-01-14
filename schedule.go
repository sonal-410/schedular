package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type scheduler struct {
	TaskID    int   `json:"taskId"`
	TimeStamp int64 `json:"timeStamp"`
	Flag      bool  `json:"flag"`
}

var db *sql.DB

func schedule(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var q scheduler
		_ = json.NewDecoder(r.Body).Decode(&q)
		fmt.Println(q)
		stmtIns, er := db.Prepare("insert into scheduler values(?,?,?)")
		if er != nil {
			panic(er.Error())
		} else {
			fmt.Println("entity add")
		}
		_, err := stmtIns.Exec(q.TimeStamp, q.Flag, q.TaskID)
		if err != nil {
			panic(err.Error())
		}
	case "GET":
		t := time.Now().Unix()
		//fmt.Println(t)
		// keys := r.URL.Query()
		vars := mux.Vars(r)
		id := vars["id"]

		var num int64
		err := db.QueryRow("select TimeStamp from scheduler where TaskID = ?", id).Scan(&num)
		if err != nil {
			log.Fatal(err)
		}
		if num > t {
			fmt.Fprintf(w, "You have %d seconds remaining for the task", num-t)
		}
		if num < t {
			var task int
			fmt.Fprintf(w, "The task time has been passed")
			rows, err := db.Query("UPDATE scheduler SET Flag = true WHERE TaskID = ?", id)
			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(&task)
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(w)
			}
			//fmt.Println(rows)
			fmt.Println(task)
			err = rows.Err()
			if err != nil {
				panic(err)
			}
			//fmt.Println("GET success")
			if err != nil {
				log.Fatal(err)
			}
			_, er := db.Query("Delete from scheduler where TaskID = ?", id)
			if er != nil {
				log.Fatal(er)
			}
		}
		if num == t {
			var task int
			fmt.Fprintf(w, "Its time for the task")
			rows, err := db.Query("UPDATE scheduler SET Flag = true WHERE TaskID = ?", id)
			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(&task)
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(w)
			}
			//fmt.Println(rows)
			//fmt.Println(task)
			err = rows.Err()
			if err != nil {
				panic(err)
			}
			//fmt.Println("GET success")
			if err != nil {
				log.Fatal(err)
			}
			_, er := db.Query("Delete from scheduler where TaskID = ?", id)
			if er != nil {
				log.Fatal(er)
			}
		}
		//	fmt.Println(num)
	}
}
func finishedTasks() {
	var task int64

	t := time.Now().Unix()
	rows, err := db.Query("select TaskID from  scheduler where Flag = false and ?>TimeStamp ", t)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&task)
		if err != nil {
			panic(err)
		}
		fmt.Println(task)
	}
	//fmt.Println(rows)
	fmt.Println(task)
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	//fmt.Println("GET success")
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	r := mux.NewRouter()
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/todo")
	defer db.Close()
	r.HandleFunc("/schedule/{id:[0-9]+}", schedule)
	finishedTasks()
	http.ListenAndServe(":8000", r)

}
