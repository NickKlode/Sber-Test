package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"sber-test"
	"sber-test/internal/service"
	service_mocks "sber-test/internal/service/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_createList(t *testing.T) {

	type mockBehavior func(s *service_mocks.MockTodoList, list sber.TodoList)
	tests := []struct {
		name                 string
		mock                 mockBehavior
		listBody             string
		list                 sber.TodoList
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "OK",
			listBody: `{"id" : 1, "title":"title", "description":"description", "date":"2020-20-20", "done":false}`,
			mock: func(s *service_mocks.MockTodoList, list sber.TodoList) {
				s.EXPECT().Create(list).Return(1, nil)
			},
			list: sber.TodoList{
				Id:          1,
				Title:       "title",
				Description: "description",
				Date:        "2020-20-20",
				Done:        false,
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Wrong Input",
			listBody:             `{"description":"description"}`,
			list:                 sber.TodoList{},
			mock:                 func(s *service_mocks.MockTodoList, list sber.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:     "Server Error",
			listBody: `{"id" : 1, "title":"title", "description":"description", "date":"2020-20-20", "done":false}`,
			list: sber.TodoList{
				Id:          1,
				Title:       "title",
				Description: "description",
				Date:        "2020-20-20",
				Done:        false,
			},
			mock: func(s *service_mocks.MockTodoList, list sber.TodoList) {
				s.EXPECT().Create(list).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockTodoList(c)
			tt.mock(repo, tt.list)

			services := &service.Service{TodoList: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/lists", handler.createList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/lists", bytes.NewBufferString(tt.listBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestHandler_getAll(t *testing.T) {
	type mockBehavior func(s *service_mocks.MockTodoList, list []sber.TodoList)
	tests := []struct {
		name                 string
		mock                 mockBehavior
		list                 []sber.TodoList
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mock: func(s *service_mocks.MockTodoList, list []sber.TodoList) {
				s.EXPECT().GetAll().Return(list, nil)
			},
			list: []sber.TodoList{{
				Id:          1,
				Title:       "title1",
				Description: "description1",
				Date:        "2020-20-20",
				Done:        false,
			},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":1,"title":"title1","description":"description1","date":"2020-20-20","done":false}]`,
		},
		{
			name: "Empty List",
			mock: func(s *service_mocks.MockTodoList, list []sber.TodoList) {
				s.EXPECT().GetAll().Return(list, nil)
			},
			list:                 []sber.TodoList{},
			expectedStatusCode:   200,
			expectedResponseBody: `[]`,
		},
		{
			name: "Server Error",
			mock: func(s *service_mocks.MockTodoList, list []sber.TodoList) {
				s.EXPECT().GetAll().Return(nil, errors.New("something went wrong"))
			},
			list: []sber.TodoList{{
				Id:          1,
				Title:       "title1",
				Description: "description1",
				Date:        "2020-20-20",
				Done:        false,
			},
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockTodoList(c)
			tt.mock(repo, tt.list)

			services := &service.Service{TodoList: repo}
			handler := Handler{services}

			r := gin.New()
			r.GET("/lists", handler.getAll)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/lists", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestHandler_getByDate(t *testing.T) {
	type mockBehavior func(s *service_mocks.MockTodoList, date sber.FindInput, data []sber.TodoList)
	tests := []struct {
		name                 string
		mock                 mockBehavior
		inputBody            string
		input                sber.FindInput
		data                 []sber.TodoList
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mock: func(s *service_mocks.MockTodoList, date sber.FindInput, data []sber.TodoList) {
				s.EXPECT().GetByDate(date.Date).Return(data, nil)
			},
			inputBody: `{"date":"2020-20-20"}`,
			input: sber.FindInput{
				Date: "2020-20-20",
			},
			data: []sber.TodoList{
				{
					Id:          1,
					Title:       "title1",
					Description: "description1",
					Date:        "2020-20-20",
					Done:        false,
				},
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":1,"title":"title1","description":"description1","date":"2020-20-20","done":false}]`,
		},
		{
			name:      "No Input",
			inputBody: "",
			input: sber.FindInput{
				Date: "",
			},
			data: []sber.TodoList{},

			mock: func(s *service_mocks.MockTodoList, date sber.FindInput, data []sber.TodoList) {},

			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name: "Server Error",
			mock: func(s *service_mocks.MockTodoList, date sber.FindInput, data []sber.TodoList) {
				s.EXPECT().GetByDate(date.Date).Return(nil, errors.New("something went wrong"))
			},
			inputBody: `{"date":"2020-20-20"}`,
			input: sber.FindInput{
				Date: "2020-20-20",
			},
			data: []sber.TodoList{
				{
					Id:          1,
					Title:       "title1",
					Description: "description1",
					Date:        "2020-20-20",
					Done:        false,
				},
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockTodoList(c)
			tt.mock(repo, tt.input, tt.data)

			services := &service.Service{TodoList: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/lists/find", handler.getByDate)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/lists/find", bytes.NewBufferString(tt.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}
