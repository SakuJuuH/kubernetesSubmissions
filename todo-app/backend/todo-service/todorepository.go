package main

import "github.com/jmoiron/sqlx"

type TodoRepository interface {
	GetTodos() ([]Todo, error)
	AddTodo(task string) (Todo, error)
	healthCheck() (bool, error)
}

type todoRepository struct {
	db *sqlx.DB
}

func (t todoRepository) healthCheck() (bool, error) {
	err := t.db.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewTodoRepository(db *sqlx.DB) TodoRepository {
	return &todoRepository{db}
}

func (t todoRepository) GetTodos() ([]Todo, error) {
	todos := make([]Todo, 0)
	err := t.db.Select(&todos, "SELECT id, task, done FROM todos")
	return todos, err
}

func (t todoRepository) AddTodo(task string) (Todo, error) {
	var todo Todo
	err := t.db.Get(&todo, "INSERT INTO todos (task) VALUES ($1) RETURNING id, task, done", task)
	return todo, err
}
