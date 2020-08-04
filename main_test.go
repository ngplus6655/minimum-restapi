package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestValidateOfArticle(t *testing.T) {
	article := Article{Title:"Test", Desc:"test description",Content:"test content"}
	validArticle := article.validate()
	assert.Equal(t, validArticle, true, "Articleの検証が正しくありません")
	article.Title = ""
	article.Desc = strings.Repeat("a", 101)
	invalidArticle := article.validate()
	assert.Equal(t, invalidArticle, false, "Articleの検証が正しくありません")
}

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
	str := `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"title":"Test","desc":"test description","content":"test content"}`
	assert.Equal(t, str, resBody, "Jsonが正しくパースされませんでした。")
}

func TestReturnAllArticles(t *testing.T) {
	db := setFixture()
	defer cleanUpFixture(db)
	d := fetchTestDB()
	router := mux.NewRouter()
	router.HandleFunc("/articles", withVars(withDB(d, returnAllArticles)))

	req := httptest.NewRequest("GET", "/articles", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	resBodyByte, _ := ioutil.ReadAll(resp.Body)
	var articles Articles
	json.Unmarshal(resBodyByte, &articles)

	assert.Equal(t, resp.StatusCode, 200, "StatusCodeの値が正しくありません。")
	assert.Equal(t, "test1", articles[0].Title, "returnAllArticlesが正しい値を返しませんでした。")
	assert.Equal(t, "test description2", articles[1].Desc, "returnAllArticlesが正しい値を返しませんでした。")
	assert.Equal(t, "test content3", articles[2].Content, "returnAllArticlesが正しい値を返しませんでした。")
}

func TestReturnSingleArticle(t *testing.T) {
	db := setFixture()
	defer cleanUpFixture(db)
	router := mux.NewRouter()
	d := fetchTestDB()
	router.HandleFunc("/articles/{id}", withVars(withDB(d, returnSingleArticle)))

	req := httptest.NewRequest("GET", "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()
	resBodyByte, _ := ioutil.ReadAll(resp.Body)

	var article Article
	json.Unmarshal(resBodyByte, &article)

	assert.Equal(t, resp.StatusCode, 200, "StatusCodeの値が正しくありません。")
	assert.Equal(t, "test1", article.Title, "returnAllArticlesが正しい値を返しませんでした。")
	assert.Equal(t, "test description1", article.Desc, "returnAllArticlesが正しい値を返しませんでした。")
	assert.Equal(t, "test content1", article.Content, "returnAllArticlesが正しい値を返しませんでした。")
}

func TestCreateNewArticle(t *testing.T) {
	db := connTestDB()
	defer cleanUpFixture(db)
	d := fetchTestDB()
	router := mux.NewRouter()
	router.HandleFunc("/articles", withVars(withDB(d, createNewArticle)))

	reqBody := strings.NewReader(`{"Title":"PostTest","desc":"testing POST methods","content":"Hello world!!"}`)
	
	req := httptest.NewRequest("POST", "/articles", reqBody)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()
	assert.Equal(t, resp.StatusCode, 200, "StatusCodeの値が正しくありません。")

	var article Article
	db.First(&article)
	assert.Equal(t, "PostTest", article.Title, "Articleのタイトルの値が不正です")
}

func TestUpdateArticle(t *testing.T) {
	db := setFixture()
	defer cleanUpFixture(db)
	d := fetchTestDB()
	router := mux.NewRouter()
	router.HandleFunc("/articles/{id}", withVars(withDB(d, updateArticle)))

	reqBody := strings.NewReader(`{"Title":"PutTest","desc":"testing PUT methods","content":"UPDATED!!"}`)
	req := httptest.NewRequest("PUT", "/articles/1", reqBody)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()
	assert.Equal(t, resp.StatusCode, 200, "StatusCodeの値が正しくありません。")

	var article Article
	db.Where("id = ?", 1).First(&article)
	assert.Equal(t, "UPDATED!!", article.Content, "ArticleのContentの値が不正です")
}


func TestDeleteArticle(t *testing.T) {
	db := setFixture()
	defer cleanUpFixture(db)
	d := fetchTestDB()
	router := mux.NewRouter()
	router.HandleFunc("/articles/{id}", withVars(withDB(d, deleteArticle)))

	req := httptest.NewRequest("DELETE", "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()
	assert.Equal(t, resp.StatusCode, 200, "StatusCodeの値が正しくありません。")

	var article Article
	db.Where("id = ?", 1).First(&article)
	assert.Equal(t, uint(0), article.ID, "ArticleのIDの値が不正です")
}