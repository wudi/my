package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"log"
	orm "my/database"
	"my/models"
	"my/result"
	"strconv"
)

type Document struct {
	Id         int    `json:"id"`
	BookName   string `json:"bookName"`
	BookIntro  string `json:"bookIntro"`
	BookAuthor string `json:"bookAuthor"`
}

type DocumentResponse struct {
	BookName   string `json:"bookName"`
	BookIntro  string `json:"bookIntro"`
	BookAuthor string `json:"bookAuthor"`
}
type SearchResponse struct {
	Time      string             `json:"time"`
	Hits      string             `json:"hits"`
	Documents []DocumentResponse `json:"documents"`
}

func Books(c *gin.Context) {
	pageNum := 1
	pageSize := 10
	if i, err := strconv.Atoi(c.Query("pageNum")); err == nil {
		pageNum = i
	}
	if i, err := strconv.Atoi(c.Query("pageSize")); err == nil {
		pageSize = i
	}
	results, err := models.Books(pageNum, pageSize)
	if err != nil {
		fmt.Print(err)
		result.Fail(c, err.Error())
		return
	}
	result.SuccessObj(c, results)
}

func BookQuery(c *gin.Context) {
	key := c.Query("query")
	if key == "" {
		result.Fail(c, "Query not specified")
		return
	}
	pageNum := 1
	pageSize := 10
	if i, err := strconv.Atoi(c.Query("pageNum")); err == nil {
		pageNum = i
	}
	if i, err := strconv.Atoi(c.Query("pageSize")); err == nil {
		pageSize = i
	}
	esQuery := elastic.NewMultiMatchQuery(key, "BookName", "BookIntro", "BookAuthor").
		Fuzziness("2")
	searchResult, err := orm.Es.Search().Index("book").Query(esQuery).
		From(pageNum - 1).Size(pageSize).
		Do(c.Request.Context())
	if err != nil {
		log.Println(err)
		result.Fail(c, "Something went wrong")
		return
	}
	res := SearchResponse{
		Time: fmt.Sprintf("%d", searchResult.TookInMillis),
		Hits: fmt.Sprintf("%d", searchResult.Hits.TotalHits),
	}
	docs := make([]DocumentResponse, 0)
	for _, hit := range searchResult.Hits.Hits {
		var doc DocumentResponse
		json.Unmarshal(hit.Source, &doc)
		docs = append(docs, doc)
	}
	res.Documents = docs
	result.SuccessObj(c, res)
	return
}

func BookSectionByNum(c *gin.Context) {
	num := c.Param("num")
	if num == "" {
		result.Fail(c, "num not specified")
		return
	}
	result.Success(c, num)
}

func BookById(c *gin.Context) {
	bookId := c.Param("bookId")
	if bookId == "" {
		result.Fail(c, "bookId not specified")
		return
	}
	book, err := models.GetBookByIdFromEs(bookId)
	if err != nil {
		result.Fail(c, err.Error())
	}
	result.SuccessObj(c, book)
}
