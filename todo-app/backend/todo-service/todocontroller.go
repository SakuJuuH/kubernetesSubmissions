package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type TodosController struct {
	repo TodoRepository
}

func NewTodosController(repo TodoRepository) *TodosController {
	return &TodosController{repo: repo}
}

func (c *TodosController) getTodos(ctx *gin.Context) {
	todos, err := c.repo.GetTodos()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get todos")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info().
		Str("path", ctx.FullPath()).
		Int("count", len(todos)).
		Msg("Todos received")
	ctx.JSON(http.StatusOK, todos)
}

func (c *TodosController) createTodo(ctx *gin.Context) {
	var requestTodo struct {
		Task string `json:"task" binding:"required"`
	}

	if err := ctx.BindJSON(&requestTodo); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, ok := validateAndLogTask(ctx, requestTodo.Task)
	if !ok {
		return
	}

	newTodo, err := c.repo.AddTodo(task)
	if err != nil {
		log.Error().Err(err).Msg("todo insert failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, newTodo)
}

func (c *TodosController) welcome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Welcome to the Todo API! Use /api/todos to manage your tasks.",
		"status_code": http.StatusOK,
	})
}

func (c *TodosController) createRandomTodo(ctx *gin.Context) {
	var randomArticleURL = os.Getenv("RANDOM_ARTICLE_URL")

	if randomArticleURL == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "RANDOM_ARTICLE_URL is not set"})
		return
	}

	resp, err := http.Get(randomArticleURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get random article")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	articleURL := resp.Request.URL.String()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}(resp.Body)

	task := fmt.Sprintf("Read: %s", articleURL)

	task, ok := validateAndLogTask(ctx, task)
	if !ok {
		return
	}

	createdTodo, err := c.repo.AddTodo(task)
	if err != nil {
		log.Error().Err(err).Msg("random todo insert failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"New todo created": createdTodo,
	})
}

func validateAndLogTask(ctx *gin.Context, task string) (string, bool) {
	const maxLen = 140
	length := len(task)

	if task == "" {
		log.Warn().
			Str("path", ctx.FullPath()).
			Int("length", length).
			Msg("todo rejected: empty task")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task cannot be empty"})
		return "", false
	}

	if length > maxLen {
		log.Warn().
			Str("path", ctx.FullPath()).
			Int("length", length).
			Int("max_length", maxLen).
			Str("task", task).
			Msg("todo rejected: too long")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task cannot exceed 140 characters"})
		return "", false
	}

	log.Info().
		Str("path", ctx.FullPath()).
		Int("length", length).
		Str("task", task).
		Msg("todo received")

	return task, true
}
