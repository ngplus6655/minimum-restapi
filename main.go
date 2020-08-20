package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"encoding/base64"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

func (a Article) validate() (valid bool) {
	valid = true
	if len(a.Title) == 0 || len(a.Title) > 30 {
		valid = false
	}
	if len(a.Desc) == 0 || len(a.Desc) > 100 {
		valid = false
	}
	if len(a.Content) == 0 || len(a.Content) > 100 {
		valid = false
	}
	return valid
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
	return article
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
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

func setFlashMessage(w http.ResponseWriter, str string){
	value64 := base64.StdEncoding.EncodeToString([]byte(str))
	cookie := &http.Cookie{
		Name: "message",
		Value: value64,
	}
	http.SetCookie(w, cookie)
}

func articlesCORSHandling(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set( "Access-Control-Allow-Headers", "Origin, Content-Type,")
	w.Header().Set( "Access-Control-Allow-Methods","GET, POST, OPTIONS" )

	if r.Method == http.MethodGet {
		log.Println("called GET /articles")
		d := GetVar(r, "db").(Database)
		db := d.init()
		var articles Articles
		db.Find(&articles)
		json.NewEncoder(w).Encode(articles)
	}

	if r.Method == http.MethodPost {
		error := false
		log.Println("called POST /articles")
		d := GetVar(r, "db").(Database)
		db := d.init()
		article := ParseJsonArticle(w, r)
		if valid := article.validate(); valid {
			db.Create(&article)
			if db.NewRecord(article) == false {
				log.Println("新規articleの保存に成功しました")
			} else if db.NewRecord(article) == true {
				error = true
			}
		} else {
			error = true
		}
		if error {
			log.Println("新規articleの保存に失敗しました")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func articlesCORSHandlingWithID(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS" )

	if r.Method == http.MethodGet {
		uid := idParamToUint(r)
		log.Println("called GET /article/" + strconv.FormatUint(uint64(uid), 10))
		d := GetVar(r, "db").(Database)
		db := d.init()
		var article Article
		db.Where("id = ?", uid).First(&article)
		json.NewEncoder(w).Encode(article)
	}

	if r.Method == http.MethodPut {
		uid := idParamToUint(r)
		log.Println("called PUT article/" + strconv.FormatUint(uint64(uid), 10))
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

	if r.Method == http.MethodDelete {
		uid := idParamToUint(r)
		log.Println("called DELETE article/" + strconv.FormatUint(uint64(uid), 10))
		d := GetVar(r, "db").(Database)
		db := d.init()
		db.Delete(Article{}, "id = ?", uid)
	}
}

func handleRequests() {
	db := setDevDB()
	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/articles", withVars(withDB(db, articlesCORSHandling))).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.HandleFunc("/articles/{id}", withVars(withDB(db, articlesCORSHandlingWithID))).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))
	log.Fatal(http.ListenAndServe(":10000", r))
}

func main() {
	log.Println("Rest API v2.0 - Mux Routers")

	d := setDevDB()
	db := d.migrate()
	defer db.Close()

	handleRequests()
}
