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

// Data структура для представления данных из базы данных
type Data struct {
	Id      int    `json:"id"`
	Content string `json:"title"`
	Status  bool   `json:"completed"`
}

// data обработчик для получения данных из базы данных
func data(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Установка соединения с базой данных
	db, err := sql.Open("mysql", "your_user:your_password@tcp(host:port)/your_data_base")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Выполнение запроса к базе данных для получения данных
	rows, err := db.Query("SELECT * FROM Todos")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var dataList []Data

	// Итерация по результатам запроса и заполнение структуры Data
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

	// Преобразование данных в формат JSON и отправка клиенту
	bytes, err := json.Marshal(dataList)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)
}

// update обработчик для обновления данных в базе данных
func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Чтение тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Разбор JSON из тела запроса в структуру Data
	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
		return
	}

	// Установка соединения с базой данных
	db, err := sql.Open("mysql", "your_user:your_password@tcp(host:port)/your_data_base")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Преобразование статуса в строку
	var num string = strconv.Itoa(data.Id)
	var status string
	if data.Status {
		status = "1"
	} else {
		status = "0"
	}

	// Формирование и выполнение SQL-запроса на обновление данных
	var query = "UPDATE `todos` SET `completed`= " + status + " WHERE `id` =" + num
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Отправка успешного статуса
	w.WriteHeader(http.StatusOK)
}

// createReq обработчик для создания новой записи в базе данных
func createReq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Чтение тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Разбор JSON из тела запроса в структуру Data
	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
		return
	}

	// Установка соединения с базой данных
	db, err := sql.Open("mysql", "your_user:your_password@tcp(host:port)/your_data_base")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Формирование и выполнение SQL-запроса на создание новой записи
	query := "INSERT INTO `todos` (title, completed) VALUES ('" + data.Content + "', '" + strconv.FormatBool(data.Status) + "')"
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}

// delete обработчик для удаления записи из базы данных
func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Чтение тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Разбор JSON из тела запроса в структуру Data
	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
		return
	}

	// Установка соединения с базой данных
	db, err := sql.Open("mysql", "your_user:your_password@tcp(host:port)/your_data_base")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Преобразование ID в строку
	var num string = strconv.Itoa(data.Id)

	// Формирование и выполнение SQL-запроса на удаление записи
	var query = "DELETE FROM `todos` WHERE `id` =" + num
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// Отправка успешного статуса
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Создание маршрутизатора
	mux := http.NewServeMux()

	// Установка обработчиков для различных запросов
	mux.HandleFunc("/data", data)
	mux.HandleFunc("/update", update)
	mux.HandleFunc("/create", createReq)
	mux.HandleFunc("/delete", delete)

	// Настройка обработчика CORS
	handler := cors.Default().Handler(mux)

	// Запуск сервера
	err := http.ListenAndServe(PORT, handler)
	if err != nil {
		log.Println(err)
	}
}
