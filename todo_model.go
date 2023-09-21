package sber

import (
	"errors"
)

type TodoList struct {
	Id          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Done        bool   `json:"done"`
}

type UpdateInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Date        *string `json:"date"`
	Done        *bool   `json:"done"`
}

type FindInput struct {
	Date string `json:"date"`
}

func (i UpdateInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Date == nil && i.Done == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
