package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	s "example.com/web-service-gin/structs"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var db *sqlx.DB

func init() {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	driver, IsExistDbDriver := os.LookupEnv("DB_DRIVER")
	db_login, IsExistDbLogin := os.LookupEnv("DB_LOGIN")
	db_pass, IsExistDbPass := os.LookupEnv("DB_PASS")
	db_table, IsExistDbTable := os.LookupEnv("DB_TABLE")

	if !IsExistDbDriver && !IsExistDbLogin && !IsExistDbPass && !IsExistDbTable {
		log.Fatal("failed to connect .evn file")
	}
	db, err = sqlx.Open(driver, db_login+":"+db_pass+"@/"+db_table)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	if err != nil {
		log.Fatal("failed to connect database")
	}
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

}

func UpdateNew(psNew s.PostNews, newId string) {

	tx := db.MustBegin()
	if psNew.Title != "" {
		tx.MustExec("UPDATE News SET Title=? WHERE Id=?", psNew.Title, newId)
	}
	if psNew.Content != "" {
		tx.MustExec("UPDATE News SET Content=? WHERE Id=?", psNew.Content, newId)
	}

	if psNew.Categories != nil {
		tx.MustExec("DELETE FROM NewsCategories WHERE NewsId=?", newId)
		for _, value := range psNew.Categories {
			tx.MustExec("INSERT INTO NewsCategories (NewsId, CategoryId) VALUES (?, ?)", newId, value)
		}

	}
	tx.Commit()

}

func IssetNew(id string) bool {

	var new s.News2
	err := db.Get(&new, "SELECT * FROM News WHERE Id =?", id)

	return err == nil
}

func main() {

	router := gin.Default()

	router.GET("/list", getListNews)
	router.POST("/edit/:id", editNew)

	router.Run(":8001")
}
func editNew(c *gin.Context) {
	newId := c.Param("id")
	var psNew s.PostNews
	c.Bind(&psNew)

	if IssetNew(newId) {
		UpdateNew(psNew, newId)
		i, _ := strconv.ParseInt(newId, 10, 64)
		c.JSON(http.StatusOK, getNewById(i))
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No new found"})
		return
	}
}
func getListNews(c *gin.Context) {
	var newChange []s.NewsResult

	err := db.Select(&newChange, "SELECT * FROM News")

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(newChange) > 0 {
		for index, item := range newChange {
			newChange[index].Categories = getListCategories(item.Id)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Success": true,
		"News":    newChange,
	})

}

func getListCategories(id int64) []int64 {
	var cat []int64
	var curCat []int64
	err := db.Select(&curCat, "SELECT CategoryId FROM NewsCategories WHERE NewsId =?", id)
	log.Println(curCat)
	if err != nil {
		log.Println(err)
	}
	cat = append(cat, curCat...)

	return cat
}

func getNewById(id int64) s.NewsResult {
	var new s.News2
	var newChange s.NewsResult
	err := db.Get(&new, "SELECT * FROM News WHERE Id =?", id)

	if err != nil {
		log.Println(err)
	}

	newChange.Id = new.Id
	newChange.Title = new.Title
	newChange.Content = new.Content
	newChange.Categories = getListCategories(id)

	return newChange

}
