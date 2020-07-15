package main

import (
    "fmt"
    "log"
		"net/http"
		"encoding/json"
		"io/ioutil"
		"strconv"
		"os"

		"github.com/gorilla/mux"
		"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Title string `json:"Title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
}

type Articles []Article

func idParamToUint(r *http.Request) uint {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	var uid uint = uint(id)
	return uid
}

func ParseJsonArticle(w http.ResponseWriter ,r *http.Request) Article {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	json.NewEncoder(w).Encode(article)
	return article
}

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: returnAllArticles")
	db := DBConn()
	var articles Articles
	db.Find(&articles)
	json.NewEncoder(w).Encode(articles)
}	

func returnSingleArticle(w http.ResponseWriter, r *http.Request){
	fmt.Println("called returnSingleArticle")
	uid := idParamToUint(r)
	db := DBConn()
	var article Article
	db.Where("id = ?", uid).First(&article)
	json.NewEncoder(w).Encode(article)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called createNewArticle")
	db := DBConn()
	article := ParseJsonArticle(w, r)
	db.Create(&article)
	if db.NewRecord(article) {
		log.Fatalln("新規articleの保存に失敗しました。")
	}
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called deleteAtricle")
	uid := idParamToUint(r)
	db := DBConn()
	db.Delete(Article{}, "id = ?", uid)
}

func updateArticle(w http.ResponseWriter, r *http.Request){
	fmt.Println("called updateAtricle")
	uid := idParamToUint(r)
	db := DBConn()

	updatedArticle := ParseJsonArticle(w, r)

	var article Article
	db.Where("id = ?", uid).First(&article)

	article.Title = updatedArticle.Title
	article.Desc = updatedArticle.Desc
	article.Content = updatedArticle.Content
	db.Save(&article)
}


func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/article/{id}", returnSingleArticle).Methods("GET")
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

var (
	dbservice = "mysql"
	dbuser = os.Getenv("MINIMUM_APP_DATABASE_USER")
	dbpass = os.Getenv("MINIMUM_APP_DATABASE_PASS")
	dbname = os.Getenv("MINIMUM_APP_DEV_DATABASE_NAME")
)

func DBConn() *gorm.DB {
	d := Database{
		Service: dbservice,
		User: dbuser,
		Pass: dbpass,
		DatabaseName: dbname,
	}
	db := d.init()
	return db
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")

	d := Database{
		Service: dbservice,
		User: dbuser,
		Pass: dbpass,
		DatabaseName: dbname,
	}

	db, err := d.connect()
	if err != nil {
		log.Fatalln("データベースの接続に失敗しました。")
	}
	defer db.Close()
	db.AutoMigrate(&Article{})
	
	handleRequests()
}