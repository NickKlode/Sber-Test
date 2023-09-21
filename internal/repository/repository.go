package repository

import (
	"sber-test"

	"github.com/jmoiron/sqlx"
)

type TodoList interface {
	Create(list sber.TodoList) (int, error)
	GetAll() ([]sber.TodoList, error)
	Delete(listId int) error
	Update(listId int, input sber.UpdateInput) error
	GetByDate(date string) ([]sber.TodoList, error)
}

type Repository struct {
	TodoList
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TodoList: NewTodoListPostgres(db),
	}
}
