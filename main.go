package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

func prep8079() {
	port := 8079

	serverMuxA := http.NewServeMux()
	serverMuxA.HandleFunc("/simpleJson", simpleJson)

	go func(message string) {
		log.Printf("Server starting on port %v\n", 8079)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), serverMuxA))
	}("second_server")
}

func prep8081() {
	port := 8081

	http.HandleFunc("/simpleJson", simpleJson2)
	http.HandleFunc("/dbResponse", dbResponse)

	log.Printf("Server starting on port %v\n", 8081)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func main() {
	prep8079()
	prep8081()
}

func simpleJson(res http.ResponseWriter, req *http.Request) {
	type fruitsResponse struct {
		Page   int
		Fruits []string
	}

	response := fruitsResponse{
		Page:   1,
		Fruits: []string{"Aaa", "bbb"},
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusBadGateway)

	res.Write(jsonResponse)
}

func decodeJsonData(model interface{}, req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(&model)
}

func sendJsonResponse(models interface{}, res http.ResponseWriter) {
	jsonResponse, err := json.Marshal(models)
	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write(jsonResponse)
}

type Owner struct {
	Owner string
	Name  string
}

type Author struct {
	Author string
	Name   string
}

func simpleJson2(res http.ResponseWriter, req *http.Request) {

	var model Author
	err := decodeJsonData(&model, req)

	if err != nil {
		panic(err)
	}

	log.Printf("Owner" + model.Author)

	jsonResponse, err := json.Marshal(model)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)

	res.Write(jsonResponse)
}

func dbResponse(res http.ResponseWriter, req *http.Request) {

	type Language struct {
		gorm.Model
		Name string
	}

	type Product struct {
		gorm.Model
		Code      string
		Price     uint
		Languages []Language `gorm:"many2many:product_languages;"`
	}

	const addr = "postgresql://maxroach@localhost:26257/gorm?ssl=false&sslmode=disable"
	db, err := gorm.Open("postgres", addr)

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.LogMode(true)
	//db.AutoMigrate(&Product{})
	//db.AutoMigrate(&Language{})
	//
	//db.Create(&Product{Code: "L1212", Price: 1000})

	//db.Create(&Product{Code: "L1212", Price: 1000}).Association("Languages").Append(&Language{Name: "EN"})

	var products []Product
	db.Preload("Languages").Where("ID > ?", "1").Find(&products)

	sendJsonResponse(products, res)
	//db.First(&product, "code = ?", "L1212")

	//
	//db.Model(&product).Update("Price", 2000)
	//
	//db.Delete(&product)
}
