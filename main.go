package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/cors"

	_ "github.com/go-sql-driver/mysql"
)

const PORT = ":8080"

type Data struct {
	Id      int    `json:"id"`
	Content string `json:"title"`
	Status  bool   `json:"completed"`
}

func data(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, err := sql.Open("mysql", "your_user:your_password@tcp(host:port)/your_data_base")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM Todos")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var dataList []Data

	for rows.Next() {
		var column1 int
		var column2 string
		var column3 bool
		err := rows.Scan(&column1, &column2, &column3)
		if err != nil {
			panic(err)
		}
		data := Data{
			Id:      column1,
			Content: column2,
			Status:  column3,
		}
		dataList = append(dataList, data)
	}

	bytes, err := json.Marshal(dataList)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("mysql", "your_user:your_password@tcp(host:port)/your_data_base")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err != nil {
		panic(err)
	}
	var num string = strconv.Itoa(data.Id)
	var status string
	if data.Status {
		status = "1"
	} else {
		status = "0"
	}
	var query = "UPDATE `todos` SET `completed`= " + status + " WHERE `id` =" + num
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	w.WriteHeader(http.StatusOK)
}

func createReq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("mysql", "your_user:your_password@tcp(host:port)/your_data_base")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err != nil {
		log.Fatal(err)
	}
	query := "INSERT INTO `todos` (title, completed) VALUES ('" + data.Content + "', '" + strconv.FormatBool(data.Status) + "')"

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "your_user:your_password@tcp(host:port)/your_data_base")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err != nil {
		log.Fatal(err)
	}

	var num string = strconv.Itoa(data.Id)
	var query = "DELETE FROM `todos` WHERE `id` =" + num
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	db.Close()
	w.WriteHeader(http.StatusOK)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/data", data)
	mux.HandleFunc("/update", update)
	mux.HandleFunc("/create", createReq)
	mux.HandleFunc("/delete", delete)
	handler := cors.Default().Handler(mux)
	err := http.ListenAndServe(PORT, handler)
	if err != nil {
		log.Println(err)
	}
}
