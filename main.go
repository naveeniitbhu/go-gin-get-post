package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	R  *gin.Engine
	Db *sql.DB
}

type Quiz struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type Question struct {
	ID             int64  `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	Options        string `json:"options" db:"options"`
	Correct_Option int64  `json:"correct_option" db:"correct_option"`
	Quiz           int64  `json:"quiz" db:"quiz"`
	Points         int64  `json:"points" db:"points"`
}

type Id struct {
	ID int64 `json:"id,omitempty" db:"id"`
}

const (
	failureMsgNotFound = "no rows in result set"
	status             = "status"
	failure            = "failure"
	reason             = "reason"
	explaination       = "explaination"
)

func main() {

	db, err := sql.Open("sqlite3", "./db/migration-goose/memory-sqlite.db")

	if err != nil {
		log.Fatalf("cannot open an sqlite3 based database: %v", err)
	} else {
		log.Println("INFO: successful, connected to database")
	}

	defer db.Close()

	app := App{
		R:  gin.Default(),
		Db: db,
	}

	app.R.GET("/api/quiz/:quiz_id", app.GetQuizId)
	app.R.POST("/api/quiz/", app.PostQuizDetails)

	app.R.GET("/api/question/:question_id", app.GetQuestion)
	app.R.GET("/api/questions/:question_id", app.GetQuestion)
	app.R.POST("/api/questions/", app.PostQuestionDetails)

	app.R.GET("/api/quiz-questions/:quiz_id", app.GetAllQuestions)

	app.R.Run(":8080")
}

func (app *App) GetQuizId(c *gin.Context) {
	id := c.Param("quiz_id")

	db := app.Db
	var quiz Quiz

	err := db.QueryRow(`Select id,name,description From quiz WHERE id=?`, id).Scan(&quiz.ID, &quiz.Name, &quiz.Description)
	if err != nil {
		log.Println(err.Error())
		c.JSON(404, gin.H{
			status: failure,
			reason: failureMsgNotFound,
		})
	} else {
		if quiz.ID < 1 {
			c.JSON(404, gin.H{
				status: failure,
				reason: failureMsgNotFound,
			})
		} else {
			log.Println("INFO: successfully retrieved quiz details")
			c.JSON(http.StatusOK, gin.H{
				"id":          id,
				"name":        quiz.Name,
				"description": quiz.Description,
			})
		}
	}
}

func (app *App) PostQuizDetails(c *gin.Context) {
	var (
		quiz Quiz
		id   int64
	)
	db := app.Db
	if err := c.ShouldBindJSON(&quiz); err == nil {
		if quiz.Name != "" && quiz.Description != "" {
			log.Println("Error: json binding is successful")
		} else {
			log.Println("Error in json binding")
			c.JSON(http.StatusBadRequest, gin.H{
				status:       failure,
				explaination: "name and description fields are required"})
			panic(err.Error())
		}
	} else {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			status:       failure,
			explaination: fmt.Sprintf("Invalid Input: %s", err.Error()),
		})
		panic(err)
	}

	res, err := db.Exec(`INSERT INTO quiz(name, description) VALUES(?, ?)`, quiz.Name, quiz.Description)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			status:       failure,
			explaination: err.Error(),
		})
	}
	id, err = res.LastInsertId()

	if err == nil {
		log.Printf("INFO: Quiz details inserted with id:%d & name:%s & description:%s\n", id, quiz.Name, quiz.Description)
		c.JSON(http.StatusCreated, gin.H{
			"id":          id,
			"name":        quiz.Name,
			"description": quiz.Description,
		})
	} else {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			status:       failure,
			explaination: fmt.Sprintf("Quiz Insert Error: %s", err.Error()),
		})
		panic(err)
	}
}

func (app *App) GetQuestion(c *gin.Context) {
	id := c.Param("question_id")

	db := app.Db
	var question Question

	err := db.QueryRow(`Select name, options, correct_option, quiz, points From questions WHERE id=?`, id).Scan(&question.Name, &question.Options, &question.Correct_Option, &question.Quiz, &question.Points)
	if err != nil {
		log.Println(err.Error())
		c.JSON(404, gin.H{
			status: failure,
			reason: failureMsgNotFound,
		})
	} else {
		log.Println("INFO: successful retrieved question details")
		c.JSON(200, gin.H{
			"id":             id,
			"name":           question.Name,
			"options":        question.Options,
			"correct_option": question.Correct_Option,
			"quiz":           question.Quiz,
			"points":         question.Points,
		})
	}
}

func (app *App) PostQuestionDetails(c *gin.Context) {
	var (
		question Question
		id       int64
		quiz     Quiz
	)
	db := app.Db

	if err := c.ShouldBindJSON(&question); err == nil {
		if question.Name != "" && question.Options != "" &&
			question.Correct_Option != 0 && question.Quiz != 0 && question.Points != 0 {
			log.Println("INFO: json binding is successful")
		} else {
			log.Println("Error in json binding")
			c.JSON(http.StatusBadRequest, gin.H{
				status:       failure,
				explaination: "all the fields are required i.e strings cannot be empty and integer cannot be zero",
			})
			panic(err.Error())
		}
	} else {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			status:       failure,
			explaination: fmt.Sprintf("Invalid Input: %s", err.Error()),
		})
		panic(err)
	}

	err := db.QueryRow(`Select name,description From quiz WHERE id=?`, question.Quiz).Scan(&quiz.Name, &quiz.Description)
	if err != nil {
		log.Println(err.Error())
		c.JSON(400, gin.H{
			status:       failure,
			explaination: fmt.Sprintf("Quiz Not found: %s", err.Error()),
		})
		panic(err)
	}

	res, err := db.Exec(`INSERT INTO questions(name, options, correct_option, quiz, points) VALUES(?, ?, ?, ?, ?)`, question.Name, question.Options, question.Correct_Option, question.Quiz, question.Points)
	if err != nil {
		log.Println(err.Error())
		c.JSON(400, gin.H{
			status:       failure,
			explaination: err.Error(),
		})
	}
	id, err = res.LastInsertId()

	if err == nil {
		c.JSON(http.StatusCreated, gin.H{
			"id":             id,
			"name":           question.Name,
			"options":        question.Options,
			"correct_option": question.Correct_Option,
			"quiz":           question.Quiz,
			"points":         question.Points,
		})
	} else {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			status:       failure,
			explaination: fmt.Sprintf("Question Insert Error: %s", err.Error()),
		})
		panic(err)
	}
}

func (app *App) GetAllQuestions(c *gin.Context) {
	quiz_id, err := strconv.ParseInt(c.Param("quiz_id"), 10, 64)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			status: failure,
			reason: err.Error(),
		})
		panic(err.Error())
	}
	db := app.Db
	var (
		question Question
		quiz     Quiz
	)

	err = db.QueryRow(`Select name, description From quiz WHERE id=?`, quiz_id).Scan(&quiz.Name, &quiz.Description)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			status: failure,
			reason: "no rows in result set",
		})
		panic(err)
	} else {
		log.Println("INFO: quiz details gathered")
	}
	rows, err := db.Query(`Select id, name, options, correct_option, points From questions WHERE quiz=?`, quiz_id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			status:       failure,
			explaination: err.Error(),
		})
		panic(err)
	}

	question.Quiz = quiz_id

	type questiondetailmap map[string]interface{}

	var questionsslice []questiondetailmap

	for rows.Next() {

		err := rows.Scan(&question.ID, &question.Name, &question.Options, &question.Correct_Option, &question.Points)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			questionsslice = append(questionsslice, map[string]interface{}{
				"id":             question.ID,
				"name":           question.Name,
				"options":        question.Options,
				"correct_option": question.Correct_Option,
				"quiz":           question.Quiz,
				"points":         question.Points,
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"name":        quiz.Name,
		"description": quiz.Description,
		"questions":   questionsslice,
	})
}
