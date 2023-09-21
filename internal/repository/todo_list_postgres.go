package repository

import (
	"fmt"
	"sber-test"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (p *TodoListPostgres) Create(list sber.TodoList) (int, error) {
	query := "INSERT INTO todolist (title, description, date) VALUES ($1, $2, $3) RETURNING id"
	var id int
	row := p.db.QueryRow(query, list.Title, list.Description, list.Date)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil

}

func (p *TodoListPostgres) GetAll() ([]sber.TodoList, error) {
	var list []sber.TodoList

	query := "SELECT * FROM todolist"
	err := p.db.Select(&list, query)
	return list, err
}

func (p *TodoListPostgres) Delete(listId int) error {
	query := "DELETE FROM todolist WHERE id=$1"
	_, err := p.db.Exec(query, listId)
	return err
}

func (p *TodoListPostgres) Update(listId int, input sber.UpdateInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	if input.Date != nil {
		setValues = append(setValues, fmt.Sprintf("date=$%d", argId))
		args = append(args, *input.Date)
		argId++
	}
	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE todolist SET %s WHERE id=$%d", setQuery, argId)
	args = append(args, listId)
	_, err := p.db.Exec(query, args...)
	return err
}

func (p *TodoListPostgres) GetByDate(date string) ([]sber.TodoList, error) {
	var list []sber.TodoList

	query := "SELECT * FROM todolist WHERE date=$1"
	err := p.db.Select(&list, query, date)
	return list, err
}
