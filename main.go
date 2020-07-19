package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type Articles []Article

func idParamToUint(r *http.Request) uint {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	var uid uint = uint(id)
	return uid
}

func ParseJsonArticle(w http.ResponseWriter, r *http.Request) Article {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	json.NewEncoder(w).Encode(article)
	return article
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	log.Println("called returnSingleArticle")
	d := GetVar(r, "db").(Database)
	db := d.init()
	var articles Articles
	db.Find(&articles)
	json.NewEncoder(w).Encode(articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	log.Println("called returnSingleArticle")
	uid := idParamToUint(r)
	d := GetVar(r, "db").(Database)
	db := d.init()
	var article Article
	db.Where("id = ?", uid).First(&article)
	json.NewEncoder(w).Encode(article)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	log.Println("called createNewArticle")
	d := GetVar(r, "db").(Database)
	db := d.init()
	article := ParseJsonArticle(w, r)
	db.Create(&article)
	if db.NewRecord(article) {
		log.Println("新規articleの保存に失敗しました。")
	}
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	log.Println("called updateAtricle")
	uid := idParamToUint(r)
	d := GetVar(r, "db").(Database)
	db := d.init()

	updatedArticle := ParseJsonArticle(w, r)

	var article Article
	db.Where("id = ?", uid).First(&article)

	article.Title = updatedArticle.Title
	article.Desc = updatedArticle.Desc
	article.Content = updatedArticle.Content
	db.Save(&article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	log.Println("called deleteAtricle")
	uid := idParamToUint(r)
	d := GetVar(r, "db").(Database)
	db := d.init()
	db.Delete(Article{}, "id = ?", uid)
}

func setDevDB() Database {
	var (
		dbservice = "mysql"
		dbuser    = os.Getenv("MINIMUM_APP_DATABASE_USER")
		dbpass    = os.Getenv("MINIMUM_APP_DATABASE_PASS")
		dbname    = os.Getenv("MINIMUM_APP_DEV_DATABASE_NAME")
	)

	d := Database{
		Service:      dbservice,
		User:         dbuser,
		Pass:         dbpass,
		DatabaseName: dbname,
	}
	return d
}

func handleRequests() {
	db := setDevDB()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", withVars(withDB(db, returnAllArticles)))
	myRouter.HandleFunc("/article/{id}", withVars(withDB(db, returnSingleArticle))).Methods("GET")
	myRouter.HandleFunc("/article", withVars(withDB(db, createNewArticle))).Methods("POST")
	myRouter.HandleFunc("/article/{id}", withVars(withDB(db, updateArticle))).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", withVars(withDB(db, deleteArticle))).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	log.Println("Rest API v2.0 - Mux Routers")

	d := setDevDB()
	db := d.migrate()
	defer db.Close()

	handleRequests()
}
