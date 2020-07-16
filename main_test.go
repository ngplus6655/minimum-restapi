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

func TestReturnSingleArticle(t *testing.T) {
	db := SetFixture()
	router := mux.NewRouter()
	router.HandleFunc("/article/{id}", returnSingleArticle)

	req := httptest.NewRequest("GET", "/article/1", nil)
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

	defer CleanUpFixture(db)
}

