package sber

import (
	"errors"
	"regexp"
	"time"
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
func ValidateDate(date string) error {
	if date == "" {
		return errors.New("date field is empty")
	}
	re := regexp.MustCompile("^202[3-9]{1}-[0-1][0-2]-[0-3][0-9]$")

	d := re.MatchString(date)
	if d {
		return nil
	}

	return errors.New("wrong date")
}

func SetDate() string {
	layout := "2006-01-02"
	t := time.Now()
	dateString := t.Format(layout)
	date := dateString
	return date
}
