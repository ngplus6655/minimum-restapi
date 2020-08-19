package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"net/http"
	"encoding/base64"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestValidateOfArticle(t *testing.T) {
	article := Article{Title: "Test", Desc: "test description", Content: "test content"}
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
	getArticle := ParseJsonArticle(w, req)
	article := Article{
		Title: "Test",
		Desc: "test description",
		Content: "test content",
	}
	assert.Equal(t, article, getArticle, "Jsonが正しくパースされませんでした。")
}

func TestSetFlashMessage(t *testing.T) {
	w := httptest.NewRecorder()
	str := "message"
	setFlashMessage(w, str)
	cookie := []string([]string{"message=bWVzc2FnZQ=="})
	assert.Equal(t, cookie, w.HeaderMap["Set-Cookie"], "cookieがうまく設定されませんでした")
}

func TestGETAllArticles(t *testing.T) {
	db := setFixture()
	defer cleanUpFixture(db)
	d := fetchTestDB()
	r := mux.NewRouter()
	r.HandleFunc("/articles", withVars(withDB(d, articlesCORSHandling))).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	req := httptest.NewRequest("GET", "/articles", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()

	resBodyByte, _ := ioutil.ReadAll(resp.Body)
	var articles Articles
	json.Unmarshal(resBodyByte, &articles)

	assert.Equal(t, resp.StatusCode, 200, "StatusCodeの値が正しくありません。")
	assert.Equal(t, "test1", articles[0].Title, "returnAllArticlesが正しい値を返しませんでした。")
	assert.Equal(t, "test description2", articles[1].Desc, "returnAllArticlesが正しい値を返しませんでした。")
	assert.Equal(t, "test content3", articles[2].Content, "returnAllArticlesが正しい値を返しませんでした。")
}

func TestGETSingleArticle(t *testing.T) {
	db := setFixture()
	defer cleanUpFixture(db)
	r := mux.NewRouter()
	d := fetchTestDB()
	r.HandleFunc("/articles/{id}", withVars(withDB(d, articlesCORSHandlingWithID))).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)

	req := httptest.NewRequest("GET", "/articles/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	resBodyByte, _ := ioutil.ReadAll(resp.Body)

	var article Article
	json.Unmarshal(resBodyByte, &article)

	assert.Equal(t, resp.StatusCode, 200, "StatusCodeの値が正しくありません。")
	assert.Equal(t, "test1", article.Title, "returnAllArticlesが正しい値を返しませんでした。")
	assert.Equal(t, "test description1", article.Desc, "returnAllArticlesが正しい値を返しませんでした。")
	assert.Equal(t, "test content1", article.Content, "returnAllArticlesが正しい値を返しませんでした。")
}

func TestPOSTNewArticle(t *testing.T) {
	db := connTestDB()
	defer cleanUpFixture(db)
	d := fetchTestDB()
	r := mux.NewRouter()
	r.HandleFunc("/articles", withVars(withDB(d, articlesCORSHandling))).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	reqBody := strings.NewReader(`{"Title":"PostTest","desc":"testing POST methods","content":"Hello world!!"}`)

	req := httptest.NewRequest("POST", "/articles", reqBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	value64 := base64.StdEncoding.EncodeToString([]byte("保存に成功しました"))
	cookie := []string([]string{"message=" + value64})
	assert.Equal(t, resp.Header["Set-Cookie"], cookie, "StatusCodeの値が正しくありません。")

	var article Article
	db.First(&article)
	assert.Equal(t, "PostTest", article.Title, "Articleのタイトルの値が不正です")
}

func TestPUTArticle(t *testing.T) {
	db := setFixture()
	defer cleanUpFixture(db)
	d := fetchTestDB()
	r := mux.NewRouter()
	r.HandleFunc("/articles/{id}", withVars(withDB(d, articlesCORSHandlingWithID))).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)

	reqBody := strings.NewReader(`{"Title":"PutTest","desc":"testing PUT methods","content":"UPDATED!!"}`)
	req := httptest.NewRequest("PUT", "/articles/1", reqBody)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	assert.Equal(t, resp.StatusCode, 200, "StatusCodeの値が正しくありません。")

	var article Article
	db.Where("id = ?", 1).First(&article)
	assert.Equal(t, "UPDATED!!", article.Content, "ArticleのContentの値が不正です")
}

func TestDELETEArticle(t *testing.T) {
	db := setFixture()
	defer cleanUpFixture(db)
	d := fetchTestDB()
	r := mux.NewRouter()
	r.HandleFunc("/articles/{id}", withVars(withDB(d, articlesCORSHandlingWithID))).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)

	req := httptest.NewRequest("DELETE", "/articles/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	resp := w.Result()
	assert.Equal(t, resp.StatusCode, 200, "StatusCodeの値が正しくありません。")

	var article Article
	db.Where("id = ?", 1).First(&article)
	assert.Equal(t, uint(0), article.ID, "ArticleのIDの値が不正です")
}
