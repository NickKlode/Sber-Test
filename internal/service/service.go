package service

import (
	"sber-test"
	"sber-test/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type TodoList interface {
	Create(list sber.TodoList) (int, error)
	GetAll() ([]sber.TodoList, error)
	Delete(listId int) error
	Update(listId int, input sber.UpdateInput) error
	GetByDate(date string) ([]sber.TodoList, error)
}

type Service struct {
	TodoList
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		TodoList: NewTodoListService(repo.TodoList),
	}
}
