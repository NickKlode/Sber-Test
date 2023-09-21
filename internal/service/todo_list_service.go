package service

import (
	"sber-test"
	"sber-test/internal/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(list sber.TodoList) (int, error) {
	return s.repo.Create(list)
}

func (s *TodoListService) GetAll() ([]sber.TodoList, error) {
	return s.repo.GetAll()
}

func (s *TodoListService) Delete(listId int) error {
	return s.repo.Delete(listId)
}

func (s *TodoListService) Update(listId int, input sber.UpdateInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(listId, input)
}
func (s *TodoListService) GetByDate(date string) ([]sber.TodoList, error) {
	return s.repo.GetByDate(date)
}
