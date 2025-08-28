package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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
	randomArticleURL := os.Getenv("RANDOM_ARTICLE_URL")
	if randomArticleURL == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "RANDOM_ARTICLE_URL is not set"})
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", randomArticleURL, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get random article")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	redirectedURL := resp.Request.URL.String()

	if redirectedURL == randomArticleURL {
		log.Warn().
			Str("path", ctx.FullPath()).
			Str("url", redirectedURL).
			Msg("random article rejected: same as RANDOM_ARTICLE_URL")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read the same article twice"})
		return
	}

	task := fmt.Sprintf("Read: %s", redirectedURL)
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

func (c *TodosController) healthCheck(ctx *gin.Context) {
	health, err := c.repo.healthCheck()
	if err != nil {
		log.Error().Err(err).Msg("Database health check failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database is not reachable"})
		return
	}
	if !health {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database is not reachable"})
	}
	msg := "Database is reachable"
	log.Info().Msg(msg)
	ctx.JSON(http.StatusOK, gin.H{"message": msg})
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
