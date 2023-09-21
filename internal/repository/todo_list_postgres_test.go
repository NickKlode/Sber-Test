package repository

import (
	"database/sql"
	"sber-test"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTodoListPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("error while opening db connection. %s", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		list sber.TodoList
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "OK",
			input: args{
				list: sber.TodoList{
					Title:       "title",
					Description: "description",
					Date:        "2020-20-20",
				},
			},
			want: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO todolist").WithArgs("title", "description", "2020-20-20").WillReturnRows(rows)
			},
		},
		{
			name: "Empty Fields",
			input: args{
				list: sber.TodoList{
					Title:       "",
					Description: "description",
					Date:        "2020-20-20",
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO todolist").WithArgs("", "description", "2020-20-20").WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Create(tt.input.list)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("error while opening db connection. %s", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		want    []sber.TodoList
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "done"}).
					AddRow(1, "title1", "description1", "2020-20-20", false).
					AddRow(2, "title2", "description2", "2020-20-20", false)
				mock.ExpectQuery("SELECT (.+) FROM todolist").WillReturnRows(rows)
			},
			want: []sber.TodoList{
				{
					Id:          1,
					Title:       "title1",
					Description: "description1",
					Date:        "2020-20-20",
					Done:        false,
				},
				{
					Id:          2,
					Title:       "title2",
					Description: "description2",
					Date:        "2020-20-20",
					Done:        false,
				},
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "done"})
				mock.ExpectQuery("SELECT (.+) FROM todolist").WillReturnRows(rows)
			},
			want: []sber.TodoList(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAll()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}

func TestTodoListPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("error while opening db connection. %s", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		listId int
		input  sber.UpdateInput
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectExec("UPDATE todolist SET(.+) WHERE (.+)").
					WithArgs("new title", "new description", true, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				input: sber.UpdateInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
					Done:        boolPointer(true),
				},
			},
		},
		{
			name: "No Input",
			mock: func() {
				mock.ExpectExec("UPDATE todolist SET WHERE (.+)").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Update(tt.input.listId, tt.input.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("error while opening db connection. %s", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		listId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectExec("DELETE FROM todolist WHERE (.+)").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec("DELETE FROM todolist WHERE (.+)").WithArgs(100).WillReturnError(sql.ErrNoRows)
			},
			input: args{
				listId: 100,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Delete(tt.input.listId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_GetByDate(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("error while opening db connection. %s", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		date sber.FindInput
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []sber.TodoList
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {

				rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "done"}).
					AddRow(1, "title1", "description1", "2020-20-20", false).
					AddRow(2, "title2", "description2", "2020-20-20", false)
				mock.ExpectQuery("SELECT (.+) FROM todolist WHERE (.+)").WithArgs("2020-20-20").WillReturnRows(rows)
			},
			input: args{
				date: sber.FindInput{
					Date: "2020-20-20",
				},
			},
			want: []sber.TodoList{
				{
					Id:          1,
					Title:       "title1",
					Description: "description1",
					Date:        "2020-20-20",
					Done:        false,
				},
				{
					Id:          2,
					Title:       "title2",
					Description: "description2",
					Date:        "2020-20-20",
					Done:        false,
				},
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "done"})
				mock.ExpectQuery("SELECT (.+) FROM todolist WHERE (.+)").WithArgs("2000-20-20").WillReturnRows(rows)
			},
			input: args{
				date: sber.FindInput{
					Date: "2000-20-20",
				},
			},
			want: []sber.TodoList(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetByDate(tt.input.date.Date)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
