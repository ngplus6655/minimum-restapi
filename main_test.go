package main

import (
	"net/http/httptest"
	"reflect"
	"testing"
	"io/ioutil"
	"strings"
	"log"
	"os"
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
)

func TestIdParamToUint(t *testing.T) {
	request := httptest.NewRequest("GET", "/article/123", nil)
	id := idParamToUint(request)
	assert.Equal(t, reflect.Uint, reflect.TypeOf(id).Kind(), "Uint型が返されませんでした。")
}

func TestParseJsonArticle(t *testing.T) {
	reqBody := strings.NewReader(`{"Title": "Test","desc": "test description","content": "test content"}`)
	req := httptest.NewRequest("GET", "/", reqBody)
	w := httptest.NewRecorder()
	ParseJsonArticle(w, req)
	resp := w.Result()
	resBodyByte, _ := ioutil.ReadAll(resp.Body)
	resBody := strings.ReplaceAll(string(resBodyByte), "\n", "")
	str := `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"Title":"Test","desc":"test description","content":"test content"}`
	assert.Equal(t, str, resBody, "Jsonが正しくパースされませんでした。")
}

func SetFixture() *gorm.DB {
	dbname = os.Getenv("MINIMUM_APP_TEST_DATABASE_NAME")
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
	db.AutoMigrate(&Article{})

	articles := Articles{
		Article{Title:"test1",Desc:"test description1",Content: "test content1",},
		Article{Title:"test2",Desc:"test description2",Content: "test content2",},
		Article{Title:"test1",Desc:"test description3",Content: "test content3",},
	}
	for _, article := range articles {
		db.Create(&article)
	}
	return db
}

func CleanUpFixture(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE articles;")
	db.Close()
}

func TestReturnAllArticles(t *testing.T){
	db := SetFixture()
	req := httptest.NewRequest("GET", "/all", nil)
	w := httptest.NewRecorder()
	returnAllArticles(w, req)
	resp := w.Result()
	resBodyByte, _ := ioutil.ReadAll(resp.Body)
	var articles Articles
	json.Unmarshal(resBodyByte, &articles)

	assert.Equal(t, resp.StatusCode, 200, "StatusCodeの値が正しくありません。")
	assert.Equal(t, "test1", articles[0].Title, "returnAllArticlesが正しい値を返しませんでした。")
	assert.Equal(t, "test description2", articles[1].Desc, "returnAllArticlesが正しい値を返しませんでした。")
	assert.Equal(t, "test content3", articles[2].Content, "returnAllArticlesが正しい値を返しませんでした。")
	
	defer CleanUpFixture(db)
}
